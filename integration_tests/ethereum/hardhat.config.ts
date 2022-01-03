import '@nomiclabs/hardhat-ethers';
import '@nomiclabs/hardhat-waffle';
import {task} from "hardhat/config";
import * as constants from "./addresses";

task(
    'integration_test_setup',
    'Sets up contracts for the integration test',
    async (args, hre) => {

        // Take over vitalik.eth
        await hre.network.provider.request({
            method: 'hardhat_impersonateAccount',
            params: [constants.WHALE],
        });

        // Send ETH to needed parties
        const whaleSigner = await hre.ethers.getSigner(constants.WHALE);

        for (let addr of constants.VALIDATORS) {
            await whaleSigner.sendTransaction({
                to: addr,
                value: hre.ethers.utils.parseEther('100'),
            });
        }

        let powers: number[] = [1073741823, 1073741823, 1073741823, 1073741823];
        let powerThreshold: number = 6666;

        const Gravity = await hre.ethers.getContractFactory("Gravity");
        const gravity = (await Gravity.deploy(
            hre.ethers.utils.formatBytes32String("gravitytest"),
            powerThreshold,
            constants.VALIDATORS,
            powers
        ));

        await gravity.deployed();
        console.log(`gravity contract deployed at - ${gravity.address}`)

        const TestERC20 = await hre.ethers.getContractFactory("TestERC20GB");
        const testERC20 = (await TestERC20.deploy());
        await testERC20.deployed();
        console.log(`test ERC20 TestGB TGB deployed at - ${testERC20.address}`)

        await hre.network.provider.send("evm_setIntervalMining", [1000]);

        await hre.run('node');
    });


/**
 * @type import('hardhat/config').HardhatUserConfig
 */
const ARCHIVE_NODE_URL = process.env.ARCHIVE_NODE_URL;

module.exports = {
    networks: {
        hardhat: {
            forking: {
                url: ARCHIVE_NODE_URL,
                blockNumber: 13405367,
            },
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
