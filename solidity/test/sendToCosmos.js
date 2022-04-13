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
        var signers, gravityId, powers, validators, powerThreshold, _a, gravity, testERC20, deployCheckpoint, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l;
        return __generator(this, function (_m) {
            switch (_m.label) {
                case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                case 1:
                    signers = _m.sent();
                    gravityId = hardhat_1.ethers.utils.formatBytes32String("foo");
                    powers = (0, pure_1.examplePowers)();
                    validators = signers.slice(0, powers.length);
                    powerThreshold = 6666;
                    return [4 /*yield*/, (0, test_utils_1.deployContracts)(gravityId, validators, powers, powerThreshold)];
                case 2:
                    _a = _m.sent(), gravity = _a.gravity, testERC20 = _a.testERC20, deployCheckpoint = _a.checkpoint;
                    // Transfer out to Cosmos, locking coins
                    // =====================================
                    return [4 /*yield*/, testERC20.functions.approve(gravity.address, 1000)];
                case 3:
                    // Transfer out to Cosmos, locking coins
                    // =====================================
                    _m.sent();
                    _c = (_b = expect(gravity.functions.sendToCosmos(testERC20.address, hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"), 1000)).to.emit(gravity, 'SendToCosmosEvent')).withArgs;
                    _d = [testERC20.address];
                    return [4 /*yield*/, signers[0].getAddress()];
                case 4: return [4 /*yield*/, _c.apply(_b, _d.concat([_m.sent(), hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"),
                        1000,
                        2]))];
                case 5:
                    _m.sent();
                    _e = expect;
                    return [4 /*yield*/, testERC20.functions.balanceOf(gravity.address)];
                case 6:
                    _e.apply(void 0, [(_m.sent())[0]]).to.equal(1000);
                    _f = expect;
                    return [4 /*yield*/, gravity.functions.state_lastEventNonce()];
                case 7:
                    _f.apply(void 0, [(_m.sent())[0]]).to.equal(2);
                    // Do it again
                    // =====================================
                    return [4 /*yield*/, testERC20.functions.approve(gravity.address, 1000)];
                case 8:
                    // Do it again
                    // =====================================
                    _m.sent();
                    _h = (_g = expect(gravity.functions.sendToCosmos(testERC20.address, hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"), 1000)).to.emit(gravity, 'SendToCosmosEvent')).withArgs;
                    _j = [testERC20.address];
                    return [4 /*yield*/, signers[0].getAddress()];
                case 9: return [4 /*yield*/, _h.apply(_g, _j.concat([_m.sent(), hardhat_1.ethers.utils.formatBytes32String("myCosmosAddress"),
                        1000,
                        3]))];
                case 10:
                    _m.sent();
                    _k = expect;
                    return [4 /*yield*/, testERC20.functions.balanceOf(gravity.address)];
                case 11:
                    _k.apply(void 0, [(_m.sent())[0]]).to.equal(2000);
                    _l = expect;
                    return [4 /*yield*/, gravity.functions.state_lastEventNonce()];
                case 12:
                    _l.apply(void 0, [(_m.sent())[0]]).to.equal(3);
                    return [2 /*return*/];
            }
        });
    });
}
describe("sendToCosmos tests", function () {
    it("works right", function () {
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
