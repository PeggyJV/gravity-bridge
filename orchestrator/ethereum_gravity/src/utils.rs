use ethers::core::abi::{self, Token};
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use ethers::utils::keccak256;
use gravity_utils::error::GravityError;
use gravity_utils::ethereum::downcast_to_u64;
use gravity_utils::gravity::*;
use gravity_utils::types::*;
use std::panic;

pub type EthSignerMiddleware = SignerMiddleware<Provider<Http>, LocalWallet>;
pub type EthClient = Arc<EthSignerMiddleware>;

pub fn get_checkpoint_abi_encode(
    valset: &Valset,
    gravity_id: &str,
) -> Result<Vec<u8>, GravityError> {
    let (eth_addresses, powers) = valset.filter_empty_addresses();
    let eth_addresses = eth_addresses.iter().map (|address| Token::Address(*address)).collect();
    let powers = powers.iter().map(|power| Token::Uint((*power).into())).collect();

    Ok(abi::encode(&[
        Token::FixedBytes(gravity_id.into_bytes()),
        Token::FixedBytes("checkpoint".to_string().into_bytes()),
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
    contract_address: EthAddress,
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, EthereumRestError> {
    let (price, limit) = get_contract_call_gas(eth_client, caller_address).await?;
    let contract: Gravity<EthSignerMiddleware> = Gravity::new(contract_address, eth_client);
    let contract_call = contract.state_lastValsetNonce()
        .from(caller_address)
        .gas(limit)
        .gas_price(price)
        .value(0u8.into());

    let val = contract_call.call().await?;

    // TODO (bolten): do we actually want to halt the bridge as the original comment implies?
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    Ok(downcast_to_u64(val).expect("Valset nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_tx_batch_nonce(
    gravity_contract_address: EthAddress,
    erc20_contract_address: EthAddress,
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, EthereumRestError> {
    let (price, limit) = get_contract_call_gas(eth_client, caller_address).await?;
    let contract: Gravity<EthSignerMiddleware> = Gravity::new(contract_address, eth_client);
    let contract_call = contract.last_batch_nonce(erc20_contract_address)
        .from(caller_address)
        .gas(limit)
        .gas_price(price)
        .value(0u8.into());

    let val = contract_call.call().await?;

    // TODO (bolten): do we actually want to halt the bridge as the original comment implies?
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    Ok(downcast_to_u64(val).expect("TxBatch nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_logic_call_nonce(
    gravity_contract_address: EthAddress,
    invalidation_id: Vec<u8>,
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, EthereumRestError> {
    let (price, limit) = get_contract_call_gas(eth_client, caller_address).await?;
    let contract: Gravity<EthSignerMiddleware> = Gravity::new(contract_address, eth_client);
    let contract_call = contract.last_logic_call_nonce(invalidation_id.as_slice())
        .from(caller_address)
        .gas(limit)
        .gas_price(price)
        .value(0u8.into());

    let val = contract_call.call().await?;

    // TODO (bolten): do we actually want to halt the bridge as the original comment implies?
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    Ok(downcast_to_u64(val).expect("LogicCall nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_event_nonce(
    gravity_contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let (price, limit) = get_contract_call_gas(eth_client, caller_address).await?;
    let contract: Gravity<EthSignerMiddleware> = Gravity::new(contract_address, eth_client);
    let contract_call = contract.state_last_event_nonce()
        .from(caller_address)
        .gas(limit)
        .gas_price(price)
        .value(0u8.into());

    let val = contract_call.call().await?;

    // TODO (bolten): do we actually want to halt the bridge as the original comment implies?
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    Ok(downcast_to_u64(real_num).expect("EventNonce nonce overflow! Bridge Halt!"))
}

/// Gets the gravityID
pub async fn get_gravity_id(
    contract_address: EthAddress,
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<String, GravityError> {
    let (price, limit) = get_contract_call_gas(eth_client, caller_address).await?;
    let contract: Gravity<EthSignerMiddleware> = Gravity::new(contract_address, eth_client);
    let contract_call = contract.state_gravity_id()
        .from(caller_address)
        .gas(limit)
        .gas_price(price)
        .value(0u8.into());

    String::from_utf8(contract_call.call().await?.to_vec())
}

pub async fn get_contract_call_gas(
    eth_client: EthClient,
    caller_address: EthAddress
) -> Result<(price, limit), GravityError> {
    const GAS_LIMIT: u128 = 12450000; // the most Hardhat will allow, will work on Geth

    let caller_balance = eth_client.get_balance(caller_address, None).await?;
    let latest_block = eth_client.get_block(BlockNumber::Latest).await?;
    let price = latest_block.base_fee_per_gas.ok_or(1u8.into()); // shouldn't happen unless pre-London
    let limit = min(GAS_LIMIT.into(), caller_balance / price.clone());

    Ok((price, limit))
}

/// Just a helper struct to represent the cost of actions on Ethereum
#[derive(Debug, Default, Clone)]
pub struct GasCost {
    pub gas: Uint256,
    pub gas_price: Uint256,
}

impl GasCost {
    pub fn get_total(&self) -> Uint256 {
        self.gas.clone() * self.gas_price.clone()
    }
}
