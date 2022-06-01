mod cosmos;
mod eth;

use abscissa_core::{clap::Parser, Command, Runnable};

use crate::commands::keys::cosmos::CosmosKeysCmd;
use crate::commands::keys::eth::EthKeysCmd;

/// Key management commands for Ethereum and Cosmos

#[derive(Command, Debug, Parser, Runnable)]
pub enum KeysCmd {
    /// Manage Ethereum keys.
    #[clap(subcommand)]
    Eth(EthKeysCmd),
    /// Manage Cosmos keys.
    #[clap(subcommand)]
    Cosmos(CosmosKeysCmd),
}
