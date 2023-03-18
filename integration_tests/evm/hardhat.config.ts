import '@nomiclabs/hardhat-ethers';
import '@nomiclabs/hardhat-waffle';
import {task} from "hardhat/config";
import * as constants from "./addresses";

task(
    'integration_test_setup',
    'Sets up contracts for the integration test',
    async (args, hre) => {

        console.log(`starting with network: ${hre.network.name}, chain id: ${hre.network.config.chainId}`);

        // Take over vitalik.eth
        // await hre.network.provider.request({
        //     method: 'hardhat_impersonateAccount',
        //     params: [constants.WHALE],
        // });

        // Send ETH to needed parties
        // const whaleSigner = await hre.ethers.getSigner(constants.WHALE);

        console.log(`funding validators`)
        const [sender] = await hre.ethers.getSigners();
        for (let addr of constants.VALIDATORS) {
            console.log(`sending 100 to ${addr}`)
            await sender.sendTransaction({
                to: addr,
                value: hre.ethers.utils.parseEther('100'),
            });
        }

        let powers: number[] = [1073741823, 1073741823, 1073741823, 1073741823];
        let powerThreshold: number = 6666;

        console.log(`deploying Gravity contract`)
        const Gravity = await hre.ethers.getContractFactory("Gravity");
        const gravityID = "gravitytest-" + EVM_CHAIN_ID;
        console.log(`gravity ID: ${gravityID}`);
        const gravity = (await Gravity.deploy(
            hre.ethers.utils.formatBytes32String("gravitytest-" + EVM_CHAIN_ID),
            powerThreshold,
            constants.VALIDATORS,
            powers
        ));

        await gravity.deployed();
        console.log(`gravity contract deployed at - ${gravity.address}`)

        console.log(`deploying TestERC20 contract`)
        const TestERC20 = await hre.ethers.getContractFactory("TestERC20GB");
        const testERC20 = (await TestERC20.deploy());
        await testERC20.deployed();
        console.log(`test ERC20 TestGB TGB deployed at - ${testERC20.address}`)

        console.log(`running node`)
        await hre.network.provider.send("evm_setIntervalMining", [1000]);
        await hre.run('node');
    });


/**
 * @type import('hardhat/config').HardhatUserConfig
 */
// const ARCHIVE_NODE_URL = process.env.ARCHIVE_NODE_URL;
const EVM_CHAIN_ID = Number(process.env.EVM_CHAIN_ID);

module.exports = {
    networks: {
        hardhat: {
            chainId: EVM_CHAIN_ID,
        },
    },
    solidity: {
        compilers: [
            {
                version: '0.8.0',
                settings: {
                    optimizer: {
                        enabled: true,
                    },
                },
            },
            {
                version: '0.8.10',
                settings: {
                    optimizer: {
                        enabled: true,
                    },
                },
            },
        ],
    },
    typechain: {
        outDir: 'typechain',
        target: 'ethers-v5',
        runOnCompile: true,
    },
    gasReporter: {
        enabled: true,
    },
};
