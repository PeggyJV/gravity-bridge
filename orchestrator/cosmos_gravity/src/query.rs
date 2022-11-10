use ethers::types::Address as EthAddress;
use gravity_proto::gravity::*;
use gravity_utils::error::GravityError;
use gravity_utils::ethereum::format_eth_address;
use gravity_utils::types::*;
use ocular::cosmrs::AccountId;
use ocular::GrpcClient;
use ocular_somm_gravity::SommGravityExt;

/// get the valset for a given nonce (block) height
pub async fn get_valset(
    client: &mut GrpcClient,
    nonce: u64,
) -> Result<Option<Valset>, GravityError> {
    Ok(client
        .query_signer_set_tx(nonce)
        .await
        .map_err(|e| GravityError::CosmosGrpcError(format!("failed to query signer set: {:?}", e)))?
        .signer_set
        .map(Into::into))
}

/// This hits the /pending_valset_requests endpoint and will provide
/// an array of validator sets we have not already signed
pub async fn get_oldest_unsigned_valsets(
    client: &mut GrpcClient,
    address: &AccountId,
) -> Result<Vec<Valset>, GravityError> {
    let valsets = client
        .query_unsigned_signer_set_txs(address.as_ref())
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!("failed to query unsigned signer sets: {:?}", e))
        })?
        .signer_sets;
    // convert from proto valset type to rust valset type
    let valsets = valsets.iter().map(|v| v.clone().into()).collect();
    Ok(valsets)
}

/// this input views the last five signer set txs that have been made, useful if you're
/// a relayer looking to ferry confirmations
pub async fn get_latest_valset(client: &mut GrpcClient) -> Result<Option<Valset>, GravityError> {
    let valset = client
        .query_latest_signer_set_tx()
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!("failed to query latest signer set: {:?}", e))
        })?
        .signer_set
        .map(Into::into);
    Ok(valset)
}

/// get all valset confirmations for a given nonce
pub async fn get_all_valset_confirms(
    client: &mut GrpcClient,
    nonce: u64,
) -> Result<Vec<ValsetConfirmResponse>, GravityError> {
    let confirms = client
        .query_signer_set_tx_confirmations(nonce)
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!(
                "failed to query signer set tx confirmations: {:?}",
                e
            ))
        })?
        .signatures;
    let mut parsed_confirms = Vec::new();
    for item in confirms {
        parsed_confirms.push(ValsetConfirmResponse::from_proto(item)?)
    }
    Ok(parsed_confirms)
}

pub async fn get_oldest_unsigned_transaction_batch(
    client: &mut GrpcClient,
    address: &AccountId,
) -> Result<Option<TransactionBatch>, GravityError> {
    let batches = client
        .query_unsigned_batch_txs(address.as_ref())
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!("failed to query unsigned batch txs: {:?}", e))
        })?
        .batches;
    let batches = extract_valid_batches(batches);
    let batch = batches.get(0);
    match batch {
        Some(batch) => Ok(Some(batch.clone())),
        None => Ok(None),
    }
}

/// gets the latest 100 transaction batches, regardless of token type
/// for relayers to consider relaying
pub async fn get_latest_transaction_batches(
    client: &mut GrpcClient,
) -> Result<Vec<TransactionBatch>, GravityError> {
    let batches = client
        .query_batch_txs(None)
        .await
        .map_err(|e| GravityError::CosmosGrpcError(format!("failed to query batch txs: {:?}", e)))?
        .batches;
    Ok(extract_valid_batches(batches))
}

// If we can't serialize a batch from a proto, but it was committed to the chain,
// we should just ignore it. It should eventually time out and be canceled.
fn extract_valid_batches(batches: Vec<BatchTx>) -> Vec<TransactionBatch> {
    let mut valid_batches = Vec::new();
    for batch in batches {
        match TransactionBatch::from_proto(batch.clone()) {
            Ok(valid_batch) => valid_batches.push(valid_batch),
            Err(e) => warn!("{}, skipping invalid batch: {:?}", e, batch),
        }
    }
    valid_batches
}

/// get all batch confirmations for a given nonce and denom
pub async fn get_transaction_batch_signatures(
    client: &mut GrpcClient,
    nonce: u64,
    contract_address: EthAddress,
) -> Result<Vec<BatchConfirmResponse>, GravityError> {
    let batch_confirms = client
        .query_batch_tx_confirmations(nonce, &format_eth_address(contract_address))
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!(
                "failed to query transaction batch signatures: {:?}",
                e
            ))
        })?
        .signatures;
    let mut out = Vec::new();
    for confirm in batch_confirms {
        out.push(BatchConfirmResponse::from_proto(confirm)?)
    }
    Ok(out)
}

/// Gets the last event nonce that a given validator has attested to, this lets us
/// catch up with what the current event nonce should be if a oracle is restarted
pub async fn get_last_event_nonce(
    client: &mut GrpcClient,
    address: &AccountId,
) -> Result<u64, GravityError> {
    Ok(client
        .query_last_submitted_ethereum_event(address.as_ref())
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!("failed to query last submitted event: {:?}", e))
        })?
        .event_nonce)
}

/// Gets the 100 latest logic calls for a relayer to consider relaying
pub async fn get_latest_logic_calls(
    client: &mut GrpcClient,
) -> Result<Vec<LogicCall>, GravityError> {
    let calls = client
        .query_contract_call_txs(None)
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!("failed to query contract call txs: {:?}", e))
        })?
        .calls;
    let mut out = Vec::new();
    for call in calls {
        out.push(LogicCall::from_proto(call)?);
    }

    // as these calls are expected to be in oldest -> newest order, but
    // the chain does not provide them as such, we will sort using the
    // invalidation nonces as keys such that for any given scope, calls will
    // be processed in nonce order
    out.sort_by_key(|call| call.invalidation_nonce);

    Ok(out)
}

pub async fn get_logic_call_signatures(
    client: &mut GrpcClient,
    invalidation_scope: Vec<u8>,
    invalidation_nonce: u64,
) -> Result<Vec<LogicCallConfirmResponse>, GravityError> {
    let call_confirms = client
        .query_contract_call_tx_confirmations(invalidation_scope, invalidation_nonce)
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!(
                "failed to query contract call tx confirmations: {:?}",
                e
            ))
        })?
        .signatures;
    let mut out = Vec::new();
    for confirm in call_confirms {
        out.push(LogicCallConfirmResponse::from_proto(confirm)?)
    }
    Ok(out)
}

pub async fn get_oldest_unsigned_logic_call(
    client: &mut GrpcClient,
    address: &AccountId,
) -> Result<Vec<LogicCall>, GravityError> {
    let calls = client
        .query_unsigned_contract_call_txs(address.as_ref())
        .await
        .map_err(|e| {
            GravityError::CosmosGrpcError(format!(
                "failed to query unsigned contract call txs: {:?}",
                e
            ))
        })?
        .calls;
    let mut out = Vec::new();
    for call in calls {
        out.push(LogicCall::from_proto(call)?)
    }
    Ok(out)
}

#[test]
fn extract_valid_batches_test() {
    let erc20_addr = "0x0635FF793Edf48cf5dB294916720A78e6e490E40".to_string();
    let token_contract = "0xC26eFfa98B8A2632141562Ae7E34953Cfe5B4888".to_string();
    let transactions = vec![SendToEthereum {
        id: 1,
        sender: "cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7".to_string(),
        ethereum_recipient: "0x64D110e00064F2b428476cD64295d8E35836ffd6".to_string(),
        erc20_token: Some(gravity_proto::gravity::Erc20Token {
            contract: erc20_addr.clone(),
            amount: "1".to_string(),
        }),
        erc20_fee: Some(gravity_proto::gravity::Erc20Token {
            contract: erc20_addr,
            amount: "1".to_string(),
        }),
    }];

    let valid_batch = BatchTx {
        batch_nonce: 1,
        timeout: 3,
        transactions,
        token_contract: token_contract.clone(),
        height: 2,
    };

    let invalid_batch = BatchTx {
        batch_nonce: 2,
        timeout: 3,
        transactions: Vec::new(),
        token_contract,
        height: 2,
    };

    let valid_batches = extract_valid_batches(vec![valid_batch, invalid_batch.clone()]);
    assert_eq!(valid_batches.len(), 1);
    assert_eq!(valid_batches.get(0).unwrap().nonce, 1);

    let should_be_empty = extract_valid_batches(vec![invalid_batch]);
    assert_eq!(should_be_empty.len(), 0);
}
