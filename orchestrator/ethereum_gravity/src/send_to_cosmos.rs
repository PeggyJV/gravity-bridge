//! Helper functions for sending tokens to Cosmos

use crate::utils::EthClient;
use deep_space::address::Address as CosmosAddress;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_utils::error::GravityError;
use std::time::Duration;

const SEND_TO_COSMOS_GAS_LIMIT: u128 = 100_000;

#[allow(clippy::too_many_arguments)]
pub async fn send_to_cosmos(
    erc20: Address,
    gravity_contract: Address,
    amount: U256,
    cosmos_destination: CosmosAddress,
    wait_timeout: Option<Duration>,
    eth_client: EthClient,
) -> Result<U256, GravityError> {
    let sender_address = eth_client.address();
    let mut approve_nonce = None;

    let approved =
        .check_erc20_approved(erc20, sender_address, gravity_contract)
        .await?;
    if !approved {
        let mut options = options.clone();
        let nonce = web3.eth_get_transaction_count(sender_address).await?;
        options.push(SendTxOption::Nonce(nonce.clone()));
        approve_nonce = Some(nonce);
        let txid = web3
            .approve_erc20_transfers(erc20, sender_secret, gravity_contract, None, options)
            .await?;
        trace!(
            "We are not approved for ERC20 transfers, approving txid: {:#066x}",
            txid
        );
        if let Some(timeout) = wait_timeout {
            web3.wait_for_transaction(txid, timeout, None).await?;
            trace!("Approval finished!")
        }
    }

    // if the user sets a gas limit we should honor it, if they don't we
    // should add the default
    let mut has_gas_limit = false;
    let mut options = options;
    for option in options.iter() {
        if let SendTxOption::GasLimit(_) = option {
            has_gas_limit = true;
            break;
        }
    }
    if !has_gas_limit {
        options.push(SendTxOption::GasLimit(SEND_TO_COSMOS_GAS_LIMIT.into()));
    }
    // if we have run an approval we should increment our nonce by one so that
    // we can be sure our actual tx can go in immediately behind
    if let Some(nonce) = approve_nonce {
        options.push(SendTxOption::Nonce(nonce + 1u8.into()));
    }

    // This code deals with some specifics of Ethereum byte encoding, Ethereum is BigEndian
    // so small values like addresses that don't take up the full length of the byte vector
    // are pushed up to the top. This duplicates the way Ethereum encodes it's own addresses
    // as closely as possible.
    let mut cosmos_dest_address_bytes = cosmos_destination.as_bytes().to_vec();
    while cosmos_dest_address_bytes.len() < 32 {
        cosmos_dest_address_bytes.insert(0, 0u8);
    }
    let encoded_destination_address = Token::Bytes(cosmos_dest_address_bytes);

    let tx_hash = web3
        .send_transaction(
            gravity_contract,
            encode_call(
                "sendToCosmos(address,bytes32,uint256)",
                &[
                    erc20.into(),
                    encoded_destination_address,
                    amount.clone().into(),
                ],
            )?,
            0u32.into(),
            sender_address,
            sender_secret,
            options,
        )
        .await?;

    if let Some(timeout) = wait_timeout {
        web3.wait_for_transaction(tx_hash.clone(), timeout, None)
            .await?;
    }

    Ok(tx_hash)
}
