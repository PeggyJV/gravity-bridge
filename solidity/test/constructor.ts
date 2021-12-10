import chai from "chai";
import { ethers } from "hardhat";
import { solidity } from "ethereum-waffle";
import { deployContracts } from "../test-utils";
import {
  examplePowers
} from "../test-utils/pure";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";

chai.use(solidity);
const { expect } = chai;

describe("constructor tests", function() {
  it("throws on malformed valset", async function() {
    const signers = await ethers.getSigners();
    const gravityId = ethers.utils.formatBytes32String("foo");

    // This is the power distribution on the Cosmos hub as of 7/14/2020
    let powers = examplePowers();
    let validators = signers.slice(0, powers.length - 1);

    const powerThreshold = 6666;

    await expect(
      deployContracts(gravityId, validators, powers,powerThreshold)
    ).to.be.revertedWith("MalformedCurrentValidatorSet()");
  });

  it("throws on insufficient power", async function() {
    const signers = await ethers.getSigners();
    const gravityId = ethers.utils.formatBytes32String("foo");

    // This is the power distribution on the Cosmos hub as of 7/14/2020
    let powers = examplePowers();
    let validators = signers.slice(0, powers.length);

    const powerThreshold = 666666666;

    await expect(
      deployContracts(gravityId, validators, powers, powerThreshold)
    ).to.be.revertedWith(
      "InsufficientPower(10000, 666666666)"
    );
  });

  it("throws on empty validator set", async function () {
    const signers = await ethers.getSigners();
    const gravityId = ethers.utils.formatBytes32String("foo");

    // This is the power distribution on the Cosmos hub as of 7/14/2020
    let powers: number[] = [];
    let validators: SignerWithAddress[] = [];


    await expect(
      deployContracts(gravityId, validators, powers,0)
    ).to.be.revertedWith(
      "InsufficientPower(0, 0)"
    );
  });
});
