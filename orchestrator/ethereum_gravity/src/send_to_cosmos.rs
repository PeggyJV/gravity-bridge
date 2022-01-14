//! Helper functions for sending tokens to Cosmos

use crate::{
    erc20_utils::{approve_erc20_transfers, check_erc20_approved},
    types::EthClient,
};
use deep_space::address::Address as CosmosAddress;
use ethers::prelude::*;
use gravity_abi::gravity::*;
use gravity_utils::error::GravityError;
use std::{result::Result, time::Duration};

const SEND_TO_COSMOS_GAS_LIMIT: u128 = 100_000;

#[allow(clippy::too_many_arguments)]
pub async fn send_to_cosmos(
    erc20: Address,
    gravity_contract: Address,
    amount: U256,
    cosmos_destination: CosmosAddress,
    wait_timeout: Option<Duration>,
    eth_client: EthClient,
) -> Result<TxHash, GravityError> {
    // TODO(bolten): this value is ported from web30, does it match our expectations?
    // Check if the allowance remaining is greater than half of a U256 - it's as good
    // a test as any.
    let allowance_threshold = U256::MAX.div_mod(2u32.into()).0;
    let approved = check_erc20_approved(
        erc20,
        gravity_contract,
        eth_client.address(),
        allowance_threshold,
        eth_client.clone(),
    )
    .await?;
    if !approved {
        let txid =
            approve_erc20_transfers(erc20, gravity_contract, wait_timeout, eth_client.clone())
                .await?;
        trace!("ERC-20 approval for {} finished with txid {}", erc20, txid);
    }

    // This code deals with some specifics of Ethereum byte encoding, Ethereum is BigEndian
    // so small values like addresses that don't take up the full length of the byte vector
    // are pushed up to the top. This duplicates the way Ethereum encodes it's own addresses
    // as closely as possible.
    let mut cosmos_dest_address_bytes = cosmos_destination.as_bytes().to_vec();
    while cosmos_dest_address_bytes.len() < 32 {
        cosmos_dest_address_bytes.insert(0, 0u8);
    }
    // TODO(bolten): have to convert back from what was done above for the contract call,
    // there's probably a cleaner way to do this
    let mut cosmos_dest_address_bytes_slice: [u8; 32] = Default::default();
    cosmos_dest_address_bytes_slice.copy_from_slice(&cosmos_dest_address_bytes[..]);

    let contract_call = Gravity::new(gravity_contract, eth_client.clone())
        .send_to_cosmos(erc20, cosmos_dest_address_bytes_slice, amount)
        .gas(SEND_TO_COSMOS_GAS_LIMIT)
        .legacy();

    let pending_tx = contract_call.send().await?;
    let tx_hash = *pending_tx;
    info!("Sending to Cosmos with txid {}", tx_hash);
    // TODO(bolten): ethers interval default is 7s, this mirrors what web30 was doing, should we adjust?
    // additionally we are mirroring only waiting for 1 confirmation by leaving that as default
    let pending_tx = pending_tx.interval(Duration::from_secs(1));
    let potential_error = GravityError::GravityContractError(format!(
        "Did not receive transaction receipt when sending to Cosmos: {}",
        tx_hash
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
