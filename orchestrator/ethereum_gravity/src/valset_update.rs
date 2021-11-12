use crate::utils::{EthSignerMiddleware, get_valset_nonce, GasCost};
use ethers::contract::builders::ContractCall;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_utils::types::*;
use gravity_utils::{error::GravityError, gravity::*, message_signatures::encode_valset_confirm_hashed};
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
    let contract_call = build_valset_update_contract_call(
        new_valset, old_valset, confirms, gravity_contract_address, gravity_id, eth_client
    ).await?;
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
    let our_eth_address = eth_client.address();
    let our_balance = eth_client.get_balance(our_eth_address, None).await?;
    let gas_limit = min((u64::MAX - 1).into(), our_balance);
    let gas_price = eth_client.get_gas_price().await?;

    let contract_call = build_valset_update_contract_call(
        new_valset, old_valset, confirms, gravity_contract_address, gravity_id, eth_client
    ).await?;
    let contract_call = contract_call.gas(gas_limit).gas_price(gas_price);

    // TODO(bolten): estimate gas only takes a transaction request, and doesn't
    // care if a specific block is passed as a parameter when creating the contract
    // call...the old code set the nonce manually...I think that the value that
    // ContractCall will put in the TransactionRequest is by default the next
    // available nonce, and the only way to change it would be to reach directly into
    // the object as the ContractCall has no method for setting the nonce...is the
    // default behavior acceptable?
    //
    //let our_nonce = eth_client.get_transaction_count(our_eth_address, None).await?;
    //contract_call.tx.set_nonce(our_nonce);

    Ok(GasCost {
        gas: contract_call.estimate_gas().await?,
        gas_price
    })
}

pub async fn build_valset_update_contract_call(
    new_valset: Valset,
    old_valset: Valset,
    confirms: &[ValsetConfirmResponse],
    gravity_contract_address: EthAddress,
    gravity_id: String,
    eth_client: EthClient,
) -> Result<ContractCall<EthSignerMiddleware, ()>, GravityError> {
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
        .from(eth_address)
        .value(0u8.into()))
}
