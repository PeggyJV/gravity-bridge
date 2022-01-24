use crate::{
    one_eth, one_hundred_eth, one_hundred_eth_uint256, utils::*, MINER_CLIENT, TOTAL_TIMEOUT,
};
use clarity::Uint256;
use cosmos_gravity::send::{send_request_batch_tx, send_to_eth};
use deep_space::coin::Coin;
use deep_space::Contact;
use ethereum_gravity::{
    erc20_utils::get_erc20_balance, send_to_cosmos::send_to_cosmos, utils::get_tx_batch_nonce,
};
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use futures::future::join_all;
use gravity_utils::ethereum::downcast_to_u64;
use std::{collections::HashSet, str::FromStr, sync::Arc, time::Duration};

const TIMEOUT: Duration = Duration::from_secs(120);

/// The number of users we will be simulating for this test, each user
/// will get one token from each token type in erc20_addresses and send it
/// across the bridge to Cosmos as a deposit and then send it back to a different
/// Ethereum address in a transaction batch
/// So the total number of
/// Ethereum sends = (2 * NUM_USERS)
/// ERC20 sends = (erc20_addresses.len() * NUM_USERS)
/// Gravity Deposits = (erc20_addresses.len() * NUM_USERS)
/// Batches executed = erc20_addresses.len() * (NUM_USERS / 100)
const NUM_USERS: usize = 100;

/// Perform a stress test by sending thousands of
/// transactions and producing large batches
#[allow(clippy::too_many_arguments)]
pub async fn transaction_stress_test(
    eth_provider: &Provider<Http>,
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
    gravity_address: EthAddress,
    erc20_addresses: Vec<EthAddress>,
) {
    // Generate 100 user keys to send ETH and multiple types of tokens
    let mut user_keys = Vec::new();
    for _ in 0..NUM_USERS {
        user_keys.push(get_user_key());
    }
    // the sending eth addresses need Ethereum to send ERC20 tokens to the bridge
    let sending_eth_addresses: Vec<EthAddress> = user_keys.iter().map(|i| i.eth_address).collect();
    // the destination eth addresses need Ethereum to perform a contract call and get their erc20 balances
    let dest_eth_addresses: Vec<EthAddress> =
        user_keys.iter().map(|i| i.eth_dest_address).collect();
    let mut eth_destinations = Vec::new();
    eth_destinations.extend(sending_eth_addresses.clone());
    eth_destinations.extend(dest_eth_addresses);
    send_eth_bulk(one_eth(), &eth_destinations, (*MINER_CLIENT).clone()).await;
    info!("Sent {} addresses 1 ETH", NUM_USERS);

    // now we need to send all the sending eth addresses erc20's to send
    for token in erc20_addresses.iter() {
        send_erc20_bulk(
            one_hundred_eth(),
            *token,
            &sending_eth_addresses,
            (*MINER_CLIENT).clone(),
        )
        .await;
        info!("Sent {} addresses 100 {}", NUM_USERS, token);
    }
    for token in erc20_addresses.iter() {
        let mut sends = Vec::new();
        for keys in user_keys.iter() {
            let eth_wallet = LocalWallet::from(keys.eth_key.clone());
            let provider = eth_provider.clone();
            let chain_id = provider
                .get_chainid()
                .await
                .expect("Could not retrieve chain ID");
            let chain_id =
                downcast_to_u64(chain_id).expect("Chain ID overflowed when downcasting to u64");
            let eth_client = Arc::new(SignerMiddleware::new(
                provider,
                eth_wallet.with_chain_id(chain_id),
            ));
            let fut = send_to_cosmos(
                *token,
                gravity_address,
                one_hundred_eth(),
                keys.cosmos_address,
                Some(TIMEOUT),
                eth_client.clone(),
            );
            sends.push(fut);
        }
        let results = join_all(sends).await;
        for result in results {
            result.unwrap();
        }
        info!(
            "Locked 100 {} from {} into the Gravity Ethereum Contract",
            token, NUM_USERS
        );
    }

    let mut good = true;
    match tokio::time::timeout(TOTAL_TIMEOUT, async {
        loop {
            good = true;
            for keys in user_keys.iter() {
                let c_addr = keys.cosmos_address;
                let balances = contact.get_balances(c_addr).await.unwrap();
                for token in erc20_addresses.iter() {
                    let mut found = false;
                    for balance in balances.iter() {
                        if balance.denom.contains(&token.to_string())
                            && balance.amount == one_hundred_eth_uint256()
                        {
                            found = true;
                        }
                    }
                    if !found {
                        good = false;
                    }
                }
            }
            if good {
                break;
            }
            tokio::time::sleep(Duration::from_secs(5)).await;
        }
    })
    .await
    {
        Ok(_) => {
            info!(
                "All {} deposits bridged to Cosmos successfully!",
                user_keys.len() * erc20_addresses.len()
            );
        }
        Err(_) => {
            panic!(
                "Failed to perform all {} deposits to Cosmos!",
                user_keys.len() * erc20_addresses.len()
            );
        }
    }

    let send_amount = one_hundred_eth_uint256() - 500u16.into();

    let mut denoms = HashSet::new();
    for token in erc20_addresses.iter() {
        let mut futs = Vec::new();
        for keys in user_keys.iter() {
            let c_addr = keys.cosmos_address;
            let c_key = keys.cosmos_key;
            let e_dest_addr = keys.eth_dest_address;
            let balances = contact.get_balances(c_addr).await.unwrap();
            // this way I don't have to hardcode a denom and we can change the way denoms are formed
            // without changing this test.
            let mut send_coin = None;
            for balance in balances {
                if balance.denom.contains(&token.to_string()) {
                    send_coin = Some(balance.clone());
                    denoms.insert(balance.denom);
                }
            }
            let mut send_coin = send_coin.unwrap();
            send_coin.amount = send_amount.clone();
            let send_fee = Coin {
                denom: send_coin.denom.clone(),
                amount: 1u8.into(),
            };
            let res = send_to_eth(
                c_key,
                e_dest_addr,
                send_coin,
                send_fee,
                (0f64, "".to_string()),
                contact,
                1.0,
            );
            futs.push(res);
        }
        let results = join_all(futs).await;
        for result in results {
            let result = result.unwrap();
            trace!("SendToEth result {:?}", result);
        }
        info!(
            "Successfully placed {} {} into the tx pool",
            NUM_USERS, token
        );
    }

    for denom in denoms {
        info!("Requesting batch for {}", denom);
        let res = send_request_batch_tx(
            keys[0].validator_key,
            denom,
            (0f64, "".to_string()),
            contact,
            1.0,
        )
        .await
        .unwrap();
        info!("batch request response is {:?}", res);
    }

    match tokio::time::timeout(TOTAL_TIMEOUT, async {
        loop {
            good = true;
            for keys in user_keys.iter() {
                let e_dest_addr = keys.eth_dest_address;
                for token in erc20_addresses.iter() {
                    let bal = get_erc20_balance(*token, e_dest_addr, (*MINER_CLIENT).clone())
                        .await
                        .unwrap();
                    let bal = Uint256::from_str(bal.to_string().as_str()).unwrap();
                    if bal != send_amount.clone() {
                        good = false;
                    }
                }
            }
            if good {
                break;
            }
            tokio::time::sleep(Duration::from_secs(5)).await;
        }
    })
    .await
    {
        Ok(_) => {
            info!(
                "All {} withdraws to Ethereum bridged successfully!",
                NUM_USERS * erc20_addresses.len()
            );
        }
        Err(_) => {
            panic!(
                "Failed to perform all {} withdraws to Ethereum!",
                NUM_USERS * erc20_addresses.len()
            );
        }
    }

    // we should find a batch nonce greater than zero since all the batches
    // executed
    let eth_wallet = LocalWallet::from(keys[0].eth_key.clone());
    let eth_client = Arc::new(SignerMiddleware::new(eth_provider.clone(), eth_wallet));
    for token in erc20_addresses {
        assert!(
            get_tx_batch_nonce(gravity_address, token, eth_client.clone())
                .await
                .unwrap()
                > 0
        )
    }
}
