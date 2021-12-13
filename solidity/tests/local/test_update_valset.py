#!/usr/bin/python3

from conftest import *

def test_throws_on_malformed_new_valset(signers):
    method_id = web3.keccak(b"MalformedNewValidatorSet()")[:4].hex()
    revert_string = "typed error: " + str(method_id)
    with brownie.reverts(revert_string):
        run_test(signers, malformedNewValset=True)

def test_throws_on_malformed_current_valset(signers):
    method_id = web3.keccak(b"MalformedCurrentValidatorSet()")[:4].hex()
    revert_string = "typed error: " + str(method_id)
    with brownie.reverts(revert_string):
        run_test(signers, malformedCurrentValset=True)

def test_throws_on_non_matching_checkpoint_for_current_valset(signers):
    method_id = web3.keccak(b"IncorrectCheckpoint()")[:4].hex()
    revert_string = "typed error: " + str(method_id)
    with brownie.reverts(revert_string):
        run_test(signers, nonMatchingCurrentValset=True)

def test_throws_on_new_valset_nonce_not_incremented(signers):
    method_id = web3.keccak(b"InvalidValsetNonce(uint256,uint256)")[:4].hex()
    encode_data = encode_abi(['uint256', 'uint256'], [0, 0]).hex()
    revert_string = "typed error: " + str(method_id + encode_data)
    with brownie.reverts(revert_string):
        run_test(signers, nonceNotIncremented=True)

# # confirmed correct revert but ganache-cli issue
# def test_throws_on_bad_validator_sig(signers):
#     method_id = web3.keccak(b"InvalidSignature()")[:4].hex()
#     revert_string = "typed error: " + str(method_id)
#     with brownie.reverts(revert_string):
#         run_test(signers, badValidatorSig=True)

def test_allows_zeroed_sig(signers):
    run_test(signers, zeroedValidatorSig=True)

# # confirmed correct revert but ganache-cli issue
# def test_throws_on_not_enough_signatures(signers):
#     method_id = web3.keccak(b"InsufficientPower(uint256,uint256)")[:4].hex()
#     encode_data = encode_abi(['uint256', 'uint256'], [2807621889, 2863311530]).hex()
#     revert_string = "typed error: " + str(method_id + encode_data)
#     with brownie.reverts(revert_string):
#         run_test(signers, notEnoughPower=True)

def test_happy_path(signers):
    gravity, checkpoint = run_test(signers)
    assert gravity.state_lastValsetCheckpoint() == checkpoint.hex()

def run_test(signers, malformedNewValset=False, malformedCurrentValset=False, nonMatchingCurrentValset=False, nonceNotIncremented=False, badValidatorSig=False, zeroedValidatorSig=False, notEnoughPower=False):
    # Prep and deploy contract
    # ========================
    gravityId = bstring2bytes32(b"foo")

    # This is the power distribution on the Cosmos hub as of 7/14/2020
    powers = examplePowers()
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)

    newPowers = examplePowers()
    newPowers[0] -= 3
    newPowers[1] += 3

    newValidators = signers[:len(newPowers)]

    if malformedNewValset:
        newValidators = signers[:len(newPowers) - 1]
    
    currentValsetNonce = 0
    if nonMatchingCurrentValset:
        powers[0] = 78
    
    newValsetNonce = 1
    if nonceNotIncremented:
        newValsetNonce = 0

    checkpoint = makeCheckpoint(getSignerAddresses(newValidators), newPowers, newValsetNonce, gravityId)

    sigs = signHash(validators, checkpoint)

    if badValidatorSig:
        sigs[1][0] = sigs[0][0]
        sigs[1][1] = sigs[0][1]
        sigs[1][2] = sigs[0][2]
    
    if zeroedValidatorSig:
        sigs[1][0] = sigs[0][0]
        sigs[1][1] = sigs[0][1]
        sigs[1][2] = sigs[0][2]
        sigs[1][0] = 0
    
    if notEnoughPower:
        sigs[1][0] = 0
        sigs[2][0] = 0
        sigs[3][0] = 0
        sigs[5][0] = 0
        sigs[6][0] = 0
        sigs[7][0] = 0
        sigs[9][0] = 0
        sigs[11][0] = 0
        sigs[13][0] = 0

    if malformedCurrentValset:
        powers.pop()

    
    gravity.updateValset([getSignerAddresses(newValidators), newPowers, newValsetNonce, 0, "0x0000000000000000000000000000000000000000"], [getSignerAddresses(validators), powers, currentValsetNonce, 0, "0x0000000000000000000000000000000000000000"], sigs)
    return gravity, checkpoint