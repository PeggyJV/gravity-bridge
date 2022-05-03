use deep_space::address::Address;
use ethers::types::Address as EthAddress;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_proto::gravity::*;
use gravity_utils::error::GravityError;
use gravity_utils::ethereum::format_eth_address;
use gravity_utils::types::*;
use tonic::transport::Channel;

/// get the valset for a given nonce (block) height
pub async fn get_valset(
    client: &mut GravityQueryClient<Channel>,
    nonce: u64,
) -> Result<Option<Valset>, GravityError> {
    let response = client
        .signer_set_tx(SignerSetTxRequest {
            signer_set_nonce: nonce,
        })
        .await?;
    let valset = response.into_inner().signer_set.map(Into::into);
    Ok(valset)
}

/// This hits the /pending_valset_requests endpoint and will provide
/// an array of validator sets we have not already signed
pub async fn get_oldest_unsigned_valsets(
    client: &mut GravityQueryClient<Channel>,
    address: Address,
) -> Result<Vec<Valset>, GravityError> {
    let response = client
        .unsigned_signer_set_txs(UnsignedSignerSetTxsRequest {
            address: address.to_string(),
        })
        .await?;
    let valsets = response.into_inner().signer_sets;
    // convert from proto valset type to rust valset type
    let valsets = valsets.iter().map(|v| v.clone().into()).collect();
    Ok(valsets)
}

/// this input views the last five signer set txs that have been made, useful if you're
/// a relayer looking to ferry confirmations
pub async fn get_latest_valset(
    client: &mut GravityQueryClient<Channel>,
) -> Result<Option<Valset>, GravityError> {
    let response = client
        .latest_signer_set_tx(LatestSignerSetTxRequest {})
        .await?;
    let valset = response.into_inner().signer_set.map(Into::into);
    Ok(valset)
}

/// get all valset confirmations for a given nonce
pub async fn get_all_valset_confirms(
    client: &mut GravityQueryClient<Channel>,
    nonce: u64,
) -> Result<Vec<ValsetConfirmResponse>, GravityError> {
    let request = client
        .signer_set_tx_confirmations(SignerSetTxConfirmationsRequest {
            signer_set_nonce: nonce,
        })
        .await?;
    let confirms = request.into_inner().signatures;
    let mut parsed_confirms = Vec::new();
    for item in confirms {
        parsed_confirms.push(ValsetConfirmResponse::from_proto(item)?)
    }
    Ok(parsed_confirms)
}

pub async fn get_oldest_unsigned_transaction_batch(
    client: &mut GravityQueryClient<Channel>,
    address: Address,
) -> Result<Option<TransactionBatch>, GravityError> {
    let request = client
        .unsigned_batch_txs(UnsignedBatchTxsRequest {
            address: address.to_string(),
        })
        .await?;
    let batches = extract_valid_batches(request.into_inner().batches);
    let batch = batches.get(0);
    match batch {
        Some(batch) => Ok(Some(batch.clone())),
        None => Ok(None),
    }
}

/// gets the latest 100 transaction batches, regardless of token type
/// for relayers to consider relaying
pub async fn get_latest_transaction_batches(
    client: &mut GravityQueryClient<Channel>,
) -> Result<Vec<TransactionBatch>, GravityError> {
    let request = client
        .batch_txs(BatchTxsRequest { pagination: None })
        .await?;
    Ok(extract_valid_batches(request.into_inner().batches))
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
    client: &mut GravityQueryClient<Channel>,
    nonce: u64,
    contract_address: EthAddress,
) -> Result<Vec<BatchConfirmResponse>, GravityError> {
    let request = client
        .batch_tx_confirmations(BatchTxConfirmationsRequest {
            batch_nonce: nonce,
            token_contract: format_eth_address(contract_address),
        })
        .await?;
    let batch_confirms = request.into_inner().signatures;
    let mut out = Vec::new();
    for confirm in batch_confirms {
        out.push(BatchConfirmResponse::from_proto(confirm)?)
    }
    Ok(out)
}

/// Gets the last event nonce that a given validator has attested to, this lets us
/// catch up with what the current event nonce should be if a oracle is restarted
pub async fn get_last_event_nonce(
    client: &mut GravityQueryClient<Channel>,
    address: Address,
) -> Result<u64, GravityError> {
    let request = client
        .last_submitted_ethereum_event(LastSubmittedEthereumEventRequest {
            address: address.to_string(),
        })
        .await?;
    Ok(request.into_inner().event_nonce)
}

/// Gets the 100 latest logic calls for a relayer to consider relaying
pub async fn get_latest_logic_calls(
    client: &mut GravityQueryClient<Channel>,
) -> Result<Vec<LogicCall>, GravityError> {
    let request = client
        .contract_call_txs(ContractCallTxsRequest { pagination: None })
        .await?;
    let calls = request.into_inner().calls;
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
    client: &mut GravityQueryClient<Channel>,
    invalidation_scope: Vec<u8>,
    invalidation_nonce: u64,
) -> Result<Vec<LogicCallConfirmResponse>, GravityError> {
    let request = client
        .contract_call_tx_confirmations(ContractCallTxConfirmationsRequest {
            invalidation_scope,
            invalidation_nonce,
        })
        .await?;
    let call_confirms = request.into_inner().signatures;
    let mut out = Vec::new();
    for confirm in call_confirms {
        out.push(LogicCallConfirmResponse::from_proto(confirm)?)
    }
    Ok(out)
}

pub async fn get_oldest_unsigned_logic_call(
    client: &mut GravityQueryClient<Channel>,
    address: Address,
) -> Result<Vec<LogicCall>, GravityError> {
    let request = client
        .unsigned_contract_call_txs(UnsignedContractCallTxsRequest {
            address: address.to_string(),
        })
        .await?;
    let calls = request.into_inner().calls;
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
