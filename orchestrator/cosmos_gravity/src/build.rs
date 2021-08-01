use bytes::BytesMut;
use clarity::PrivateKey as EthPrivateKey;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use deep_space::utils::bytes_to_hex_str;
use deep_space::Contact;
use deep_space::Msg;
use ethereum_gravity::utils::downcast_uint256;
use gravity_proto::gravity as proto;
use gravity_utils::message_signatures::{
    encode_logic_call_confirm, encode_tx_batch_confirm, encode_valset_confirm,
};
use gravity_utils::types::*;
use prost::Message;
use prost_types::Any;

pub fn signer_set_tx_confirmation_messages(
    contact: &Contact,
    ethereum_key: EthPrivateKey,
    valsets: Vec<Valset>,
    cosmos_key: CosmosPrivateKey,
    gravity_id: String,
) -> Vec<Msg> {
    let our_cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let our_eth_address = ethereum_key.to_public_key().unwrap();

    let mut messages = Vec::new();
    for valset in valsets {
        let message = encode_valset_confirm(gravity_id.clone(), valset.clone());
        let eth_signature = ethereum_key.sign_ethereum_msg(&message);
        let confirm = proto::SignerSetTxConfirmation {
            ethereum_signer: our_eth_address.to_string(),
            signer_set_nonce: valset.nonce,
            signature: eth_signature.to_bytes().to_vec(),
        };
        let size = Message::encoded_len(&confirm);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&confirm, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough

        let wrapper = proto::MsgSubmitEthereumTxConfirmation {
            signer: our_cosmos_address.to_string(),
            confirmation: Some(Any {
                type_url: "/gravity.v1.SignerSetTxConfirmation".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumTxConfirmation", wrapper);
        messages.push(msg);
    }
    messages
}

pub fn batch_tx_confirmation_messages(
    contact: &Contact,
    ethereum_key: EthPrivateKey,
    transaction_batches: Vec<TransactionBatch>,
    cosmos_key: CosmosPrivateKey,
    gravity_id: String,
) -> Vec<Msg> {
    let our_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let our_eth_address = ethereum_key.to_public_key().unwrap();

    let mut messages = Vec::new();
    for batch in transaction_batches {
        info!("Submitting signature for batch {:?}", batch);
        let message = encode_tx_batch_confirm(gravity_id.clone(), batch.clone());
        let eth_signature = ethereum_key.sign_ethereum_msg(&message);
        info!(
            "Sending batch update address {} sig {} hash {}",
            our_eth_address,
            bytes_to_hex_str(&eth_signature.to_bytes()),
            bytes_to_hex_str(&message),
        );
        let confirm = proto::BatchTxConfirmation {
            token_contract: batch.token_contract.to_string(),
            batch_nonce: batch.nonce,
            ethereum_signer: our_eth_address.to_string(),
            signature: eth_signature.to_bytes().to_vec(),
        };
        let size = Message::encoded_len(&confirm);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&confirm, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.BatchTxConfirmation".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumTxConfirmation", wrapper);
        messages.push(msg);
    }
    messages
}

pub fn contract_call_tx_confirmation_messages(
    contact: &Contact,
    ethereum_key: EthPrivateKey,
    logic_calls: Vec<LogicCall>,
    cosmos_key: CosmosPrivateKey,
    gravity_id: String,
) -> Vec<Msg> {
    let our_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

    let our_eth_address = ethereum_key.to_public_key().unwrap();

    let mut messages = Vec::new();
    for call in logic_calls {
        let message = encode_logic_call_confirm(gravity_id.clone(), call.clone());
        let eth_signature = ethereum_key.sign_ethereum_msg(&message);
        let confirm = proto::ContractCallTxConfirmation {
            ethereum_signer: our_eth_address.to_string(),
            signature: eth_signature.to_bytes().to_vec(),
            invalidation_scope: bytes_to_hex_str(&call.invalidation_id).as_bytes().to_vec(),
            invalidation_nonce: call.invalidation_nonce,
        };
        let size = Message::encoded_len(&confirm);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&confirm, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumTxConfirmation {
            signer: our_address.to_string(),
            confirmation: Some(Any {
                type_url: "/gravity.v1.ContractCallTxConfirmation".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumTxConfirmation", wrapper);
        messages.push(msg);
    }
    messages
}

pub fn submit_ethereum_event_messages(
    contact: &Contact,
    cosmos_key: CosmosPrivateKey,
    deposits: Vec<SendToCosmosEvent>,
    withdraws: Vec<TransactionBatchExecutedEvent>,
    erc20_deploys: Vec<Erc20DeployedEvent>,
    logic_calls: Vec<LogicCallExecutedEvent>,
    valsets: Vec<ValsetUpdatedEvent>,
) -> Vec<Msg> {
    let our_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

    // This sorts oracle messages by event nonce before submitting them. It's not a pretty implementation because
    // we're missing an intermediary layer of abstraction. We could implement 'EventTrait' and then implement sort
    // for it, but then when we go to transform 'EventTrait' objects into GravityMsg enum values we'll have all sorts
    // of issues extracting the inner object from the TraitObject. Likewise we could implement sort of GravityMsg but that
    // would require a truly horrendous (nearly 100 line) match statement to deal with all combinations. That match statement
    // could be reduced by adding two traits to sort against but really this is the easiest option.
    //
    // We index the events by event nonce in an unordered hashmap and then play them back in order into a vec
    let mut unordered_msgs = std::collections::HashMap::new();
    for deposit in deposits {
        let event = proto::SendToCosmosEvent {
            event_nonce: downcast_uint256(deposit.event_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(deposit.block_height).unwrap(),
            token_contract: deposit.erc20.to_string(),
            amount: deposit.amount.to_string(),
            cosmos_receiver: deposit.destination.to_string(),
            ethereum_sender: deposit.sender.to_string(),
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.SendToCosmosEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(deposit.event_nonce, msg);
    }
    for withdraw in withdraws {
        let event = proto::BatchExecutedEvent {
            event_nonce: downcast_uint256(withdraw.event_nonce.clone()).unwrap(),
            batch_nonce: downcast_uint256(withdraw.batch_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(withdraw.block_height).unwrap(),
            token_contract: withdraw.erc20.to_string(),
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.BatchExecutedEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(withdraw.event_nonce, msg);
    }
    for deploy in erc20_deploys {
        let event = proto::Erc20DeployedEvent {
            event_nonce: downcast_uint256(deploy.event_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(deploy.block_height).unwrap(),
            cosmos_denom: deploy.cosmos_denom,
            token_contract: deploy.erc20_address.to_string(),
            erc20_name: deploy.name,
            erc20_symbol: deploy.symbol,
            erc20_decimals: deploy.decimals as u64,
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.ERC20DeployedEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(deploy.event_nonce, msg);
    }
    for call in logic_calls {
        let event = proto::ContractCallExecutedEvent {
            event_nonce: downcast_uint256(call.event_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(call.block_height).unwrap(),
            invalidation_id: call.invalidation_id,
            invalidation_nonce: downcast_uint256(call.invalidation_nonce).unwrap(),
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.ContractCallExecutedEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(call.event_nonce, msg);
    }
    for valset in valsets {
        let event = proto::SignerSetTxExecutedEvent {
            event_nonce: downcast_uint256(valset.event_nonce.clone()).unwrap(),
            signer_set_tx_nonce: downcast_uint256(valset.valset_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(valset.block_height).unwrap(),
            members: valset.members.iter().map(|v| v.into()).collect(),
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.SignerSetTxExecutedEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(valset.event_nonce, msg);
    }

    let mut keys = Vec::new();
    for (key, _) in unordered_msgs.iter() {
        keys.push(key.clone());
    }
    keys.sort();

    let mut messages = Vec::new();
    for i in keys {
        messages.push(unordered_msgs.remove_entry(&i).unwrap().1);
    }
    messages
}
