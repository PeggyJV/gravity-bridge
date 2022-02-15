//! The Gravity deployERC20 endpoint deploys an ERC20 contract representing a Cosmos asset onto the Ethereum blockchain
//! the event for this deployment is then ferried over to Cosmos where the validators will accept the ERC20 contract address
//! as the representation of this asset on Ethereum

use crate::{types::EthClient, utils::get_gas_price};
use ethers::prelude::*;
use gravity_abi::gravity::*;
use gravity_utils::{
    error::GravityError,
    ethereum::{downcast_to_f64, format_eth_hash},
};
use std::{result::Result, time::Duration};

/// Calls the Gravity ethereum contract to deploy the ERC20 representation of the given Cosmos asset
/// denom. If an existing contract is already deployed representing this asset this call will cost
/// Gas but not actually do anything. Returns the new contract address or an error
#[allow(clippy::too_many_arguments)]
pub async fn deploy_erc20(
    cosmos_denom: String,
    erc20_name: String,
    erc20_symbol: String,
    decimals: u8,
    gravity_contract: Address,
    wait_timeout: Option<Duration>,
    gas_multiplier: f64,
    eth_client: EthClient,
) -> Result<TxHash, GravityError> {
    let contract_call = Gravity::new(gravity_contract, eth_client.clone()).deploy_erc20(
        cosmos_denom,
        erc20_name,
        erc20_symbol.clone(),
        decimals,
    );
    let gas_price = get_gas_price(eth_client.clone()).await?;
    let gas = contract_call.estimate_gas().await?;
    let gas_as_f64 = downcast_to_f64(gas);
    if gas_as_f64.is_none() {
        return Err(GravityError::GravityContractError(format!(
            "Gas estimate too large to downcast to f64: {}",
            gas
        )));
    }
    let gas = (gas_as_f64.unwrap() * gas_multiplier) as u128;

    // TODO(bolten): it seems like a bug in ethers will replace manually set gas limits with
    // a gas estimate if no access list is defined for EIP1559 transactions, so we're forcing a
    // legacy transaction here to allow for the multiplier to take effect
    let contract_call = contract_call.gas_price(gas_price).gas(gas).legacy();

    let pending_tx = contract_call.send().await?;
    let tx_hash = *pending_tx;
    info!("Deploying ERC-20 with tx hash {}", format_eth_hash(tx_hash));
    // TODO(bolten): ethers interval default is 7s, this mirrors what web30 was doing, should we adjust?
    // additionally we are mirroring only waiting for 1 confirmation by leaving that as default
    let pending_tx = pending_tx.interval(Duration::from_secs(1));
    let potential_error = GravityError::GravityContractError(format!(
        "Did not receive transaction receipt when deploying ERC-20 {}: {}",
        erc20_symbol,
        format_eth_hash(tx_hash)
    ));

    if let Some(timeout) = wait_timeout {
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
