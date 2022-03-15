//! Ethereum Event watcher watches for events such as a deposit to the Gravity Ethereum contract or a validator set update
//! or a transaction batch update. It then responds to these events by performing actions on the Cosmos chain if required

use crate::get_with_retry::get_block_number_with_retry;
use crate::get_with_retry::get_chain_id_with_retry;
use crate::metrics;
use cosmos_gravity::build;
use cosmos_gravity::query::get_last_event_nonce;
use cosmos_gravity::crypto::PrivateKey as CosmosPrivateKey;
use deep_space::{Contact, Msg};
use ethereum_gravity::types::EthClient;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_abi::gravity::*;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::ethereum::downcast_to_u64;
use gravity_utils::types::EventNonceFilter;
use gravity_utils::types::{FromLogs, FromLogsWithPrefix};
use gravity_utils::{
    error::GravityError,
    ethereum::bytes_to_hex_str,
    types::{
        Erc20DeployedEvent, LogicCallExecutedEvent, SendToCosmosEvent,
        TransactionBatchExecutedEvent, ValsetUpdatedEvent,
    },
};
use std::{result::Result, time};
use tonic::transport::Channel;

#[allow(clippy::too_many_arguments)]
pub async fn check_for_events(
    eth_client: EthClient,
    contact: &Contact,
    grpc_client: &mut GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    cosmos_key: CosmosPrivateKey,
    starting_block: U64,
    block_delay: U64,
    msg_sender: tokio::sync::mpsc::Sender<Vec<Msg>>,
) -> Result<U64, GravityError> {
    let prefix = contact.get_prefix();
    let our_cosmos_address = cosmos_key.to_address(&prefix).unwrap();
    let latest_block = get_block_number_with_retry(eth_client.clone()).await;
    let latest_block = latest_block - block_delay;

    metrics::set_ethereum_check_for_events_starting_block(starting_block.as_u64());
    metrics::set_ethereum_check_for_events_end_block(latest_block.as_u64());

    let filter_gravity_contract_address = ValueOrArray::Value(gravity_contract_address);

    let mut erc20_deployed_filter = Filter::new()
        .address(filter_gravity_contract_address.clone())
        .event(&Erc20DeployedEventFilter::abi_signature());
    let mut logic_call_filter = Filter::new()
        .address(filter_gravity_contract_address.clone())
        .event(&LogicCallEventFilter::abi_signature());
    let mut send_to_cosmos_filter = Filter::new()
        .address(filter_gravity_contract_address.clone())
        .event(&SendToCosmosEventFilter::abi_signature());
    let mut transaction_batch_filter = Filter::new()
        .address(filter_gravity_contract_address.clone())
        .event(&TransactionBatchExecutedEventFilter::abi_signature());
    let mut valset_updated_filter = Filter::new()
        .address(filter_gravity_contract_address.clone())
        .event(&ValsetUpdatedEventFilter::abi_signature());

    let search_range = starting_block..latest_block;

    // select uses an inclusive version of the range
    erc20_deployed_filter = erc20_deployed_filter.select(search_range.clone());
    logic_call_filter = logic_call_filter.select(search_range.clone());
    send_to_cosmos_filter = send_to_cosmos_filter.select(search_range.clone());
    transaction_batch_filter = transaction_batch_filter.select(search_range.clone());
    valset_updated_filter = valset_updated_filter.select(search_range.clone());

    let erc20_deployed_events = eth_client.get_logs(&erc20_deployed_filter).await?;
    debug!("ERC20 events detected {:?}", erc20_deployed_events);
    let erc20_deployed_events = Erc20DeployedEvent::from_logs(&erc20_deployed_events)?;
    debug!("parsed erc20 deploys {:?}", erc20_deployed_events);

    let logic_call_events = eth_client.get_logs(&logic_call_filter).await?;
    debug!("Logic call events detected {:?}", logic_call_events);
    let logic_call_events = LogicCallExecutedEvent::from_logs(&logic_call_events)?;
    debug!("parsed logic call executions {:?}", logic_call_events);

    let send_to_cosmos_events = eth_client.get_logs(&send_to_cosmos_filter).await?;
    debug!("Send to Cosmos events detected {:?}", send_to_cosmos_events);
    let send_to_cosmos_events = SendToCosmosEvent::from_logs(&send_to_cosmos_events, &prefix)?;
    debug!("parsed send to cosmos events {:?}", send_to_cosmos_events);

    let transaction_batch_events = eth_client.get_logs(&transaction_batch_filter).await?;
    debug!("Batch events detected {:?}", transaction_batch_events);
    let transaction_batch_events =
        TransactionBatchExecutedEvent::from_logs(&transaction_batch_events)?;
    debug!("parsed batches {:?}", transaction_batch_events);

    let valset_updated_events = eth_client.get_logs(&valset_updated_filter).await?;
    debug!("Valset events detected {:?}", valset_updated_events);
    let valset_updated_events = ValsetUpdatedEvent::from_logs(&valset_updated_events)?;
    debug!("parsed valsets {:?}", valset_updated_events);

    // note that starting block overlaps with our last checked block, because we have to deal with
    // the possibility that the relayer was killed after relaying only one of multiple events in a single
    // block, so we also need this routine so make sure we don't send in the first event in this hypothetical
    // multi event block again. In theory we only send all events for every block and that will pass of fail
    // atomicly but lets not take that risk.
    let last_event_nonce = get_last_event_nonce(grpc_client, our_cosmos_address).await?;
    metrics::set_cosmos_last_event_nonce(last_event_nonce);

    let erc20_deployed_events: Vec<Erc20DeployedEvent> =
        Erc20DeployedEvent::filter_by_event_nonce(last_event_nonce, &erc20_deployed_events);
    let logic_call_events: Vec<LogicCallExecutedEvent> =
        LogicCallExecutedEvent::filter_by_event_nonce(last_event_nonce, &logic_call_events);
    let send_to_cosmos_events: Vec<SendToCosmosEvent> =
        SendToCosmosEvent::filter_by_event_nonce(last_event_nonce, &send_to_cosmos_events);
    let transaction_batch_events: Vec<TransactionBatchExecutedEvent> =
        TransactionBatchExecutedEvent::filter_by_event_nonce(
            last_event_nonce,
            &transaction_batch_events,
        );
    let valset_updated_events: Vec<ValsetUpdatedEvent> =
        ValsetUpdatedEvent::filter_by_event_nonce(last_event_nonce, &valset_updated_events);

    for erc20_deployed_event in erc20_deployed_events.iter() {
        info!(
            "Oracle observed ERC20 deploy with denom {}, erc20 name {}, symbol {}, and event_nonce {}",
            erc20_deployed_event.cosmos_denom,
            erc20_deployed_event.name,
            erc20_deployed_event.symbol,
            erc20_deployed_event.event_nonce,
        )
    }

    for logic_call_event in logic_call_events.iter() {
        info!(
            "Oracle observed logic call execution with invalidation_id {}, invalidation_nonce {}, and event_nonce {}",
            bytes_to_hex_str(&logic_call_event.invalidation_id),
            logic_call_event.invalidation_nonce,
            logic_call_event.event_nonce
        );
    }

    for send_to_cosmos_event in send_to_cosmos_events.iter() {
        info!(
            "Oracle observed send to cosmos event with ethereum sender {}, cosmos receiver {}, amount {}, and event nonce {}",
            send_to_cosmos_event.sender,
            send_to_cosmos_event.destination,
            send_to_cosmos_event.amount,
            send_to_cosmos_event.event_nonce
        );
    }

    for transaction_batch_event in transaction_batch_events.iter() {
        info!(
            "Oracle observed batch with batch_nonce {}, erc20 {}, and event_nonce {}",
            transaction_batch_event.batch_nonce,
            transaction_batch_event.erc20,
            transaction_batch_event.event_nonce
        );
    }

    for valset_updated_event in valset_updated_events.iter() {
        info!(
            "Oracle observed valset update with valset_nonce {}, event_nonce {}, block_height {}, and members {:?}",
            valset_updated_event.valset_nonce,
            valset_updated_event.event_nonce,
            valset_updated_event.block_height,
            valset_updated_event.members
        )
    }

    if !erc20_deployed_events.is_empty()
        || !logic_call_events.is_empty()
        || !send_to_cosmos_events.is_empty()
        || !transaction_batch_events.is_empty()
        || !valset_updated_events.is_empty()
    {
        let messages = build::ethereum_event_messages(
            contact,
            cosmos_key,
            send_to_cosmos_events.to_owned(),
            transaction_batch_events.to_owned(),
            erc20_deployed_events.to_owned(),
            logic_call_events.to_owned(),
            valset_updated_events.to_owned(),
        );

        info!("Sending {} messages to cosmos", messages.len());

        if let Some(erc20_deployed_event) = erc20_deployed_events.last() {
            metrics::set_ethereum_last_erc20_event(erc20_deployed_event.event_nonce);
            metrics::set_ethereum_last_erc20_block(erc20_deployed_event.block_height);
        }

        if let Some(send_to_cosmos_event) = send_to_cosmos_events.last() {
            metrics::set_ethereum_last_deposit_event(send_to_cosmos_event.event_nonce);
            metrics::set_ethereum_last_deposit_block(send_to_cosmos_event.block_height);
        }

        if let Some(transaction_batch_event) = transaction_batch_events.last() {
            metrics::set_ethereum_last_batch_event(transaction_batch_event.event_nonce);
            metrics::set_ethereum_last_batch_nonce(transaction_batch_event.batch_nonce);
        }

        if let Some(valset_updated_event) = valset_updated_events.last() {
            metrics::set_ethereum_last_valset_event(valset_updated_event.event_nonce);
            metrics::set_ethereum_last_valset_nonce(valset_updated_event.valset_nonce);
        }

        if let Some(logic_call_event) = logic_call_events.last() {
            metrics::set_ethereum_last_logic_call_event(logic_call_event.event_nonce);
            metrics::set_ethereum_last_logic_call_nonce(logic_call_event.invalidation_nonce);
        }

        msg_sender
            .send(messages)
            .await
            .expect("Could not send messages");

        let timeout = time::Duration::from_secs(30);
        contact.wait_for_next_block(timeout).await?;

        // TODO(bolten): we are only waiting one block, is it possible if we are sending multiple
        // events via the sender, they could be received over the block boundary and thus our new
        // event nonce does not reflect full processing of the above events?
        let new_event_nonce = get_last_event_nonce(grpc_client, our_cosmos_address).await?;
        if new_event_nonce == last_event_nonce {
            return Err(GravityError::InvalidBridgeStateError(
                format!("Claims did not process, trying to update but still on {}, trying again in a moment", last_event_nonce),
            ));
        } else {
            info!("Claims processed, new nonce: {}", new_event_nonce);
        }
    }

    Ok(latest_block)
}

/// The number of blocks behind the 'latest block' on Ethereum our event checking should be.
/// Ethereum does not have finality and as such is subject to chain reorgs and temporary forks
/// if we check for events up to the very latest block we may process an event which did not
/// 'actually occur' in the longest POW chain.
///
/// Obviously we must chose some delay in order to prevent incorrect events from being claimed
///
/// For EVM chains with finality the correct value for this is zero. As there's no need
/// to concern ourselves with re-orgs or forking. This function checks the netID of the
/// provided Ethereum RPC and adjusts the block delay accordingly
///
/// The value used here for Ethereum is a balance between being reasonably fast and reasonably secure
/// As you can see on https://etherscan.io/blocks_forked uncles (one block deep reorgs)
/// occur once every few minutes. Two deep once or twice a day.
/// https://etherscan.io/chart/uncles
/// Let's make a conservative assumption of 1% chance of an uncle being a two block deep reorg
/// (actual is closer to 0.3%) and assume that continues as we increase the depth.
/// Given an uncle every 2.8 minutes, a 6 deep reorg would be 2.8 minutes * (100^4) or one
/// 6 deep reorg every 53,272 years.
///
pub async fn get_block_delay(eth_client: EthClient) -> Result<U64, GravityError> {
    // TODO(bolten): get_net_version() exists on the version of ethers we are currently
    // depending on, but it's broken, so we're relying on chain ID
    let chain_id_result = get_chain_id_with_retry(eth_client.clone()).await;
    let chain_id = downcast_to_u64(chain_id_result);
    if chain_id.is_none() {
        return Err(GravityError::EthereumBadDataError(format!(
            "Chain ID is larger than u64 max: {}",
            chain_id_result
        )));
    }

    match chain_id.unwrap() {
        // Mainline Ethereum, Ethereum classic, or Ropsten, Mordor testnets
        // all POW Chains
        1 | 3 | 61 | 63 => Ok(13u8.into()),
        // Dev, our own Gravity Ethereum testnet, Hardhat
        // all non-pow chains
        2018 | 15 | 31337 => Ok(0u8.into()),
        // Rinkeby, Goerli, Kotti
        // Clique (POA) Consensus
        4 | 5 | 6 => Ok(10u8.into()),
        // assume the safe option (POW) where we don't know
        _ => Ok(13u8.into()),
    }
}
