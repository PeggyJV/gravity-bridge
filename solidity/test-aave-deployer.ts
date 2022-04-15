import { ethers } from "ethers";
import { MockAaveV2StablecoinCellar } from "./typechain/MockAaveV2StablecoinCellar";
import fs from "fs";
import commandLineArgs from "command-line-args";
import { exit } from "process";

const args = commandLineArgs([
  // the ethernum node used to deploy the contract
  { name: "eth-node", type: String },
  // the Ethereum private key that will contain the gas required to pay for the contact deployment
  { name: "eth-privkey", type: String },
  // the Gravity contract address
  { name: "gravity-address", type: String},
]);

async function deploy() {
  var startTime = new Date();
  const provider = await new ethers.providers.JsonRpcProvider(args["eth-node"]);
  let wallet = new ethers.Wallet(args["eth-privkey"], provider);

  const aave_path =
    "artifacts/contracts/MockAaveV2StablecoinCellar.sol/MockAaveV2StablecoinCellar.json";

  const { abi, bytecode } = getContractArtifacts(aave_path);
  const aaveFactory = new ethers.ContractFactory(abi, bytecode, wallet);
  const aave = (await aaveFactory.deploy()) as MockAaveV2StablecoinCellar;
  await aave.deployed();
  await aave.transferOwnership(args["gravity-address"]);
  const aaveAddress = aave.address;
  console.log("Aave deployed at Address - ", aaveAddress);
}

function getContractArtifacts(path: string): { bytecode: string; abi: string } {
  var { bytecode, abi } = JSON.parse(fs.readFileSync(path, "utf8").toString());
  return { bytecode, abi };
}
const decode = (str: string): string =>
  Buffer.from(str, "base64").toString("binary");

async function main() {
  await deploy();
}

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

main();
