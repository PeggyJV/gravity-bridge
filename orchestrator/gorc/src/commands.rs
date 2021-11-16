//! Gorc Subcommands
//! This is where you specify the subcommands of your application.

mod cosmos_to_eth;
mod deploy;
mod eth_to_cosmos;
mod keys;
mod orchestrator;
mod print_config;
mod query;
mod sign_delegate_keys;
mod tests;
mod tx;

use crate::config::GorcConfig;
use abscissa_core::{Clap, Command, Configurable, FrameworkError, Runnable};
use std::path::PathBuf;

/// Gorc Configuration Filename
pub const CONFIG_FILE: &str = "gorc.toml";

/// Gorc Subcommands
#[derive(Command, Debug, Clap, Runnable)]
pub enum GorcCmd {
    CosmosToEth(cosmos_to_eth::CosmosToEthCmd),

    #[clap(subcommand)]
    Deploy(deploy::DeployCmd),

    EthToCosmos(eth_to_cosmos::EthToCosmosCmd),

    #[clap(subcommand)]
    Keys(keys::KeysCmd),

    #[clap(subcommand)]
    Orchestrator(orchestrator::OrchestratorCmd),

    PrintConfig(print_config::PrintConfigCmd),

    #[clap(subcommand)]
    Query(query::QueryCmd),

    SignDelegateKeys(sign_delegate_keys::SignDelegateKeysCmd),

    #[clap(subcommand)]
    Tests(tests::TestsCmd),

    #[clap(subcommand)]
    Tx(tx::TxCmd),
}

/// Entry point for the application. It needs to be a struct to allow using subcommands!
#[derive(Command, Debug, Clap)]
#[clap(author, about, version)]
pub struct EntryPoint {
    #[clap(subcommand)]
    cmd: GorcCmd,

    /// Enable verbose logging
    #[clap(short, long)]
    pub verbose: bool,

    /// Use the specified config file
    #[clap(short, long)]
    pub config: Option<String>,
}

impl Runnable for EntryPoint {
    fn run(&self) {
        self.cmd.run()
    }
}

/// This trait allows you to define how application configuration is loaded.
impl Configurable<GorcConfig> for EntryPoint {
    /// Location of the configuration file
    fn config_path(&self) -> Option<PathBuf> {
        // Check if the config file exists, and if it does not, ignore it.
        // If you'd like for a missing configuration file to be a hard error
        // instead, always return `Some(CONFIG_FILE)` here.
        let filename = self
            .config
            .as_ref()
            .map(PathBuf::from)
            .unwrap_or_else(|| CONFIG_FILE.into());

        if filename.exists() {
            Some(filename)
        } else {
            None
        }
    }

    /// Apply changes to the config after it's been loaded, e.g. overriding
    /// values in a config file using command-line options.
    ///
    /// This can be safely deleted if you don't want to override config
    /// settings from command-line options.
    fn process_config(&self, config: GorcConfig) -> Result<GorcConfig, FrameworkError> {
        match &self.cmd {
            _ => Ok(config),
        }
    }
}
