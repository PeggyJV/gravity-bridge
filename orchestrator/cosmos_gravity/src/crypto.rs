use deep_space::private_key::TxParts;
use std::str::FromStr;

#[cfg(not(feature = "ethermint"))]
use deep_space::public_key::COSMOS_PUBKEY_URL;
use deep_space::{
    error::{Bip39Error, HdWalletError, PrivateKeyError},
    private_key::{PrivateKey as InnerPrivateKey, SignType},
    Address, MessageArgs, Msg,
};

#[cfg(feature = "ethermint")]
pub const DEFAULT_HD_PATH: &str = "m/44'/60'/0'/0/0";
#[cfg(not(feature = "ethermint"))]
pub const DEFAULT_HD_PATH: &str = "m/44'/118'/0'/0/0";

/// PrivateKey wraps cosmos private key, switch between cosmos and ethermint behavior according to cargo features.
#[derive(Debug, Copy, Clone)]
pub struct PrivateKey(InnerPrivateKey);

impl FromStr for PrivateKey {
    type Err = PrivateKeyError;
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        InnerPrivateKey::from_str(s).map(Self)
    }
}

impl Into<InnerPrivateKey> for PrivateKey {
    fn into(self) -> InnerPrivateKey {
        self.0
    }
}

impl PrivateKey {
    pub fn from_hd_wallet_path(
        hd_path: &str,
        phrase: &str,
        passphrase: &str,
    ) -> Result<Self, PrivateKeyError> {
        InnerPrivateKey::from_hd_wallet_path(hd_path, phrase, passphrase).map(Self)
    }

    pub fn from_phrase(phrase: &str, passphrase: &str) -> Result<Self, PrivateKeyError> {
        if phrase.is_empty() {
            return Err(HdWalletError::Bip39Error(Bip39Error::BadWordCount(0)).into());
        }
        Self::from_hd_wallet_path(DEFAULT_HD_PATH, phrase, passphrase)
    }

    pub fn from_secret(secret: &[u8]) -> Self {
        Self(InnerPrivateKey::from_secret(secret))
    }

    pub fn to_address(&self, prefix: &str) -> Result<Address, PrivateKeyError> {
        #[cfg(feature = "ethermint")]
        let result = {
            let pubkey = self.0.to_public_key("")?;
            Ok(pubkey.to_ethermint_address_with_prefix(prefix)?)
        };
        #[cfg(not(feature = "ethermint"))]
        let result = self.0.to_address(prefix);

        result
    }

    pub fn sign_std_msg(
        &self,
        messages: &[Msg],
        args: MessageArgs,
        memo: impl Into<String>,
    ) -> Result<Vec<u8>, PrivateKeyError> {
        #[cfg(feature = "ethermint")]
        let result = self.0.sign_std_msg_ethermint(
            messages,
            args,
            memo,
            "/ethermint.crypto.v1.ethsecp256k1.PubKey",
        );
        #[cfg(not(feature = "ethermint"))]
        let result = self.0.sign_std_msg(messages, args, memo);

        result
    }

    pub fn build_tx(
        &self,
        messages: &[Msg],
        args: MessageArgs,
        memo: impl Into<String>,
    ) -> Result<TxParts, PrivateKeyError> {
        #[cfg(feature = "ethermint")]
        return self.0.build_tx(
            messages,
            args,
            memo,
            "/ethermint.crypto.v1.ethsecp256k1.PubKey",
            SignType::Ethermint,
        );
        #[cfg(not(feature = "ethermint"))]
        return self
            .0
            .build_tx(messages, args, memo, COSMOS_PUBKEY_URL, SignType::Cosmos);
    }
}
