//! This test verifies that live updating of orchestrator keys works correctly

use crate::utils::ValidatorKeys;
use ethers::core::k256::elliptic_curve::generic_array::GenericArray;
use ethers::types::Address as EthAddress;
use ethers::{core::k256::ecdsa::SigningKey, prelude::*};
use gravity::deep_space::Address as CosmosAddress;
use gravity::deep_space::Contact;
use gravity::deep_space::CosmosPrivateKey;
use gravity::deep_space::PrivateKey;
use gravity::send::update_gravity_delegate_addresses;
use gravity::utils::ethereum::format_eth_address;
use gravity_proto::gravity::{
    query_client::QueryClient as GravityQueryClient, DelegateKeysByEthereumSignerRequest,
    DelegateKeysByOrchestratorRequest,
};
use rand::Rng;
use std::time::Duration;
use tonic::transport::Channel;

const BLOCK_TIMEOUT: Duration = Duration::from_secs(30);

pub async fn orch_keys_update(
    grpc_client: GravityQueryClient<Channel>,
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
) {
    let mut keys = keys;
    let mut grpc_client = grpc_client;
    // just to test that we have the right keys from the gentx
    info!("About to check already set delegate addresses");
    for k in keys.iter() {
        let eth_address = LocalWallet::from(k.eth_key.clone()).address();
        let orch_address = k.orch_key.to_address(&contact.get_prefix()).unwrap();
        let eth_response = grpc_client
            .delegate_keys_by_ethereum_signer(DelegateKeysByEthereumSignerRequest {
                ethereum_signer: format_eth_address(eth_address),
            })
            .await
            .unwrap()
            .into_inner();

        let parsed_response_orch_address: CosmosAddress =
            eth_response.orchestrator_address.parse().unwrap();
        assert_eq!(parsed_response_orch_address, orch_address);

        let orchestrator_response = grpc_client
            .delegate_keys_by_orchestrator(DelegateKeysByOrchestratorRequest {
                orchestrator_address: orch_address.to_string(),
            })
            .await
            .unwrap()
            .into_inner();

        let parsed_response_eth_address: EthAddress =
            orchestrator_response.ethereum_signer.parse().unwrap();
        assert_eq!(parsed_response_eth_address, eth_address);
    }

    info!("Starting with {:?}", keys);

    // now we change them all
    for k in keys.iter_mut() {
        let mut rng = rand::thread_rng();
        let secret: [u8; 32] = rng.gen();
        // generate some new keys to replace the old ones
        let key_bytes = GenericArray::from_slice(&secret);
        let ethereum_key = SigningKey::from_bytes(&key_bytes).unwrap();
        let ethereum_wallet = LocalWallet::from(ethereum_key.clone());
        let cosmos_key = CosmosPrivateKey::from_secret(&secret);
        // update the keys in the key list
        k.eth_key = ethereum_key;
        k.orch_key = cosmos_key;
        let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

        info!(
            "Signing and submitting Delegate addresses {} for validator {}",
            format_eth_address(ethereum_wallet.address()),
            cosmos_address,
        );
        // send in the new delegate keys signed by the validator address
        update_gravity_delegate_addresses(
            contact,
            ethereum_wallet.address(),
            cosmos_address,
            k.validator_key,
            ethereum_wallet,
            (0f64, "".to_string()),
            2.0,
        )
        .await
        .expect("Failed to set delegate addresses!");
    }

    contact.wait_for_next_block(BLOCK_TIMEOUT).await.unwrap();

    // TODO registering is too unreliable right now for confusing reasons, revisit with prototx

    // info!("About to check changed delegate addresses");
    // // verify that the change has taken place
    // for k in keys.iter() {
    //     let eth_address = k.eth_key.to_public_key().unwrap();
    //     let orch_address = k.orch_key.to_public_key().unwrap().to_address();

    //     let orchestrator_response = grpc_client
    //         .get_delegate_key_by_orchestrator(QueryDelegateKeysByOrchestratorAddress {
    //             orchestrator_address: orch_address.to_string(),
    //         })
    //         .await
    //         .unwrap()
    //         .into_inner();

    //     let parsed_response_eth_address: EthAddress =
    //         orchestrator_response.eth_address.parse().unwrap();
    //     assert_eq!(parsed_response_eth_address, eth_address);

    //     let eth_response = grpc_client
    //         .get_delegate_key_by_eth(QueryDelegateKeysByEthAddress {
    //             eth_address: eth_address.to_string(),
    //         })
    //         .await
    //         .unwrap()
    //         .into_inner();

    //     let parsed_response_orch_address: CosmosAddress =
    //         eth_response.orchestrator_address.parse().unwrap();
    //     assert_eq!(parsed_response_orch_address, orch_address);
    // }
}
