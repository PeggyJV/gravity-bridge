use crate::{application::APP, prelude::*};
use abscissa_core::{clap::Parser, Command, Runnable};
use ethers::{prelude::*, types::Address as EthAddress};
use gravity_utils::{
    connection_prep::{
        check_delegate_addresses, check_for_eth, check_for_fee_denom, create_rpc_connections,
        wait_for_cosmos_node_ready,
    },
    ethereum::{downcast_to_u64, format_eth_address},
};
use orchestrator::main_loop::{
    orchestrator_main_loop, ETH_ORACLE_LOOP_SPEED, ETH_SIGNER_LOOP_SPEED,
};
use relayer::main_loop::LOOP_SPEED as RELAYER_LOOP_SPEED;
use std::{cmp::min, sync::Arc};

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

        let cosmos_key = config.load_deep_space_key(self.cosmos_key.clone());
        let cosmos_address = cosmos_key.to_address(&cosmos_prefix).unwrap();

        let ethereum_wallet = config.load_ethers_wallet(self.ethereum_key.clone());
        let ethereum_address = ethereum_wallet.address();

        let contract_address: EthAddress = config
            .gravity
            .contract
            .parse()
            .expect("Could not parse gravity contract address");

        let fees_denom = config.gravity.fees_denom.clone();

        let timeout = min(
            min(ETH_SIGNER_LOOP_SPEED, ETH_ORACLE_LOOP_SPEED),
            RELAYER_LOOP_SPEED,
        );

        abscissa_tokio::run_with_actix(&APP, async {
            let connections = create_rpc_connections(
                cosmos_prefix,
                Some(config.cosmos.grpc.clone()),
                Some(config.ethereum.rpc.clone()),
                timeout,
            )
            .await;

            let mut grpc = connections.grpc.clone().unwrap();
            let contact = connections.contact.clone().unwrap();
            let provider = connections.eth_provider.clone().unwrap();
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
            wait_for_cosmos_node_ready(&contact).await;

            // check if the delegate addresses are correctly configured
            check_delegate_addresses(
                &mut grpc,
                ethereum_address,
                cosmos_address,
                &contact.get_prefix(),
            )
            .await;

            // check if we actually have the promised balance of tokens to pay fees
            check_for_fee_denom(&fees_denom, cosmos_address, &contact).await;
            check_for_eth(ethereum_address, eth_client.clone()).await;

            let gas_price = config.cosmos.gas_price.as_tuple();

            orchestrator_main_loop(
                cosmos_key,
                contact,
                eth_client,
                grpc,
                contract_address,
                gas_price,
                &config.metrics.listen_addr,
                config.ethereum.gas_price_multiplier,
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
