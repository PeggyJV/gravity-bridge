// use crate::get_chain_id;
use crate::get_fee;
use crate::get_gas_price;
use crate::utils::*;
use crate::MINER_ADDRESS;
use crate::MINER_CLIENT;
use crate::OPERATION_TIMEOUT;
use crate::TOTAL_TIMEOUT;
use clarity::Uint256;
use ethers::core::k256::ecdsa::SigningKey;
use ethers::core::k256::elliptic_curve::generic_array::GenericArray;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity::deep_space::address::Address as CosmosAddress;
use gravity::deep_space::coin::Coin;
use gravity::deep_space::Contact;
use gravity::deep_space::CosmosPrivateKey;
use gravity::deep_space::PrivateKey;
use gravity::ethereum::erc20_utils::get_erc20_balance;
use gravity::ethereum::utils::get_valset_nonce;
use gravity::ethereum::{send_to_cosmos::send_to_cosmos, utils::get_tx_batch_nonce};
use gravity::send::send_to_eth;
use gravity::utils::types::SendToCosmosEvent;
use gravity::{build, query::get_oldest_unsigned_transaction_batch, send};
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use rand::Rng;
use std::str::FromStr;
use std::time::Duration;
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;

pub async fn happy_path_test(
    grpc_client: GravityQueryClient<Channel>,
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
    gravity_address: EthAddress,
    erc20_address: EthAddress,
    validator_out: bool,
) {
    let mut grpc_client = grpc_client;

    // bootstrapping tests finish here and we move into operational tests

    // send 3 valset updates to make sure the process works back to back
    // don't do this in the validator out test because it changes powers
    // randomly and may actually make it impossible for that test to pass
    // by random re-allocation of powers. If we had 5 or 10 validators
    // instead of 3 this wouldn't be a problem. But with 3 not even 1% of
    // power can be reallocated to the down validator before things stop
    // working. We'll settle for testing that the initial valset (generated
    // with the first block) is successfully updated

    if !validator_out {
        for _ in 0u32..2 {
            test_valset_update(contact, &keys, gravity_address).await;
        }
    } else {
        wait_for_nonzero_valset(gravity_address).await;
    }

    // generate an address for coin sending tests, this ensures test imdepotency
    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    let dest_cosmos_private_key = CosmosPrivateKey::from_secret(&secret);
    let dest_cosmos_address = dest_cosmos_private_key.to_address("cosmos").unwrap();
    let key_bytes = GenericArray::from_slice(&secret);
    let dest_eth_private_key = SigningKey::from_bytes(&key_bytes).unwrap();
    let dest_eth_wallet = LocalWallet::from(dest_eth_private_key.clone());
    let dest_eth_address = dest_eth_wallet.address();

    // the denom and amount of the token bridged from Ethereum -> Cosmos
    // so the denom is the gravity<hash> token name
    // Send a token 3 times
    for _ in 0u32..3 {
        test_erc20_deposit(
            contact,
            dest_cosmos_address,
            gravity_address,
            erc20_address,
            100u64.into(),
        )
        .await;
    }

    // We are going to submit a duplicate tx with nonce 1
    // This had better not increase the balance again
    // this test may have false positives if the timeout is not
    // long enough. TODO check for an error on the cosmos send response
    submit_duplicate_erc20_send(
        1u64.into(),
        contact,
        erc20_address,
        1u64.into(),
        dest_cosmos_address,
        &keys,
    )
    .await;

    // we test a batch by sending a transaction
    test_batch(
        contact,
        &mut grpc_client,
        dest_eth_address,
        gravity_address,
        keys[0].validator_key,
        dest_cosmos_private_key,
        erc20_address,
    )
    .await;
}

pub async fn wait_for_nonzero_valset(gravity_address: EthAddress) {
    match tokio::time::timeout(TOTAL_TIMEOUT, async {
        let mut current_eth_valset_nonce =
            get_valset_nonce(gravity_address, (*MINER_CLIENT).clone())
                .await
                .expect("Failed to get current eth valset");

        while 0 == current_eth_valset_nonce {
            info!("Validator set is not yet updated to >0, waiting");
            current_eth_valset_nonce = get_valset_nonce(gravity_address, (*MINER_CLIENT).clone())
                .await
                .expect("Failed to get current eth valset");
            delay_for(Duration::from_secs(4)).await;
        }
    })
    .await
    {
        Ok(_) => {
            println!("Success")
        }
        Err(_) => {
            panic!("Failed to update validator set");
        }
    }
}

pub async fn test_valset_update(
    contact: &Contact,
    keys: &[ValidatorKeys],
    gravity_address: EthAddress,
) {
    // if we don't do this the orchestrators may run ahead of us and we'll be stuck here after
    // getting credit for two loops when we did one
    let starting_eth_valset_nonce = get_valset_nonce(gravity_address, (*MINER_CLIENT).clone())
        .await
        .expect("Failed to get starting eth valset");
    // this is hacky and not really a good way to test validator set updates in a highly
    // repeatable fashion. What we really need to do is be aware of the total staking state
    // and manipulate the validator set very intentionally rather than kinda blindly like
    // we are here. For example the more your run this function the less this fixed amount
    // makes any difference, eventually it will fail because the change to the total staked
    // percentage is too small.
    let mut rng = rand::thread_rng();
    let keys_to_change = rng.gen_range(0..keys.len());
    let keys_to_change = &keys[keys_to_change];

    let validator_to_change = keys_to_change.validator_key;
    let delegate_address = validator_to_change
        .to_address("cosmosvaloper")
        .unwrap()
        .to_string();

    // should be about 4% of the total power to start
    // let amount = crate::STARTING_STAKE_PER_VALIDATOR / 4; // 12.5B
    let amount = crate::STAKE_SUPPLY_PER_VALIDATOR / 4; // 25B
    let amount = gravity::deep_space::coin::Coin {
        amount: amount.into(),
        denom: "stake".to_string(),
    };
    match tokio::time::timeout(TOTAL_TIMEOUT, async {
        // now we send a valset request that the orchestrators will pick up on
        // in this case we send it as the first validator because they can pay the fee
        info!(
            "Sending in valset request (starting_eth_valset_nonce {})",
            starting_eth_valset_nonce
        );

        info!(
            "Delegating {} to {} in order to generate a validator set update",
            amount, delegate_address
        );
        loop {
            let res = contact
                .delegate_to_validator(
                    delegate_address.parse().unwrap(),
                    amount.clone(),
                    get_fee(),
                    keys_to_change.orch_key,
                    Some(OPERATION_TIMEOUT),
                )
                .await;

            if res.is_err() {
                warn!("Delegate to validator failed (will retry) {:?}", res);
                continue; // retry
            }
            break;
        }
    })
    .await
    {
        Ok(_) => {
            info!("Delegated {} to {}", amount, delegate_address);
        }
        Err(_) => {
            panic!("Delegate to validator timed out.");
        }
    }

    let mut current_eth_valset_nonce = get_valset_nonce(gravity_address, (*MINER_CLIENT).clone())
        .await
        .expect("Failed to get current eth valset");

    match tokio::time::timeout(TOTAL_TIMEOUT, async {
        while starting_eth_valset_nonce == current_eth_valset_nonce {
            info!(
                "Validator set is not yet updated to >{}, waiting",
                starting_eth_valset_nonce
            );
            current_eth_valset_nonce = get_valset_nonce(gravity_address, (*MINER_CLIENT).clone())
                .await
                .expect("Failed to get current eth valset");
            delay_for(Duration::from_secs(4)).await;
        }
    })
    .await
    {
        Ok(_) => {
            assert!(starting_eth_valset_nonce != current_eth_valset_nonce);
            info!("Validator set successfully updated!");
        }
        Err(_) => {
            panic!("Failed to update validator set");
        }
    }
}

/// this function tests Ethereum -> Cosmos
async fn test_erc20_deposit(
    contact: &Contact,
    dest: CosmosAddress,
    gravity_address: EthAddress,
    erc20_address: EthAddress,
    amount: U256,
) {
    let amount_uint256 = Uint256::from_str(amount.to_string().as_str()).unwrap();
    let start_coin = check_cosmos_balance("gravity", dest, contact).await;
    info!(
        "Sending to Cosmos from {} to {} with amount {}",
        *MINER_ADDRESS, dest, amount
    );
    // we send some erc20 tokens to the gravity contract to register a deposit
    let tx_id = send_to_cosmos(
        erc20_address,
        gravity_address,
        amount,
        dest,
        Some(TOTAL_TIMEOUT),
        (*MINER_CLIENT).clone(),
    )
    .await
    .expect("Failed to send tokens to Cosmos");
    info!("Send to Cosmos txid: {:#066x}", tx_id);

    match tokio::time::timeout(TOTAL_TIMEOUT, async {
        match (
            start_coin.clone(),
            check_cosmos_balance("gravity", dest, contact).await,
        ) {
            (Some(start_coin), Some(end_coin)) => {
                if start_coin.amount + amount_uint256.clone() == end_coin.amount
                    && start_coin.denom == end_coin.denom
                {
                    info!(
                        "Successfully bridged ERC20 {}{} to Cosmos! Balance is now {}{}",
                        amount, start_coin.denom, end_coin.amount, end_coin.denom
                    );
                }
            }
            (None, Some(end_coin)) => {
                if amount_uint256 == end_coin.amount {
                    info!(
                        "Successfully bridged ERC20 {}{} to Cosmos! Balance is now {}{}",
                        amount, end_coin.denom, end_coin.amount, end_coin.denom
                    );
                } else {
                    panic!("Failed to bridge ERC20!")
                }
            }
            _ => {}
        }
    })
    .await
    {
        Ok(_) => {
            info!("Waiting for ERC20 deposit");
            contact.wait_for_next_block(TOTAL_TIMEOUT).await.unwrap();
        }
        Err(_) => {
            panic!("Failed to bridge ERC20!");
        }
    }
}

#[allow(clippy::too_many_arguments)]
async fn test_batch(
    contact: &Contact,
    grpc_client: &mut GravityQueryClient<Channel>,
    dest_eth_address: EthAddress,
    gravity_address: EthAddress,
    requester_cosmos_private_key: CosmosPrivateKey,
    dest_cosmos_private_key: CosmosPrivateKey,
    erc20_contract: EthAddress,
) {
    let dest_cosmos_address = dest_cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap();
    let coin = check_cosmos_balance("gravity", dest_cosmos_address, contact)
        .await
        .unwrap();
    let token_name = coin.denom;
    let amount = coin.amount;

    let bridge_denom_fee = Coin {
        denom: token_name.clone(),
        amount: 1u64.into(),
    };
    let amount = amount - 5u64.into();
    info!(
        "Sending {}{} from {} on Cosmos back to Ethereum",
        amount, token_name, dest_cosmos_address
    );
    let res = send_to_eth(
        dest_cosmos_private_key,
        dest_eth_address,
        Coin {
            denom: token_name.clone(),
            amount: amount.clone(),
        },
        bridge_denom_fee.clone(),
        (10f64, "footoken".to_string()),
        contact,
        1.0,
    )
    .await
    .unwrap();
    info!("Sent tokens to Ethereum with {:?}", res);

    contact.wait_for_next_block(TOTAL_TIMEOUT).await.unwrap();
    let requester_address = requester_cosmos_private_key
        .to_address(&contact.get_prefix())
        .unwrap();
    get_oldest_unsigned_transaction_batch(grpc_client, requester_address)
        .await
        .expect("Failed to get batch to sign");

    let mut current_eth_batch_nonce =
        get_tx_batch_nonce(gravity_address, erc20_contract, (*MINER_CLIENT).clone())
            .await
            .expect("Failed to get current eth valset");
    let starting_batch_nonce = current_eth_batch_nonce;

    match tokio::time::timeout(TOTAL_TIMEOUT, async {
        while starting_batch_nonce == current_eth_batch_nonce {
            info!(
                "Batch is not yet submitted {}>, waiting",
                starting_batch_nonce
            );
            current_eth_batch_nonce =
                get_tx_batch_nonce(gravity_address, erc20_contract, (*MINER_CLIENT).clone())
                    .await
                    .expect("Failed to get current eth tx batch nonce");
            delay_for(Duration::from_secs(4)).await;
        }
    })
    .await
    {
        Ok(_) => {
            println!("Submitted transaction batch set");
        }
        Err(_) => {
            panic!("Failed to submit transaction batch set");
        }
    }

    let eth_client = (*MINER_CLIENT).clone();
    let tx = TransactionRequest {
        from: Some(eth_client.address()),
        to: Some(NameOrAddress::Address(dest_eth_address)),
        gas: None,
        gas_price: None,
        value: Some(1_000_000_000_000_000_000u128.into()),
        data: Some(Vec::new().into()),
        nonce: None,
        chain_id: None,
    };

    let pending_tx = eth_client.send_transaction(tx, None).await.unwrap();
    pending_tx.await.unwrap();

    let amount_u256 = U256::from_str(amount.to_string().as_str()).unwrap();

    // we have to send this address one eth so that it can perform contract calls
    send_one_eth(dest_eth_address, eth_client.clone()).await;
    assert_eq!(
        get_erc20_balance(erc20_contract, dest_eth_address, eth_client.clone())
            .await
            .unwrap(),
        amount_u256
    );
    info!(
        "Successfully updated txbatch nonce to {} and sent {}{} tokens to Ethereum!",
        current_eth_batch_nonce, amount, token_name
    );
}

// this function submits a EthereumBridgeDepositClaim to the module with a given nonce. This can be set to be a nonce that has
// already been submitted to test the nonce functionality.
#[allow(clippy::too_many_arguments)]
async fn submit_duplicate_erc20_send(
    nonce: U256,
    contact: &Contact,
    erc20_address: EthAddress,
    amount: U256,
    receiver: CosmosAddress,
    keys: &[ValidatorKeys],
) {
    let start_coin = check_cosmos_balance("gravity", receiver, contact)
        .await
        .expect("Did not find coins!");

    let ethereum_sender = "0x912fd21d7a69678227fe6d08c64222db41477ba0"
        .parse()
        .unwrap();
    let event = SendToCosmosEvent {
        event_nonce: nonce,
        block_height: 500u16.into(),
        erc20: erc20_address,
        sender: ethereum_sender,
        destination: receiver,
        amount,
    };

    // iterate through all validators and try to send an event with duplicate nonce
    for k in keys.iter() {
        let cosmos_key = k.validator_key;

        let messages = build::ethereum_event_messages(
            contact,
            cosmos_key,
            vec![event.clone()],
            vec![],
            vec![],
            vec![],
            vec![],
        );

        let gas_price = get_gas_price();
        let res = send::send_messages(contact, cosmos_key, gas_price, messages, 1.0).await;
        let res = res.unwrap();
        trace!("Submitted duplicate sendToCosmos event: {:?}", res);
    }

    if let Some(end_coin) = check_cosmos_balance("gravity", receiver, contact).await {
        if start_coin.amount == end_coin.amount && start_coin.denom == end_coin.denom {
            info!("Successfully failed to duplicate ERC20!");
        } else {
            panic!("Duplicated ERC20!")
        }
    } else {
        panic!("Duplicate test failed for unknown reasons!");
    }
}
