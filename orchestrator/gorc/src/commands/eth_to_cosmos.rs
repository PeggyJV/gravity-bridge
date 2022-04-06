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
#[clap(
    long_about = "DESCRIPTION \n\n Send Eth token to Cosmos chain.\n This command sends Eth token to the Cosmos chain via the Gravity bridge. \n It takes the tx amount, Eth keyname, contract address, Cosmos token destination, number of times \n and the ERC20 token contract address."
)]
pub struct EthToCosmosCmd {
    /// Erc20 contract address.
    #[clap(short = 'E', long)]
    erc20_address: String,

    /// Tx amount.
    #[clap(short, long)]
    init_amount: String,

    /// Eth keyname.
    #[clap(short, long)]
    ethereum_key: String,

    /// Cosmos address
    #[clap(short, long)]
    cosmos_dest: String,

    /// The number of times transactions should repeat itself, default is 1.
    #[clap(short, long, default_value = "1")]
    times: String,

    /// Contract address.
    #[clap(short = 'C', long)]
    contract_address: String,
}

impl Runnable for EthToCosmosCmd {
    fn run(&self) {
        let config = APP.config();
        let erc20_address = self.erc20_address.clone();
        let erc20_address: EthAddress = erc20_address
            .parse()
            .expect("Invalid ERC20 contract address!");

        let ethereum_wallet = config.load_ethers_wallet(self.ethereum_key.clone());

        let contract_address: EthAddress = self
            .contract_address
            .parse()
            .expect("Invalid contract address!");

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
            let cosmos_dest: CosmosAddress = self.cosmos_dest.parse().unwrap();
            let ethereum_address = eth_client.address();
            check_for_eth(ethereum_address, eth_client.clone()).await;
            let amount =
                U256::from_dec_str(&self.init_amount).expect("cannot parse amount as U256");

            let erc20_balance =
                get_erc20_balance(erc20_address, ethereum_address, eth_client.clone())
                    .await
                    .expect("Failed to get balance, check ERC20 contract address");

            let times = self.times.clone();
            let times_usize = times.parse::<usize>().expect("cannot parse times as usize");
            let times_u256 = U256::from_dec_str(&times).expect("cannot parse times as U256");

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
                    self.init_amount.parse::<f64>().unwrap(),
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
