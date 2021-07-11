// This library should provide functions that support command line workflows that
// - loads keys by name from a keystore
//   on failure, because the named key does not exist, it offers to create it
//   - offers choices here;
//     - using a new mnemonic
//     - using an existing mnemonic (aka recovery)
// - loads and converts keys to different representations used by different tools
// - leverages the signatory crate to manage the keystore
// - generates keys the hkd32 crate
// - provides abscissa commands
//
// todos for this library:
// [ ] Function that loads a key from keystore by name then converts it to a given format
//     exact list of formats need to be pinned down; generally aiming to play nice with:
//     - clarity
//     - ethers
//     - contact / deep_space

#[cfg(test)]
mod tests {
    use hkd32::mnemonic;
    use rand_core;
    use signatory::ecdsa::secp256k1;
    use signatory::FsKeyStore;
    use signatory::GeneratePkcs8;

    // NOTE these are not _real_ tests. I'm using them to bring together dependencies and learn how to use them.
    // NOTE run w/ `cargo test -- --nocapture` to see the println

    #[test]
    fn load_key_by_name() {
        let tmp = tempfile::tempdir().expect("Failed to create temp dir!");
        let keystore = FsKeyStore::create_or_open(&tmp.path().join("key-store"))
            .expect("Failed to create or open key-store");

        let my_key = secp256k1::SigningKey::generate_pkcs8();

        let key_name = "my_key".parse().expect("Failed to parse key_name");

        keystore
            .store(&key_name, &my_key)
            .expect("Failed to store my_key!");

        let key_info = keystore.info(&key_name).expect("Failed to info my_key!");
        println!(">> {:?}", key_info);
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
