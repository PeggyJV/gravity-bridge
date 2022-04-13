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
        var signers, gravityId, powers, validators, powerThreshold, _a, gravity, testERC20, deployCheckpoint, SimpleLogicBatchMiddleware, logicBatch, TestLogicContract, logicContract, numTxs, txPayloads, txAmounts, i, _b, _c, _d, _e, _f, invalidationNonce, timeOut, methodName, logicCallArgs, digest, sigs, currentValsetNonce, valset, logicCallSubmitResult, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s;
        var _t;
        return __generator(this, function (_u) {
            switch (_u.label) {
                case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                case 1:
                    signers = _u.sent();
                    gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                    powers = (0, pure_1.examplePowers)();
                    validators = signers.slice(0, powers.length);
                    powerThreshold = 6666;
                    return [4 /*yield*/, (0, test_utils_1.deployContracts)(gravityId, validators, powers, powerThreshold)];
                case 2:
                    _a = _u.sent(), gravity = _a.gravity, testERC20 = _a.testERC20, deployCheckpoint = _a.checkpoint;
                    return [4 /*yield*/, hardhat_1.ethers.getContractFactory("SimpleLogicBatchMiddleware")];
                case 3:
                    SimpleLogicBatchMiddleware = _u.sent();
                    return [4 /*yield*/, SimpleLogicBatchMiddleware.deploy()];
                case 4:
                    logicBatch = (_u.sent());
                    // We set the ownership to gravity so that nobody else can call it.
                    return [4 /*yield*/, logicBatch.transferOwnership(gravity.address)];
                case 5:
                    // We set the ownership to gravity so that nobody else can call it.
                    _u.sent();
                    return [4 /*yield*/, hardhat_1.ethers.getContractFactory("TestLogicContract")];
                case 6:
                    TestLogicContract = _u.sent();
                    return [4 /*yield*/, TestLogicContract.deploy(testERC20.address)];
                case 7:
                    logicContract = (_u.sent());
                    // We set its owner to the batch contract. 
                    return [4 /*yield*/, logicContract.transferOwnership(logicBatch.address)];
                case 8:
                    // We set its owner to the batch contract. 
                    _u.sent();
                    // Transfer out to Cosmos, locking coins
                    // =====================================
                    return [4 /*yield*/, testERC20.functions.approve(gravity.address, 1000)];
                case 9:
                    // Transfer out to Cosmos, locking coins
                    // =====================================
                    _u.sent();
                    return [4 /*yield*/, gravity.functions.sendToCosmos(testERC20.address, hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"), 1000)];
                case 10:
                    _u.sent();
                    numTxs = 10;
                    txPayloads = new Array(numTxs);
                    txAmounts = new Array(numTxs);
                    i = 0;
                    _u.label = 11;
                case 11:
                    if (!(i < numTxs)) return [3 /*break*/, 14];
                    txAmounts[i] = 5;
                    _b = txPayloads;
                    _c = i;
                    _e = (_d = logicContract.interface).encodeFunctionData;
                    _f = ["transferTokens"];
                    return [4 /*yield*/, signers[20].getAddress()];
                case 12:
                    _b[_c] = _e.apply(_d, _f.concat([[_u.sent(), 2, 2]]));
                    _u.label = 13;
                case 13:
                    i++;
                    return [3 /*break*/, 11];
                case 14:
                    invalidationNonce = 1;
                    if (opts.invalidationNonceNotHigher) {
                        invalidationNonce = 0;
                    }
                    timeOut = 4766922941000;
                    if (opts.timedOut) {
                        timeOut = 0;
                    }
                    methodName = hardhat_1.ethers.utils.formatBytes32String("logicCall");
                    logicCallArgs = {
                        transferAmounts: [numTxs * 5],
                        transferTokenContracts: [testERC20.address],
                        feeAmounts: [numTxs],
                        feeTokenContracts: [testERC20.address],
                        logicContractAddress: logicBatch.address,
                        payload: logicBatch.interface.encodeFunctionData("logicBatch", [txAmounts, txPayloads, logicContract.address, testERC20.address]),
                        timeOut: timeOut,
                        invalidationId: hardhat_1.ethers.utils.hexZeroPad(testERC20.address, 32),
                        invalidationNonce: invalidationNonce // invalidationNonce
                    };
                    digest = hardhat_1.ethers.utils.keccak256(hardhat_1.ethers.utils.defaultAbiCoder.encode([
                        "bytes32",
                        "bytes32",
                        "uint256[]",
                        "address[]",
                        "uint256[]",
                        "address[]",
                        "address",
                        "bytes",
                        "uint256",
                        "bytes32",
                        "uint256" // invalidationNonce
                    ], [
                        gravityId,
                        methodName,
                        logicCallArgs.transferAmounts,
                        logicCallArgs.transferTokenContracts,
                        logicCallArgs.feeAmounts,
                        logicCallArgs.feeTokenContracts,
                        logicCallArgs.logicContractAddress,
                        logicCallArgs.payload,
                        logicCallArgs.timeOut,
                        logicCallArgs.invalidationId,
                        logicCallArgs.invalidationNonce
                    ]));
                    return [4 /*yield*/, (0, pure_1.signHash)(validators, digest)];
                case 15:
                    sigs = _u.sent();
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
                    _t = {};
                    return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                case 16:
                    valset = (_t.validators = _u.sent(),
                        _t.powers = powers,
                        _t.valsetNonce = currentValsetNonce,
                        _t.rewardAmount = 0,
                        _t.rewardToken = pure_1.ZeroAddress,
                        _t);
                    return [4 /*yield*/, gravity.submitLogicCall(valset, sigs, logicCallArgs)];
                case 17:
                    logicCallSubmitResult = _u.sent();
                    // check that the relayer was paid
                    _g = expect;
                    _j = (_h = testERC20.functions).balanceOf;
                    return [4 /*yield*/, logicCallSubmitResult.from];
                case 18: return [4 /*yield*/, _j.apply(_h, [_u.sent()])];
                case 19: return [4 /*yield*/, (_u.sent())[0].toNumber()];
                case 20:
                    // check that the relayer was paid
                    _g.apply(void 0, [_u.sent()]).to.equal(9010);
                    _k = expect;
                    _m = (_l = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[20].getAddress()];
                case 21: return [4 /*yield*/, _m.apply(_l, [_u.sent()])];
                case 22:
                    _k.apply(void 0, [(_u.sent())[0].toNumber()]).to.equal(40);
                    _o = expect;
                    return [4 /*yield*/, testERC20.functions.balanceOf(gravity.address)];
                case 23:
                    _o.apply(void 0, [(_u.sent())[0].toNumber()]).to.equal(940);
                    _p = expect;
                    return [4 /*yield*/, testERC20.functions.balanceOf(logicContract.address)];
                case 24:
                    _p.apply(void 0, [(_u.sent())[0].toNumber()]).to.equal(10);
                    _q = expect;
                    _s = (_r = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[0].getAddress()];
                case 25: return [4 /*yield*/, _s.apply(_r, [_u.sent()])];
                case 26:
                    _q.apply(void 0, [(_u.sent())[0].toNumber()]).to.equal(9010);
                    return [2 /*return*/];
            }
        });
    });
}
describe("submitLogicCall tests", function () {
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
    it("throws on invalidation nonce not incremented", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ invalidationNonceNotHigher: true })).to.be.revertedWith("InvalidLogicCallNonce(0, 0)")];
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
    it("throws on timeout", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ timedOut: true })).to.be.revertedWith("LogicCallTimedOut()")];
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
describe("logicCall Go test hash", function () {
    it("produces good hash", function () {
        return __awaiter(this, void 0, void 0, function () {
            var signers, gravityId, powers, validators, powerThreshold, _a, gravity, testERC20, deployCheckpoint, methodName, numTxs, invalidationNonce, timeOut, logicCallArgs, abiEncodedLogicCall, logicCallDigest, sigs, currentValsetNonce, valset, res, _b, _c, _d, _e;
            var _f, _g;
            return __generator(this, function (_h) {
                switch (_h.label) {
                    case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                    case 1:
                        signers = _h.sent();
                        gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                        powers = [6667];
                        validators = signers.slice(0, powers.length);
                        powerThreshold = 6666;
                        return [4 /*yield*/, (0, test_utils_1.deployContracts)(gravityId, validators, powers, powerThreshold)];
                    case 2:
                        _a = _h.sent(), gravity = _a.gravity, testERC20 = _a.testERC20, deployCheckpoint = _a.checkpoint;
                        // Transfer out to Cosmos, locking coins
                        // =====================================
                        return [4 /*yield*/, testERC20.functions.approve(gravity.address, 1000)];
                    case 3:
                        // Transfer out to Cosmos, locking coins
                        // =====================================
                        _h.sent();
                        return [4 /*yield*/, gravity.functions.sendToCosmos(testERC20.address, hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"), 1000)];
                    case 4:
                        _h.sent();
                        methodName = hardhat_1.ethers.utils.formatBytes32String("logicCall");
                        numTxs = 10;
                        invalidationNonce = 1;
                        timeOut = 4766922941000;
                        logicCallArgs = {
                            transferAmounts: [1],
                            transferTokenContracts: [testERC20.address],
                            feeAmounts: [1],
                            feeTokenContracts: [testERC20.address],
                            logicContractAddress: "0x17c1736CcF692F653c433d7aa2aB45148C016F68",
                            payload: hardhat_1.ethers.utils.formatBytes32String("testingPayload"),
                            timeOut: timeOut,
                            invalidationId: hardhat_1.ethers.utils.formatBytes32String("invalidationId"),
                            invalidationNonce: invalidationNonce // invalidationNonce
                        };
                        abiEncodedLogicCall = hardhat_1.ethers.utils.defaultAbiCoder.encode([
                            "bytes32",
                            "bytes32",
                            "uint256[]",
                            "address[]",
                            "uint256[]",
                            "address[]",
                            "address",
                            "bytes",
                            "uint256",
                            "bytes32",
                            "uint256" // invalidationNonce
                        ], [
                            gravityId,
                            methodName,
                            logicCallArgs.transferAmounts,
                            logicCallArgs.transferTokenContracts,
                            logicCallArgs.feeAmounts,
                            logicCallArgs.feeTokenContracts,
                            logicCallArgs.logicContractAddress,
                            logicCallArgs.payload,
                            logicCallArgs.timeOut,
                            logicCallArgs.invalidationId,
                            logicCallArgs.invalidationNonce
                        ]);
                        logicCallDigest = hardhat_1.ethers.utils.keccak256(abiEncodedLogicCall);
                        return [4 /*yield*/, (0, pure_1.signHash)(validators, logicCallDigest)];
                    case 5:
                        sigs = _h.sent();
                        currentValsetNonce = 0;
                        _f = {};
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                    case 6:
                        valset = (_f.validators = _h.sent(),
                            _f.powers = powers,
                            _f.valsetNonce = currentValsetNonce,
                            _f.rewardAmount = 0,
                            _f.rewardToken = pure_1.ZeroAddress,
                            _f);
                        return [4 /*yield*/, gravity.populateTransaction.submitLogicCall(valset, sigs, logicCallArgs)];
                    case 7:
                        res = _h.sent();
                        console.log("elements in logic call digest:", {
                            "gravityId": gravityId,
                            "logicMethodName": methodName,
                            "transferAmounts": logicCallArgs.transferAmounts,
                            "transferTokenContracts": logicCallArgs.transferTokenContracts,
                            "feeAmounts": logicCallArgs.feeAmounts,
                            "feeTokenContracts": logicCallArgs.feeTokenContracts,
                            "logicContractAddress": logicCallArgs.logicContractAddress,
                            "payload": logicCallArgs.payload,
                            "timeout": logicCallArgs.timeOut,
                            "invalidationId": logicCallArgs.invalidationId,
                            "invalidationNonce": logicCallArgs.invalidationNonce
                        });
                        console.log("abiEncodedCall:", abiEncodedLogicCall);
                        console.log("callDigest:", logicCallDigest);
                        _c = (_b = console).log;
                        _d = ["elements in logic call function call:"];
                        _g = {};
                        _e = "currentValidators";
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                    case 8:
                        _c.apply(_b, _d.concat([(_g[_e] = _h.sent(),
                                _g["currentPowers"] = powers,
                                _g["currentValsetNonce"] = currentValsetNonce,
                                _g["sigs"] = sigs,
                                _g)]));
                        console.log("Function call bytes:", res.data);
                        return [2 /*return*/];
                }
            });
        });
    });
});
