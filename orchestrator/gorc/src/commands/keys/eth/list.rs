use super::show::ShowEthKeyCmd;
use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use std::path;

/// List all Eth Keys
#[derive(Command, Debug, Default, Parser)]
#[clap(
    long_about = "DESCRIPTION \n\n List all Eth keys in keystore.\n This command lists all Eth keys and their addresses from the keystore."
)]
pub struct ListEthKeyCmd {
    /// Show private key when set to true. Takes a Boolean.
    #[clap(short, long)]
    pub show_private_key: bool,
}

// Entry point for `gorc keys eth list`
impl Runnable for ListEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);

        for entry in keystore.read_dir().expect("Could not read keystore") {
            let path = entry.unwrap().path();
            if path.is_file() {
                if let Some(extension) = path.extension() {
                    if extension == "pem" {
                        let name = path.file_stem().unwrap();
                        let name = name.to_str().unwrap();
                        let show_cmd = ShowEthKeyCmd {
                            name: name.to_string(),
                            show_private_key: self.show_private_key,
                            show_name: true,
                        };
                        show_cmd.run();
                    }
                }
            }
        }
    }
}
