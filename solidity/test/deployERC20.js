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
var ethers_1 = require("ethers");
chai_1["default"].use(ethereum_waffle_1.solidity);
var expect = chai_1["default"].expect;
function parseEvent(contract, txReceipt, eventOrder) {
    return __awaiter(this, void 0, void 0, function () {
        var receipt, args;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, txReceipt];
                case 1:
                    receipt = _a.sent();
                    if (receipt.events) {
                        args = receipt.events[eventOrder].args;
                        return [2 /*return*/, args];
                    }
                    return [2 /*return*/, undefined];
            }
        });
    });
}
function runTest(opts) {
    return __awaiter(this, void 0, void 0, function () {
        var signers, gravityId, powers, validators, powerThreshold, _a, gravity, testERC20, deployCheckpoint, tx, eventArgs, ERC20contract, maxUint256, _b, numTxs, txDestinationsInt, txFees, txAmounts, i, txDestinations, batchNonce, batchTimeout, methodName, abiEncoded, digest, sigs, currentValsetNonce, valset, _c, _d, _e, _f;
        var _g;
        return __generator(this, function (_h) {
            switch (_h.label) {
                case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                case 1:
                    signers = _h.sent();
                    gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                    powers = (0, pure_1.examplePowers)();
                    validators = signers.slice(0, powers.length);
                    powerThreshold = 6666;
                    return [4 /*yield*/, (0, test_utils_1.deployContracts)(gravityId, validators, powers, powerThreshold)];
                case 2:
                    _a = _h.sent(), gravity = _a.gravity, testERC20 = _a.testERC20, deployCheckpoint = _a.checkpoint;
                    return [4 /*yield*/, gravity.deployERC20('uatom', 'Atom', 'ATOM', 6)];
                case 3:
                    tx = _h.sent();
                    return [4 /*yield*/, parseEvent(gravity, tx.wait(), 1)];
                case 4:
                    eventArgs = _h.sent();
                    if (eventArgs == undefined) {
                        throw new Error("No event args");
                    }
                    expect(eventArgs._cosmosDenom).to.equal('uatom');
                    ERC20contract = new hardhat_1.ethers.Contract(eventArgs._tokenContract, [
                        "function balanceOf(address account) view returns (uint256 balance)"
                    ], gravity.provider);
                    maxUint256 = ethers_1.BigNumber.from(2).pow(256).sub(1);
                    // Check that gravity balance is correct
                    _b = expect;
                    return [4 /*yield*/, ERC20contract.functions.balanceOf(gravity.address)];
                case 5:
                    // Check that gravity balance is correct
                    _b.apply(void 0, [(_h.sent()).toString()]).to.equal(maxUint256.toString());
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
                case 6:
                    txDestinations = _h.sent();
                    batchNonce = 1;
                    batchTimeout = 10000;
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
                        eventArgs._tokenContract,
                        batchTimeout
                    ]);
                    digest = hardhat_1.ethers.utils.keccak256(abiEncoded);
                    return [4 /*yield*/, (0, pure_1.signHash)(validators, digest)];
                case 7:
                    sigs = _h.sent();
                    currentValsetNonce = 0;
                    _g = {};
                    return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                case 8:
                    valset = (_g.validators = _h.sent(),
                        _g.powers = powers,
                        _g.valsetNonce = currentValsetNonce,
                        _g.rewardAmount = 0,
                        _g.rewardToken = pure_1.ZeroAddress,
                        _g);
                    return [4 /*yield*/, gravity.submitBatch(valset, sigs, txAmounts, txDestinations, txFees, batchNonce, eventArgs._tokenContract, batchTimeout)];
                case 9:
                    _h.sent();
                    // Check that Gravity's balance is correct
                    _c = expect;
                    return [4 /*yield*/, ERC20contract.functions.balanceOf(gravity.address)];
                case 10:
                    // Check that Gravity's balance is correct
                    _c.apply(void 0, [(_h.sent()).toString()]).to.equal(maxUint256.sub(200).toString());
                    // Check that one of the recipient's balance is correct
                    _d = expect;
                    _f = (_e = ERC20contract.functions).balanceOf;
                    return [4 /*yield*/, signers[6].getAddress()];
                case 11: return [4 /*yield*/, _f.apply(_e, [_h.sent()])];
                case 12:
                    // Check that one of the recipient's balance is correct
                    _d.apply(void 0, [(_h.sent()).toString()]).to.equal('1');
                    return [2 /*return*/];
            }
        });
    });
}
describe("deployERC20 tests", function () {
    // There is no way for this function to throw so there are
    // no throwing tests
    it("runs", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, runTest({})];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
});
