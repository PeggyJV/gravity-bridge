//! This file parses the Gravity contract ethereum events. Note that there is no Ethereum ABI unpacking implementation. Instead each event
//! is parsed directly from it's binary representation. This is technical debt within this implementation. It's quite easy to parse any
//! individual event manually but a generic decoder can be quite challenging to implement. A proper implementation would probably closely
//! mirror Serde and perhaps even become a serde crate for Ethereum ABI decoding
//! For now reference the ABI encoding document here https://docs.soliditylang.org/en/v0.8.3/abi-spec.html

use super::ValsetMember;
use crate::error::GravityError;
use crate::ethereum::downcast_to_u64;
use deep_space::Address as CosmosAddress;
use ethers::abi::RawLog;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use gravity_abi::gravity::*;
use std::result::Result;

// given a Log retrieved by querying the Ethereum chain, decode it into one of
// the generated event types for the Gravity contract
fn log_to_ethers_event<T: EthLogDecode>(log: &Log) -> Result<T, ethers::abi::Error> {
    T::decode_log(&RawLog {
        topics: log.topics.clone(),
        data: log.data.to_vec(),
    })
    .map_err(From::from)
}

// our event model structs use U256 to represent block height, but Logs provide it
// to us as a U64 (strangely, no direct conversion from U64, so we have to do this type dance)
fn block_height_from_log(log: &Log) -> Result<U256, GravityError> {
    match log.block_number {
        Some(block_height) => Ok(block_height.as_u64().into()),
        None => Err(GravityError::InvalidEventLogError(format!(
            "Log does not have block number, we only search logs already in blocks? {:?}",
            log
        ))),
    }
}

// some traits to avoid code duplication

pub trait FromLog: Sized {
    fn from_log(input: &Log) -> Result<Self, GravityError>;
}

pub trait FromLogWithPrefix: Sized {
    fn from_log(input: &Log, prefix: &str) -> Result<Self, GravityError>;
}

pub trait EventNonce {
    fn get_event_nonce(&self) -> U256;
}

pub trait FromLogs {
    fn from_logs<T: FromLog>(input: &[Log]) -> Result<Vec<T>, GravityError> {
        let mut res = Vec::new();
        for item in input {
            res.push(T::from_log(item)?);
        }
        Ok(res)
    }
}

pub trait FromLogsWithPrefix {
    fn from_logs<T: FromLogWithPrefix>(
        input: &[Log],
        prefix: &str,
    ) -> Result<Vec<T>, GravityError> {
        let mut res = Vec::new();
        for item in input {
            res.push(T::from_log(item, prefix)?);
        }
        Ok(res)
    }
}

pub trait EventNonceFilter: Sized {
    /// returns all values in the array with event nonces greater
    /// than the provided value
    fn filter_by_event_nonce<T: Clone + EventNonce>(event_nonce: u64, input: &[T]) -> Vec<T> {
        input
            .iter()
            .filter(|item| item.get_event_nonce() > event_nonce.into())
            .map(|item| (*item).clone())
            .collect()
    }
}

/// A parsed struct representing the Ethereum event fired by the Gravity contract
/// when the validator set is updated. Reward amount and reward token are included
/// as part of the contract-defined type, but currently they will always be zeroed
/// out.
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ValsetUpdatedEvent {
    pub valset_nonce: U256,
    pub event_nonce: U256,
    pub reward_amount: U256,
    pub reward_token: EthAddress,
    pub block_height: U256,
    pub members: Vec<ValsetMember>,
}

impl FromLog for ValsetUpdatedEvent {
    fn from_log(input: &Log) -> Result<ValsetUpdatedEvent, GravityError> {
        let event: ValsetUpdatedEventFilter = log_to_ethers_event(input)?;
        if event.powers.len() != event.validators.len() {
            return Err(GravityError::InvalidEventLogError(format!(
                "ValsetUpdatedEvent powers and validators have different length: {:?}",
                event
            )));
        }

        let mut powers: Vec<u64> = Vec::new();
        for power in &event.powers {
            if let Some(downcast_power) = downcast_to_u64(*power) {
                powers.push(downcast_power);
            } else {
                return Err(GravityError::InvalidEventLogError(format!(
                    "ValsetUpdatedEvent contains powers that cannot be downcast to u64: {:?}",
                    event
                )));
            }
        }

        let validators: Vec<ValsetMember> = powers
            .iter()
            .zip(event.validators.iter())
            .map(|(power, validator)| ValsetMember {
                power: *power,
                eth_address: Some(*validator),
            })
            .collect();

        let mut check = validators.clone();
        check.sort();
        check.reverse();
        // if the validator set is not sorted we're in a bad spot
        // TODO(bolten): we often get warnings with the sorting of our bootstrapped valset, must figure out why
        if validators != check {
            warn!(
                "Someone submitted an unsorted validator set, this means all updates will fail until someone feeds in this unsorted value by hand {:?} instead of {:?}",
                validators, check
            );
        }

        Ok(ValsetUpdatedEvent {
            valset_nonce: event.new_valset_nonce,
            event_nonce: event.event_nonce,
            reward_amount: event.reward_amount,
            reward_token: event.reward_token,
            block_height: block_height_from_log(input)?,
            members: validators,
        })
    }
}

impl FromLogs for ValsetUpdatedEvent {}
impl EventNonce for ValsetUpdatedEvent {
    fn get_event_nonce(&self) -> U256 {
        self.event_nonce
    }
}
impl EventNonceFilter for ValsetUpdatedEvent {}

/// A parsed struct representing the Ethereum event fired by the Gravity contract when
/// a transaction batch is executed.
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct TransactionBatchExecutedEvent {
    /// the nonce attached to the transaction batch that follows
    /// it throughout it's lifecycle
    pub batch_nonce: U256,
    /// The block height this event occurred at
    pub block_height: U256,
    /// The ERC20 token contract address for the batch executed, since batches are uniform
    /// in token type there is only one
    pub erc20: EthAddress,
    /// the event nonce representing a unique ordering of events coming out
    /// of the Gravity solidity contract. Ensuring that these events can only be played
    /// back in order
    pub event_nonce: U256,
}

impl FromLog for TransactionBatchExecutedEvent {
    fn from_log(input: &Log) -> Result<TransactionBatchExecutedEvent, GravityError> {
        let event: TransactionBatchExecutedEventFilter = log_to_ethers_event(input)?;

        Ok(TransactionBatchExecutedEvent {
            batch_nonce: event.batch_nonce,
            block_height: block_height_from_log(input)?,
            erc20: event.token,
            event_nonce: event.event_nonce,
        })
    }
}

impl FromLogs for TransactionBatchExecutedEvent {}
impl EventNonce for TransactionBatchExecutedEvent {
    fn get_event_nonce(&self) -> U256 {
        self.event_nonce
    }
}
impl EventNonceFilter for TransactionBatchExecutedEvent {}

/// A parsed struct representing the Ethereum event fired when someone makes a deposit
/// on the Gravity contract
#[derive(Serialize, Deserialize, Debug, Clone, Eq, PartialEq, Hash)]
pub struct SendToCosmosEvent {
    /// The token contract address for the deposit
    pub erc20: EthAddress,
    /// The Ethereum Sender
    pub sender: EthAddress,
    /// The Cosmos destination
    pub destination: CosmosAddress,
    /// The amount of the erc20 token that is being sent
    pub amount: U256,
    /// The transaction's nonce, used to make sure there can be no accidental duplication
    pub event_nonce: U256,
    /// The block height this event occurred at
    pub block_height: U256,
}

impl FromLogWithPrefix for SendToCosmosEvent {
    fn from_log(input: &Log, prefix: &str) -> Result<SendToCosmosEvent, GravityError> {
        let event: SendToCosmosEventFilter = log_to_ethers_event(input)?;

        Ok(SendToCosmosEvent {
            erc20: event.token_contract,
            sender: event.sender,
            destination: CosmosAddress::from_slice(&event.destination[12..32], prefix)?,
            amount: event.amount,
            event_nonce: event.event_nonce,
            block_height: block_height_from_log(input)?,
        })
    }
}

impl FromLogsWithPrefix for SendToCosmosEvent {}
impl EventNonce for SendToCosmosEvent {
    fn get_event_nonce(&self) -> U256 {
        self.event_nonce
    }
}
impl EventNonceFilter for SendToCosmosEvent {}

/// A parsed struct representing the Ethereum event fired when someone uses the Gravity
/// contract to deploy a new ERC20 contract representing a Cosmos asset
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct Erc20DeployedEvent {
    /// The denom on the Cosmos chain this contract is intended to represent
    pub cosmos_denom: String,
    /// The ERC20 address of the deployed contract, this may or may not be adopted
    /// by the Cosmos chain as the contract for this asset
    pub erc20_address: EthAddress,
    /// The name of the token in the ERC20 contract, should match the Cosmos denom
    /// but it is up to the Cosmos module to check that
    pub name: String,
    /// The symbol for the token in the ERC20 contract
    pub symbol: String,
    /// The number of decimals required to represent the smallest unit of this token
    pub decimals: u8,
    pub event_nonce: U256,
    pub block_height: U256,
}

impl FromLog for Erc20DeployedEvent {
    fn from_log(input: &Log) -> Result<Erc20DeployedEvent, GravityError> {
        let event: Erc20DeployedEventFilter = log_to_ethers_event(input)?;

        Ok(Erc20DeployedEvent {
            cosmos_denom: event.cosmos_denom,
            erc20_address: event.token_contract,
            name: event.name,
            symbol: event.symbol,
            decimals: event.decimals,
            event_nonce: event.event_nonce,
            block_height: block_height_from_log(input)?,
        })
    }
}

impl FromLogs for Erc20DeployedEvent {}
impl EventNonce for Erc20DeployedEvent {
    fn get_event_nonce(&self) -> U256 {
        self.event_nonce
    }
}
impl EventNonceFilter for Erc20DeployedEvent {}

/// A parsed struct representing the Ethereum event fired when someone uses the Gravity
/// contract to send an arbitrary logic call.
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct LogicCallExecutedEvent {
    pub invalidation_id: Vec<u8>,
    pub invalidation_nonce: U256,
    pub return_data: Vec<u8>,
    pub event_nonce: U256,
    pub block_height: U256,
}

impl FromLog for LogicCallExecutedEvent {
    fn from_log(input: &Log) -> Result<LogicCallExecutedEvent, GravityError> {
        let event: LogicCallEventFilter = log_to_ethers_event(input)?;

        Ok(LogicCallExecutedEvent {
            invalidation_id: event.invalidation_id.into(),
            invalidation_nonce: event.invalidation_nonce,
            return_data: event.return_data.to_vec(),
            event_nonce: event.event_nonce,
            block_height: block_height_from_log(input)?,
        })
    }
}

impl FromLogs for LogicCallExecutedEvent {}
impl EventNonce for LogicCallExecutedEvent {
    fn get_event_nonce(&self) -> U256 {
        self.event_nonce
    }
}
impl EventNonceFilter for LogicCallExecutedEvent {}
