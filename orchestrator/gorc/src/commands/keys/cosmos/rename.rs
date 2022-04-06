use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use std::path;

/// Rename a Cosmos Key
#[derive(Command, Debug, Default, Parser)]
#[clap(
    long_about = "DESCRIPTION \n\n Rename a Cosmos key.\n This command will rename a Cosmos key in the keystore. It takes the existing keyname and new \n keyname."
)]
pub struct RenameCosmosKeyCmd {
    /// Cosmos keyname in keystore.
    pub name: String,

    /// New keyname to replace name in keystore.
    pub new_name: String,

    /// Overwrite key with the same name in the keystore when set to true. Takes a Boolean.
    #[clap(short, long)]
    pub overwrite: bool,
}

/// The `gorc keys cosmos rename [name] [new-name]` subcommand: show keys
impl Runnable for RenameCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();

        let name = self.name.parse().expect("Could not parse name");

        let new_name = self.new_name.parse().expect("Could not parse new_name");
        if let Ok(_info) = keystore.info(&new_name) {
            if !self.overwrite {
                println!("Key already exists, exiting.");
                return;
            }
        }

        let key = keystore.load(&name).expect("Could not load key");
        keystore.store(&new_name, &key).unwrap();
        keystore.delete(&name).unwrap();
    }
}
