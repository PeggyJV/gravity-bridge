"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
exports.__esModule = true;
exports.QueryClientImpl = exports.UnbatchedSendToEthereumsResponse = exports.UnbatchedSendToEthereumsRequest = exports.BatchedSendToEthereumsResponse = exports.BatchedSendToEthereumsRequest = exports.DelegateKeysByOrchestratorResponse = exports.DelegateKeysByOrchestratorRequest = exports.DelegateKeysByEthereumSignerResponse = exports.DelegateKeysByEthereumSignerRequest = exports.DelegateKeysByValidatorResponse = exports.DelegateKeysByValidatorRequest = exports.DenomToERC20Response = exports.DenomToERC20Request = exports.ERC20ToDenomResponse = exports.ERC20ToDenomRequest = exports.LastSubmittedEthereumEventResponse = exports.LastSubmittedEthereumEventRequest = exports.BatchTxConfirmationsResponse = exports.BatchTxConfirmationsRequest = exports.ContractCallTxConfirmationsResponse = exports.ContractCallTxConfirmationsRequest = exports.BatchTxFeesResponse = exports.BatchTxFeesRequest = exports.UnsignedContractCallTxsResponse = exports.UnsignedContractCallTxsRequest = exports.UnsignedBatchTxsResponse = exports.UnsignedBatchTxsRequest = exports.UnsignedSignerSetTxsResponse = exports.UnsignedSignerSetTxsRequest = exports.ContractCallTxsResponse = exports.ContractCallTxsRequest = exports.BatchTxsResponse = exports.BatchTxsRequest = exports.SignerSetTxsResponse = exports.SignerSetTxsRequest = exports.SignerSetTxConfirmationsResponse = exports.SignerSetTxConfirmationsRequest = exports.ContractCallTxResponse = exports.ContractCallTxRequest = exports.BatchTxResponse = exports.BatchTxRequest = exports.SignerSetTxResponse = exports.LatestSignerSetTxRequest = exports.SignerSetTxRequest = exports.ParamsResponse = exports.ParamsRequest = exports.protobufPackage = void 0;
/* eslint-disable */
var long_1 = require("long");
var minimal_1 = require("protobufjs/minimal");
var genesis_1 = require("../../gravity/v1/genesis");
var gravity_1 = require("../../gravity/v1/gravity");
var pagination_1 = require("../../cosmos/base/query/v1beta1/pagination");
var msgs_1 = require("../../gravity/v1/msgs");
var coin_1 = require("../../cosmos/base/v1beta1/coin");
exports.protobufPackage = "gravity.v1";
var baseParamsRequest = {};
exports.ParamsRequest = {
    encode: function (_, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseParamsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (_) {
        var message = __assign({}, baseParamsRequest);
        return message;
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
    fromPartial: function (_) {
        var message = __assign({}, baseParamsRequest);
        return message;
    }
};
var baseParamsResponse = {};
exports.ParamsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.params !== undefined) {
            genesis_1.Params.encode(message.params, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseParamsResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.params = genesis_1.Params.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseParamsResponse);
        if (object.params !== undefined && object.params !== null) {
            message.params = genesis_1.Params.fromJSON(object.params);
        }
        else {
            message.params = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.params !== undefined &&
            (obj.params = message.params ? genesis_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseParamsResponse);
        if (object.params !== undefined && object.params !== null) {
            message.params = genesis_1.Params.fromPartial(object.params);
        }
        else {
            message.params = undefined;
        }
        return message;
    }
};
var baseSignerSetTxRequest = { signerSetNonce: long_1["default"].UZERO };
exports.SignerSetTxRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.signerSetNonce.isZero()) {
            writer.uint32(8).uint64(message.signerSetNonce);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSignerSetTxRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signerSetNonce = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSignerSetTxRequest);
        if (object.signerSetNonce !== undefined && object.signerSetNonce !== null) {
            message.signerSetNonce = long_1["default"].fromString(object.signerSetNonce);
        }
        else {
            message.signerSetNonce = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.signerSetNonce !== undefined &&
            (obj.signerSetNonce = (message.signerSetNonce || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSignerSetTxRequest);
        if (object.signerSetNonce !== undefined && object.signerSetNonce !== null) {
            message.signerSetNonce = object.signerSetNonce;
        }
        else {
            message.signerSetNonce = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseLatestSignerSetTxRequest = {};
exports.LatestSignerSetTxRequest = {
    encode: function (_, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseLatestSignerSetTxRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (_) {
        var message = __assign({}, baseLatestSignerSetTxRequest);
        return message;
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
    fromPartial: function (_) {
        var message = __assign({}, baseLatestSignerSetTxRequest);
        return message;
    }
};
var baseSignerSetTxResponse = {};
exports.SignerSetTxResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.signerSet !== undefined) {
            gravity_1.SignerSetTx.encode(message.signerSet, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSignerSetTxResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signerSet = gravity_1.SignerSetTx.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSignerSetTxResponse);
        if (object.signerSet !== undefined && object.signerSet !== null) {
            message.signerSet = gravity_1.SignerSetTx.fromJSON(object.signerSet);
        }
        else {
            message.signerSet = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.signerSet !== undefined &&
            (obj.signerSet = message.signerSet
                ? gravity_1.SignerSetTx.toJSON(message.signerSet)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSignerSetTxResponse);
        if (object.signerSet !== undefined && object.signerSet !== null) {
            message.signerSet = gravity_1.SignerSetTx.fromPartial(object.signerSet);
        }
        else {
            message.signerSet = undefined;
        }
        return message;
    }
};
var baseBatchTxRequest = {
    tokenContract: "",
    batchNonce: long_1["default"].UZERO
};
exports.BatchTxRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.tokenContract !== "") {
            writer.uint32(10).string(message.tokenContract);
        }
        if (!message.batchNonce.isZero()) {
            writer.uint32(16).uint64(message.batchNonce);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTxRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tokenContract = reader.string();
                    break;
                case 2:
                    message.batchNonce = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchTxRequest);
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = String(object.tokenContract);
        }
        else {
            message.tokenContract = "";
        }
        if (object.batchNonce !== undefined && object.batchNonce !== null) {
            message.batchNonce = long_1["default"].fromString(object.batchNonce);
        }
        else {
            message.batchNonce = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.tokenContract !== undefined &&
            (obj.tokenContract = message.tokenContract);
        message.batchNonce !== undefined &&
            (obj.batchNonce = (message.batchNonce || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchTxRequest);
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = object.tokenContract;
        }
        else {
            message.tokenContract = "";
        }
        if (object.batchNonce !== undefined && object.batchNonce !== null) {
            message.batchNonce = object.batchNonce;
        }
        else {
            message.batchNonce = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseBatchTxResponse = {};
exports.BatchTxResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.batch !== undefined) {
            gravity_1.BatchTx.encode(message.batch, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTxResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.batch = gravity_1.BatchTx.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchTxResponse);
        if (object.batch !== undefined && object.batch !== null) {
            message.batch = gravity_1.BatchTx.fromJSON(object.batch);
        }
        else {
            message.batch = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.batch !== undefined &&
            (obj.batch = message.batch ? gravity_1.BatchTx.toJSON(message.batch) : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchTxResponse);
        if (object.batch !== undefined && object.batch !== null) {
            message.batch = gravity_1.BatchTx.fromPartial(object.batch);
        }
        else {
            message.batch = undefined;
        }
        return message;
    }
};
var baseContractCallTxRequest = { invalidationNonce: long_1["default"].UZERO };
exports.ContractCallTxRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.invalidationScope.length !== 0) {
            writer.uint32(10).bytes(message.invalidationScope);
        }
        if (!message.invalidationNonce.isZero()) {
            writer.uint32(16).uint64(message.invalidationNonce);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseContractCallTxRequest);
        message.invalidationScope = new Uint8Array();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.invalidationScope = reader.bytes();
                    break;
                case 2:
                    message.invalidationNonce = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseContractCallTxRequest);
        message.invalidationScope = new Uint8Array();
        if (object.invalidationScope !== undefined &&
            object.invalidationScope !== null) {
            message.invalidationScope = bytesFromBase64(object.invalidationScope);
        }
        if (object.invalidationNonce !== undefined &&
            object.invalidationNonce !== null) {
            message.invalidationNonce = long_1["default"].fromString(object.invalidationNonce);
        }
        else {
            message.invalidationNonce = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.invalidationScope !== undefined &&
            (obj.invalidationScope = base64FromBytes(message.invalidationScope !== undefined
                ? message.invalidationScope
                : new Uint8Array()));
        message.invalidationNonce !== undefined &&
            (obj.invalidationNonce = (message.invalidationNonce || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseContractCallTxRequest);
        if (object.invalidationScope !== undefined &&
            object.invalidationScope !== null) {
            message.invalidationScope = object.invalidationScope;
        }
        else {
            message.invalidationScope = new Uint8Array();
        }
        if (object.invalidationNonce !== undefined &&
            object.invalidationNonce !== null) {
            message.invalidationNonce = object.invalidationNonce;
        }
        else {
            message.invalidationNonce = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseContractCallTxResponse = {};
exports.ContractCallTxResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.logicCall !== undefined) {
            gravity_1.ContractCallTx.encode(message.logicCall, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseContractCallTxResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.logicCall = gravity_1.ContractCallTx.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseContractCallTxResponse);
        if (object.logicCall !== undefined && object.logicCall !== null) {
            message.logicCall = gravity_1.ContractCallTx.fromJSON(object.logicCall);
        }
        else {
            message.logicCall = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.logicCall !== undefined &&
            (obj.logicCall = message.logicCall
                ? gravity_1.ContractCallTx.toJSON(message.logicCall)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseContractCallTxResponse);
        if (object.logicCall !== undefined && object.logicCall !== null) {
            message.logicCall = gravity_1.ContractCallTx.fromPartial(object.logicCall);
        }
        else {
            message.logicCall = undefined;
        }
        return message;
    }
};
var baseSignerSetTxConfirmationsRequest = {
    signerSetNonce: long_1["default"].UZERO
};
exports.SignerSetTxConfirmationsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.signerSetNonce.isZero()) {
            writer.uint32(8).uint64(message.signerSetNonce);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSignerSetTxConfirmationsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signerSetNonce = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSignerSetTxConfirmationsRequest);
        if (object.signerSetNonce !== undefined && object.signerSetNonce !== null) {
            message.signerSetNonce = long_1["default"].fromString(object.signerSetNonce);
        }
        else {
            message.signerSetNonce = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.signerSetNonce !== undefined &&
            (obj.signerSetNonce = (message.signerSetNonce || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSignerSetTxConfirmationsRequest);
        if (object.signerSetNonce !== undefined && object.signerSetNonce !== null) {
            message.signerSetNonce = object.signerSetNonce;
        }
        else {
            message.signerSetNonce = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseSignerSetTxConfirmationsResponse = {};
exports.SignerSetTxConfirmationsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.signatures; _i < _a.length; _i++) {
            var v = _a[_i];
            msgs_1.SignerSetTxConfirmation.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSignerSetTxConfirmationsResponse);
        message.signatures = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signatures.push(msgs_1.SignerSetTxConfirmation.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSignerSetTxConfirmationsResponse);
        message.signatures = [];
        if (object.signatures !== undefined && object.signatures !== null) {
            for (var _i = 0, _a = object.signatures; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signatures.push(msgs_1.SignerSetTxConfirmation.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.signatures) {
            obj.signatures = message.signatures.map(function (e) {
                return e ? msgs_1.SignerSetTxConfirmation.toJSON(e) : undefined;
            });
        }
        else {
            obj.signatures = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSignerSetTxConfirmationsResponse);
        message.signatures = [];
        if (object.signatures !== undefined && object.signatures !== null) {
            for (var _i = 0, _a = object.signatures; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signatures.push(msgs_1.SignerSetTxConfirmation.fromPartial(e));
            }
        }
        return message;
    }
};
var baseSignerSetTxsRequest = {};
exports.SignerSetTxsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSignerSetTxsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSignerSetTxsRequest);
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? pagination_1.PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSignerSetTxsRequest);
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    }
};
var baseSignerSetTxsResponse = {};
exports.SignerSetTxsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.signerSets; _i < _a.length; _i++) {
            var v = _a[_i];
            gravity_1.SignerSetTx.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSignerSetTxsResponse);
        message.signerSets = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signerSets.push(gravity_1.SignerSetTx.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSignerSetTxsResponse);
        message.signerSets = [];
        if (object.signerSets !== undefined && object.signerSets !== null) {
            for (var _i = 0, _a = object.signerSets; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signerSets.push(gravity_1.SignerSetTx.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.signerSets) {
            obj.signerSets = message.signerSets.map(function (e) {
                return e ? gravity_1.SignerSetTx.toJSON(e) : undefined;
            });
        }
        else {
            obj.signerSets = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? pagination_1.PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSignerSetTxsResponse);
        message.signerSets = [];
        if (object.signerSets !== undefined && object.signerSets !== null) {
            for (var _i = 0, _a = object.signerSets; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signerSets.push(gravity_1.SignerSetTx.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    }
};
var baseBatchTxsRequest = {};
exports.BatchTxsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTxsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchTxsRequest);
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? pagination_1.PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchTxsRequest);
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    }
};
var baseBatchTxsResponse = {};
exports.BatchTxsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.batches; _i < _a.length; _i++) {
            var v = _a[_i];
            gravity_1.BatchTx.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTxsResponse);
        message.batches = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.batches.push(gravity_1.BatchTx.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchTxsResponse);
        message.batches = [];
        if (object.batches !== undefined && object.batches !== null) {
            for (var _i = 0, _a = object.batches; _i < _a.length; _i++) {
                var e = _a[_i];
                message.batches.push(gravity_1.BatchTx.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.batches) {
            obj.batches = message.batches.map(function (e) {
                return e ? gravity_1.BatchTx.toJSON(e) : undefined;
            });
        }
        else {
            obj.batches = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? pagination_1.PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchTxsResponse);
        message.batches = [];
        if (object.batches !== undefined && object.batches !== null) {
            for (var _i = 0, _a = object.batches; _i < _a.length; _i++) {
                var e = _a[_i];
                message.batches.push(gravity_1.BatchTx.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    }
};
var baseContractCallTxsRequest = {};
exports.ContractCallTxsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseContractCallTxsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseContractCallTxsRequest);
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? pagination_1.PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseContractCallTxsRequest);
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    }
};
var baseContractCallTxsResponse = {};
exports.ContractCallTxsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.calls; _i < _a.length; _i++) {
            var v = _a[_i];
            gravity_1.ContractCallTx.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseContractCallTxsResponse);
        message.calls = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.calls.push(gravity_1.ContractCallTx.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseContractCallTxsResponse);
        message.calls = [];
        if (object.calls !== undefined && object.calls !== null) {
            for (var _i = 0, _a = object.calls; _i < _a.length; _i++) {
                var e = _a[_i];
                message.calls.push(gravity_1.ContractCallTx.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.calls) {
            obj.calls = message.calls.map(function (e) {
                return e ? gravity_1.ContractCallTx.toJSON(e) : undefined;
            });
        }
        else {
            obj.calls = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? pagination_1.PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseContractCallTxsResponse);
        message.calls = [];
        if (object.calls !== undefined && object.calls !== null) {
            for (var _i = 0, _a = object.calls; _i < _a.length; _i++) {
                var e = _a[_i];
                message.calls.push(gravity_1.ContractCallTx.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    }
};
var baseUnsignedSignerSetTxsRequest = { address: "" };
exports.UnsignedSignerSetTxsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseUnsignedSignerSetTxsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseUnsignedSignerSetTxsRequest);
        if (object.address !== undefined && object.address !== null) {
            message.address = String(object.address);
        }
        else {
            message.address = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.address !== undefined && (obj.address = message.address);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseUnsignedSignerSetTxsRequest);
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = "";
        }
        return message;
    }
};
var baseUnsignedSignerSetTxsResponse = {};
exports.UnsignedSignerSetTxsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.signerSets; _i < _a.length; _i++) {
            var v = _a[_i];
            gravity_1.SignerSetTx.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseUnsignedSignerSetTxsResponse);
        message.signerSets = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signerSets.push(gravity_1.SignerSetTx.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseUnsignedSignerSetTxsResponse);
        message.signerSets = [];
        if (object.signerSets !== undefined && object.signerSets !== null) {
            for (var _i = 0, _a = object.signerSets; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signerSets.push(gravity_1.SignerSetTx.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.signerSets) {
            obj.signerSets = message.signerSets.map(function (e) {
                return e ? gravity_1.SignerSetTx.toJSON(e) : undefined;
            });
        }
        else {
            obj.signerSets = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseUnsignedSignerSetTxsResponse);
        message.signerSets = [];
        if (object.signerSets !== undefined && object.signerSets !== null) {
            for (var _i = 0, _a = object.signerSets; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signerSets.push(gravity_1.SignerSetTx.fromPartial(e));
            }
        }
        return message;
    }
};
var baseUnsignedBatchTxsRequest = { address: "" };
exports.UnsignedBatchTxsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseUnsignedBatchTxsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseUnsignedBatchTxsRequest);
        if (object.address !== undefined && object.address !== null) {
            message.address = String(object.address);
        }
        else {
            message.address = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.address !== undefined && (obj.address = message.address);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseUnsignedBatchTxsRequest);
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = "";
        }
        return message;
    }
};
var baseUnsignedBatchTxsResponse = {};
exports.UnsignedBatchTxsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.batches; _i < _a.length; _i++) {
            var v = _a[_i];
            gravity_1.BatchTx.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseUnsignedBatchTxsResponse);
        message.batches = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.batches.push(gravity_1.BatchTx.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseUnsignedBatchTxsResponse);
        message.batches = [];
        if (object.batches !== undefined && object.batches !== null) {
            for (var _i = 0, _a = object.batches; _i < _a.length; _i++) {
                var e = _a[_i];
                message.batches.push(gravity_1.BatchTx.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.batches) {
            obj.batches = message.batches.map(function (e) {
                return e ? gravity_1.BatchTx.toJSON(e) : undefined;
            });
        }
        else {
            obj.batches = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseUnsignedBatchTxsResponse);
        message.batches = [];
        if (object.batches !== undefined && object.batches !== null) {
            for (var _i = 0, _a = object.batches; _i < _a.length; _i++) {
                var e = _a[_i];
                message.batches.push(gravity_1.BatchTx.fromPartial(e));
            }
        }
        return message;
    }
};
var baseUnsignedContractCallTxsRequest = { address: "" };
exports.UnsignedContractCallTxsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseUnsignedContractCallTxsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseUnsignedContractCallTxsRequest);
        if (object.address !== undefined && object.address !== null) {
            message.address = String(object.address);
        }
        else {
            message.address = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.address !== undefined && (obj.address = message.address);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseUnsignedContractCallTxsRequest);
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = "";
        }
        return message;
    }
};
var baseUnsignedContractCallTxsResponse = {};
exports.UnsignedContractCallTxsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.calls; _i < _a.length; _i++) {
            var v = _a[_i];
            gravity_1.ContractCallTx.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseUnsignedContractCallTxsResponse);
        message.calls = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.calls.push(gravity_1.ContractCallTx.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseUnsignedContractCallTxsResponse);
        message.calls = [];
        if (object.calls !== undefined && object.calls !== null) {
            for (var _i = 0, _a = object.calls; _i < _a.length; _i++) {
                var e = _a[_i];
                message.calls.push(gravity_1.ContractCallTx.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.calls) {
            obj.calls = message.calls.map(function (e) {
                return e ? gravity_1.ContractCallTx.toJSON(e) : undefined;
            });
        }
        else {
            obj.calls = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseUnsignedContractCallTxsResponse);
        message.calls = [];
        if (object.calls !== undefined && object.calls !== null) {
            for (var _i = 0, _a = object.calls; _i < _a.length; _i++) {
                var e = _a[_i];
                message.calls.push(gravity_1.ContractCallTx.fromPartial(e));
            }
        }
        return message;
    }
};
var baseBatchTxFeesRequest = {};
exports.BatchTxFeesRequest = {
    encode: function (_, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTxFeesRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (_) {
        var message = __assign({}, baseBatchTxFeesRequest);
        return message;
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
    fromPartial: function (_) {
        var message = __assign({}, baseBatchTxFeesRequest);
        return message;
    }
};
var baseBatchTxFeesResponse = {};
exports.BatchTxFeesResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.fees; _i < _a.length; _i++) {
            var v = _a[_i];
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTxFeesResponse);
        message.fees = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.fees.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchTxFeesResponse);
        message.fees = [];
        if (object.fees !== undefined && object.fees !== null) {
            for (var _i = 0, _a = object.fees; _i < _a.length; _i++) {
                var e = _a[_i];
                message.fees.push(coin_1.Coin.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.fees) {
            obj.fees = message.fees.map(function (e) { return (e ? coin_1.Coin.toJSON(e) : undefined); });
        }
        else {
            obj.fees = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchTxFeesResponse);
        message.fees = [];
        if (object.fees !== undefined && object.fees !== null) {
            for (var _i = 0, _a = object.fees; _i < _a.length; _i++) {
                var e = _a[_i];
                message.fees.push(coin_1.Coin.fromPartial(e));
            }
        }
        return message;
    }
};
var baseContractCallTxConfirmationsRequest = {
    invalidationNonce: long_1["default"].UZERO
};
exports.ContractCallTxConfirmationsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.invalidationScope.length !== 0) {
            writer.uint32(10).bytes(message.invalidationScope);
        }
        if (!message.invalidationNonce.isZero()) {
            writer.uint32(16).uint64(message.invalidationNonce);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseContractCallTxConfirmationsRequest);
        message.invalidationScope = new Uint8Array();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.invalidationScope = reader.bytes();
                    break;
                case 2:
                    message.invalidationNonce = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseContractCallTxConfirmationsRequest);
        message.invalidationScope = new Uint8Array();
        if (object.invalidationScope !== undefined &&
            object.invalidationScope !== null) {
            message.invalidationScope = bytesFromBase64(object.invalidationScope);
        }
        if (object.invalidationNonce !== undefined &&
            object.invalidationNonce !== null) {
            message.invalidationNonce = long_1["default"].fromString(object.invalidationNonce);
        }
        else {
            message.invalidationNonce = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.invalidationScope !== undefined &&
            (obj.invalidationScope = base64FromBytes(message.invalidationScope !== undefined
                ? message.invalidationScope
                : new Uint8Array()));
        message.invalidationNonce !== undefined &&
            (obj.invalidationNonce = (message.invalidationNonce || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseContractCallTxConfirmationsRequest);
        if (object.invalidationScope !== undefined &&
            object.invalidationScope !== null) {
            message.invalidationScope = object.invalidationScope;
        }
        else {
            message.invalidationScope = new Uint8Array();
        }
        if (object.invalidationNonce !== undefined &&
            object.invalidationNonce !== null) {
            message.invalidationNonce = object.invalidationNonce;
        }
        else {
            message.invalidationNonce = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseContractCallTxConfirmationsResponse = {};
exports.ContractCallTxConfirmationsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.signatures; _i < _a.length; _i++) {
            var v = _a[_i];
            msgs_1.ContractCallTxConfirmation.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseContractCallTxConfirmationsResponse);
        message.signatures = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signatures.push(msgs_1.ContractCallTxConfirmation.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseContractCallTxConfirmationsResponse);
        message.signatures = [];
        if (object.signatures !== undefined && object.signatures !== null) {
            for (var _i = 0, _a = object.signatures; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signatures.push(msgs_1.ContractCallTxConfirmation.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.signatures) {
            obj.signatures = message.signatures.map(function (e) {
                return e ? msgs_1.ContractCallTxConfirmation.toJSON(e) : undefined;
            });
        }
        else {
            obj.signatures = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseContractCallTxConfirmationsResponse);
        message.signatures = [];
        if (object.signatures !== undefined && object.signatures !== null) {
            for (var _i = 0, _a = object.signatures; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signatures.push(msgs_1.ContractCallTxConfirmation.fromPartial(e));
            }
        }
        return message;
    }
};
var baseBatchTxConfirmationsRequest = {
    batchNonce: long_1["default"].UZERO,
    tokenContract: ""
};
exports.BatchTxConfirmationsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.batchNonce.isZero()) {
            writer.uint32(8).uint64(message.batchNonce);
        }
        if (message.tokenContract !== "") {
            writer.uint32(18).string(message.tokenContract);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTxConfirmationsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.batchNonce = reader.uint64();
                    break;
                case 2:
                    message.tokenContract = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchTxConfirmationsRequest);
        if (object.batchNonce !== undefined && object.batchNonce !== null) {
            message.batchNonce = long_1["default"].fromString(object.batchNonce);
        }
        else {
            message.batchNonce = long_1["default"].UZERO;
        }
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = String(object.tokenContract);
        }
        else {
            message.tokenContract = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.batchNonce !== undefined &&
            (obj.batchNonce = (message.batchNonce || long_1["default"].UZERO).toString());
        message.tokenContract !== undefined &&
            (obj.tokenContract = message.tokenContract);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchTxConfirmationsRequest);
        if (object.batchNonce !== undefined && object.batchNonce !== null) {
            message.batchNonce = object.batchNonce;
        }
        else {
            message.batchNonce = long_1["default"].UZERO;
        }
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = object.tokenContract;
        }
        else {
            message.tokenContract = "";
        }
        return message;
    }
};
var baseBatchTxConfirmationsResponse = {};
exports.BatchTxConfirmationsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.signatures; _i < _a.length; _i++) {
            var v = _a[_i];
            msgs_1.BatchTxConfirmation.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTxConfirmationsResponse);
        message.signatures = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signatures.push(msgs_1.BatchTxConfirmation.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchTxConfirmationsResponse);
        message.signatures = [];
        if (object.signatures !== undefined && object.signatures !== null) {
            for (var _i = 0, _a = object.signatures; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signatures.push(msgs_1.BatchTxConfirmation.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.signatures) {
            obj.signatures = message.signatures.map(function (e) {
                return e ? msgs_1.BatchTxConfirmation.toJSON(e) : undefined;
            });
        }
        else {
            obj.signatures = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchTxConfirmationsResponse);
        message.signatures = [];
        if (object.signatures !== undefined && object.signatures !== null) {
            for (var _i = 0, _a = object.signatures; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signatures.push(msgs_1.BatchTxConfirmation.fromPartial(e));
            }
        }
        return message;
    }
};
var baseLastSubmittedEthereumEventRequest = { address: "" };
exports.LastSubmittedEthereumEventRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseLastSubmittedEthereumEventRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseLastSubmittedEthereumEventRequest);
        if (object.address !== undefined && object.address !== null) {
            message.address = String(object.address);
        }
        else {
            message.address = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.address !== undefined && (obj.address = message.address);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseLastSubmittedEthereumEventRequest);
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = "";
        }
        return message;
    }
};
var baseLastSubmittedEthereumEventResponse = {
    eventNonce: long_1["default"].UZERO
};
exports.LastSubmittedEthereumEventResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.eventNonce.isZero()) {
            writer.uint32(8).uint64(message.eventNonce);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseLastSubmittedEthereumEventResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.eventNonce = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseLastSubmittedEthereumEventResponse);
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = long_1["default"].fromString(object.eventNonce);
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.eventNonce !== undefined &&
            (obj.eventNonce = (message.eventNonce || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseLastSubmittedEthereumEventResponse);
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = object.eventNonce;
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseERC20ToDenomRequest = { erc20: "" };
exports.ERC20ToDenomRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.erc20 !== "") {
            writer.uint32(10).string(message.erc20);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseERC20ToDenomRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.erc20 = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseERC20ToDenomRequest);
        if (object.erc20 !== undefined && object.erc20 !== null) {
            message.erc20 = String(object.erc20);
        }
        else {
            message.erc20 = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.erc20 !== undefined && (obj.erc20 = message.erc20);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseERC20ToDenomRequest);
        if (object.erc20 !== undefined && object.erc20 !== null) {
            message.erc20 = object.erc20;
        }
        else {
            message.erc20 = "";
        }
        return message;
    }
};
var baseERC20ToDenomResponse = { denom: "", cosmosOriginated: false };
exports.ERC20ToDenomResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.denom !== "") {
            writer.uint32(10).string(message.denom);
        }
        if (message.cosmosOriginated === true) {
            writer.uint32(16).bool(message.cosmosOriginated);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseERC20ToDenomResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.denom = reader.string();
                    break;
                case 2:
                    message.cosmosOriginated = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseERC20ToDenomResponse);
        if (object.denom !== undefined && object.denom !== null) {
            message.denom = String(object.denom);
        }
        else {
            message.denom = "";
        }
        if (object.cosmosOriginated !== undefined &&
            object.cosmosOriginated !== null) {
            message.cosmosOriginated = Boolean(object.cosmosOriginated);
        }
        else {
            message.cosmosOriginated = false;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.denom !== undefined && (obj.denom = message.denom);
        message.cosmosOriginated !== undefined &&
            (obj.cosmosOriginated = message.cosmosOriginated);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseERC20ToDenomResponse);
        if (object.denom !== undefined && object.denom !== null) {
            message.denom = object.denom;
        }
        else {
            message.denom = "";
        }
        if (object.cosmosOriginated !== undefined &&
            object.cosmosOriginated !== null) {
            message.cosmosOriginated = object.cosmosOriginated;
        }
        else {
            message.cosmosOriginated = false;
        }
        return message;
    }
};
var baseDenomToERC20Request = { denom: "" };
exports.DenomToERC20Request = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.denom !== "") {
            writer.uint32(10).string(message.denom);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseDenomToERC20Request);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.denom = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseDenomToERC20Request);
        if (object.denom !== undefined && object.denom !== null) {
            message.denom = String(object.denom);
        }
        else {
            message.denom = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.denom !== undefined && (obj.denom = message.denom);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseDenomToERC20Request);
        if (object.denom !== undefined && object.denom !== null) {
            message.denom = object.denom;
        }
        else {
            message.denom = "";
        }
        return message;
    }
};
var baseDenomToERC20Response = { erc20: "", cosmosOriginated: false };
exports.DenomToERC20Response = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.erc20 !== "") {
            writer.uint32(10).string(message.erc20);
        }
        if (message.cosmosOriginated === true) {
            writer.uint32(16).bool(message.cosmosOriginated);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseDenomToERC20Response);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.erc20 = reader.string();
                    break;
                case 2:
                    message.cosmosOriginated = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseDenomToERC20Response);
        if (object.erc20 !== undefined && object.erc20 !== null) {
            message.erc20 = String(object.erc20);
        }
        else {
            message.erc20 = "";
        }
        if (object.cosmosOriginated !== undefined &&
            object.cosmosOriginated !== null) {
            message.cosmosOriginated = Boolean(object.cosmosOriginated);
        }
        else {
            message.cosmosOriginated = false;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.erc20 !== undefined && (obj.erc20 = message.erc20);
        message.cosmosOriginated !== undefined &&
            (obj.cosmosOriginated = message.cosmosOriginated);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseDenomToERC20Response);
        if (object.erc20 !== undefined && object.erc20 !== null) {
            message.erc20 = object.erc20;
        }
        else {
            message.erc20 = "";
        }
        if (object.cosmosOriginated !== undefined &&
            object.cosmosOriginated !== null) {
            message.cosmosOriginated = object.cosmosOriginated;
        }
        else {
            message.cosmosOriginated = false;
        }
        return message;
    }
};
var baseDelegateKeysByValidatorRequest = { validatorAddress: "" };
exports.DelegateKeysByValidatorRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.validatorAddress !== "") {
            writer.uint32(10).string(message.validatorAddress);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseDelegateKeysByValidatorRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseDelegateKeysByValidatorRequest);
        if (object.validatorAddress !== undefined &&
            object.validatorAddress !== null) {
            message.validatorAddress = String(object.validatorAddress);
        }
        else {
            message.validatorAddress = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.validatorAddress !== undefined &&
            (obj.validatorAddress = message.validatorAddress);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseDelegateKeysByValidatorRequest);
        if (object.validatorAddress !== undefined &&
            object.validatorAddress !== null) {
            message.validatorAddress = object.validatorAddress;
        }
        else {
            message.validatorAddress = "";
        }
        return message;
    }
};
var baseDelegateKeysByValidatorResponse = {
    ethAddress: "",
    orchestratorAddress: ""
};
exports.DelegateKeysByValidatorResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.ethAddress !== "") {
            writer.uint32(10).string(message.ethAddress);
        }
        if (message.orchestratorAddress !== "") {
            writer.uint32(18).string(message.orchestratorAddress);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseDelegateKeysByValidatorResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.ethAddress = reader.string();
                    break;
                case 2:
                    message.orchestratorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseDelegateKeysByValidatorResponse);
        if (object.ethAddress !== undefined && object.ethAddress !== null) {
            message.ethAddress = String(object.ethAddress);
        }
        else {
            message.ethAddress = "";
        }
        if (object.orchestratorAddress !== undefined &&
            object.orchestratorAddress !== null) {
            message.orchestratorAddress = String(object.orchestratorAddress);
        }
        else {
            message.orchestratorAddress = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.ethAddress !== undefined && (obj.ethAddress = message.ethAddress);
        message.orchestratorAddress !== undefined &&
            (obj.orchestratorAddress = message.orchestratorAddress);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseDelegateKeysByValidatorResponse);
        if (object.ethAddress !== undefined && object.ethAddress !== null) {
            message.ethAddress = object.ethAddress;
        }
        else {
            message.ethAddress = "";
        }
        if (object.orchestratorAddress !== undefined &&
            object.orchestratorAddress !== null) {
            message.orchestratorAddress = object.orchestratorAddress;
        }
        else {
            message.orchestratorAddress = "";
        }
        return message;
    }
};
var baseDelegateKeysByEthereumSignerRequest = { ethereumSigner: "" };
exports.DelegateKeysByEthereumSignerRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.ethereumSigner !== "") {
            writer.uint32(10).string(message.ethereumSigner);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseDelegateKeysByEthereumSignerRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.ethereumSigner = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseDelegateKeysByEthereumSignerRequest);
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = String(object.ethereumSigner);
        }
        else {
            message.ethereumSigner = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.ethereumSigner !== undefined &&
            (obj.ethereumSigner = message.ethereumSigner);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseDelegateKeysByEthereumSignerRequest);
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = object.ethereumSigner;
        }
        else {
            message.ethereumSigner = "";
        }
        return message;
    }
};
var baseDelegateKeysByEthereumSignerResponse = {
    validatorAddress: "",
    orchestratorAddress: ""
};
exports.DelegateKeysByEthereumSignerResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.validatorAddress !== "") {
            writer.uint32(10).string(message.validatorAddress);
        }
        if (message.orchestratorAddress !== "") {
            writer.uint32(18).string(message.orchestratorAddress);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseDelegateKeysByEthereumSignerResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddress = reader.string();
                    break;
                case 2:
                    message.orchestratorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseDelegateKeysByEthereumSignerResponse);
        if (object.validatorAddress !== undefined &&
            object.validatorAddress !== null) {
            message.validatorAddress = String(object.validatorAddress);
        }
        else {
            message.validatorAddress = "";
        }
        if (object.orchestratorAddress !== undefined &&
            object.orchestratorAddress !== null) {
            message.orchestratorAddress = String(object.orchestratorAddress);
        }
        else {
            message.orchestratorAddress = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.validatorAddress !== undefined &&
            (obj.validatorAddress = message.validatorAddress);
        message.orchestratorAddress !== undefined &&
            (obj.orchestratorAddress = message.orchestratorAddress);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseDelegateKeysByEthereumSignerResponse);
        if (object.validatorAddress !== undefined &&
            object.validatorAddress !== null) {
            message.validatorAddress = object.validatorAddress;
        }
        else {
            message.validatorAddress = "";
        }
        if (object.orchestratorAddress !== undefined &&
            object.orchestratorAddress !== null) {
            message.orchestratorAddress = object.orchestratorAddress;
        }
        else {
            message.orchestratorAddress = "";
        }
        return message;
    }
};
var baseDelegateKeysByOrchestratorRequest = {
    orchestratorAddress: ""
};
exports.DelegateKeysByOrchestratorRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.orchestratorAddress !== "") {
            writer.uint32(10).string(message.orchestratorAddress);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseDelegateKeysByOrchestratorRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.orchestratorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseDelegateKeysByOrchestratorRequest);
        if (object.orchestratorAddress !== undefined &&
            object.orchestratorAddress !== null) {
            message.orchestratorAddress = String(object.orchestratorAddress);
        }
        else {
            message.orchestratorAddress = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.orchestratorAddress !== undefined &&
            (obj.orchestratorAddress = message.orchestratorAddress);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseDelegateKeysByOrchestratorRequest);
        if (object.orchestratorAddress !== undefined &&
            object.orchestratorAddress !== null) {
            message.orchestratorAddress = object.orchestratorAddress;
        }
        else {
            message.orchestratorAddress = "";
        }
        return message;
    }
};
var baseDelegateKeysByOrchestratorResponse = {
    validatorAddress: "",
    ethereumSigner: ""
};
exports.DelegateKeysByOrchestratorResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.validatorAddress !== "") {
            writer.uint32(10).string(message.validatorAddress);
        }
        if (message.ethereumSigner !== "") {
            writer.uint32(18).string(message.ethereumSigner);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseDelegateKeysByOrchestratorResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddress = reader.string();
                    break;
                case 2:
                    message.ethereumSigner = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseDelegateKeysByOrchestratorResponse);
        if (object.validatorAddress !== undefined &&
            object.validatorAddress !== null) {
            message.validatorAddress = String(object.validatorAddress);
        }
        else {
            message.validatorAddress = "";
        }
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = String(object.ethereumSigner);
        }
        else {
            message.ethereumSigner = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.validatorAddress !== undefined &&
            (obj.validatorAddress = message.validatorAddress);
        message.ethereumSigner !== undefined &&
            (obj.ethereumSigner = message.ethereumSigner);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseDelegateKeysByOrchestratorResponse);
        if (object.validatorAddress !== undefined &&
            object.validatorAddress !== null) {
            message.validatorAddress = object.validatorAddress;
        }
        else {
            message.validatorAddress = "";
        }
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = object.ethereumSigner;
        }
        else {
            message.ethereumSigner = "";
        }
        return message;
    }
};
var baseBatchedSendToEthereumsRequest = { senderAddress: "" };
exports.BatchedSendToEthereumsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.senderAddress !== "") {
            writer.uint32(10).string(message.senderAddress);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchedSendToEthereumsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.senderAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchedSendToEthereumsRequest);
        if (object.senderAddress !== undefined && object.senderAddress !== null) {
            message.senderAddress = String(object.senderAddress);
        }
        else {
            message.senderAddress = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.senderAddress !== undefined &&
            (obj.senderAddress = message.senderAddress);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchedSendToEthereumsRequest);
        if (object.senderAddress !== undefined && object.senderAddress !== null) {
            message.senderAddress = object.senderAddress;
        }
        else {
            message.senderAddress = "";
        }
        return message;
    }
};
var baseBatchedSendToEthereumsResponse = {};
exports.BatchedSendToEthereumsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.sendToEthereums; _i < _a.length; _i++) {
            var v = _a[_i];
            gravity_1.SendToEthereum.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchedSendToEthereumsResponse);
        message.sendToEthereums = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sendToEthereums.push(gravity_1.SendToEthereum.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchedSendToEthereumsResponse);
        message.sendToEthereums = [];
        if (object.sendToEthereums !== undefined &&
            object.sendToEthereums !== null) {
            for (var _i = 0, _a = object.sendToEthereums; _i < _a.length; _i++) {
                var e = _a[_i];
                message.sendToEthereums.push(gravity_1.SendToEthereum.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.sendToEthereums) {
            obj.sendToEthereums = message.sendToEthereums.map(function (e) {
                return e ? gravity_1.SendToEthereum.toJSON(e) : undefined;
            });
        }
        else {
            obj.sendToEthereums = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchedSendToEthereumsResponse);
        message.sendToEthereums = [];
        if (object.sendToEthereums !== undefined &&
            object.sendToEthereums !== null) {
            for (var _i = 0, _a = object.sendToEthereums; _i < _a.length; _i++) {
                var e = _a[_i];
                message.sendToEthereums.push(gravity_1.SendToEthereum.fromPartial(e));
            }
        }
        return message;
    }
};
var baseUnbatchedSendToEthereumsRequest = { senderAddress: "" };
exports.UnbatchedSendToEthereumsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.senderAddress !== "") {
            writer.uint32(10).string(message.senderAddress);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseUnbatchedSendToEthereumsRequest);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.senderAddress = reader.string();
                    break;
                case 2:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseUnbatchedSendToEthereumsRequest);
        if (object.senderAddress !== undefined && object.senderAddress !== null) {
            message.senderAddress = String(object.senderAddress);
        }
        else {
            message.senderAddress = "";
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.senderAddress !== undefined &&
            (obj.senderAddress = message.senderAddress);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? pagination_1.PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseUnbatchedSendToEthereumsRequest);
        if (object.senderAddress !== undefined && object.senderAddress !== null) {
            message.senderAddress = object.senderAddress;
        }
        else {
            message.senderAddress = "";
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    }
};
var baseUnbatchedSendToEthereumsResponse = {};
exports.UnbatchedSendToEthereumsResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        for (var _i = 0, _a = message.sendToEthereums; _i < _a.length; _i++) {
            var v = _a[_i];
            gravity_1.SendToEthereum.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseUnbatchedSendToEthereumsResponse);
        message.sendToEthereums = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sendToEthereums.push(gravity_1.SendToEthereum.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseUnbatchedSendToEthereumsResponse);
        message.sendToEthereums = [];
        if (object.sendToEthereums !== undefined &&
            object.sendToEthereums !== null) {
            for (var _i = 0, _a = object.sendToEthereums; _i < _a.length; _i++) {
                var e = _a[_i];
                message.sendToEthereums.push(gravity_1.SendToEthereum.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.sendToEthereums) {
            obj.sendToEthereums = message.sendToEthereums.map(function (e) {
                return e ? gravity_1.SendToEthereum.toJSON(e) : undefined;
            });
        }
        else {
            obj.sendToEthereums = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? pagination_1.PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseUnbatchedSendToEthereumsResponse);
        message.sendToEthereums = [];
        if (object.sendToEthereums !== undefined &&
            object.sendToEthereums !== null) {
            for (var _i = 0, _a = object.sendToEthereums; _i < _a.length; _i++) {
                var e = _a[_i];
                message.sendToEthereums.push(gravity_1.SendToEthereum.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    }
};
var QueryClientImpl = /** @class */ (function () {
    function QueryClientImpl(rpc) {
        this.rpc = rpc;
        this.Params = this.Params.bind(this);
        this.SignerSetTx = this.SignerSetTx.bind(this);
        this.LatestSignerSetTx = this.LatestSignerSetTx.bind(this);
        this.BatchTx = this.BatchTx.bind(this);
        this.ContractCallTx = this.ContractCallTx.bind(this);
        this.SignerSetTxs = this.SignerSetTxs.bind(this);
        this.BatchTxs = this.BatchTxs.bind(this);
        this.ContractCallTxs = this.ContractCallTxs.bind(this);
        this.SignerSetTxConfirmations = this.SignerSetTxConfirmations.bind(this);
        this.BatchTxConfirmations = this.BatchTxConfirmations.bind(this);
        this.ContractCallTxConfirmations = this.ContractCallTxConfirmations.bind(this);
        this.UnsignedSignerSetTxs = this.UnsignedSignerSetTxs.bind(this);
        this.UnsignedBatchTxs = this.UnsignedBatchTxs.bind(this);
        this.UnsignedContractCallTxs = this.UnsignedContractCallTxs.bind(this);
        this.LastSubmittedEthereumEvent = this.LastSubmittedEthereumEvent.bind(this);
        this.BatchTxFees = this.BatchTxFees.bind(this);
        this.ERC20ToDenom = this.ERC20ToDenom.bind(this);
        this.DenomToERC20 = this.DenomToERC20.bind(this);
        this.BatchedSendToEthereums = this.BatchedSendToEthereums.bind(this);
        this.UnbatchedSendToEthereums = this.UnbatchedSendToEthereums.bind(this);
        this.DelegateKeysByValidator = this.DelegateKeysByValidator.bind(this);
        this.DelegateKeysByEthereumSigner = this.DelegateKeysByEthereumSigner.bind(this);
        this.DelegateKeysByOrchestrator = this.DelegateKeysByOrchestrator.bind(this);
    }
    QueryClientImpl.prototype.Params = function (request) {
        var data = exports.ParamsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "Params", data);
        return promise.then(function (data) { return exports.ParamsResponse.decode(new minimal_1["default"].Reader(data)); });
    };
    QueryClientImpl.prototype.SignerSetTx = function (request) {
        var data = exports.SignerSetTxRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "SignerSetTx", data);
        return promise.then(function (data) {
            return exports.SignerSetTxResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.LatestSignerSetTx = function (request) {
        var data = exports.LatestSignerSetTxRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "LatestSignerSetTx", data);
        return promise.then(function (data) {
            return exports.SignerSetTxResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.BatchTx = function (request) {
        var data = exports.BatchTxRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "BatchTx", data);
        return promise.then(function (data) { return exports.BatchTxResponse.decode(new minimal_1["default"].Reader(data)); });
    };
    QueryClientImpl.prototype.ContractCallTx = function (request) {
        var data = exports.ContractCallTxRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "ContractCallTx", data);
        return promise.then(function (data) {
            return exports.ContractCallTxResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.SignerSetTxs = function (request) {
        var data = exports.SignerSetTxsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "SignerSetTxs", data);
        return promise.then(function (data) {
            return exports.SignerSetTxsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.BatchTxs = function (request) {
        var data = exports.BatchTxsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "BatchTxs", data);
        return promise.then(function (data) {
            return exports.BatchTxsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.ContractCallTxs = function (request) {
        var data = exports.ContractCallTxsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "ContractCallTxs", data);
        return promise.then(function (data) {
            return exports.ContractCallTxsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.SignerSetTxConfirmations = function (request) {
        var data = exports.SignerSetTxConfirmationsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "SignerSetTxConfirmations", data);
        return promise.then(function (data) {
            return exports.SignerSetTxConfirmationsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.BatchTxConfirmations = function (request) {
        var data = exports.BatchTxConfirmationsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "BatchTxConfirmations", data);
        return promise.then(function (data) {
            return exports.BatchTxConfirmationsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.ContractCallTxConfirmations = function (request) {
        var data = exports.ContractCallTxConfirmationsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "ContractCallTxConfirmations", data);
        return promise.then(function (data) {
            return exports.ContractCallTxConfirmationsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.UnsignedSignerSetTxs = function (request) {
        var data = exports.UnsignedSignerSetTxsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "UnsignedSignerSetTxs", data);
        return promise.then(function (data) {
            return exports.UnsignedSignerSetTxsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.UnsignedBatchTxs = function (request) {
        var data = exports.UnsignedBatchTxsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "UnsignedBatchTxs", data);
        return promise.then(function (data) {
            return exports.UnsignedBatchTxsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.UnsignedContractCallTxs = function (request) {
        var data = exports.UnsignedContractCallTxsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "UnsignedContractCallTxs", data);
        return promise.then(function (data) {
            return exports.UnsignedContractCallTxsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.LastSubmittedEthereumEvent = function (request) {
        var data = exports.LastSubmittedEthereumEventRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "LastSubmittedEthereumEvent", data);
        return promise.then(function (data) {
            return exports.LastSubmittedEthereumEventResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.BatchTxFees = function (request) {
        var data = exports.BatchTxFeesRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "BatchTxFees", data);
        return promise.then(function (data) {
            return exports.BatchTxFeesResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.ERC20ToDenom = function (request) {
        var data = exports.ERC20ToDenomRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "ERC20ToDenom", data);
        return promise.then(function (data) {
            return exports.ERC20ToDenomResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.DenomToERC20 = function (request) {
        var data = exports.DenomToERC20Request.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "DenomToERC20", data);
        return promise.then(function (data) {
            return exports.DenomToERC20Response.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.BatchedSendToEthereums = function (request) {
        var data = exports.BatchedSendToEthereumsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "BatchedSendToEthereums", data);
        return promise.then(function (data) {
            return exports.BatchedSendToEthereumsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.UnbatchedSendToEthereums = function (request) {
        var data = exports.UnbatchedSendToEthereumsRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "UnbatchedSendToEthereums", data);
        return promise.then(function (data) {
            return exports.UnbatchedSendToEthereumsResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.DelegateKeysByValidator = function (request) {
        var data = exports.DelegateKeysByValidatorRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "DelegateKeysByValidator", data);
        return promise.then(function (data) {
            return exports.DelegateKeysByValidatorResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.DelegateKeysByEthereumSigner = function (request) {
        var data = exports.DelegateKeysByEthereumSignerRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "DelegateKeysByEthereumSigner", data);
        return promise.then(function (data) {
            return exports.DelegateKeysByEthereumSignerResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    QueryClientImpl.prototype.DelegateKeysByOrchestrator = function (request) {
        var data = exports.DelegateKeysByOrchestratorRequest.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Query", "DelegateKeysByOrchestrator", data);
        return promise.then(function (data) {
            return exports.DelegateKeysByOrchestratorResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    return QueryClientImpl;
}());
exports.QueryClientImpl = QueryClientImpl;
var globalThis = (function () {
    if (typeof globalThis !== "undefined")
        return globalThis;
    if (typeof self !== "undefined")
        return self;
    if (typeof window !== "undefined")
        return window;
    if (typeof global !== "undefined")
        return global;
    throw "Unable to locate global object";
})();
var atob = globalThis.atob ||
    (function (b64) { return globalThis.Buffer.from(b64, "base64").toString("binary"); });
function bytesFromBase64(b64) {
    var bin = atob(b64);
    var arr = new Uint8Array(bin.length);
    for (var i = 0; i < bin.length; ++i) {
        arr[i] = bin.charCodeAt(i);
    }
    return arr;
}
var btoa = globalThis.btoa ||
    (function (bin) { return globalThis.Buffer.from(bin, "binary").toString("base64"); });
function base64FromBytes(arr) {
    var bin = [];
    for (var i = 0; i < arr.byteLength; ++i) {
        bin.push(String.fromCharCode(arr[i]));
    }
    return btoa(bin.join(""));
}
if (minimal_1["default"].util.Long !== long_1["default"]) {
    minimal_1["default"].util.Long = long_1["default"];
    minimal_1["default"].configure();
}
