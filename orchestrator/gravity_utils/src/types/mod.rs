mod batches;
mod ethereum_events;
mod logic_call;
mod signatures;
mod valsets;
use crate::error::GravityError;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;

pub use batches::*;
pub use ethereum_events::*;
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
        warn!("Converting input.amount: {:?}", input.amount);
        warn!("Converting input.contract: {:?}", input.contract);
        Ok(Erc20Token {
            amount: U256::from_dec_str(input.amount.as_str())?,
            token_contract_address: input.contract.parse()?,
        })
    }
}
