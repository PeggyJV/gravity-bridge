#!/usr/bin/python3

from conftest import *

def test_throws_on_malformed_current_valset(signers):
    try:
        run_test(signers, malformedCurrentValset=True)
    except ValueError as err:
        assert err.args[0] == "Malformed current validator set"
    else:
        raise "Error"

def test_throws_on_malformed_txbatch(signers):
    try:
        run_test(signers, malformedTxBatch=True)
    except ValueError as err:
        assert err.args[0] == "Malformed batch of transactions"
    else:
        raise "Error"

def test_throws_on_batch_nonce_incremented(signers):
    try:
        run_test(signers, batchNonceNotHigher=True)
    except ValueError as err:
        assert err.args[0] == "New batch nonce must be greater than the current nonce"
    else:
        raise "Error"

def test_throws_on_timeout_batch(signers):
    try:
        run_test(signers, batchTimedOut=True)
    except ValueError as err:
        assert err.args[0] == "Batch timeout must be greater than the current block height"
    else:
        raise "Error"

def test_throws_on_non_matching_checkpoint_for_current_valset(signers):
    try:
        run_test(signers, nonMatchingCurrentValset=True)
    except ValueError as err:
        assert err.args[0] == "Supplied current validators and powers do not match checkpoint."
    else:
        raise "Error"

def test_throws_on_bad_validator_sig(signers):
    try:
        run_test(signers, badValidatorSig=True)
    except ValueError as err:
        assert err.args[0] == "Validator signature does not match."
    else:
        raise "Error"

def test_allows_zeroed_sig(signers):
    run_test(signers, zeroedValidatorSig=True)

def test_throws_on_not_enough_signatures(signers):
    try:
        run_test(signers, notEnoughPower=True)
    except ValueError as err:
        assert err.args[0] == "Submitted validator set signatures do not have enough power."
    else:
        raise "Error"

def test_does_not_throw_on_barely_enough_signatures(signers):
    run_test(signers, barelyEnoughPower=True)

def test_produces_good_hash(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = [6667]
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)

    txAmounts = [1]
    txFees = [1]
    txDestinations = getSignerAddresses([signers[5]])
    batchNonce = 1
    batchTimeout = web3.eth.get_block('latest').number + 1000

    testERC20.approve(gravity, 1000)
    gravity.sendToCosmos(testERC20, bstring2bytes32(b"myCosmosAddress"), 1000)

    batchMethodName = bstring2bytes32(b"transactionBatch")

    abiEncodedBatch = encode_abi([
            "bytes32",
            "bytes32",
            "uint256[]",
            "address[]",
            "uint256[]",
            "uint256",
            "address",
            "uint256"
        ],
        [
            gravityId,
            batchMethodName,
            txAmounts,
            txDestinations,
            txFees,
            batchNonce,
            testERC20.address,
            batchTimeout
        ])
    batchDigest = web3.keccak(abiEncodedBatch)

    print("elements in batch digest:")
    print({"gravityId":gravityId,
        "logicMethodName": batchMethodName,
        "txAmounts": txAmounts,
        "txDestinations": txDestinations,
        "txFees": txFees,
        "batchNonce": batchNonce,
        "batchTimeout": batchTimeout,
        "tokenContract": testERC20.address})
    print("abiEncodedBatch:")
    print(abiEncodedBatch)
    print("batchDigest:")
    print(batchDigest)

    sigs = signHash(validators, batchDigest)
    currentValsetNonce = 0

    res = gravity.submitBatch(getSignerAddresses(validators), powers, currentValsetNonce, sigs, txAmounts, txDestinations, txFees, batchNonce, testERC20, batchTimeout)

def run_test(signers, batchNonceNotHigher=False, malformedTxBatch=False, nonMatchingCurrentValset=False, badValidatorSig=False, zeroedValidatorSig=False, notEnoughPower=False, barelyEnoughPower=False, malformedCurrentValset=False, batchTimedOut=False):
    # Prep and deploy contract
    # ========================
    gravityId = bstring2bytes32(b"foo")

    # This is the power distribution on the Cosmos hub as of 7/14/2020
    powers = examplePowers()
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)

    # Transfer out to Cosmos, locking coins
    # =====================================
    testERC20.approve(gravity, 1000, {"from": signers[0]})
    gravity.sendToCosmos(testERC20, bstring2bytes32(b"myCosmosAddress"), 1000, {"from": signers[0]})

    # Prepare batch
    # ===============================

    numTxs = 100
    txDestinationsInt = []
    txFees = []

    txAmounts = []
    for i in range(numTxs):
        txFees.append(1)
        txAmounts.append(1)
        txDestinationsInt.append(signers[i + 5])
    txDestinations = getSignerAddresses(txDestinationsInt)
    
    if malformedTxBatch:
        txFees.pop()
    
    block = web3.eth.get_block('latest')
    batchTimeout = block.number + 1000
    if batchTimedOut:
        batchTimeout = block.number - 1
    
    batchNonce = 1
    if batchNonceNotHigher:
        batchNonce = 0
    
    methodName = bstring2bytes32(b"transactionBatch")

    digest = web3.keccak(
        encode_abi([
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
            txAmounts,
            txDestinations,
            txFees,
            batchNonce,
            testERC20.address,
            batchTimeout
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

    tx_data = gravity.submitBatch.encode_input(getSignerAddresses(validators), powers, currentValsetNonce, sigs, txAmounts, txDestinations, txFees, batchNonce, testERC20, batchTimeout)
    try:
        web3.eth.estimate_gas({"to": gravity.address, "from": signers[0].address, "data": tx_data})
    except ValueError as err:
        raise ValueError(err.args[0]["message"][50:])
    except BaseException as err:
        print(f"Unexpected {err=}, {type(err)=}")

    gravity.submitBatch(getSignerAddresses(validators), powers, currentValsetNonce, sigs, txAmounts, txDestinations, txFees, batchNonce, testERC20, batchTimeout)
