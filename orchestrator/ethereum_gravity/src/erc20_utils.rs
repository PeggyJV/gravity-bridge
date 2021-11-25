use crate::types::EthClient;
use ethers::prelude::*;
use gravity_abi::erc20::ERC20;
use gravity_utils::error::GravityError;
use std::{result::Result, time::Duration};

/// Checks if any given contract is approved to spend money from any given erc20 contract
/// using any given address. What exactly this does can be hard to grok, essentially when
/// you want contract A to be able to spend your erc20 contract funds you need to call 'approve'
/// on the ERC20 contract with your own address and A's address so that in the future when you call
/// contract A it can manipulate your ERC20 balances. This function checks if that has already been done
/// and that the allowed amount is greater than the provided allowance threshold.
pub async fn check_erc20_approved(
    erc20: Address,
    target_contract: Address,
    address: Address,
    allowance_threshold: U256,
    eth_client: EthClient,
) -> Result<bool, GravityError> {
    let erc20_contract = ERC20::new(erc20, eth_client.clone());
    let contract_call = erc20_contract.allowance(address, target_contract);
    let allowance = contract_call.call().await?;

    Ok(allowance > allowance_threshold)
}

/// Approves a given contract to spend erc20 funds from the given address from the erc20 contract provided.
/// What exactly this does can be hard to grok, essentially when you want contract A to be able to spend
/// your erc20 contract funds you need to call 'approve' on the ERC20 contract with your own address and A's
/// address so that in the future when you call contract A it can manipulate your ERC20 balances.
/// This function performs that action and waits for it to complete for up to Timeout duration
pub async fn approve_erc20_transfers(
    erc20: Address,
    target_contract: Address,
    timeout_option: Option<Duration>,
    eth_client: EthClient,
) -> Result<TxHash, GravityError> {
    let erc20_contract = ERC20::new(erc20, eth_client.clone());
    let contract_call = erc20_contract.approve(target_contract, U256::MAX);

    let pending_tx = contract_call.send().await?;
    let tx_hash = *pending_tx;
    info!("Approving ERC-20 {} with txid {}", erc20, tx_hash);
    // TODO(bolten): ethers interval default is 7s, this mirrors what web30 was doing, should we adjust?
    // additionally we are mirroring only waiting for 1 confirmation by leaving that as default
    let pending_tx = pending_tx.interval(Duration::from_secs(1));
    let potential_error = GravityError::GravityContractError(format!(
        "Did not receive transaction receipt when approving ERC-20 {}: {}",
        erc20, tx_hash
    ));

    if let Some(timeout) = timeout_option {
        match tokio::time::timeout(timeout, pending_tx).await?? {
            Some(receipt) => Ok(receipt.transaction_hash),
            None => Err(potential_error),
        }
    } else {
        match pending_tx.await? {
            Some(receipt) => Ok(receipt.transaction_hash),
            None => Err(potential_error),
        }
    }
}

pub async fn get_erc20_balance(
    erc20: Address,
    address: Address,
    eth_client: EthClient,
) -> Result<U256, GravityError> {
    let erc20_contract = ERC20::new(erc20, eth_client.clone());
    let contract_call = erc20_contract.balance_of(address);

    Ok(contract_call.call().await?)
}
