use crate::one_eth;
use deep_space::address::Address as CosmosAddress;
use deep_space::coin::Coin;
use cosmos_gravity::crypto::PrivateKey as CosmosPrivateKey;
use deep_space::Contact;
use ethereum_gravity::{erc20_utils::get_erc20_balance, types::EthClient};
use ethers::core::k256::ecdsa::SigningKey;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use futures::future::join_all;
use gravity_abi::erc20::ERC20;
use rand::Rng;

pub async fn send_one_eth(dest: EthAddress, eth_client: EthClient) {
    let tx = TransactionRequest {
        from: Some(eth_client.address()),
        to: Some(NameOrAddress::Address(dest)),
        gas: None,
        gas_price: None,
        value: Some(one_eth()),
        data: Some(Vec::new().into()),
        nonce: None,
    };

    let pending_tx = eth_client.send_transaction(tx, None).await.unwrap();
    pending_tx.await.unwrap();
}

pub async fn check_cosmos_balance(
    denom: &str,
    address: CosmosAddress,
    contact: &Contact,
) -> Option<Coin> {
    let account_info = contact.get_balances(address).await.unwrap();
    trace!("Cosmos balance {:?}", account_info);
    for coin in account_info {
        // make sure the name and amount is correct
        if coin.denom.starts_with(denom) {
            return Some(coin);
        }
    }
    None
}

/// This function efficiently distributes ERC20 tokens to a large number of provided Ethereum addresses
/// the real problem here is that you can't do more than one send operation at a time from a
/// single address without your sequence getting out of whack. By manually setting the nonce
/// here we can send thousands of transactions in only a few blocks
pub async fn send_erc20_bulk(
    amount: U256,
    erc20: EthAddress,
    destinations: &[EthAddress],
    eth_client: EthClient,
) {
    let miner_balance = get_erc20_balance(erc20, eth_client.address(), eth_client.clone())
        .await
        .unwrap();
    assert!(miner_balance > amount.checked_mul(destinations.len().into()).unwrap());

    let mut nonce = eth_client
        .get_transaction_count(eth_client.address(), None)
        .await
        .unwrap();
    let mut transactions = Vec::new();

    for address in destinations {
        let data = ERC20::new(erc20, eth_client.clone())
            .transfer(*address, amount)
            .calldata()
            .unwrap();

        let tx = TransactionRequest {
            from: Some(eth_client.address()),
            to: Some(NameOrAddress::Address(erc20)),
            gas: Some(100_000u32.into()),
            gas_price: None,
            value: Some(0u32.into()),
            data: Some(data),
            nonce: Some(nonce),
        };

        let tx = eth_client.send_transaction(tx, None);
        transactions.push(tx);
        nonce += 1u64.into();
    }

    let pending_tx_results = join_all(transactions).await;
    let mut pending_txs = Vec::new();
    for pending_tx_result in pending_tx_results {
        let pending_tx = pending_tx_result.unwrap();
        pending_txs.push(pending_tx);
    }
    join_all(pending_txs).await;

    for address in destinations {
        let new_balance = get_erc20_balance(erc20, *address, eth_client.clone())
            .await
            .unwrap();
        assert!(new_balance >= amount);
    }
}

/// This function efficiently distributes ETH to a large number of provided Ethereum addresses
/// the real problem here is that you can't do more than one send operation at a time from a
/// single address without your sequence getting out of whack. By manually setting the nonce
/// here we can quickly send thousands of transactions in only a few blocks
pub async fn send_eth_bulk(amount: U256, destinations: &[EthAddress], eth_client: EthClient) {
    let mut nonce = eth_client
        .get_transaction_count(eth_client.address(), None)
        .await
        .unwrap();
    let mut transactions = Vec::new();

    for address in destinations {
        let tx = TransactionRequest {
            from: Some(eth_client.address()),
            to: Some(NameOrAddress::Address(*address)),
            gas: Some(24_000u64.into()),
            gas_price: Some(1_000_000_000u64.into()),
            value: Some(amount),
            data: Some(Vec::new().into()),
            nonce: Some(nonce),
        };

        let tx = eth_client.send_transaction(tx, None);
        transactions.push(tx);
        nonce += 1u64.into();
    }

    let pending_tx_results = join_all(transactions).await;
    let mut pending_txs = Vec::new();
    for pending_tx_result in pending_tx_results {
        let pending_tx = pending_tx_result.unwrap();
        pending_txs.push(pending_tx);
    }
    join_all(pending_txs).await;
}

pub fn get_user_key() -> BridgeUserKey {
    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    // the starting location of the funds
    let eth_key = SigningKey::from_bytes(&secret).unwrap();
    let eth_address = LocalWallet::from(eth_key.clone()).address();
    // the destination on cosmos that sends along to the final ethereum destination
    let cosmos_key = CosmosPrivateKey::from_secret(&secret);
    let cosmos_address = cosmos_key
        .to_address(CosmosAddress::DEFAULT_PREFIX)
        .unwrap();
    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    // the final destination of the tokens back on Ethereum
    let eth_dest_key = SigningKey::from_bytes(&secret).unwrap();
    let eth_dest_address = LocalWallet::from(eth_dest_key.clone()).address();
    BridgeUserKey {
        eth_address,
        eth_key,
        cosmos_address,
        cosmos_key,
        eth_dest_key,
        eth_dest_address,
    }
}
#[derive(Debug)]
pub struct BridgeUserKey {
    // the starting addresses that get Eth balances to send across the bridge
    pub eth_address: EthAddress,
    pub eth_key: SigningKey,
    // the cosmos addresses that get the funds and send them on to the dest eth addresses
    pub cosmos_address: CosmosAddress,
    pub cosmos_key: CosmosPrivateKey,
    // the location tokens are sent back to on Ethereum
    pub eth_dest_address: EthAddress,
    pub eth_dest_key: SigningKey,
}

#[derive(Debug, Clone)]
pub struct ValidatorKeys {
    /// The Ethereum key used by this validator to sign Gravity bridge messages
    pub eth_key: SigningKey,
    /// The Orchestrator key used by this validator to submit oracle messages and signatures
    /// to the cosmos chain
    pub orch_key: CosmosPrivateKey,
    /// The validator key used by this validator to actually sign and produce blocks
    pub validator_key: CosmosPrivateKey,
}
