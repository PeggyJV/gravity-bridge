use crate::types::EthClient;
use ethers::core::abi::{self, Token};
use ethers::middleware::gas_oracle::{Etherscan, GasCategory};
use ethers::prelude::gas_oracle::GasOracle;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use ethers::utils::keccak256;
use gravity_abi::gravity::*;
use gravity_utils::error::GravityError;
use gravity_utils::ethereum::downcast_to_u64;
use gravity_utils::types::*;
use std::cmp::min;

pub fn get_checkpoint_abi_encode(
    valset: &Valset,
    gravity_id: &str,
) -> Result<Vec<u8>, GravityError> {
    let (eth_addresses, powers) = valset.filter_empty_addresses();
    let eth_addresses = eth_addresses
        .iter()
        .map(|address| Token::Address(*address))
        .collect();
    let powers = powers
        .iter()
        .map(|power| Token::Uint((*power).into()))
        .collect();

    Ok(abi::encode(&[
        Token::FixedBytes(gravity_id.as_bytes().to_vec()),
        Token::FixedBytes("checkpoint".as_bytes().to_vec()),
        Token::Uint(valset.nonce.into()),
        Token::Array(eth_addresses),
        Token::Array(powers),
    ]))
}

pub fn get_checkpoint_hash(valset: &Valset, gravity_id: &str) -> Result<Vec<u8>, GravityError> {
    let locally_computed_abi_encode = get_checkpoint_abi_encode(&valset, &gravity_id)?;
    let locally_computed_digest = keccak256(locally_computed_abi_encode.as_slice());
    Ok(locally_computed_digest.to_vec())
}

/// Gets the latest validator set nonce
pub async fn get_valset_nonce(
    gravity_contract_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, GravityError> {
    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .state_last_valset_nonce()
        .from(eth_client.address())
        .value(U256::zero());
    let gas_cost = get_call_gas_cost(eth_client.clone()).await?;
    let contract_call = contract_call
        .gas(gas_cost.gas)
        .gas_price(gas_cost.gas_price);

    let valset_nonce = contract_call.call().await?;

    // TODO (bolten): do we actually want to halt the bridge as the original comment implies?
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    Ok(downcast_to_u64(valset_nonce).expect("Valset nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_tx_batch_nonce(
    gravity_contract_address: EthAddress,
    erc20_contract_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, GravityError> {
    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .last_batch_nonce(erc20_contract_address)
        .from(eth_client.address())
        .value(U256::zero());
    let gas_cost = get_call_gas_cost(eth_client.clone()).await?;
    let contract_call = contract_call
        .gas(gas_cost.gas)
        .gas_price(gas_cost.gas_price);

    let tx_batch_nonce = contract_call.call().await?;

    // TODO (bolten): do we actually want to halt the bridge as the original comment implies?
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    Ok(downcast_to_u64(tx_batch_nonce).expect("TxBatch nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_logic_call_nonce(
    gravity_contract_address: EthAddress,
    invalidation_id: Vec<u8>,
    eth_client: EthClient,
) -> Result<u64, GravityError> {
    let invalidation_id = convert_invalidation_id_to_fixed_array(invalidation_id)?;

    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .last_logic_call_nonce(invalidation_id)
        .from(eth_client.address())
        .value(U256::zero());
    let gas_cost = get_call_gas_cost(eth_client.clone()).await?;
    let contract_call = contract_call
        .gas(gas_cost.gas)
        .gas_price(gas_cost.gas_price);

    let logic_call_nonce = contract_call.call().await?;

    // TODO (bolten): do we actually want to halt the bridge as the original comment implies?
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    Ok(downcast_to_u64(logic_call_nonce).expect("LogicCall nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_event_nonce(
    gravity_contract_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, GravityError> {
    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .state_last_event_nonce()
        .from(eth_client.address())
        .value(U256::zero());
    let gas_cost = get_call_gas_cost(eth_client.clone()).await?;
    let contract_call = contract_call
        .gas(gas_cost.gas)
        .gas_price(gas_cost.gas_price);

    let event_nonce = contract_call.call().await?;

    // TODO (bolten): do we actually want to halt the bridge as the original comment implies?
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    Ok(downcast_to_u64(event_nonce).expect("EventNonce nonce overflow! Bridge Halt!"))
}

/// Gets the gravityID
pub async fn get_gravity_id(
    gravity_contract_address: EthAddress,
    eth_client: EthClient,
) -> Result<String, GravityError> {
    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .state_gravity_id()
        .from(eth_client.address())
        .value(U256::zero());
    let gas_cost = get_call_gas_cost(eth_client.clone()).await?;
    let contract_call = contract_call
        .gas(gas_cost.gas)
        .gas_price(gas_cost.gas_price);

    let gravity_id = contract_call.call().await?;
    let id_as_string = String::from_utf8(gravity_id.to_vec());

    match id_as_string {
        Ok(id) => Ok(id),
        Err(err) => Err(GravityError::GravityContractError(format!(
            "Received invalid utf8 when getting gravity id {:?}: {}",
            &gravity_id, err
        ))),
    }
}

/// Retrieve gas price and limit in a similar fashion to web30's simulate_transaction.
/// These values are intended to be used in conjunction with eth_call rather than
/// eth_sendtransaction. In ethers this is represented by `call()` on a ContractCall rather
/// than `send()`. Using `call()` will not send a transaction from the caller account or
/// spend gas.
pub async fn get_call_gas_cost(eth_client: EthClient) -> Result<GasCost, GravityError> {
    const GAS_LIMIT: u128 = 12450000; // the most Hardhat will allow, will work on Geth

    let caller_balance = eth_client.get_balance(eth_client.address(), None).await?;
    let latest_block = eth_client.get_block(BlockNumber::Latest).await?.unwrap();
    let gas_price = latest_block.base_fee_per_gas.unwrap_or(1u8.into()); // "or" clause shouldn't happen unless pre-London
    if gas_price == U256::zero() {
        return Err(GravityError::EthereumBadDataError(
            "Latest block returned base fee per gas of zero".to_string(),
        ));
    }

    let gas = min(GAS_LIMIT.into(), caller_balance.div_mod(gas_price).0);

    Ok(GasCost { gas, gas_price })
}

/// If ETHERSCAN_API_KEY env var is set, we'll call out to Etherscan for a gas estimate.
/// Otherwise, just call eth_gasPrice.
pub async fn get_send_transaction_gas_price(eth_client: EthClient) -> Result<U256, GravityError> {
    if let Ok(api_key) = std::env::var("ETHERSCAN_API_KEY") {
        let etherscan_oracle =
            Etherscan::new(Some(api_key.as_str())).category(GasCategory::Standard);
        return Ok(etherscan_oracle.fetch().await?);
    }

    Ok(eth_client.get_gas_price().await?)
}

pub fn convert_invalidation_id_to_fixed_array(
    invalidation_id: Vec<u8>,
) -> Result<[u8; 32], GravityError> {
    if invalidation_id.len() != 32 {
        return Err(GravityError::InvalidArgumentError(format!(
            "Error getting logic call nonce, invalidation id is not 32 bytes: {:?}",
            invalidation_id
        )));
    }

    let mut invalidation_id_slice: [u8; 32] = Default::default();
    invalidation_id_slice.copy_from_slice(&invalidation_id[..]);
    Ok(invalidation_id_slice)
}

/// Just a helper struct to represent the cost of actions on Ethereum
#[derive(Debug, Default, Clone)]
pub struct GasCost {
    pub gas: U256,
    pub gas_price: U256,
}

impl GasCost {
    pub fn get_total(&self) -> U256 {
        self.gas * self.gas_price
    }
}
