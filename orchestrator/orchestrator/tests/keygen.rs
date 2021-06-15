use deep_space::private_key::PrivateKey as CosmosPrivateKey;

#[test]
fn check_expected_keygen() {
    let mnemonic = "weasel lunch attack blossom tone drum unfair worry risk level negative height sight nation inside task oyster client shiver aware neck mansion gun dune";
    let expected_addr = "cosmos18umn8nad5m8vcr3567v0ylu0m3q2ksrp5c3zf5";

    let cosmos_key = CosmosPrivateKey::from_phrase(mnemonic, "").expect("Invalid Private Cosmos Key!");
    let public_cosmos_address =  cosmos_key.to_address("cosmos").unwrap();
    assert_eq!(public_cosmos_address.to_string(), expected_addr);
}