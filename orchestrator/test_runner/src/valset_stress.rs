use crate::happy_path::test_valset_update;
use crate::utils::ValidatorKeys;
use ethers::types::Address as EthAddress;
use gravity::deep_space::Contact;

pub async fn validator_set_stress_test(
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
    gravity_address: EthAddress,
) {
    for _ in 0u32..10 {
        test_valset_update(contact, &keys, gravity_address).await;
    }
}
