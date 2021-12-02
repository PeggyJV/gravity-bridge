#!/usr/bin/python3

from conftest import *

def test_uniswap_logic_happy_path_tests(accounts, signers):

    lp_signer = accounts.at("0x0c731Fb0D03211DD32A456370AD2ec3fFad46520")
    print(lp_signer)

    starting_lp_eth_balance = lp_signer.balance()
    print("starting LP eth balance: " + str(starting_lp_eth_balance))

    uniswap_router_address = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"

    usdc_eth_lp = Contract.from_abi("UniswapV2Pair", "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc", [{"inputs":[],"payable":False,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":False,"inputs":[{"indexed":True,"internalType":"address","name":"owner","type":"address"},{"indexed":True,"internalType":"address","name":"spender","type":"address"},{"indexed":False,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":False,"inputs":[{"indexed":True,"internalType":"address","name":"sender","type":"address"},{"indexed":False,"internalType":"uint256","name":"amount0","type":"uint256"},{"indexed":False,"internalType":"uint256","name":"amount1","type":"uint256"},{"indexed":True,"internalType":"address","name":"to","type":"address"}],"name":"Burn","type":"event"},{"anonymous":False,"inputs":[{"indexed":True,"internalType":"address","name":"sender","type":"address"},{"indexed":False,"internalType":"uint256","name":"amount0","type":"uint256"},{"indexed":False,"internalType":"uint256","name":"amount1","type":"uint256"}],"name":"Mint","type":"event"},{"anonymous":False,"inputs":[{"indexed":True,"internalType":"address","name":"sender","type":"address"},{"indexed":False,"internalType":"uint256","name":"amount0In","type":"uint256"},{"indexed":False,"internalType":"uint256","name":"amount1In","type":"uint256"},{"indexed":False,"internalType":"uint256","name":"amount0Out","type":"uint256"},{"indexed":False,"internalType":"uint256","name":"amount1Out","type":"uint256"},{"indexed":True,"internalType":"address","name":"to","type":"address"}],"name":"Swap","type":"event"},{"anonymous":False,"inputs":[{"indexed":False,"internalType":"uint112","name":"reserve0","type":"uint112"},{"indexed":False,"internalType":"uint112","name":"reserve1","type":"uint112"}],"name":"Sync","type":"event"},{"anonymous":False,"inputs":[{"indexed":True,"internalType":"address","name":"from","type":"address"},{"indexed":True,"internalType":"address","name":"to","type":"address"},{"indexed":False,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"constant":True,"inputs":[],"name":"DOMAIN_SEPARATOR","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":True,"inputs":[],"name":"MINIMUM_LIQUIDITY","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":True,"inputs":[],"name":"PERMIT_TYPEHASH","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":True,"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":False,"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":False,"stateMutability":"nonpayable","type":"function"},{"constant":True,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":False,"inputs":[{"internalType":"address","name":"to","type":"address"}],"name":"burn","outputs":[{"internalType":"uint256","name":"amount0","type":"uint256"},{"internalType":"uint256","name":"amount1","type":"uint256"}],"payable":False,"stateMutability":"nonpayable","type":"function"},{"constant":True,"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":True,"inputs":[],"name":"factory","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":True,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"_reserve0","type":"uint112"},{"internalType":"uint112","name":"_reserve1","type":"uint112"},{"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":False,"inputs":[{"internalType":"address","name":"_token0","type":"address"},{"internalType":"address","name":"_token1","type":"address"}],"name":"initialize","outputs":[],"payable":False,"stateMutability":"nonpayable","type":"function"},{"constant":True,"inputs":[],"name":"kLast","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":False,"inputs":[{"internalType":"address","name":"to","type":"address"}],"name":"mint","outputs":[{"internalType":"uint256","name":"liquidity","type":"uint256"}],"payable":False,"stateMutability":"nonpayable","type":"function"},{"constant":True,"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":True,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"nonces","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":False,"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"uint256","name":"deadline","type":"uint256"},{"internalType":"uint8","name":"v","type":"uint8"},{"internalType":"bytes32","name":"r","type":"bytes32"},{"internalType":"bytes32","name":"s","type":"bytes32"}],"name":"permit","outputs":[],"payable":False,"stateMutability":"nonpayable","type":"function"},{"constant":True,"inputs":[],"name":"price0CumulativeLast","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":True,"inputs":[],"name":"price1CumulativeLast","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":False,"inputs":[{"internalType":"address","name":"to","type":"address"}],"name":"skim","outputs":[],"payable":False,"stateMutability":"nonpayable","type":"function"},{"constant":False,"inputs":[{"internalType":"uint256","name":"amount0Out","type":"uint256"},{"internalType":"uint256","name":"amount1Out","type":"uint256"},{"internalType":"address","name":"to","type":"address"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"swap","outputs":[],"payable":False,"stateMutability":"nonpayable","type":"function"},{"constant":True,"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":False,"inputs":[],"name":"sync","outputs":[],"payable":False,"stateMutability":"nonpayable","type":"function"},{"constant":True,"inputs":[],"name":"token0","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":True,"inputs":[],"name":"token1","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":True,"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":False,"stateMutability":"view","type":"function"},{"constant":False,"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":False,"stateMutability":"nonpayable","type":"function"},{"constant":False,"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":False,"stateMutability":"nonpayable","type":"function"}])

    total_lp_supply = usdc_eth_lp.totalSupply()
    lp_provider_balance = usdc_eth_lp.balanceOf(lp_signer)
    reserve = usdc_eth_lp.getReserves()
    reserve0 = reserve[0]
    reserve1 = reserve[1]

    print("TotalSupply:" + str(total_lp_supply) + " Whale Balance:" + str(lp_provider_balance) + " Reserve0:" + str(reserve0) + " Reserve1:" + str(reserve1))

    lp_balance_to_send = 2_000_000_000_000
    eth_per_lp_unit = reserve1 / total_lp_supply

    print("Eth per lp token:" + str(eth_per_lp_unit))

    usdc_address = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
    
    gravityId = bstring2bytes32(b"foo")

    powers = examplePowers()
    validators = signers[:len(powers)]
    powerThreshold = 6666
    gravity, testERC20, checkpoint = deployContracts(signers, gravityId, validators, powers, powerThreshold)

    logicBatch = SimpleLogicBatchMiddleware.deploy({"from": signers[0]})

    logicBatch.transferOwnership(gravity)

    logicContract = TestUniswapLiquidity.deploy(uniswap_router_address, {"from": signers[0]})

    logicContract.transferOwnership(logicBatch)

    logic_contract_balance_start = usdc_eth_lp.balanceOf(logicContract)

    print("Logic Contract Balance:" + str(logic_contract_balance_start))

    usdc_eth_lp.approve(gravity, lp_provider_balance, {"from": lp_signer})

    gravity.sendToCosmos(usdc_eth_lp, bstring2bytes32(b"myCosmosAddress"), lp_balance_to_send * 500, {"from": lp_signer})

    post_gas_balance = lp_signer.balance()

    print("Post_gas_balance:" + str(post_gas_balance))

    numTxs = 3
    txPayloads = []
    txAmounts = []
    for i in range(numTxs):
        txAmounts.append(lp_balance_to_send)
        txPayloads.append(logicContract.redeemLiquidityETH.encode_input(usdc_address, lp_balance_to_send, 0, 0, lp_signer, 4766922941000))

    invalidationNonce = 1
    timeOut = 4766922941000
    
    methodName = bstring2bytes32(b"logicCall")

    payload = logicBatch.logicBatch.encode_input(txAmounts, txPayloads, logicContract.address, usdc_eth_lp)
    logicCallArgs = [
        [lp_balance_to_send * 400], # transferAmounts
        [usdc_eth_lp.address], # transferTokenContracts
        [numTxs], # feeAmounts
        [usdc_eth_lp.address], # feeTokenContracts
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

    tx_data = gravity.submitLogicCall.encode_input(getSignerAddresses(validators), powers, currentValsetNonce, sigs, logicCallArgs)
    try:
        gas = web3.eth.estimate_gas({"to": gravity.address, "from": signers[0].address, "data": tx_data})
    except ValueError as err:
        raise ValueError(err.args[0]["message"][50:])

    gravity.submitLogicCall(getSignerAddresses(validators), powers, currentValsetNonce, sigs, logicCallArgs, {"from": signers[0]})

    ending_lp_eth_balance = lp_signer.balance()

    print("Ending LP eth balance:" + str(ending_lp_eth_balance))

    balance_difference = ending_lp_eth_balance - post_gas_balance

    print("Ending LP eth balance difference:" + str(balance_difference))

    expect_gains = eth_per_lp_unit * lp_balance_to_send * numTxs

    print("Expected LP eth balance difference:" + str(expect_gains))

    logic_contract_balance_end = usdc_eth_lp.balanceOf(logicContract)

    print("Logic Contract Balance:" + str(logic_contract_balance_end))

    assert logic_contract_balance_end == 0
    assert balance_difference - expect_gains > 0
