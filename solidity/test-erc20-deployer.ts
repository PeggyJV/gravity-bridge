import { ethers } from "ethers";
import { Gren1 } from "./typechain/Gren1";
import { Gren2 } from "./typechain/Gren2";
import fs from "fs";
import commandLineArgs from "command-line-args";
import { exit } from "process";

const args = commandLineArgs([
  // the ethernum node used to deploy the contract
  { name: "eth-node", type: String },
  // the Ethereum private key that will contain the gas required to pay for the contact deployment
  { name: "eth-privkey", type: String },
]);

async function deploy() {
  var startTime = new Date();
  const provider = await new ethers.providers.JsonRpcProvider(args["eth-node"]);
  let wallet = new ethers.Wallet(args["eth-privkey"], provider);

  const gren1_path =
    "artifacts/contracts/Gren1.sol/Gren1.json";
  const gren2_path =
    "artifacts/contracts/Gren2.sol/Gren2.json";

  const { abi, bytecode } = getContractArtifacts(gren1_path);
  const erc20Factory = new ethers.ContractFactory(abi, bytecode, wallet);
  const testERC20 = (await erc20Factory.deploy()) as Gren1;
  await testERC20.deployed();
  const erc20TestAddress = testERC20.address;
  console.log("Gren1 deployed at Address - ", erc20TestAddress);

  const { abi: abi1, bytecode: bytecode1 } = getContractArtifacts(gren2_path);
  const erc20Factory1 = new ethers.ContractFactory(abi1, bytecode1, wallet);
  const testERC201 = (await erc20Factory1.deploy()) as Gren2;
  await testERC201.deployed();
  const erc20TestAddress1 = testERC201.address;
  console.log("Gren2 deployed at Address - ", erc20TestAddress1);
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
