#!/usr/bin/python3

from conftest import *

def test_throws_on_malformed_new_valset(signers):
    try:
        run_test(signers, malformedNewValset=True)
    except ValueError as err:
        assert err.args[0] == "MalformedNewValidatorSet()"
    else:
        raise "Error"

def test_throws_on_malformed_current_valset(signers):
    try:
        run_test(signers, malformedCurrentValset=True)
    except ValueError as err:
        assert err.args[0] == "MalformedNewValidatorSet()"
    else:
        raise "Error"

def test_throws_on_non_matching_checkpoint_for_current_valset(signers):
    try:
        run_test(signers, nonMatchingCurrentValset=True)
    except ValueError as err:
        assert err.args[0] == "IncorrectCheckpoint()"
    else:
        raise "Error"

def test_throws_on_new_valset_nonce_not_incremented(signers):
    try:
        run_test(signers, nonceNotIncremented=True)
    except ValueError as err:
        assert err.args[0] == "InvalidValsetNonce(0, 0)"
    else:
        raise "Error"

def test_throws_on_bad_validator_sig(signers):
    try:
        run_test(signers, badValidatorSig=True)
    except ValueError as err:
        assert err.args[0] == "InvalidSignature()"
    else:
        raise "Error"

def test_allows_zeroed_sig(signers):
    run_test(signers, zeroedValidatorSig=True)

def test_throws_on_not_enough_signatures(signers):
    try:
        run_test(signers, notEnoughPower=True)
    except ValueError as err:
        assert err.args[0] == "InsufficientPower(2807621889, 2863311530)"
    else:
        raise "Error"

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

    
    tx_data = gravity.updateValset.encode_input(getSignerAddresses(newValidators), newPowers, newValsetNonce, getSignerAddresses(validators), powers, currentValsetNonce, sigs)
    try:
        gas = web3.eth.estimate_gas({"to": gravity.address, "from": signers[0].address, "data": tx_data})
    except ValueError as err:
        raise ValueError(err.args[0]["message"][50:])
    except BaseException as err:
        print(f"Unexpected {err=}, {type(err)=}")

    gravity.updateValset(getSignerAddresses(newValidators), newPowers, newValsetNonce, getSignerAddresses(validators), powers, currentValsetNonce, sigs)
    return gravity, checkpoint