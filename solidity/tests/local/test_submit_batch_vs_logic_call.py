#!/usr/bin/python3

from conftest import *

def test_large_batch(signers):
    runSubmitBatchTest(signers, 10)
    runLogicCallTest(signers, 10)

def test_small_batch(signers):
    runSubmitBatchTest(signers, 1)
    runLogicCallTest(signers, 1)

def test_reentrant(signers):
    try:
        runLogicCallTest(signers, 1, True)
    except ValueError as err:
        assert err.args[0] == "ReentrancyGuard: reentrant call"
    else:
        raise "Error"

def prepareTxBatch(batchSize, signers):
    numTxs = batchSize
    destinations = [signers[0].address] * numTxs
    fees = [0] * numTxs
    amounts = [0] * numTxs
    for i in range(numTxs):
        fees[i] = 1
        amounts[i] = 1
        destinations[i] = signers[i + 5].address
    return numTxs, destinations, fees, amounts

def sendToCosmos(gravity, testERC20, numCoins):
    testERC20.approve(gravity, numCoins)
    gravity.sendToCosmos(testERC20, bstring2bytes32(b"myCosmosAddress"), numCoins)

def prep(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)
    reentrantERC20 = ReentrantERC20.deploy(gravity.address, {"from": signers[0]})
    return gravityId, powers, validators, gravity, testERC20, reentrantERC20

def runSubmitBatchTest(signers, batchSize):
    gravityId, powers, validators, gravity, testERC20, reentrantERC20 = prep(signers)
    sendToCosmos(gravity, testERC20, 1000)
    assert testERC20.balanceOf(gravity) == 1000
    assert testERC20.balanceOf(signers[0].address) == 9000

    numTxs, destinations, fees, amounts = prepareTxBatch(batchSize, signers)
    batchNonce = 1
    batchTimeout = 10000
    methodName = bstring2bytes32(b"transactionBatch")
    digest = web3.keccak(encode_abi(
        [
            "bytes32",
            "bytes32",
            "uint256[]",
            "address[]",
            "uint256[]",
            "uint256",
            "address",
            "uint256",
        ],
        [
            gravityId,
            methodName,
            amounts,
            destinations,
            fees,
            batchNonce,
            testERC20.address,
            batchTimeout,
        ]
    ))

    sigs = signHash(validators, digest)
    gravity.submitBatch([getSignerAddresses(validators), powers, 0, 0, "0x0000000000000000000000000000000000000000"], sigs, amounts, destinations, fees, 1, testERC20, batchTimeout)

    assert testERC20.balanceOf(signers[5]) == 1
    assert testERC20.balanceOf(signers[5 + numTxs - 1]) == 1
    assert testERC20.balanceOf(gravity) == 1000 - numTxs * 2
    assert testERC20.balanceOf(signers[0]) == 9000 + numTxs

def runLogicCallTest(signers, batchSize, reentrant=False):
    gravityId, powers, validators, gravity, testERC20, reentrantERC20 = prep(signers)
    tokenBatchMiddleware = TestTokenBatchMiddleware.deploy({"from": signers[0]})
    tokenBatchMiddleware.transferOwnership(gravity)
    sendToCosmos(gravity, testERC20, 1000)
    assert testERC20.balanceOf(gravity) == 1000
    assert testERC20.balanceOf(signers[0]) == 9000
    
    txBatch = prepareTxBatch(batchSize, signers)
    batchNonce = 1
    
    methodName = bstring2bytes32(b"logicCall")
    logicCallArgs = [
        [txBatch[0]],
        [testERC20.address],
        [txBatch[0]],
        [testERC20.address],
        tokenBatchMiddleware.address,
        bytes.fromhex(tokenBatchMiddleware.submitBatch.encode_input(txBatch[3], txBatch[1], reentrantERC20 if reentrant else testERC20)[2:]),
        4766922941000,
        encode_abi(["uint256"], [web3.toInt(hexstr=testERC20.address)]),
        1
    ]
    digest = web3.keccak(
        encode_abi([
            "bytes32", # gravityId
            "bytes32", # methodName
            "uint256[]", # transferAmounts
            "address[]", # transferTokenContracts
            "uint256[]", # feeAmounts
            "address[]", # feeTokenContracts
            "address", # logicContractAddress
            "bytes", # payload
            "uint256", # timeOut
            "bytes32", # invalidationId
            "uint256" # invalidationNonce
        ],
        [
            gravityId,
            methodName,
            logicCallArgs[0],
            logicCallArgs[1],
            logicCallArgs[2],
            logicCallArgs[3],
            logicCallArgs[4],
            logicCallArgs[5],
            logicCallArgs[6],
            logicCallArgs[7],
            logicCallArgs[8]
        ])
    )
    sigs = signHash(validators, digest)

    tx_data = gravity.submitLogicCall.encode_input([getSignerAddresses(validators), powers, 0, 0, "0x0000000000000000000000000000000000000000"], sigs, logicCallArgs)
    try:
        gas = web3.eth.estimate_gas({"to": gravity.address, "from": signers[0].address, "data": tx_data})
    except ValueError as err:
        raise ValueError(err.args[0]["message"][50:])
    except BaseException as err:
        print(f"Unexpected {err=}, {type(err)=}")

    gravity.submitLogicCall(
        [
            getSignerAddresses(validators),
            powers,
            0,
            0,
            "0x0000000000000000000000000000000000000000"
        ],
        sigs,
        logicCallArgs
    )

    assert testERC20.balanceOf(signers[5]) == 1
    assert testERC20.balanceOf(signers[5 + txBatch[0] - 1]) == 1
    assert testERC20.balanceOf(gravity) == 1000 - txBatch[0] * 2
    assert testERC20.balanceOf(signers[0]) == 9000 + txBatch[0]

