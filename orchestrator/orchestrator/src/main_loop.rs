//! This file contains the main loops for two distinct functions that just happen to reside int his same binary for ease of use. The Ethereum Signer and the Ethereum Oracle are both roles in Gravity
//! that can only be run by a validator. This single binary the 'Orchestrator' runs not only these two rules but also the untrusted role of a relayer, that does not need any permissions and has it's
//! own crate and binary so that anyone may run it.

use crate::metrics;
use crate::{
    ethereum_event_watcher::check_for_events, metrics::metrics_main_loop,
    oracle_resync::get_last_checked_block,
};
use clarity::{address::Address as EthAddress, Uint256};
use clarity::{utils::bytes_to_hex_str, PrivateKey as EthPrivateKey};
use cosmos_gravity::send::send_main_loop;
use cosmos_gravity::{
    build,
    query::{
        get_oldest_unsigned_logic_call, get_oldest_unsigned_transaction_batch,
        get_oldest_unsigned_valsets,
    },
};
use deep_space::client::ChainStatus;
use deep_space::error::CosmosGrpcError;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use deep_space::{Contact, Msg};
use ethereum_gravity::utils::get_gravity_id;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use relayer::main_loop::relayer_main_loop;
use std::convert::TryInto;
use std::{net, time::Duration};
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

/// The execution speed governing all loops in this file
/// which is to say all loops started by Orchestrator main
/// loop except the relayer loop
pub const ETH_SIGNER_LOOP_SPEED: Duration = Duration::from_secs(11);
pub const ETH_ORACLE_LOOP_SPEED: Duration = Duration::from_secs(13);

/// This loop combines the three major roles required to make
/// up the 'Orchestrator', all three of these are async loops
/// meaning they will occupy the same thread, but since they do
/// very little actual cpu bound work and spend the vast majority
/// of all execution time sleeping this shouldn't be an issue at all.
#[allow(clippy::many_single_char_names)]
#[allow(clippy::too_many_arguments)]
pub async fn orchestrator_main_loop(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: Web3,
    contact: Contact,
    grpc_client: GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    gas_price: (f64, String),
    metrics_listen: &net::SocketAddr,
    eth_gas_multiplier: f32,
    blocks_to_search: u128,
    gas_adjustment: f64,
    relayer_opt_out: bool,
    cosmos_msg_batch_size: u32,
) {
    let (tx, rx) = tokio::sync::mpsc::channel(1);

    let a = send_main_loop(
        &contact,
        cosmos_key,
        gas_price,
        rx,
        gas_adjustment,
        cosmos_msg_batch_size.try_into().unwrap(),
    );

    let b = eth_oracle_main_loop(
        cosmos_key,
        web3.clone(),
        contact.clone(),
        grpc_client.clone(),
        gravity_contract_address,
        blocks_to_search,
        tx.clone(),
    );

    let c = eth_signer_main_loop(
        cosmos_key,
        ethereum_key,
        web3.clone(),
        contact.clone(),
        grpc_client.clone(),
        gravity_contract_address,
        tx.clone(),
    );

    let d = metrics_main_loop(metrics_listen);

    if !relayer_opt_out {
        let e = relayer_main_loop(
            ethereum_key,
            web3,
            grpc_client.clone(),
            gravity_contract_address,
            eth_gas_multiplier,
        );
        futures::future::join5(a, b, c, d, e).await;
    } else {
        futures::future::join4(a, b, c, d).await;
    }
}

const DELAY: Duration = Duration::from_secs(5);

/// This function is responsible for making sure that Ethereum events are retrieved from the Ethereum blockchain
/// and ferried over to Cosmos where they will be used to issue tokens or process batches.
#[allow(unused_variables)]
pub async fn eth_oracle_main_loop(
    cosmos_key: CosmosPrivateKey,
    web3: Web3,
    contact: Contact,
    grpc_client: GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    blocks_to_search: u128,
    msg_sender: tokio::sync::mpsc::Sender<Vec<Msg>>,
) {
    let our_cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let long_timeout_web30 = Web3::new(&web3.get_url(), Duration::from_secs(120));
    let mut last_checked_block: Uint256 = get_last_checked_block(
        grpc_client.clone(),
        our_cosmos_address,
        gravity_contract_address,
        &long_timeout_web30,
        blocks_to_search,
    )
    .await;
    info!("Oracle resync complete, Oracle now operational");
    let mut grpc_client = grpc_client;

    loop {
        let (async_resp, _) = tokio::join!(
            async {
                let latest_eth_block = web3.eth_block_number().await;
                let latest_cosmos_block = contact.get_chain_status().await;
                match (latest_eth_block, latest_cosmos_block) {
                    (Ok(latest_eth_block), Ok(ChainStatus::Moving { block_height })) => {
                        metrics::set_cosmos_block_height(block_height.clone());
                        metrics::set_ethereum_block_height(latest_eth_block.clone());
                        trace!(
                            "Latest Eth block {} Latest Cosmos block {}",
                            latest_eth_block,
                            block_height,
                        );
                    }
                    (Ok(_latest_eth_block), Ok(ChainStatus::Syncing)) => {
                        warn!("Cosmos node syncing, Eth oracle paused");
                        delay_for(DELAY).await;
                    }
                    (Ok(_latest_eth_block), Ok(ChainStatus::WaitingToStart)) => {
                        warn!("Cosmos node syncing waiting for chain start, Eth oracle paused");
                        delay_for(DELAY).await;
                    }
                    (Ok(_), Err(_)) => {
                        metrics::COSMOS_UNAVAILABLE.inc();
                        warn!("Could not contact Cosmos grpc, trying again");
                        delay_for(DELAY).await;
                    }
                    (Err(_), Ok(_)) => {
                        metrics::ETHEREUM_UNAVAILABLE.inc();
                        warn!("Could not contact Eth node, trying again");
                        delay_for(DELAY).await;
                    }
                    (Err(_), Err(_)) => {
                        metrics::COSMOS_UNAVAILABLE.inc();
                        metrics::ETHEREUM_UNAVAILABLE.inc();
                        error!("Could not reach Ethereum or Cosmos rpc!");
                        delay_for(DELAY).await;
                    }
                }

                // Relays events from Ethereum -> Cosmos
                match check_for_events(
                    &web3,
                    &contact,
                    &mut grpc_client,
                    gravity_contract_address,
                    cosmos_key,
                    last_checked_block.clone(),
                    msg_sender.clone(),
                )
                .await
                {
                    Ok(new_block) => last_checked_block = new_block,
                    Err(e) => {
                        metrics::ETHEREUM_EVENT_CHECK_FAILURES.inc();
                        error!("Failed to get events for block range, Check your Eth node and Cosmos gRPC {:?}", e);
                        if let gravity_utils::error::GravityError::CosmosGrpcError(err) = e {
                            if let CosmosGrpcError::TransactionFailed { tx: _, time: _ } = err {
                                delay_for(Duration::from_secs(10)).await;
                            }
                        }
                    }
                }
            },
            delay_for(ETH_ORACLE_LOOP_SPEED)
        );
    }
}

/// The eth_signer simply signs off on any batches or validator sets provided by the validator
/// since these are provided directly by a trusted Cosmsos node they can simply be assumed to be
/// valid and signed off on.
#[allow(unused_variables)]
pub async fn eth_signer_main_loop(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: Web3,
    contact: Contact,
    grpc_client: GravityQueryClient<Channel>,
    contract_address: EthAddress,
    msg_sender: tokio::sync::mpsc::Sender<Vec<Msg>>,
) {
    let our_cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();
    let mut grpc_client = grpc_client;

    let gravity_id = get_gravity_id(contract_address, our_ethereum_address, &web3).await;
    if gravity_id.is_err() {
        error!("Failed to get GravityID, check your Eth node");
        return;
    }
    let gravity_id = gravity_id.unwrap();

    loop {
        let (async_resp, _) = tokio::join!(
            async {
                let latest_eth_block = web3.eth_block_number().await;
                let latest_cosmos_block = contact.get_chain_status().await;
                match (latest_eth_block, latest_cosmos_block) {
                    (Ok(latest_eth_block), Ok(ChainStatus::Moving { block_height })) => {
                        metrics::set_cosmos_block_height(block_height.clone());
                        metrics::set_ethereum_block_height(latest_eth_block.clone());
                        trace!(
                            "Latest Eth block {} Latest Cosmos block {}",
                            latest_eth_block,
                            block_height,
                        );
                    }
                    (Ok(_latest_eth_block), Ok(ChainStatus::Syncing)) => {
                        warn!("Cosmos node syncing, Eth signer paused");
                        delay_for(DELAY).await;
                    }
                    (Ok(_latest_eth_block), Ok(ChainStatus::WaitingToStart)) => {
                        warn!("Cosmos node syncing waiting for chain start, Eth signer paused");
                        delay_for(DELAY).await;
                    }
                    (Ok(_), Err(_)) => {
                        metrics::COSMOS_UNAVAILABLE.inc();
                        warn!("Could not contact Cosmos grpc, trying again");
                        delay_for(DELAY).await;
                    }
                    (Err(_), Ok(_)) => {
                        metrics::ETHEREUM_UNAVAILABLE.inc();
                        warn!("Could not contact Eth node, trying again");
                        delay_for(DELAY).await;
                    }
                    (Err(_), Err(_)) => {
                        metrics::COSMOS_UNAVAILABLE.inc();
                        metrics::ETHEREUM_UNAVAILABLE.inc();
                        error!("Could not reach Ethereum or Cosmos rpc!");
                        delay_for(DELAY).await;
                    }
                }

                // sign the last unsigned valsets
                match get_oldest_unsigned_valsets(&mut grpc_client, our_cosmos_address).await {
                    Ok(valsets) => {
                        if valsets.is_empty() {
                            trace!("No validator sets to sign, node is caught up!")
                        } else {
                            info!(
                                "Sending {} valset confirms starting with {}",
                                valsets.len(),
                                valsets[0].nonce
                            );
                            let messages = build::signer_set_tx_confirmation_messages(
                                &contact,
                                ethereum_key,
                                valsets,
                                cosmos_key,
                                gravity_id.clone(),
                            );
                            msg_sender
                                .send(messages)
                                .await
                                .expect("Could not send messages");
                        }
                    }
                    Err(e) => {
                        metrics::UNSIGNED_VALSET_FAILURES.inc();
                        error!(
                            "Failed to get unsigned valsets, check your Cosmos gRPC {:?}",
                            e
                        );
                    }
                }

                // sign the last unsigned batch, TODO check if we already have signed this
                match get_oldest_unsigned_transaction_batch(&mut grpc_client, our_cosmos_address)
                    .await
                {
                    Ok(Some(last_unsigned_batch)) => {
                        info!(
                            "Sending batch confirm for {}:{} fees {} timeout {}",
                            last_unsigned_batch.token_contract,
                            last_unsigned_batch.nonce,
                            last_unsigned_batch.total_fee.amount,
                            last_unsigned_batch.batch_timeout,
                        );
                        let transaction_batches = vec![last_unsigned_batch];
                        let messages = build::batch_tx_confirmation_messages(
                            &contact,
                            ethereum_key,
                            transaction_batches,
                            cosmos_key,
                            gravity_id.clone(),
                        );
                        msg_sender
                            .send(messages)
                            .await
                            .expect("Could not send messages");
                    }
                    Ok(None) => info!("No unsigned batches! Everything good!"),
                    Err(e) => {
                        metrics::UNSIGNED_BATCH_FAILURES.inc();
                        error!(
                            "Failed to get unsigned Batches, check your Cosmos gRPC {:?}",
                            e
                        );
                    }
                }

                let logic_calls =
                    get_oldest_unsigned_logic_call(&mut grpc_client, our_cosmos_address).await;
                if let Ok(logic_calls) = logic_calls {
                    for logic_call in logic_calls {
                        info!(
                            "Sending Logic call confirm for {}:{}",
                            bytes_to_hex_str(&logic_call.invalidation_id),
                            logic_call.invalidation_nonce
                        );
                        let logic_calls = vec![logic_call];
                        let messages = build::contract_call_tx_confirmation_messages(
                            &contact,
                            ethereum_key,
                            logic_calls,
                            cosmos_key,
                            gravity_id.clone(),
                        );
                        msg_sender
                            .send(messages)
                            .await
                            .expect("Could not send messages");
                    }
                } else if let Err(e) = logic_calls {
                    metrics::UNSIGNED_LOGIC_CALL_FAILURES.inc();
                    error!(
                        "Failed to get unsigned Logic Calls, check your Cosmos gRPC {:?}",
                        e
                    );
                }
            },
            delay_for(ETH_SIGNER_LOOP_SPEED)
        );
    }
}

#[allow(dead_code)]
pub async fn check_for_eth(orchestrator_address: EthAddress, web3: Web3) {
    let balance = web3.eth_get_balance(orchestrator_address).await.unwrap();
    if balance == 0u8.into() {
        warn!("You don't have any Ethereum! You will need to send some to {} for this program to work. Dust will do for basic operations, more info about average relaying costs will be presented as the program runs", orchestrator_address);
    }
    metrics::set_ethereum_bal(balance);
}
