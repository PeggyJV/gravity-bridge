use crate::{application::APP, prelude::*};
use abscissa_core::{clap::Parser, Command, Runnable};
use ethereum_gravity::deploy_erc20::deploy_erc20;
use ethers::prelude::*;
use gravity_utils::{
    connection_prep::{check_for_eth, create_eth_provider},
    ethereum::{downcast_to_u64, format_eth_hash},
};
use ocular::GrpcClient;
use ocular_somm_gravity::SommGravityExt;
use std::convert::TryFrom;
use std::process::exit;
use std::{sync::Arc, time::Duration};
use tokio::time::sleep as delay_for;

/// Deploy Erc20
#[derive(Command, Debug, Parser)]
pub struct Erc20 {
    args: Vec<String>,

    #[clap(short, long)]
    ethereum_key: String,

    #[clap(short, long, default_value = "1.0")]
    gas_multiplier: f64,
}

impl Runnable for Erc20 {
    fn run(&self) {
        abscissa_tokio::run_with_actix(&APP, async {
            self.deploy().await;
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            exit(1);
        });
    }
}

impl Erc20 {
    async fn deploy(&self) {
        let denom = self.args.get(0).expect("denom is required");

        let config = APP.config();

        let ethereum_wallet = config.load_ethers_wallet(self.ethereum_key.clone());
        let contract_address = config
            .gravity
            .contract
            .parse()
            .expect("Could not parse gravity contract address");

        let timeout = Duration::from_secs(500);
        let cosmos_client = GrpcClient::new(&config.cosmos.grpc)
            .await
            .expect("failed to construct GrpcClient");
        let provider = create_eth_provider(config.ethereum.rpc.clone())
            .await
            .expect("error creating eth provider");
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

        let res = cosmos_client
            .query_denom_to_erc20_params(denom)
            .await
            .expect("Couldn't get erc-20 params");

        println!("Starting deploy of ERC20");

        let res = deploy_erc20(
            res.base_denom,
            res.erc20_name,
            res.erc20_symbol,
            u8::try_from(res.erc20_decimals).unwrap(),
            contract_address,
            Some(timeout),
            self.gas_multiplier,
            eth_client.clone(),
        )
        .await
        .expect("Could not deploy ERC20");

        println!("We have deployed ERC20 contract at tx hash {}, waiting to see if the Cosmos chain choses to adopt it",
            format_eth_hash(res));

        match tokio::time::timeout(Duration::from_secs(300), async {
            loop {
                let res = cosmos_client.query_denom_to_erc20(denom).await;

                if let Ok(erc20) = res {
                    break erc20;
                }
                delay_for(Duration::from_secs(1)).await;
            }
        })
        .await
        {
            Ok(erc20) => {
                println!(
                    "Asset {} has accepted new ERC20 representation {}",
                    denom, erc20
                );
                exit(0);
            }
            Err(_) => {
                println!(
                    "Your ERC20 contract was not adopted, double check the metadata and try again"
                );
                exit(1);
            }
        }
    }
}
