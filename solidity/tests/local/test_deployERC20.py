#!/usr/bin/python3

from conftest import *

def test_deployERC20_tests(signers):
    gravityId = bstring2bytes32(b"foo")
    powers = examplePowers()
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)
    tx = gravity.deployERC20("uatom", "Atom", "ATOM", 6)

    assert tx.events["ERC20DeployedEvent"]["_cosmosDenom"] == "uatom"
    assert tx.events["ERC20DeployedEvent"]["_name"] == "Atom"
    assert tx.events["ERC20DeployedEvent"]["_symbol"] == "ATOM"
    assert tx.events["ERC20DeployedEvent"]["_decimals"] == 6
    assert tx.events["ERC20DeployedEvent"]["_eventNonce"] == 2

    ERC20contract = Contract.from_abi("TokenContract", tx.events["ERC20DeployedEvent"]["_tokenContract"], [{"inputs":[{"internalType":"address","name":"_gravityAddress","type":"address"},{"internalType":"string","name":"_name","type":"string"},{"internalType":"string","name":"_symbol","type":"string"},{"internalType":"uint8","name":"_decimals","type":"uint8"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":False,"inputs":[{"indexed":True,"internalType":"address","name":"owner","type":"address"},{"indexed":True,"internalType":"address","name":"spender","type":"address"},{"indexed":False,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":False,"inputs":[{"indexed":True,"internalType":"address","name":"from","type":"address"},{"indexed":True,"internalType":"address","name":"to","type":"address"},{"indexed":False,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}])
    maxUint256 = 2 ** 256 - 1


    assert ERC20contract.balanceOf(gravity) == maxUint256

    numTxs = 100
    txDestinationsInt = [signers[0]] * numTxs
    txFees = [0] * numTxs
    txAmounts = [0] * numTxs
    for i in range(numTxs):
        txFees[i] = 1
        txAmounts[i] = 1
        txDestinationsInt[i] = signers[i + 1]

    txDestinations = getSignerAddresses(txDestinationsInt)
    batchNonce = 1
    batchTimeout = 10000

    methodName = bstring2bytes32(b"transactionBatch")
    abiEncoded = encode_abi(["bytes32", "bytes32", "uint256[]", "address[]", "uint256[]", "uint256", "address", "uint256"], [gravityId, methodName, txAmounts, txDestinations, txFees, batchNonce, tx.events["ERC20DeployedEvent"]["_tokenContract"], batchTimeout])
    digest = web3.keccak(abiEncoded)
    sigs = signHash(validators, digest)
    currentValsetNonce = 0

    valset ={
        "validators": getSignerAddresses(validators),
        "powers": powers,
        "valSetNonce": currentValsetNonce,
        "rewardAmount": 0,
        "rewardToken": "0x0000000000000000000000000000000000000000",
    }
    
    gravity.submitBatch([valset["validators"], valset["powers"], valset["valSetNonce"], valset["rewardAmount"], valset["rewardToken"]], sigs, txAmounts, txDestinations, txFees, batchNonce, tx.events["ERC20DeployedEvent"]["_tokenContract"], batchTimeout)
    assert ERC20contract.balanceOf(gravity) == maxUint256 - 200
    assert ERC20contract.balanceOf(signers[6].address) == 1