use crate::types::{LogicCall, TransactionBatch, Valset};
use ethers::core::abi::{self, Token};
use ethers::prelude::*;
use ethers::utils::hash_message;
use ethers::utils::keccak256;

/// takes the required input data and produces the required signature to confirm a validator
/// set update on the Gravity Ethereum contract. This value will then be signed before being
/// submitted to Cosmos, verified, and then relayed to Ethereum
/// Note: This is the message, you need to run Keccak256::digest() in order to get the 32byte
/// digest that is normally signed or may be used as a 'hash of the message'
pub fn encode_valset_confirm(gravity_id: String, valset: Valset) -> Vec<u8> {
    let (eth_addresses, powers) = valset.filter_empty_addresses();
    let eth_addresses = eth_addresses
        .iter()
        .map(|address| Token::Address(*address))
        .collect();
    let powers = powers
        .iter()
        .map(|power| Token::Uint((*power).into()))
        .collect();

    abi::encode(&[
        Token::FixedBytes(gravity_id.into_bytes()),
        Token::FixedBytes("checkpoint".to_string().into_bytes()),
        Token::Uint(valset.nonce.into()),
        Token::Array(eth_addresses),
        Token::Array(powers),
        Token::Uint(U256::zero()),
        Token::Address(H160::zero()),
    ])
}

pub fn encode_valset_confirm_hashed(gravity_id: String, valset: Valset) -> Vec<u8> {
    let digest = keccak256(encode_valset_confirm(gravity_id, valset).as_slice());
    hash_message(digest).as_bytes().to_vec()
}

#[test]
fn test_valset_signature() {
    use crate::{ethereum::hex_str_to_bytes, types::ValsetMember};
    use ethers::utils::keccak256;

    let correct_hash: Vec<u8> =
        hex_str_to_bytes("0x8cd4cc7f06bd39d4f77d94643a9ae6b3bdc3d3b78263683933cdd5c088452b9d")
            .unwrap();

    // a validator set
    let valset = Valset {
        nonce: 0,
        members: vec![
            ValsetMember {
                eth_address: Some(
                    "0xc783df8a850f42e7F7e57013759C285caa701eB6"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
            ValsetMember {
                eth_address: Some(
                    "0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
            ValsetMember {
                eth_address: Some(
                    "0xE5904695748fe4A84b40b3fc79De2277660BD1D3"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
        ],
    };
    let checkpoint = encode_valset_confirm("foo".to_string(), valset);
    let checkpoint_hash = keccak256(&checkpoint);
    assert_eq!(correct_hash, checkpoint_hash);

    // the same valset, except with an intentionally incorrect hash
    let valset = Valset {
        nonce: 1,
        members: vec![
            ValsetMember {
                eth_address: Some(
                    "0xc783df8a850f42e7F7e57013759C285caa701eB6"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
            ValsetMember {
                eth_address: Some(
                    "0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
            ValsetMember {
                eth_address: Some(
                    "0xE5904695748fe4A84b40b3fc79De2277660BD1D3"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
        ],
    };
    let checkpoint = encode_valset_confirm("foo".to_string(), valset);
    let checkpoint_hash = keccak256(&checkpoint);
    assert_ne!(correct_hash, checkpoint_hash)
}

/// takes the required input data and produces the required signature to confirm a transaction
/// batch on the Gravity Ethereum contract. This value will then be signed before being
/// submitted to Cosmos, verified, and then relayed to Ethereum
/// Note: This is the message, you need to run Keccak256::digest() in order to get the 32byte
/// digest that is normally signed or may be used as a 'hash of the message'
pub fn encode_tx_batch_confirm(gravity_id: String, batch: TransactionBatch) -> Vec<u8> {
    let (amounts, destinations, fees) = batch.get_checkpoint_values_tokens();

    abi::encode(&[
        Token::FixedBytes(gravity_id.into_bytes()),
        Token::FixedBytes("transactionBatch".to_string().into_bytes()),
        amounts,
        destinations,
        fees,
        Token::Uint(batch.nonce.into()),
        Token::Address(batch.token_contract),
        Token::Uint(batch.batch_timeout.into()),
    ])
}

pub fn encode_tx_batch_confirm_hashed(gravity_id: String, batch: TransactionBatch) -> Vec<u8> {
    let digest = keccak256(encode_tx_batch_confirm(gravity_id, batch).as_slice());
    hash_message(digest).as_bytes().to_vec()
}

#[tokio::test]
async fn test_batch_signature() {
    use crate::{
        ethereum::{hex_str_to_bytes, u8_slice_to_fixed_32},
        types::{BatchTransaction, Erc20Token},
    };
    use ethers::core::k256::ecdsa::SigningKey;
    use ethers::prelude::*;
    use ethers::utils::keccak256;
    use rand::Rng;

    let correct_hash: Vec<u8> =
        hex_str_to_bytes("0xa3a7ee0a363b8ad2514e7ee8f110d7449c0d88f3b0913c28c1751e6e0079a9b2")
            .unwrap();
    let erc20_addr = "0x835973768750b3ED2D5c3EF5AdcD5eDb44d12aD4"
        .parse()
        .unwrap();
    let sender_addr = "althea1c8nkaxk3d0p2gd7ummvmyqpdvqd6pkehqhwnnt"
        .parse()
        .unwrap();

    let token = Erc20Token {
        amount: 1u64.into(),
        token_contract_address: erc20_addr,
    };

    let batch = TransactionBatch {
        batch_timeout: 2111u64,
        nonce: 1u64,
        transactions: vec![BatchTransaction {
            id: 1u64,
            ethereum_recipient: "0x9FC9C2DfBA3b6cF204C37a5F690619772b926e39"
                .parse()
                .unwrap(),
            sender: sender_addr,
            erc20_fee: token.clone(),
            erc20_token: token.clone(),
        }],
        total_fee: token,
        token_contract: erc20_addr,
    };

    let checkpoint = encode_tx_batch_confirm("foo".to_string(), batch.clone());
    let checkpoint_hash = keccak256(&checkpoint);
    assert_eq!(correct_hash.len(), checkpoint_hash.len());
    assert_eq!(correct_hash, checkpoint_hash);

    // checkpoint is correct lets make sure our signature code works
    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    let eth_key = SigningKey::from_bytes(&secret).unwrap();
    let eth_wallet = LocalWallet::from(eth_key);
    let eth_address = eth_wallet.address();
    let checkpoint =
        keccak256(encode_tx_batch_confirm("foo".to_string(), batch.clone()).as_slice());
    let checkpoint_hash = encode_tx_batch_confirm_hashed("foo".to_string(), batch.clone());
    let checkpoint_hash = u8_slice_to_fixed_32(&checkpoint_hash).unwrap();

    let eth_signature = eth_wallet.sign_message(checkpoint).await.unwrap();

    assert_eq!(eth_address, eth_signature.recover(checkpoint_hash).unwrap());
}

#[tokio::test]
async fn test_specific_batch_signature() {
    use crate::{
        ethereum::u8_slice_to_fixed_32,
        types::{BatchTransaction, Erc20Token},
    };
    use ethers::core::k256::ecdsa::SigningKey;
    use ethers::prelude::*;
    use ethers::utils::keccak256;
    use rand::Rng;

    let erc20_addr = "0x0635FF793Edf48cf5dB294916720A78e6e490E40"
        .parse()
        .unwrap();
    let sender_addr = "cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7"
        .parse()
        .unwrap();

    let token = Erc20Token {
        amount: 1u64.into(),
        token_contract_address: erc20_addr,
    };

    let batch = TransactionBatch {
        batch_timeout: 4427201u64,
        nonce: 15u64,
        transactions: vec![BatchTransaction {
            id: 1301u64,
            ethereum_recipient: "0x64D110e00064F2b428476cD64295d8E35836ffd6"
                .parse()
                .unwrap(),
            sender: sender_addr,
            erc20_fee: token.clone(),
            erc20_token: token.clone(),
        }],
        total_fee: token,
        token_contract: erc20_addr,
    };

    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    // the starting location of the funds
    let eth_key = SigningKey::from_bytes(&secret).unwrap();
    let eth_wallet = LocalWallet::from(eth_key);
    let eth_address = eth_wallet.address();

    let checkpoint =
        keccak256(encode_tx_batch_confirm("foo".to_string(), batch.clone()).as_slice());
    let checkpoint_hash = encode_tx_batch_confirm_hashed("foo".to_string(), batch.clone());
    let checkpoint_hash = u8_slice_to_fixed_32(&checkpoint_hash).unwrap();

    let eth_signature = eth_wallet.sign_message(checkpoint).await.unwrap();

    assert_eq!(eth_address, eth_signature.recover(checkpoint_hash).unwrap());
}

/// takes the required input data and produces the required signature to confirm a logic
/// call on the Gravity Ethereum contract. This value will then be signed before being
/// submitted to Cosmos, verified, and then relayed to Ethereum
/// Note: This is the message, you need to run Keccak256::digest() in order to get the 32byte
/// digest that is normally signed or may be used as a 'hash of the message'
pub fn encode_logic_call_confirm(gravity_id: String, call: LogicCall) -> Vec<u8> {
    let transfer_amounts = call
        .transfers
        .iter()
        .map(|transfer| Token::Uint(transfer.amount))
        .collect();
    let transfer_token_contracts = call
        .transfers
        .iter()
        .map(|transfer| Token::Address(transfer.token_contract_address))
        .collect();
    let fee_amounts = call
        .fees
        .iter()
        .map(|fee| Token::Uint(fee.amount))
        .collect();
    let fee_token_contracts = call
        .fees
        .iter()
        .map(|fee| Token::Address(fee.token_contract_address))
        .collect();

    abi::encode(&[
        Token::FixedBytes(gravity_id.into_bytes()), // Gravity Instance ID
        Token::FixedBytes("logicCall".to_string().into_bytes()), // Function Name
        Token::Array(transfer_amounts),             // Array of Transfer amounts
        Token::Array(transfer_token_contracts),     // ERC-20 contract for transfers
        Token::Array(fee_amounts),                  // Array of Fees
        Token::Array(fee_token_contracts),          // ERC-20 contract for fee payments
        Token::Address(call.logic_contract_address), // Address of a logic contract
        Token::Bytes(call.payload),                 // Encoded arguments to logic contract
        Token::Uint(call.timeout.into()),           // Timeout on batch
        Token::FixedBytes(call.invalidation_id),    // Scope of logic batch
        Token::Uint(call.invalidation_nonce.into()), // Nonce of logic batch. See 2-d nonce scheme.
    ])
}

pub fn encode_logic_call_confirm_hashed(gravity_id: String, call: LogicCall) -> Vec<u8> {
    let digest = keccak256(encode_logic_call_confirm(gravity_id, call).as_slice());
    hash_message(digest).as_bytes().to_vec()
}

#[test]
fn test_logic_call_signature() {
    use crate::{
        ethereum::hex_str_to_bytes,
        types::{Erc20Token, LogicCall},
    };
    use ethers::utils::keccak256;

    let correct_hash: Vec<u8> =
        hex_str_to_bytes("0x1de95c9ace999f8ec70c6dc8d045942da2612950567c4861aca959c0650194da")
            .unwrap();
    let token_contract_address = "0xC26eFfa98B8A2632141562Ae7E34953Cfe5B4888"
        .parse()
        .unwrap();
    let logic_contract_address = "0x17c1736CcF692F653c433d7aa2aB45148C016F68"
        .parse()
        .unwrap();
    let token = vec![Erc20Token {
        amount: 1u8.into(),
        token_contract_address,
    }];

    let logic_call = LogicCall {
        transfers: token.clone(),
        fees: token,
        logic_contract_address,
        payload: hex_str_to_bytes(
            "0x74657374696e675061796c6f6164000000000000000000000000000000000000",
        )
        .unwrap(),
        timeout: 4766922941000,
        invalidation_id: hex_str_to_bytes(
            "0x696e76616c69646174696f6e4964000000000000000000000000000000000000",
        )
        .unwrap(),
        invalidation_nonce: 1u8.into(),
    };
    let checkpoint = encode_logic_call_confirm("foo".to_string(), logic_call);
    println!("{}", checkpoint.len() / 32);

    let checkpoint_hash = keccak256(&checkpoint);

    assert_eq!(correct_hash.len(), checkpoint_hash.len());
    assert_eq!(correct_hash, checkpoint_hash)
}
