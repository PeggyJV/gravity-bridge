//! `cosmos subcommands` subcommand

use crate::{application::APP, prelude::*};
use abscissa_core::{clap::Parser, Command, Runnable};

/// Query cosmos chain
#[derive(Command, Debug, Parser)]
pub enum Cosmos {
    /// Query the balance in a Cosmos account.
    Balance(Balance),

    /// Query gravity keys.
    GravityKeys(GravityKeys),
}

impl Runnable for Cosmos {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Parser)]
pub struct Balance {
    /// Cosmos keyname
    key_name: String,

    #[clap(short, long)]
    help: bool,
}

impl Runnable for Balance {
    fn run(&self) {
        let _key_name = self.key_name.clone();
    }
}

#[derive(Command, Debug, Parser)]
pub struct GravityKeys {
    /// Gravity keyname
    key_name: String,

    #[clap(short, long)]
    help: bool,
}

impl Runnable for GravityKeys {
    /// Start the application.
    fn run(&self) {
        let _key_name = self.key_name.clone();

        abscissa_tokio::run(&APP, async { unimplemented!() }).unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}
