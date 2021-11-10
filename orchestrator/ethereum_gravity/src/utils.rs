use clarity::abi::{Token, encode_call};
use clarity::Uint256;
use clarity::{abi::encode_tokens, Address as EthAddress};
use ethers::prelude::*;
use gravity_utils::error::GravityError;
use gravity_utils::types::*;
use sha3::{Digest, Keccak256};
use std::panic;
use web30::{client::Web3, jsonrpc::error::Web3Error};

pub type EthClient = Arc<SignerMiddleware<Provider<Http>, LocalWallet>>;

pub fn get_checkpoint_abi_encode(
    valset: &Valset,
    gravity_id: &str,
) -> Result<Vec<u8>, GravityError> {
    let (eth_addresses, powers) = valset.filter_empty_addresses();
    Ok(encode_tokens(&[
        Token::FixedString(gravity_id.to_string()),
        Token::FixedString("checkpoint".to_string()),
        valset.nonce.into(),
        eth_addresses.into(),
        powers.into(),
    ]))
}

pub fn get_checkpoint_hash(valset: &Valset, gravity_id: &str) -> Result<Vec<u8>, GravityError> {
    let locally_computed_abi_encode = get_checkpoint_abi_encode(&valset, &gravity_id);
    let locally_computed_digest = Keccak256::digest(&locally_computed_abi_encode?);
    Ok(locally_computed_digest.to_vec())
}

/// Gets the latest validator set nonce
pub async fn get_valset_nonce(
    contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {

    let payload = encode_call("state_lastValsetNonce()", &[]).unwrap();

    let val = web3
        .simulate_transaction(contract_address, 0u8.into(), payload, caller_address, None)
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_to_u64(real_num).expect("Valset nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_tx_batch_nonce(
    gravity_contract_address: EthAddress,
    erc20_contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let payload = encode_call("lastBatchNonce(address)", &[erc20_contract_address.into()]).unwrap();
    let val = web3
        .simulate_transaction(
            gravity_contract_address,
            0u8.into(),
            payload,
            caller_address,
            None,
        )
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_to_u64(real_num).expect("TxBatch nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_logic_call_nonce(
    gravity_contract_address: EthAddress,
    invalidation_id: Vec<u8>,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let payload = encode_call(
        "lastLogicCallNonce(bytes32)",
        &[Token::Bytes(invalidation_id)],
    )
    .unwrap();
    let val = web3
        .simulate_transaction(
            gravity_contract_address,
            0u8.into(),
            payload,
            caller_address,
            None,
        )
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_to_u64(real_num).expect("LogicCall nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_event_nonce(
    gravity_contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let payload = encode_call("state_lastEventNonce()", &[]).unwrap();
    let val = web3
        .simulate_transaction(
            gravity_contract_address,
            0u8.into(),
            payload,
            caller_address,
            None,
        )
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_to_u64(real_num).expect("EventNonce nonce overflow! Bridge Halt!"))
}

/// Gets the gravityID
pub async fn get_gravity_id(
    contract_address: EthAddress,
    caller_address: EthAddress,
    eth_client: EthClient,
) -> Result<String, GravityError> {
    const GAS_LIMIT: u128 = 12450000; // the most Hardhat will allow, will work on Geth

    let caller_balance = eth_client.get_balance(caller_address, None).await?;
    let latest_block = eth_client.get_block(BlockNumber::Latest).await?;
    let price = latest_block.base_fee_per_gas.ok_or(1u8.into()); // shouldn't happen unless pre-London
    let limit = min(GAS_LIMIT.into(), caller_balance / price.clone());

    let contract = Gravity::new(contract_address, eth_client);
    let contract_call = contract.state_gravityId()
        .from(caller_address)
        .gas(limit.into())
        .gas_price(price.into())
        .value(0u8.into());

    String::from_utf8(contract_call.call().await?)
}

/// Gets the ERC20 symbol, should maybe be upstreamed
pub async fn get_erc20_symbol(
    contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<String, GravityError> {

    let payload = encode_call("symbol()", &[]).unwrap();

    let val_symbol = web3
    .simulate_transaction(contract_address, 0u8.into(), payload, caller_address, None)
    .await?;
    // Pardon the unwrap, but this is temporary code, intended only for the tests, to help them
    // deal with a deprecated feature (the symbol), which will be removed soon
    Ok(String::from_utf8(val_symbol).unwrap())
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
