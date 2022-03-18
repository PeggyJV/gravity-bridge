"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
require("@nomiclabs/hardhat-ethers");
require("@nomiclabs/hardhat-waffle");
var config_1 = require("hardhat/config");
var constants = require("./addresses");
config_1.task('integration_test_setup', 'Sets up contracts for the integration test', function (args, hre) { return __awaiter(void 0, void 0, void 0, function () {
    var whaleSigner, _i, _a, addr, powers, powerThreshold, Gravity, gravity, TestERC20, testERC20;
    return __generator(this, function (_b) {
        switch (_b.label) {
            case 0: 
            // Take over vitalik.eth
            return [4 /*yield*/, hre.network.provider.request({
                    method: 'hardhat_impersonateAccount',
                    params: [constants.WHALE],
                })];
            case 1:
                // Take over vitalik.eth
                _b.sent();
                return [4 /*yield*/, hre.ethers.getSigner(constants.WHALE)];
            case 2:
                whaleSigner = _b.sent();
                _i = 0, _a = constants.VALIDATORS;
                _b.label = 3;
            case 3:
                if (!(_i < _a.length)) return [3 /*break*/, 6];
                addr = _a[_i];
                return [4 /*yield*/, whaleSigner.sendTransaction({
                        to: addr,
                        value: hre.ethers.utils.parseEther('100'),
                    })];
            case 4:
                _b.sent();
                _b.label = 5;
            case 5:
                _i++;
                return [3 /*break*/, 3];
            case 6:
                powers = [1073741823, 1073741823, 1073741823, 1073741823];
                powerThreshold = 6666;
                return [4 /*yield*/, hre.ethers.getContractFactory("Gravity")];
            case 7:
                Gravity = _b.sent();
                return [4 /*yield*/, Gravity.deploy(hre.ethers.utils.formatBytes32String("gravitytest"), powerThreshold, constants.VALIDATORS, powers)];
            case 8:
                gravity = (_b.sent());
                return [4 /*yield*/, gravity.deployed()];
            case 9:
                _b.sent();
                console.log("gravity contract deployed at - " + gravity.address);
                return [4 /*yield*/, hre.ethers.getContractFactory("TestERC20GB")];
            case 10:
                TestERC20 = _b.sent();
                return [4 /*yield*/, TestERC20.deploy()];
            case 11:
                testERC20 = (_b.sent());
                return [4 /*yield*/, testERC20.deployed()];
            case 12:
                _b.sent();
                console.log("test ERC20 TestGB TGB deployed at - " + testERC20.address);
                return [4 /*yield*/, hre.network.provider.send("evm_setIntervalMining", [1000])];
            case 13:
                _b.sent();
                return [4 /*yield*/, hre.run('node')];
            case 14:
                _b.sent();
                return [2 /*return*/];
        }
    });
}); });
/**
 * @type import('hardhat/config').HardhatUserConfig
 */
var ARCHIVE_NODE_URL = process.env.ARCHIVE_NODE_URL;
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
