use crate::{
    types::{EthClient, EthSignerMiddleware},
    utils::{get_logic_call_nonce, get_send_transaction_gas_price, GasCost},
};
use ethers::contract::builders::ContractCall;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_abi::gravity::*;
use gravity_utils::ethereum::{bytes_to_hex_str, vec_u8_to_fixed_32};
use gravity_utils::types::*;
use gravity_utils::{error::GravityError, message_signatures::encode_logic_call_confirm_hashed};
use std::{result::Result, time::Duration};

/// this function generates an appropriate Ethereum transaction
/// to submit the provided logic call
#[allow(clippy::too_many_arguments)]
pub async fn send_eth_logic_call(
    current_valset: Valset,
    call: LogicCall,
    confirms: &[LogicCallConfirmResponse],
    timeout: Duration,
    gravity_contract_address: EthAddress,
    gravity_id: String,
    gas_cost: GasCost,
    eth_client: EthClient,
) -> Result<(), GravityError> {
    let new_call_nonce = call.invalidation_nonce;
    info!(
        "Ordering signatures and submitting LogicCall {}:{} to Ethereum",
        bytes_to_hex_str(&call.invalidation_id),
        new_call_nonce
    );
    trace!("Call {:?}", call);

    let before_nonce = get_logic_call_nonce(
        gravity_contract_address,
        call.invalidation_id.clone(),
        eth_client.clone(),
    )
    .await?;

    let current_block_height = eth_client.get_block_number().await?;
    if before_nonce >= new_call_nonce {
        info!(
            "Someone else updated the LogicCall to {}, exiting early",
            before_nonce
        );
        return Ok(());
    } else if current_block_height > call.timeout.into() {
        info!(
            "This LogicCall is timed out. timeout block: {} current block: {}, exiting early",
            current_block_height, call.timeout
        );
        return Ok(());
    }

    let contract_call = build_send_logic_call_contract_call(
        current_valset,
        &call,
        confirms,
        gravity_contract_address,
        gravity_id,
        eth_client.clone(),
    )?;

    let contract_call = contract_call
        .gas(gas_cost.gas)
        .gas_price(gas_cost.gas_price);

    let pending_tx = contract_call.send().await?;
    let tx_hash = *pending_tx;
    info!("Sent logic call with txid {}", tx_hash);
    // TODO(bolten): ethers interval default is 7s, this mirrors what web30 was doing, should we adjust?
    // additionally we are mirroring only waiting for 1 confirmation by leaving that as default
    let pending_tx = pending_tx.interval(Duration::from_secs(1));

    match tokio::time::timeout(timeout, pending_tx).await?? {
        Some(_) => (),
        None => error!(
            "Did not receive transaction receipt when submitting batch: {}",
            tx_hash
        ),
    }

    let last_nonce = get_logic_call_nonce(
        gravity_contract_address,
        call.invalidation_id,
        eth_client.clone(),
    )
    .await?;

    if last_nonce != new_call_nonce {
        error!(
            "Current nonce is {} expected to update to nonce {}",
            last_nonce, new_call_nonce
        );
    } else {
        info!(
            "Successfully updated LogicCall with new Nonce {:?}",
            last_nonce
        );
    }
    Ok(())
}

/// Returns the cost in Eth of sending this batch
pub async fn estimate_logic_call_cost(
    current_valset: Valset,
    call: LogicCall,
    confirms: &[LogicCallConfirmResponse],
    gravity_contract_address: EthAddress,
    gravity_id: String,
    eth_client: EthClient,
) -> Result<GasCost, GravityError> {
    let contract_call = build_send_logic_call_contract_call(
        current_valset,
        &call,
        confirms,
        gravity_contract_address,
        gravity_id,
        eth_client.clone(),
    )?;

    Ok(GasCost {
        gas: contract_call.estimate_gas().await?,
        gas_price: get_send_transaction_gas_price(eth_client.clone()).await?,
    })
}

pub fn build_send_logic_call_contract_call(
    current_valset: Valset,
    call: &LogicCall,
    confirms: &[LogicCallConfirmResponse],
    gravity_contract_address: EthAddress,
    gravity_id: String,
    eth_client: EthClient,
) -> Result<ContractCall<EthSignerMiddleware, ()>, GravityError> {
    let (current_addresses, current_powers) = current_valset.filter_empty_addresses();
    let current_powers: Vec<U256> = current_powers.iter().map(|power| (*power).into()).collect();
    let current_valset_nonce = current_valset.nonce;
    let hash = encode_logic_call_confirm_hashed(gravity_id, call.clone());
    let sig_data = current_valset.order_sigs(&hash, confirms)?;
    let sig_arrays = to_arrays(sig_data);

    let transfer_amounts = call
        .transfers
        .iter()
        .map(|transfer| transfer.amount)
        .collect();
    let transfer_token_contracts = call
        .transfers
        .iter()
        .map(|transfer| transfer.token_contract_address)
        .collect();
    let fee_amounts = call.fees.iter().map(|fee| fee.amount).collect();
    let fee_token_contracts = call
        .fees
        .iter()
        .map(|fee| fee.token_contract_address)
        .collect();
    let invalidation_id = vec_u8_to_fixed_32(call.invalidation_id.clone())?;

    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .submit_logic_call(
            current_addresses,
            current_powers,
            current_valset_nonce.into(),
            sig_arrays.v,
            sig_arrays.r,
            sig_arrays.s,
            LogicCallArgs {
                transfer_amounts,
                transfer_token_contracts,
                fee_amounts,
                fee_token_contracts,
                logic_contract_address: call.logic_contract_address,
                payload: call.payload.clone().into(),
                time_out: call.timeout.into(),
                invalidation_id,
                invalidation_nonce: call.invalidation_nonce.into(),
            },
        )
        .from(eth_client.address())
        .value(U256::zero());

    Ok(contract_call)
}
