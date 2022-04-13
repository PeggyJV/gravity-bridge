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
exports.__esModule = true;
var chai_1 = require("chai");
var hardhat_1 = require("hardhat");
var ethereum_waffle_1 = require("ethereum-waffle");
var test_utils_1 = require("../test-utils");
var pure_1 = require("../test-utils/pure");
chai_1["default"].use(ethereum_waffle_1.solidity);
var expect = chai_1["default"].expect;
function runTest(opts) {
    return __awaiter(this, void 0, void 0, function () {
        var signers, gravityId, powers, validators, powerThreshold, _a, gravity, testERC20, deployCheckpoint, numTxs, txDestinationsInt, txFees, txAmounts, i, txDestinations, batchTimeout, batchNonce, methodName, abiEncoded, digest, sigs, currentValsetNonce, valset, batchSubmitTx;
        var _b;
        return __generator(this, function (_c) {
            switch (_c.label) {
                case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                case 1:
                    signers = _c.sent();
                    gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                    powers = (0, pure_1.examplePowers)();
                    validators = signers.slice(0, powers.length);
                    powerThreshold = 6666;
                    return [4 /*yield*/, (0, test_utils_1.deployContracts)(gravityId, validators, powers, powerThreshold)];
                case 2:
                    _a = _c.sent(), gravity = _a.gravity, testERC20 = _a.testERC20, deployCheckpoint = _a.checkpoint;
                    // Transfer out to Cosmos, locking coins
                    // =====================================
                    return [4 /*yield*/, testERC20.functions.approve(gravity.address, 1000)];
                case 3:
                    // Transfer out to Cosmos, locking coins
                    // =====================================
                    _c.sent();
                    return [4 /*yield*/, gravity.functions.sendToCosmos(testERC20.address, hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"), 1000)];
                case 4:
                    _c.sent();
                    numTxs = 100;
                    txDestinationsInt = new Array(numTxs);
                    txFees = new Array(numTxs);
                    txAmounts = new Array(numTxs);
                    for (i = 0; i < numTxs; i++) {
                        txFees[i] = 1;
                        txAmounts[i] = 1;
                        txDestinationsInt[i] = signers[i + 5];
                    }
                    return [4 /*yield*/, (0, pure_1.getSignerAddresses)(txDestinationsInt)];
                case 5:
                    txDestinations = _c.sent();
                    if (opts.malformedTxBatch) {
                        // Make the fees array the wrong size
                        txFees.pop();
                    }
                    batchTimeout = hardhat_1.ethers.provider.blockNumber + 1000;
                    if (opts.batchTimeout) {
                        batchTimeout = hardhat_1.ethers.provider.blockNumber - 1;
                    }
                    batchNonce = 1;
                    if (opts.batchNonceNotHigher) {
                        batchNonce = 0;
                    }
                    methodName = hardhat_1.ethers.utils.formatBytes32String("transactionBatch");
                    abiEncoded = hardhat_1.ethers.utils.defaultAbiCoder.encode([
                        "bytes32",
                        "bytes32",
                        "uint256[]",
                        "address[]",
                        "uint256[]",
                        "uint256",
                        "address",
                        "uint256",
                    ], [
                        gravityId,
                        methodName,
                        txAmounts,
                        txDestinations,
                        txFees,
                        batchNonce,
                        testERC20.address,
                        batchTimeout,
                    ]);
                    digest = hardhat_1.ethers.utils.keccak256(abiEncoded);
                    return [4 /*yield*/, (0, pure_1.signHash)(validators, digest)];
                case 6:
                    sigs = _c.sent();
                    currentValsetNonce = 0;
                    if (opts.nonMatchingCurrentValset) {
                        // Wrong nonce
                        currentValsetNonce = 420;
                    }
                    if (opts.malformedCurrentValset) {
                        // Remove one of the powers to make the length not match
                        powers.pop();
                    }
                    if (opts.badValidatorSig) {
                        // Switch the first sig for the second sig to screw things up
                        sigs[1].v = sigs[0].v;
                        sigs[1].r = sigs[0].r;
                        sigs[1].s = sigs[0].s;
                    }
                    if (opts.zeroedValidatorSig) {
                        // Switch the first sig for the second sig to screw things up
                        sigs[1].v = sigs[0].v;
                        sigs[1].r = sigs[0].r;
                        sigs[1].s = sigs[0].s;
                        // Then zero it out to skip evaluation
                        sigs[1].v = 0;
                    }
                    if (opts.notEnoughPower) {
                        // zero out enough signatures that we dip below the threshold
                        sigs[1].v = 0;
                        sigs[2].v = 0;
                        sigs[3].v = 0;
                        sigs[5].v = 0;
                        sigs[6].v = 0;
                        sigs[7].v = 0;
                        sigs[9].v = 0;
                        sigs[11].v = 0;
                        sigs[13].v = 0;
                    }
                    if (opts.barelyEnoughPower) {
                        // Stay just above the threshold
                        sigs[1].v = 0;
                        sigs[2].v = 0;
                        sigs[3].v = 0;
                        sigs[5].v = 0;
                        sigs[6].v = 0;
                        sigs[7].v = 0;
                        sigs[9].v = 0;
                        sigs[11].v = 0;
                    }
                    _b = {};
                    return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                case 7:
                    valset = (_b.validators = _c.sent(),
                        _b.powers = powers,
                        _b.valsetNonce = currentValsetNonce,
                        _b.rewardAmount = 0,
                        _b.rewardToken = pure_1.ZeroAddress,
                        _b);
                    return [4 /*yield*/, gravity.submitBatch(valset, sigs, txAmounts, txDestinations, txFees, batchNonce, testERC20.address, batchTimeout)];
                case 8:
                    batchSubmitTx = _c.sent();
                    return [2 /*return*/];
            }
        });
    });
}
describe("submitBatch tests", function () {
    it("throws on malformed current valset", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ malformedCurrentValset: true })).to.be.revertedWith("MalformedCurrentValidatorSet()")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("throws on malformed txbatch", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ malformedTxBatch: true })).to.be.revertedWith("MalformedBatch()")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("throws on batch nonce not incremented", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ batchNonceNotHigher: true })).to.be.revertedWith("InvalidBatchNonce(0, 0)")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("throws on timeout batch", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ batchTimeout: true })).to.be.revertedWith("BatchTimedOut()")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("throws on non matching checkpoint for current valset", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ nonMatchingCurrentValset: true })).to.be.revertedWith("IncorrectCheckpoint()")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("throws on bad validator sig", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ badValidatorSig: true })).to.be.revertedWith("InvalidSignature()")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("allows zeroed sig", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, runTest({ zeroedValidatorSig: true })];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("throws on not enough signatures", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ notEnoughPower: true })).to.be.revertedWith("InsufficientPower(6537, 6666)")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("does not throw on barely enough signatures", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, runTest({ barelyEnoughPower: true })];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
});
// This test produces a hash for the contract which should match what is being used in the Go unit tests. It's here for
// the use of anyone updating the Go tests.
describe("submitBatch Go test hash", function () {
    it("produces good hash", function () {
        return __awaiter(this, void 0, void 0, function () {
            var signers, gravityId, powers, validators, powerThreshold, _a, gravity, testERC20, deployCheckpoint, txAmounts, txFees, txDestinations, batchNonce, batchTimeout, batchMethodName, abiEncodedBatch, batchDigest, sigs, currentValsetNonce, valset;
            var _b;
            return __generator(this, function (_c) {
                switch (_c.label) {
                    case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                    case 1:
                        signers = _c.sent();
                        gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                        powers = [6667];
                        validators = signers.slice(0, powers.length);
                        powerThreshold = 6666;
                        return [4 /*yield*/, (0, test_utils_1.deployContracts)(gravityId, validators, powers, powerThreshold)];
                    case 2:
                        _a = _c.sent(), gravity = _a.gravity, testERC20 = _a.testERC20, deployCheckpoint = _a.checkpoint;
                        txAmounts = [1];
                        txFees = [1];
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)([signers[5]])];
                    case 3:
                        txDestinations = _c.sent();
                        batchNonce = 1;
                        batchTimeout = hardhat_1.ethers.provider.blockNumber + 1000;
                        // Transfer out to Cosmos, locking coins
                        // =====================================
                        return [4 /*yield*/, testERC20.functions.approve(gravity.address, 1000)];
                    case 4:
                        // Transfer out to Cosmos, locking coins
                        // =====================================
                        _c.sent();
                        return [4 /*yield*/, gravity.functions.sendToCosmos(testERC20.address, hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"), 1000)];
                    case 5:
                        _c.sent();
                        batchMethodName = hardhat_1.ethers.utils.formatBytes32String("transactionBatch");
                        abiEncodedBatch = hardhat_1.ethers.utils.defaultAbiCoder.encode([
                            "bytes32",
                            "bytes32",
                            "uint256[]",
                            "address[]",
                            "uint256[]",
                            "uint256",
                            "address",
                            "uint256",
                        ], [
                            gravityId,
                            batchMethodName,
                            txAmounts,
                            txDestinations,
                            txFees,
                            batchNonce,
                            testERC20.address,
                            batchTimeout,
                        ]);
                        batchDigest = hardhat_1.ethers.utils.keccak256(abiEncodedBatch);
                        console.log("elements in batch digest:", {
                            gravityId: gravityId,
                            batchMethodName: batchMethodName,
                            txAmounts: txAmounts,
                            txDestinations: txDestinations,
                            txFees: txFees,
                            batchNonce: batchNonce,
                            batchTimeout: batchTimeout,
                            tokenContract: testERC20.address
                        });
                        console.log("abiEncodedBatch:", abiEncodedBatch);
                        console.log("batchDigest:", batchDigest);
                        return [4 /*yield*/, (0, pure_1.signHash)(validators, batchDigest)];
                    case 6:
                        sigs = _c.sent();
                        currentValsetNonce = 0;
                        _b = {};
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                    case 7:
                        valset = (_b.validators = _c.sent(),
                            _b.powers = powers,
                            _b.valsetNonce = currentValsetNonce,
                            _b.rewardAmount = 0,
                            _b.rewardToken = pure_1.ZeroAddress,
                            _b);
                        return [4 /*yield*/, gravity.submitBatch(valset, sigs, txAmounts, txDestinations, txFees, batchNonce, testERC20.address, batchTimeout)];
                    case 8:
                        _c.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
});
