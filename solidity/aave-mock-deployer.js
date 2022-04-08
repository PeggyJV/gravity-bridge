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
var ethers_1 = require("ethers");
var fs_1 = require("fs");
var command_line_args_1 = require("command-line-args");
var assert_1 = require("assert");
var args = (0, command_line_args_1["default"])([
    // the ethernum node used to deploy the contract
    { name: "eth-node", type: String },
    // the Ethereum private key that will contain the gas required to pay for the contact deployment
    { name: "eth-privkey", type: String },
    // the cellar contract .json file
    { name: "contract", type: String },
    // gravity contract address to transfer ownership to
    { name: "gravity-address", type: String }
]);
function deploy() {
    return __awaiter(this, void 0, void 0, function () {
        var gravityAddress, provider, wallet, _a, bytecode, abi, cellarFactory, cellar;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0:
                    gravityAddress = args["gravity-address"];
                    (0, assert_1.doesNotThrow)(function () { return ethers_1.ethers.utils.getAddress(gravityAddress); }, "the provided gravity address is invalid.");
                    return [4 /*yield*/, new ethers_1.ethers.providers.JsonRpcProvider(args["eth-node"])];
                case 1:
                    provider = _b.sent();
                    wallet = new ethers_1.ethers.Wallet(args["eth-privkey"], provider);
                    _a = getContractArtifacts(args["contract"]), bytecode = _a.bytecode, abi = _a.abi;
                    cellarFactory = new ethers_1.ethers.ContractFactory(abi, bytecode, wallet);
                    return [4 /*yield*/, cellarFactory.deploy()];
                case 2:
                    cellar = (_b.sent());
                    return [4 /*yield*/, cellar.deployed()];
                case 3:
                    _b.sent();
                    return [4 /*yield*/, cellar.transferOwnership(args["gravity-address"])];
                case 4:
                    _b.sent();
                    return [2 /*return*/];
            }
        });
    });
}
function getContractArtifacts(path) {
    var _a = JSON.parse(fs_1["default"].readFileSync(path, "utf8").toString()), bytecode = _a.bytecode, abi = _a.abi;
    return { bytecode: bytecode, abi: abi };
}
deploy();
