use ethers::{
    prelude::*,
    types::transaction::{eip2718::TypedTransaction, eip712::Eip712},
};
use ethers_gcp_kms_signer::{CKMSError, GcpKeyRingRef, GcpKmsProvider, GcpKmsSigner};
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
    /// Construct a signer backed by a Google Cloud KMS asymmetric signing key.
    ///
    /// The private key never leaves KMS: every signature is produced by an
    /// AsymmetricSign API call, so no key material is present in argv, in
    /// process memory, or on disk. This is the difference that matters versus
    /// `SignerType::Local`, whose key is passed in as raw bytes and is
    /// therefore readable from `/proc/<pid>/cmdline` by any local user when
    /// supplied on a command line.
    ///
    /// The key must be `EC_SIGN_SECP256K1_SHA256`; other algorithms will not
    /// produce recoverable Ethereum signatures. Ambient GCP credentials
    /// (workload identity, attached service account, or
    /// GOOGLE_APPLICATION_CREDENTIALS) must hold
    /// `cloudkms.cryptoKeyVersions.useToSign` and
    /// `cloudkms.cryptoKeyVersions.viewPublicKey` on the key version.
    ///
    /// `chain_id` is bound at construction because GcpKmsSigner captures it;
    /// callers must therefore resolve the chain ID before building the signer.
    pub async fn new_gcp_kms(
        project_id: &str,
        location: &str,
        key_ring: &str,
        key_name: String,
        key_version: u64,
        chain_id: u64,
    ) -> Result<Self, CKMSError> {
        let key_ring_ref = GcpKeyRingRef::new(project_id, location, key_ring);
        let provider = GcpKmsProvider::new(key_ring_ref).await?;
        let signer = GcpKmsSigner::new(provider, key_name, key_version, chain_id).await?;

        Ok(SignerType::GcpKms(signer))
    }

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

                sig.v = 0;

                let sig0_address = sig
                    .recover(message.as_ref())
                    .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string()))?;

                if sig0_address.cmp(&expected_address) == Ordering::Equal {
                    return Ok(sig);
                }

                sig.v = 1;

                let sig1_address = sig
                    .recover(message.as_ref())
                    .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string()))?;

                if sig1_address.cmp(&expected_address) == Ordering::Equal {
                    return Ok(sig);
                }

                Err(ethers::providers::ProviderError::CustomError(
                    "Invalid signature while normalizing".to_string(),
                ))
            }
            // Don't need to correct v for LocalWallet
            _ => Ok(sig.clone()),
        }
    }
}

/// Maximum time to wait for a single Google Cloud KMS signing call.
///
/// KMS signing is a network RPC. The relayer's main loop drives valset, batch
/// and logic-call relaying under a single `tokio::join!`, so a signing future
/// that never resolves stops *all* relaying indefinitely rather than failing
/// the one submission. The underlying client sets no deadline of its own, so
/// we impose one here and surface a normal signer error, which the existing
/// submission paths already log and retry on a later loop iteration.
const KMS_SIGN_TIMEOUT: std::time::Duration = std::time::Duration::from_secs(30);

/// Apply [`KMS_SIGN_TIMEOUT`] to a KMS signing future, mapping both the
/// signer's own error and a timeout into a `ProviderError`.
async fn with_kms_timeout<T, E: std::fmt::Display>(
    what: &str,
    fut: impl std::future::Future<Output = Result<T, E>>,
) -> Result<T, ethers::providers::ProviderError> {
    match tokio::time::timeout(KMS_SIGN_TIMEOUT, fut).await {
        Ok(Ok(v)) => Ok(v),
        Ok(Err(e)) => Err(ethers::providers::ProviderError::CustomError(e.to_string())),
        Err(_) => Err(ethers::providers::ProviderError::CustomError(format!(
            "GCP KMS {} timed out after {}s",
            what,
            KMS_SIGN_TIMEOUT.as_secs()
        ))),
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
            SignerType::GcpKms(signer) => {
                with_kms_timeout("sign_message", signer.sign_message(message)).await
            }
        }?;

        self.normalize(msg, &sig)
    }

    async fn sign_transaction(&self, tx: &TypedTransaction) -> Result<Signature, Self::Error> {
        let sig = match self {
            SignerType::Local(wallet) => wallet
                .sign_transaction(tx)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
            SignerType::GcpKms(signer) => {
                with_kms_timeout("sign_transaction", signer.sign_transaction(tx)).await
            }
        }?;

        // Get the transaction hash for recovery
        let tx_hash = tx.sighash();
        self.normalize(&tx_hash, &sig)
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
        let sig = match self {
            SignerType::Local(wallet) => wallet
                .sign_typed_data(payload)
                .await
                .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string())),
            SignerType::GcpKms(signer) => {
                with_kms_timeout("sign_typed_data", signer.sign_typed_data(payload)).await
            }
        }?;

        // Get the typed data hash for recovery
        let hash = payload
            .encode_eip712()
            .map_err(|e| ethers::providers::ProviderError::CustomError(e.to_string()))?;
        self.normalize(&hash, &sig)
    }
}
