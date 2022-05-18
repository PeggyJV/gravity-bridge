use crate::main_loop::LOOP_SPEED;
use cosmos_gravity::query::{get_latest_logic_calls, get_logic_call_signatures};
use ethereum_gravity::logic_call::LogicCallSkips;
use ethereum_gravity::one_eth_f32;
use ethereum_gravity::utils::handle_contract_error;
use ethereum_gravity::{
    logic_call::send_eth_logic_call, types::EthClient, utils::get_logic_call_nonce,
};
use ethers::types::Address as EthAddress;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::ethereum::{bytes_to_hex_str, downcast_to_f32};
use gravity_utils::types::{LogicCallConfirmResponse, Valset};
use gravity_utils::{message_signatures::encode_logic_call_confirm_hashed, types::LogicCall};
use std::time::Duration;
use tonic::transport::Channel;

#[allow(clippy::too_many_arguments)]
pub async fn relay_logic_calls(
    // the validator set currently in the contract on Ethereum
    current_valset: Valset,
    eth_client: EthClient,
    grpc_client: &mut GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    gravity_id: String,
    timeout: Duration,
    eth_gas_price_multiplier: f32,
    logic_call_skips: &mut LogicCallSkips,
) {
    let latest_calls = match get_latest_logic_calls(grpc_client).await {
        Ok(calls) => {
            debug!("Latest Logic calls {:?}", calls);
            calls
        }
        Err(err) => {
            error!("Error while retrieving latest logic calls: {:?}", err);
            return;
        }
    };
    let mut oldest_signed_call: Option<LogicCall> = None;
    let mut oldest_signatures: Option<Vec<LogicCallConfirmResponse>> = None;
    for call in latest_calls {
        if logic_call_skips.permanently_skipped(&call) {
            info!("LogicCall {}/{} permanently skipped until oracle confirms or on-chain timeout after eth height {}",
                bytes_to_hex_str(&call.invalidation_id), call.invalidation_nonce, call.timeout
            );
            continue;
        }

        let skips_left: u64 = logic_call_skips.skips_left(&call).into();
        if skips_left > 0 {
            warn!(
                "Skipping LogicCall {}/{} with eth timeout {}, estimated next retry after minimum of {} seconds",
                bytes_to_hex_str(&call.invalidation_id), call.invalidation_nonce, call.timeout, skips_left * LOOP_SPEED.as_secs()
            );
            logic_call_skips.skip(&call);
            continue;
        }

        let sigs = get_logic_call_signatures(
            grpc_client,
            call.invalidation_id.clone(),
            call.invalidation_nonce,
        )
        .await;
        debug!("Got sigs {:?}", sigs);
        if let Ok(sigs) = sigs {
            let hash = encode_logic_call_confirm_hashed(gravity_id.clone(), call.clone());
            // this checks that the signatures for the logic call are actually possible to submit to the chain
            if current_valset.order_sigs(&hash, &sigs).is_ok() {
                oldest_signed_call = Some(call);
                oldest_signatures = Some(sigs);
            } else {
                warn!(
                    "LogicCall {}/{} can not be submitted yet, waiting for more signatures",
                    bytes_to_hex_str(&call.invalidation_id),
                    call.invalidation_nonce
                );
            }
        } else {
            error!(
                "could not get signatures for {}/{} with {:?}",
                bytes_to_hex_str(&call.invalidation_id),
                call.invalidation_nonce,
                sigs
            );
        }
    }
    if oldest_signed_call.is_none() {
        debug!("Could not find Call with signatures! exiting");
        return;
    }
    let oldest_signed_call = oldest_signed_call.unwrap();
    let oldest_signatures = oldest_signatures.unwrap();

    let latest_ethereum_call = get_logic_call_nonce(
        gravity_contract_address,
        oldest_signed_call.invalidation_id.clone(),
        eth_client.clone(),
    )
    .await;
    if latest_ethereum_call.is_err() {
        error!(
            "Failed to get latest Ethereum LogicCall with {:?}",
            latest_ethereum_call
        );
        return;
    }
    let latest_ethereum_call = latest_ethereum_call.unwrap();
    let latest_cosmos_call_nonce = oldest_signed_call.clone().invalidation_nonce;
    if latest_cosmos_call_nonce > latest_ethereum_call {
        let cost = ethereum_gravity::logic_call::estimate_logic_call_cost(
            current_valset.clone(),
            oldest_signed_call.clone(),
            &oldest_signatures,
            gravity_contract_address,
            gravity_id.clone(),
            eth_client.clone(),
        )
        .await;

        if cost.is_err() {
            warn!("LogicCall cost estimate failed");
            let should_permanently_skip = handle_contract_error(cost.unwrap_err());
            if should_permanently_skip {
                logic_call_skips.skip_permanently(&oldest_signed_call);
            } else {
                logic_call_skips.skip(&oldest_signed_call);
            }
            return;
        }

        let mut cost = cost.unwrap();
        let total_cost = downcast_to_f32(cost.get_total());
        if total_cost.is_none() {
            error!(
                "Total gas cost greater than f32 max, skipping logic call submission: {}",
                oldest_signed_call.invalidation_nonce
            );
            logic_call_skips.skip(&oldest_signed_call);
            return;
        }
        let total_cost = total_cost.unwrap();
        let gas_price_as_f32 = downcast_to_f32(cost.gas_price).unwrap(); // if the total cost isn't greater, this isn't

        info!(
            "We have detected latest LogicCall {} but latest on Ethereum is {} This LogicCall is estimated to cost {} Gas / {:.4} ETH to submit",
            latest_cosmos_call_nonce,
            latest_ethereum_call,
            cost.gas_price.clone(),
            total_cost / one_eth_f32(),
        );

        cost.gas_price = ((gas_price_as_f32 * eth_gas_price_multiplier) as u128).into();

        let res = send_eth_logic_call(
            current_valset,
            oldest_signed_call.clone(),
            &oldest_signatures,
            timeout,
            gravity_contract_address,
            gravity_id.clone(),
            cost,
            eth_client.clone(),
            logic_call_skips,
        )
        .await;

        if res.is_err() {
            warn!("LogicCall submission failed");
            let should_permanently_skip = handle_contract_error(res.unwrap_err());
            if should_permanently_skip {
                logic_call_skips.skip_permanently(&oldest_signed_call);
            } else {
                logic_call_skips.skip(&oldest_signed_call);
            }
        }
    }
}
