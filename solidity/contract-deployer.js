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
var process_1 = require("process");
var stargate_1 = require("@cosmjs/stargate");
var tendermint_rpc_1 = require("@cosmjs/tendermint-rpc");
var query_1 = require("./gen/gravity/v1/query");
var args = (0, command_line_args_1["default"])([
    // the ethernum node used to deploy the contract
    { name: "eth-node", type: String },
    // the cosmos node that will be used to grab the validator set via RPC (TODO),
    { name: "cosmos-node", type: String },
    // the Ethereum private key that will contain the gas required to pay for the contact deployment
    { name: "eth-privkey", type: String },
    // the gravity contract .json file
    { name: "contract", type: String },
    // test mode, if enabled this script deploys three ERC20 contracts for testing
    { name: "test-mode", type: String },
]);
function getQueryService() {
    return __awaiter(this, void 0, void 0, function () {
        var cosmosNode, tendermintClient, queryClient, rpcClient;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    cosmosNode = args["cosmos-node"];
                    return [4 /*yield*/, tendermint_rpc_1.Tendermint34Client.connect(cosmosNode)];
                case 1:
                    tendermintClient = _a.sent();
                    queryClient = new stargate_1.QueryClient(tendermintClient);
                    rpcClient = (0, stargate_1.createProtobufRpcClient)(queryClient);
                    return [2 /*return*/, new query_1.QueryClientImpl(rpcClient)];
            }
        });
    });
}
function deploy() {
    return __awaiter(this, void 0, void 0, function () {
        var startTime, provider, wallet, success, present, timeDiff, erc20_a_path, erc20_b_path, erc20_c_path, main_location_a, main_location_b, main_location_c, alt_location_1_a, alt_location_1_b, alt_location_1_c, alt_location_2_a, alt_location_2_b, alt_location_2_c, solidity_dir_a, solidity_dir_b, solidity_dir_c, docker_location_a, docker_location_b, docker_location_c, _a, abi_1, bytecode_1, erc20Factory, testERC20, erc20TestAddress, _b, abi1, bytecode1, erc20Factory1, testERC201, erc20TestAddress1, _c, abi2, bytecode2, erc20Factory2, testERC202, erc20TestAddress2, gravityIdString, gravityId, _d, abi, bytecode, factory, latestValset, eth_addresses, powers, powers_sum, i, vote_power, gravity;
        return __generator(this, function (_e) {
            switch (_e.label) {
                case 0:
                    startTime = new Date();
                    return [4 /*yield*/, new ethers_1.ethers.providers.JsonRpcProvider(args["eth-node"])];
                case 1:
                    provider = _e.sent();
                    wallet = new ethers_1.ethers.Wallet(args["eth-privkey"], provider);
                    if (!(args["test-mode"] == "True" || args["test-mode"] == "true")) return [3 /*break*/, 4];
                    success = false;
                    _e.label = 2;
                case 2:
                    if (!!success) return [3 /*break*/, 4];
                    present = new Date();
                    timeDiff = present.getTime() - startTime.getTime();
                    timeDiff = timeDiff / 1000;
                    provider
                        .getBlockNumber()
                        .then(function (_) { return (success = true); })["catch"](function (_) { return console.log("Ethereum RPC error, trying again"); });
                    if (timeDiff > 600) {
                        console.log("Could not contact Ethereum RPC after 10 minutes, check the URL!");
                        (0, process_1.exit)(1);
                    }
                    return [4 /*yield*/, sleep(1000)];
                case 3:
                    _e.sent();
                    return [3 /*break*/, 2];
                case 4:
                    if (!(args["test-mode"] == "True" || args["test-mode"] == "true")) return [3 /*break*/, 11];
                    console.log("Test mode, deploying ERC20 contracts");
                    main_location_a = "/gravity/solidity/artifacts/contracts/TestERC20A.sol/TestERC20A.json";
                    main_location_b = "/gravity/solidity/artifacts/contracts/TestERC20B.sol/TestERC20B.json";
                    main_location_c = "/gravity/solidity/artifacts/contracts/TestERC20C.sol/TestERC20C.json";
                    alt_location_1_a = "/solidity/TestERC20A.json";
                    alt_location_1_b = "/solidity/TestERC20B.json";
                    alt_location_1_c = "/solidity/TestERC20C.json";
                    alt_location_2_a = "TestERC20A.json";
                    alt_location_2_b = "TestERC20B.json";
                    alt_location_2_c = "TestERC20C.json";
                    solidity_dir_a = "artifacts/contracts/TestERC20A.sol/TestERC20A.json";
                    solidity_dir_b = "artifacts/contracts/TestERC20B.sol/TestERC20B.json";
                    solidity_dir_c = "artifacts/contracts/TestERC20C.sol/TestERC20C.json";
                    docker_location_a = "/artifacts/contracts/TestERC20A.sol/TestERC20A.json";
                    docker_location_b = "/artifacts/contracts/TestERC20B.sol/TestERC20B.json";
                    docker_location_c = "/artifacts/contracts/TestERC20C.sol/TestERC20C.json";
                    if (fs_1["default"].existsSync(main_location_a)) {
                        erc20_a_path = main_location_a;
                        erc20_b_path = main_location_b;
                        erc20_c_path = main_location_c;
                    }
                    else if (fs_1["default"].existsSync(alt_location_1_a)) {
                        erc20_a_path = alt_location_1_a;
                        erc20_b_path = alt_location_1_b;
                        erc20_c_path = alt_location_1_c;
                    }
                    else if (fs_1["default"].existsSync(alt_location_2_a)) {
                        erc20_a_path = alt_location_2_a;
                        erc20_b_path = alt_location_2_b;
                        erc20_c_path = alt_location_2_c;
                    }
                    else if (fs_1["default"].existsSync(solidity_dir_a)) {
                        erc20_a_path = solidity_dir_a;
                        erc20_b_path = solidity_dir_b;
                        erc20_c_path = solidity_dir_c;
                    }
                    else if (fs_1["default"].existsSync(docker_location_a)) {
                        erc20_a_path = docker_location_a;
                        erc20_b_path = docker_location_b;
                        erc20_c_path = docker_location_c;
                    }
                    else {
                        console.log("Test mode was enabled but the ERC20 contracts can't be found!");
                        (0, process_1.exit)(1);
                    }
                    _a = getContractArtifacts(erc20_a_path), abi_1 = _a.abi, bytecode_1 = _a.bytecode;
                    erc20Factory = new ethers_1.ethers.ContractFactory(abi_1, bytecode_1, wallet);
                    return [4 /*yield*/, erc20Factory.deploy()];
                case 5:
                    testERC20 = (_e.sent());
                    return [4 /*yield*/, testERC20.deployed()];
                case 6:
                    _e.sent();
                    erc20TestAddress = testERC20.address;
                    console.log("ERC20 deployed at Address - ", erc20TestAddress);
                    _b = getContractArtifacts(erc20_b_path), abi1 = _b.abi, bytecode1 = _b.bytecode;
                    erc20Factory1 = new ethers_1.ethers.ContractFactory(abi1, bytecode1, wallet);
                    return [4 /*yield*/, erc20Factory1.deploy()];
                case 7:
                    testERC201 = (_e.sent());
                    return [4 /*yield*/, testERC201.deployed()];
                case 8:
                    _e.sent();
                    erc20TestAddress1 = testERC201.address;
                    console.log("ERC20 deployed at Address - ", erc20TestAddress1);
                    _c = getContractArtifacts(erc20_c_path), abi2 = _c.abi, bytecode2 = _c.bytecode;
                    erc20Factory2 = new ethers_1.ethers.ContractFactory(abi2, bytecode2, wallet);
                    return [4 /*yield*/, erc20Factory2.deploy()];
                case 9:
                    testERC202 = (_e.sent());
                    return [4 /*yield*/, testERC202.deployed()];
                case 10:
                    _e.sent();
                    erc20TestAddress2 = testERC202.address;
                    console.log("ERC20 deployed at Address - ", erc20TestAddress2);
                    _e.label = 11;
                case 11: return [4 /*yield*/, getGravityId()];
                case 12:
                    gravityIdString = _e.sent();
                    gravityId = ethers_1.ethers.utils.formatBytes32String(gravityIdString);
                    console.log("Starting Gravity contract deploy");
                    _d = getContractArtifacts(args["contract"]), abi = _d.abi, bytecode = _d.bytecode;
                    factory = new ethers_1.ethers.ContractFactory(abi, bytecode, wallet);
                    console.log("About to get latest Gravity valset");
                    return [4 /*yield*/, getLatestValset()];
                case 13:
                    latestValset = _e.sent();
                    eth_addresses = [];
                    powers = [];
                    powers_sum = 0;
                    // this MUST be sorted uniformly across all components of Gravity in this
                    // case we perform the sorting in module/x/gravity/keeper/types.go to the
                    // output of the endpoint should always be sorted correctly. If you're
                    // having strange problems with updating the validator set you should go
                    // look there.
                    for (i = 0; i < latestValset.signers.length; i++) {
                        if (latestValset.signers[i].ethereumAddress == "") {
                            continue;
                        }
                        eth_addresses.push(latestValset.signers[i].ethereumAddress);
                        powers.push(latestValset.signers[i].power.toNumber());
                        powers_sum += latestValset.signers[i].power.toNumber();
                    }
                    vote_power = 2834678415;
                    if (powers_sum < vote_power) {
                        console.log("Refusing to deploy! Incorrect power! Please inspect the validator set below");
                        console.log("If less than 66% of the current voting power has unset Ethereum Addresses we refuse to deploy");
                        console.log(latestValset);
                        (0, process_1.exit)(1);
                    }
                    return [4 /*yield*/, factory.deploy(
                        // todo generate this randomly at deployment time that way we can avoid
                        // anything but intentional conflicts
                        gravityId, vote_power, eth_addresses, powers)];
                case 14:
                    gravity = (_e.sent());
                    return [4 /*yield*/, gravity.deployed()];
                case 15:
                    _e.sent();
                    console.log("Gravity deployed at Address - ", gravity.address);
                    return [4 /*yield*/, submitGravityAddress(gravity.address)];
                case 16:
                    _e.sent();
                    return [2 /*return*/];
            }
        });
    });
}
function getContractArtifacts(path) {
    var _a = JSON.parse(fs_1["default"].readFileSync(path, "utf8").toString()), bytecode = _a.bytecode, abi = _a.abi;
    return { bytecode: bytecode, abi: abi };
}
var decode = function (str) {
    return Buffer.from(str, "base64").toString("binary");
};
function getLatestValset() {
    return __awaiter(this, void 0, void 0, function () {
        var queryService, res;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, getQueryService()];
                case 1:
                    queryService = _a.sent();
                    return [4 /*yield*/, queryService.LatestSignerSetTx({})];
                case 2:
                    res = _a.sent();
                    if (!res.signerSet) {
                        console.log("Could not retrieve signer set");
                        (0, process_1.exit)(1);
                    }
                    return [2 /*return*/, res.signerSet];
            }
        });
    });
}
function getGravityId() {
    return __awaiter(this, void 0, void 0, function () {
        var queryService, res;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, getQueryService()];
                case 1:
                    queryService = _a.sent();
                    return [4 /*yield*/, queryService.Params({})];
                case 2:
                    res = _a.sent();
                    if (!res.params) {
                        console.log("Could not retrieve params");
                        (0, process_1.exit)(1);
                    }
                    return [2 /*return*/, res.params.gravityId];
            }
        });
    });
}
function submitGravityAddress(address) {
    return __awaiter(this, void 0, void 0, function () { return __generator(this, function (_a) {
        return [2 /*return*/];
    }); });
}
function main() {
    return __awaiter(this, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, deploy()];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    });
}
function sleep(ms) {
    return new Promise(function (resolve) { return setTimeout(resolve, ms); });
}
main();
