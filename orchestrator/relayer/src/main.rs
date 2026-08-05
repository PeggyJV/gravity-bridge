use std::sync::Arc;

use crate::main_loop::relayer_main_loop;
use crate::main_loop::LOOP_SPEED;
use docopt::Docopt;
use env_logger::Env;
use ethers::prelude::*;
use ethers::signers::LocalWallet as EthWallet;
use ethers::types::Address as EthAddress;
use gravity::ethereum::types::SignerType;
use gravity::utils::connection_prep::check_for_eth;
use gravity::utils::connection_prep::create_rpc_connections;
use gravity::utils::connection_prep::wait_for_cosmos_node_ready;
use gravity::utils::ethereum::downcast_to_u64;
use gravity::utils::ethereum::format_eth_address;

pub mod batch_relaying;
pub mod find_latest_valset;
pub mod logic_call_relaying;
pub mod main_loop;
pub mod valset_relaying;

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate lazy_static;
#[macro_use]
extern crate log;

#[derive(Debug, Deserialize)]
struct Args {
    flag_ethereum_key: Option<String>,
    flag_gcp_kms_project: Option<String>,
    flag_gcp_kms_location: Option<String>,
    flag_gcp_kms_key_ring: Option<String>,
    flag_gcp_kms_key_name: Option<String>,
    flag_gcp_kms_key_version: Option<u64>,
    flag_cosmos_grpc: String,
    flag_address_prefix: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage: {} [--ethereum-key=<key>] [--gcp-kms-project=<id>] [--gcp-kms-location=<loc>] [--gcp-kms-key-ring=<ring>] [--gcp-kms-key-name=<name>] [--gcp-kms-key-version=<ver>] --cosmos-grpc=<url> --address-prefix=<prefix> --ethereum-rpc=<url> --contract-address=<addr>
        Options:
            -h --help                    Show this screen.
            --ethereum-key=<key>         An Ethereum private key containing non-trivial funds
            --gcp-kms-project=<id>       Google Cloud project id holding the Kms key ring
            --gcp-kms-location=<loc>     Google Cloud location of the Kms key ring
            --gcp-kms-key-ring=<ring>    Kms key ring name
            --gcp-kms-key-name=<name>    Kms key name, must be a secp256k1 signing key
            --gcp-kms-key-version=<ver>  Kms key version, required, each version has its own address
            --cosmos-grpc=<gurl>         The Cosmos gRPC url
            --address-prefix=<prefix>    The prefix for addresses on this Cosmos chain
            --ethereum-rpc=<eurl>        The Ethereum RPC url, Geth light clients work and sync fast
            --contract-address=<addr>    The Ethereum contract address for Gravity
        About:
            The Gravity relayer component, responsible for relaying data from the Cosmos blockchain
            to the Ethereum blockchain, cosmos key and fees are optional since they are only used
            to request the creation of batches or validator sets to relay.
            for Althea-Gravity.
            Signing. Supply exactly one signing method. An Ethereum private key is
            discouraged, because it is visible in the process table to every local user
            on the host. Prefer the Google Cloud Kms options, where the private key never
            leaves Kms. All five Kms options are required together. Each Kms key version
            has its own Ethereum address, which pays relayer gas and must be funded
            before use.
            Written By: {}
            Version {}",
            env!("CARGO_PKG_NAME"),
            env!("CARGO_PKG_AUTHORS"),
            env!("CARGO_PKG_VERSION"),
        );
}

/// How the relayer will obtain Ethereum signatures.
///
/// Resolved from CLI flags up front so that an invalid combination aborts
/// before any network connections are made.
enum SigningConfig {
    /// Raw private key supplied on the command line. Readable via
    /// /proc/<pid>/cmdline by any local user; retained for compatibility.
    ///
    /// Held already-parsed so that a malformed key is rejected during argument
    /// handling, before any network work, as it was before KMS support moved
    /// signer construction below chain-ID resolution.
    Local(Box<EthWallet>),
    /// Google Cloud KMS. No key material is held by this process.
    GcpKms {
        project: String,
        location: String,
        key_ring: String,
        key_name: String,
        key_version: u64,
    },
}

/// Hand-written rather than derived: the `Local` variant holds a raw private
/// key, and a derived `Debug` would print it in full anywhere this value is
/// formatted, which is the exact class of leak this whole change exists to
/// close. The KMS variant holds only a key reference, so it prints in full.
impl std::fmt::Debug for SigningConfig {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            SigningConfig::Local(_) => f.write_str("Local(<redacted>)"),
            SigningConfig::GcpKms {
                project,
                location,
                key_ring,
                key_name,
                key_version,
            } => write!(
                f,
                "GcpKms(projects/{}/locations/{}/keyRings/{}/cryptoKeys/{}/cryptoKeyVersions/{})",
                project, location, key_ring, key_name, key_version
            ),
        }
    }
}

impl SigningConfig {
    /// The GCP KMS options, which are all required together. The key version is
    /// deliberately required rather than defaulted: a KMS signing address is
    /// version-specific, so silently pinning to version 1 would both fail to
    /// follow a key rotation (relaying stops once version 1 is disabled) and
    /// hide the fact that changing version changes the Ethereum address, which
    /// must be funded for gas before it can relay.
    const KMS_FLAGS: [&'static str; 5] = [
        "--gcp-kms-project",
        "--gcp-kms-location",
        "--gcp-kms-key-ring",
        "--gcp-kms-key-name",
        "--gcp-kms-key-version",
    ];

    /// Resolve the signing method from CLI flags.
    ///
    /// Returns `Err` with an operator-facing message rather than panicking, so
    /// that a misconfiguration produces a clean exit instead of a backtrace.
    fn from_args(args: &Args) -> Result<Self, String> {
        let supplied: Vec<&str> = Self::KMS_FLAGS
            .iter()
            .zip([
                args.flag_gcp_kms_project.is_some(),
                args.flag_gcp_kms_location.is_some(),
                args.flag_gcp_kms_key_ring.is_some(),
                args.flag_gcp_kms_key_name.is_some(),
                args.flag_gcp_kms_key_version.is_some(),
            ])
            .filter(|(_, present)| *present)
            .map(|(name, _)| *name)
            .collect();
        let missing: Vec<&str> = Self::KMS_FLAGS
            .iter()
            .filter(|f| !supplied.contains(f))
            .copied()
            .collect();

        match (args.flag_ethereum_key.as_ref(), supplied.is_empty()) {
            (Some(_), false) => Err(format!(
                "Both signing methods configured. --ethereum-key was supplied together \
                 with {}. Supply exactly one signing method so it is unambiguous which \
                 key signs.",
                supplied.join(", ")
            )),
            (None, true) => Err(format!(
                "No signing method configured. Supply either --ethereum-key or all of {}.",
                Self::KMS_FLAGS.join(", ")
            )),
            (Some(key), true) => {
                warn!(
                    "Signing with a raw --ethereum-key. This key is visible in the process \
                     table to every local user on this host. Prefer {}, so that the key \
                     never leaves KMS.",
                    Self::KMS_FLAGS.join(", ")
                );
                let wallet: EthWallet = key
                    .parse()
                    .map_err(|e| format!("Invalid Ethereum private key: {}", e))?;
                Ok(SigningConfig::Local(Box::new(wallet)))
            }
            (None, false) => {
                if !missing.is_empty() {
                    return Err(format!(
                        "Incomplete GCP KMS configuration. Supplied {}, but these are also \
                         required: {}.",
                        supplied.join(", "),
                        missing.join(", ")
                    ));
                }
                Ok(SigningConfig::GcpKms {
                    project: args.flag_gcp_kms_project.clone().unwrap(),
                    location: args.flag_gcp_kms_location.clone().unwrap(),
                    key_ring: args.flag_gcp_kms_key_ring.clone().unwrap(),
                    key_name: args.flag_gcp_kms_key_name.clone().unwrap(),
                    key_version: args.flag_gcp_kms_key_version.unwrap(),
                })
            }
        }
    }

    async fn into_signer(self, chain_id: u64) -> Result<SignerType, String> {
        match self {
            SigningConfig::Local(wallet) => Ok(SignerType::Local(*wallet).with_chain_id(chain_id)),
            SigningConfig::GcpKms {
                project,
                location,
                key_ring,
                key_name,
                key_version,
            } => {
                info!(
                    "Using GCP KMS signer: projects/{}/locations/{}/keyRings/{}/cryptoKeys/{}/cryptoKeyVersions/{}",
                    project, location, key_ring, key_name, key_version
                );
                SignerType::new_gcp_kms(
                    &project,
                    &location,
                    &key_ring,
                    key_name,
                    key_version,
                    chain_id,
                )
                .await
                .map_err(|e| format!("Failed to construct GCP KMS signer: {}", e))
            }
        }
    }
}

#[tokio::main]
async fn main() {
    env_logger::Builder::from_env(Env::default().default_filter_or("info")).init();
    // On Linux static builds we need to probe ssl certs path to be able to
    // do TLS stuff.
    #[allow(deprecated)]
    openssl_probe::init_ssl_cert_env_vars();

    let args: Args = Docopt::new(USAGE.as_str())
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());
    // Decide how we will sign before doing any network work, so a
    // misconfiguration fails immediately rather than after connecting.
    let signing_config = SigningConfig::from_args(&args).unwrap_or_else(|e| {
        error!("{}", e);
        std::process::exit(1);
    });

    let gravity_contract_address: EthAddress = args
        .flag_contract_address
        .parse()
        .expect("Invalid contract address!");

    let connections = create_rpc_connections(
        args.flag_address_prefix,
        Some(args.flag_cosmos_grpc),
        Some(args.flag_ethereum_rpc),
        LOOP_SPEED,
    )
    .await;
    let provider = connections.eth_provider.clone().unwrap();
    let chain_id = provider
        .get_chainid()
        .await
        .expect("Could not retrieve chain ID during relayer start");
    let chain_id = downcast_to_u64(chain_id).expect("Chain ID overflowed when downcasting to u64");

    // GcpKmsSigner binds chain_id at construction, so the signer can only be
    // built once the chain ID is known.
    let ethereum_wallet = signing_config
        .into_signer(chain_id)
        .await
        .unwrap_or_else(|e| {
            error!("{}", e);
            std::process::exit(1);
        });

    let eth_client = SignerMiddleware::new(provider, ethereum_wallet);
    let eth_client = Arc::new(eth_client);

    let public_eth_key = eth_client.address();
    info!("Starting Gravity Relayer");
    info!("Ethereum Address: {}", format_eth_address(public_eth_key));

    let contact = connections.contact.clone().unwrap();

    // check if the cosmos node is syncing, if so wait for it
    // we can't move any steps above this because they may fail on an incorrect
    // historic chain state while syncing occurs
    wait_for_cosmos_node_ready(&contact).await;
    check_for_eth(public_eth_key, eth_client.clone()).await;

    relayer_main_loop(
        eth_client,
        connections.grpc.unwrap(),
        gravity_contract_address,
        1.1f32,
        1.1f32,
    )
    .await
}

#[cfg(test)]
mod tests {
    use super::*;

    /// Args with neither signing method configured. Non-signing fields are
    /// irrelevant to SigningConfig resolution but must be present.
    fn bare_args() -> Args {
        Args {
            flag_ethereum_key: None,
            flag_gcp_kms_project: None,
            flag_gcp_kms_location: None,
            flag_gcp_kms_key_ring: None,
            flag_gcp_kms_key_name: None,
            flag_gcp_kms_key_version: None,
            flag_cosmos_grpc: "http://localhost:9090".to_string(),
            flag_address_prefix: "somm".to_string(),
            flag_ethereum_rpc: "http://localhost:8545".to_string(),
            flag_contract_address: "0x69592e6f9d21989a043646fE8225da2600e5A0f7".to_string(),
        }
    }

    fn full_kms_args() -> Args {
        Args {
            flag_gcp_kms_project: Some("somm-production-validator".to_string()),
            flag_gcp_kms_location: Some("us-central1".to_string()),
            flag_gcp_kms_key_ring: Some("somm-production-keyring".to_string()),
            flag_gcp_kms_key_name: Some("sommelier-primary-eth-signer".to_string()),
            flag_gcp_kms_key_version: Some(1),
            ..bare_args()
        }
    }

    #[test]
    fn no_signing_method_is_rejected() {
        let err = SigningConfig::from_args(&bare_args()).unwrap_err();
        assert!(err.contains("No signing method configured"), "{}", err);
    }

    /// A valid 32-byte secp256k1 key.
    const TEST_KEY: &str = "0x0000000000000000000000000000000000000000000000000000000000000001";

    #[test]
    fn raw_key_selects_local() {
        let args = Args {
            flag_ethereum_key: Some(TEST_KEY.to_string()),
            ..bare_args()
        };
        match SigningConfig::from_args(&args).unwrap() {
            SigningConfig::Local(_) => {}
            other => panic!("expected Local signer, got {:?}", other),
        }
    }

    /// The `Local` variant holds a private key, so its Debug must never print
    /// it. This is the leak class the whole change exists to close.
    #[test]
    fn local_debug_redacts_the_key() {
        let args = Args {
            flag_ethereum_key: Some(TEST_KEY.to_string()),
            ..bare_args()
        };
        let cfg = SigningConfig::from_args(&args).unwrap();
        let rendered = format!("{:?}", cfg);
        assert!(!rendered.contains("0000000000000000"), "{}", rendered);
        assert!(rendered.contains("redacted"), "{}", rendered);
    }

    #[test]
    fn full_kms_flags_select_kms() {
        match SigningConfig::from_args(&full_kms_args()).unwrap() {
            SigningConfig::GcpKms {
                project,
                location,
                key_ring,
                key_name,
                key_version,
            } => {
                assert_eq!(project, "somm-production-validator");
                assert_eq!(location, "us-central1");
                assert_eq!(key_ring, "somm-production-keyring");
                assert_eq!(key_name, "sommelier-primary-eth-signer");
                assert_eq!(key_version, 1);
            }
            _ => panic!("expected GcpKms signer"),
        }
    }

    #[test]
    fn both_methods_are_rejected() {
        let args = Args {
            flag_ethereum_key: Some("0xdeadbeef".to_string()),
            ..full_kms_args()
        };
        let err = SigningConfig::from_args(&args).unwrap_err();
        assert!(err.contains("Both signing methods configured"), "{}", err);
    }

    /// Every proper subset of the KMS flags must be rejected, and the message
    /// must name what is missing. This is the case most likely to bite an
    /// operator mid-migration, so it is checked exhaustively rather than for a
    /// single representative subset.
    #[test]
    fn every_partial_kms_combination_is_rejected() {
        let full = full_kms_args();
        // Bit i set means "supply flag i". 0 is the no-method case, which is
        // covered separately; 31 is the complete set.
        for mask in 1..31u8 {
            let args = Args {
                flag_gcp_kms_project: (mask & 1 != 0)
                    .then(|| full.flag_gcp_kms_project.clone().unwrap()),
                flag_gcp_kms_location: (mask & 2 != 0)
                    .then(|| full.flag_gcp_kms_location.clone().unwrap()),
                flag_gcp_kms_key_ring: (mask & 4 != 0)
                    .then(|| full.flag_gcp_kms_key_ring.clone().unwrap()),
                flag_gcp_kms_key_name: (mask & 8 != 0)
                    .then(|| full.flag_gcp_kms_key_name.clone().unwrap()),
                flag_gcp_kms_key_version: (mask & 16 != 0).then_some(1),
                ..bare_args()
            };
            let err = SigningConfig::from_args(&args)
                .err()
                .unwrap_or_else(|| panic!("mask {} should not have produced a signer", mask));
            assert!(
                err.contains("Incomplete GCP KMS configuration"),
                "mask {}: {}",
                mask,
                err
            );
        }
    }

    /// The key version must never be silently defaulted: a KMS signing address
    /// is version-specific, so a default would pin to a version the operator
    /// did not choose and whose Ethereum address may hold no gas.
    #[test]
    fn missing_key_version_is_rejected() {
        let args = Args {
            flag_gcp_kms_key_version: None,
            ..full_kms_args()
        };
        let err = SigningConfig::from_args(&args).unwrap_err();
        assert!(err.contains("--gcp-kms-key-version"), "{}", err);
    }

    /// An invalid raw key must surface as a clean error rather than a panic,
    /// and must do so during argument handling, before any network work. The
    /// original code parsed the key at the top of main; KMS support moved
    /// signer construction below chain-ID resolution, so this guards against
    /// the key check drifting back after connection setup.
    #[test]
    fn invalid_raw_key_is_rejected_during_arg_handling() {
        let args = Args {
            flag_ethereum_key: Some("not-a-key".to_string()),
            ..bare_args()
        };
        let err = SigningConfig::from_args(&args).unwrap_err();
        assert!(err.contains("Invalid Ethereum private key"), "{}", err);
    }
}
