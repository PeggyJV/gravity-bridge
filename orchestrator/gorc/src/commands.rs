//! Gorc Subcommands
//!
//! This is where you specify the subcommands of your application.
//!
//! The default application comes with two subcommands:
//!
//! - `start`: launches the application
//! - `version`: print application version
//!
//! See the `impl Configurable` below for how to specify the path to the
//! application's configuration file.

mod start;
mod version;
mod tx;
mod query;
mod tests;
mod keys;

use self::{
    start::StartCmd, version::VersionCmd, tx::TxCmd,
    query::QueryCmd, tests::TestsCmd, keys::KeysCmd, 
};
use crate::config::GorcConfig;
use abscissa_core::{
    config::Override, Command, Configurable, FrameworkError, Help, Options, Runnable,
};
use std::path::PathBuf;

/// Gorc Configuration Filename
pub const CONFIG_FILE: &str = "gorc.toml";

/// Gorc Subcommands
#[derive(Command, Debug, Options, Runnable)]
pub enum GorcCmd {
    #[options(help = "create transactions on either ethereum or cosmos chains")]
    Tx(TxCmd),
    
    #[options(help = "query state on either ethereum or cosmos chains")]
    Query(QueryCmd),

    #[options(help = "run tests against configured chains")]
    Tests(TestsCmd),

    #[options(help = "start the application")]
    Start(StartCmd),

    #[options(help = "start the application")]
    Keys(KeysCmd),

    #[options(help = "get usage information")]
    Help(Help<Self>),
    
    #[options(help = "display version information")]
    Version(VersionCmd),
}

/// This trait allows you to define how application configuration is loaded.
impl Configurable<GorcConfig> for GorcCmd {
    /// Location of the configuration file
    fn config_path(&self) -> Option<PathBuf> {
        // Check if the config file exists, and if it does not, ignore it.
        // If you'd like for a missing configuration file to be a hard error
        // instead, always return `Some(CONFIG_FILE)` here.
        let filename = PathBuf::from(CONFIG_FILE);

        if filename.exists() {
            Some(filename)
        } else {
            None
        }
    }

    /// Apply changes to the config after it's been loaded, e.g. overriding
    /// values in a config file using command-line options.
    ///
    /// This can be safely deleted if you don't want to override config
    /// settings from command-line options.
    fn process_config(
        &self,
        config: GorcConfig,
    ) -> Result<GorcConfig, FrameworkError> {
        match self {
            GorcCmd::Start(cmd) => cmd.override_config(config),
            _ => Ok(config),
        }
    }
}
