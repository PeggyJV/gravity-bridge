//! This module contains code for the validator update lifecycle. Functioning as a way for this validator to observe
//! the state of both chains and perform the required operations.

use std::time::Duration;

use cosmos_gravity::query::get_latest_valset;
use cosmos_gravity::query::{get_all_valset_confirms, get_valset};
use ethereum_gravity::{one_eth_f32, types::EthClient, valset_update::send_eth_valset_update};
use ethers::types::Address as EthAddress;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::{
    ethereum::bytes_to_hex_str, ethereum::downcast_to_f32,
    message_signatures::encode_valset_confirm_hashed, types::Valset,
};
use tonic::transport::Channel;

/// Check the last validator set on Ethereum, if it's lower than our latest validator
/// set then we should package and submit the update as an Ethereum transaction
pub async fn relay_valsets(
    // the validator set currently in the contract on Ethereum
    current_eth_valset: Valset,
    eth_client: EthClient,
    grpc_client: &mut GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    gravity_id: String,
    timeout: Duration,
) {
    // we have to start with the current ethereum valset, we need to know what's currently
    // in the contract in order to determine if a new validator set is valid.
    // For example the contract has set A which contains validators x/y/z the
    // latest valset has set C which has validators z/e/f in order to have enough
    // power we actually need to submit validator set B with validators x/y/e in
    // order to know that we need a set from the history

    // we should determine if we need to relay one
    // to Ethereum for that we will find the latest confirmed valset and compare it to the ethereum chain
    let latest_valset = get_latest_valset(grpc_client).await;
    if latest_valset.is_err() {
        error!("Failed to get latest valset! {:?}", latest_valset);
        return;
    }
    let latest_valset = latest_valset.unwrap();
    if latest_valset.is_none() {
        return;
    }

    let latest_valset = latest_valset.unwrap();

    // we only use the latest valsets endpoint to get a starting point, from there we will iterate
    // backwards until we find the newest validator set that we can submit to the bridge. So if we
    // have sets A-Z and it's possible to submit only A, L, and Q before reaching Z this code will do
    // so.
    let mut latest_nonce = latest_valset.nonce;
    let mut latest_confirmed = None;
    let mut latest_cosmos_valset = None;
    // this is used to display the state of the last validator set to fail signature checks
    let mut last_error = None;
    while latest_nonce > 0 {
        match get_valset(grpc_client, latest_nonce).await {
            Ok(Some(cosmos_valset)) => {
                assert_eq!(cosmos_valset.nonce, latest_nonce);

                match get_all_valset_confirms(grpc_client, latest_nonce).await {
                    Ok(confirms) => {
                        debug!(
                            "Considering cosmos_valset {:?} confirms {:?}",
                            cosmos_valset, confirms
                        );

                        for confirm in confirms.iter() {
                            assert_eq!(cosmos_valset.nonce, confirm.nonce);
                        }

                        let hash =
                            encode_valset_confirm_hashed(gravity_id.clone(), cosmos_valset.clone());

                        // there are two possible encoding problems that could cause the very rare sig failure bug,
                        // one of them is that the hash is incorrect, that's not probable considering that
                        // both Geth and Clarity agree on it. but this lets us check
                        info!("New valset hash {}", bytes_to_hex_str(&hash));

                        // order valset sigs prepares signatures for submission, notice we compare
                        // them to the 'current' set in the bridge, this confirms for us that the validator set
                        // we have here can be submitted to the bridge in it's current state
                        match current_eth_valset.order_sigs(&hash, &confirms) {
                            Ok(_) => {
                                info!("Consideration: looks good");
                                latest_confirmed = Some(confirms);
                                latest_cosmos_valset = Some(cosmos_valset);
                                // once we have the latest validator set we can submit exit
                                break;
                            }
                            Err(err) => {
                                warn!("Consideration: looks bad. {}", err);
                                last_error = Some(err);
                            }
                        }
                    }
                    Err(err) => {
                        error!(
                            "Error getting Cosmos valset confirmations for nonce {}. {:?}",
                            latest_nonce, err
                        );
                    }
                }
            }
            // Bootstrap condition
            Ok(None) => (),
            Err(err) => {
                error!(
                    "Error getting Cosmos valset for nonce {}. {:?}",
                    latest_nonce, err
                );
            }
        }

        latest_nonce -= 1
    }

    debug!("Relaying latest_confirmed {:?}", latest_confirmed);
    debug!("Relaying latest_cosmos_valset {:?}", latest_cosmos_valset);

    if latest_confirmed.is_none() {
        error!("We don't have a latest confirmed valset?");
        return;
    }
    // the latest cosmos validator set that it is possible to submit given the constraints
    // of the validator set currently in the bridge
    let latest_cosmos_valset = latest_cosmos_valset.unwrap();
    // the signatures for the above
    let latest_cosmos_confirmed = latest_confirmed.unwrap();

    // this will print a message indicating the signing state of the latest validator
    // set if the latest available validator set is not the latest one that is possible
    // to submit. AKA if the bridge is behind where it should be
    if latest_nonce > latest_cosmos_valset.nonce && last_error.is_some() {
        warn!("{:?}", last_error)
    }

    let should_relay = latest_cosmos_valset.nonce > current_eth_valset.nonce;
    info!(
        "Considering nonces: latest_cosmos_valset.nonce {} current_eth_valset.nonce {} should_relay {}",
        latest_cosmos_valset.nonce, current_eth_valset.nonce,
        should_relay,
    );

    if should_relay {
        let cost = ethereum_gravity::valset_update::estimate_valset_cost(
            &latest_cosmos_valset,
            &current_eth_valset,
            &latest_cosmos_confirmed,
            gravity_contract_address,
            gravity_id.clone(),
            eth_client.clone(),
        )
        .await;

        if cost.is_err() {
            error!(
                "Valset cost estimate for Nonce {} failed with {:?}",
                latest_cosmos_valset.nonce, cost
            );
            return;
        }
        let cost = cost.unwrap();
        let total_cost = downcast_to_f32(cost.get_total());
        if total_cost.is_none() {
            error!(
                "Total gas cost greater than f32 max, skipping valset submission: {}",
                latest_cosmos_valset.nonce
            );
            return;
        }
        let total_cost = total_cost.unwrap();

        info!(
           "We have detected latest valset {} but latest on Ethereum is {} This valset is estimated to cost {} Gas / {:.4} ETH to submit",
            latest_cosmos_valset.nonce, current_eth_valset.nonce,
            cost.gas_price.clone(),
            total_cost / one_eth_f32()
        );

        let relay_response = send_eth_valset_update(
            latest_cosmos_valset.clone(),
            current_eth_valset.clone(),
            &latest_cosmos_confirmed,
            timeout,
            gravity_contract_address,
            gravity_id,
            cost,
            eth_client.clone(),
        )
        .await;

        info!(
            "relay_response {:?} (latest_cosmos_valset.nonce {} current_eth_valset.nonce {})",
            relay_response, latest_cosmos_valset.nonce, current_eth_valset.nonce,
        );
    }
}
