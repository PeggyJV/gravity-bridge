use abscissa_core::{clap::Parser, Command, Runnable};

/// Run tests against configured chains
#[derive(Command, Debug, Parser)]
pub enum TestsCmd {
    Runner(Runner),
}

impl Runnable for TestsCmd {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Parser)]
pub struct Runner {
    free: Vec<String>,

    #[clap(short, long)]
    help: bool,
}

impl Runnable for Runner {
    fn run(&self) {}
}
