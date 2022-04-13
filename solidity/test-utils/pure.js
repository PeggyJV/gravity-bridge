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
exports.parseEvent = exports.examplePowers = exports.makeTxBatchHash = exports.signHash = exports.makeCheckpoint = exports.getSignerAddresses = exports.ZeroAddress = void 0;
var hardhat_1 = require("hardhat");
exports.ZeroAddress = "0x0000000000000000000000000000000000000000";
function getSignerAddresses(signers) {
    return __awaiter(this, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, Promise.all(signers.map(function (signer) { return signer.getAddress(); }))];
                case 1: return [2 /*return*/, _a.sent()];
            }
        });
    });
}
exports.getSignerAddresses = getSignerAddresses;
function makeCheckpoint(validators, powers, valsetNonce, rewardAmount, rewardToken, gravityId) {
    var methodName = hardhat_1.ethers.utils.formatBytes32String("checkpoint");
    var abiEncoded = hardhat_1.ethers.utils.defaultAbiCoder.encode(["bytes32", "bytes32", "uint256", "address[]", "uint256[]", "uint256", "address"], [gravityId, methodName, valsetNonce, validators, powers, rewardAmount, rewardToken]);
    var checkpoint = hardhat_1.ethers.utils.keccak256(abiEncoded);
    return checkpoint;
}
exports.makeCheckpoint = makeCheckpoint;
function signHash(signers, hash) {
    return __awaiter(this, void 0, void 0, function () {
        var sigs, i, sig, address, splitSig;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    sigs = [];
                    i = 0;
                    _a.label = 1;
                case 1:
                    if (!(i < signers.length)) return [3 /*break*/, 5];
                    return [4 /*yield*/, signers[i].signMessage(hardhat_1.ethers.utils.arrayify(hash))];
                case 2:
                    sig = _a.sent();
                    return [4 /*yield*/, signers[i].getAddress()];
                case 3:
                    address = _a.sent();
                    splitSig = hardhat_1.ethers.utils.splitSignature(sig);
                    sigs.push({ v: splitSig.v, r: splitSig.r, s: splitSig.s });
                    _a.label = 4;
                case 4:
                    i = i + 1;
                    return [3 /*break*/, 1];
                case 5: return [2 /*return*/, sigs];
            }
        });
    });
}
exports.signHash = signHash;
function makeTxBatchHash(amounts, destinations, fees, nonces, gravityId) {
    var methodName = hardhat_1.ethers.utils.formatBytes32String("transactionBatch");
    var abiEncoded = hardhat_1.ethers.utils.defaultAbiCoder.encode(["bytes32", "bytes32", "uint256[]", "address[]", "uint256[]", "uint256[]"], [gravityId, methodName, amounts, destinations, fees, nonces]);
    // console.log(abiEncoded);
    var txHash = hardhat_1.ethers.utils.keccak256(abiEncoded);
    return txHash;
}
exports.makeTxBatchHash = makeTxBatchHash;
function examplePowers() {
    return [
        707,
        621,
        608,
        439,
        412,
        407,
        319,
        312,
        311,
        303,
        246,
        241,
        224,
        213,
        194,
        175,
        173,
        170,
        154,
        149,
        139,
        123,
        119,
        113,
        110,
        107,
        105,
        104,
        92,
        90,
        88,
        88,
        88,
        85,
        85,
        84,
        82,
        70,
        67,
        64,
        59,
        58,
        56,
        55,
        52,
        52,
        52,
        50,
        49,
        44,
        42,
        40,
        39,
        38,
        37,
        37,
        36,
        35,
        34,
        33,
        33,
        33,
        32,
        31,
        30,
        30,
        29,
        28,
        27,
        26,
        25,
        24,
        23,
        23,
        22,
        22,
        22,
        21,
        21,
        20,
        19,
        18,
        17,
        16,
        14,
        14,
        13,
        13,
        11,
        10,
        10,
        10,
        10,
        10,
        9,
        8,
        8,
        7,
        7,
        7,
        6,
        6,
        5,
        5,
        5,
        5,
        5,
        5,
        4,
        4,
        3,
        2,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1
    ];
}
exports.examplePowers = examplePowers;
function parseEvent(contract, txPromise, eventOrder) {
    return __awaiter(this, void 0, void 0, function () {
        var tx, receipt, args, acc;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, txPromise];
                case 1:
                    tx = _a.sent();
                    return [4 /*yield*/, contract.provider.getTransactionReceipt(tx.hash)];
                case 2:
                    receipt = _a.sent();
                    args = contract.interface.parseLog(receipt.logs[eventOrder]).args;
                    acc = {};
                    args = Object.keys(args).reduce(function (acc, key) {
                        if (Number.isNaN(parseInt(key, 10)) && key !== 'length') {
                            acc[key] = args[key];
                        }
                        return acc;
                    }, acc);
                    return [2 /*return*/, args];
            }
        });
    });
}
exports.parseEvent = parseEvent;
