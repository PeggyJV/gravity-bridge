use crate::{
    types::{EthClient, EthSignerMiddleware},
    utils::{get_gas_price, get_logic_call_nonce, GasCost},
};
use ethers::contract::builders::ContractCall;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_abi::gravity::*;
use gravity_utils::ethereum::{bytes_to_hex_str, vec_u8_to_fixed_32};
use gravity_utils::types::*;
use gravity_utils::{error::GravityError, message_signatures::encode_logic_call_confirm_hashed};
use std::{result::Result, time::Duration, collections::HashMap};

/// this function generates an appropriate Ethereum transaction
/// to submit the provided logic call
#[allow(clippy::too_many_arguments)]
pub async fn send_eth_logic_call(
    current_valset: Valset,
    call: LogicCall,
    confirms: &[LogicCallConfirmResponse],
    timeout: Duration,
    gravity_contract_address: EthAddress,
    gravity_id: String,
    gas_cost: GasCost,
    eth_client: EthClient,
    logic_call_skips: &mut LogicCallSkips,
) -> Result<(), GravityError> {
    let new_call_nonce = call.invalidation_nonce;
    info!(
        "Ordering signatures and submitting LogicCall {}:{} to Ethereum",
        bytes_to_hex_str(&call.invalidation_id),
        new_call_nonce
    );
    trace!("Call {:?}", call);

    let before_nonce = get_logic_call_nonce(
        gravity_contract_address,
        call.invalidation_id.clone(),
        eth_client.clone(),
    )
    .await?;

    let current_block_height = eth_client.get_block_number().await?;
    logic_call_skips.clear_old_calls(current_block_height.as_u64());

    if before_nonce >= new_call_nonce {
        info!(
            "Someone else updated the LogicCall to {}, exiting early",
            before_nonce
        );

        logic_call_skips.skip(&call);
        return Ok(());
    } else if current_block_height > call.timeout.into() {
        info!(
            "This LogicCall is timed out. timeout block: {} current block: {}, exiting early",
            current_block_height, call.timeout
        );

        logic_call_skips.skip(&call);
        return Ok(());
    }

    let contract_call = build_send_logic_call_contract_call(
        current_valset,
        &call,
        confirms,
        gravity_contract_address,
        gravity_id,
        eth_client.clone(),
    )?;

    let contract_call = contract_call
        .gas(gas_cost.gas)
        .gas_price(gas_cost.gas_price);

    let pending_tx = contract_call.send().await?;
    let tx_hash = *pending_tx;
    info!("Sent logic call with txid {}", tx_hash);
    // TODO(bolten): ethers interval default is 7s, this mirrors what web30 was doing, should we adjust?
    // additionally we are mirroring only waiting for 1 confirmation by leaving that as default
    let pending_tx = pending_tx.interval(Duration::from_secs(1));

    match tokio::time::timeout(timeout, pending_tx).await?? {
        Some(_) => (),
        None => error!(
            "Did not receive transaction receipt when submitting batch: {}",
            tx_hash
        ),
    }

    let last_nonce = get_logic_call_nonce(
        gravity_contract_address,
        call.invalidation_id,
        eth_client.clone(),
    )
    .await?;

    if last_nonce != new_call_nonce {
        error!(
            "Current nonce is {} expected to update to nonce {}",
            last_nonce, new_call_nonce
        );
    } else {
        info!(
            "Successfully updated LogicCall with new Nonce {:?}",
            last_nonce
        );
    }
    Ok(())
}

/// Returns the cost in Eth of sending this batch
pub async fn estimate_logic_call_cost(
    current_valset: Valset,
    call: LogicCall,
    confirms: &[LogicCallConfirmResponse],
    gravity_contract_address: EthAddress,
    gravity_id: String,
    eth_client: EthClient,
) -> Result<GasCost, GravityError> {
    let contract_call = build_send_logic_call_contract_call(
        current_valset,
        &call,
        confirms,
        gravity_contract_address,
        gravity_id,
        eth_client.clone(),
    )?;

    Ok(GasCost {
        gas: contract_call.estimate_gas().await?,
        gas_price: get_gas_price(eth_client.clone()).await?,
    })
}

pub fn build_send_logic_call_contract_call(
    current_valset: Valset,
    call: &LogicCall,
    confirms: &[LogicCallConfirmResponse],
    gravity_contract_address: EthAddress,
    gravity_id: String,
    eth_client: EthClient,
) -> Result<ContractCall<EthSignerMiddleware, ()>, GravityError> {
    let (current_addresses, current_powers) = current_valset.filter_empty_addresses();
    let current_powers: Vec<U256> = current_powers.iter().map(|power| (*power).into()).collect();
    let current_valset_nonce = current_valset.nonce;
    let hash = encode_logic_call_confirm_hashed(gravity_id, call.clone());
    let sig_data = current_valset.order_sigs(&hash, confirms)?;

    let transfer_amounts = call
        .transfers
        .iter()
        .map(|transfer| transfer.amount)
        .collect();
    let transfer_token_contracts = call
        .transfers
        .iter()
        .map(|transfer| transfer.token_contract_address)
        .collect();
    let fee_amounts = call.fees.iter().map(|fee| fee.amount).collect();
    let fee_token_contracts = call
        .fees
        .iter()
        .map(|fee| fee.token_contract_address)
        .collect();
    let invalidation_id = vec_u8_to_fixed_32(call.invalidation_id.clone())?;

    let contract_call = Gravity::new(gravity_contract_address, eth_client.clone())
        .submit_logic_call(
            ValsetArgs {
                validators: current_addresses,
                powers: current_powers,
                valset_nonce: current_valset_nonce.into(),
                reward_amount: U256::zero(),
                reward_token: H160::zero(),
            },
            sig_data
                .iter()
                .map(|sig_data| sig_data.to_val_sig())
                .collect(),
            LogicCallArgs {
                transfer_amounts,
                transfer_token_contracts,
                fee_amounts,
                fee_token_contracts,
                logic_contract_address: call.logic_contract_address,
                payload: call.payload.clone().into(),
                time_out: call.timeout.into(),
                invalidation_id,
                invalidation_nonce: call.invalidation_nonce.into(),
            },
        )
        .from(eth_client.address())
        .value(U256::zero());

    Ok(contract_call)
}

pub struct LogicCallSkips {
    skip_map: HashMap<Vec<u8>, HashMap<u64, LogicCall>>,
}

impl LogicCallSkips {
    pub fn new() -> Self {
        LogicCallSkips {
            skip_map: HashMap::new(),
        }
    }

    pub fn should_skip(&self, call: &LogicCall) -> bool {
        let id_skips = self.skip_map.get(&call.invalidation_id);
        if id_skips.is_some() {
            let nonce_skips = id_skips.unwrap().get(&call.invalidation_nonce);
            if nonce_skips.is_some() {
                return true;
            }
        }

        false
    }

    pub fn skip(&mut self, call: &LogicCall) {
        let id_skips = self.skip_map.get_mut(&call.invalidation_id);
        if id_skips.is_none() {
            let new_id_skips = HashMap::from([(call.invalidation_nonce, call.clone())]);
            self.skip_map.insert(call.invalidation_id.clone(), new_id_skips);
        } else {
            id_skips.unwrap().insert(call.invalidation_nonce.clone(), call.clone());
        }
    }

    pub fn clear_old_calls(&mut self, ethereum_height: u64) {
        for id_skip_map in self.skip_map.iter_mut() {
            let nonce_map = id_skip_map.1;
            for nonce_skip_map in nonce_map.clone() {
                let call = nonce_skip_map.1;
                if call.timeout < ethereum_height {
                    nonce_map.remove(&call.invalidation_nonce);
                }
            }
        }
    }
}

#[test]
fn test_logic_call_skips() {
    let logic_call_1_nonce_1 = LogicCall {
        transfers: Vec::new(),
        fees: Vec::new(),
        logic_contract_address: EthAddress::default(),
        payload: Vec::new(),
        timeout: 800,
        invalidation_id: vec![0, 1, 2],
        invalidation_nonce: 1,
    };

    let logic_call_1_nonce_2 = LogicCall {
        transfers: Vec::new(),
        fees: Vec::new(),
        logic_contract_address: EthAddress::default(),
        payload: Vec::new(),
        timeout: 900,
        invalidation_id: vec![0, 1, 2],
        invalidation_nonce: 2,
    };

    let logic_call_2 = LogicCall {
        transfers: Vec::new(),
        fees: Vec::new(),
        logic_contract_address: EthAddress::default(),
        payload: Vec::new(),
        timeout: 1000,
        invalidation_id: vec![3, 4, 5],
        invalidation_nonce: 1,
    };

    let mut skips = LogicCallSkips::new();

    assert_eq!(skips.should_skip(&logic_call_1_nonce_1), false);
    assert_eq!(skips.should_skip(&logic_call_1_nonce_2), false);
    assert_eq!(skips.should_skip(&logic_call_2), false);

    skips.skip(&logic_call_1_nonce_1);
    skips.skip(&logic_call_2);

    assert_eq!(skips.should_skip(&logic_call_1_nonce_1), true);
    assert_eq!(skips.should_skip(&logic_call_1_nonce_2), false);
    assert_eq!(skips.should_skip(&logic_call_2), true);

    skips.skip(&logic_call_1_nonce_2);

    assert_eq!(skips.should_skip(&logic_call_1_nonce_1), true);
    assert_eq!(skips.should_skip(&logic_call_1_nonce_2), true);
    assert_eq!(skips.should_skip(&logic_call_2), true);

    skips.clear_old_calls(850);

    assert_eq!(skips.should_skip(&logic_call_1_nonce_1), false);
    assert_eq!(skips.should_skip(&logic_call_1_nonce_2), true);
    assert_eq!(skips.should_skip(&logic_call_2), true);

    skips.clear_old_calls(980);

    assert_eq!(skips.should_skip(&logic_call_1_nonce_1), false);
    assert_eq!(skips.should_skip(&logic_call_1_nonce_2), false);
    assert_eq!(skips.should_skip(&logic_call_2), true);

    skips.clear_old_calls(1001);

    assert_eq!(skips.should_skip(&logic_call_1_nonce_1), false);
    assert_eq!(skips.should_skip(&logic_call_1_nonce_2), false);
    assert_eq!(skips.should_skip(&logic_call_2), false);
}
