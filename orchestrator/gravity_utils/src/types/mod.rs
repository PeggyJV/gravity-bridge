mod batches;
mod ethereum_events;
mod gravity_contract_errors;
mod logic_call;
mod signatures;
mod valsets;

use crate::error::GravityError;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use std::result::Result;

pub use batches::*;
pub use ethereum_events::*;
pub use gravity_contract_errors::*;
pub use logic_call::*;
pub use signatures::*;
pub use valsets::*;

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct Erc20Token {
    pub amount: U256,
    #[serde(rename = "contract")]
    pub token_contract_address: EthAddress,
}

impl Erc20Token {
    pub fn from_proto(input: gravity_proto::gravity::Erc20Token) -> Result<Self, GravityError> {
        Ok(Erc20Token {
            amount: U256::from_dec_str(input.amount.as_str())?,
            token_contract_address: input.contract.parse()?,
        })
    }
}
