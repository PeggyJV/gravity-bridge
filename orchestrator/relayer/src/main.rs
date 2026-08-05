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
    flag_ethereum_key: String,
    flag_cosmos_grpc: String,
    flag_address_prefix: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage: {} --ethereum-key=<key> --cosmos-grpc=<url> --address-prefix=<prefix> --ethereum-rpc=<url> --contract-address=<addr>
        Options:
            -h --help                    Show this screen.
            --ethereum-key=<ekey>        An Ethereum private key containing non-trivial funds, or the literal value kms
            --cosmos-grpc=<gurl>         The Cosmos gRPC url
            --address-prefix=<prefix>    The prefix for addresses on this Cosmos chain
            --ethereum-rpc=<eurl>        The Ethereum RPC url, Geth light clients work and sync fast
            --contract-address=<addr>    The Ethereum contract address for Gravity
        About:
            The Gravity relayer component, responsible for relaying data from the Cosmos blockchain
            to the Ethereum blockchain, cosmos key and fees are optional since they are only used
            to request the creation of batches or validator sets to relay.
            for Althea-Gravity.
            Signing. Passing a raw Ethereum private key is discouraged, because it is
            visible in the process table to every local user on the host. Instead pass
            the literal value kms, which signs through Google Cloud Kms so that the
            private key never leaves Kms. That mode reads its configuration from these
            environment variables, all of which are required.
            Gravity_Gcp_Kms_Project, Gravity_Gcp_Kms_Location, Gravity_Gcp_Kms_Key_Ring,
            Gravity_Gcp_Kms_Key_Name and Gravity_Gcp_Kms_Key_Version, upper cased.
            The key must be a secp256k1 signing key. Each key version has its own
            Ethereum address, which pays relayer gas and must be funded before use.
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
    Local(String),
    /// Google Cloud KMS. No key material is held by this process.
    GcpKms {
        project: String,
        location: String,
        key_ring: String,
        key_name: String,
        key_version: u64,
    },
}

impl SigningConfig {
    /// All required. The key version is deliberately required rather than
    /// defaulted: a GCP KMS signing address is version-specific, so silently
    /// pinning to version 1 would both survive a key rotation without picking
    /// it up (relaying stops when version 1 is disabled) and hide the fact
    /// that moving versions changes the Ethereum address, which must be
    /// funded for gas before it can relay.
    const KMS_VARS: [&'static str; 5] = [
        "GRAVITY_GCP_KMS_PROJECT",
        "GRAVITY_GCP_KMS_LOCATION",
        "GRAVITY_GCP_KMS_KEY_RING",
        "GRAVITY_GCP_KMS_KEY_NAME",
        "GRAVITY_GCP_KMS_KEY_VERSION",
    ];

    /// Sentinel value for --ethereum-key that selects the GCP KMS signer.
    const KMS_SENTINEL: &'static str = "kms";

    fn from_args(args: &Args) -> Self {
        if args.flag_ethereum_key != Self::KMS_SENTINEL {
            warn!(
                "Signing with a raw --ethereum-key. This key is visible in the process \
                 table to every local user on this host. Prefer GCP KMS: pass \
                 --ethereum-key={} and set {}, so the key never leaves KMS.",
                Self::KMS_SENTINEL,
                Self::KMS_VARS.join(", ")
            );
            return SigningConfig::Local(args.flag_ethereum_key.clone());
        }

        let missing: Vec<&str> = Self::KMS_VARS
            .iter()
            .filter(|v| std::env::var(v).is_err())
            .copied()
            .collect();
        if !missing.is_empty() {
            panic!(
                "--ethereum-key={} selects the GCP KMS signer, but these required \
                 environment variables are unset: {}.",
                Self::KMS_SENTINEL,
                missing.join(", ")
            );
        }

        let get = |v: &str| std::env::var(v).unwrap();
        let key_version = get("GRAVITY_GCP_KMS_KEY_VERSION")
            .parse::<u64>()
            .expect("GRAVITY_GCP_KMS_KEY_VERSION must be a positive integer");

        SigningConfig::GcpKms {
            project: get("GRAVITY_GCP_KMS_PROJECT"),
            location: get("GRAVITY_GCP_KMS_LOCATION"),
            key_ring: get("GRAVITY_GCP_KMS_KEY_RING"),
            key_name: get("GRAVITY_GCP_KMS_KEY_NAME"),
            key_version,
        }
    }

    async fn into_signer(self, chain_id: u64) -> SignerType {
        match self {
            SigningConfig::Local(key) => {
                let wallet: EthWallet = key.parse().expect("Invalid Ethereum private key!");
                SignerType::Local(wallet).with_chain_id(chain_id)
            }
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
                .expect("Failed to construct GCP KMS signer")
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
    let signing_config = SigningConfig::from_args(&args);

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
    let ethereum_wallet = signing_config.into_signer(chain_id).await;

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
