use ethers::abi::Detokenize;
use ethers::contract::builders::ContractCall;
use ethers::core::abi::{self, Token};
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use ethers::utils::keccak256;
use gravity_abi::gravity::*;
use gravity_utils::error::GravityError;
use gravity_utils::ethereum::downcast_to_u64;
use gravity_utils::types::*;
use std::cmp::min;
use std::sync::Arc;

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
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, GravityError> {
    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .state_last_valset_nonce();
    let contract_call = build_contract_eth_call(contract_call, eth_client.clone()).await?;
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
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, GravityError> {
    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .last_batch_nonce(erc20_contract_address);
    let contract_call = build_contract_eth_call(contract_call, eth_client.clone()).await?;
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
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, GravityError> {
    if invalidation_id.len() != 32 {
        return Err(GravityError::InvalidArgumentError(format!(
            "Error getting logic call nonce, invalidation id is not 32 bytes: {:?}", invalidation_id)))
    }

    let mut invalidation_id_slice: [u8; 32] = Default::default();
    invalidation_id_slice.copy_from_slice(&invalidation_id[..]);

    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .last_logic_call_nonce(invalidation_id_slice);
    let contract_call = build_contract_eth_call(contract_call, eth_client.clone()).await?;
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
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<u64, GravityError> {
    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .state_last_event_nonce();
    let contract_call = build_contract_eth_call(contract_call, eth_client.clone()).await?;
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
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<String, GravityError> {
    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .state_gravity_id();
    let contract_call = build_contract_eth_call(contract_call, eth_client.clone()).await?;
    let gravity_id = contract_call.call().await?;
    let id_as_string = String::from_utf8(gravity_id.to_vec());

    match id_as_string {
        Ok(id) => Ok(id),
        Err(err) => Err(GravityError::GravityContractError(format!(
            "Received invalid utf8 when getting gravity id: {:?}", &gravity_id
        )))
    }
}

/// Since all the contract eth_calls here use the same gas and value settings, use a common
/// function to append them to the ContractCall builder.
///
/// Retrieve gas price and limit in a similar fashion to web30's simulate_transaction.
/// These values are intended to be used in conjunction with eth_call rather than
/// eth_sendtransaction. In ethers this is represented by `call()` on a ContractCall rather
/// than `send()`. Using `call()` will not send a transaction from the caller account or
/// spend gas.
pub async fn build_contract_eth_call<D: Detokenize>(
    contract_call: ContractCall<EthSignerMiddleware, D>,
    eth_client: EthClient,
) -> Result<ContractCall<EthSignerMiddleware, D>, GravityError> {
    const GAS_LIMIT: u128 = 12450000; // the most Hardhat will allow, will work on Geth

    let caller_balance = eth_client.get_balance(eth_client.address(), None).await?;
    let latest_block = eth_client.get_block(BlockNumber::Latest).await?.unwrap();
    let price = latest_block.base_fee_per_gas.unwrap_or(1u8.into()); // "or" clause shouldn't happen unless pre-London
    if price == U256::zero() {
        return Err(GravityError::EthereumBadDataError("Latest block returned base fee per gas of zero".to_string()));
    }

    let limit = min(GAS_LIMIT.into(), caller_balance.div_mod(price).0);

    Ok(contract_call.from(eth_client.address())
        .gas(limit)
        .gas_price(price)
        .value(U256::zero()))
}

pub async fn get_max_gas_cost(eth_client: EthClient) -> Result<GasCost, GravityError> {
    let our_balance = eth_client.get_balance(eth_client.address(), None).await?;

    Ok(GasCost {
        gas: min((u64::MAX - 1).into(), our_balance),
        gas_price: eth_client.get_gas_price().await?,
    })
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
