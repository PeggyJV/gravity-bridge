use crate::{
    batch_relaying::relay_batches, find_latest_valset::find_latest_valset,
    logic_call_relaying::relay_logic_calls, valset_relaying::relay_valsets,
};
use ethereum_gravity::{logic_call::LogicCallSkips, types::EthClient, utils::get_gravity_id};
use ethers::types::Address as EthAddress;
use ocular::GrpcClient;
use std::time::Duration;

pub const LOOP_SPEED: Duration = Duration::from_secs(17);
pub const PENDING_TX_TIMEOUT: Duration = Duration::from_secs(120);

/// This function contains the orchestrator primary loop, it is broken out of the main loop so that
/// it can be called in the test runner for easier orchestration of multi-node tests
#[allow(unused_variables)]
pub async fn relayer_main_loop(
    eth_client: EthClient,
    cosmos_client: GrpcClient,
    gravity_contract_address: EthAddress,
    eth_gas_price_multiplier: f32,
    eth_gas_multiplier: f32,
) {
    let mut cosmos_client = cosmos_client;
    let gravity_id = get_gravity_id(gravity_contract_address, eth_client.clone()).await;
    if gravity_id.is_err() {
        error!("Failed to get GravityID, check your Eth node");
        return;
    }
    let gravity_id = gravity_id.unwrap();
    let mut logic_call_skips = LogicCallSkips::new();

    loop {
        let (async_resp, _) = tokio::join!(
            async {
                let current_eth_valset = find_latest_valset(
                    &mut cosmos_client,
                    gravity_contract_address,
                    eth_client.clone(),
                )
                .await;
                if current_eth_valset.is_err() {
                    error!("Could not get current valset! {:?}", current_eth_valset);
                    return;
                }
                let current_eth_valset = current_eth_valset.unwrap();

                relay_valsets(
                    current_eth_valset.clone(),
                    eth_client.clone(),
                    &mut cosmos_client,
                    gravity_contract_address,
                    gravity_id.clone(),
                    PENDING_TX_TIMEOUT,
                    eth_gas_price_multiplier,
                    eth_gas_multiplier,
                )
                .await;

                relay_batches(
                    current_eth_valset.clone(),
                    eth_client.clone(),
                    &mut cosmos_client,
                    gravity_contract_address,
                    gravity_id.clone(),
                    PENDING_TX_TIMEOUT,
                    eth_gas_price_multiplier,
                    eth_gas_multiplier,
                )
                .await;

                relay_logic_calls(
                    current_eth_valset,
                    eth_client.clone(),
                    &mut cosmos_client,
                    gravity_contract_address,
                    gravity_id.clone(),
                    PENDING_TX_TIMEOUT,
                    eth_gas_price_multiplier,
                    eth_gas_multiplier,
                    &mut logic_call_skips,
                )
                .await;
            },
            tokio::time::sleep(LOOP_SPEED)
        );
    }
}
