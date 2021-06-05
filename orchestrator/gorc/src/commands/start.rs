//! `start` subcommand - example of how to write a subcommand

use crate::{application::APP, prelude::*, utils::*};
/// App-local prelude includes `app_reader()`/`app_writer()`/`app_config()`
/// accessors along with logging macros. Customize as you see fit.
use abscissa_core::{Command, Options, Runnable};
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use env_logger::Env;
use gravity_utils::connection_prep::{
    check_delegate_addresses, check_for_eth, check_for_fee_denom, create_rpc_connections,
    wait_for_cosmos_node_ready,
};
use orchestrator::main_loop::{
    orchestrator_main_loop, ETH_ORACLE_LOOP_SPEED, ETH_SIGNER_LOOP_SPEED,
};
use relayer::main_loop::{relayer_main_loop, LOOP_SPEED as RELAYER_LOOP_SPEED};
use std::cmp::min;

/// `start` subcommand

#[derive(Command, Debug, Options)]
pub enum StartCmd {
    /// To whom are we saying hello?
    #[options(help = "orchestrator [contract-address] [fee-denom]")]
    Orchestrator(Orchestrator),

    #[options(help = "relayer")]
    Relayer(Relayer),
}

impl Runnable for StartCmd {
    /// Start the application.
    fn run(&self) {
        //Your code goes here
    }
}

#[derive(Command, Debug, Options)]
pub struct Orchestrator {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

fn get_cosmos_key(key_name: &str) -> CosmosPrivateKey {
    unimplemented!()
}

fn lookup_eth_key(key: String) -> EthPrivateKey {
    todo!()
}

impl Runnable for Orchestrator {
    fn run(&self) {
        let timeout = min(
            min(ETH_SIGNER_LOOP_SPEED, ETH_ORACLE_LOOP_SPEED),
            RELAYER_LOOP_SPEED,
        );

        let config = APP.config();
        let cosmos_prefix = config.cosmos.prefix.clone();
        let cosmso_grpc = config.cosmos.grpc.clone();
        let cosmos_key = get_cosmos_key(&config.cosmos.key);
        let eth_key = lookup_eth_key(config.ethereum.key.clone());
        let ethereum_public_key = eth_key.to_public_key().unwrap();
        let fee_denom = config.gravity.fee_denom.clone();
        let contract_address: EthAddress = config
            .gravity
            .contract
            .clone()
            .parse()
            .expect("Expected config.gravity.contract to be an Eth ddress");

        abscissa_tokio::run(&APP, async {
            trace!("Probing RPC connections");
            // probe all rpc connections and see if they are valid

            let connections =
                create_rpc_connections(cosmos_prefix, Some(cosmso_grpc), None, TIMEOUT).await;
            let mut grpc = connections.grpc.clone().unwrap();
            let contact = connections.contact.clone().unwrap();
            let web3 = connections.web3.clone().unwrap();
            let public_cosmos_key = cosmos_key.to_address(&contact.get_prefix()).unwrap();
            info!("Starting Gravity Validator companion binary Relayer + Oracle + Eth Signer");
            info!(
                "Ethereum Address: {} Cosmos Address {}",
                ethereum_public_key, public_cosmos_key
            );

            // check if the cosmos node is syncing, if so wait for it
            // we can't move any steps above this because they may fail on an incorrect
            // historic chain state while syncing occurs
            wait_for_cosmos_node_ready(&contact).await;

            // check if the delegate addresses are correctly configured
            check_delegate_addresses(
                &mut grpc,
                ethereum_public_key,
                public_cosmos_key,
                &contact.get_prefix(),
            )
            .await;

            // check if we actually have the promised balance of tokens to pay fees
            check_for_fee_denom(&fee_denom, public_cosmos_key, &contact).await;
            check_for_eth(ethereum_public_key, &web3).await;

            orchestrator_main_loop(
                cosmos_key,
                eth_key,
                connections.web3.unwrap(),
                connections.contact.unwrap(),
                connections.grpc.unwrap(),
                contract_address,
                fee_denom,
            )
            .await;
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}

#[derive(Command, Debug, Options)]
pub struct Relayer {
    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Relayer {
    /// Start the application.
    fn run(&self) {
        env_logger::Builder::from_env(Env::default().default_filter_or("info")).init();
        // On Linux static builds we need to probe ssl certs path to be able to
        // do TLS stuff.
        openssl_probe::init_ssl_cert_env_vars();
        let config = APP.config();
        let cosmos_prefix = config.cosmos.prefix.clone();
        let cosmso_grpc = config.cosmos.grpc.clone();
        let cosmos_key = get_cosmos_key(&config.cosmos.key);
        let eth_key = lookup_eth_key(config.ethereum.key.clone());
        let ethereum_public_key = eth_key.to_public_key().unwrap();
        let fee_denom = config.gravity.fee_denom.clone();
        let contract_address: EthAddress = config
            .gravity
            .contract
            .clone()
            .parse()
            .expect("Expected config.gravity.contract to be an Eth ddress");

        abscissa_tokio::run(&APP, async {
            trace!("Probing RPC connections");
            // probe all rpc connections and see if they are valid

            let connections =
                create_rpc_connections(cosmos_prefix, Some(cosmso_grpc), None, TIMEOUT).await;
            let mut grpc = connections.grpc.clone().unwrap();
            let contact = connections.contact.clone().unwrap();
            let web3 = connections.web3.clone().unwrap();
            let public_cosmos_key = cosmos_key.to_address(&contact.get_prefix()).unwrap();
            info!("Starting Gravity Validator companion binary Relayer + Oracle + Eth Signer");
            info!(
                "Ethereum Address: {} Cosmos Address {}",
                ethereum_public_key, public_cosmos_key
            );

            // check if the cosmos node is syncing, if so wait for it
            // we can't move any steps above this because they may fail on an incorrect
            // historic chain state while syncing occurs
            wait_for_cosmos_node_ready(&contact).await;
            check_for_eth(ethereum_public_key, &web3).await;

            relayer_main_loop(
                eth_key,
                connections.web3.unwrap(),
                connections.grpc.unwrap(),
                contract_address,
            )
            .await
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}
