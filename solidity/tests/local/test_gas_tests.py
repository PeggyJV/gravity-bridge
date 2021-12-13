#!/usr/bin/python3

from conftest import *

def test_make_checkpoint_in_isolation(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)
    gravity.testMakeCheckpoint([getSignerAddresses(validators), powers, 0, 0, "0x0000000000000000000000000000000000000000"], gravityId)

def test_check_validator_signatures_in_isolation(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)
    sigs = signHash(validators, 0x7bc422a00c175cae98cf2f4c36f2f8b63ec51ab8c57fecda9bccf0987ae2d67d)
    gravity.testCheckValidatorSignatures([getSignerAddresses(validators), powers, 0, 0, "0x0000000000000000000000000000000000000000"], sigs, 0x7bc422a00c175cae98cf2f4c36f2f8b63ec51ab8c57fecda9bccf0987ae2d67d, 6666)
