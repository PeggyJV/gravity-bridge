use std::sync::Arc;

use crate::main_loop::relayer_main_loop;
use crate::main_loop::LOOP_SPEED;
use docopt::Docopt;
use env_logger::Env;
use ethers::prelude::*;
use ethers::signers::LocalWallet as EthWallet;
use ethers::types::Address as EthAddress;
use gravity_utils::{
    connection_prep::{check_for_eth, create_rpc_connections, wait_for_cosmos_node_ready},
    ethereum::{downcast_to_u64, format_eth_address},
};

pub mod batch_relaying;
pub mod find_latest_valset;
pub mod logic_call_relaying;
pub mod main_loop;
pub mod valset_relaying;

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate lazy_static;
#[macro_use]
extern crate log;

#[derive(Debug, Deserialize)]
struct Args {
    flag_ethereum_key: String,
    flag_cosmos_grpc: String,
    flag_address_prefix: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage: {} --ethereum-key=<key> --cosmos-grpc=<url> --address-prefix=<prefix> --ethereum-rpc=<url> --contract-address=<addr>
        Options:
            -h --help                    Show this screen.
            --ethereum-key=<ekey>        An Ethereum private key containing non-trivial funds
            --cosmos-grpc=<gurl>         The Cosmos gRPC url
            --address-prefix=<prefix>    The prefix for addresses on this Cosmos chain
            --ethereum-grpc=<eurl>       The Ethereum RPC url, Geth light clients work and sync fast
            --contract-address=<addr>    The Ethereum contract address for Gravity
        About:
            The Gravity relayer component, responsible for relaying data from the Cosmos blockchain
            to the Ethereum blockchain, cosmos key and fees are optional since they are only used
            to request the creation of batches or validator sets to relay.
            for Althea-Gravity.
            Written By: {}
            Version {}",
            env!("CARGO_PKG_NAME"),
            env!("CARGO_PKG_AUTHORS"),
            env!("CARGO_PKG_VERSION"),
        );
}

#[actix_rt::main]
async fn main() {
    env_logger::Builder::from_env(Env::default().default_filter_or("info")).init();
    // On Linux static builds we need to probe ssl certs path to be able to
    // do TLS stuff.
    openssl_probe::init_ssl_cert_env_vars();

    let args: Args = Docopt::new(USAGE.as_str())
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());
    let ethereum_wallet: EthWallet = args
        .flag_ethereum_key
        .parse()
        .expect("Invalid Ethereum private key!");
    let gravity_contract_address: EthAddress = args
        .flag_contract_address
        .parse()
        .expect("Invalid contract address!");

    let connections = create_rpc_connections(
        args.flag_address_prefix,
        Some(args.flag_cosmos_grpc),
        Some(args.flag_ethereum_rpc),
        LOOP_SPEED,
    )
    .await;
    let provider = connections.eth_provider.clone().unwrap();
    let chain_id = provider
        .get_chainid()
        .await
        .expect("Could not retrieve chain ID during relayer start");
    let chain_id = downcast_to_u64(chain_id).expect("Chain ID overflowed when downcasting to u64");
    let eth_client =
        SignerMiddleware::new(provider, ethereum_wallet.clone().with_chain_id(chain_id));
    let eth_client = Arc::new(eth_client);

    let public_eth_key = eth_client.address();
    info!("Starting Gravity Relayer");
    info!("Ethereum Address: {}", format_eth_address(public_eth_key));

    let contact = connections.contact.clone().unwrap();

    // check if the cosmos node is syncing, if so wait for it
    // we can't move any steps above this because they may fail on an incorrect
    // historic chain state while syncing occurs
    wait_for_cosmos_node_ready(&contact).await;
    check_for_eth(public_eth_key, eth_client.clone()).await;

    relayer_main_loop(
        eth_client,
        connections.grpc.unwrap(),
        gravity_contract_address,
        1f32,
    )
    .await
}
