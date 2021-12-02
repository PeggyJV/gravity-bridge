#!/usr/bin/python3

from conftest import *

def test_send_to_cosmos_tests(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)
    testERC20.approve(gravity, 1000)
    tx = gravity.sendToCosmos(testERC20, bstring2bytes32(b"myCosmosAddress"), 1000)
    assert tx.events["SendToCosmosEvent"]["_tokenContract"] == testERC20.address
    assert tx.events["SendToCosmosEvent"]["_sender"] == signers[0].address
    assert tx.events["SendToCosmosEvent"]["_destination"].hex() == bstring2bytes32(b"myCosmosAddress").hex()
    assert tx.events["SendToCosmosEvent"]["_amount"] == 1000
    assert tx.events["SendToCosmosEvent"]["_eventNonce"] == 2

    assert testERC20.balanceOf(gravity) == 1000
    assert gravity.state_lastEventNonce() == 2

    testERC20.approve(gravity, 1000)
    tx = gravity.sendToCosmos(testERC20, bstring2bytes32(b"myCosmosAddress"), 1000)
    assert tx.events["SendToCosmosEvent"]["_tokenContract"] == testERC20.address
    assert tx.events["SendToCosmosEvent"]["_sender"] == signers[0].address
    assert tx.events["SendToCosmosEvent"]["_destination"].hex() == bstring2bytes32(b"myCosmosAddress").hex()
    assert tx.events["SendToCosmosEvent"]["_amount"] == 1000
    assert tx.events["SendToCosmosEvent"]["_eventNonce"] == 3

    assert testERC20.balanceOf(gravity) == 2000
    assert gravity.state_lastEventNonce() == 3