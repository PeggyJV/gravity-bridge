// This library should provide functions that support command line workflows that
// - loads keys by name from a keystore
//   on failure, because the named key does not exist, it offers to create it
//   - offers choices here;
//     - using a new mnemonic
//     - using an existing mnemonic (aka recovery)
// - loads and converts keys to different representations used by different tools
// - leverages the signatory crate to manage the keystore
// - generates keys the hkd32 crate
// - provides abscissa commands (maybe)
//
// todos for this library:
// [ ] Function that loads a key from keystore by name then converts it to a given format
//     exact list of formats need to be pinned down; generally aiming to play nice with:
//     - clarity
//     - ethers
//     - contact / deep_space

// --------------------

// Wondering if this "lib" should become a binary / stand-alone cli tool.
// This create could still provide a reusable abscissa command, but I think
// but a standalone tool with usage like this feels worthwhile:
//   `key-tool show my-key`
//     - uses default location for the keystore (flag with default value)
//     - prints several representations of the key (clarity, ethers, etc.)
//     - if the keys isn't found
//       - Key not found. Do you want to add it?
//         1) With a recovery phrase?
//         2) Generate a new key?
//         3) Exit

#[cfg(test)]
mod tests {
    // NOTE these are not _real_ tests. I'm using them to bring together dependencies and learn how to use them.
    // NOTE run w/ `cargo test -- --nocapture` to see the println

    // Questions:
    // - How do I get from a mnemonic::Phrase to a pkcs8::PrivateKeyDocument (for keystore.store)?
    //   -- maybe derive_subkey?? or possibly to_seed?? (latter requires a password)
    // - What "display formats" do we care about?
    //   -- and where can I find examples??
    use hkd32;
    use rand_core::OsRng;
    use signatory::ecdsa::secp256k1;
    use signatory::FsKeyStore;
    use signatory::GeneratePkcs8;
    use pkcs8::ToPrivateKey;
    use bip32::{Mnemonic, XPrv};
    use std::path::Path;

    const RECOVERY_PHRASE: &str = "save able shop proud seek reflect prepare mechanic armor car core shuffle room axis file diet axis try secret evolve opinion prosper flush buyer";

    #[test]
    fn convert_pkcs8_to_clarity() {


        // clarity::private_key::PrivateKey::from_slice(_pkcs8_key.private_key_info().private_key).unwrap();

        // eg: let clarity_key = clarity::private_key::PrivateKey::from_slice(???).expect("Could not create clarity key");
    }

    #[test]
    fn convert_hkd32_mnemonic_phrase_to_pkcs8() {
        let mnemonic = Mnemonic::random(&mut OsRng, Default::default());
        println!{"**Important** write this mnemonic in a safe place.\n"}

        println!{"{}", mnemonic.phrase()};

        let seed = mnemonic.to_seed("MYPASSWORD"); // todo: password argument
        let xprv = XPrv::new(seed).unwrap();
        let private_key_der = k256::SecretKey::from(xprv.private_key()).to_pkcs8_der();
        let keystore_path = Path::new("keystore");
        if !keystore_path.exists() {
            FsKeyStore::create(keystore_path).unwrap();
        }
        let keystore = FsKeyStore::open(keystore_path).unwrap();
        keystore.store(&name: String, &private_key_der).unwrap();
        // key_material.as_bytes()

        // let _signatory_key: signatory::pkcs8::PrivateKeyDocument;

        // TODO(ugochi): figure out how to init _signatory_key from _phrase
    }

    #[test]
    fn store_load_and_delete_my_key() {
        let tempdir = tempfile::tempdir().expect("Could not create tempdir");

        let keystore = tempdir.path().join("keystore");
        println!(">> {:?}", keystore);

        let keystore =
            FsKeyStore::create_or_open(&keystore).expect("Could not create or open keystore");

        let key_name = "my_key".parse().expect("Could not parse key name");

        keystore
            .store(&key_name, &secp256k1::SigningKey::generate_pkcs8())
            .expect("Could not store key");

        let key_info = keystore.info(&key_name).expect("Could not load key info");
        println!(">> {:?}", key_info);

        keystore.delete(&key_name).expect("Could not delete key");
    }

    #[test]
    fn explore_hkd32_mnemonic_create_phrase() {
        let lang = hkd32::mnemonic::Language::English;
        let p = hkd32::mnemonic::Phrase::random(&mut rand_core::OsRng, lang);
        println!(">> {}", p.phrase());
    }

    #[test]
    fn explore_hkd32_mnemonic_recover_phrase() {
        let lang = hkd32::mnemonic::Language::English;
        let p = hkd32::mnemonic::Phrase::new(RECOVERY_PHRASE, lang);
        let p = p.expect("Failed to create phrase!");
        assert_eq!(p.phrase(), RECOVERY_PHRASE);

        // Things you can do with p:
        // p.derive_subkey(path: impl AsRef<Path>)
        // p.to_seed(password: &str)
        // p.entropy();
        // p.language()
    }

    #[test]
    fn explore_deep_space_from_phrase() {
        let pass_phrase = "";
        let _deep_space_key = deep_space::PrivateKey::from_phrase(RECOVERY_PHRASE, pass_phrase)
            .expect("Could not generate key from phrase");

        // _deep_space_key.0
    }

    #[test]
    fn explore_xxx() {
        // use signatory::ecdsa::elliptic_curve::AlgorithmParameters;
        // let oid = signatory::ecdsa::Secp256k1::OID;

        let oid = signatory::ecdsa::elliptic_curve::ALGORITHM_OID;
        let oid = signatory::pkcs8::ObjectIdentifier::from(oid);

        let aid = signatory::pkcs8::AlgorithmIdentifier {
            oid,
            parameters: None,
        };

        signatory::pkcs8::PrivateKeyInfo::new(aid, &[]);
    }
}
