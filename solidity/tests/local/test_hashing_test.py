#!/usr/bin/python3

from conftest import *

def test_hashing_test(signers):
    gravityId = bstring2bytes32(b"foo")
    validators = []
    powers = []
    for i in range(100):
        validators.append(signers[i])
        powers.append(5000)
    
    hashingContract = HashingTest.deploy({"from": signers[0]})
    hashingContract.IterativeHash(getSignerAddresses(validators), powers, 1, gravityId)

    hashingContract.ConcatHash(getSignerAddresses(validators), powers, 1, gravityId)

    hashingContract.ConcatHash2(getSignerAddresses(validators), powers, 1, gravityId)

    contractCheckpoint = hashingContract.lastCheckpoint()
    externalCheckpoint = makeTestCheckpoint(getSignerAddresses(validators), powers, 1, gravityId)
    assert contractCheckpoint == externalCheckpoint.hex()

    hashingContract.JustSaveEverything(getSignerAddresses(validators), powers, 1)
    hashingContract.JustSaveEverythingAgain(getSignerAddresses(validators), powers, 1)

def makeTestCheckpoint(validators, powers, valsetNonce, gravityId):
    methodName = b"checkpoint"
    abiEncoded = encode_abi(["bytes32", "bytes32", "uint256", "address[]", "uint256[]"], [gravityId, methodName, valsetNonce, validators, powers])
    checkpoint = web3.keccak(abiEncoded)
    return checkpoint