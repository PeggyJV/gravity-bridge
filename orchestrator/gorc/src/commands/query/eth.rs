//! `eth subcommands` subcommand

use crate::{application::APP, prelude::*};
use abscissa_core::{Clap, Command, Runnable};

/// Query Eth chain
#[derive(Command, Debug, Clap)]
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

#[derive(Command, Debug, Clap)]
pub struct Balance {
    free: Vec<String>,

    #[clap(short, long)]
    help: bool,
}

impl Runnable for Balance {
    fn run(&self) {
        assert!(self.free.len() == 1);
        let _key_name = self.free[0].clone();

        abscissa_tokio::run(&APP, async { unimplemented!() }).unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}

#[derive(Command, Debug, Clap)]
pub struct Contract {
    #[clap(short, long)]
    help: bool,
}

impl Runnable for Contract {
    /// Start the application.
    fn run(&self) {}
}
