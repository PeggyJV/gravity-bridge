use bytes::BytesMut;
use deep_space::address::Address;
use deep_space::coin::Coin;
use deep_space::error::CosmosGrpcError;
use deep_space::Contact;
use deep_space::Fee;
use deep_space::Msg;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use ethers::utils::keccak256;
use gravity_proto::cosmos_sdk_proto::cosmos::base::abci::v1beta1::TxResponse;
use gravity_proto::cosmos_sdk_proto::cosmos::tx::v1beta1::BroadcastMode;
use gravity_proto::gravity as proto;
use gravity_utils::error::GravityError;
use gravity_utils::ethereum::format_eth_address;
use prost::Message;
use std::cmp;
use std::collections::HashSet;
use std::{result::Result, time::Duration};

use crate::crypto::PrivateKey as CosmosPrivateKey;

pub const MEMO: &str = "Sent using Gravity Bridge Orchestrator";
pub const TIMEOUT: Duration = Duration::from_secs(60);

/// Send a transaction updating the eth address for the sending
/// Cosmos address. The sending Cosmos address should be a validator
pub async fn update_gravity_delegate_addresses(
    contact: &Contact,
    delegate_eth_address: EthAddress,
    delegate_cosmos_address: Address,
    cosmos_key: CosmosPrivateKey,
    ethereum_wallet: LocalWallet,
    gas_price: (f64, String),
    gas_adjustment: f64,
) -> Result<TxResponse, GravityError> {
    let our_valoper_address = cosmos_key
        .to_address(&contact.get_prefix())
        .unwrap()
        // This works so long as the format set by the cosmos hub is maintained
        // having a main prefix followed by a series of titles for specific keys
        // this will not work if that convention is broken. This will be resolved when
        // GRPC exposes prefix endpoints (coming to upstream cosmos sdk soon)
        .to_bech32(format!("{}valoper", contact.get_prefix()))
        .unwrap();

    let nonce = contact
        .get_account_info(cosmos_key.to_address(&contact.get_prefix()).unwrap())
        .await?
        .sequence;

    let eth_sign_msg = proto::DelegateKeysSignMsg {
        validator_address: our_valoper_address.clone(),
        nonce,
    };

    let mut data = BytesMut::with_capacity(eth_sign_msg.encoded_len());
    Message::encode(&eth_sign_msg, &mut data).expect("encoding failed");

    let data_hash = keccak256(data);
    let eth_signature = ethereum_wallet.sign_message(data_hash).await?;
    let msg = proto::MsgDelegateKeys {
        validator_address: our_valoper_address.to_string(),
        orchestrator_address: delegate_cosmos_address.to_string(),
        ethereum_address: format_eth_address(delegate_eth_address),
        eth_signature: eth_signature.to_vec(),
    };
    let msg = Msg::new("/gravity.v1.MsgDelegateKeys", msg);

    send_messages(contact, cosmos_key, gas_price, vec![msg], gas_adjustment).await
}

/// Sends tokens from Cosmos to Ethereum. These tokens will not be sent immediately instead
/// they will require some time to be included in a batch
pub async fn send_to_eth(
    cosmos_key: CosmosPrivateKey,
    destination: EthAddress,
    amount: Coin,
    bridge_fee: Coin,
    gas_price: (f64, String),
    contact: &Contact,
    gas_adjustment: f64,
) -> Result<TxResponse, GravityError> {
    if amount.denom != bridge_fee.denom {
        return Err(GravityError::CosmosGrpcError(CosmosGrpcError::BadInput(
            format!(
                "The amount ({}) and bridge_fee ({}) denominations do not match.",
                amount.denom, bridge_fee.denom,
            ),
        )));
    }

    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

    let msg = proto::MsgSendToEthereum {
        sender: cosmos_address.to_string(),
        ethereum_recipient: format_eth_address(destination),
        amount: Some(amount.into()),
        bridge_fee: Some(bridge_fee.clone().into()),
    };
    let msg = Msg::new("/gravity.v1.MsgSendToEthereum", msg);
    send_messages(contact, cosmos_key, gas_price, vec![msg], gas_adjustment).await
}

pub async fn send_request_batch_tx(
    cosmos_key: CosmosPrivateKey,
    denom: String,
    gas_price: (f64, String),
    contact: &Contact,
    gas_adjustment: f64,
) -> Result<TxResponse, GravityError> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let msg_request_batch = proto::MsgRequestBatchTx {
        signer: cosmos_address.to_string(),
        denom,
    };
    let msg = Msg::new("/gravity.v1.MsgRequestBatchTx", msg_request_batch);
    send_messages(contact, cosmos_key, gas_price, vec![msg], gas_adjustment).await
}

pub async fn send_messages(
    contact: &Contact,
    cosmos_key: CosmosPrivateKey,
    gas_price: (f64, String),
    messages: Vec<Msg>,
    gas_adjustment: f64,
) -> Result<TxResponse, GravityError> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

    let fee_amount = Coin {
        denom: gas_price.1.clone(),
        amount: 0u32.into(),
    };

    let fee = Fee {
        amount: vec![fee_amount],
        gas_limit: 0,
        granter: None,
        payer: None,
    };

    let mut args = contact.get_message_args(cosmos_address, fee).await?;

    let tx_parts = cosmos_key.build_tx(&messages, args.clone(), MEMO)?;
    let gas = contact.simulate_tx(tx_parts).await?;

    // multiply the estimated gas by the configured gas adjustment
    let gas_limit: f64 = (gas.gas_used as f64) * gas_adjustment;
    args.fee.gas_limit = cmp::max(gas_limit as u64, 500000 * messages.len() as u64);

    // compute the fee as fee=ceil(gas_limit * gas_price)
    let fee_amount: f64 = args.fee.gas_limit as f64 * gas_price.0;
    let fee_amount: u64 = fee_amount.abs().ceil() as u64;
    let fee_amount = Coin {
        denom: gas_price.1,
        amount: fee_amount.into(),
    };
    args.fee.amount = vec![fee_amount];

    let msg_bytes = cosmos_key.sign_std_msg(&messages, args, MEMO)?;
    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;

    Ok(contact.wait_for_tx(response, TIMEOUT).await?)
}

pub async fn send_main_loop(
    contact: &Contact,
    cosmos_key: CosmosPrivateKey,
    gas_price: (f64, String),
    mut rx: tokio::sync::mpsc::Receiver<Vec<Msg>>,
    gas_adjustment: f64,
    msg_batch_size: usize,
) {
    while let Some(messages) = rx.recv().await {
        for msg_chunk in messages.chunks(msg_batch_size) {
            match send_messages(
                contact,
                cosmos_key,
                gas_price.to_owned(),
                msg_chunk.to_vec(),
                gas_adjustment,
            )
            .await
            {
                Ok(res) => trace!("okay: {:?}", res),
                Err(err) => {
                    let msg_types = msg_chunk
                        .iter()
                        .map(|msg| prost_types::Any::from(msg.clone()).type_url)
                        .collect::<HashSet<String>>();

                    error!(
                        "Error during gRPC call to Cosmos containing {} messages of types {:?}: {:?}",
                        msg_chunk.len(),
                        msg_types,
                        err
                    );
                }
            }
        }
    }
}
