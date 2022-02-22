//! This is the happy path test for Cosmos to Ethereum asset transfers, meaning assets originated on Cosmos
use crate::utils::get_user_key;
use crate::utils::send_one_eth;
use crate::MINER_CLIENT;
use crate::TOTAL_TIMEOUT;
use crate::{get_fee, utils::ValidatorKeys};
use clarity::Uint256;
use cosmos_gravity::send::{send_request_batch_tx, send_to_eth};
use deep_space::coin::Coin;
use deep_space::Contact;
use ethereum_gravity::erc20_utils::get_erc20_balance;
use ethereum_gravity::{deploy_erc20::deploy_erc20, utils::get_event_nonce};
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_proto::gravity::{
    query_client::QueryClient as GravityQueryClient, DenomToErc20Request,
};
use gravity_utils::ethereum::downcast_to_u64;
use std::str::FromStr;
use std::sync::Arc;
use tonic::transport::Channel;

pub async fn happy_path_test_v2(
    eth_provider: &Provider<Http>,
    grpc_client: GravityQueryClient<Channel>,
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
    gravity_address: EthAddress,
) {
    let mut grpc_client = grpc_client;
    let eth_wallet = LocalWallet::from(keys[0].eth_key.clone());
    let provider = eth_provider.clone();
    let chain_id = provider
        .get_chainid()
        .await
        .expect("Could not retrieve chain ID");
    let chain_id = downcast_to_u64(chain_id).expect("Chain ID overflowed when downcasting to u64");
    let eth_client = Arc::new(SignerMiddleware::new(
        provider,
        eth_wallet.with_chain_id(chain_id),
    ));
    let starting_event_nonce = get_event_nonce(gravity_address, eth_client.clone())
        .await
        .unwrap();

    let token_to_send_to_eth = "footoken".to_string();
    let token_to_send_to_eth_display_name = "mfootoken".to_string();

    deploy_erc20(
        token_to_send_to_eth.clone(),
        token_to_send_to_eth_display_name.clone(),
        token_to_send_to_eth_display_name.clone(),
        6,
        gravity_address,
        Some(TOTAL_TIMEOUT),
        1.0,
        eth_client.clone(),
    )
    .await
    .unwrap();
    let ending_event_nonce = get_event_nonce(gravity_address, eth_client.clone())
        .await
        .unwrap();

    assert!(starting_event_nonce != ending_event_nonce);
    info!(
        "Successfully deployed new ERC20 representing FooToken on Cosmos with event nonce {}",
        ending_event_nonce
    );
    // the erc20 representing the cosmos asset on Ethereum
    let erc20_contract = match tokio::time::timeout(TOTAL_TIMEOUT, async {
        loop {
            let res = grpc_client
                .denom_to_erc20(DenomToErc20Request {
                    denom: token_to_send_to_eth.clone(),
                })
                .await;

            if let Ok(res) = res {
                break res;
            }
            tokio::time::sleep(std::time::Duration::from_secs(1)).await;
        }
    })
    .await
    {
        Ok(res) => {
            let erc20 = res.into_inner().erc20;
            info!(
                "Successfully adopted {} token contract of {}",
                token_to_send_to_eth, erc20
            );
            Some(erc20)
        }
        Err(_) => {
            panic!(
                "Cosmos did not adopt the ERC20 contract for {} it must be invalid in some way",
                token_to_send_to_eth
            );
        }
    };

    let erc20_contract: EthAddress = erc20_contract.unwrap().parse().unwrap();

    // one foo token
    let amount_to_bridge: Uint256 = 1_000_000u64.into();
    let send_to_user_coin = Coin {
        denom: token_to_send_to_eth.clone(),
        amount: amount_to_bridge.clone() + 100u8.into(),
    };
    let send_to_eth_coin = Coin {
        denom: token_to_send_to_eth.clone(),
        amount: amount_to_bridge.clone(),
    };

    let user = get_user_key();
    // send the user some footoken
    contact
        .send_tokens(
            send_to_user_coin.clone(),
            Some(get_fee()),
            user.cosmos_address,
            keys[0].validator_key.into(),
            Some(TOTAL_TIMEOUT),
        )
        .await
        .unwrap();

    let balances = contact.get_balances(user.cosmos_address).await.unwrap();
    let mut found = false;
    for coin in balances {
        if coin.denom == token_to_send_to_eth.clone() {
            found = true;
            break;
        }
    }
    if !found {
        panic!(
            "Failed to send {} to the user address",
            token_to_send_to_eth
        );
    }
    info!(
        "Sent some {} to user address {}",
        token_to_send_to_eth, user.cosmos_address
    );
    // send the user some eth, they only need this to check their
    // erc20 balance, so a pretty minor usecase
    send_one_eth(user.eth_address, (*MINER_CLIENT).clone()).await;
    info!("Sent 1 eth to user address {}", user.eth_address);

    let res = send_to_eth(
        user.cosmos_key,
        user.eth_address,
        send_to_eth_coin,
        get_fee(),
        (10f64, "footoken".to_string()),
        contact,
        1.0,
    )
    .await
    .unwrap();
    info!("Send to eth res {:?}", res);
    info!(
        "Locked up {} {} to send to Cosmos",
        amount_to_bridge, token_to_send_to_eth
    );

    let res = send_request_batch_tx(
        keys[0].validator_key,
        token_to_send_to_eth.clone(),
        (10f64, "footoken".to_string()),
        contact,
        1.0,
    )
    .await
    .unwrap();
    info!("Batch request res {:?}", res);
    info!("Sent batch request to move things along");

    info!("Waiting for batch to be signed and relayed to Ethereum");
    match tokio::time::timeout(TOTAL_TIMEOUT, async {
        loop {
            let balance =
                get_erc20_balance(erc20_contract, user.eth_address, (*MINER_CLIENT).clone()).await;
            if balance.is_err() {
                continue;
            }
            let balance = balance.unwrap();
            let balance = Uint256::from_str(balance.to_string().as_str()).unwrap();
            if balance == amount_to_bridge {
                break;
            } else if balance != 0u8.into() {
                panic!(
                    "Expected {} {} but got {} instead",
                    amount_to_bridge, token_to_send_to_eth, balance
                );
            }
            tokio::time::sleep(std::time::Duration::from_secs(1)).await;
        }
    })
    .await
    {
        Ok(_) => {
            info!(
                "Successfully bridged {} Cosmos asset {} to Ethereum!",
                amount_to_bridge, token_to_send_to_eth
            );
        }
        Err(_) => {
            panic!(
                "An error occured while bridging {} Cosmos asset {} to Ethereum!",
                amount_to_bridge, token_to_send_to_eth
            );
        }
    }
}
