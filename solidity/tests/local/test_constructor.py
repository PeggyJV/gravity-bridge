#!/usr/bin/python3

from conftest import *

def test_malformed_valset(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:(len(powers) - 1)]
    powerThreshold = 6666
    method_id = web3.keccak(b"MalformedCurrentValidatorSet()")[:4].hex()
    revert_string = "typed error: " + str(method_id)

    with brownie.reverts(revert_string):
        deployContracts(signers, gravityId, validators, powers, powerThreshold)

def test_insufficient_power(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:(len(powers))]
    powerThreshold = 666666666
    method_id = web3.keccak(b"InsufficientPower(uint256,uint256)")[:4].hex()
    encode_data = encode_abi(['uint256', 'uint256'], [10000, 666666666]).hex()
    revert_string = "typed error: " + str(method_id + encode_data)

    with brownie.reverts(revert_string):
        deployContracts(signers, gravityId, validators, powers, powerThreshold)

