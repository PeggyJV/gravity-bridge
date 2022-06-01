use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};

/// Show an Eth Key
#[derive(Command, Debug, Default, Parser)]
#[clap(
    long_about = "DESCRIPTION \n\n Show details of an Eth key in the keystore.\n This command shows details of an Eth key in the keystore, it takes the name of the key."
)]
pub struct ShowEthKeyCmd {
    /// Cosmos keyname
    pub name: String,

    /// Show private key when set to true. Takes a Boolean.
    #[clap(short, long)]
    pub show_private_key: bool,

    /// Show key name when set to true. Takes a Boolean.
    #[clap(short = 'n', long)]
    pub show_name: bool,
}

// Entry point for `gorc keys eth show [name]`
impl Runnable for ShowEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let name = self.name.clone();
        // TODO(bolten): is ethers wallet even capable of printing the public and
        // private keys? for now, leaving load_clarity_key in config.rs and
        // maintaining the functionality here
        let key = config.load_clarity_key(name.clone());

        let pub_key = key.to_public_key().expect("Could not build public key");

        if self.show_name {
            print!("{}\t", name);
        }

        if self.show_private_key {
            println!("{}\t{}", pub_key, key);
        } else {
            println!("{}", pub_key);
        }
    }
}
