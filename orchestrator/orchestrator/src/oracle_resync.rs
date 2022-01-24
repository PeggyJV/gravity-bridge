use deep_space::address::Address as CosmosAddress;
use ethereum_gravity::types::EthClient;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_abi::gravity::*;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::types::{
    Erc20DeployedEvent, LogicCallExecutedEvent, SendToCosmosEvent, TransactionBatchExecutedEvent,
    ValsetUpdatedEvent,
};
use gravity_utils::types::{FromLog, FromLogWithPrefix};
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;

use crate::get_with_retry::get_block_number_with_retry;
use crate::get_with_retry::get_last_event_nonce_with_retry;
use crate::get_with_retry::RETRY_TIME;

/// This function retrieves the last event nonce that we have relayed to Cosmos
/// it then uses the Ethereum indexes to find what block the last event we relayed is in
pub async fn get_last_checked_block(
    grpc_client: GravityQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
    gravity_contract_address: EthAddress,
    eth_client: EthClient,
    blocks_to_search: u64,
) -> U64 {
    // TODO(bolten): original version of this used a 120 second timeout when querying
    // the eth chain, should we replicate that in eth_client?
    let mut grpc_client = grpc_client;
    let mut last_event_nonce: U256 =
        get_last_event_nonce_with_retry(&mut grpc_client, our_cosmos_address)
            .await
            .into();

    // zero indicates this oracle has never submitted an event before since there is no
    // zero event nonce (it's pre-incremented in the solidity contract) we have to go
    // and look for event nonce one.
    if last_event_nonce == 0u8.into() {
        last_event_nonce = 1u8.into();
    }

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

    let mut end_search_block = get_block_number_with_retry(eth_client.clone()).await;
    let blocks_to_search: U64 = blocks_to_search.into();

    while end_search_block > 0u8.into() {
        info!(
            "Oracle is resyncing, looking back into the history to find our last event nonce {}, on block {}",
            last_event_nonce, end_search_block
        );

        let start_search_block = end_search_block.saturating_sub(blocks_to_search);
        let search_range = start_search_block..end_search_block;

        // select uses an inclusive version of the range
        erc20_deployed_filter = erc20_deployed_filter.select(search_range.clone());
        logic_call_filter = logic_call_filter.select(search_range.clone());
        send_to_cosmos_filter = send_to_cosmos_filter.select(search_range.clone());
        transaction_batch_filter = transaction_batch_filter.select(search_range.clone());
        valset_updated_filter = valset_updated_filter.select(search_range.clone());

        let erc20_deployed_events = eth_client.get_logs(&erc20_deployed_filter).await;
        let logic_call_events = eth_client.get_logs(&logic_call_filter).await;
        let send_to_cosmos_events = eth_client.get_logs(&send_to_cosmos_filter).await;
        let transaction_batch_events = eth_client.get_logs(&transaction_batch_filter).await;
        // valset update events have one special property that is useful to us in this handler:
        // a valset update event for nonce 0 is emitted in the contract constructor meaning once you
        // find that event you can exit the search with confidence that you have not missed any events
        // without searching the entire blockchain history
        let valset_updated_events = eth_client.get_logs(&valset_updated_filter).await;

        if erc20_deployed_events.is_err()
            || logic_call_events.is_err()
            || send_to_cosmos_events.is_err()
            || transaction_batch_events.is_err()
            || valset_updated_events.is_err()
        {
            error!("Failed to get blockchain events while resyncing, is your Eth node working? If you see only one of these it's fine");
            delay_for(RETRY_TIME).await;
            continue;
        }

        let erc20_deployed_events = erc20_deployed_events.unwrap();
        let logic_call_events = logic_call_events.unwrap();
        let send_to_cosmos_events = send_to_cosmos_events.unwrap();
        let transaction_batch_events = transaction_batch_events.unwrap();
        let mut valset_updated_events = valset_updated_events.unwrap();

        // look for and return the block number of the event last seen on the Cosmos chain
        // then we will play events from that block (including that block, just in case
        // there is more than one event there) onwards. We use valset nonce 0 as an indicator
        // of what block the contract was deployed on.
        for event in erc20_deployed_events {
            match Erc20DeployedEvent::from_log(&event) {
                Ok(deploy) => {
                    trace!(
                        "{} ERC20 deploy event nonce, {} last event nonce",
                        deploy.event_nonce,
                        last_event_nonce
                    );
                    if deploy.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got ERC20DeployedEvent that we can't parse: {}", e),
            }
        }

        for event in logic_call_events {
            match LogicCallExecutedEvent::from_log(&event) {
                Ok(call) => {
                    trace!(
                        "{} logic call event nonce, {} last event nonce",
                        call.event_nonce,
                        last_event_nonce
                    );
                    if call.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got LogicCallExecutedEvent that we can't parse: {}", e),
            }
        }

        for event in send_to_cosmos_events {
            let prefix = our_cosmos_address.get_prefix();
            match SendToCosmosEvent::from_log(&event, prefix.as_str()) {
                Ok(send) => {
                    trace!(
                        "{} send to Cosmos event nonce, {} last event nonce",
                        send.event_nonce,
                        last_event_nonce
                    );
                    if send.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got SendToCosmosEvent that we can't parse: {}", e),
            }
        }

        for event in transaction_batch_events {
            match TransactionBatchExecutedEvent::from_log(&event) {
                Ok(batch) => {
                    trace!(
                        "{} transaction batch event nonce, {} last event nonce",
                        batch.event_nonce,
                        last_event_nonce
                    );
                    if batch.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!(
                    "Got TransactionBatchExecutedEvent that we can't parse: {}",
                    e
                ),
            }
        }

        // this reverse solves a very specific bug, we use the properties of the first valsets for edgecase
        // handling here, but events come in chronological order, so if we don't reverse the iterator
        // we will encounter the first validator sets first and exit early and incorrectly.
        // note that reversing everything won't actually get you that much of a performance gain
        // because this only involves events within the searching block range.
        valset_updated_events.reverse();
        for event in valset_updated_events {
            match ValsetUpdatedEvent::from_log(&event) {
                Ok(valset) => {
                    // if we've found this event it is the first possible event from the contract
                    // no other events can come before it, therefore either there's been a parsing error
                    // or no events have been submitted on this chain yet.
                    let bootstrapping =
                        valset.valset_nonce == 0u32.into() && last_event_nonce == 1u8.into();
                    // our last event was a valset update event, treat as normal case
                    let common_case =
                        valset.event_nonce == last_event_nonce && event.block_number.is_some();
                    trace!(
                        "{} valset updated event nonce, {} last event nonce",
                        valset.event_nonce,
                        last_event_nonce
                    );
                    if common_case || bootstrapping {
                        return event.block_number.unwrap();
                    }
                    // if we're looking for a later event nonce and we find the deployment of the contract
                    // we must have failed to parse the event we're looking for. The oracle can not start
                    else if valset.valset_nonce == 0u32.into() && last_event_nonce > 1u8.into() {
                        panic!("Could not find the last event relayed by {}, Last Event nonce is {} but no event matching that could be found!", our_cosmos_address, last_event_nonce)
                    }
                }
                Err(e) => error!("Got ValsetUpdatedEvent that we can't parse: {}", e),
            }
        }

        end_search_block = start_search_block.saturating_sub(1u8.into()); // filter ranges are inclusive, avoid searching same block
    }

    // we should exit above when we find the zero valset, if we have the wrong contract address through we could be at it a while as we go over
    // the entire history to 'prove' it.
    panic!("You have reached the end of block history without finding the Gravity contract deploy event! You must have the wrong contract address!");
}
