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
function prepareTxBatch(batchSize, signers) {
    return __awaiter(this, void 0, void 0, function () {
        var numTxs, destinations, fees, amounts, i, _a, _b;
        return __generator(this, function (_c) {
            switch (_c.label) {
                case 0:
                    numTxs = batchSize;
                    destinations = new Array(numTxs);
                    fees = new Array(numTxs);
                    amounts = new Array(numTxs);
                    i = 0;
                    _c.label = 1;
                case 1:
                    if (!(i < numTxs)) return [3 /*break*/, 4];
                    fees[i] = 1;
                    amounts[i] = 1;
                    _a = destinations;
                    _b = i;
                    return [4 /*yield*/, signers[i + 5].getAddress()];
                case 2:
                    _a[_b] = _c.sent();
                    _c.label = 3;
                case 3:
                    i++;
                    return [3 /*break*/, 1];
                case 4: return [2 /*return*/, {
                        numTxs: numTxs,
                        destinations: destinations,
                        fees: fees,
                        amounts: amounts
                    }];
            }
        });
    });
}
function sendToCosmos(gravity, testERC20, numCoins) {
    return __awaiter(this, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: 
                // Transfer out to Cosmos, locking coins
                // =====================================
                return [4 /*yield*/, testERC20.functions.approve(gravity.address, numCoins)];
                case 1:
                    // Transfer out to Cosmos, locking coins
                    // =====================================
                    _a.sent();
                    return [4 /*yield*/, gravity.functions.sendToCosmos(testERC20.address, hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"), numCoins)];
                case 2:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    });
}
function prep() {
    return __awaiter(this, void 0, void 0, function () {
        var signers, gravityId, powers, validators, powerThreshold, _a, gravity, testERC20, ReentrantERC20Contract, reentrantERC20;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                case 1:
                    signers = _b.sent();
                    gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                    powers = (0, pure_1.examplePowers)();
                    validators = signers.slice(0, powers.length);
                    powerThreshold = 6666;
                    return [4 /*yield*/, (0, test_utils_1.deployContracts)(gravityId, validators, powers, powerThreshold)];
                case 2:
                    _a = _b.sent(), gravity = _a.gravity, testERC20 = _a.testERC20;
                    return [4 /*yield*/, hardhat_1.ethers.getContractFactory("ReentrantERC20")];
                case 3:
                    ReentrantERC20Contract = _b.sent();
                    return [4 /*yield*/, ReentrantERC20Contract.deploy(gravity.address)];
                case 4:
                    reentrantERC20 = (_b.sent());
                    return [2 /*return*/, {
                            signers: signers,
                            gravityId: gravityId,
                            powers: powers,
                            validators: validators,
                            gravity: gravity,
                            testERC20: testERC20,
                            reentrantERC20: reentrantERC20
                        }];
            }
        });
    });
}
function runSubmitBatchTest(opts) {
    return __awaiter(this, void 0, void 0, function () {
        var _a, signers, gravityId, powers, validators, gravity, testERC20, _b, _c, _d, _e, txBatch, batchNonce, batchTimeout, methodName, digest, sigs, valset, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q;
        var _r;
        return __generator(this, function (_s) {
            switch (_s.label) {
                case 0: return [4 /*yield*/, prep()];
                case 1:
                    _a = _s.sent(), signers = _a.signers, gravityId = _a.gravityId, powers = _a.powers, validators = _a.validators, gravity = _a.gravity, testERC20 = _a.testERC20;
                    // Lock tokens in gravity
                    // ====================
                    return [4 /*yield*/, sendToCosmos(gravity, testERC20, 1000)];
                case 2:
                    // Lock tokens in gravity
                    // ====================
                    _s.sent();
                    _b = expect;
                    return [4 /*yield*/, testERC20.functions.balanceOf(gravity.address)];
                case 3:
                    _b.apply(void 0, [(_s.sent())[0].toNumber(),
                        "gravity does not have correct balance after sendToCosmos"]).to.equal(1000);
                    _c = expect;
                    _e = (_d = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[0].getAddress()];
                case 4: return [4 /*yield*/, _e.apply(_d, [_s.sent()])];
                case 5:
                    _c.apply(void 0, [(_s.sent())[0].toNumber(),
                        "msg.sender does not have correct balance after sendToCosmos"]).to.equal(9000);
                    return [4 /*yield*/, prepareTxBatch(opts.batchSize, signers)];
                case 6:
                    txBatch = _s.sent();
                    batchNonce = 1;
                    batchTimeout = 10000;
                    methodName = hardhat_1.ethers.utils.formatBytes32String("transactionBatch");
                    digest = hardhat_1.ethers.utils.keccak256(hardhat_1.ethers.utils.defaultAbiCoder.encode([
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
                        txBatch.amounts,
                        txBatch.destinations,
                        txBatch.fees,
                        batchNonce,
                        testERC20.address,
                        batchTimeout,
                    ]));
                    return [4 /*yield*/, (0, pure_1.signHash)(validators, digest)];
                case 7:
                    sigs = _s.sent();
                    _r = {};
                    return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                case 8:
                    valset = (_r.validators = _s.sent(),
                        _r.powers = powers,
                        _r.valsetNonce = 0,
                        _r.rewardAmount = 0,
                        _r.rewardToken = pure_1.ZeroAddress,
                        _r);
                    return [4 /*yield*/, gravity.submitBatch(valset, sigs, txBatch.amounts, txBatch.destinations, txBatch.fees, 1, testERC20.address, batchTimeout)];
                case 9:
                    _s.sent();
                    _f = expect;
                    _h = (_g = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[5].getAddress()];
                case 10: return [4 /*yield*/, _h.apply(_g, [_s.sent()])];
                case 11:
                    _f.apply(void 0, [(_s.sent())[0].toNumber(),
                        "first address in tx batch does not have correct balance after submitBatch"]).to.equal(1);
                    _j = expect;
                    _l = (_k = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[5 + txBatch.numTxs - 1].getAddress()];
                case 12: return [4 /*yield*/, _l.apply(_k, [_s.sent()])];
                case 13:
                    _j.apply(void 0, [(_s.sent())[0].toNumber(),
                        "last address in tx batch does not have correct balance after submitBatch"]).to.equal(1);
                    _m = expect;
                    return [4 /*yield*/, testERC20.functions.balanceOf(gravity.address)];
                case 14:
                    _m.apply(void 0, [(_s.sent())[0].toNumber(),
                        "gravity does not have correct balance after submitBatch"
                        // Each tx in batch is worth 1 coin sent + 1 coin fee
                    ]).to.equal(1000 - txBatch.numTxs * 2);
                    _o = expect;
                    _q = (_p = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[0].getAddress()];
                case 15: return [4 /*yield*/, _q.apply(_p, [_s.sent()])];
                case 16:
                    _o.apply(void 0, [(_s.sent())[0].toNumber(),
                        "msg.sender does not have correct balance after submitBatch"
                        // msg.sender has received 1 coin in fees for each tx
                    ]).to.equal(9000 + txBatch.numTxs);
                    return [2 /*return*/];
            }
        });
    });
}
function runLogicCallTest(opts) {
    return __awaiter(this, void 0, void 0, function () {
        var _a, signers, gravityId, powers, validators, gravity, testERC20, reentrantERC20, TestTokenBatchMiddleware, tokenBatchMiddleware, _b, _c, _d, _e, txBatch, batchNonce, methodName, logicCallArgs, digest, sigs, valset, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q;
        var _r;
        return __generator(this, function (_s) {
            switch (_s.label) {
                case 0: return [4 /*yield*/, prep()];
                case 1:
                    _a = _s.sent(), signers = _a.signers, gravityId = _a.gravityId, powers = _a.powers, validators = _a.validators, gravity = _a.gravity, testERC20 = _a.testERC20, reentrantERC20 = _a.reentrantERC20;
                    return [4 /*yield*/, hardhat_1.ethers.getContractFactory("TestTokenBatchMiddleware")];
                case 2:
                    TestTokenBatchMiddleware = _s.sent();
                    return [4 /*yield*/, TestTokenBatchMiddleware.deploy()];
                case 3:
                    tokenBatchMiddleware = (_s.sent());
                    return [4 /*yield*/, tokenBatchMiddleware.transferOwnership(gravity.address)];
                case 4:
                    _s.sent();
                    // Lock tokens in gravity
                    // ====================
                    return [4 /*yield*/, sendToCosmos(gravity, testERC20, 1000)];
                case 5:
                    // Lock tokens in gravity
                    // ====================
                    _s.sent();
                    _b = expect;
                    return [4 /*yield*/, testERC20.functions.balanceOf(gravity.address)];
                case 6:
                    _b.apply(void 0, [(_s.sent())[0].toNumber(),
                        "gravity does not have correct balance after sendToCosmos"]).to.equal(1000);
                    _c = expect;
                    _e = (_d = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[0].getAddress()];
                case 7: return [4 /*yield*/, _e.apply(_d, [_s.sent()])];
                case 8:
                    _c.apply(void 0, [(_s.sent())[0].toNumber(),
                        "msg.sender does not have correct balance after sendToCosmos"]).to.equal(9000);
                    return [4 /*yield*/, prepareTxBatch(opts.batchSize, signers)];
                case 9:
                    txBatch = _s.sent();
                    batchNonce = 1;
                    methodName = hardhat_1.ethers.utils.formatBytes32String("logicCall");
                    logicCallArgs = {
                        transferAmounts: [txBatch.numTxs],
                        transferTokenContracts: [testERC20.address],
                        feeAmounts: [txBatch.numTxs],
                        feeTokenContracts: [testERC20.address],
                        logicContractAddress: tokenBatchMiddleware.address,
                        payload: tokenBatchMiddleware.interface.encodeFunctionData("submitBatch", [
                            txBatch.amounts,
                            txBatch.destinations,
                            opts.reentrant ? reentrantERC20.address : testERC20.address,
                        ]),
                        timeOut: 4766922941000,
                        invalidationId: hardhat_1.ethers.utils.hexZeroPad(testERC20.address, 32),
                        invalidationNonce: 1
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
                        "uint256", // invalidationNonce
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
                        logicCallArgs.invalidationNonce,
                    ]));
                    return [4 /*yield*/, (0, pure_1.signHash)(validators, digest)];
                case 10:
                    sigs = _s.sent();
                    _r = {};
                    return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                case 11:
                    valset = (_r.validators = _s.sent(),
                        _r.powers = powers,
                        _r.valsetNonce = 0,
                        _r.rewardAmount = 0,
                        _r.rewardToken = pure_1.ZeroAddress,
                        _r);
                    return [4 /*yield*/, gravity.submitLogicCall(valset, sigs, logicCallArgs)];
                case 12:
                    _s.sent();
                    _f = expect;
                    _h = (_g = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[5].getAddress()];
                case 13: return [4 /*yield*/, _h.apply(_g, [_s.sent()])];
                case 14:
                    _f.apply(void 0, [(_s.sent())[0].toNumber(),
                        "first address in tx batch does not have correct balance after submitLogicCall"]).to.equal(1);
                    _j = expect;
                    _l = (_k = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[5 + txBatch.numTxs - 1].getAddress()];
                case 15: return [4 /*yield*/, _l.apply(_k, [_s.sent()])];
                case 16:
                    _j.apply(void 0, [(_s.sent())[0].toNumber(),
                        "last address in tx batch does not have correct balance after submitLogicCall"]).to.equal(1);
                    _m = expect;
                    return [4 /*yield*/, testERC20.functions.balanceOf(gravity.address)];
                case 17:
                    _m.apply(void 0, [(_s.sent())[0].toNumber(),
                        "gravity does not have correct balance after submitLogicCall"
                        // Each tx in batch is worth 1 coin sent + 1 coin fee
                    ]).to.equal(1000 - txBatch.numTxs * 2);
                    _o = expect;
                    _q = (_p = testERC20.functions).balanceOf;
                    return [4 /*yield*/, signers[0].getAddress()];
                case 18: return [4 /*yield*/, _q.apply(_p, [_s.sent()])];
                case 19:
                    _o.apply(void 0, [(_s.sent())[0].toNumber(),
                        "msg.sender does not have correct balance after submitLogicCall"
                        // msg.sender has received 1 coin in fees for each tx
                    ]).to.equal(9000 + txBatch.numTxs);
                    return [2 /*return*/];
            }
        });
    });
}
describe("Compare gas usage of old submitBatch method vs new logicCall method submitting one batch", function () {
    it("Large batch", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, runSubmitBatchTest({ batchSize: 10 })];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, runLogicCallTest({ batchSize: 10 })];
                    case 2:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("Small batch", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, runSubmitBatchTest({ batchSize: 1 })];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, runLogicCallTest({ batchSize: 1 })];
                    case 2:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("Reentrant", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runLogicCallTest({ batchSize: 1, reentrant: true })).to.be.revertedWith("ReentrancyGuard: reentrant call")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
});
