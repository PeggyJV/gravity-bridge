use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use signatory;
use std::path::Path;
use deep_space;

#[derive(Command, Debug, Default, Options)]
pub struct ShowCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "show [name]")]
    pub name: String,
}

/// The `gorc keys cosmos show [name]` subcommand: show keys
impl Runnable for ShowCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();
        let key_name = &self.name.parse().unwrap();

        let key = keystore.load(&key_name).expect("Could not load key");
        let key = key
            .to_pem()
            .parse::<k256::elliptic_curve::SecretKey<k256::Secp256k1>>()
            .expect("Could not parse key");
        let key = deep_space::private_key::PrivateKey::from_secret(&key.to_bytes());

        let address = key.to_address("cosmos").expect("Could not generate public key");
        println!("{:?}", address)
    }
}