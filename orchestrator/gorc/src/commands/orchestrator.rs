mod start;

use abscissa_core::{clap::Parser, Command, Runnable};

/// Management commannds for the orchestrator
#[derive(Command, Debug, Parser, Runnable)]
pub enum OrchestratorCmd {
    Start(start::StartCommand),
}
