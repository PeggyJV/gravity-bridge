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
        var signers, gravityId, powers, validators, powerThreshold, _a, gravity, testERC20, deployCheckpoint, newPowers, newValidators, i, currentValsetNonce, newValsetNonce, currentValset, newValset, ERC20contract, eventArgs, checkpoint, sigs, valsetUpdateTx, _b, _c, _d;
        var _e, _f;
        return __generator(this, function (_g) {
            switch (_g.label) {
                case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                case 1:
                    signers = _g.sent();
                    gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                    powers = (0, pure_1.examplePowers)();
                    validators = signers.slice(0, powers.length);
                    powerThreshold = 6666;
                    return [4 /*yield*/, (0, test_utils_1.deployContracts)(gravityId, validators, powers, powerThreshold)];
                case 2:
                    _a = _g.sent(), gravity = _a.gravity, testERC20 = _a.testERC20, deployCheckpoint = _a.checkpoint;
                    newPowers = (0, pure_1.examplePowers)();
                    newPowers[0] -= 3;
                    newPowers[1] += 3;
                    newValidators = signers.slice(0, newPowers.length);
                    if (opts.malformedNewValset) {
                        // Validators and powers array don't match
                        newValidators = signers.slice(0, newPowers.length - 1);
                    }
                    else if (opts.zeroLengthValset) {
                        newValidators = [];
                        newPowers = [];
                    }
                    else if (opts.notEnoughPowerNewSet) {
                        for (i in newPowers) {
                            newPowers[i] = 5;
                        }
                    }
                    currentValsetNonce = 0;
                    if (opts.nonMatchingCurrentValset) {
                        powers[0] = 78;
                    }
                    newValsetNonce = 1;
                    if (opts.nonceNotIncremented) {
                        newValsetNonce = 0;
                    }
                    _e = {};
                    return [4 /*yield*/, (0, pure_1.getSignerAddresses)(validators)];
                case 3:
                    currentValset = (_e.validators = _g.sent(),
                        _e.powers = powers,
                        _e.valsetNonce = currentValsetNonce,
                        _e.rewardAmount = 0,
                        _e.rewardToken = pure_1.ZeroAddress,
                        _e);
                    _f = {};
                    return [4 /*yield*/, (0, pure_1.getSignerAddresses)(newValidators)];
                case 4:
                    newValset = (_f.validators = _g.sent(),
                        _f.powers = newPowers,
                        _f.valsetNonce = newValsetNonce,
                        _f.rewardAmount = 0,
                        _f.rewardToken = pure_1.ZeroAddress,
                        _f);
                    if (!opts.badReward) return [3 /*break*/, 5];
                    // some amount of a reward, in a random token that's not in the bridge
                    // should panic because the token doesn't exist
                    newValset.rewardAmount = 5000000;
                    newValset.rewardToken = "0x8bcd7D3532CB626A7138962Bdb859737e5B6d4a7";
                    return [3 /*break*/, 10];
                case 5:
                    if (!opts.withReward) return [3 /*break*/, 7];
                    return [4 /*yield*/, (0, pure_1.parseEvent)(gravity, gravity.deployERC20('uatom', 'Atom', 'ATOM', 6), 1)];
                case 6:
                    eventArgs = _g.sent();
                    newValset.rewardToken = eventArgs._tokenContract;
                    // five atom, issued as an inflationary reward
                    newValset.rewardAmount = 5000000;
                    // connect with the contract to check balances later
                    ERC20contract = new hardhat_1.ethers.Contract(eventArgs._tokenContract, [
                        "function balanceOf(address account) view returns (uint256 balance)"
                    ], gravity.provider);
                    return [3 /*break*/, 10];
                case 7:
                    if (!opts.notEnoughReward) return [3 /*break*/, 10];
                    // send in 1000 tokens, then have a reward of five million
                    return [4 /*yield*/, testERC20.functions.approve(gravity.address, 1000)];
                case 8:
                    // send in 1000 tokens, then have a reward of five million
                    _g.sent();
                    return [4 /*yield*/, gravity.functions.sendToCosmos(testERC20.address, hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"), 1000)];
                case 9:
                    _g.sent();
                    newValset.rewardToken = testERC20.address;
                    newValset.rewardAmount = 5000000;
                    _g.label = 10;
                case 10:
                    checkpoint = (0, pure_1.makeCheckpoint)(newValset.validators, newValset.powers, newValset.valsetNonce, newValset.rewardAmount, newValset.rewardToken, gravityId);
                    return [4 /*yield*/, (0, pure_1.signHash)(validators, checkpoint)];
                case 11:
                    sigs = _g.sent();
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
                    if (opts.malformedCurrentValset) {
                        // Remove one of the powers to make the length not match
                        powers.pop();
                    }
                    return [4 /*yield*/, gravity.updateValset(newValset, currentValset, sigs)];
                case 12:
                    valsetUpdateTx = _g.sent();
                    if (!opts.withReward) return [3 /*break*/, 16];
                    // panic if we failed to deploy the contract earlier
                    expect(ERC20contract);
                    if (!ERC20contract) return [3 /*break*/, 16];
                    _b = expect;
                    _d = (_c = ERC20contract.functions).balanceOf;
                    return [4 /*yield*/, valsetUpdateTx.from];
                case 13: return [4 /*yield*/, _d.apply(_c, [_g.sent()])];
                case 14: return [4 /*yield*/, (_g.sent())[0].toNumber()];
                case 15:
                    _b.apply(void 0, [_g.sent()]).to.equal(5000000);
                    _g.label = 16;
                case 16: return [2 /*return*/, { gravity: gravity, checkpoint: checkpoint }];
            }
        });
    });
}
describe("updateValset tests", function () {
    it("throws on malformed new valset", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ malformedNewValset: true })).to.be.revertedWith("MalformedNewValidatorSet()")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("throws on empty new valset", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ zeroLengthValset: true })).to.be.revertedWith("MalformedNewValidatorSet()")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
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
    it("throws on new valset nonce not incremented", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ nonceNotIncremented: true })).to.be.revertedWith("InvalidValsetNonce(0, 0)")];
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
    it("throws on not enough power in new set", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ notEnoughPowerNewSet: true })).to.be.revertedWith("InsufficientPower(625, 6666)")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("throws on bad reward ", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ badReward: true })).to.be.revertedWith("Address: call to non-contract")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("throws on not enough reward ", function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, expect(runTest({ notEnoughReward: true })).to.be.revertedWith("transfer amount exceeds balance")];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    });
    it("pays reward correctly", function () {
        return __awaiter(this, void 0, void 0, function () {
            var _a, gravity, checkpoint, _b;
            return __generator(this, function (_c) {
                switch (_c.label) {
                    case 0: return [4 /*yield*/, runTest({ withReward: true })];
                    case 1:
                        _a = _c.sent(), gravity = _a.gravity, checkpoint = _a.checkpoint;
                        _b = expect;
                        return [4 /*yield*/, gravity.functions.state_lastValsetCheckpoint()];
                    case 2:
                        _b.apply(void 0, [(_c.sent())[0]]).to.equal(checkpoint);
                        return [2 /*return*/];
                }
            });
        });
    });
    it("happy path", function () {
        return __awaiter(this, void 0, void 0, function () {
            var _a, gravity, checkpoint, _b;
            return __generator(this, function (_c) {
                switch (_c.label) {
                    case 0: return [4 /*yield*/, runTest({})];
                    case 1:
                        _a = _c.sent(), gravity = _a.gravity, checkpoint = _a.checkpoint;
                        _b = expect;
                        return [4 /*yield*/, gravity.functions.state_lastValsetCheckpoint()];
                    case 2:
                        _b.apply(void 0, [(_c.sent())[0]]).to.equal(checkpoint);
                        return [2 /*return*/];
                }
            });
        });
    });
});
// This test produces a hash for the contract which should match what is being used in the Go unit tests. It's here for
// the use of anyone updating the Go tests.
describe("updateValset Go test hash", function () {
    it("produces good hash", function () {
        return __awaiter(this, void 0, void 0, function () {
            var gravityId, methodName, validators, powers, newValset, checkpoint, abiEncodedValset, valsetDigest;
            return __generator(this, function (_a) {
                gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                methodName = hardhat_1.ethers.utils.formatBytes32String("checkpoint");
                validators = ["0xE5904695748fe4A84b40b3fc79De2277660BD1D3",
                    "0xc783df8a850f42e7F7e57013759C285caa701eB6",
                    "0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4",
                ];
                powers = [1431655765, 1431655765, 1431655765];
                newValset = {
                    validators: validators,
                    powers: powers,
                    valsetNonce: 0,
                    rewardAmount: 0,
                    rewardToken: pure_1.ZeroAddress
                };
                checkpoint = (0, pure_1.makeCheckpoint)(newValset.validators, newValset.powers, newValset.valsetNonce, newValset.rewardAmount, newValset.rewardToken, gravityId);
                abiEncodedValset = hardhat_1.ethers.utils.defaultAbiCoder.encode([
                    "bytes32",
                    "bytes32",
                    "uint256",
                    "address[]",
                    "uint256[]",
                    "uint256",
                    "address" // rewardToken
                ], [
                    gravityId,
                    methodName,
                    newValset.valsetNonce,
                    newValset.validators,
                    newValset.powers,
                    newValset.rewardAmount,
                    newValset.rewardToken,
                ]);
                valsetDigest = hardhat_1.ethers.utils.keccak256(abiEncodedValset);
                // these should be equal, otherwise either our abi encoding here
                // or over in test-utils/pure.ts is incorrect
                expect(valsetDigest).equal(checkpoint);
                console.log("elements in Valset digest:", {
                    "gravityId": gravityId,
                    "validators": validators,
                    "powers": powers,
                    "valsetNonce": newValset.valsetNonce,
                    "rewardAmount": newValset.rewardAmount,
                    "rewardToken": newValset.rewardToken
                });
                console.log("abiEncodedValset:", abiEncodedValset);
                console.log("valsetDigest:", valsetDigest);
                return [2 /*return*/];
            });
        });
    });
});
