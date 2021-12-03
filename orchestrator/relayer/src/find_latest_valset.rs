use ethereum_gravity::types::EthClient;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_abi::gravity::*;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::types::{FromLog, ValsetUpdatedEvent};
use gravity_utils::{error::GravityError, ethereum::downcast_to_u64, types::Valset};
use std::{panic, result::Result};
use tonic::transport::Channel;

/// This function finds the latest valset on the Gravity contract by looking back through the event
/// history and finding the most recent ValsetUpdatedEvent. Most of the time this will be very fast
/// as the latest update will be in recent blockchain history and the search moves from the present
/// backwards in time. In the case that the validator set has not been updated for a very long time
/// this will take longer.
pub async fn find_latest_valset(
    grpc_client: &mut GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    eth_client: EthClient,
) -> Result<Valset, GravityError> {
    // calculate some constant U64 values only once
    const BLOCKS_TO_SEARCH: u64 = 5_000u64;

    let mut filter = Filter::new()
        .address(ValueOrArray::Value(gravity_contract_address))
        .event(&ValsetUpdatedEventFilter::abi_signature());
    let mut end_filter_block = eth_client.get_block_number().await?;

    while end_filter_block > 0u64.into() {
        debug!("About to submit a Valset or Batch, looking back into the history to find the last Valset Update, on block {}", end_filter_block);

        let start_filter_block = end_filter_block.saturating_sub(BLOCKS_TO_SEARCH.into());
        filter = filter.select(start_filter_block..end_filter_block);

        let mut filtered_logged_events = eth_client.get_logs(&filter).await?;
        filtered_logged_events.reverse(); // we'll process these in reverse order to start from the most recent and work backwards

        // TODO(bolten): the original logic only checked one valset event, even if there may have been multiple within the
        // filtered blockspace...need more clarity on how severe an error it is if one of these events is malformed, and if
        // we should return early with an error or just log it the way the previous version did
        for logged_event in filtered_logged_events {
            debug!("Found event {:?}", logged_event);

            match ValsetUpdatedEvent::from_log(&logged_event) {
                Ok(valset_updated_event) => {
                    let downcast_nonce = downcast_to_u64(valset_updated_event.valset_nonce);
                    if downcast_nonce.is_none() {
                        error!(
                            "ValsetUpdatedEvent has nonce larger than u64: {:?}",
                            valset_updated_event
                        );
                        continue;
                    }

                    let latest_eth_valset = Valset {
                        nonce: downcast_nonce.unwrap(),
                        members: valset_updated_event.members,
                    };
                    let cosmos_chain_valset =
                        cosmos_gravity::query::get_valset(grpc_client, latest_eth_valset.nonce)
                            .await?;
                    check_if_valsets_differ(cosmos_chain_valset, &latest_eth_valset)?;
                    return Ok(latest_eth_valset);
                }
                Err(e) => error!("Got valset event that we can't parse {}", e),
            }
        }

        end_filter_block = start_filter_block.saturating_sub(1u64.into()); // filter ranges are inclusive, avoid searching same block
    }

    panic!("Could not find the last validator set for contract {}, probably not a valid Gravity contract!", gravity_contract_address)
}

/// This function exists to provide a warning if Cosmos and Ethereum have different validator sets
/// for a given nonce. In the mundane version of this warning the validator sets disagree on sorting order
/// which can happen if some relayer uses an unstable sort, or in a case of a mild griefing attack.
/// The Gravity contract validates signatures in order of highest to lowest power. That way it can exit
/// the loop early once a vote has enough power, if a relayer where to submit things in the reverse order
/// they could grief users of the contract into paying more in gas.
/// The other (and far worse) way a disagreement here could occur is if validators are colluding to steal
/// funds from the Gravity contract and have submitted a highjacking update. If slashing for off Cosmos chain
/// Ethereum signatures is implemented you would put that handler here.
fn check_if_valsets_differ(
    cosmos_valset: Option<Valset>,
    ethereum_valset: &Valset,
) -> Result<(), GravityError> {
    if cosmos_valset.is_none() && ethereum_valset.nonce == 0 {
        // bootstrapping case
        return Ok(());
    } else if cosmos_valset.is_none() {
        error!("Cosmos does not have a valset for nonce {} but that is the one on the Ethereum chain! Possible bridge highjacking!", ethereum_valset.nonce);
        return Ok(());
    }
    let cosmos_valset = cosmos_valset.unwrap();
    if cosmos_valset != *ethereum_valset {
        let mut c_valset = cosmos_valset.members;
        let mut e_valset = ethereum_valset.members.clone();
        c_valset.sort();
        e_valset.sort();
        if c_valset == e_valset {
            warn!(
                "Sorting disagreement between Cosmos and Ethereum on Valset nonce {}",
                ethereum_valset.nonce
            );
        } else {
            error!("Validator sets for nonce {} Cosmos and Ethereum differ. Possible bridge highjacking!", ethereum_valset.nonce)
        }
    }

    Ok(())
}
