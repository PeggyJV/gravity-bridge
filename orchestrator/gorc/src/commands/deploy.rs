mod erc20;
use erc20::Erc20;

use abscissa_core::{clap::Parser, Command, Runnable};

/// Provides tools for contract deployment
#[derive(Command, Debug, Parser, Runnable)]
pub enum DeployCmd {
    Erc20(Erc20),
}
