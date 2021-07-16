use serde::{Deserialize, Serialize};
use signatory::FsKeyStore;
use std::path::Path;

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct GorcConfig {
    pub keystore: String,
    pub gravity: GravitySection,
    pub ethereum: EthereumSection,
    pub cosmos: CosmosSection,
}

impl GorcConfig {
    pub fn load_deep_space_key(&self, name: String) -> deep_space::private_key::PrivateKey {
        let keystore = Path::new(&self.keystore);
        let keystore = FsKeyStore::create_or_open(keystore).expect("Could not open keystore");
        let name = name.parse().expect("Could not parse name");
        let key = keystore.load(&name).expect("Could not load key");

        let key = key
            .to_pem()
            .parse::<k256::elliptic_curve::SecretKey<k256::Secp256k1>>()
            .expect("Could not parse key");

        let key = deep_space::utils::bytes_to_hex_str(&key.to_bytes());

        return key
            .parse::<deep_space::private_key::PrivateKey>()
            .expect("Could not parse private key");
    }
}

impl Default for GorcConfig {
    fn default() -> Self {
        Self {
            keystore: "/tmp/keystore".to_owned(),
            gravity: GravitySection::default(),
            ethereum: EthereumSection::default(),
            cosmos: CosmosSection::default(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct GravitySection {
    pub contract: String,
}

impl Default for GravitySection {
    fn default() -> Self {
        Self {
            contract: "0x6b175474e89094c44da98b954eedeac495271d0f".to_owned(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct EthereumSection {
    pub key_derivation_path: String,
    pub rpc: String,
}

impl Default for EthereumSection {
    fn default() -> Self {
        Self {
            key_derivation_path: "m/44'/60'/0'/0/0".to_owned(),
            rpc: "http://localhost:8545".to_owned(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct CosmosSection {
    pub key_derivation_path: String,
    pub grpc: String,
    pub prefix: String,
}

impl Default for CosmosSection {
    fn default() -> Self {
        Self {
            key_derivation_path: "m/44'/118'/0'/0/0".to_owned(),
            grpc: "http://localhost:9090".to_owned(),
            prefix: "cosmos".to_owned(),
        }
    }
}
