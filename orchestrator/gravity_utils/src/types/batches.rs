use super::*;
use crate::error::GravityError;
use deep_space::Address as CosmosAddress;
use ethers::core::abi::Token;
use ethers::types::{Address as EthAddress, Signature as EthSignature};
use std::{convert::TryFrom, result::Result};

/// This represents an individual transaction being bridged over to Ethereum
/// parallel is the OutgoingTransferTx in x/gravity/types/batch.go
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct BatchTransaction {
    pub id: u64,
    pub sender: CosmosAddress,
    pub ethereum_recipient: EthAddress,
    pub erc20_token: Erc20Token,
    pub erc20_fee: Erc20Token,
}

impl BatchTransaction {
    pub fn from_proto(input: gravity_proto::gravity::SendToEthereum) -> Result<Self, GravityError> {
        if input.erc20_fee.is_none() || input.erc20_token.is_none() {
            return Err(GravityError::InvalidBridgeStateError(
                "Can not have tx with null erc20_token!".to_string(),
            ));
        }

        Ok(BatchTransaction {
            id: input.id,
            sender: input.sender.parse()?,
            ethereum_recipient: input.ethereum_recipient.parse()?,
            erc20_token: Erc20Token::from_proto(input.erc20_token.unwrap())?,
            erc20_fee: Erc20Token::from_proto(input.erc20_fee.unwrap())?,
        })
    }
}

/// the response we get when querying for a valset confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct TransactionBatch {
    pub nonce: u64,
    pub batch_timeout: u64,
    pub transactions: Vec<BatchTransaction>,
    pub total_fee: Erc20Token,
    pub token_contract: EthAddress,
}

impl TransactionBatch {
    /// extracts the amounts, destinations and fees as submitted to the Ethereum contract
    /// and used for signatures
    pub fn get_checkpoint_values(&self) -> (Vec<U256>, Vec<EthAddress>, Vec<U256>) {
        let amounts: Vec<U256> = self
            .transactions
            .iter()
            .map(|tx| tx.erc20_token.amount)
            .collect();
        let destinations: Vec<EthAddress> = self
            .transactions
            .iter()
            .map(|tx| tx.ethereum_recipient)
            .collect();
        let fees: Vec<U256> = self
            .transactions
            .iter()
            .map(|tx| tx.erc20_fee.amount)
            .collect();

        (amounts, destinations, fees)
    }

    pub fn get_checkpoint_values_tokens(&self) -> (Token, Token, Token) {
        let (amounts, destinations, fees) = self.get_checkpoint_values();
        let amounts_tokens = amounts.iter().map(|amount| Token::Uint(*amount)).collect();
        let destinations_tokens = destinations
            .iter()
            .map(|destination| Token::Address(*destination))
            .collect();
        let fees_tokens = fees.iter().map(|fee| Token::Uint(*fee)).collect();

        (
            Token::Array(amounts_tokens),
            Token::Array(destinations_tokens),
            Token::Array(fees_tokens),
        )
    }

    pub fn from_proto(input: gravity_proto::gravity::BatchTx) -> Result<Self, GravityError> {
        let mut transactions = Vec::new();
        let mut running_total_fee: Option<Erc20Token> = None;
        for tx in input.transactions {
            let tx = BatchTransaction::from_proto(tx)?;
            if let Some(total_fee) = running_total_fee {
                let running_amount = total_fee.amount.checked_add(tx.erc20_fee.amount);
                if running_amount.is_none() {
                    return Err(GravityError::OverflowError(
                        format!("U256 overflow when adding all fees together for transaction batch with nonce {}", input.batch_nonce)
                    ));
                }

                running_total_fee = Some(Erc20Token {
                    token_contract_address: total_fee.token_contract_address,
                    amount: running_amount.unwrap(),
                });
            } else {
                running_total_fee = Some(tx.erc20_fee.clone())
            }
            transactions.push(tx);
        }

        if let Some(total_fee) = running_total_fee {
            Ok(TransactionBatch {
                batch_timeout: input.timeout,
                nonce: input.batch_nonce,
                transactions,
                token_contract: total_fee.token_contract_address,
                total_fee,
            })
        } else {
            Err(GravityError::InvalidBridgeStateError(
                "Transaction batch containing no transactions!".to_string(),
            ))
        }
    }
}

/// the response we get when querying for a batch confirmation
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct BatchConfirmResponse {
    pub nonce: u64,
    pub token_contract: EthAddress,
    pub ethereum_signer: EthAddress,
    pub eth_signature: EthSignature,
}

impl BatchConfirmResponse {
    pub fn from_proto(
        input: gravity_proto::gravity::BatchTxConfirmation,
    ) -> Result<Self, GravityError> {
        Ok(BatchConfirmResponse {
            nonce: input.batch_nonce,
            token_contract: input.token_contract.parse()?,
            ethereum_signer: input.ethereum_signer.parse()?,
            eth_signature: EthSignature::try_from(input.signature.as_slice())?,
        })
    }
}

impl Confirm for BatchConfirmResponse {
    fn get_eth_address(&self) -> EthAddress {
        self.ethereum_signer
    }
    fn get_signature(&self) -> EthSignature {
        self.eth_signature
    }
}
