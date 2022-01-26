//! `eth subcommands` subcommand

use crate::{application::APP, prelude::*, utils::*};
use abscissa_core::{clap::Parser, Command, Runnable};
use deep_space::address::Address as CosmosAddress;
use ethereum_gravity::erc20_utils::get_erc20_balance;
use ethereum_gravity::send_to_cosmos::send_to_cosmos;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_utils::{
    connection_prep::{check_for_eth, create_rpc_connections},
    ethereum::{downcast_to_u64, format_eth_address},
};
use std::sync::Arc;

/// Create transactions in Eth chain
#[derive(Command, Debug, Parser)]
pub enum Eth {
    SendToCosmos(SendToCosmos),

    Send(Send),
}

impl Runnable for Eth {
    fn run(&self) {}
}

#[derive(Command, Debug, Parser)]
pub struct SendToCosmos {
    free: Vec<String>,

    #[clap(short, long)]
    help: bool,
}

// TODO(bolten): I guess this command is also non-functional?
// are the commands under tx dead code?
fn lookup_eth_key(_key: String) -> LocalWallet {
    todo!()
}

impl Runnable for SendToCosmos {
    fn run(&self) {
        assert!(self.free.len() == 4);
        let from_eth_key = self.free[0].clone();
        let to_cosmos_addr: CosmosAddress = self.free[1]
            .clone()
            .parse()
            .expect("Expected a valid Cosmos Address");
        let erc20_contract: EthAddress = self.free[2]
            .clone()
            .parse()
            .expect("Expected a valid Eth Address");
        let erc20_amount = self.free[3].clone();
        let ethereum_wallet = lookup_eth_key(from_eth_key);

        println!(
            "Sending from Eth address {}",
            format_eth_address(ethereum_wallet.address())
        );
        let config = APP.config();
        let cosmos_prefix = config.cosmos.prefix.clone();
        let cosmso_grpc = config.cosmos.grpc.clone();
        let eth_rpc = config.ethereum.rpc.clone();
        let contract_address: EthAddress = config
            .gravity
            .contract
            .clone()
            .parse()
            .expect("Expected config.gravity.contract to be an Eth ddress");

        abscissa_tokio::run_with_actix(&APP, async {
            let connections =
                create_rpc_connections(cosmos_prefix, Some(cosmso_grpc), Some(eth_rpc), TIMEOUT)
                    .await;
            let provider = connections.eth_provider.clone().unwrap();
            let chain_id = provider
                .get_chainid()
                .await
                .expect("Could not retrieve chain ID");
            let chain_id =
                downcast_to_u64(chain_id).expect("Chain ID overflowed when downcasting to u64");
            let eth_client =
                SignerMiddleware::new(provider, ethereum_wallet.clone().with_chain_id(chain_id));
            let eth_client = Arc::new(eth_client);
            check_for_eth(eth_client.address(), eth_client.clone()).await;

            let amount =
                U256::from_dec_str(erc20_amount.as_str()).expect("Could not parse amount to U256");

            let erc20_balance =
                get_erc20_balance(erc20_contract, eth_client.address(), eth_client.clone())
                    .await
                    .expect("Failed to get balance, check ERC20 contract address");

            if erc20_balance == 0u8.into() {
                panic!(
                    "You have zero {} tokens, please double check your sender and erc20 addresses!",
                    erc20_contract
                );
            }
            println!(
                "Sending {} / {} to Cosmos from {} to {}",
                amount,
                erc20_contract,
                eth_client.address(),
                to_cosmos_addr
            );
            // we send some erc20 tokens to the gravity contract to register a deposit
            let res = send_to_cosmos(
                erc20_contract,
                contract_address,
                amount,
                to_cosmos_addr,
                Some(TIMEOUT),
                eth_client.clone(),
            )
            .await;
            match res {
                Ok(tx_id) => println!("Send to Cosmos txid: {:#066x}", tx_id),
                Err(e) => println!("Failed to send tokens! {:?}", e),
            }
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}

#[derive(Command, Debug, Parser)]
pub struct Send {
    free: Vec<String>,

    #[clap(short, long)]
    help: bool,
}

impl Runnable for Send {
    fn run(&self) {}
}
