use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use signatory::FsKeyStore;
use std::path::Path;

/// Delete a Cosmos Key
#[derive(Command, Debug, Default, Parser)]
#[clap(
    long_about = "DESCRIPTION \n\n Delete a Cosmos Key.\n This command deletes a Cosmos key from your keystore when provided with the keyname."
)]
pub struct DeleteCosmosKeyCmd {
    /// Cosmos keyname in keystore.
    pub name: String,
}

/// The `gorc keys cosmos delete [name] ` subcommand: delete the given key
impl Runnable for DeleteCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        // Path where key is stored.
        let keystore = Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();
        // Collect key name from args.
        let name = self.name.parse().expect("Could not parse name");
        // Delete keyname after locating file from path and key name.
        let _delete_key = FsKeyStore::delete(&keystore, &name).unwrap();
    }
}
