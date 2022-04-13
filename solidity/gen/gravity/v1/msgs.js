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
exports.MsgClientImpl = exports.MsgDelegateKeysResponse = exports.MsgDelegateKeys = exports.MsgSubmitEthereumEventResponse = exports.SignerSetTxExecutedEvent = exports.ERC20DeployedEvent = exports.ContractCallExecutedEvent = exports.BatchExecutedEvent = exports.SendToCosmosEvent = exports.MsgSubmitEthereumEvent = exports.MsgSubmitEthereumTxConfirmationResponse = exports.SignerSetTxConfirmation = exports.BatchTxConfirmation = exports.ContractCallTxConfirmation = exports.MsgSubmitEthereumTxConfirmation = exports.MsgRequestBatchTxResponse = exports.MsgRequestBatchTx = exports.MsgCancelSendToEthereumResponse = exports.MsgCancelSendToEthereum = exports.MsgSendToEthereumResponse = exports.MsgSendToEthereum = exports.protobufPackage = void 0;
/* eslint-disable */
var long_1 = require("long");
var minimal_1 = require("protobufjs/minimal");
var coin_1 = require("../../cosmos/base/v1beta1/coin");
var any_1 = require("../../google/protobuf/any");
var gravity_1 = require("../../gravity/v1/gravity");
exports.protobufPackage = "gravity.v1";
var baseMsgSendToEthereum = { sender: "", ethereumRecipient: "" };
exports.MsgSendToEthereum = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.ethereumRecipient !== "") {
            writer.uint32(18).string(message.ethereumRecipient);
        }
        if (message.amount !== undefined) {
            coin_1.Coin.encode(message.amount, writer.uint32(26).fork()).ldelim();
        }
        if (message.bridgeFee !== undefined) {
            coin_1.Coin.encode(message.bridgeFee, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgSendToEthereum);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.ethereumRecipient = reader.string();
                    break;
                case 3:
                    message.amount = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.bridgeFee = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseMsgSendToEthereum);
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = String(object.sender);
        }
        else {
            message.sender = "";
        }
        if (object.ethereumRecipient !== undefined &&
            object.ethereumRecipient !== null) {
            message.ethereumRecipient = String(object.ethereumRecipient);
        }
        else {
            message.ethereumRecipient = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = coin_1.Coin.fromJSON(object.amount);
        }
        else {
            message.amount = undefined;
        }
        if (object.bridgeFee !== undefined && object.bridgeFee !== null) {
            message.bridgeFee = coin_1.Coin.fromJSON(object.bridgeFee);
        }
        else {
            message.bridgeFee = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.ethereumRecipient !== undefined &&
            (obj.ethereumRecipient = message.ethereumRecipient);
        message.amount !== undefined &&
            (obj.amount = message.amount ? coin_1.Coin.toJSON(message.amount) : undefined);
        message.bridgeFee !== undefined &&
            (obj.bridgeFee = message.bridgeFee
                ? coin_1.Coin.toJSON(message.bridgeFee)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseMsgSendToEthereum);
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = object.sender;
        }
        else {
            message.sender = "";
        }
        if (object.ethereumRecipient !== undefined &&
            object.ethereumRecipient !== null) {
            message.ethereumRecipient = object.ethereumRecipient;
        }
        else {
            message.ethereumRecipient = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = coin_1.Coin.fromPartial(object.amount);
        }
        else {
            message.amount = undefined;
        }
        if (object.bridgeFee !== undefined && object.bridgeFee !== null) {
            message.bridgeFee = coin_1.Coin.fromPartial(object.bridgeFee);
        }
        else {
            message.bridgeFee = undefined;
        }
        return message;
    }
};
var baseMsgSendToEthereumResponse = { id: long_1["default"].UZERO };
exports.MsgSendToEthereumResponse = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.id.isZero()) {
            writer.uint32(8).uint64(message.id);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgSendToEthereumResponse);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseMsgSendToEthereumResponse);
        if (object.id !== undefined && object.id !== null) {
            message.id = long_1["default"].fromString(object.id);
        }
        else {
            message.id = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.id !== undefined &&
            (obj.id = (message.id || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseMsgSendToEthereumResponse);
        if (object.id !== undefined && object.id !== null) {
            message.id = object.id;
        }
        else {
            message.id = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseMsgCancelSendToEthereum = { id: long_1["default"].UZERO, sender: "" };
exports.MsgCancelSendToEthereum = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.id.isZero()) {
            writer.uint32(8).uint64(message.id);
        }
        if (message.sender !== "") {
            writer.uint32(18).string(message.sender);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgCancelSendToEthereum);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.uint64();
                    break;
                case 2:
                    message.sender = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseMsgCancelSendToEthereum);
        if (object.id !== undefined && object.id !== null) {
            message.id = long_1["default"].fromString(object.id);
        }
        else {
            message.id = long_1["default"].UZERO;
        }
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = String(object.sender);
        }
        else {
            message.sender = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.id !== undefined &&
            (obj.id = (message.id || long_1["default"].UZERO).toString());
        message.sender !== undefined && (obj.sender = message.sender);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseMsgCancelSendToEthereum);
        if (object.id !== undefined && object.id !== null) {
            message.id = object.id;
        }
        else {
            message.id = long_1["default"].UZERO;
        }
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = object.sender;
        }
        else {
            message.sender = "";
        }
        return message;
    }
};
var baseMsgCancelSendToEthereumResponse = {};
exports.MsgCancelSendToEthereumResponse = {
    encode: function (_, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgCancelSendToEthereumResponse);
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
        var message = __assign({}, baseMsgCancelSendToEthereumResponse);
        return message;
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
    fromPartial: function (_) {
        var message = __assign({}, baseMsgCancelSendToEthereumResponse);
        return message;
    }
};
var baseMsgRequestBatchTx = { denom: "", signer: "" };
exports.MsgRequestBatchTx = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.denom !== "") {
            writer.uint32(10).string(message.denom);
        }
        if (message.signer !== "") {
            writer.uint32(18).string(message.signer);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgRequestBatchTx);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.denom = reader.string();
                    break;
                case 2:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseMsgRequestBatchTx);
        if (object.denom !== undefined && object.denom !== null) {
            message.denom = String(object.denom);
        }
        else {
            message.denom = "";
        }
        if (object.signer !== undefined && object.signer !== null) {
            message.signer = String(object.signer);
        }
        else {
            message.signer = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.denom !== undefined && (obj.denom = message.denom);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseMsgRequestBatchTx);
        if (object.denom !== undefined && object.denom !== null) {
            message.denom = object.denom;
        }
        else {
            message.denom = "";
        }
        if (object.signer !== undefined && object.signer !== null) {
            message.signer = object.signer;
        }
        else {
            message.signer = "";
        }
        return message;
    }
};
var baseMsgRequestBatchTxResponse = {};
exports.MsgRequestBatchTxResponse = {
    encode: function (_, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgRequestBatchTxResponse);
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
        var message = __assign({}, baseMsgRequestBatchTxResponse);
        return message;
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
    fromPartial: function (_) {
        var message = __assign({}, baseMsgRequestBatchTxResponse);
        return message;
    }
};
var baseMsgSubmitEthereumTxConfirmation = { signer: "" };
exports.MsgSubmitEthereumTxConfirmation = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.confirmation !== undefined) {
            any_1.Any.encode(message.confirmation, writer.uint32(10).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(18).string(message.signer);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgSubmitEthereumTxConfirmation);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.confirmation = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseMsgSubmitEthereumTxConfirmation);
        if (object.confirmation !== undefined && object.confirmation !== null) {
            message.confirmation = any_1.Any.fromJSON(object.confirmation);
        }
        else {
            message.confirmation = undefined;
        }
        if (object.signer !== undefined && object.signer !== null) {
            message.signer = String(object.signer);
        }
        else {
            message.signer = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.confirmation !== undefined &&
            (obj.confirmation = message.confirmation
                ? any_1.Any.toJSON(message.confirmation)
                : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseMsgSubmitEthereumTxConfirmation);
        if (object.confirmation !== undefined && object.confirmation !== null) {
            message.confirmation = any_1.Any.fromPartial(object.confirmation);
        }
        else {
            message.confirmation = undefined;
        }
        if (object.signer !== undefined && object.signer !== null) {
            message.signer = object.signer;
        }
        else {
            message.signer = "";
        }
        return message;
    }
};
var baseContractCallTxConfirmation = {
    invalidationNonce: long_1["default"].UZERO,
    ethereumSigner: ""
};
exports.ContractCallTxConfirmation = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.invalidationScope.length !== 0) {
            writer.uint32(10).bytes(message.invalidationScope);
        }
        if (!message.invalidationNonce.isZero()) {
            writer.uint32(16).uint64(message.invalidationNonce);
        }
        if (message.ethereumSigner !== "") {
            writer.uint32(26).string(message.ethereumSigner);
        }
        if (message.signature.length !== 0) {
            writer.uint32(34).bytes(message.signature);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseContractCallTxConfirmation);
        message.invalidationScope = new Uint8Array();
        message.signature = new Uint8Array();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.invalidationScope = reader.bytes();
                    break;
                case 2:
                    message.invalidationNonce = reader.uint64();
                    break;
                case 3:
                    message.ethereumSigner = reader.string();
                    break;
                case 4:
                    message.signature = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseContractCallTxConfirmation);
        message.invalidationScope = new Uint8Array();
        message.signature = new Uint8Array();
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
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = String(object.ethereumSigner);
        }
        else {
            message.ethereumSigner = "";
        }
        if (object.signature !== undefined && object.signature !== null) {
            message.signature = bytesFromBase64(object.signature);
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
        message.ethereumSigner !== undefined &&
            (obj.ethereumSigner = message.ethereumSigner);
        message.signature !== undefined &&
            (obj.signature = base64FromBytes(message.signature !== undefined ? message.signature : new Uint8Array()));
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseContractCallTxConfirmation);
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
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = object.ethereumSigner;
        }
        else {
            message.ethereumSigner = "";
        }
        if (object.signature !== undefined && object.signature !== null) {
            message.signature = object.signature;
        }
        else {
            message.signature = new Uint8Array();
        }
        return message;
    }
};
var baseBatchTxConfirmation = {
    tokenContract: "",
    batchNonce: long_1["default"].UZERO,
    ethereumSigner: ""
};
exports.BatchTxConfirmation = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.tokenContract !== "") {
            writer.uint32(10).string(message.tokenContract);
        }
        if (!message.batchNonce.isZero()) {
            writer.uint32(16).uint64(message.batchNonce);
        }
        if (message.ethereumSigner !== "") {
            writer.uint32(26).string(message.ethereumSigner);
        }
        if (message.signature.length !== 0) {
            writer.uint32(34).bytes(message.signature);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTxConfirmation);
        message.signature = new Uint8Array();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tokenContract = reader.string();
                    break;
                case 2:
                    message.batchNonce = reader.uint64();
                    break;
                case 3:
                    message.ethereumSigner = reader.string();
                    break;
                case 4:
                    message.signature = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchTxConfirmation);
        message.signature = new Uint8Array();
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
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = String(object.ethereumSigner);
        }
        else {
            message.ethereumSigner = "";
        }
        if (object.signature !== undefined && object.signature !== null) {
            message.signature = bytesFromBase64(object.signature);
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.tokenContract !== undefined &&
            (obj.tokenContract = message.tokenContract);
        message.batchNonce !== undefined &&
            (obj.batchNonce = (message.batchNonce || long_1["default"].UZERO).toString());
        message.ethereumSigner !== undefined &&
            (obj.ethereumSigner = message.ethereumSigner);
        message.signature !== undefined &&
            (obj.signature = base64FromBytes(message.signature !== undefined ? message.signature : new Uint8Array()));
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchTxConfirmation);
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
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = object.ethereumSigner;
        }
        else {
            message.ethereumSigner = "";
        }
        if (object.signature !== undefined && object.signature !== null) {
            message.signature = object.signature;
        }
        else {
            message.signature = new Uint8Array();
        }
        return message;
    }
};
var baseSignerSetTxConfirmation = {
    signerSetNonce: long_1["default"].UZERO,
    ethereumSigner: ""
};
exports.SignerSetTxConfirmation = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.signerSetNonce.isZero()) {
            writer.uint32(8).uint64(message.signerSetNonce);
        }
        if (message.ethereumSigner !== "") {
            writer.uint32(18).string(message.ethereumSigner);
        }
        if (message.signature.length !== 0) {
            writer.uint32(26).bytes(message.signature);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSignerSetTxConfirmation);
        message.signature = new Uint8Array();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signerSetNonce = reader.uint64();
                    break;
                case 2:
                    message.ethereumSigner = reader.string();
                    break;
                case 3:
                    message.signature = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSignerSetTxConfirmation);
        message.signature = new Uint8Array();
        if (object.signerSetNonce !== undefined && object.signerSetNonce !== null) {
            message.signerSetNonce = long_1["default"].fromString(object.signerSetNonce);
        }
        else {
            message.signerSetNonce = long_1["default"].UZERO;
        }
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = String(object.ethereumSigner);
        }
        else {
            message.ethereumSigner = "";
        }
        if (object.signature !== undefined && object.signature !== null) {
            message.signature = bytesFromBase64(object.signature);
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.signerSetNonce !== undefined &&
            (obj.signerSetNonce = (message.signerSetNonce || long_1["default"].UZERO).toString());
        message.ethereumSigner !== undefined &&
            (obj.ethereumSigner = message.ethereumSigner);
        message.signature !== undefined &&
            (obj.signature = base64FromBytes(message.signature !== undefined ? message.signature : new Uint8Array()));
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSignerSetTxConfirmation);
        if (object.signerSetNonce !== undefined && object.signerSetNonce !== null) {
            message.signerSetNonce = object.signerSetNonce;
        }
        else {
            message.signerSetNonce = long_1["default"].UZERO;
        }
        if (object.ethereumSigner !== undefined && object.ethereumSigner !== null) {
            message.ethereumSigner = object.ethereumSigner;
        }
        else {
            message.ethereumSigner = "";
        }
        if (object.signature !== undefined && object.signature !== null) {
            message.signature = object.signature;
        }
        else {
            message.signature = new Uint8Array();
        }
        return message;
    }
};
var baseMsgSubmitEthereumTxConfirmationResponse = {};
exports.MsgSubmitEthereumTxConfirmationResponse = {
    encode: function (_, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgSubmitEthereumTxConfirmationResponse);
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
        var message = __assign({}, baseMsgSubmitEthereumTxConfirmationResponse);
        return message;
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
    fromPartial: function (_) {
        var message = __assign({}, baseMsgSubmitEthereumTxConfirmationResponse);
        return message;
    }
};
var baseMsgSubmitEthereumEvent = { signer: "" };
exports.MsgSubmitEthereumEvent = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.event !== undefined) {
            any_1.Any.encode(message.event, writer.uint32(10).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(18).string(message.signer);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgSubmitEthereumEvent);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.event = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseMsgSubmitEthereumEvent);
        if (object.event !== undefined && object.event !== null) {
            message.event = any_1.Any.fromJSON(object.event);
        }
        else {
            message.event = undefined;
        }
        if (object.signer !== undefined && object.signer !== null) {
            message.signer = String(object.signer);
        }
        else {
            message.signer = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.event !== undefined &&
            (obj.event = message.event ? any_1.Any.toJSON(message.event) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseMsgSubmitEthereumEvent);
        if (object.event !== undefined && object.event !== null) {
            message.event = any_1.Any.fromPartial(object.event);
        }
        else {
            message.event = undefined;
        }
        if (object.signer !== undefined && object.signer !== null) {
            message.signer = object.signer;
        }
        else {
            message.signer = "";
        }
        return message;
    }
};
var baseSendToCosmosEvent = {
    eventNonce: long_1["default"].UZERO,
    tokenContract: "",
    amount: "",
    ethereumSender: "",
    cosmosReceiver: "",
    ethereumHeight: long_1["default"].UZERO
};
exports.SendToCosmosEvent = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.eventNonce.isZero()) {
            writer.uint32(8).uint64(message.eventNonce);
        }
        if (message.tokenContract !== "") {
            writer.uint32(18).string(message.tokenContract);
        }
        if (message.amount !== "") {
            writer.uint32(26).string(message.amount);
        }
        if (message.ethereumSender !== "") {
            writer.uint32(34).string(message.ethereumSender);
        }
        if (message.cosmosReceiver !== "") {
            writer.uint32(42).string(message.cosmosReceiver);
        }
        if (!message.ethereumHeight.isZero()) {
            writer.uint32(48).uint64(message.ethereumHeight);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSendToCosmosEvent);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.eventNonce = reader.uint64();
                    break;
                case 2:
                    message.tokenContract = reader.string();
                    break;
                case 3:
                    message.amount = reader.string();
                    break;
                case 4:
                    message.ethereumSender = reader.string();
                    break;
                case 5:
                    message.cosmosReceiver = reader.string();
                    break;
                case 6:
                    message.ethereumHeight = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSendToCosmosEvent);
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = long_1["default"].fromString(object.eventNonce);
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = String(object.tokenContract);
        }
        else {
            message.tokenContract = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = String(object.amount);
        }
        else {
            message.amount = "";
        }
        if (object.ethereumSender !== undefined && object.ethereumSender !== null) {
            message.ethereumSender = String(object.ethereumSender);
        }
        else {
            message.ethereumSender = "";
        }
        if (object.cosmosReceiver !== undefined && object.cosmosReceiver !== null) {
            message.cosmosReceiver = String(object.cosmosReceiver);
        }
        else {
            message.cosmosReceiver = "";
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = long_1["default"].fromString(object.ethereumHeight);
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.eventNonce !== undefined &&
            (obj.eventNonce = (message.eventNonce || long_1["default"].UZERO).toString());
        message.tokenContract !== undefined &&
            (obj.tokenContract = message.tokenContract);
        message.amount !== undefined && (obj.amount = message.amount);
        message.ethereumSender !== undefined &&
            (obj.ethereumSender = message.ethereumSender);
        message.cosmosReceiver !== undefined &&
            (obj.cosmosReceiver = message.cosmosReceiver);
        message.ethereumHeight !== undefined &&
            (obj.ethereumHeight = (message.ethereumHeight || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSendToCosmosEvent);
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = object.eventNonce;
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = object.tokenContract;
        }
        else {
            message.tokenContract = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = object.amount;
        }
        else {
            message.amount = "";
        }
        if (object.ethereumSender !== undefined && object.ethereumSender !== null) {
            message.ethereumSender = object.ethereumSender;
        }
        else {
            message.ethereumSender = "";
        }
        if (object.cosmosReceiver !== undefined && object.cosmosReceiver !== null) {
            message.cosmosReceiver = object.cosmosReceiver;
        }
        else {
            message.cosmosReceiver = "";
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = object.ethereumHeight;
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseBatchExecutedEvent = {
    tokenContract: "",
    eventNonce: long_1["default"].UZERO,
    ethereumHeight: long_1["default"].UZERO,
    batchNonce: long_1["default"].UZERO
};
exports.BatchExecutedEvent = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.tokenContract !== "") {
            writer.uint32(10).string(message.tokenContract);
        }
        if (!message.eventNonce.isZero()) {
            writer.uint32(16).uint64(message.eventNonce);
        }
        if (!message.ethereumHeight.isZero()) {
            writer.uint32(24).uint64(message.ethereumHeight);
        }
        if (!message.batchNonce.isZero()) {
            writer.uint32(32).uint64(message.batchNonce);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchExecutedEvent);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tokenContract = reader.string();
                    break;
                case 2:
                    message.eventNonce = reader.uint64();
                    break;
                case 3:
                    message.ethereumHeight = reader.uint64();
                    break;
                case 4:
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
        var message = __assign({}, baseBatchExecutedEvent);
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = String(object.tokenContract);
        }
        else {
            message.tokenContract = "";
        }
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = long_1["default"].fromString(object.eventNonce);
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = long_1["default"].fromString(object.ethereumHeight);
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
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
        message.eventNonce !== undefined &&
            (obj.eventNonce = (message.eventNonce || long_1["default"].UZERO).toString());
        message.ethereumHeight !== undefined &&
            (obj.ethereumHeight = (message.ethereumHeight || long_1["default"].UZERO).toString());
        message.batchNonce !== undefined &&
            (obj.batchNonce = (message.batchNonce || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchExecutedEvent);
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = object.tokenContract;
        }
        else {
            message.tokenContract = "";
        }
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = object.eventNonce;
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = object.ethereumHeight;
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
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
var baseContractCallExecutedEvent = {
    eventNonce: long_1["default"].UZERO,
    invalidationNonce: long_1["default"].UZERO,
    ethereumHeight: long_1["default"].UZERO
};
exports.ContractCallExecutedEvent = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.eventNonce.isZero()) {
            writer.uint32(8).uint64(message.eventNonce);
        }
        if (message.invalidationId.length !== 0) {
            writer.uint32(18).bytes(message.invalidationId);
        }
        if (!message.invalidationNonce.isZero()) {
            writer.uint32(24).uint64(message.invalidationNonce);
        }
        if (!message.ethereumHeight.isZero()) {
            writer.uint32(32).uint64(message.ethereumHeight);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseContractCallExecutedEvent);
        message.invalidationId = new Uint8Array();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.eventNonce = reader.uint64();
                    break;
                case 2:
                    message.invalidationId = reader.bytes();
                    break;
                case 3:
                    message.invalidationNonce = reader.uint64();
                    break;
                case 4:
                    message.ethereumHeight = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseContractCallExecutedEvent);
        message.invalidationId = new Uint8Array();
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = long_1["default"].fromString(object.eventNonce);
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.invalidationId !== undefined && object.invalidationId !== null) {
            message.invalidationId = bytesFromBase64(object.invalidationId);
        }
        if (object.invalidationNonce !== undefined &&
            object.invalidationNonce !== null) {
            message.invalidationNonce = long_1["default"].fromString(object.invalidationNonce);
        }
        else {
            message.invalidationNonce = long_1["default"].UZERO;
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = long_1["default"].fromString(object.ethereumHeight);
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.eventNonce !== undefined &&
            (obj.eventNonce = (message.eventNonce || long_1["default"].UZERO).toString());
        message.invalidationId !== undefined &&
            (obj.invalidationId = base64FromBytes(message.invalidationId !== undefined
                ? message.invalidationId
                : new Uint8Array()));
        message.invalidationNonce !== undefined &&
            (obj.invalidationNonce = (message.invalidationNonce || long_1["default"].UZERO).toString());
        message.ethereumHeight !== undefined &&
            (obj.ethereumHeight = (message.ethereumHeight || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseContractCallExecutedEvent);
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = object.eventNonce;
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.invalidationId !== undefined && object.invalidationId !== null) {
            message.invalidationId = object.invalidationId;
        }
        else {
            message.invalidationId = new Uint8Array();
        }
        if (object.invalidationNonce !== undefined &&
            object.invalidationNonce !== null) {
            message.invalidationNonce = object.invalidationNonce;
        }
        else {
            message.invalidationNonce = long_1["default"].UZERO;
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = object.ethereumHeight;
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseERC20DeployedEvent = {
    eventNonce: long_1["default"].UZERO,
    cosmosDenom: "",
    tokenContract: "",
    erc20Name: "",
    erc20Symbol: "",
    erc20Decimals: long_1["default"].UZERO,
    ethereumHeight: long_1["default"].UZERO
};
exports.ERC20DeployedEvent = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.eventNonce.isZero()) {
            writer.uint32(8).uint64(message.eventNonce);
        }
        if (message.cosmosDenom !== "") {
            writer.uint32(18).string(message.cosmosDenom);
        }
        if (message.tokenContract !== "") {
            writer.uint32(26).string(message.tokenContract);
        }
        if (message.erc20Name !== "") {
            writer.uint32(34).string(message.erc20Name);
        }
        if (message.erc20Symbol !== "") {
            writer.uint32(42).string(message.erc20Symbol);
        }
        if (!message.erc20Decimals.isZero()) {
            writer.uint32(48).uint64(message.erc20Decimals);
        }
        if (!message.ethereumHeight.isZero()) {
            writer.uint32(56).uint64(message.ethereumHeight);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseERC20DeployedEvent);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.eventNonce = reader.uint64();
                    break;
                case 2:
                    message.cosmosDenom = reader.string();
                    break;
                case 3:
                    message.tokenContract = reader.string();
                    break;
                case 4:
                    message.erc20Name = reader.string();
                    break;
                case 5:
                    message.erc20Symbol = reader.string();
                    break;
                case 6:
                    message.erc20Decimals = reader.uint64();
                    break;
                case 7:
                    message.ethereumHeight = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseERC20DeployedEvent);
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = long_1["default"].fromString(object.eventNonce);
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.cosmosDenom !== undefined && object.cosmosDenom !== null) {
            message.cosmosDenom = String(object.cosmosDenom);
        }
        else {
            message.cosmosDenom = "";
        }
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = String(object.tokenContract);
        }
        else {
            message.tokenContract = "";
        }
        if (object.erc20Name !== undefined && object.erc20Name !== null) {
            message.erc20Name = String(object.erc20Name);
        }
        else {
            message.erc20Name = "";
        }
        if (object.erc20Symbol !== undefined && object.erc20Symbol !== null) {
            message.erc20Symbol = String(object.erc20Symbol);
        }
        else {
            message.erc20Symbol = "";
        }
        if (object.erc20Decimals !== undefined && object.erc20Decimals !== null) {
            message.erc20Decimals = long_1["default"].fromString(object.erc20Decimals);
        }
        else {
            message.erc20Decimals = long_1["default"].UZERO;
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = long_1["default"].fromString(object.ethereumHeight);
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.eventNonce !== undefined &&
            (obj.eventNonce = (message.eventNonce || long_1["default"].UZERO).toString());
        message.cosmosDenom !== undefined &&
            (obj.cosmosDenom = message.cosmosDenom);
        message.tokenContract !== undefined &&
            (obj.tokenContract = message.tokenContract);
        message.erc20Name !== undefined && (obj.erc20Name = message.erc20Name);
        message.erc20Symbol !== undefined &&
            (obj.erc20Symbol = message.erc20Symbol);
        message.erc20Decimals !== undefined &&
            (obj.erc20Decimals = (message.erc20Decimals || long_1["default"].UZERO).toString());
        message.ethereumHeight !== undefined &&
            (obj.ethereumHeight = (message.ethereumHeight || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseERC20DeployedEvent);
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = object.eventNonce;
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.cosmosDenom !== undefined && object.cosmosDenom !== null) {
            message.cosmosDenom = object.cosmosDenom;
        }
        else {
            message.cosmosDenom = "";
        }
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = object.tokenContract;
        }
        else {
            message.tokenContract = "";
        }
        if (object.erc20Name !== undefined && object.erc20Name !== null) {
            message.erc20Name = object.erc20Name;
        }
        else {
            message.erc20Name = "";
        }
        if (object.erc20Symbol !== undefined && object.erc20Symbol !== null) {
            message.erc20Symbol = object.erc20Symbol;
        }
        else {
            message.erc20Symbol = "";
        }
        if (object.erc20Decimals !== undefined && object.erc20Decimals !== null) {
            message.erc20Decimals = object.erc20Decimals;
        }
        else {
            message.erc20Decimals = long_1["default"].UZERO;
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = object.ethereumHeight;
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseSignerSetTxExecutedEvent = {
    eventNonce: long_1["default"].UZERO,
    signerSetTxNonce: long_1["default"].UZERO,
    ethereumHeight: long_1["default"].UZERO
};
exports.SignerSetTxExecutedEvent = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.eventNonce.isZero()) {
            writer.uint32(8).uint64(message.eventNonce);
        }
        if (!message.signerSetTxNonce.isZero()) {
            writer.uint32(16).uint64(message.signerSetTxNonce);
        }
        if (!message.ethereumHeight.isZero()) {
            writer.uint32(24).uint64(message.ethereumHeight);
        }
        for (var _i = 0, _a = message.members; _i < _a.length; _i++) {
            var v = _a[_i];
            gravity_1.EthereumSigner.encode(v, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSignerSetTxExecutedEvent);
        message.members = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.eventNonce = reader.uint64();
                    break;
                case 2:
                    message.signerSetTxNonce = reader.uint64();
                    break;
                case 3:
                    message.ethereumHeight = reader.uint64();
                    break;
                case 4:
                    message.members.push(gravity_1.EthereumSigner.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSignerSetTxExecutedEvent);
        message.members = [];
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = long_1["default"].fromString(object.eventNonce);
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.signerSetTxNonce !== undefined &&
            object.signerSetTxNonce !== null) {
            message.signerSetTxNonce = long_1["default"].fromString(object.signerSetTxNonce);
        }
        else {
            message.signerSetTxNonce = long_1["default"].UZERO;
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = long_1["default"].fromString(object.ethereumHeight);
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        if (object.members !== undefined && object.members !== null) {
            for (var _i = 0, _a = object.members; _i < _a.length; _i++) {
                var e = _a[_i];
                message.members.push(gravity_1.EthereumSigner.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.eventNonce !== undefined &&
            (obj.eventNonce = (message.eventNonce || long_1["default"].UZERO).toString());
        message.signerSetTxNonce !== undefined &&
            (obj.signerSetTxNonce = (message.signerSetTxNonce || long_1["default"].UZERO).toString());
        message.ethereumHeight !== undefined &&
            (obj.ethereumHeight = (message.ethereumHeight || long_1["default"].UZERO).toString());
        if (message.members) {
            obj.members = message.members.map(function (e) {
                return e ? gravity_1.EthereumSigner.toJSON(e) : undefined;
            });
        }
        else {
            obj.members = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSignerSetTxExecutedEvent);
        message.members = [];
        if (object.eventNonce !== undefined && object.eventNonce !== null) {
            message.eventNonce = object.eventNonce;
        }
        else {
            message.eventNonce = long_1["default"].UZERO;
        }
        if (object.signerSetTxNonce !== undefined &&
            object.signerSetTxNonce !== null) {
            message.signerSetTxNonce = object.signerSetTxNonce;
        }
        else {
            message.signerSetTxNonce = long_1["default"].UZERO;
        }
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = object.ethereumHeight;
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        if (object.members !== undefined && object.members !== null) {
            for (var _i = 0, _a = object.members; _i < _a.length; _i++) {
                var e = _a[_i];
                message.members.push(gravity_1.EthereumSigner.fromPartial(e));
            }
        }
        return message;
    }
};
var baseMsgSubmitEthereumEventResponse = {};
exports.MsgSubmitEthereumEventResponse = {
    encode: function (_, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgSubmitEthereumEventResponse);
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
        var message = __assign({}, baseMsgSubmitEthereumEventResponse);
        return message;
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
    fromPartial: function (_) {
        var message = __assign({}, baseMsgSubmitEthereumEventResponse);
        return message;
    }
};
var baseMsgDelegateKeys = {
    validatorAddress: "",
    orchestratorAddress: "",
    ethereumAddress: ""
};
exports.MsgDelegateKeys = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.validatorAddress !== "") {
            writer.uint32(10).string(message.validatorAddress);
        }
        if (message.orchestratorAddress !== "") {
            writer.uint32(18).string(message.orchestratorAddress);
        }
        if (message.ethereumAddress !== "") {
            writer.uint32(26).string(message.ethereumAddress);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgDelegateKeys);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddress = reader.string();
                    break;
                case 2:
                    message.orchestratorAddress = reader.string();
                    break;
                case 3:
                    message.ethereumAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseMsgDelegateKeys);
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
        if (object.ethereumAddress !== undefined &&
            object.ethereumAddress !== null) {
            message.ethereumAddress = String(object.ethereumAddress);
        }
        else {
            message.ethereumAddress = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.validatorAddress !== undefined &&
            (obj.validatorAddress = message.validatorAddress);
        message.orchestratorAddress !== undefined &&
            (obj.orchestratorAddress = message.orchestratorAddress);
        message.ethereumAddress !== undefined &&
            (obj.ethereumAddress = message.ethereumAddress);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseMsgDelegateKeys);
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
        if (object.ethereumAddress !== undefined &&
            object.ethereumAddress !== null) {
            message.ethereumAddress = object.ethereumAddress;
        }
        else {
            message.ethereumAddress = "";
        }
        return message;
    }
};
var baseMsgDelegateKeysResponse = {};
exports.MsgDelegateKeysResponse = {
    encode: function (_, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseMsgDelegateKeysResponse);
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
        var message = __assign({}, baseMsgDelegateKeysResponse);
        return message;
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
    fromPartial: function (_) {
        var message = __assign({}, baseMsgDelegateKeysResponse);
        return message;
    }
};
var MsgClientImpl = /** @class */ (function () {
    function MsgClientImpl(rpc) {
        this.rpc = rpc;
        this.SendToEthereum = this.SendToEthereum.bind(this);
        this.CancelSendToEthereum = this.CancelSendToEthereum.bind(this);
        this.RequestBatchTx = this.RequestBatchTx.bind(this);
        this.SubmitEthereumTxConfirmation = this.SubmitEthereumTxConfirmation.bind(this);
        this.SubmitEthereumEvent = this.SubmitEthereumEvent.bind(this);
        this.SetDelegateKeys = this.SetDelegateKeys.bind(this);
    }
    MsgClientImpl.prototype.SendToEthereum = function (request) {
        var data = exports.MsgSendToEthereum.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Msg", "SendToEthereum", data);
        return promise.then(function (data) {
            return exports.MsgSendToEthereumResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    MsgClientImpl.prototype.CancelSendToEthereum = function (request) {
        var data = exports.MsgCancelSendToEthereum.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Msg", "CancelSendToEthereum", data);
        return promise.then(function (data) {
            return exports.MsgCancelSendToEthereumResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    MsgClientImpl.prototype.RequestBatchTx = function (request) {
        var data = exports.MsgRequestBatchTx.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Msg", "RequestBatchTx", data);
        return promise.then(function (data) {
            return exports.MsgRequestBatchTxResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    MsgClientImpl.prototype.SubmitEthereumTxConfirmation = function (request) {
        var data = exports.MsgSubmitEthereumTxConfirmation.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Msg", "SubmitEthereumTxConfirmation", data);
        return promise.then(function (data) {
            return exports.MsgSubmitEthereumTxConfirmationResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    MsgClientImpl.prototype.SubmitEthereumEvent = function (request) {
        var data = exports.MsgSubmitEthereumEvent.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Msg", "SubmitEthereumEvent", data);
        return promise.then(function (data) {
            return exports.MsgSubmitEthereumEventResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    MsgClientImpl.prototype.SetDelegateKeys = function (request) {
        var data = exports.MsgDelegateKeys.encode(request).finish();
        var promise = this.rpc.request("gravity.v1.Msg", "SetDelegateKeys", data);
        return promise.then(function (data) {
            return exports.MsgDelegateKeysResponse.decode(new minimal_1["default"].Reader(data));
        });
    };
    return MsgClientImpl;
}());
exports.MsgClientImpl = MsgClientImpl;
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
