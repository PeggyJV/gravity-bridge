use super::*;
use crate::error::GravityError;
use ethers::types::{Address as EthAddress, Signature as EthSignature};
use std::{convert::TryFrom, result::Result};

/// the response we get when querying for a valset confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct LogicCall {
    pub transfers: Vec<Erc20Token>,
    pub fees: Vec<Erc20Token>,
    pub logic_contract_address: EthAddress,
    pub payload: Vec<u8>,
    pub timeout: u64,
    pub invalidation_id: Vec<u8>,
    pub invalidation_nonce: u64,
}

impl LogicCall {
    pub fn from_proto(input: gravity_proto::gravity::ContractCallTx) -> Result<Self, GravityError> {
        let mut transfers: Vec<Erc20Token> = Vec::new();
        let mut fees: Vec<Erc20Token> = Vec::new();
        for token in input.tokens {
            transfers.push(Erc20Token::from_proto(token)?)
        }
        for fee in input.fees {
            fees.push(Erc20Token::from_proto(fee)?)
        }

        Ok(LogicCall {
            transfers,
            fees,
            logic_contract_address: input.address.parse()?,
            payload: input.payload,
            timeout: input.timeout,
            invalidation_id: input.invalidation_scope,
            invalidation_nonce: input.invalidation_nonce,
        })
    }
}

/// the response we get when querying for a logic call confirmation
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct LogicCallConfirmResponse {
    pub invalidation_id: Vec<u8>,
    pub invalidation_nonce: u64,
    pub ethereum_signer: EthAddress,
    pub eth_signature: EthSignature,
}

impl LogicCallConfirmResponse {
    pub fn from_proto(
        input: gravity_proto::gravity::ContractCallTxConfirmation,
    ) -> Result<Self, GravityError> {
        Ok(LogicCallConfirmResponse {
            invalidation_id: input.invalidation_scope,
            invalidation_nonce: input.invalidation_nonce,
            ethereum_signer: input.ethereum_signer.parse()?,
            eth_signature: EthSignature::try_from(input.signature.as_slice())?,
        })
    }
}

impl Confirm for LogicCallConfirmResponse {
    fn get_eth_address(&self) -> EthAddress {
        self.ethereum_signer
    }
    fn get_signature(&self) -> EthSignature {
        self.eth_signature
    }
}
