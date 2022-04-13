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
describe("Gravity happy path valset update + batch submit", function () {
    it("Happy path", function () {
        return __awaiter(this, void 0, void 0, function () {
            var signers, gravityId, valset0, powerThreshold, _a, gravity, testERC20, deployCheckpoint, valset1, valset0_str, valset1_str, checkpoint1, sigs1, _b, numTxs, txDestinationsInt, txFees, totalFees, txAmounts, i, txDestinations, batchNonce, batchTimeout, methodName, abiEncoded, digest, sigs, batchSubmitTx, _c, _d, _e, _f, _g, _h;
            var _j, _k;
            return __generator(this, function (_l) {
                switch (_l.label) {
                    case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                    case 1:
                        signers = _l.sent();
                        gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                        valset0 = {
                            // This is the power distribution on the Cosmos hub as of 7/14/2020
                            powers: (0, pure_1.examplePowers)(),
                            validators: signers.slice(0, (0, pure_1.examplePowers)().length),
                            valsetNonce: 0,
                            rewardAmount: 0,
                            rewardToken: pure_1.ZeroAddress
                        };
                        powerThreshold = 6666;
                        return [4 /*yield*/, (0, test_utils_1.deployContracts)(gravityId, valset0.validators, valset0.powers, powerThreshold)];
                    case 2:
                        _a = _l.sent(), gravity = _a.gravity, testERC20 = _a.testERC20, deployCheckpoint = _a.checkpoint;
                        valset1 = (function () {
                            // Make new valset by modifying some powers
                            var powers = (0, pure_1.examplePowers)();
                            powers[0] -= 3;
                            powers[1] += 3;
                            var validators = signers.slice(0, powers.length);
                            return {
                                powers: powers,
                                validators: validators,
                                valsetNonce: 1,
                                rewardAmount: 0,
                                rewardToken: pure_1.ZeroAddress
                            };
                        })();
                        _j = {
                            powers: valset0.powers
                        };
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(valset0.validators)];
                    case 3:
                        valset0_str = (_j.validators = _l.sent(),
                            _j.valsetNonce = valset0.valsetNonce,
                            _j.rewardAmount = valset0.rewardAmount,
                            _j.rewardToken = valset0.rewardToken,
                            _j);
                        _k = {
                            powers: valset1.powers
                        };
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(valset1.validators)];
                    case 4:
                        valset1_str = (_k.validators = _l.sent(),
                            _k.valsetNonce = valset1.valsetNonce,
                            _k.rewardAmount = valset1.rewardAmount,
                            _k.rewardToken = valset1.rewardToken,
                            _k);
                        checkpoint1 = (0, pure_1.makeCheckpoint)(valset1_str.validators, valset1_str.powers, valset1_str.valsetNonce, valset1_str.rewardAmount, valset1_str.rewardToken, gravityId);
                        return [4 /*yield*/, (0, pure_1.signHash)(valset0.validators, checkpoint1)];
                    case 5:
                        sigs1 = _l.sent();
                        return [4 /*yield*/, gravity.updateValset(valset1_str, valset0_str, sigs1)];
                    case 6:
                        _l.sent();
                        _b = expect;
                        return [4 /*yield*/, gravity.functions.state_lastValsetCheckpoint()];
                    case 7:
                        _b.apply(void 0, [(_l.sent())[0]]).to.equal(checkpoint1);
                        // SUBMITBATCH
                        // ==========================
                        // Transfer out to Cosmos, locking coins
                        return [4 /*yield*/, testERC20.functions.approve(gravity.address, 1000)];
                    case 8:
                        // SUBMITBATCH
                        // ==========================
                        // Transfer out to Cosmos, locking coins
                        _l.sent();
                        return [4 /*yield*/, gravity.functions.sendToCosmos(testERC20.address, hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"), 1000)];
                    case 9:
                        _l.sent();
                        numTxs = 100;
                        txDestinationsInt = new Array(numTxs);
                        txFees = new Array(numTxs);
                        totalFees = 0;
                        txAmounts = new Array(numTxs);
                        for (i = 0; i < numTxs; i++) {
                            txFees[i] = 1;
                            totalFees += 1;
                            txAmounts[i] = 1;
                            txDestinationsInt[i] = signers[i + 5];
                        }
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(txDestinationsInt)];
                    case 10:
                        txDestinations = _l.sent();
                        batchNonce = 1;
                        batchTimeout = 10000000;
                        methodName = hardhat_1.ethers.utils.formatBytes32String("transactionBatch");
                        abiEncoded = hardhat_1.ethers.utils.defaultAbiCoder.encode([
                            "bytes32",
                            "bytes32",
                            "uint256[]",
                            "address[]",
                            "uint256[]",
                            "uint256",
                            "address",
                            "uint256"
                        ], [
                            gravityId,
                            methodName,
                            txAmounts,
                            txDestinations,
                            txFees,
                            batchNonce,
                            testERC20.address,
                            batchTimeout
                        ]);
                        digest = hardhat_1.ethers.utils.keccak256(abiEncoded);
                        return [4 /*yield*/, (0, pure_1.signHash)(valset1.validators, digest)];
                    case 11:
                        sigs = _l.sent();
                        return [4 /*yield*/, gravity.submitBatch(valset1_str, sigs, txAmounts, txDestinations, txFees, batchNonce, testERC20.address, batchTimeout)];
                    case 12:
                        batchSubmitTx = _l.sent();
                        // check that the transfer was successful
                        _c = expect;
                        _e = (_d = testERC20.functions).balanceOf;
                        return [4 /*yield*/, signers[6].getAddress()];
                    case 13: return [4 /*yield*/, _e.apply(_d, [_l.sent()])];
                    case 14: return [4 /*yield*/, (_l.sent())[0].toNumber()];
                    case 15:
                        // check that the transfer was successful
                        _c.apply(void 0, [_l.sent()]).to.equal(1);
                        // check that the relayer was paid
                        _f = expect;
                        _h = (_g = testERC20.functions).balanceOf;
                        return [4 /*yield*/, batchSubmitTx.from];
                    case 16: return [4 /*yield*/, _h.apply(_g, [_l.sent()])];
                    case 17: return [4 /*yield*/, (_l.sent())[0].toNumber()];
                    case 18:
                        // check that the relayer was paid
                        _f.apply(void 0, [_l.sent()]).to.equal(9000 + totalFees);
                        return [2 /*return*/];
                }
            });
        });
    });
});
