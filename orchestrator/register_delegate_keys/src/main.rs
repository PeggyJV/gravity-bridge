//! This file is a single use binary that will allow you to register your validator ethereum key

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate lazy_static;

use clarity::PrivateKey as EthPrivateKey;
use docopt::Docopt;
use ethers::core::k256::ecdsa::SigningKey;
use ethers::core::k256::elliptic_curve::generic_array::GenericArray;
use ethers::prelude::*;
use gravity::deep_space::mnemonic::Mnemonic;
use gravity::deep_space::{CosmosPrivateKey, PrivateKey};
use gravity::send::update_gravity_delegate_addresses;
use gravity::utils::connection_prep::check_for_fee_denom;
use gravity::utils::connection_prep::{create_rpc_connections, wait_for_cosmos_node_ready};
use log::error;
use rand::{thread_rng, Rng};
use std::time::Duration;

#[derive(Debug, Deserialize)]
struct Args {
    flag_validator_phrase: String,
    flag_cosmos_phrase: Option<String>,
    flag_ethereum_key: Option<String>,
    flag_address_prefix: String,
    flag_cosmos_grpc: String,
    flag_fees: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
        "Usage: {} --validator-phrase=<key> --address-prefix=<prefix> [--cosmos-phrase=<key>] [--ethereum-key=<key>] --cosmos-grpc=<url> --fees=<denom>
        Options:
            -h --help                 Show this screen.
            --validator-phrase=<vkey> The Cosmos private key of the validator. Must be saved when you generate your key
            --ethereum-key=<ekey>     (Optional) The Ethereum private key to register, will be generated if not provided
            --cosmos-phrase=<ckey>    (Optional) The phrase for the Cosmos key to register, will be generated if not provided.
            --address-prefix=<prefix> The prefix for Addresses on this chain (eg 'cosmos')
            --cosmos-grpc=<curl>      The Cosmos RPC url, usually the validator. This will need to be manually enabled
            --fees=<denom>            The Cosmos Denom in which to pay Cosmos chain fees
        About:
            Special purpose binary for bootstrapping Gravity chains. This will submit and optionally
            generate an Ethereum key that will be used to sign messages on behalf of your Validator
            on the Cosmos blockchain running the Gravity module. Be aware this Ethereum key must be kept
            safe as you can be slashed for losing it.
            Written By: {}
            Version {}",
        env!("CARGO_PKG_NAME"),
        env!("CARGO_PKG_AUTHORS"),
        env!("CARGO_PKG_VERSION"),
    );
}

const TIMEOUT: Duration = Duration::from_secs(60);

#[tokio::main]
async fn main() {
    env_logger::init();
    // On Linux static builds we need to probe ssl certs path to be able to
    // do TLS stuff.
    #[allow(deprecated)]
    openssl_probe::init_ssl_cert_env_vars();

    let args: Args = Docopt::new(USAGE.as_str())
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());

    let fee_denom = args.flag_fees;

    let connections = create_rpc_connections(
        args.flag_address_prefix,
        Some(args.flag_cosmos_grpc),
        None,
        TIMEOUT,
    )
    .await;
    let contact = connections.contact.unwrap();
    wait_for_cosmos_node_ready(&contact).await;

    let validator_key = CosmosPrivateKey::from_phrase(&args.flag_validator_phrase, "")
        .expect("Failed to parse validator key");
    let validator_addr = validator_key.to_address(&contact.get_prefix()).unwrap();
    check_for_fee_denom(&fee_denom, validator_addr, &contact).await;

    let cosmos_key = if let Some(cosmos_phrase) = args.flag_cosmos_phrase {
        CosmosPrivateKey::from_phrase(&cosmos_phrase, "").expect("Failed to parse cosmos key")
    } else {
        let new_phrase = Mnemonic::generate(24).unwrap();
        let key = CosmosPrivateKey::from_phrase(new_phrase.as_str(), "").unwrap();
        println!(
            "No Cosmos key provided, your generated key is\n {} -> {}",
            new_phrase.as_str(),
            key.to_address(&contact.get_prefix()).unwrap()
        );
        key
    };
    let ethereum_key = if let Some(key) = args.flag_ethereum_key {
        key.parse().expect("Invalid Ethereum Private key!")
    } else {
        let mut rng = thread_rng();
        let key: [u8; 32] = rng.gen();
        let key = EthPrivateKey::from_bytes(key).unwrap();
        println!(
            "No Ethereum key provided, your generated key is\n {} -> {}",
            key,
            key.to_address()
        );
        key
    };

    // TODO(bolten): left clarity in place for the above bit because it seems like
    // SigningKey/VerifyingKey don't implement the Display trait
    let ethereum_key_bytes = ethereum_key.to_bytes();
    let key_bytes = GenericArray::from_slice(&ethereum_key_bytes);
    let signing_key = SigningKey::from_bytes(&key_bytes).unwrap();
    let ethereum_wallet = LocalWallet::from(signing_key);

    let ethereum_address = ethereum_wallet.address();
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

    let res = update_gravity_delegate_addresses(
        &contact,
        ethereum_address,
        cosmos_address,
        validator_key,
        ethereum_wallet,
        (0f64, "".to_string()),
        1.0f64,
    )
    .await
    .expect("Failed to update Eth address");

    let res = contact.wait_for_tx(res.into(), TIMEOUT).await;

    if let Err(e) = res {
        error!("Failed trying to register delegate addresses error {:?}, correct the error and try again", e);
        std::process::exit(1);
    }

    let eth_address = ethereum_key.to_address();
    println!(
        "Registered Delegate Ethereum address {} and Cosmos address {}",
        eth_address, cosmos_address
    )
}
