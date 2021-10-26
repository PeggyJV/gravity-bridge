mod start;

use abscissa_core::{Clap, Command, Runnable};

/// Management commannds for the orchestrator
#[derive(Command, Debug, Clap, Runnable)]
pub enum OrchestratorCmd {
    Start(start::StartCommand),
}
