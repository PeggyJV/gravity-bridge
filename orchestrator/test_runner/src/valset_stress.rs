use crate::happy_path::test_valset_update;
use crate::utils::ValidatorKeys;
use deep_space::Contact;
use ethers::types::Address as EthAddress;

pub async fn validator_set_stress_test(
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
    gravity_address: EthAddress,
) {
    for _ in 0u32..10 {
        test_valset_update(contact, &keys, gravity_address).await;
    }
}
