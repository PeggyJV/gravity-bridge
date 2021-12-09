#!/usr/bin/python3

from conftest import *

def test_malformed_valset(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:(len(powers) - 1)]
    powerThreshold = 6666

    try:
        deployContracts(signers, gravityId, validators, powers, powerThreshold)
    except ValueError as err:
        assert err.args[0] == "Malformed current validator set"
    else:
        raise "Error"

def test_insufficient_power(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:(len(powers))]
    powerThreshold = 666666666

    try:
        deployContracts(signers, gravityId, validators, powers, powerThreshold)
    except ValueError as err:
        assert err.args[0] == "Submitted validator set signatures do not have enough power."
    else:
        raise "Error"
