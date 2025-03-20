use crate::utils::error::GravityError;
use crate::utils::ethereum::format_eth_address;
use crate::Coin as FromCoin;
use bytes::BytesMut;
use cosmos_sdk_proto_althea::cosmos::base::abci::v1beta1::TxResponse;
use cosmos_sdk_proto_althea::cosmos::tx::v1beta1::BroadcastMode;
use deep_space::address::Address;
use deep_space::coin::Coin;
use deep_space::error::CosmosGrpcError;
use deep_space::Contact;
use deep_space::CosmosPrivateKey;
use deep_space::Fee;
use deep_space::Msg;
use deep_space::PrivateKey;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use ethers::utils::keccak256;
use gravity_proto::gravity as proto;
use prost::Message;
use std::cmp;
use std::collections::HashSet;
use std::{result::Result, time::Duration};

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
    let from_amount_coin: FromCoin = amount.into();
    let from_bridge_fee_coin: FromCoin = bridge_fee.into();
    let msg = proto::MsgSendToEthereum {
        sender: cosmos_address.to_string(),
        ethereum_recipient: format_eth_address(destination),
        amount: Some(from_amount_coin.clone().into()),
        bridge_fee: Some(from_bridge_fee_coin.clone().into()),
    };
    let msg = Msg::new("/gravity.v1.MsgSendToEthereum", msg);
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

    let mut args = contact.get_message_args(cosmos_address, fee, None).await?;
    let gas = contact
        .simulate_tx(messages.as_slice(), None, cosmos_key)
        .await?;

    // multiply the estimated gas by the configured gas adjustment
    if let Some(gas_info) = gas.gas_info {
        let gas_limit: f64 = (gas_info.gas_used as f64) * gas_adjustment;
        args.fee.gas_limit = cmp::max(gas_limit as u64, 500000 * messages.len() as u64);
    } else {
        // Is this a good gas limit max?
        args.fee.gas_limit = cmp::max(2500000 as u64, 500000 * messages.len() as u64);
    }

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

    let response = contact.wait_for_tx(response, TIMEOUT).await?;
    let tx_response: TxResponse = response.into();

    Ok(tx_response)
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
            let batch = msg_chunk.to_vec();
            match send_messages(
                contact,
                cosmos_key,
                gas_price.to_owned(),
                msg_chunk.to_vec(),
                gas_adjustment,
            )
            .await
            {
                Ok(res) => debug!("message batch sent: {:?}", res),
                Err(err) => {
                    log_send_error(&batch, err);

                    // multiple messages in a single Cosmos transaction will be rejected
                    // atomically if that transaction cannot be delivered, so retry each
                    // element separately
                    info!("Trying each message in batch individually");
                    for msg in batch {
                        let msg_vec = vec![msg];
                        match send_messages(
                            contact,
                            cosmos_key,
                            gas_price.to_owned(),
                            msg_vec.clone(),
                            gas_adjustment,
                        )
                        .await
                        {
                            Ok(res) => debug!("message sent: {:?}", res),
                            Err(err) => log_send_error(&msg_vec, err),
                        }
                    }
                }
            }
        }
    }
}

fn log_send_error(messages: &Vec<Msg>, err: GravityError) {
    let msg_types = messages
        .iter()
        .map(|msg| prost_types::Any::from(msg.clone()).type_url)
        .collect::<HashSet<String>>();

    error!(
        "Error during gRPC call to Cosmos containing {} messages of types {:?}: {:?}",
        messages.len(),
        msg_types,
        err
    );
}
