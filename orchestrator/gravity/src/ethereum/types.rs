use ethers::{
    prelude::*,
    types::transaction::{eip2718::TypedTransaction, eip712::Eip712},
};
use ethers_gcp_kms_signer::GcpKmsSigner;
use std::{cmp::Ordering, sync::Arc};

pub type EthSignerMiddleware = SignerMiddleware<Provider<Http>, SignerType>;
pub type EthClient = Arc<EthSignerMiddleware>;

/// Wrapper enum for different signer types
#[derive(Clone, Debug)]
pub enum SignerType {
    Local(LocalWallet),
    GcpKms(GcpKmsSigner),
}

impl SignerType {
    // Normalizes signature for Gravity's signature validation
    pub fn normalize(
        &self,
        message: impl AsRef<[u8]>,
        sig: &Signature,
    ) -> Result<Signature, ethers::providers::ProviderError> {
        match self {
            // Gravity does not implement eip155 modifications to v so we need to recompute v to the allowed range of 0 or 1
            SignerType::GcpKms(signer) => {
                let expected_address = signer.address();
                let mut sig = sig.to_owned();
                let message = message.as_ref();

                // secp256k1 curve order from go-ethereum/crypto/crypto.go
                let n = U256::from_str_radix(
                    "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141",
                    16,
                )
                .unwrap();
                let n_half = n / U256::from(2);

                // Normalize s to be less than n/2 (homestead rule)
                if sig.s > n_half {
                    sig.s = n - sig.s;
                }

                sig.v = 0;

                let sig0_address = sig
                    .recover(message)
                    .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string()))?;

                if sig0_address.cmp(&expected_address) == Ordering::Equal {
                    return Ok(sig);
                }

                sig.v = 1;

                let sig1_address = sig
                    .recover(message)
                    .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string()))?;

                if sig1_address.cmp(&expected_address) == Ordering::Equal {
                    return Ok(sig);
                }

                Err(ethers::providers::ProviderError::CustomError(
                    "Invalid signature, unable to normalize".to_string(),
                ))
            }
            // Don't need to correct v for LocalWallet
            _ => Ok(sig.clone()),
        }
    }
}

#[async_trait::async_trait]
impl Signer for SignerType {
    type Error = ethers::providers::ProviderError;

    async fn sign_message<S: Send + Sync + AsRef<[u8]>>(
        &self,
        message: S,
    ) -> Result<Signature, Self::Error> {
        // We copy the msg because it's not Clone
        let mut msg = vec![0u8; message.as_ref().len()];

        msg.copy_from_slice(message.as_ref());

        let sig = match self {
            SignerType::Local(wallet) => wallet
                .sign_message(message)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
            SignerType::GcpKms(signer) => signer
                .sign_message(message)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
        }?;

        self.normalize(&msg, &sig)
    }

    async fn sign_transaction(&self, tx: &TypedTransaction) -> Result<Signature, Self::Error> {
        match self {
            SignerType::Local(wallet) => wallet
                .sign_transaction(tx)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
            SignerType::GcpKms(signer) => signer
                .sign_transaction(tx)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
        }
    }

    fn address(&self) -> Address {
        match self {
            SignerType::Local(wallet) => wallet.address(),
            SignerType::GcpKms(signer) => signer.address(),
        }
    }

    fn chain_id(&self) -> u64 {
        match self {
            SignerType::Local(wallet) => wallet.chain_id(),
            SignerType::GcpKms(signer) => signer.chain_id(),
        }
    }

    fn with_chain_id<T: Into<u64>>(self, chain_id: T) -> Self {
        match self {
            SignerType::Local(wallet) => SignerType::Local(wallet.with_chain_id(chain_id)),
            SignerType::GcpKms(signer) => SignerType::GcpKms(signer.with_chain_id(chain_id)),
        }
    }

    async fn sign_typed_data<T: Eip712 + Send + Sync>(
        &self,
        payload: &T,
    ) -> Result<Signature, Self::Error> {
        match self {
            SignerType::Local(wallet) => wallet
                .sign_typed_data(payload)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
            SignerType::GcpKms(signer) => signer
                .sign_typed_data(payload)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
        }
    }
}
