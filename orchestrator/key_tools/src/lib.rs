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
    use hkd32::mnemonic;
    use rand_core;
    use signatory::ecdsa::secp256k1;
    use signatory::FsKeyStore;
    use signatory::GeneratePkcs8;

    // NOTE these are not _real_ tests. I'm using them to bring together dependencies and learn how to use them.
    // NOTE run w/ `cargo test -- --nocapture` to see the println

    // Questions:
    // - How do I get from a mnemonic::Phrase to a pkcs8::PrivateKeyDocument (for keystore.store)?
    //   -- maybe derive_subkey?? or possibly to_seed?? (latter requires a password)
    // - What "display formats" do we care about?
    //   -- and where can I find examples??

    #[test]
    fn convert_signatory_to_clarity_pubkey() {
        let signatory_key = &secp256k1::SigningKey::generate_pkcs8();
        println!("{}", signatory_key.to_pem().as_str());


        println!(">> algorithm {:?}", signatory_key.private_key_info().algorithm);

        // let clarity_key = clarity::private_key::PrivateKey::from_slice(
        // signatory_key.private_key_info().private_key,
        // )
        // .expect("Could not create clarity key");
        //
        // let clarity_pubkey = clarity_key
        // .to_public_key()
        // .expect("Could not create clarity pub key");
        //
        // println!(">> {}", clarity_pubkey);
    }

    #[test]
    fn convert_phrase_to_private_key_document() {
        let phrase = mnemonic::Phrase::random(&mut rand_core::OsRng, mnemonic::Language::English);

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

        let key_info = keystore.info(&key_name).expect("Could not lookup key");
        println!(">> {:?}", key_info);

        keystore.delete(&key_name).expect("Could not delete key");
    }

    #[test]
    fn create_random_mnemonic_phrase() {
        let p = mnemonic::Phrase::random(&mut rand_core::OsRng, mnemonic::Language::English);
        println!(">> {}", p.phrase());
    }

    #[test]
    fn recover_phrase_from_mnemonic() {
        const MNEMONIC: &str = "save able shop proud seek reflect prepare mechanic armor car core shuffle room axis file diet axis try secret evolve opinion prosper flush buyer";
        let p = mnemonic::Phrase::new(MNEMONIC, mnemonic::Language::English)
            .expect("Failed to create phrase!");
        assert_eq!(MNEMONIC, p.phrase())
    }
}
