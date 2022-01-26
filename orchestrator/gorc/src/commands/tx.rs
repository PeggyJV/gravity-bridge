//! `tx` subcommand

mod cosmos;

mod eth;

use abscissa_core::{clap::Parser, Command, Runnable};

/// Create transactions on either ethereum or cosmos chains
#[derive(Command, Debug, Parser)]
pub enum TxCmd {
    #[clap(subcommand)]
    Cosmos(cosmos::Cosmos),

    #[clap(subcommand)]
    Eth(eth::Eth),
}

impl Runnable for TxCmd {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}
