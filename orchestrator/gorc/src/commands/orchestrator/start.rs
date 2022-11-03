use crate::{application::APP, prelude::*};
use abscissa_core::{clap::Parser, Command, Runnable};
use ethers::{prelude::*, types::Address as EthAddress};
use gravity_utils::{
    connection_prep::{
        check_delegate_addresses, check_for_eth, check_for_fee_denom, create_eth_provider,
        wait_for_cosmos_node_ready,
    },
    ethereum::{downcast_to_u64, format_eth_address},
};
use ocular::{chain::ChainContext, GrpcClient};
use orchestrator::main_loop::orchestrator_main_loop;
use std::sync::Arc;

/// Start the Orchestrator
#[derive(Command, Debug, Parser)]
pub struct StartCommand {
    #[clap(short, long)]
    cosmos_key: String,

    #[clap(short, long)]
    ethereum_key: String,

    #[clap(short, long)]
    orchestrator_only: bool,
}

impl Runnable for StartCommand {
    fn run(&self) {
        openssl_probe::init_ssl_cert_env_vars();

        let config = APP.config();
        let cosmos_prefix = config.cosmos.prefix.clone();

        let cosmos_account = config.load_account(self.cosmos_key.clone());
        let cosmos_address = cosmos_account.id(&cosmos_prefix).unwrap();

        let context = ChainContext {
            prefix: cosmos_prefix,
            ..Default::default()
        };

        let ethereum_wallet = config.load_ethers_wallet(self.ethereum_key.clone());
        let ethereum_address = ethereum_wallet.address();

        let contract_address: EthAddress = config
            .gravity
            .contract
            .parse()
            .expect("Could not parse gravity contract address");

        let fees_denom = config.gravity.fees_denom.clone();

        abscissa_tokio::run_with_actix(&APP, async {
            let mut grpc_client = GrpcClient::new(&config.cosmos.grpc)
                .await
                .expect("failed to construct GrpcClient");
            let provider = create_eth_provider(config.ethereum.rpc.clone())
                .await
                .expect("error creating eth provider");
            let chain_id = provider
                .get_chainid()
                .await
                .expect("Could not retrieve chain ID during orchestrator start");
            let chain_id =
                downcast_to_u64(chain_id).expect("Chain ID overflowed when downcasting to u64");
            let eth_client =
                SignerMiddleware::new(provider, ethereum_wallet.clone().with_chain_id(chain_id));
            let eth_client = Arc::new(eth_client);

            info!("Starting Relayer + Oracle + Ethereum Signer");
            info!("Ethereum Address: {}", format_eth_address(ethereum_address));
            info!("Cosmos Address {}", cosmos_address);

            // check if the cosmos node is syncing, if so wait for it
            // we can't move any steps above this because they may fail on an incorrect
            // historic chain state while syncing occurs
            wait_for_cosmos_node_ready(&mut grpc_client).await;

            // check if the delegate addresses are correctly configured
            check_delegate_addresses(&mut grpc_client, ethereum_address, &cosmos_address).await;

            // check if we actually have the promised balance of tokens to pay fees
            check_for_fee_denom(&fees_denom, &cosmos_address, &mut grpc_client)
                .await
                .expect("failed to check for fee denom");
            check_for_eth(ethereum_address, eth_client.clone()).await;

            let gas_price = config.cosmos.gas_price.as_tuple();

            orchestrator_main_loop(
                &config.cosmos.grpc,
                &cosmos_account,
                context,
                eth_client,
                contract_address,
                gas_price,
                &config.metrics.listen_addr,
                config.ethereum.gas_price_multiplier,
                config.ethereum.gas_multiplier,
                config.ethereum.blocks_to_search,
                config.cosmos.gas_adjustment,
                self.orchestrator_only,
                config.cosmos.msg_batch_size,
            )
            .await;
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}
