use ethers::{
    prelude::*,
    types::transaction::{eip2718::TypedTransaction, eip712::Eip712},
};
use ethers_gcp_kms_signer::GcpKmsSigner;
use std::sync::Arc;

pub type EthSignerMiddleware = SignerMiddleware<Provider<Http>, SignerType>;
pub type EthClient = Arc<EthSignerMiddleware>;

/// Wrapper enum for different signer types
#[derive(Clone, Debug)]
pub enum SignerType {
    Local(LocalWallet),
    GcpKms(GcpKmsSigner),
}

#[async_trait::async_trait]
impl Signer for SignerType {
    type Error = ethers::providers::ProviderError;

    async fn sign_message<S: Send + Sync + AsRef<[u8]>>(
        &self,
        message: S,
    ) -> Result<Signature, Self::Error> {
        match self {
            SignerType::Local(wallet) => wallet
                .sign_message(message)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
            SignerType::GcpKms(signer) => signer
                .sign_message(message)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
        }
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
