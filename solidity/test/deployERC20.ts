import chai from "chai";
import { ethers } from "hardhat";
import { solidity } from "ethereum-waffle";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  makeTxBatchHash,
  examplePowers,
  ZeroAddress
} from "../test-utils/pure";
import {ContractReceipt, utils} from 'ethers';
import { BigNumber } from "ethers";
chai.use(solidity);
const { expect } = chai;

async function parseEvent(contract: any, txReceipt: Promise<ContractReceipt>, eventOrder: number) {
  const receipt = await txReceipt;

  if (receipt.events){
    let args = receipt.events[eventOrder].args;

    return args

  }

  return undefined
}

async function runTest(opts: {}) {



  // Prep and deploy Gravity contract
  // ========================
  const signers = await ethers.getSigners();
  const gravityId = ethers.utils.formatBytes32String("foo");
  // This is the power distribution on the Cosmos hub as of 7/14/2020
  let powers = examplePowers();
  let validators = signers.slice(0, powers.length);
  const powerThreshold = 6666;
  const {
    gravity,
    testERC20,
    checkpoint: deployCheckpoint
  } = await deployContracts(gravityId, validators, powers, powerThreshold);




  // Deploy ERC20 contract representing Cosmos asset
  // ===============================================

  let tx = await gravity.deployERC20('uatom', 'Atom', 'ATOM', 6);



  const eventArgs = await parseEvent(gravity,tx.wait(), 1)

  if (eventArgs == undefined)
   {
    throw new Error("No event args");
  }


  expect(eventArgs._cosmosDenom).to.equal( 'uatom');




  // Connect to deployed contract for testing
  // ========================================
  let ERC20contract = new ethers.Contract(eventArgs._tokenContract, [
    "function balanceOf(address account) view returns (uint256 balance)"
  ], gravity.provider);


  const maxUint256 = BigNumber.from(2).pow(256).sub(1)

  // Check that gravity balance is correct
  expect((await ERC20contract.functions.balanceOf(gravity.address)).toString()).to.equal(maxUint256.toString());


  // Prepare batch
  // ===============================
  const numTxs = 100;
  const txDestinationsInt = new Array(numTxs);
  const txFees = new Array(numTxs);

  const txAmounts = new Array(numTxs);
  for (let i = 0; i < numTxs; i++) {
    txFees[i] = 1;
    txAmounts[i] = 1;
    txDestinationsInt[i] = signers[i + 5];
  }
  const txDestinations = await getSignerAddresses(txDestinationsInt);
  let batchNonce = 1
  let batchTimeout = 10000




  // Call method
  // ===========
  const methodName = ethers.utils.formatBytes32String(
    "transactionBatch"
  );
  let abiEncoded = ethers.utils.defaultAbiCoder.encode(
    [
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
      methodName,
      txAmounts,
      txDestinations,
      txFees,
      batchNonce,
      eventArgs._tokenContract,
      batchTimeout
    ]
  );
  let digest = ethers.utils.keccak256(abiEncoded);
  let sigs = await signHash(validators, digest);
  let currentValsetNonce = 0;

  let valset = {
    validators: await getSignerAddresses(validators),
    powers,
    valsetNonce: currentValsetNonce,
    rewardAmount: 0,
    rewardToken: ZeroAddress
  }

  await gravity.submitBatch(
    valset,

    sigs,

    txAmounts,
    txDestinations,
    txFees,
    batchNonce,
    eventArgs._tokenContract,
    batchTimeout
  );

  // Check that Gravity's balance is correct
  expect((await ERC20contract.functions.balanceOf(gravity.address)).toString()).to.equal(maxUint256.sub(200).toString())

  // Check that one of the recipient's balance is correct
  expect((await ERC20contract.functions.balanceOf(await signers[6].getAddress())).toString()).to.equal('1')
}

describe("deployERC20 tests", function () {
  // There is no way for this function to throw so there are
  // no throwing tests
  it("runs", async function () {
    await runTest({})
  });
});
