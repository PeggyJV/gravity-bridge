mod erc20;
use erc20::Erc20;

use abscissa_core::{Clap, Command, Runnable};

/// Provides tools for contract deployment
#[derive(Command, Debug, Clap, Runnable)]
pub enum DeployCmd {
    Erc20(Erc20),
}
