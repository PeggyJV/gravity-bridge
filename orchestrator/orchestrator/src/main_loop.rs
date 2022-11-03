//! This file contains the main loops for two distinct functions that just happen to reside in this same binary for ease of use.
//! The Ethereum Signer and the Ethereum Oracle are both roles in Gravity that can only be run by a validator. This single binary
//! the 'Orchestrator' runs not only these two roles but also the untrusted role of a relayer, that does not need any permissions
//! and has its own crate and binary so that anyone may run it.

use crate::ethereum_event_watcher::get_block_delay;
use crate::metrics;
use crate::{
    ethereum_event_watcher::check_for_events, metrics::metrics_main_loop,
    oracle_resync::get_last_checked_block,
};
use cosmos_gravity::send::send_main_loop;
use cosmos_gravity::{
    build,
    query::{
        get_oldest_unsigned_logic_call, get_oldest_unsigned_transaction_batch,
        get_oldest_unsigned_valsets,
    },
};
use ethereum_gravity::types::EthClient;
use ethereum_gravity::utils::get_gravity_id;
use ethers::{prelude::*, types::Address as EthAddress};
use gravity_proto::cosmos_sdk_proto::Any;
use gravity_utils::connection_prep::{get_chain_status, ChainStatus};
use gravity_utils::error::GravityError;
use gravity_utils::ethereum::bytes_to_hex_str;
use ocular::chain::ChainContext;
use ocular::cosmrs::AccountId;
use ocular::prelude::*;
use relayer::main_loop::relayer_main_loop;
use std::convert::TryInto;
use std::process::exit;
use std::{net, time::Duration};
use tokio::time::sleep as delay_for;

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
    grpc_endpoint: &str,
    cosmos_signer: &AccountInfo,
    context: ChainContext,
    eth_client: EthClient,
    gravity_contract_address: EthAddress,
    gas_price: (f64, String),
    metrics_listen: &net::SocketAddr,
    eth_gas_price_multiplier: f32,
    eth_gas_multiplier: f32,
    blocks_to_search: u64,
    gas_adjustment: f64,
    relayer_opt_out: bool,
    cosmos_msg_batch_size: u32,
) {
    let (tx, rx) = tokio::sync::mpsc::channel(1);
    let cosmos_account = cosmos_signer.id(&context.prefix).unwrap();
    let send_cosmos_client = GrpcClient::new(grpc_endpoint).await.unwrap();
    let a = send_main_loop(
        send_cosmos_client,
        cosmos_signer,
        &context,
        gas_price,
        gas_adjustment,
        rx,
        cosmos_msg_batch_size.try_into().unwrap(),
    );

    let eth_oracle_cosmos_client = GrpcClient::new(grpc_endpoint).await.unwrap();
    let b = eth_oracle_main_loop(
        eth_oracle_cosmos_client,
        &cosmos_account,
        eth_client.clone(),
        gravity_contract_address,
        blocks_to_search,
        tx.clone(),
    );

    let eth_signer_cosmos_client = GrpcClient::new(grpc_endpoint).await.unwrap();
    let c = eth_signer_main_loop(
        eth_signer_cosmos_client,
        &cosmos_account,
        eth_client.clone(),
        gravity_contract_address,
        tx.clone(),
    );

    let d = metrics_main_loop(metrics_listen);

    if !relayer_opt_out {
        let relayer_cosmos_client = GrpcClient::new(grpc_endpoint).await.unwrap();
        let e = relayer_main_loop(
            eth_client.clone(),
            relayer_cosmos_client,
            gravity_contract_address,
            eth_gas_price_multiplier,
            eth_gas_multiplier,
        );
        futures::future::join5(a, b, c, d, e).await;
    } else {
        futures::future::join4(a, b, c, d).await;
    }
}

// the amount of time to wait when encountering error conditions
const DELAY: Duration = Duration::from_secs(5);

// the number of loop iterations to wait between sending height update messages
const HEIGHT_UPDATE_INTERVAL: u32 = 50;

/// This function is responsible for making sure that Ethereum events are retrieved from the Ethereum blockchain
/// and ferried over to Cosmos where they will be used to issue tokens or process batches.
pub async fn eth_oracle_main_loop(
    cosmos_client: GrpcClient,
    cosmos_account: &AccountId,
    eth_client: EthClient,
    gravity_contract_address: EthAddress,
    blocks_to_search: u64,
    msg_sender: tokio::sync::mpsc::Sender<Vec<Any>>,
) {
    let mut cosmos_client = cosmos_client;
    let block_delay = match get_block_delay(eth_client.clone()).await {
        Ok(block_delay) => block_delay,
        Err(e) => {
            error!(
                "Error encountered when retrieving block delay, cannot continue: {}",
                e
            );
            exit(1);
        }
    };
    let mut last_checked_block = get_last_checked_block(
        &mut cosmos_client,
        cosmos_account,
        gravity_contract_address,
        eth_client.clone(),
        blocks_to_search,
    )
    .await;
    info!("Oracle resync complete, Oracle now operational");
    let mut loop_count: u32 = 0;

    loop {
        let (_async_resp, _) = tokio::join!(
            async {
                let latest_eth_block = eth_client.get_block_number().await;
                let latest_cosmos_block = get_chain_status(&mut cosmos_client).await;
                if let (Ok(latest_eth_block), ChainStatus::Moving(block)) =
                    (&latest_eth_block, latest_cosmos_block.clone())
                {
                    let block_height = block.header.unwrap().height as u64;
                    metrics::set_cosmos_block_height(block_height);
                    metrics::set_ethereum_block_height(latest_eth_block.as_u64());
                    trace!(
                        "Latest Eth block {} Latest Cosmos block {}",
                        latest_eth_block,
                        block_height,
                    );

                    // send latest Ethereum height to the Cosmos chain periodically
                    // subtract the block delay based on the environment, in order to have
                    // more confidence we are attesting to a height that has not been re-orged
                    if loop_count % HEIGHT_UPDATE_INTERVAL == 0 {
                        let messages = build::ethereum_vote_height_messages(
                            cosmos_account,
                            (latest_eth_block - block_delay).as_u64(),
                        )
                        .await;

                        msg_sender
                            .send(messages)
                            .await
                            .expect("Could not send Ethereum height votes");
                    }
                } else {
                    match latest_cosmos_block {
                        ChainStatus::Syncing => warn!("Cosmos node is syncing."),
                        ChainStatus::WaitingToStart(error) => {
                            if error.is_some() {
                                warn!("Cosmos node is waiting to start: {}", error.unwrap());
                                metrics::COSMOS_UNAVAILABLE.inc();
                            } else {
                                warn!("Cosmos node is waiting to start.")
                            }
                        }
                        _ => (),
                    }

                    if let Err(e) = latest_eth_block {
                        warn!("Could not contact Eth node, trying again: {}", e);
                        metrics::ETHEREUM_UNAVAILABLE.inc();
                    }

                    delay_for(DELAY).await;
                }

                // Relays events from Ethereum -> Cosmos
                match check_for_events(
                    &mut cosmos_client,
                    cosmos_account,
                    eth_client.clone(),
                    gravity_contract_address,
                    last_checked_block,
                    blocks_to_search.into(),
                    block_delay,
                    msg_sender.clone(),
                )
                .await
                {
                    Ok(new_block) => last_checked_block = new_block,
                    Err(e) => {
                        metrics::ETHEREUM_EVENT_CHECK_FAILURES.inc();
                        error!("Failed to get events for block range, Check your Eth node and Cosmos gRPC {:?}", e);
                        if let GravityError::CosmosGrpcError(_) = e {
                            delay_for(Duration::from_secs(10)).await;
                        }
                    }
                }
            },
            delay_for(ETH_ORACLE_LOOP_SPEED)
        );

        loop_count += 1;
    }
}

/// The eth_signer simply signs off on any batches or validator sets provided by the validator
/// since these are provided directly by a trusted Cosmsos node they can simply be assumed to be
/// valid and signed off on.
#[allow(unused_variables)]
pub async fn eth_signer_main_loop(
    cosmos_client: GrpcClient,
    cosmos_account: &AccountId,
    eth_client: EthClient,
    contract_address: EthAddress,
    msg_sender: tokio::sync::mpsc::Sender<Vec<Any>>,
) {
    let mut cosmos_client = cosmos_client;
    let gravity_id = get_gravity_id(contract_address, eth_client.clone()).await;
    if gravity_id.is_err() {
        error!("Failed to get GravityID, check your Eth node");
        return;
    }
    let gravity_id = gravity_id.unwrap();

    loop {
        let (async_resp, _) = tokio::join!(
            async {
                let latest_eth_block = eth_client.get_block_number().await;
                let latest_cosmos_block = get_chain_status(&mut cosmos_client).await;
                if let (Ok(latest_eth_block), ChainStatus::Moving(block)) =
                    (&latest_eth_block, latest_cosmos_block.clone())
                {
                    let block_height = block.header.unwrap().height as u64;
                    metrics::set_cosmos_block_height(block_height);
                    metrics::set_ethereum_block_height(latest_eth_block.as_u64());
                    trace!(
                        "Latest Eth block {} Latest Cosmos block {}",
                        latest_eth_block,
                        block_height,
                    );
                } else {
                    match latest_cosmos_block {
                        ChainStatus::Syncing => warn!("Cosmos node is syncing."),
                        ChainStatus::WaitingToStart(error) => {
                            if error.is_some() {
                                warn!("Cosmos node is waiting to start: {}", error.unwrap());
                                metrics::COSMOS_UNAVAILABLE.inc();
                            } else {
                                warn!("Cosmos node is waiting to start.")
                            }
                        }
                        _ => (),
                    }

                    if let Err(e) = latest_eth_block {
                        warn!("Could not contact Eth node, trying again: {}", e);
                        metrics::ETHEREUM_UNAVAILABLE.inc();
                    }

                    delay_for(DELAY).await;
                }

                // sign the last unsigned valsets
                match get_oldest_unsigned_valsets(&mut cosmos_client, cosmos_account).await {
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
                                cosmos_account,
                                eth_client.clone(),
                                valsets,
                                gravity_id.clone(),
                            )
                            .await;
                            msg_sender
                                .send(messages)
                                .await
                                .expect("Could not send messages");
                        }
                    }
                    Err(e) => {
                        metrics::UNSIGNED_VALSET_FAILURES.inc();
                        error!(
                            "Failed to get unsigned valset, check your Cosmos gRPC {:?}",
                            e
                        );
                    }
                }

                // sign the last unsigned batch, TODO check if we already have signed this
                match get_oldest_unsigned_transaction_batch(&mut cosmos_client, cosmos_account)
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
                            cosmos_account,
                            eth_client.clone(),
                            transaction_batches,
                            gravity_id.clone(),
                        )
                        .await;
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
                    get_oldest_unsigned_logic_call(&mut cosmos_client, cosmos_account).await;
                if let Ok(logic_calls) = logic_calls {
                    for logic_call in logic_calls {
                        info!(
                            "Sending Logic call confirm for {}:{}",
                            bytes_to_hex_str(&logic_call.invalidation_id),
                            logic_call.invalidation_nonce
                        );
                        let logic_calls = vec![logic_call];
                        let messages = build::contract_call_tx_confirmation_messages(
                            cosmos_account,
                            eth_client.clone(),
                            logic_calls,
                            gravity_id.clone(),
                        )
                        .await;
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
                    )
                }
            },
            delay_for(ETH_SIGNER_LOOP_SPEED)
        );
    }
}

pub async fn check_for_eth(orchestrator_address: EthAddress, eth_client: EthClient) {
    let balance = eth_client
        .get_balance(orchestrator_address, None)
        .await
        .unwrap();
    if balance == 0u8.into() {
        warn!("You don't have any Ethereum! You will need to send some to {} for this program to work. Dust will do for basic operations, more info about average relaying costs will be presented as the program runs", orchestrator_address);
    }
    metrics::set_ethereum_bal(balance);
}
