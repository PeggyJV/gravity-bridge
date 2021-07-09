
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
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}