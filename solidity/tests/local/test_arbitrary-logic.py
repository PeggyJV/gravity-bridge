#!/usr/bin/python3

from conftest import *

def test_submitLogicCall_malformed_valset_reverts(signers):
    try:
        run_test(signers, malformedCurrentValset=True)
    except ValueError as err:
        assert err.args[0] == "Malformed current validator set"
    else:
        raise "Error"

def test_submitLogicCall_invalidation_nonce_not_incremented_reverts(signers):
    try:
        run_test(signers, invalidationNonceNotHigher=True)
    except ValueError as err:
        assert err.args[0] == "New invalidation nonce must be greater than the current nonce"
    else:
        raise "Error"

def test_submitLogicCall_non_matching_checkpoint_for_current_valset_reverts(signers):
    try:
        run_test(signers, nonMatchingCurrentValset=True)
    except ValueError as err:
        assert err.args[0] == "Supplied current validators and powers do not match checkpoint."
    else:
        raise "Error"

def test_submitLogicCall_bad_validator_sig_reverts(signers):
    try:
        run_test(signers, badValidatorSig=True)
    except ValueError as err:
        assert err.args[0] == "Validator signature does not match."
    else:
        raise "Error"

def test_allows_zeroed_sig(signers):
    run_test(signers, zeroedValidatorSig=True)

def test_not_enough_signatures_reverts(signers):
    try:
        run_test(signers, notEnoughPower=True)
    except ValueError as err:
        assert err.args[0] == "Submitted validator set signatures do not have enough power."
    else:
        raise "Error"

def test_barely_enough_signatures(signers):
    run_test(signers, barelyEnoughPower=True)

def test_timeout_reverts(signers):
    try:
        run_test(signers, timedOut=True)
    except ValueError as err:
        assert err.args[0] == "Timed out"
    else:
        raise "Error"

def test_good_hash(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = [6667]
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)
    testERC20.approve(gravity, 1000)
    gravity.sendToCosmos(testERC20, bstring2bytes32(b"myCosmosAddress"), 1000)

    methodName = bstring2bytes32(b"logicCall")
    numTxs = 10
    invalidationNonce = 1
    timeOut = 4766922941000
    logicCallArgs = [
        [1], # transferAmounts
        [testERC20.address], # transferTokenContracts
        [1], # feeAmounts
        [testERC20.address], # feeTokenContracts
        "0x17c1736CcF692F653c433d7aa2aB45148C016F68", # logicContractAddress
        bstring2bytes32(b"testingPayload"), # payloads
        timeOut,
        bstring2bytes32(b"invalidationId"), # invalidationId
        invalidationNonce # invalidationNonce
    ]

    abiEncodedLogicCall = encode_abi([
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
    logicCallDigest = web3.keccak(abiEncodedLogicCall)
    sigs = signHash(validators, logicCallDigest)
    currentValsetNonce = 0

    res = gravity.submitLogicCall.encode_input(getSignerAddresses(validators), powers, currentValsetNonce, sigs, logicCallArgs)
    print("elements in logic call digest:")
    print({"gravityId":gravityId,
        "logicMethodName": methodName,
        "transferAmounts": logicCallArgs[0],
        "transferTokenContracts": logicCallArgs[1],
        "feeAmounts": logicCallArgs[2],
        "feeTokenContracts": logicCallArgs[3],
        "logicContractAddress": logicCallArgs[4],
        "payload": logicCallArgs[5],
        "timeout": logicCallArgs[6],
        "invalidationId": logicCallArgs[7],
        "invalidationNonce": logicCallArgs[8]}
        )
    print("abitEncodedCall:")
    print(abiEncodedLogicCall)
    print("callDigest:")
    print(logicCallDigest)
    print("elements in logic call function call:")
    print({
        "currentValidators": getSignerAddresses(validators),
        "currentPowers": powers,
        "currentValsetNonce": currentValsetNonce,
        "sigs": [sigs],
    })
    print(res)

def run_test(signers, invalidationNonceNotHigher=False, malformedTxBatch=False, nonMatchingCurrentValset=False, badValidatorSig=False, zeroedValidatorSig=False, notEnoughPower=False, barelyEnoughPower=False, malformedCurrentValset=False, timedOut=False):
    # Prep and deploy contract
    # ========================
    gravityId = bstring2bytes32(b"foo")

    # This is the power distribution on the Cosmos hub as of 7/14/2020
    powers = examplePowers()
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)

    # First we deploy the logic batch middleware contract. This makes it easy to call a logic
    # contract a bunch of times in a batch.
    logicBatch = SimpleLogicBatchMiddleware.deploy({"from": signers[0]})

    # We set the ownership to gravity so that nobody else can call it.
    logicBatch.transferOwnership(gravity, {"from": signers[0]})

    # Then we deploy the actual logic contract.
    logicContract = TestLogicContract.deploy(testERC20, {"from": signers[0]})

    # We set its owner to the batch contract. 
    logicContract.transferOwnership(logicBatch, {"from": signers[0]})

    # Transfer out to Cosmos, locking coins
    # =====================================
    testERC20.approve(gravity, 1000, {"from": signers[0]})
    gravity.sendToCosmos(testERC20, bstring2bytes32(b"myCosmosAddress"), 1000, {"from": signers[0]})

    # Prepare batch
    # ===============================
    # This code prepares the batch of transactions by encoding the arguments to the logicContract.
    # This batch contains 10 transactions which each:
    # - Transfer 5 coins to the logic contract
    # - Call transferTokens on the logic contract, transferring 2+2 coins to signer 20
    # After the batch runs, signer 20 should have 40 coins, Gravity should have 940 coins,
    # and the logic contract should have 10 coins

    numTxs = 10
    txPayloads = []
    txAmounts = []
    for i in range(numTxs):
        txAmounts.append(5)
        txPayloads.append(logicContract.transferTokens.encode_input(signers[20].address, 2, 2))

    invalidationNonce = 1
    if invalidationNonceNotHigher:
        invalidationNonce = 0
    timeOut = 4766922941000
    if timedOut:
        timeOut = 0
    
    # Call method
    # ===========
    # We have to give the logicBatch contract 5 coins for each tx, since it will transfer that
    # much to the logic contract.
    # We give msg.sender 1 coin in fees for each tx.
    methodName = bstring2bytes32(b"logicCall")

    payload = logicBatch.logicBatch.encode_input(txAmounts, txPayloads, logicContract.address, testERC20.address)
    logicCallArgs = [
        [numTxs * 5], # transferAmounts
        [testERC20.address], # transferTokenContracts
        [numTxs], # feeAmounts
        [testERC20.address], # feeTokenContracts
        logicBatch.address, # logicContractAddress
        web3.toBytes(hexstr=payload), # payloads
        timeOut,
        encode_abi(["uint256"], [web3.toInt(hexstr=testERC20.address)]), # invalidationId
        invalidationNonce # invalidationNonce
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

    currentValsetNonce = 0
    if nonMatchingCurrentValset:
        # Wrong nonce
        currentValsetNonce = 420
    
    if malformedCurrentValset:
        # Remove one of the powers to make the length not match
        powers.pop()

    if badValidatorSig:
        # Switch the first sig for the second sig to screw things up
        sigs[1][0] = sigs[0][0]
        sigs[1][1] = sigs[0][1]
        sigs[1][2] = sigs[0][2]
    
    if zeroedValidatorSig:
        # Switch the first sig for the second sig to screw things up
        sigs[1][0] = sigs[0][0]
        sigs[1][0] = sigs[0][0]
        sigs[1][0] = sigs[0][0]
        # Then zero it out to skip evaluation
        sigs[1][0] = 0
    
    if notEnoughPower:
        # zero out enough signatures that we dip below the threshold
        sigs[1][0] = 0
        sigs[2][0] = 0
        sigs[3][0] = 0
        sigs[5][0] = 0
        sigs[6][0] = 0
        sigs[7][0] = 0
        sigs[9][0] = 0
        sigs[11][0] = 0
        sigs[13][0] = 0
    
    if barelyEnoughPower:
        # Stay just above the threshold
        sigs[1][0] = 0
        sigs[2][0] = 0
        sigs[3][0] = 0
        sigs[5][0] = 0
        sigs[6][0] = 0
        sigs[7][0] = 0
        sigs[9][0] = 0
        sigs[11][0] = 0
    
    tx_data = gravity.submitLogicCall.encode_input(getSignerAddresses(validators), powers, currentValsetNonce, sigs, logicCallArgs)
    try:
        gas = web3.eth.estimate_gas({"to": gravity.address, "from": signers[0].address, "data": tx_data})
    except ValueError as err:
        raise ValueError(err.args[0]["message"][50:])
    except BaseException as err:
        print(f"Unexpected {err=}, {type(err)=}")

    gravity.submitLogicCall(getSignerAddresses(validators), powers, currentValsetNonce, sigs, logicCallArgs, {"from": signers[0]})

    assert testERC20.balanceOf(signers[20]) == 40
    assert testERC20.balanceOf(gravity) == 940
    assert testERC20.balanceOf(logicContract) == 10
    assert testERC20.balanceOf(signers[0]) == 9010
