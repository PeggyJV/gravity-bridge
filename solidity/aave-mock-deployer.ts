import { MockAaveV2StablecoinCellar } from "./typechain/MockAaveV2StablecoinCellar";
import { ethers } from "ethers";
import fs from "fs";
import commandLineArgs from "command-line-args";
import { assert } from "console";
import { doesNotThrow } from "assert";

const args = commandLineArgs([
    // the ethernum node used to deploy the contract
    { name: "eth-node", type: String },
    // the Ethereum private key that will contain the gas required to pay for the contact deployment
    { name: "eth-privkey", type: String },
    // the cellar contract .json file
    { name: "contract", type: String },
    // gravity contract address to transfer ownership to
    { name: "gravity-address", type: String }
]);

async function deploy() {
    // validate contract address
    const gravityAddress = args["gravity-address"];
    doesNotThrow(() => ethers.utils.getAddress(gravityAddress), "the provided gravity address is invalid.");

    const provider = await new ethers.providers.JsonRpcProvider(args["eth-node"]);
    const wallet = new ethers.Wallet(args["eth-privkey"], provider);
    const { abi, bytecode } = getContractArtifacts(args["contract"]);
    const cellarFactory = new ethers.ContractFactory(abi, bytecode, wallet);
    const cellar = (await cellarFactory.deploy()) as MockAaveV2StablecoinCellar;

    await cellar.deployed();
    await cellar.transferOwnership(args["gravity-address"]);
}

function getContractArtifacts(path: string): { bytecode: string; abi: string } {
    var { bytecode, abi } = JSON.parse(fs.readFileSync(path, "utf8").toString());
    return { bytecode, abi };
}

deploy();
