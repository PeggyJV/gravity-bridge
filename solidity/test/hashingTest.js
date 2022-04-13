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
exports.makeCheckpoint = void 0;
var chai_1 = require("chai");
var hardhat_1 = require("hardhat");
var ethereum_waffle_1 = require("ethereum-waffle");
var pure_1 = require("../test-utils/pure");
chai_1["default"].use(ethereum_waffle_1.solidity);
var expect = chai_1["default"].expect;
describe("Hashing test", function () {
    it("Hashing test", function () {
        return __awaiter(this, void 0, void 0, function () {
            var signers, gravityId, validators, powers, i, HashingTest, hashingContract, _a, _b, _c, _d, _e, _f, contractCheckpoint, externalCheckpoint, _g, _h, _j, _k, _l;
            return __generator(this, function (_m) {
                switch (_m.label) {
                    case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                    case 1:
                        signers = _m.sent();
                        gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                        validators = [];
                        powers = [];
                        for (i = 0; i < 100; i++) {
                            validators.push(signers[i]);
                            powers.push(5000);
                        }
                        return [4 /*yield*/, hardhat_1.ethers.getContractFactory("HashingTest")];
                    case 2:
                        HashingTest = _m.sent();
                        return [4 /*yield*/, HashingTest.deploy()];
                    case 3:
                        hashingContract = (_m.sent());
                        return [4 /*yield*/, hashingContract.deployed()];
                    case 4:
                        _m.sent();
                        _b = (_a = hashingContract).IterativeHash;
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                    case 5: return [4 /*yield*/, _b.apply(_a, [_m.sent(), powers,
                            1,
                            gravityId])];
                    case 6:
                        _m.sent();
                        _d = (_c = hashingContract).ConcatHash;
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                    case 7: return [4 /*yield*/, _d.apply(_c, [_m.sent(), powers,
                            1,
                            gravityId])];
                    case 8:
                        _m.sent();
                        _f = (_e = hashingContract).ConcatHash2;
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                    case 9: return [4 /*yield*/, _f.apply(_e, [_m.sent(), powers,
                            1,
                            gravityId])];
                    case 10:
                        _m.sent();
                        return [4 /*yield*/, hashingContract.lastCheckpoint()];
                    case 11:
                        contractCheckpoint = _m.sent();
                        _g = makeCheckpoint;
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                    case 12:
                        externalCheckpoint = _g.apply(void 0, [_m.sent(), powers,
                            1,
                            gravityId]);
                        expect(contractCheckpoint === externalCheckpoint);
                        _j = (_h = hashingContract).JustSaveEverything;
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                    case 13: return [4 /*yield*/, _j.apply(_h, [_m.sent(), powers,
                            1])];
                    case 14:
                        _m.sent();
                        _l = (_k = hashingContract).JustSaveEverythingAgain;
                        return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                    case 15: return [4 /*yield*/, _l.apply(_k, [_m.sent(), powers,
                            1])];
                    case 16:
                        _m.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
});
function makeCheckpoint(validators, powers, valsetNonce, gravityId) {
    var methodName = hardhat_1.ethers.utils.formatBytes32String("checkpoint");
    var abiEncoded = hardhat_1.ethers.utils.defaultAbiCoder.encode(["bytes32", "bytes32", "uint256", "address[]", "uint256[]"], [gravityId, methodName, valsetNonce, validators, powers]);
    var checkpoint = hardhat_1.ethers.utils.keccak256(abiEncoded);
    return checkpoint;
}
exports.makeCheckpoint = makeCheckpoint;
