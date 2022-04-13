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
exports.ERC20ToDenom = exports.GenesisState = exports.Params = exports.protobufPackage = void 0;
/* eslint-disable */
var long_1 = require("long");
var minimal_1 = require("protobufjs/minimal");
var any_1 = require("../../google/protobuf/any");
var gravity_1 = require("../../gravity/v1/gravity");
var msgs_1 = require("../../gravity/v1/msgs");
exports.protobufPackage = "gravity.v1";
var baseParams = {
    gravityId: "",
    contractSourceHash: "",
    bridgeEthereumAddress: "",
    bridgeChainId: long_1["default"].UZERO,
    signedSignerSetTxsWindow: long_1["default"].UZERO,
    signedBatchesWindow: long_1["default"].UZERO,
    ethereumSignaturesWindow: long_1["default"].UZERO,
    targetBatchTimeout: long_1["default"].UZERO,
    averageBlockTime: long_1["default"].UZERO,
    averageEthereumBlockTime: long_1["default"].UZERO,
    unbondSlashingSignerSetTxsWindow: long_1["default"].UZERO
};
exports.Params = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.gravityId !== "") {
            writer.uint32(10).string(message.gravityId);
        }
        if (message.contractSourceHash !== "") {
            writer.uint32(18).string(message.contractSourceHash);
        }
        if (message.bridgeEthereumAddress !== "") {
            writer.uint32(34).string(message.bridgeEthereumAddress);
        }
        if (!message.bridgeChainId.isZero()) {
            writer.uint32(40).uint64(message.bridgeChainId);
        }
        if (!message.signedSignerSetTxsWindow.isZero()) {
            writer.uint32(48).uint64(message.signedSignerSetTxsWindow);
        }
        if (!message.signedBatchesWindow.isZero()) {
            writer.uint32(56).uint64(message.signedBatchesWindow);
        }
        if (!message.ethereumSignaturesWindow.isZero()) {
            writer.uint32(64).uint64(message.ethereumSignaturesWindow);
        }
        if (!message.targetBatchTimeout.isZero()) {
            writer.uint32(80).uint64(message.targetBatchTimeout);
        }
        if (!message.averageBlockTime.isZero()) {
            writer.uint32(88).uint64(message.averageBlockTime);
        }
        if (!message.averageEthereumBlockTime.isZero()) {
            writer.uint32(96).uint64(message.averageEthereumBlockTime);
        }
        if (message.slashFractionSignerSetTx.length !== 0) {
            writer.uint32(106).bytes(message.slashFractionSignerSetTx);
        }
        if (message.slashFractionBatch.length !== 0) {
            writer.uint32(114).bytes(message.slashFractionBatch);
        }
        if (message.slashFractionEthereumSignature.length !== 0) {
            writer.uint32(122).bytes(message.slashFractionEthereumSignature);
        }
        if (message.slashFractionConflictingEthereumSignature.length !== 0) {
            writer
                .uint32(130)
                .bytes(message.slashFractionConflictingEthereumSignature);
        }
        if (!message.unbondSlashingSignerSetTxsWindow.isZero()) {
            writer.uint32(136).uint64(message.unbondSlashingSignerSetTxsWindow);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseParams);
        message.slashFractionSignerSetTx = new Uint8Array();
        message.slashFractionBatch = new Uint8Array();
        message.slashFractionEthereumSignature = new Uint8Array();
        message.slashFractionConflictingEthereumSignature = new Uint8Array();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.gravityId = reader.string();
                    break;
                case 2:
                    message.contractSourceHash = reader.string();
                    break;
                case 4:
                    message.bridgeEthereumAddress = reader.string();
                    break;
                case 5:
                    message.bridgeChainId = reader.uint64();
                    break;
                case 6:
                    message.signedSignerSetTxsWindow = reader.uint64();
                    break;
                case 7:
                    message.signedBatchesWindow = reader.uint64();
                    break;
                case 8:
                    message.ethereumSignaturesWindow = reader.uint64();
                    break;
                case 10:
                    message.targetBatchTimeout = reader.uint64();
                    break;
                case 11:
                    message.averageBlockTime = reader.uint64();
                    break;
                case 12:
                    message.averageEthereumBlockTime = reader.uint64();
                    break;
                case 13:
                    message.slashFractionSignerSetTx = reader.bytes();
                    break;
                case 14:
                    message.slashFractionBatch = reader.bytes();
                    break;
                case 15:
                    message.slashFractionEthereumSignature = reader.bytes();
                    break;
                case 16:
                    message.slashFractionConflictingEthereumSignature = reader.bytes();
                    break;
                case 17:
                    message.unbondSlashingSignerSetTxsWindow = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseParams);
        message.slashFractionSignerSetTx = new Uint8Array();
        message.slashFractionBatch = new Uint8Array();
        message.slashFractionEthereumSignature = new Uint8Array();
        message.slashFractionConflictingEthereumSignature = new Uint8Array();
        if (object.gravityId !== undefined && object.gravityId !== null) {
            message.gravityId = String(object.gravityId);
        }
        else {
            message.gravityId = "";
        }
        if (object.contractSourceHash !== undefined &&
            object.contractSourceHash !== null) {
            message.contractSourceHash = String(object.contractSourceHash);
        }
        else {
            message.contractSourceHash = "";
        }
        if (object.bridgeEthereumAddress !== undefined &&
            object.bridgeEthereumAddress !== null) {
            message.bridgeEthereumAddress = String(object.bridgeEthereumAddress);
        }
        else {
            message.bridgeEthereumAddress = "";
        }
        if (object.bridgeChainId !== undefined && object.bridgeChainId !== null) {
            message.bridgeChainId = long_1["default"].fromString(object.bridgeChainId);
        }
        else {
            message.bridgeChainId = long_1["default"].UZERO;
        }
        if (object.signedSignerSetTxsWindow !== undefined &&
            object.signedSignerSetTxsWindow !== null) {
            message.signedSignerSetTxsWindow = long_1["default"].fromString(object.signedSignerSetTxsWindow);
        }
        else {
            message.signedSignerSetTxsWindow = long_1["default"].UZERO;
        }
        if (object.signedBatchesWindow !== undefined &&
            object.signedBatchesWindow !== null) {
            message.signedBatchesWindow = long_1["default"].fromString(object.signedBatchesWindow);
        }
        else {
            message.signedBatchesWindow = long_1["default"].UZERO;
        }
        if (object.ethereumSignaturesWindow !== undefined &&
            object.ethereumSignaturesWindow !== null) {
            message.ethereumSignaturesWindow = long_1["default"].fromString(object.ethereumSignaturesWindow);
        }
        else {
            message.ethereumSignaturesWindow = long_1["default"].UZERO;
        }
        if (object.targetBatchTimeout !== undefined &&
            object.targetBatchTimeout !== null) {
            message.targetBatchTimeout = long_1["default"].fromString(object.targetBatchTimeout);
        }
        else {
            message.targetBatchTimeout = long_1["default"].UZERO;
        }
        if (object.averageBlockTime !== undefined &&
            object.averageBlockTime !== null) {
            message.averageBlockTime = long_1["default"].fromString(object.averageBlockTime);
        }
        else {
            message.averageBlockTime = long_1["default"].UZERO;
        }
        if (object.averageEthereumBlockTime !== undefined &&
            object.averageEthereumBlockTime !== null) {
            message.averageEthereumBlockTime = long_1["default"].fromString(object.averageEthereumBlockTime);
        }
        else {
            message.averageEthereumBlockTime = long_1["default"].UZERO;
        }
        if (object.slashFractionSignerSetTx !== undefined &&
            object.slashFractionSignerSetTx !== null) {
            message.slashFractionSignerSetTx = bytesFromBase64(object.slashFractionSignerSetTx);
        }
        if (object.slashFractionBatch !== undefined &&
            object.slashFractionBatch !== null) {
            message.slashFractionBatch = bytesFromBase64(object.slashFractionBatch);
        }
        if (object.slashFractionEthereumSignature !== undefined &&
            object.slashFractionEthereumSignature !== null) {
            message.slashFractionEthereumSignature = bytesFromBase64(object.slashFractionEthereumSignature);
        }
        if (object.slashFractionConflictingEthereumSignature !== undefined &&
            object.slashFractionConflictingEthereumSignature !== null) {
            message.slashFractionConflictingEthereumSignature = bytesFromBase64(object.slashFractionConflictingEthereumSignature);
        }
        if (object.unbondSlashingSignerSetTxsWindow !== undefined &&
            object.unbondSlashingSignerSetTxsWindow !== null) {
            message.unbondSlashingSignerSetTxsWindow = long_1["default"].fromString(object.unbondSlashingSignerSetTxsWindow);
        }
        else {
            message.unbondSlashingSignerSetTxsWindow = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.gravityId !== undefined && (obj.gravityId = message.gravityId);
        message.contractSourceHash !== undefined &&
            (obj.contractSourceHash = message.contractSourceHash);
        message.bridgeEthereumAddress !== undefined &&
            (obj.bridgeEthereumAddress = message.bridgeEthereumAddress);
        message.bridgeChainId !== undefined &&
            (obj.bridgeChainId = (message.bridgeChainId || long_1["default"].UZERO).toString());
        message.signedSignerSetTxsWindow !== undefined &&
            (obj.signedSignerSetTxsWindow = (message.signedSignerSetTxsWindow || long_1["default"].UZERO).toString());
        message.signedBatchesWindow !== undefined &&
            (obj.signedBatchesWindow = (message.signedBatchesWindow || long_1["default"].UZERO).toString());
        message.ethereumSignaturesWindow !== undefined &&
            (obj.ethereumSignaturesWindow = (message.ethereumSignaturesWindow || long_1["default"].UZERO).toString());
        message.targetBatchTimeout !== undefined &&
            (obj.targetBatchTimeout = (message.targetBatchTimeout || long_1["default"].UZERO).toString());
        message.averageBlockTime !== undefined &&
            (obj.averageBlockTime = (message.averageBlockTime || long_1["default"].UZERO).toString());
        message.averageEthereumBlockTime !== undefined &&
            (obj.averageEthereumBlockTime = (message.averageEthereumBlockTime || long_1["default"].UZERO).toString());
        message.slashFractionSignerSetTx !== undefined &&
            (obj.slashFractionSignerSetTx = base64FromBytes(message.slashFractionSignerSetTx !== undefined
                ? message.slashFractionSignerSetTx
                : new Uint8Array()));
        message.slashFractionBatch !== undefined &&
            (obj.slashFractionBatch = base64FromBytes(message.slashFractionBatch !== undefined
                ? message.slashFractionBatch
                : new Uint8Array()));
        message.slashFractionEthereumSignature !== undefined &&
            (obj.slashFractionEthereumSignature = base64FromBytes(message.slashFractionEthereumSignature !== undefined
                ? message.slashFractionEthereumSignature
                : new Uint8Array()));
        message.slashFractionConflictingEthereumSignature !== undefined &&
            (obj.slashFractionConflictingEthereumSignature = base64FromBytes(message.slashFractionConflictingEthereumSignature !== undefined
                ? message.slashFractionConflictingEthereumSignature
                : new Uint8Array()));
        message.unbondSlashingSignerSetTxsWindow !== undefined &&
            (obj.unbondSlashingSignerSetTxsWindow = (message.unbondSlashingSignerSetTxsWindow || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseParams);
        if (object.gravityId !== undefined && object.gravityId !== null) {
            message.gravityId = object.gravityId;
        }
        else {
            message.gravityId = "";
        }
        if (object.contractSourceHash !== undefined &&
            object.contractSourceHash !== null) {
            message.contractSourceHash = object.contractSourceHash;
        }
        else {
            message.contractSourceHash = "";
        }
        if (object.bridgeEthereumAddress !== undefined &&
            object.bridgeEthereumAddress !== null) {
            message.bridgeEthereumAddress = object.bridgeEthereumAddress;
        }
        else {
            message.bridgeEthereumAddress = "";
        }
        if (object.bridgeChainId !== undefined && object.bridgeChainId !== null) {
            message.bridgeChainId = object.bridgeChainId;
        }
        else {
            message.bridgeChainId = long_1["default"].UZERO;
        }
        if (object.signedSignerSetTxsWindow !== undefined &&
            object.signedSignerSetTxsWindow !== null) {
            message.signedSignerSetTxsWindow = object.signedSignerSetTxsWindow;
        }
        else {
            message.signedSignerSetTxsWindow = long_1["default"].UZERO;
        }
        if (object.signedBatchesWindow !== undefined &&
            object.signedBatchesWindow !== null) {
            message.signedBatchesWindow = object.signedBatchesWindow;
        }
        else {
            message.signedBatchesWindow = long_1["default"].UZERO;
        }
        if (object.ethereumSignaturesWindow !== undefined &&
            object.ethereumSignaturesWindow !== null) {
            message.ethereumSignaturesWindow = object.ethereumSignaturesWindow;
        }
        else {
            message.ethereumSignaturesWindow = long_1["default"].UZERO;
        }
        if (object.targetBatchTimeout !== undefined &&
            object.targetBatchTimeout !== null) {
            message.targetBatchTimeout = object.targetBatchTimeout;
        }
        else {
            message.targetBatchTimeout = long_1["default"].UZERO;
        }
        if (object.averageBlockTime !== undefined &&
            object.averageBlockTime !== null) {
            message.averageBlockTime = object.averageBlockTime;
        }
        else {
            message.averageBlockTime = long_1["default"].UZERO;
        }
        if (object.averageEthereumBlockTime !== undefined &&
            object.averageEthereumBlockTime !== null) {
            message.averageEthereumBlockTime = object.averageEthereumBlockTime;
        }
        else {
            message.averageEthereumBlockTime = long_1["default"].UZERO;
        }
        if (object.slashFractionSignerSetTx !== undefined &&
            object.slashFractionSignerSetTx !== null) {
            message.slashFractionSignerSetTx = object.slashFractionSignerSetTx;
        }
        else {
            message.slashFractionSignerSetTx = new Uint8Array();
        }
        if (object.slashFractionBatch !== undefined &&
            object.slashFractionBatch !== null) {
            message.slashFractionBatch = object.slashFractionBatch;
        }
        else {
            message.slashFractionBatch = new Uint8Array();
        }
        if (object.slashFractionEthereumSignature !== undefined &&
            object.slashFractionEthereumSignature !== null) {
            message.slashFractionEthereumSignature =
                object.slashFractionEthereumSignature;
        }
        else {
            message.slashFractionEthereumSignature = new Uint8Array();
        }
        if (object.slashFractionConflictingEthereumSignature !== undefined &&
            object.slashFractionConflictingEthereumSignature !== null) {
            message.slashFractionConflictingEthereumSignature =
                object.slashFractionConflictingEthereumSignature;
        }
        else {
            message.slashFractionConflictingEthereumSignature = new Uint8Array();
        }
        if (object.unbondSlashingSignerSetTxsWindow !== undefined &&
            object.unbondSlashingSignerSetTxsWindow !== null) {
            message.unbondSlashingSignerSetTxsWindow = object.unbondSlashingSignerSetTxsWindow;
        }
        else {
            message.unbondSlashingSignerSetTxsWindow = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseGenesisState = { lastObservedEventNonce: long_1["default"].UZERO };
exports.GenesisState = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.params !== undefined) {
            exports.Params.encode(message.params, writer.uint32(10).fork()).ldelim();
        }
        if (!message.lastObservedEventNonce.isZero()) {
            writer.uint32(16).uint64(message.lastObservedEventNonce);
        }
        for (var _i = 0, _a = message.outgoingTxs; _i < _a.length; _i++) {
            var v = _a[_i];
            any_1.Any.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (var _b = 0, _c = message.confirmations; _b < _c.length; _b++) {
            var v = _c[_b];
            any_1.Any.encode(v, writer.uint32(34).fork()).ldelim();
        }
        for (var _d = 0, _e = message.ethereumEventVoteRecords; _d < _e.length; _d++) {
            var v = _e[_d];
            gravity_1.EthereumEventVoteRecord.encode(v, writer.uint32(74).fork()).ldelim();
        }
        for (var _f = 0, _g = message.delegateKeys; _f < _g.length; _f++) {
            var v = _g[_f];
            msgs_1.MsgDelegateKeys.encode(v, writer.uint32(82).fork()).ldelim();
        }
        for (var _h = 0, _j = message.erc20ToDenoms; _h < _j.length; _h++) {
            var v = _j[_h];
            exports.ERC20ToDenom.encode(v, writer.uint32(90).fork()).ldelim();
        }
        for (var _k = 0, _l = message.unbatchedSendToEthereumTxs; _k < _l.length; _k++) {
            var v = _l[_k];
            gravity_1.SendToEthereum.encode(v, writer.uint32(98).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseGenesisState);
        message.outgoingTxs = [];
        message.confirmations = [];
        message.ethereumEventVoteRecords = [];
        message.delegateKeys = [];
        message.erc20ToDenoms = [];
        message.unbatchedSendToEthereumTxs = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.params = exports.Params.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.lastObservedEventNonce = reader.uint64();
                    break;
                case 3:
                    message.outgoingTxs.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.confirmations.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                case 9:
                    message.ethereumEventVoteRecords.push(gravity_1.EthereumEventVoteRecord.decode(reader, reader.uint32()));
                    break;
                case 10:
                    message.delegateKeys.push(msgs_1.MsgDelegateKeys.decode(reader, reader.uint32()));
                    break;
                case 11:
                    message.erc20ToDenoms.push(exports.ERC20ToDenom.decode(reader, reader.uint32()));
                    break;
                case 12:
                    message.unbatchedSendToEthereumTxs.push(gravity_1.SendToEthereum.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseGenesisState);
        message.outgoingTxs = [];
        message.confirmations = [];
        message.ethereumEventVoteRecords = [];
        message.delegateKeys = [];
        message.erc20ToDenoms = [];
        message.unbatchedSendToEthereumTxs = [];
        if (object.params !== undefined && object.params !== null) {
            message.params = exports.Params.fromJSON(object.params);
        }
        else {
            message.params = undefined;
        }
        if (object.lastObservedEventNonce !== undefined &&
            object.lastObservedEventNonce !== null) {
            message.lastObservedEventNonce = long_1["default"].fromString(object.lastObservedEventNonce);
        }
        else {
            message.lastObservedEventNonce = long_1["default"].UZERO;
        }
        if (object.outgoingTxs !== undefined && object.outgoingTxs !== null) {
            for (var _i = 0, _a = object.outgoingTxs; _i < _a.length; _i++) {
                var e = _a[_i];
                message.outgoingTxs.push(any_1.Any.fromJSON(e));
            }
        }
        if (object.confirmations !== undefined && object.confirmations !== null) {
            for (var _b = 0, _c = object.confirmations; _b < _c.length; _b++) {
                var e = _c[_b];
                message.confirmations.push(any_1.Any.fromJSON(e));
            }
        }
        if (object.ethereumEventVoteRecords !== undefined &&
            object.ethereumEventVoteRecords !== null) {
            for (var _d = 0, _e = object.ethereumEventVoteRecords; _d < _e.length; _d++) {
                var e = _e[_d];
                message.ethereumEventVoteRecords.push(gravity_1.EthereumEventVoteRecord.fromJSON(e));
            }
        }
        if (object.delegateKeys !== undefined && object.delegateKeys !== null) {
            for (var _f = 0, _g = object.delegateKeys; _f < _g.length; _f++) {
                var e = _g[_f];
                message.delegateKeys.push(msgs_1.MsgDelegateKeys.fromJSON(e));
            }
        }
        if (object.erc20ToDenoms !== undefined && object.erc20ToDenoms !== null) {
            for (var _h = 0, _j = object.erc20ToDenoms; _h < _j.length; _h++) {
                var e = _j[_h];
                message.erc20ToDenoms.push(exports.ERC20ToDenom.fromJSON(e));
            }
        }
        if (object.unbatchedSendToEthereumTxs !== undefined &&
            object.unbatchedSendToEthereumTxs !== null) {
            for (var _k = 0, _l = object.unbatchedSendToEthereumTxs; _k < _l.length; _k++) {
                var e = _l[_k];
                message.unbatchedSendToEthereumTxs.push(gravity_1.SendToEthereum.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.params !== undefined &&
            (obj.params = message.params ? exports.Params.toJSON(message.params) : undefined);
        message.lastObservedEventNonce !== undefined &&
            (obj.lastObservedEventNonce = (message.lastObservedEventNonce || long_1["default"].UZERO).toString());
        if (message.outgoingTxs) {
            obj.outgoingTxs = message.outgoingTxs.map(function (e) {
                return e ? any_1.Any.toJSON(e) : undefined;
            });
        }
        else {
            obj.outgoingTxs = [];
        }
        if (message.confirmations) {
            obj.confirmations = message.confirmations.map(function (e) {
                return e ? any_1.Any.toJSON(e) : undefined;
            });
        }
        else {
            obj.confirmations = [];
        }
        if (message.ethereumEventVoteRecords) {
            obj.ethereumEventVoteRecords = message.ethereumEventVoteRecords.map(function (e) {
                return e ? gravity_1.EthereumEventVoteRecord.toJSON(e) : undefined;
            });
        }
        else {
            obj.ethereumEventVoteRecords = [];
        }
        if (message.delegateKeys) {
            obj.delegateKeys = message.delegateKeys.map(function (e) {
                return e ? msgs_1.MsgDelegateKeys.toJSON(e) : undefined;
            });
        }
        else {
            obj.delegateKeys = [];
        }
        if (message.erc20ToDenoms) {
            obj.erc20ToDenoms = message.erc20ToDenoms.map(function (e) {
                return e ? exports.ERC20ToDenom.toJSON(e) : undefined;
            });
        }
        else {
            obj.erc20ToDenoms = [];
        }
        if (message.unbatchedSendToEthereumTxs) {
            obj.unbatchedSendToEthereumTxs = message.unbatchedSendToEthereumTxs.map(function (e) { return (e ? gravity_1.SendToEthereum.toJSON(e) : undefined); });
        }
        else {
            obj.unbatchedSendToEthereumTxs = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseGenesisState);
        message.outgoingTxs = [];
        message.confirmations = [];
        message.ethereumEventVoteRecords = [];
        message.delegateKeys = [];
        message.erc20ToDenoms = [];
        message.unbatchedSendToEthereumTxs = [];
        if (object.params !== undefined && object.params !== null) {
            message.params = exports.Params.fromPartial(object.params);
        }
        else {
            message.params = undefined;
        }
        if (object.lastObservedEventNonce !== undefined &&
            object.lastObservedEventNonce !== null) {
            message.lastObservedEventNonce = object.lastObservedEventNonce;
        }
        else {
            message.lastObservedEventNonce = long_1["default"].UZERO;
        }
        if (object.outgoingTxs !== undefined && object.outgoingTxs !== null) {
            for (var _i = 0, _a = object.outgoingTxs; _i < _a.length; _i++) {
                var e = _a[_i];
                message.outgoingTxs.push(any_1.Any.fromPartial(e));
            }
        }
        if (object.confirmations !== undefined && object.confirmations !== null) {
            for (var _b = 0, _c = object.confirmations; _b < _c.length; _b++) {
                var e = _c[_b];
                message.confirmations.push(any_1.Any.fromPartial(e));
            }
        }
        if (object.ethereumEventVoteRecords !== undefined &&
            object.ethereumEventVoteRecords !== null) {
            for (var _d = 0, _e = object.ethereumEventVoteRecords; _d < _e.length; _d++) {
                var e = _e[_d];
                message.ethereumEventVoteRecords.push(gravity_1.EthereumEventVoteRecord.fromPartial(e));
            }
        }
        if (object.delegateKeys !== undefined && object.delegateKeys !== null) {
            for (var _f = 0, _g = object.delegateKeys; _f < _g.length; _f++) {
                var e = _g[_f];
                message.delegateKeys.push(msgs_1.MsgDelegateKeys.fromPartial(e));
            }
        }
        if (object.erc20ToDenoms !== undefined && object.erc20ToDenoms !== null) {
            for (var _h = 0, _j = object.erc20ToDenoms; _h < _j.length; _h++) {
                var e = _j[_h];
                message.erc20ToDenoms.push(exports.ERC20ToDenom.fromPartial(e));
            }
        }
        if (object.unbatchedSendToEthereumTxs !== undefined &&
            object.unbatchedSendToEthereumTxs !== null) {
            for (var _k = 0, _l = object.unbatchedSendToEthereumTxs; _k < _l.length; _k++) {
                var e = _l[_k];
                message.unbatchedSendToEthereumTxs.push(gravity_1.SendToEthereum.fromPartial(e));
            }
        }
        return message;
    }
};
var baseERC20ToDenom = { erc20: "", denom: "" };
exports.ERC20ToDenom = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.erc20 !== "") {
            writer.uint32(10).string(message.erc20);
        }
        if (message.denom !== "") {
            writer.uint32(18).string(message.denom);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseERC20ToDenom);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.erc20 = reader.string();
                    break;
                case 2:
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
        var message = __assign({}, baseERC20ToDenom);
        if (object.erc20 !== undefined && object.erc20 !== null) {
            message.erc20 = String(object.erc20);
        }
        else {
            message.erc20 = "";
        }
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
        message.erc20 !== undefined && (obj.erc20 = message.erc20);
        message.denom !== undefined && (obj.denom = message.denom);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseERC20ToDenom);
        if (object.erc20 !== undefined && object.erc20 !== null) {
            message.erc20 = object.erc20;
        }
        else {
            message.erc20 = "";
        }
        if (object.denom !== undefined && object.denom !== null) {
            message.denom = object.denom;
        }
        else {
            message.denom = "";
        }
        return message;
    }
};
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
