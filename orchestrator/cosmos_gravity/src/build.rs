use ethereum_gravity::types::EthClient;
use ethers::prelude::*;
use ethers::utils::keccak256;
use gravity_proto::gravity as proto;
use gravity_proto::ToAny;
use gravity_utils::ethereum::{downcast_to_u64, format_eth_address};
use gravity_utils::message_signatures::{
    encode_logic_call_confirm, encode_tx_batch_confirm, encode_valset_confirm,
};
use gravity_utils::types::*;
use ocular::cosmrs::AccountId;
use ocular::tx::ModuleMsg;
use ocular_somm_gravity::SommGravity;
use prost_types::Any;
use std::collections::BTreeMap;

pub async fn signer_set_tx_confirmation_messages(
    cosmos_account: &AccountId,
    eth_client: EthClient,
    valsets: Vec<Valset>,
    gravity_id: String,
) -> Vec<Any> {
    let cosmos_address = cosmos_account.as_ref();
    let ethereum_address = eth_client.address();

    let mut msgs = Vec::new();
    for valset in valsets {
        let data = keccak256(encode_valset_confirm(gravity_id.clone(), valset.clone()).as_slice());
        // Signer trait responds with a Result, but we use a LocalWallet and it
        // will never throw an error
        let signature = eth_client.signer().sign_message(data).await.unwrap();
        let confirmation = SommGravity::SignerSetTxConfirmation {
            ethereum_signer: &format_eth_address(ethereum_address),
            signer_set_nonce: valset.nonce,
            signature: signature.into(),
        }
        .into_any()
        .unwrap();
        let msg = SommGravity::SubmitEthereumTxConfirmation {
            signer: cosmos_address,
            confirmation,
        }
        .into_any()
        .unwrap();
        msgs.push(msg);
    }

    msgs
}

pub async fn batch_tx_confirmation_messages(
    cosmos_account: &AccountId,
    eth_client: EthClient,
    batches: Vec<TransactionBatch>,
    gravity_id: String,
) -> Vec<Any> {
    let cosmos_address = cosmos_account.as_ref();
    let ethereum_address = eth_client.address();

    let mut msgs = Vec::new();
    for batch in batches {
        let data = keccak256(encode_tx_batch_confirm(gravity_id.clone(), batch.clone()).as_slice());
        // Signer trait responds with a Result, but we use a LocalWallet and it
        // will never throw an error
        let signature = eth_client.signer().sign_message(data).await.unwrap();
        let confirmation = SommGravity::BatchTxConfirmation {
            ethereum_signer: &format_eth_address(ethereum_address),
            token_contract_address: &format_eth_address(batch.token_contract),
            batch_nonce: batch.nonce,
            signature: signature.into(),
        }
        .into_any()
        .unwrap();
        let msg = SommGravity::SubmitEthereumTxConfirmation {
            signer: cosmos_address,
            confirmation,
        }
        .into_any()
        .unwrap();
        msgs.push(msg);
    }

    msgs
}

pub async fn contract_call_tx_confirmation_messages(
    cosmos_account: &AccountId,
    eth_client: EthClient,
    logic_calls: Vec<LogicCall>,
    gravity_id: String,
) -> Vec<Any> {
    let cosmos_address = cosmos_account.as_ref();
    let ethereum_address = eth_client.address();

    let mut msgs = Vec::new();
    for logic_call in logic_calls {
        let data =
            keccak256(encode_logic_call_confirm(gravity_id.clone(), logic_call.clone()).as_slice());
        // Signer trait responds with a Result, but we use a LocalWallet and it
        // will never throw an error
        let signature = eth_client.signer().sign_message(data).await.unwrap();
        let confirmation = SommGravity::ContractCallTxConfirmation {
            ethereum_signer: &format_eth_address(ethereum_address),
            invalidation_scope: logic_call.invalidation_id,
            invalidation_nonce: logic_call.invalidation_nonce,
            signature: signature.into(),
        }
        .into_any()
        .unwrap();
        let msg = SommGravity::SubmitEthereumTxConfirmation {
            signer: cosmos_address,
            confirmation,
        }
        .into_any()
        .unwrap();
        msgs.push(msg);
    }

    msgs
}

pub async fn ethereum_vote_height_messages(
    cosmos_address: &AccountId,
    ethereum_height: u64,
) -> Vec<Any> {
    let msg = SommGravity::SubmitEthereumHeightVote {
        ethereum_height,
        signer: cosmos_address.as_ref(),
    }
    .into_any()
    .unwrap();

    vec![msg]
}

pub fn ethereum_event_messages(
    cosmos_address: &AccountId,
    deposits: Vec<SendToCosmosEvent>,
    batches: Vec<TransactionBatchExecutedEvent>,
    erc20_deploys: Vec<Erc20DeployedEvent>,
    logic_calls: Vec<LogicCallExecutedEvent>,
    valsets: Vec<ValsetUpdatedEvent>,
) -> Vec<Any> {
    // This sorts oracle messages by event nonce before submitting them. It's not a pretty implementation because
    // we're missing an intermediary layer of abstraction. We could implement 'EventTrait' and then implement sort
    // for it, but then when we go to transform 'EventTrait' objects into GravityMsg enum values we'll have all sorts
    // of issues extracting the inner object from the TraitObject. Likewise we could implement sort of GravityMsg but that
    // would require a truly horrendous (nearly 100 line) match statement to deal with all combinations. That match statement
    // could be reduced by adding two traits to sort against but really this is the easiest option.
    //
    // We index the events by event nonce in an unordered hashmap and then play them back in order into a vec
    let mut unordered_msgs = BTreeMap::new();
    for deposit in deposits {
        let event = proto::SendToCosmosEvent {
            event_nonce: downcast_to_u64(deposit.event_nonce).unwrap(),
            ethereum_height: downcast_to_u64(deposit.block_height).unwrap(),
            token_contract: format_eth_address(deposit.erc20),
            amount: deposit.amount.to_string(),
            cosmos_receiver: deposit.destination.to_string(),
            ethereum_sender: format_eth_address(deposit.sender),
        };
        let msg = encode_submit_ethereum_event(cosmos_address.as_ref(), event);
        unordered_msgs.insert(deposit.event_nonce, msg);
    }
    for batch in batches {
        let event = proto::BatchExecutedEvent {
            event_nonce: downcast_to_u64(batch.event_nonce).unwrap(),
            batch_nonce: downcast_to_u64(batch.batch_nonce).unwrap(),
            ethereum_height: downcast_to_u64(batch.block_height).unwrap(),
            token_contract: format_eth_address(batch.erc20),
        };
        let msg = encode_submit_ethereum_event(cosmos_address.as_ref(), event);
        unordered_msgs.insert(batch.event_nonce, msg);
    }
    for deploy in erc20_deploys {
        let event = proto::Erc20DeployedEvent {
            event_nonce: downcast_to_u64(deploy.event_nonce).unwrap(),
            ethereum_height: downcast_to_u64(deploy.block_height).unwrap(),
            cosmos_denom: deploy.cosmos_denom,
            token_contract: format_eth_address(deploy.erc20_address),
            erc20_name: deploy.name,
            erc20_symbol: deploy.symbol,
            erc20_decimals: deploy.decimals as u64,
        };
        let msg = encode_submit_ethereum_event(cosmos_address.as_ref(), event);
        unordered_msgs.insert(deploy.event_nonce, msg);
    }
    for logic_call in logic_calls {
        let event = proto::ContractCallExecutedEvent {
            event_nonce: downcast_to_u64(logic_call.event_nonce).unwrap(),
            ethereum_height: downcast_to_u64(logic_call.block_height).unwrap(),
            invalidation_scope: logic_call.invalidation_id,
            invalidation_nonce: downcast_to_u64(logic_call.invalidation_nonce).unwrap(),
        };
        let msg = encode_submit_ethereum_event(cosmos_address.as_ref(), event);
        unordered_msgs.insert(logic_call.event_nonce, msg);
    }
    for valset in valsets {
        // note that SignerSetTxExecutedEvent does not include reward amount or
        // reward token, which is fine since we are not actually using them at the
        // moment, but it is part of the contract-defined event
        let event = proto::SignerSetTxExecutedEvent {
            event_nonce: downcast_to_u64(valset.event_nonce).unwrap(),
            signer_set_tx_nonce: downcast_to_u64(valset.valset_nonce).unwrap(),
            ethereum_height: downcast_to_u64(valset.block_height).unwrap(),
            members: valset.members.iter().map(|v| v.into()).collect(),
        };
        let msg = encode_submit_ethereum_event(cosmos_address.as_ref(), event);
        unordered_msgs.insert(valset.event_nonce, msg);
    }

    let mut msgs = Vec::new();
    for (i, _) in unordered_msgs.clone().iter() {
        msgs.push(unordered_msgs.remove_entry(i).unwrap().1);
    }

    msgs
}

fn encode_submit_ethereum_event<T>(cosmos_address: &str, event: T) -> Any
where
    T: ToAny + prost::Message,
{
    let msg = proto::MsgSubmitEthereumEvent {
        signer: cosmos_address.to_string(),
        event: event.to_any(),
    };
    let mut any = Any::default();
    prost::Message::encode(&msg, &mut any.value).expect("failed to encode MsgSubmitEthereumEvent");
    any.type_url = "/gravity.v1.MsgSubmitEthereumEvent".to_string();

    any
}
