use super::show::ShowCosmosKeyCmd;
use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use std::path::Path;

/// List all Cosmos Keys
#[derive(Command, Debug, Default, Parser)]
#[clap(
    long_about = "DESCRIPTION \n\n List all Cosmos keys in keystore.\n This command lists all Cosmos keys and their addresses from the keystore."
)]
pub struct ListCosmosKeyCmd {}

// Entry point for `gorc keys cosmos list`
impl Runnable for ListCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = Path::new(&config.keystore);

        for entry in keystore.read_dir().expect("Could not read keystore") {
            let path = entry.unwrap().path();
            if path.is_file() {
                if let Some(extension) = path.extension() {
                    if extension == "pem" {
                        let name = path.file_stem().unwrap();
                        let name = name.to_str().unwrap();
                        let args = vec![name.to_string()];
                        let show_cmd = ShowCosmosKeyCmd { args };
                        show_cmd.run();
                    }
                }
            }
        }
    }
}
