use bytes::BytesMut;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use ethers::utils::keccak256;
use gravity_proto::cosmos_sdk_proto::cosmos::base::abci::v1beta1::TxResponse;
use gravity_proto::gravity as proto;
use gravity_utils::error::GravityError;
use gravity_utils::ethereum::format_eth_address;
use ocular::chain::ChainContext;
use ocular::cosmrs::{AccountId, Coin, Denom};
use ocular::prelude::AccountInfo;
use ocular::tx::{FeeInfo, ModuleMsg, UnsignedTx};
use ocular::GrpcClient;
use ocular_somm_gravity::SommGravity;
use prost::Message;
use prost_types::Any;
use std::cmp;
use std::collections::HashSet;
use std::str::FromStr;
use std::thread::sleep;
use std::time::Instant;
use std::{result::Result, time::Duration};

pub const MEMO: &str = "Sent using Gravity Bridge Orchestrator";
pub const TIMEOUT: u64 = 60;
pub const OPERATOR_SUFFIX: &str = "valoper";

/// Send a transaction updating the eth address for the sending
/// Cosmos address. The sending Cosmos address should be a validator
pub async fn update_gravity_delegate_addresses(
    cosmos_client: &mut GrpcClient,
    signer: &AccountInfo,
    context: &ChainContext,
    delegate_eth_address: EthAddress,
    delegate_cosmos_address: AccountId,
    ethereum_wallet: LocalWallet,
    gas_price: (f64, String),
    gas_adjustment: f64,
) -> Result<TxResponse, GravityError> {
    let our_valoper_address = signer
        .address(&format!("{}{}", &context.prefix, OPERATOR_SUFFIX))
        .unwrap();

    let nonce = cosmos_client
        .query_account(&our_valoper_address)
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!(
                "Error getting account info for {}: {}",
                our_valoper_address, e
            ))
        })?
        .sequence;

    let eth_sign_msg = proto::DelegateKeysSignMsg {
        validator_address: our_valoper_address.clone(),
        nonce,
    };

    let mut data = BytesMut::with_capacity(eth_sign_msg.encoded_len());
    Message::encode(&eth_sign_msg, &mut data).expect("encoding failed");

    let data_hash = keccak256(data);
    let eth_signature = ethereum_wallet.sign_message(data_hash).await?.to_vec();
    let msg = SommGravity::SetDelegateKeys {
        validator_address: &our_valoper_address,
        orchestrator_address: delegate_cosmos_address.as_ref(),
        ethereum_address: &delegate_eth_address.to_string(),
        eth_signature,
    }
    .into_any()
    .map_err(|e| {
        GravityError::CosmosGrpcError(format!("failed to encode RequestBatchTx: {:?}", e))
    })?;

    send_messages(
        cosmos_client,
        vec![msg],
        signer,
        context,
        gas_price,
        gas_adjustment,
    )
    .await
}

/// Sends tokens from Cosmos to Ethereum. These tokens will not be sent immediately instead
/// they will require some time to be included in a batch
pub async fn send_to_eth(
    cosmos_client: &mut GrpcClient,
    signer: &AccountInfo,
    context: &ChainContext,
    destination: EthAddress,
    amount: Coin,
    bridge_fee: Coin,
    gas_price: (f64, String),
    gas_adjustment: f64,
) -> Result<TxResponse, GravityError> {
    let cosmos_address = signer.address(&context.prefix).unwrap();
    if amount.denom != bridge_fee.denom {
        return Err(GravityError::CosmosGrpcError(format!(
            "The amount ({}) and bridge_fee ({}) denominations do not match.",
            amount.denom, bridge_fee.denom,
        )));
    }

    let msg = SommGravity::SendToEthereum {
        sender: &cosmos_address,
        ethereum_recipient: &format_eth_address(destination),
        amount,
        bridge_fee,
    }
    .into_any()
    .map_err(|e| {
        GravityError::CosmosGrpcError(format!("failed to encode RequestBatchTx: {:?}", e))
    })?;

    send_messages(
        cosmos_client,
        vec![msg],
        signer,
        context,
        gas_price,
        gas_adjustment,
    )
    .await
}

pub async fn send_request_batch_tx(
    cosmos_client: &mut GrpcClient,
    signer: &AccountInfo,
    context: &ChainContext,
    denom: String,
    gas_price: (f64, String),
    gas_adjustment: f64,
) -> Result<TxResponse, GravityError> {
    let cosmos_address = signer.address(&context.prefix).unwrap();
    let msg = SommGravity::RequestBatchTx {
        denom: &denom,
        signer: &cosmos_address,
    }
    .into_any()
    .map_err(|e| {
        GravityError::CosmosGrpcError(format!("failed to encode RequestBatchTx: {:?}", e))
    })?;

    send_messages(
        cosmos_client,
        vec![msg],
        signer,
        context,
        gas_price,
        gas_adjustment,
    )
    .await
}

pub async fn wait_for_tx(
    cosmos_client: &mut GrpcClient,
    hash: &str,
) -> Result<TxResponse, GravityError> {
    let start = Instant::now();
    loop {
        let response = cosmos_client.query_tx_by_hash(hash).await.map_err(|e| {
            GravityError::CosmosGrpcError(format!("failed to query transaction: {}", e))
        })?;

        if response.tx.is_some() {
            return Ok(response.tx_response.unwrap());
        }

        let now = Instant::now();
        if now.checked_duration_since(start).unwrap().as_secs() >= TIMEOUT {
            return Err(GravityError::CosmosGrpcError(
                "timed out waiting for transaction to be included in a block".to_string(),
            ));
        }

        sleep(Duration::from_secs(1));
    }
}

pub async fn send_messages(
    cosmos_client: &mut GrpcClient,
    messages: Vec<Any>,
    signer: &AccountInfo,
    context: &ChainContext,
    gas_price: (f64, String),
    gas_adjustment: f64,
) -> Result<TxResponse, GravityError> {
    let fee_amount = Coin {
        denom: Denom::from_str(gas_price.1.as_str()).unwrap(),
        amount: 0u32.into(),
    };
    let mut fee_info = FeeInfo::new(fee_amount);
    let tx: UnsignedTx = (&messages).into();
    let sim_tx = tx
        .clone()
        .sign(signer, fee_info.clone(), context, cosmos_client)
        .await
        .map_err(|e| GravityError::CosmosGrpcError(format!("failed to sign transaction: {}", e)))?;
    let gas = cosmos_client
        .simulate(sim_tx)
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!("failed to simulate transaction: {}", e))
        })?
        .gas_info
        .unwrap();

    // multiply the estimated gas by the configured gas adjustment
    let gas_limit: f64 = (gas.gas_used as f64) * gas_adjustment;
    fee_info.gas_limit(cmp::max(gas_limit as u64, 500000 * messages.len() as u64));

    // compute the fee as fee=ceil(gas_limit * gas_price)
    let fee_amount = (gas_limit as f64 * gas_price.0).abs().ceil() as u128;
    fee_info.amount(fee_amount);

    let tx = tx
        .sign(signer, fee_info, context, cosmos_client)
        .await
        .map_err(|e| GravityError::CosmosGrpcError(format!("failed to sign transaction: {}", e)))?;

    // we block on the broadcast while we wait for DeliverTx to complete
    let response = cosmos_client
        .broadcast_commit(tx)
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!("failed to broadcast transaction: {}", e))
        })?
        .tx_response
        .unwrap();

    if response.code != 0 {
        return Err(GravityError::CosmosGrpcError(format!(
            "non-zero response code from tx broadcast. code: {}, logs: {:?}",
            response.code, response.logs
        )));
    }

    // we wait in case the tx was somehow not included in a block
    wait_for_tx(cosmos_client, &response.txhash).await
}

pub async fn send_main_loop(
    cosmos_client: GrpcClient,
    signer: &AccountInfo,
    context: &ChainContext,
    gas_price: (f64, String),
    gas_adjustment: f64,
    mut rx: tokio::sync::mpsc::Receiver<Vec<Any>>,
    msg_batch_size: usize,
) {
    let mut cosmos_client = cosmos_client;
    while let Some(messages) = rx.recv().await {
        for msg_chunk in messages.chunks(msg_batch_size) {
            match send_messages(
                &mut cosmos_client,
                msg_chunk.to_vec(),
                signer,
                context,
                gas_price.to_owned(),
                gas_adjustment,
            )
            .await
            {
                Ok(res) => trace!("okay: {:?}", res),
                Err(err) => {
                    let msg_types = msg_chunk
                        .iter()
                        .map(|msg| msg.clone().type_url)
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
