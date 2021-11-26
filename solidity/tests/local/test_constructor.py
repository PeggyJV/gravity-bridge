#!/usr/bin/python3

from conftest import *

def test_malformed_valset(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:(len(powers) - 1)]
    powerThreshold = 6666

    with brownie.reverts("Malformed current validator set"):
        deployContracts(signers, gravityId, validators, powers, powerThreshold)

def test_insufficient_power(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:(len(powers))]
    powerThreshold = 666666666

    with brownie.reverts("Submitted validator set signatures do not have enough power."):
        deployContracts(signers, gravityId, validators, powers, powerThreshold)
