use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use gravity::deep_space::PrivateKey;

/// Show a Cosmos Key
#[derive(Command, Debug, Default, Parser)]
pub struct ShowCosmosKeyCmd {
    pub args: Vec<String>,
}

// Entry point for `gorc keys cosmos show [name]`
impl Runnable for ShowCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let name = self.args.get(0).expect("name is required");
        let key = config.load_deep_space_key(name.clone());

        let address = key
            .to_address(config.cosmos.prefix.trim())
            .expect("Could not generate public key");

        println!("{}\t{}", name, address)
    }
}
