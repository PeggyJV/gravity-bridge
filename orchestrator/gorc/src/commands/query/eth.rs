//! `eth subcommands` subcommand

use crate::{application::APP, prelude::*};
use abscissa_core::{clap::Parser, Command, Runnable};

/// Query Eth chain
#[derive(Command, Debug, Parser)]
pub enum Eth {
    Balance(Balance),

    Contract(Contract),
}

impl Runnable for Eth {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Parser)]
pub struct Balance {
    /// Eth keyname
    key_name: String,

    #[clap(short, long)]
    help: bool,
}

impl Runnable for Balance {
    fn run(&self) {
        let _key_name = self.key_name.clone();

        abscissa_tokio::run(&APP, async { unimplemented!() }).unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}

#[derive(Command, Debug, Parser)]
pub struct Contract {
    #[clap(short, long)]
    help: bool,
}

impl Runnable for Contract {
    /// Start the application.
    fn run(&self) {}
}
