use crate::{
    batch_relaying::relay_batches, find_latest_valset::find_latest_valset,
    logic_call_relaying::relay_logic_calls, valset_relaying::relay_valsets,
};
use ethereum_gravity::{logic_call::LogicCallSkips, types::EthClient, utils::get_gravity_id};
use ethers::types::Address as EthAddress;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use std::time::Duration;
use tonic::transport::Channel;

pub const LOOP_SPEED: Duration = Duration::from_secs(17);
pub const PENDING_TX_TIMEOUT: Duration = Duration::from_secs(120);

#[allow(unused_variables)]
pub async fn relayer_main_loop(
    gravity_id: String,
    eth_client: EthClient,
    grpc_client: GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    eth_gas_price_multiplier: f32,
    eth_gas_multiplier: f32,
) {
    loop {
        info!("starting relayer");
        run_relayer(
            gravity_id.clone(),
            eth_client.clone(),
            grpc_client.clone(),
            gravity_contract_address,
            eth_gas_price_multiplier,
            eth_gas_multiplier,
        )
        .await;

        warn!("relayer exited unexpectedly. restarting!");
    }
}

#[allow(unused_variables)]
pub async fn run_relayer(
    gravity_id: String,
    eth_client: EthClient,
    grpc_client: GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    eth_gas_price_multiplier: f32,
    eth_gas_multiplier: f32,
) {
    let mut grpc_client = grpc_client;
    let mut logic_call_skips = LogicCallSkips::new();

    loop {
        let (async_resp, _) = tokio::join!(
            async {
                let current_eth_valset = find_latest_valset(
                    &mut grpc_client,
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
                    &mut grpc_client,
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
                    &mut grpc_client,
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
                    &mut grpc_client,
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
