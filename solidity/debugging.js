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
var stargate_1 = require("@cosmjs/stargate");
var tendermint_rpc_1 = require("@cosmjs/tendermint-rpc");
var query_1 = require("./gen/gravity/v1/query");
var long_1 = require("long");
var process_1 = require("process");
function getQueryService() {
    return __awaiter(this, void 0, void 0, function () {
        var cosmosNode, tendermintClient, queryClient, rpcClient;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    cosmosNode = "http://localhost:26657";
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
function getParams() {
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
                    return [2 /*return*/, res.params];
            }
        });
    });
}
function getValset(signerSetNonce) {
    return __awaiter(this, void 0, void 0, function () {
        var queryService, res;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, getQueryService()];
                case 1:
                    queryService = _a.sent();
                    return [4 /*yield*/, queryService.SignerSetTx({ signerSetNonce: signerSetNonce })];
                case 2:
                    res = _a.sent();
                    if (!res.signerSet) {
                        console.log("Could not retrieve signer set", res);
                        (0, process_1.exit)(1);
                    }
                    return [2 /*return*/, res.signerSet];
            }
        });
    });
}
function getSignerSetTxConfirmations(signerSetNonce) {
    return __awaiter(this, void 0, void 0, function () {
        var queryService, res;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, getQueryService()];
                case 1:
                    queryService = _a.sent();
                    return [4 /*yield*/, queryService.SignerSetTxConfirmations({ signerSetNonce: signerSetNonce })];
                case 2:
                    res = _a.sent();
                    if (!res.signatures) {
                        console.log("Could not retrieve signatures", res);
                        (0, process_1.exit)(1);
                    }
                    return [2 /*return*/, res];
            }
        });
    });
}
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
function getAllValsets() {
    return __awaiter(this, void 0, void 0, function () {
        var queryService, res;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, getQueryService()];
                case 1:
                    queryService = _a.sent();
                    return [4 /*yield*/, queryService.SignerSetTxs({})];
                case 2:
                    res = _a.sent();
                    return [2 /*return*/, res];
            }
        });
    });
}
function getDelegateKeys() {
    return __awaiter(this, void 0, void 0, function () {
        var queryService, res;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, getQueryService()];
                case 1:
                    queryService = _a.sent();
                    return [4 /*yield*/, queryService.DelegateKeysByOrchestrator({
                            orchestratorAddress: "cosmos14uvqun482ydhljwtvacy5grvgh23xrmgymg0wd"
                        })];
                case 2:
                    res = _a.sent();
                    return [2 /*return*/, res];
            }
        });
    });
}
(function () {
    return __awaiter(this, void 0, void 0, function () {
        var res;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, getValset(long_1["default"].fromInt(1))];
                case 1:
                    res = _a.sent();
                    console.log(JSON.stringify(res));
                    return [2 /*return*/];
            }
        });
    });
})();
