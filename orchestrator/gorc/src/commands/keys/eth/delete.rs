use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use signatory::FsKeyStore;
use std::path;

/// Delete an Eth Key
#[derive(Command, Debug, Default, Parser)]
#[clap(
    long_about = "DESCRIPTION \n\n Delete an Eth Key.\n This command deletes an Eth key from your keystore when provided with the keyname."
)]
pub struct DeleteEthKeyCmd {
    /// Eth keyname.
    pub name: String,
}

// Entry point for `gorc keys eth delete [name]`
// - [name] required; key name
impl Runnable for DeleteEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore = FsKeyStore::create_or_open(keystore).expect("Could not open keystore");
        let name = self.name.parse().expect("Could not parse name");
        keystore.delete(&name).expect("Could not delete key");
    }
}
