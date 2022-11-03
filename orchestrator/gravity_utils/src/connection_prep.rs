//! This module provides useful tools for handling the Contact and Web30 connections for the relayer, orchestrator and various other utilities.
//! It's a common problem to have conflicts between ipv4 and ipv6 localhost and this module is first and foremost supposed to resolve that problem
//! by trying more than one thing to handle potentially misconfigured inputs.

use crate::error::GravityError;
use crate::ethereum::format_eth_address;
use cosmos_sdk_proto::cosmos::base::tendermint::v1beta1::GetLatestBlockResponse;
use ethers::prelude::*;
use ethers::providers::Provider;
use ethers::types::Address as EthAddress;
use ocular::cosmrs::AccountId;
use ocular::GrpcClient;
use ocular_somm_gravity::SommGravityExt;
use std::convert::TryFrom;
use std::process::exit;
use std::sync::Arc;
use std::time::{Duration, Instant};
use tokio::time::sleep as delay_for;
use url::Url;

pub type CosmosBlock = cosmos_sdk_proto::tendermint::types::Block;

/// Creates an eth provider
pub async fn create_eth_provider(eth_rpc_url: String) -> Result<Provider<Http>, GravityError> {
    let url = Url::parse(&eth_rpc_url)
        .unwrap_or_else(|_| panic!("Invalid Ethereum RPC url {}", eth_rpc_url));
    check_scheme(&url, &eth_rpc_url);
    let eth_url = eth_rpc_url.trim_end_matches('/');
    // TODO(bolten): should probably set a non-default interval, but what is the appropriate
    // value?
    let base_eth_provider = Provider::<Http>::try_from(eth_url)
        .unwrap_or_else(|_| panic!("Could not instantiate Ethereum HTTP provider: {}", eth_url));
    let try_base = base_eth_provider.get_block_number().await;
    match try_base {
        // it worked, lets go!
        Ok(_) => Ok(base_eth_provider),
        // did not work, now we check if it's localhost
        Err(e) => {
            warn!(
                "Failed to access Ethereum RPC with {:?} trying fallback options",
                e
            );
            if eth_url.to_lowercase().contains("localhost") {
                let port = url.port().unwrap_or(80);
                // this should be http or https
                let prefix = url.scheme();
                let ipv6_url = format!("{}://::1:{}", prefix, port);
                let ipv4_url = format!("{}://127.0.0.1:{}", prefix, port);
                let ipv6_eth_provider = Provider::<Http>::try_from(ipv6_url.as_str())
                    .unwrap_or_else(|_| {
                        panic!(
                            "Could not instantiate Ethereum HTTP provider: {}",
                            &ipv6_url
                        )
                    });
                let ipv4_eth_provider = Provider::<Http>::try_from(ipv4_url.as_str())
                    .unwrap_or_else(|_| {
                        panic!(
                            "Could not instantiate Ethereum HTTP provider: {}",
                            &ipv4_url
                        )
                    });
                let ipv6_test = ipv6_eth_provider.get_block_number().await;
                let ipv4_test = ipv4_eth_provider.get_block_number().await;
                warn!("Trying fallback urls {} {}", ipv6_url, ipv4_url);
                match (ipv4_test, ipv6_test) {
                        (Ok(_), Err(_)) => {
                            info!("Url fallback succeeded, your Ethereum rpc url {} has been corrected to {}", eth_rpc_url, ipv4_url);
                            Ok(ipv4_eth_provider)
                        }
                        (Err(_), Ok(_)) => {
                            info!("Url fallback succeeded, your Ethereum  rpc url {} has been corrected to {}", eth_rpc_url, ipv6_url);
                            Ok(ipv6_eth_provider)
                        },
                        (Ok(_), Ok(_)) => panic!("This should never happen? Why didn't things work the first time?"),
                        (Err(_), Err(_)) => panic!("Could not connect to Ethereum rpc, are you sure it's running and on the specified port? {}", eth_rpc_url)
                    }
            } else if url.port().is_none() || url.scheme() == "http" {
                let body = url.host_str().unwrap_or_else(|| {
                    panic!("Ethereum rpc url contains no host? {}", eth_rpc_url)
                });
                // transparently upgrade to https if available, we can't transparently downgrade for obvious security reasons
                let https_on_80_url = format!("https://{}:80", body);
                let https_on_443_url = format!("https://{}:443", body);
                let https_on_80_eth_provider = Provider::<Http>::try_from(https_on_80_url.as_str())
                    .unwrap_or_else(|_| {
                        panic!(
                            "Could not instantiate Ethereum HTTP provider: {}",
                            &https_on_80_url
                        )
                    });
                let https_on_443_eth_provider =
                    Provider::<Http>::try_from(https_on_443_url.as_str()).unwrap_or_else(|_| {
                        panic!(
                            "Could not instantiate Ethereum HTTP provider: {}",
                            &https_on_443_url
                        )
                    });
                let https_on_80_test = https_on_80_eth_provider.get_block_number().await;
                let https_on_443_test = https_on_443_eth_provider.get_block_number().await;
                warn!(
                    "Trying fallback urls {} {}",
                    https_on_443_url, https_on_80_url
                );
                match (https_on_80_test, https_on_443_test) {
                        (Ok(_), Err(_)) => {
                            info!("Https upgrade succeeded, your Ethereum rpc url {} has been corrected to {}", eth_rpc_url, https_on_80_url);
                            Ok(https_on_80_eth_provider)
                        },
                        (Err(_), Ok(_)) => {
                            info!("Https upgrade succeeded, your Ethereum rpc url {} has been corrected to {}", eth_rpc_url, https_on_443_url);
                            Ok(https_on_443_eth_provider)
                        },
                        (Ok(_), Ok(_)) => panic!("This should never happen? Why didn't things work the first time?"),
                        (Err(_), Err(_)) => panic!("Could not connect to Ethereum rpc, are you sure it's running and on the specified port? {}", eth_rpc_url)
                    }
            } else {
                panic!("Could not connect to Ethereum rpc! please check your grpc url {} for errors {:?}", eth_rpc_url, e)
            }
        }
    }
}

/// Verify that a url has an http or https prefix
fn check_scheme(input: &Url, original_string: &str) {
    if !(input.scheme() == "http" || input.scheme() == "https") {
        panic!(
            "Your url {} has an invalid scheme, please chose http or https",
            original_string
        )
    }
}

#[derive(Clone, Debug)]
pub enum ChainStatus {
    Moving(Box<CosmosBlock>),
    Syncing,
    WaitingToStart(Option<String>),
}

/// Gets the chain status
pub async fn get_chain_status(cosmos_client: &mut GrpcClient) -> ChainStatus {
    match cosmos_client.query_syncing().await {
        Ok(syncing) => {
            if syncing {
                return ChainStatus::Syncing;
            }
        }
        Err(e) => return ChainStatus::WaitingToStart(Some(e.to_string())),
    };

    match cosmos_client.query_latest_block().await {
        Ok(res) => {
            if res.block.is_some() {
                ChainStatus::Moving(Box::new(res.block.unwrap()))
            } else {
                ChainStatus::WaitingToStart(None)
            }
        }
        Err(e) => ChainStatus::WaitingToStart(Some(e.to_string())),
    }
}

/// This function will wait until the Cosmos node is ready, this is intended
/// for situations such as when a node is syncing or when a node is waiting on
/// a halted chain.
pub async fn wait_for_cosmos_node_ready(cosmos_client: &mut GrpcClient) {
    loop {
        match get_chain_status(cosmos_client).await {
            ChainStatus::Moving(_) => {
                println!("Cosmos node is ready!");
                break;
            }
            ChainStatus::Syncing => {
                println!("Cosmos node is syncing, waiting for it to be ready");
                delay_for(Duration::from_secs(5)).await;
                continue;
            }
            ChainStatus::WaitingToStart(error) => {
                if error.is_some() {
                    println!(
                        "Error querying latest block from Cosmos node, waiting for it to be ready: {:?}",
                        error
                    );
                    delay_for(Duration::from_secs(5)).await;
                } else {
                    println!("Cosmos node is not ready, waiting for it to be ready");
                    delay_for(Duration::from_secs(5)).await;
                }
            }
        }
    }
}

/// Waits for the next block
pub async fn wait_for_next_block(
    cosmos_client: &mut GrpcClient,
    timeout: Duration,
) -> Result<(), GravityError> {
    let res = cosmos_client.query_latest_block().await.map_err(|e| {
        GravityError::CosmosGrpcError(format!(
            "Error querying latest block from Cosmos node: {:?}",
            e
        ))
    })?;
    let current_block = extract_height(res)?;
    let start = Instant::now();
    loop {
        let res = cosmos_client.query_latest_block().await.map_err(|e| {
            GravityError::CosmosGrpcError(format!(
                "Error querying latest block from Cosmos node: {:?}",
                e
            ))
        })?;

        let next_block = extract_height(res)?;
        if next_block > current_block {
            break;
        }
        let now = Instant::now();
        if now.checked_duration_since(start).unwrap().as_secs() >= timeout.as_secs() {
            return Err(GravityError::CosmosGrpcError(
                "timed out waiting for transaction to be included in a block".to_string(),
            ));
        }
        delay_for(Duration::from_secs(1)).await;
    }

    Ok(())
}

fn extract_height(res: GetLatestBlockResponse) -> Result<i64, GravityError> {
    Ok(res
        .block
        .ok_or_else(|| GravityError::CosmosGrpcError("Cosmos node returned no block".to_string()))?
        .header
        .unwrap()
        .height)
}

/// This function checks the orchestrator delegate addresses
/// for consistency what this means is that it takes the Ethereum
/// address and Orchestrator address from the Orchestrator and checks
/// that both are registered and internally consistent.
pub async fn check_delegate_addresses(
    cosmos_client: &mut GrpcClient,
    delegate_eth_address: EthAddress,
    delegate_orchestrator_address: &AccountId,
) {
    let eth_response = cosmos_client
        .query_delegate_keys_by_ethereum_signer(&format_eth_address(delegate_eth_address))
        .await;
    let orchestrator_response = cosmos_client
        .query_delegate_keys_by_orchestrator(delegate_orchestrator_address.as_ref())
        .await;
    trace!("{:?} {:?}", eth_response, orchestrator_response);
    match (eth_response, orchestrator_response) {
        (Ok(e), Ok(o)) => {
            let req_delegate_orchestrator_address = e.orchestrator_address;
            let req_delegate_eth_address: EthAddress = o.ethereum_signer.parse().unwrap();
            if req_delegate_eth_address != delegate_eth_address
                && req_delegate_orchestrator_address != delegate_orchestrator_address.to_string()
            {
                error!("Your Delegate Ethereum and Orchestrator addresses are both incorrect!");
                error!(
                    "You provided {}  Correct Value {}",
                    delegate_eth_address, req_delegate_eth_address
                );
                error!(
                    "You provided {}  Correct Value {}",
                    delegate_orchestrator_address, req_delegate_orchestrator_address
                );
                error!("In order to resolve this issue you should double check your input value or re-register your delegate keys");
                exit(1);
            } else if req_delegate_eth_address != delegate_eth_address {
                error!("Your Delegate Ethereum address is incorrect!");
                error!(
                    "You provided {}  Correct Value {}",
                    delegate_eth_address, req_delegate_eth_address
                );
                error!("In order to resolve this issue you should double check how you input your eth private key");
                exit(1);
            } else if req_delegate_orchestrator_address != delegate_orchestrator_address.to_string()
            {
                error!("Your Delegate Orchestrator address is incorrect!");
                error!(
                    "You provided {}  Correct Value {}",
                    delegate_eth_address, req_delegate_eth_address
                );
                error!("In order to resolve this issue you should double check how you input your Orchestrator address phrase, make sure you didn't use your Validator phrase!");
                exit(1);
            }

            if e.validator_address != o.validator_address {
                error!("You are using delegate keys from two different validator addresses!");
                error!("If you get this error message I would just blow everything away and start again");
                exit(1);
            }
        }
        (Err(e), Ok(_)) => {
            error!("Your delegate Ethereum address is incorrect, please double check you private key. If you can't locate the correct private key register your delegate keys again and use the new value {:?}", e);
            exit(1);
        }
        (Ok(_), Err(e)) => {
            error!("Your delegate Cosmos address is incorrect, please double check your phrase. If you can't locate the correct phrase register your delegate keys again and use the new value {:?}", e);
            exit(1);
        }
        (Err(_), Err(_)) => {
            error!("Delegate keys are not set! Please Register your delegate keys");
            exit(1);
        }
    }
}

/// Checks if a given denom, used for fees is in the provided address
pub async fn check_for_fee_denom(
    fee_denom: &str,
    address: &AccountId,
    cosmos_client: &mut GrpcClient,
) -> Result<(), GravityError> {
    let mut found = false;
    let balances = cosmos_client
        .query_all_balances(address.as_ref())
        .await
        .map_err(|e| GravityError::CosmosGrpcError(e.to_string()))?
        .balances;
    for balance in balances {
        if balance.denom.to_string() == *fee_denom {
            found = true;
            break;
        }
    }
    if !found {
        warn!("You have specified that fees should be paid in {} but account {} has no balance of that token!", fee_denom, address);
    }

    Ok(())
}

// TODO(bolten): is using LocalWallet too specific?
/// Checks the user has some Ethereum in their address to pay for things
pub async fn check_for_eth(
    address: EthAddress,
    eth_client: Arc<SignerMiddleware<Provider<Http>, LocalWallet>>,
) {
    let balance = eth_client.get_balance(address, None).await.unwrap();
    if balance == 0u8.into() {
        warn!("You don't have any Ethereum! You will need to send some to {} for this program to work. Dust will do for basic operations, more info about average relaying costs will be presented as the program runs", address);
    }
}
