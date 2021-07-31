use crate::build;

use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use deep_space::address::Address;
use deep_space::coin::Coin;
use deep_space::error::CosmosGrpcError;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use deep_space::Contact;
use deep_space::Fee;
use deep_space::Msg;
use gravity_proto::cosmos_sdk_proto::cosmos::base::abci::v1beta1::TxResponse;
use gravity_proto::cosmos_sdk_proto::cosmos::tx::v1beta1::service_client::ServiceClient as TxServiceClient;
use gravity_proto::cosmos_sdk_proto::cosmos::tx::v1beta1::BroadcastMode;
use gravity_proto::cosmos_sdk_proto::cosmos::tx::v1beta1::BroadcastTxRequest;
use gravity_proto::gravity as proto;

use gravity_utils::types::*;
use std::time::Duration;

use bytes::BytesMut;
use prost::Message;

pub const MEMO: &str = "Sent using Althea Orchestrator";
pub const TIMEOUT: Duration = Duration::from_secs(60);

/// Send a transaction updating the eth address for the sending
/// Cosmos address. The sending Cosmos address should be a validator
pub async fn update_gravity_delegate_addresses(
    contact: &Contact,
    delegate_eth_address: EthAddress,
    delegate_cosmos_address: Address,
    cosmos_private_key: CosmosPrivateKey,
    eth_private_key: EthPrivateKey,
    fee: Coin,
) -> Result<TxResponse, CosmosGrpcError> {
    trace!("Updating Gravity Delegate addresses");
    let our_valoper_address = cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap()
        // This works so long as the format set by the cosmos hub is maintained
        // having a main prefix followed by a series of titles for specific keys
        // this will not work if that convention is broken. This will be resolved when
        // GRPC exposes prefix endpoints (coming to upstream cosmos sdk soon)
        .to_bech32(format!("{}valoper", contact.get_prefix()))
        .unwrap();
    let our_address = cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap();

    let sequence = &contact
        .get_account_info(
            cosmos_private_key
                .to_address(&contact.get_prefix())
                .unwrap(),
        )
        .await?
        .sequence;

    let eth_sign_msg = proto::DelegateKeysSignMsg {
        validator_address: our_valoper_address.clone(),
        nonce: *sequence,
    };
    let size = Message::encoded_len(&eth_sign_msg);
    let mut buf = BytesMut::with_capacity(size);
    Message::encode(&eth_sign_msg, &mut buf).expect("Failed to encode DelegateKeysSignMsg!");

    let eth_signature = eth_private_key.sign_ethereum_msg(&buf).to_bytes().to_vec();

    let msg_set_orch_address = proto::MsgDelegateKeys {
        validator_address: our_valoper_address.to_string(),
        orchestrator_address: delegate_cosmos_address.to_string(),
        ethereum_address: delegate_eth_address.to_string(),
        eth_signature,
    };

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000u64,
        granter: None,
        payer: None,
    };

    let msg = Msg::new("/gravity.v1.MsgDelegateKeys", msg_set_orch_address);

    let args = contact.get_message_args(our_address, fee).await?;

    let msg_bytes = cosmos_private_key.sign_std_msg(&[msg], args, MEMO)?;

    let mut txrpc = TxServiceClient::connect(contact.get_url()).await?;

    let response = txrpc
        .broadcast_tx(BroadcastTxRequest {
            tx_bytes: msg_bytes,
            mode: BroadcastMode::Block.into(),
        })
        .await?;

    let response = response.into_inner();

    contact
        .wait_for_tx(response.tx_response.unwrap(), TIMEOUT)
        .await
}

/// Send in a confirmation for an array of validator sets, it's far more efficient to send these
/// as a single message
#[allow(clippy::too_many_arguments)]
pub async fn send_valset_confirms(
    contact: &Contact,
    eth_private_key: EthPrivateKey,
    fee: Coin,
    valsets: Vec<Valset>,
    cosmos_private_key: CosmosPrivateKey,
    gravity_id: String,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap();

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64,
        granter: None,
        payer: None,
    };

    let messages = build::signer_set_tx_confirmations(
        contact,
        eth_private_key,
        valsets,
        cosmos_private_key,
        gravity_id,
    );

    let args = contact.get_message_args(our_address, fee).await?;

    let msg_bytes = cosmos_private_key.sign_std_msg(&messages, args, MEMO)?;

    let mut txrpc = TxServiceClient::connect(contact.get_url()).await?;

    let response = txrpc
        .broadcast_tx(BroadcastTxRequest {
            tx_bytes: msg_bytes,
            mode: BroadcastMode::Block.into(),
        })
        .await?;

    let response = response.into_inner();

    contact
        .wait_for_tx(response.tx_response.unwrap(), TIMEOUT)
        .await
}

/// Send in a confirmation for a specific transaction batch
pub async fn send_batch_confirm(
    contact: &Contact,
    eth_private_key: EthPrivateKey,
    fee: Coin,
    transaction_batches: Vec<TransactionBatch>,
    cosmos_private_key: CosmosPrivateKey,
    gravity_id: String,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap();

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64,
        granter: None,
        payer: None,
    };

    let messages = build::batch_tx_confirmations(
        contact,
        eth_private_key,
        transaction_batches,
        cosmos_private_key,
        gravity_id,
    );

    let args = contact.get_message_args(our_address, fee).await?;

    let msg_bytes = cosmos_private_key.sign_std_msg(&messages, args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;

    contact.wait_for_tx(response, TIMEOUT).await
}

/// Send in a confirmation for a specific logic call
pub async fn send_logic_call_confirm(
    contact: &Contact,
    eth_private_key: EthPrivateKey,
    fee: Coin,
    logic_calls: Vec<LogicCall>,
    cosmos_private_key: CosmosPrivateKey,
    gravity_id: String,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap();

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64,
        granter: None,
        payer: None,
    };

    let messages = build::contract_call_tx_confirmations(
        contact,
        eth_private_key,
        logic_calls,
        cosmos_private_key,
        gravity_id,
    );

    let args = contact.get_message_args(our_address, fee).await?;

    let msg_bytes = cosmos_private_key.sign_std_msg(&messages, args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;

    contact.wait_for_tx(response, TIMEOUT).await
}

pub async fn send_ethereum_claims(
    contact: &Contact,
    cosmos_private_key: CosmosPrivateKey,
    deposits: Vec<SendToCosmosEvent>,
    withdraws: Vec<TransactionBatchExecutedEvent>,
    erc20_deploys: Vec<Erc20DeployedEvent>,
    logic_calls: Vec<LogicCallExecutedEvent>,
    valsets: Vec<ValsetUpdatedEvent>,
    fee: Coin,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap();

    let messages = build::submit_ethereum_event_messages(
        contact,
        cosmos_private_key,
        deposits,
        withdraws,
        erc20_deploys,
        logic_calls,
        valsets,
    );

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64 * (messages.len() as u64),
        granter: None,
        payer: None,
    };

    let args = contact.get_message_args(our_address, fee).await?;

    let msg_bytes = cosmos_private_key.sign_std_msg(&messages, args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;

    contact.wait_for_tx(response, TIMEOUT).await
}

/// Sends tokens from Cosmos to Ethereum. These tokens will not be sent immediately instead
/// they will require some time to be included in a batch
pub async fn send_to_eth(
    cosmos_private_key: CosmosPrivateKey,
    destination: EthAddress,
    amount: Coin,
    fee: Coin,
    contact: &Contact,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap();
    if amount.denom != fee.denom {
        return Err(CosmosGrpcError::BadInput(format!(
            "{} {} is an invalid denom set for SendToEth you must pay fees in the same token your sending",
            amount.denom, fee.denom,
        )));
    }
    let balances = contact.get_balances(our_address).await.unwrap();
    let mut found = false;
    for balance in balances {
        if balance.denom == amount.denom {
            let total_amount = amount.amount.clone() + (fee.amount.clone() * 2u8.into());
            if balance.amount < total_amount {
                return Err(CosmosGrpcError::BadInput(format!(
                    "Insufficient balance of {} to send {}",
                    amount.denom, total_amount,
                )));
            }
            found = true;
        }
    }
    if !found {
        return Err(CosmosGrpcError::BadInput(format!(
            "No balance of {} to send",
            amount.denom,
        )));
    }

    let msg_send_to_eth = proto::MsgSendToEthereum {
        sender: our_address.to_string(),
        ethereum_recipient: destination.to_string(),
        amount: Some(amount.into()),
        bridge_fee: Some(fee.clone().into()),
    };

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000u64,
        granter: None,
        payer: None,
    };

    let msg = Msg::new("/gravity.v1.MsgSendToEthereum", msg_send_to_eth);

    let args = contact.get_message_args(our_address, fee).await?;

    let msg_bytes = cosmos_private_key.sign_std_msg(&[msg], args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;

    contact.wait_for_tx(response, TIMEOUT).await
}

pub async fn send_request_batch(
    cosmos_private_key: CosmosPrivateKey,
    denom: String,
    fee: Coin,
    contact: &Contact,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap();

    let msg_request_batch = proto::MsgRequestBatchTx {
        signer: our_address.to_string(),
        denom,
    };

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64,
        granter: None,
        payer: None,
    };

    let msg = Msg::new("/gravity.v1.MsgRequestBatchTx", msg_request_batch);

    let args = contact.get_message_args(our_address, fee).await?;

    let msg_bytes = cosmos_private_key.sign_std_msg(&[msg], args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;

    contact.wait_for_tx(response, TIMEOUT).await
}
