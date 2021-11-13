use crate::utils::{EthSignerMiddleware, get_valset_nonce, GasCost, set_contract_call_gas_for_estimate};
use ethers::contract::builders::ContractCall;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_abi::gravity::*;
use gravity_utils::{error::GravityError, message_signatures::encode_valset_confirm_hashed, types::*};
use std::{cmp::min, time::Duration};

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
    eth_client: EthClient,
) -> Result<(), GravityError> {
    let old_nonce = old_valset.nonce;
    let new_nonce = new_valset.nonce;
    assert!(new_nonce > old_nonce);
    let eth_address = eth_client.address();
    info!(
        "Ordering signatures and submitting validator set {} -> {} update to Ethereum",
        old_nonce, new_nonce
    );
    let before_nonce = get_valset_nonce(gravity_contract_address, eth_address, eth_client).await?;
    if before_nonce != old_nonce {
        info!(
            "Someone else updated the valset to {}, exiting early",
            before_nonce
        );
        return Ok(());
    }

    let contract_call = build_valset_update_contract_call(
        new_valset, old_valset, confirms, gravity_contract_address, gravity_id, eth_client
    )?;
    let pending_tx = contract_call.send().await?;
    info!("Sent valset update with txid {:#066x}", pending_tx);
    // TODO(bolten): ethers interval default is 7s, this mirrors what web30 was doing, should we adjust?
    // additionally we are mirroring only waiting for 1 confirmation by leaving that as default
    pending_tx.interval(Duration::from_secs(1));

    if let Err(tx_error) = tokio::time::timeout(timeout, async { pending_tx.await? }).await {
        return Err(tx_error);
    };

    let last_nonce = get_valset_nonce(gravity_contract_address, eth_address, eth_client).await?;
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
        new_valset, old_valset, confirms, gravity_contract_address, gravity_id, eth_client
    )?;
    let contract_call =
        set_contract_call_gas_for_estimate(contract_call, eth_client).await?;

    Ok(GasCost {
        gas: contract_call.estimate_gas().await?,
        gas_price
    })
}

pub fn build_valset_update_contract_call(
    new_valset: Valset,
    old_valset: Valset,
    confirms: &[ValsetConfirmResponse],
    gravity_contract_address: EthAddress,
    gravity_id: String,
    eth_client: EthClient,
) -> Result<ContractCall<EthSignerMiddleware, ()>, GravityError> {
    let (old_addresses, old_powers) = old_valset.filter_empty_addresses();
    let (new_addresses, new_powers) = new_valset.filter_empty_addresses();

    // remember the signatures are over the new valset and therefore this is the value we must encode
    // the old valset exists only as a hash in the ethereum store
    let hash = encode_valset_confirm_hashed(gravity_id, new_valset);
    // we need to use the old valset here because our signatures need to match the current
    // members of the validator set in the contract.
    let sig_data = old_valset.order_sigs(&hash, confirms)?;
    let sig_arrays = to_arrays(sig_data);

    let contract = Gravity::new(gravity_contract_address, eth_client);
    Ok(contract.update_valset(
        new_addresses, new_powers, new_nonce.into(),
        old_addresses, old_powers, old_nonce.into(),
        sig_arrays.v, sig_arrays.r, sig_arrays.s)
        .from(eth_client.address())
        .value(0u8.into()))
}
