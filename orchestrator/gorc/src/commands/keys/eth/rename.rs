use crate::application::APP;
use abscissa_core::{Application, Clap, Command, Runnable};
use signatory::FsKeyStore;
use std::path;

/// Rename an Eth Key
#[derive(Command, Debug, Default, Clap)]
pub struct RenameEthKeyCmd {
    pub args: Vec<String>,

    #[clap(short, long)]
    pub overwrite: bool,
}

// Entry point for `gorc keys eth rename [name] [new-name]`
impl Runnable for RenameEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore = FsKeyStore::create_or_open(keystore).expect("Could not open keystore");

        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");

        let new_name = self.args.get(1).expect("new-name is required");
        let new_name = new_name.parse().expect("Could not parse new-name");
        if let Ok(_info) = keystore.info(&new_name) {
            if !self.overwrite {
                eprintln!("Key already exists, exiting.");
                return;
            }
        }

        let key = keystore.load(&name).expect("Could not load key");
        keystore
            .store(&new_name, &key)
            .expect("Could not store key");
        keystore.delete(&name).expect("Could not delete key");
    }
}
