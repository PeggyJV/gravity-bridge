use super::show::ShowCosmosKeyCmd;
use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use k256::pkcs8::ToPrivateKey;
use signatory::FsKeyStore;
use std::path;

/// Recover a Cosmos Key
#[derive(Command, Debug, Default, Parser)]
#[clap(
    long_about = "DESCRIPTION \n\n Recover an external Cosmos key.\n This command will recover a Cosmos key, storing it in the keystore. \n It takes a keyname and bip39-mnemonic."
)]
pub struct RecoverCosmosKeyCmd {
    /// Cosmos keyname.
    pub name: String,

    /// Overwrite key with the same name in the keystore when set to true. Takes a Boolean.
    #[clap(short, long)]
    pub overwrite: bool,

    /// bip39-mnemonic optional. When absent you'll be prompted to enter it.
    pub mnemonic: Option<String>,
}

// `gorc keys cosmos recover [name] (bip39-mnemonic)`
// - [name] required; key name
// - (bip39-mnemonic) optional; when absent the user will be prompted to enter it
impl Runnable for RecoverCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore = FsKeyStore::create_or_open(keystore).expect("Could not open keystore");

        let name = self.name.parse().expect("Could not parse name");
        if let Ok(_info) = keystore.info(&name) {
            if !self.overwrite {
                eprintln!("Key already exists, exiting.");
                return;
            }
        }

        let mnemonic = match self.mnemonic.clone() {
            Some(mnemonic) => mnemonic,
            None => rpassword::read_password_from_tty(Some("> Enter your bip39-mnemonic:\n"))
                .expect("Could not read mnemonic"),
        };

        let mnemonic = bip32::Mnemonic::new(mnemonic.trim(), Default::default())
            .expect("Could not parse mnemonic");

        let seed = mnemonic.to_seed("");

        let path = config.cosmos.key_derivation_path.clone();
        let path = path
            .parse::<bip32::DerivationPath>()
            .expect("Could not parse derivation path");

        let key = bip32::XPrv::derive_from_path(seed, &path).expect("Could not derive key");
        let key = k256::SecretKey::from(key.private_key());
        let key = key
            .to_pkcs8_der()
            .expect("Could not PKCS8 encod private key");

        keystore.store(&name, &key).expect("Could not store key");

        let name = name.to_string();
        let show_cmd = ShowCosmosKeyCmd { name };
        show_cmd.run();
    }
}
