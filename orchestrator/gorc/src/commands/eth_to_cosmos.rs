use crate::application::APP;
use abscissa_core::{clap::Parser, status_err, Application, Command, Runnable};
use deep_space::address::Address as CosmosAddress;
use ethereum_gravity::erc20_utils::get_erc20_balance;
use ethereum_gravity::send_to_cosmos::send_to_cosmos;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_utils::{
    connection_prep::{check_for_eth, create_rpc_connections},
    ethereum::downcast_to_u64,
};
use std::{sync::Arc, time::Duration};

const TIMEOUT: Duration = Duration::from_secs(60);

/// This command send Ethereum to Cosmos
#[derive(Command, Debug, Default, Parser)]
pub struct EthToCosmosCmd {
    pub args: Vec<String>,
}

impl Runnable for EthToCosmosCmd {
    fn run(&self) {
        let config = APP.config();
        let erc20_address = self.args.get(0).expect("erc20 address is required");
        let erc20_address: EthAddress = erc20_address
            .parse()
            .expect("Invalid ERC20 contract address!");

        let ethereum_key = self.args.get(1).expect("key is required");
        let ethereum_wallet = config.load_ethers_wallet(ethereum_key.clone());

        let contract_address = self.args.get(2).expect("contract address is required");
        let contract_address: EthAddress =
            contract_address.parse().expect("Invalid contract address!");

        let cosmos_prefix = config.cosmos.prefix.trim();
        let eth_rpc = config.ethereum.rpc.trim();

        abscissa_tokio::run_with_actix(&APP, async {
            let connections = create_rpc_connections(
                cosmos_prefix.to_string(),
                None,
                Some(eth_rpc.to_string()),
                TIMEOUT,
            )
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
            let cosmos_dest = self.args.get(3).expect("cosmos destination is required");
            let cosmos_dest: CosmosAddress = cosmos_dest.parse().unwrap();
            let ethereum_address = eth_client.address();
            check_for_eth(ethereum_address, eth_client.clone()).await;

            let init_amount = self.args.get(4).expect("amount is required");
            let amount = U256::from_dec_str(init_amount).expect("cannot parse amount as U256");

            let erc20_balance =
                get_erc20_balance(erc20_address, ethereum_address, eth_client.clone())
                    .await
                    .expect("Failed to get balance, check ERC20 contract address");

            let times = self.args.get(5).expect("times is required");
            let times_usize = times.parse::<usize>().expect("cannot parse times as usize");
            let times_u256 = U256::from_dec_str(times).expect("cannot parse times as U256");

            if erc20_balance == 0u8.into() {
                panic!(
                    "You have zero {} tokens, please double check your sender and erc20 addresses!",
                    contract_address
                );
            } else if amount * times_u256 > erc20_balance {
                panic!(
                    "Insufficient balance {} > {}",
                    amount * times_u256,
                    erc20_balance
                );
            }

            for _ in 0..times_usize {
                println!(
                    "Sending {} / {} to Cosmos from {} to {}",
                    init_amount.parse::<f64>().unwrap(),
                    erc20_address,
                    ethereum_address,
                    cosmos_dest
                );
                // we send some erc20 tokens to the gravity contract to register a deposit
                let res = send_to_cosmos(
                    erc20_address,
                    contract_address,
                    amount,
                    cosmos_dest,
                    Some(TIMEOUT),
                    eth_client.clone(),
                )
                .await;
                match res {
                    Ok(tx_id) => println!("Send to Cosmos txid: {}", tx_id),
                    Err(e) => println!("Failed to send tokens! {:?}", e),
                }
            }
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
        });
    }
}
