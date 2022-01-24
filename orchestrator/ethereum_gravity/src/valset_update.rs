use crate::{
    types::{EthClient, EthSignerMiddleware},
    utils::{get_gas_price, get_valset_nonce, GasCost},
};
use ethers::contract::builders::ContractCall;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_abi::gravity::*;
use gravity_utils::{
    error::GravityError, message_signatures::encode_valset_confirm_hashed, types::*,
};
use std::{result::Result, time::Duration};

/// this function generates an appropriate Ethereum transaction
/// to submit the provided validator set and signatures.
#[allow(clippy::too_many_arguments)]
pub async fn send_eth_valset_update(
    new_valset: Valset,
    old_valset: Valset,
    confirms: &[ValsetConfirmResponse],
    timeout: Duration,
    gravity_contract_address: EthAddress,
    gravity_id: String,
    gas_cost: GasCost,
    eth_client: EthClient,
) -> Result<(), GravityError> {
    let old_nonce = old_valset.nonce;
    let new_nonce = new_valset.nonce;

    info!(
        "Ordering signatures and submitting validator set {} -> {} update to Ethereum",
        old_nonce, new_nonce
    );
    let before_nonce = get_valset_nonce(gravity_contract_address, eth_client.clone()).await?;
    if before_nonce != old_nonce {
        info!(
            "Someone else updated the valset to {}, exiting early",
            before_nonce
        );
        return Ok(());
    }

    let contract_call = build_valset_update_contract_call(
        &new_valset,
        &old_valset,
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
    info!("Sent valset update with txid {}", tx_hash);
    // TODO(bolten): ethers interval default is 7s, this mirrors what web30 was doing, should we adjust?
    // additionally we are mirroring only waiting for 1 confirmation by leaving that as default
    let pending_tx = pending_tx.interval(Duration::from_secs(1));

    match tokio::time::timeout(timeout, pending_tx).await?? {
        Some(_) => (),
        None => error!(
            "Did not receive transaction receipt when sending valset update: {}",
            tx_hash
        ),
    }

    let last_nonce = get_valset_nonce(gravity_contract_address, eth_client.clone()).await?;
    if last_nonce != new_nonce {
        error!(
            "Current nonce is {} expected to update to nonce {}",
            last_nonce, new_nonce
        );
    } else {
        info!(
            "Successfully updated Valset with new Nonce {:?}",
            last_nonce
        );
    }

    Ok(())
}

/// Returns the cost in Eth of sending this valset update
pub async fn estimate_valset_cost(
    new_valset: &Valset,
    old_valset: &Valset,
    confirms: &[ValsetConfirmResponse],
    gravity_contract_address: EthAddress,
    gravity_id: String,
    eth_client: EthClient,
) -> Result<GasCost, GravityError> {
    let contract_call = build_valset_update_contract_call(
        new_valset,
        old_valset,
        confirms,
        gravity_contract_address,
        gravity_id,
        eth_client.clone(),
    )?;

    Ok(GasCost {
        gas: contract_call.estimate_gas().await?,
        gas_price: get_gas_price(eth_client.clone()).await?,
    })
}

pub fn build_valset_update_contract_call(
    new_valset: &Valset,
    old_valset: &Valset,
    confirms: &[ValsetConfirmResponse],
    gravity_contract_address: EthAddress,
    gravity_id: String,
    eth_client: EthClient,
) -> Result<ContractCall<EthSignerMiddleware, ()>, GravityError> {
    let (old_addresses, old_powers) = old_valset.filter_empty_addresses();
    let (new_addresses, new_powers) = new_valset.filter_empty_addresses();
    let old_powers: Vec<U256> = old_powers.iter().map(|power| (*power).into()).collect();
    let new_powers: Vec<U256> = new_powers.iter().map(|power| (*power).into()).collect();

    // remember the signatures are over the new valset and therefore this is the value we must encode
    // the old valset exists only as a hash in the ethereum store
    let hash = encode_valset_confirm_hashed(gravity_id, new_valset.clone());
    // we need to use the old valset here because our signatures need to match the current
    // members of the validator set in the contract.
    let sig_data = old_valset.order_sigs(&hash, confirms)?;

    let contract = Gravity::new(gravity_contract_address, eth_client.clone());
    Ok(contract
        .update_valset(
            ValsetArgs {
                validators: new_addresses,
                powers: new_powers,
                valset_nonce: new_valset.nonce.into(),
                reward_amount: U256::zero(),
                reward_token: H160::zero(),
            },
            ValsetArgs {
                validators: old_addresses,
                powers: old_powers,
                valset_nonce: old_valset.nonce.into(),
                reward_amount: U256::zero(),
                reward_token: H160::zero(),
            },
            sig_data
                .iter()
                .map(|sig_data| sig_data.to_val_sig())
                .collect(),
        )
        .from(eth_client.address())
        .value(U256::zero()))
}
