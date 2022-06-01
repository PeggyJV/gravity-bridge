use super::show::ShowEthKeyCmd;
use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use k256::pkcs8::ToPrivateKey;
use rand_core::OsRng;
use signatory::FsKeyStore;
use std::path;

/// Add a new Eth Key
#[derive(Command, Debug, Default, Parser)]
#[clap(
    long_about = "DESCRIPTION \n\n Create a new Eth Key.\n This command creates a new Eth key. It has an overwrite option, which if set to true, overwrites\n an existing key in the keystore with the same keyname."
)]
pub struct AddEthKeyCmd {
    /// Eth keyname.
    pub name: String,

    /// Overwrite key with the same name in the keystore when set to true. Takes a Boolean.
    #[clap(short, long)]
    pub overwrite: bool,

    /// Show private key after creation of key.
    #[clap(short, long)]
    show_private_key: bool,
}

// Entry point for `gorc keys eth add [name]`
// - [name] required; key name
impl Runnable for AddEthKeyCmd {
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

        let mnemonic = bip32::Mnemonic::random(&mut OsRng, Default::default());
        eprintln!("**Important** record this bip39-mnemonic in a safe place:");
        println!("{}", mnemonic.phrase());

        let seed = mnemonic.to_seed("");

        let path = config.ethereum.key_derivation_path.trim();
        let path = path
            .parse::<bip32::DerivationPath>()
            .expect("Could not parse derivation path");

        let key = bip32::XPrv::derive_from_path(seed, &path).expect("Could not derive key");
        let key = k256::SecretKey::from(key.private_key());
        let key = key
            .to_pkcs8_der()
            .expect("Could not PKCS8 encod private key");

        keystore.store(&name, &key).expect("Could not store key");

        let show_cmd = ShowEthKeyCmd {
            name: name.to_string(),
            show_private_key: self.show_private_key,
            show_name: false,
        };
        show_cmd.run();
    }
}
