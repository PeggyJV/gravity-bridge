//! for things that don't belong in the cosmos or ethereum libraries but also don't belong
//! in a function specific library
use clarity::Error as ClarityError;
use deep_space::error::AddressError as CosmosAddressError;
use deep_space::error::CosmosGrpcError;
use num_bigint::ParseBigIntError;
use thiserror::Error;
use tokio::time::error::Elapsed;
use tonic::Status;
use web30::jsonrpc::error::Web3Error;

#[derive(Error, Debug)]
#[allow(clippy::large_enum_variant)]
pub enum GravityError {
    #[error("Got invalid BigInt from cosmos! {0}")]
    InvalidBigInt(ParseBigIntError),

    #[error("Cosmos gRPC error {0}")]
    CosmosGrpcError(CosmosGrpcError),

    #[error("Cosmos Address error {0}")]
    CosmosAddressError(CosmosAddressError),

    #[error("Ethereum REST error {0}")]
    EthereumRestError(Web3Error),

    #[error("Invalid bridge state! {0}")]
    InvalidBridgeStateError(String),

    #[error("ValidatorSet Update Failed")]
    FailedToUpdateValset,

    #[error("Contract operation failed {0}")]
    EthereumContractError(String),

    #[error("Invalid TX options for this call {0}")]
    InvalidOptionsError(String),

    #[error("Clarity Error {0}")]
    ClarityError(ClarityError),

    #[error("Operation timed out!")]
    TimeoutError,

    #[error("InvalidEvent {0}")]
    InvalidEventLogError(String),

    #[error("Gravity gRPC error {0}")]
    GravityGrpcError(Status),

    #[error("Insufficient Voting Power to Pass")]
    InsufficientVotingPowerToPass(String),

    #[error("Failed to parse big integer {0}")]
    ParseBigIntError(ParseBigIntError),
}

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

impl From<Web3Error> for GravityError {
    fn from(error: Web3Error) -> Self {
        GravityError::EthereumRestError(error)
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
