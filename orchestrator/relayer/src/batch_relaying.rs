use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity::ethereum::one_eth_f32;
use gravity::ethereum::submit_batch::estimate_tx_batch_cost;
use gravity::ethereum::{
    submit_batch::send_eth_transaction_batch, types::EthClient, utils::get_tx_batch_nonce,
};
use gravity::query::get_latest_transaction_batches;
use gravity::query::get_transaction_batch_signatures;
use gravity::utils::ethereum::downcast_to_f32;
use gravity::utils::message_signatures::encode_tx_batch_confirm_hashed;
use gravity::utils::types::{BatchConfirmResponse, TransactionBatch, Valset};
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use std::collections::HashMap;
use std::time::Duration;
use tonic::transport::Channel;

#[derive(Debug, Clone)]
struct SubmittableBatch {
    batch: TransactionBatch,
    sigs: Vec<BatchConfirmResponse>,
}

/// This function relays batches from Cosmos to Ethereum. First we request
/// the latest transaction batches, which is a list of the latest 100 batches
/// of all types. From there we determine which batches are valid to submit as
/// far as signatures and then make requests to Ethereum to determine which are
/// valid to submit given the current chain state. From there we simulate a submission
/// and if that succeeds and we like the gas cost we complete the relaying process and
/// actually submit the data to Ethereum
#[allow(clippy::too_many_arguments)]
pub async fn relay_batches(
    // the validator set currently in the contract on Ethereum
    current_valset: Valset,
    eth_client: EthClient,
    grpc_client: &mut GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    gravity_id: String,
    timeout: Duration,
    eth_gas_price_multiplier: f32,
    eth_gas_multiplier: f32,
) {
    let possible_batches =
        get_batches_and_signatures(current_valset.clone(), grpc_client, gravity_id.clone()).await;

    debug!("possible batches {:?}", possible_batches);

    submit_batches(
        current_valset,
        eth_client.clone(),
        gravity_contract_address,
        gravity_id,
        timeout,
        eth_gas_price_multiplier,
        eth_gas_multiplier,
        possible_batches,
    )
    .await;
}

/// This function retrieves the latest batches from the Cosmos module and then
/// iterates through the signatures for each batch, determining if they are ready
/// to submit. It is possible for a batch to not have valid signatures for two reasons
/// one is that not enough signatures have been collected yet from the validators two is
/// that the batch is old enough that the signatures do not reflect the current validator
/// set on Ethereum. In both the later and the former case the correct solution is to wait
/// through timeouts, new signatures, or a later valid batch being submitted old batches will
/// always be resolved.
async fn get_batches_and_signatures(
    current_valset: Valset,
    grpc_client: &mut GravityQueryClient<Channel>,
    gravity_id: String,
) -> HashMap<EthAddress, Vec<SubmittableBatch>> {
    let latest_batches = if let Ok(lb) = get_latest_transaction_batches(grpc_client).await {
        lb
    } else {
        return HashMap::new();
    };
    debug!("Latest batches {:?}", latest_batches);

    let mut possible_batches = HashMap::new();
    for batch in latest_batches {
        let sigs =
            get_transaction_batch_signatures(grpc_client, batch.nonce, batch.token_contract).await;
        debug!("Got sigs {:?}", sigs);
        if let Ok(sigs) = sigs {
            // this checks that the signatures for the batch are actually possible to submit to the chain
            let hash = encode_tx_batch_confirm_hashed(gravity_id.clone(), batch.clone());
            if current_valset.order_sigs(&hash, &sigs).is_ok() {
                // we've found a valid batch, add it to the list for it's token type
                possible_batches
                    .entry(batch.token_contract)
                    .or_insert_with(Vec::new);

                let list = possible_batches.get_mut(&batch.token_contract).unwrap();
                list.push(SubmittableBatch { batch, sigs });
            } else {
                warn!(
                    "Batch {}/{} can not be submitted yet, waiting for more signatures",
                    batch.token_contract, batch.nonce
                );
            }
        } else {
            error!(
                "could not get signatures for {}:{} with {:?}",
                batch.token_contract, batch.nonce, sigs
            );
        }
    }

    possible_batches
}

/// Attempts to submit batches with valid signatures, checking the state
/// of the Ethereum chain to ensure that it is valid to submit a given batch
/// more specifically that the correctly signed batch has not timed out or already
/// been submitted. The goal of this function is to submit batches in chronological order
/// of their creation, submitting batches newest first will invalidate old batches and is
/// less efficient if those old batches are profitable.
/// This function estimates the cost of submitting a batch before actually submitting it
/// to Ethereum, if it is determined that the ETH cost to submit is too high the batch will
/// be skipped and a later, more profitable, batch may be submitted.
/// Keep in mind that many other relayers are making this same computation and some may have
/// different standards for their profit margin, therefore there may be a race not only to
/// submit individual batches but also batches in different orders
#[allow(clippy::too_many_arguments)]
async fn submit_batches(
    current_valset: Valset,
    eth_client: EthClient,
    gravity_contract_address: EthAddress,
    gravity_id: String,
    timeout: Duration,
    eth_gas_price_multiplier: f32,
    eth_gas_multiplier: f32,
    possible_batches: HashMap<EthAddress, Vec<SubmittableBatch>>,
) {
    let ethereum_block_height = if let Ok(bn) = eth_client.get_block_number().await {
        bn
    } else {
        error!("Failed to get eth block height, is your eth node working?");
        return;
    };

    // requests data from Ethereum only once per token type, this is valid because we are
    // iterating from oldest to newest, so submitting a batch earlier in the loop won't
    // ever invalidate submitting a batch later in the loop. Another relayer could always
    // do that though.
    for (token_type, possible_batches) in possible_batches {
        let erc20_contract = token_type;
        let latest_ethereum_batch =
            get_tx_batch_nonce(gravity_contract_address, erc20_contract, eth_client.clone()).await;
        if latest_ethereum_batch.is_err() {
            error!(
                "Failed to get latest Ethereum batch with {:?}",
                latest_ethereum_batch
            );
            return;
        }
        let latest_ethereum_batch = latest_ethereum_batch.unwrap();

        for batch in possible_batches {
            let oldest_signed_batch = batch.batch;
            let oldest_signatures = batch.sigs;

            if oldest_signed_batch.batch_timeout < ethereum_block_height.as_u64() {
                warn!(
                    "Batch {}/{} has timed out and can not be submitted",
                    oldest_signed_batch.nonce, oldest_signed_batch.token_contract
                );
                continue;
            }

            let latest_cosmos_batch_nonce = oldest_signed_batch.clone().nonce;
            if latest_cosmos_batch_nonce > latest_ethereum_batch {
                let cost = estimate_tx_batch_cost(
                    current_valset.clone(),
                    oldest_signed_batch.clone(),
                    &oldest_signatures,
                    gravity_contract_address,
                    gravity_id.clone(),
                    eth_client.clone(),
                )
                .await;

                if cost.is_err() {
                    error!("Batch cost estimate failed with {:?}", cost);
                    continue;
                }

                let mut cost = cost.unwrap();
                let total_cost = downcast_to_f32(cost.get_total());
                if total_cost.is_none() {
                    error!(
                        "Total gas cost greater than f32 max, skipping batch submission: {}",
                        oldest_signed_batch.nonce
                    );
                    continue;
                }
                let total_cost = total_cost.unwrap();
                let gas_price_as_f32 = downcast_to_f32(cost.gas_price).unwrap(); // if the total cost isn't greater, this isn't
                let gas_as_f32 = downcast_to_f32(cost.gas).unwrap(); // same as above re: total cost

                info!(
                    "We have detected latest batch {} but latest on Ethereum is {} This batch is estimated to cost {} Gas / {:.4} ETH to submit",
                    latest_cosmos_batch_nonce,
                    latest_ethereum_batch,
                    cost.gas_price.clone(),
                    total_cost / one_eth_f32()
                );

                cost.gas_price = ((gas_price_as_f32 * eth_gas_price_multiplier) as u128).into();
                cost.gas = ((gas_as_f32 * eth_gas_multiplier) as u128).into();

                let res = send_eth_transaction_batch(
                    current_valset.clone(),
                    oldest_signed_batch,
                    &oldest_signatures,
                    timeout,
                    gravity_contract_address,
                    gravity_id.clone(),
                    cost,
                    eth_client.clone(),
                )
                .await;

                if res.is_err() {
                    warn!("Batch submission failed with {:?}", res);
                }
            }
        }
    }
}
