use std::convert::TryInto;
use ethers::prelude::*;
use ethers::utils::keccak256;
use lazy_static::lazy_static;

fn err_to_selector(error: &str) -> [u8; 4] {
    keccak256(error)[0..4].try_into().unwrap()
}

// ethers abigen doesn't create functions for errors, so we're defining the selectors
// here by calculating the first hash bytes rather than computing them by hand in advance
lazy_static! {
    static ref INVALID_SIGNATURE: [u8; 4] = err_to_selector("InvalidSignature()");
    static ref INVALID_VALSET_NONCE: [u8; 4] = err_to_selector("InvalidValsetNonce(uint256,uint256)");
    static ref INVALID_BATCH_NONCE: [u8; 4] = err_to_selector("InvalidBatchNonce(uint256,uint256)");
    static ref INVALID_LOGIC_CALL_NONCE: [u8; 4] = err_to_selector("InvalidLogicCallNonce(uint256,uint256)");
    static ref INVALID_LOGIC_CALL_TRANSFERS: [u8; 4] = err_to_selector("InvalidLogicCallTransfers()");
    static ref INVALID_LOGIC_CALL_FEES: [u8; 4] = err_to_selector("InvalidLogicCallFees()");
    static ref INVALID_SEND_TO_COSMOS: [u8; 4] = err_to_selector("InvalidSendToCosmos()");
    static ref INCORRECT_CHECKPOINT: [u8; 4] = err_to_selector("IncorrectCheckpoint()");
    static ref MALFORMED_NEW_VALIDATOR_SET: [u8; 4] = err_to_selector("MalformedNewValidatorSet()");
    static ref MALFORMED_CURRENT_VALIDATOR_SET: [u8; 4] = err_to_selector("MalformedCurrentValidatorSet()");
    static ref MALFORMED_BATCH: [u8; 4] = err_to_selector("MalformedBatch()");
    static ref INSUFFICIENT_POWER: [u8; 4] = err_to_selector("InsufficientPower(uint256,uint256)");
    static ref BATCH_TIMED_OUT: [u8; 4] = err_to_selector("BatchTimedOut()");
    static ref LOGIC_CALL_TIMED_OUT: [u8; 4] = err_to_selector("LogicCallTimedOut()");
}

pub enum GravityContractError {
    InvalidSignature(InvalidSignature),
    InvalidValsetNonce(InvalidValsetNonce),
    InvalidBatchNonce(InvalidBatchNonce),
    InvalidLogicCallNonce(InvalidLogicCallNonce),
    InvalidLogicCallTransfers(InvalidLogicCallTransfers),
    InvalidLogicCallFees(InvalidLogicCallFees),
    InvalidSendToCosmos(InvalidSendToCosmos),
    IncorrectCheckpoint(IncorrectCheckpoint),
    MalformedNewValidatorSet(MalformedNewValidatorSet),
    MalformedCurrentValidatorSet(MalformedCurrentValidatorSet),
    MalformedBatch(MalformedBatch),
    InsufficientPower(InsufficientPower),
    BatchTimedOut(BatchTimedOut),
    LogicCallTimedOut(LogicCallTimedOut),
}

pub fn decode_gravity_error(data: Vec<u8>) -> Option<GravityContractError> {
    if data.len() < 4 {
        return None;
    }

    // TODO(bolten): trying to do this as a match arm was so incredibly unwieldy that
    // I gave up
    let selector: [u8; 4] = data[0..4].try_into().unwrap();

    if selector == INVALID_SIGNATURE.as_slice() {
        return Some(GravityContractError::InvalidSignature(InvalidSignature{}))
    }

    if selector == INVALID_VALSET_NONCE.as_slice() {
        if data.len() != 68 {
            return None
        }

        return Some(GravityContractError::InvalidValsetNonce(InvalidValsetNonce{
            new_nonce: data[4..36].into(),
            current_nonce: data[36..].into(),
        }))
    }

    if selector == INVALID_BATCH_NONCE.as_slice() {
        if data.len() != 68 {
            return None
        }

        return Some(GravityContractError::InvalidBatchNonce(InvalidBatchNonce{
            new_nonce: data[4..36].into(),
            current_nonce: data[36..].into(),
        }))
    }

    if selector == INVALID_LOGIC_CALL_NONCE.as_slice() {
        if data.len() != 68 {
            return None
        }

        return Some(GravityContractError::InvalidLogicCallNonce(InvalidLogicCallNonce{
            new_nonce: data[4..36].into(),
            current_nonce: data[36..].into(),
        }))
    }

    if selector == INVALID_LOGIC_CALL_TRANSFERS.as_slice() {
        return Some(GravityContractError::InvalidLogicCallTransfers(InvalidLogicCallTransfers{}))
    }

    if selector == INVALID_LOGIC_CALL_FEES.as_slice() {
        return Some(GravityContractError::InvalidLogicCallFees(InvalidLogicCallFees{}))
    }

    if selector == INVALID_SEND_TO_COSMOS.as_slice() {
        return Some(GravityContractError::InvalidSendToCosmos(InvalidSendToCosmos{}))
    }

    if selector == INCORRECT_CHECKPOINT.as_slice() {
        return Some(GravityContractError::InvalidLogicCallTransfers(InvalidLogicCallTransfers{}))
    }

    if selector == MALFORMED_NEW_VALIDATOR_SET.as_slice() {
        return Some(GravityContractError::MalformedNewValidatorSet(MalformedNewValidatorSet{}))
    }

    if selector == MALFORMED_CURRENT_VALIDATOR_SET.as_slice() {
        return Some(GravityContractError::MalformedCurrentValidatorSet(MalformedCurrentValidatorSet{}))
    }

    if selector == MALFORMED_BATCH.as_slice() {
        return Some(GravityContractError::MalformedBatch(MalformedBatch{}))
    }

    if selector == INSUFFICIENT_POWER.as_slice() {
        if data.len() != 68 {
            return None
        }

        return Some(GravityContractError::InsufficientPower(InsufficientPower{
            cumulative_power: data[4..36].into(),
            power_threshold: data[36..].into(),
        }))
    }

    if selector == BATCH_TIMED_OUT.as_slice() {
        return Some(GravityContractError::BatchTimedOut(BatchTimedOut{}))
    }

    if selector == LOGIC_CALL_TIMED_OUT.as_slice() {
        return Some(GravityContractError::LogicCallTimedOut(LogicCallTimedOut{}))
    }

    info!("Did not find gravity error");

    None
}

pub struct InvalidSignature {}

impl InvalidSignature {
    pub fn message(&self) -> String {
        "Invalid signature".to_string()
    }
}

pub struct InvalidValsetNonce {
    new_nonce: U256,
    current_nonce: U256,
}

impl InvalidValsetNonce {
    pub fn message(&self) -> String {
        format!("Invalid valset nonce, new nonce {}, current nonce {}",
            self.new_nonce, self.current_nonce)
    }
}

pub struct InvalidBatchNonce {
    new_nonce: U256,
    current_nonce: U256,
}

impl InvalidBatchNonce {
    pub fn message(&self) -> String {
        format!("Invalid batch nonce, new nonce {}, current nonce {}",
            self.new_nonce, self.current_nonce)
    }
}

pub struct InvalidLogicCallNonce {
    new_nonce: U256,
    current_nonce: U256,
}

impl InvalidLogicCallNonce {
    pub fn message(&self) -> String {
        format!("Invalid logic call nonce, new nonce {}, current nonce {}",
            self.new_nonce, self.current_nonce)
    }
}

pub struct InvalidLogicCallTransfers {}

impl InvalidLogicCallTransfers {
    pub fn message(&self) -> String {
        "Invalid logic call transfers".to_string()
    }
}

pub struct InvalidLogicCallFees {}

impl InvalidLogicCallFees {
    pub fn message(&self) -> String {
        "Invalid logic call fees".to_string()
    }
}

pub struct InvalidSendToCosmos {}

impl InvalidSendToCosmos {
    pub fn message(&self) -> String {
        "Invalid send to cosmos".to_string()
    }
}

pub struct IncorrectCheckpoint {}

impl IncorrectCheckpoint {
    pub fn message(&self) -> String {
        "Incorrect checkpoint".to_string()
    }
}

pub struct MalformedNewValidatorSet {}

impl MalformedNewValidatorSet {
    pub fn message(&self) -> String {
        "Malformed new validator set".to_string()
    }
}

pub struct MalformedCurrentValidatorSet {}

impl MalformedCurrentValidatorSet {
    pub fn message(&self) -> String {
        "Malformed current validator set".to_string()
    }
}

pub struct MalformedBatch {}

impl MalformedBatch {
    pub fn message(&self) -> String {
        "Malformed batch".to_string()
    }
}

pub struct InsufficientPower {
    cumulative_power: U256,
    power_threshold: U256,
}

impl InsufficientPower {
    pub fn message(&self) -> String {
        format!("Insufficient power, cumulative power {}, power threshold {}",
            self.cumulative_power, self.power_threshold)
    }
}

pub struct BatchTimedOut {}

impl BatchTimedOut {
    pub fn message(&self) -> String {
        "Batch timed out".to_string()
    }
}

pub struct LogicCallTimedOut {}

impl LogicCallTimedOut {
    pub fn message(&self) -> String {
        "Logic call timed out".to_string()
    }
}