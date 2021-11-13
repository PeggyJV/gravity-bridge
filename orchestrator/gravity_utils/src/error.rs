//! for things that don't belong in the cosmos or ethereum libraries but also don't belong
//! in a function specific library

use clarity::Error as ClarityError;
use deep_space::error::AddressError as CosmosAddressError;
use deep_space::error::CosmosGrpcError;
use ethers::abi::Error as EthersAbiError;
use ethers::abi::ethereum_types::FromDecStrErr as EthersParseUintError;
use ethers::prelude::*;
use ethers::prelude::signer::SignerMiddlewareError;
use ethers::types::SignatureError as EthersSignatureError;
use rustc_hex::FromHexError as EthersParseAddressError;
use num_bigint::ParseBigIntError;
use std::fmt::{self, Debug};
use tokio::time::error::Elapsed;
use tonic::Status;

#[derive(Debug)]
#[allow(clippy::large_enum_variant)]
pub enum GravityError {
    InvalidBigInt(ParseBigIntError),
    CosmosGrpcError(CosmosGrpcError),
    CosmosAddressError(CosmosAddressError),
    EthereumRestError(SignerMiddlewareError<Provider<Http>, LocalWallet>),
    EthersAbiError(EthersAbiError),
    EthersParseAddressError(EthersParseAddressError),
    EthersParseUintError(EthersParseUintError),
    EthersSignatureError(EthersSignatureError),
    InvalidBridgeStateError(String),
    FailedToUpdateValset,
    EthereumContractError(String),
    InvalidOptionsError(String),
    ClarityError(ClarityError),
    TimeoutError,
    InvalidEventLogError(String),
    GravityGrpcError(Status),
    InsufficientVotingPowerToPass(String),
    ParseBigIntError(ParseBigIntError),
}

impl fmt::Display for GravityError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            GravityError::GravityGrpcError(val) => write!(f, "Gravity gRPC error {}", val),
            GravityError::CosmosGrpcError(val) => write!(f, "Cosmos gRPC error {}", val),
            GravityError::InvalidBigInt(val) => {
                write!(f, "Got invalid BigInt from cosmos! {}", val)
            }
            GravityError::CosmosAddressError(val) => write!(f, "Cosmos Address error {}", val),
            GravityError::EthereumRestError(val) => write!(f, "Ethereum REST error {}", val),
            GravityError::EthersAbiError(val) => write!(f, "Ethers ABI error {}", val),
            GravityError::EthersParseAddressError(val) => write!(f, "Ethers H160 address parse error {}", val),
            GravityError::EthersParseUintError(val) => write!(f, "Ethers U256 parse error {}", val),
            GravityError::EthersSignatureError(val) => write!(f, "Ethers signature error {}", val),
            GravityError::InvalidOptionsError(val) => {
                write!(f, "Invalid TX options for this call {}", val)
            }
            GravityError::InvalidBridgeStateError(val) => {
                write!(f, "Invalid bridge state! {}", val)
            }
            GravityError::FailedToUpdateValset => write!(f, "ValidatorSetUpdate Failed!"),
            GravityError::TimeoutError => write!(f, "Operation timed out!"),
            GravityError::ClarityError(val) => write!(f, "Clarity Error {}", val),
            GravityError::InvalidEventLogError(val) => write!(f, "InvalidEvent: {}", val),
            GravityError::EthereumContractError(val) => {
                write!(f, "Contract operation failed: {}", val)
            }
            GravityError::InsufficientVotingPowerToPass(val) => {
                write!(f, "{}", val)
            }
            GravityError::ParseBigIntError(val) => write!(f, "Failed to parse big integer {}", val),
        }
    }
}

impl std::error::Error for GravityError {}

impl From<CosmosGrpcError> for GravityError {
    fn from(error: CosmosGrpcError) -> Self {
        GravityError::CosmosGrpcError(error)
    }
}

impl From<Elapsed> for GravityError {
    fn from(_error: Elapsed) -> Self {
        GravityError::TimeoutError
    }
}

impl From<ClarityError> for GravityError {
    fn from(error: ClarityError) -> Self {
        GravityError::ClarityError(error)
    }
}

impl From<SignerMiddlewareError<Provider<Http>, LocalWallet>> for GravityError {
    fn from(error: SignerMiddlewareError<Provider<Http>, LocalWallet>) -> Self {
        GravityError::EthereumRestError(error)
    }
}

impl From<EthersAbiError> for GravityError {
    fn from(error: EthersAbiError) -> Self {
        GravityError::EthersAbiError(error)
    }
}

impl From<EthersParseAddressError> for GravityError {
    fn from(error: EthersParseAddressError) -> Self {
        GravityError::EthersParseAddressError(error)
    }
}

impl From<EthersParseUintError> for GravityError {
    fn from(error: EthersParseUintError) -> Self {
        GravityError::EthersParseUintError(error)
    }
}

impl From<EthersSignatureError> for GravityError {
    fn from(error: EthersSignatureError) -> Self {
        GravityError::EthersSignatureError(error)
    }
}

impl From<Status> for GravityError {
    fn from(error: Status) -> Self {
        GravityError::GravityGrpcError(error)
    }
}

impl From<CosmosAddressError> for GravityError {
    fn from(error: CosmosAddressError) -> Self {
        GravityError::CosmosAddressError(error)
    }
}

impl From<ParseBigIntError> for GravityError {
    fn from(error: ParseBigIntError) -> Self {
        GravityError::InvalidBigInt(error)
    }
}
