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
exports.IDSet = exports.ERC20Token = exports.ContractCallTx = exports.SendToEthereum = exports.BatchTx = exports.SignerSetTx = exports.EthereumSigner = exports.LatestEthereumBlockHeight = exports.EthereumEventVoteRecord = exports.protobufPackage = void 0;
/* eslint-disable */
var long_1 = require("long");
var minimal_1 = require("protobufjs/minimal");
var any_1 = require("../../google/protobuf/any");
exports.protobufPackage = "gravity.v1";
var baseEthereumEventVoteRecord = { votes: "", accepted: false };
exports.EthereumEventVoteRecord = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.event !== undefined) {
            any_1.Any.encode(message.event, writer.uint32(10).fork()).ldelim();
        }
        for (var _i = 0, _a = message.votes; _i < _a.length; _i++) {
            var v = _a[_i];
            writer.uint32(18).string(v);
        }
        if (message.accepted === true) {
            writer.uint32(24).bool(message.accepted);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseEthereumEventVoteRecord);
        message.votes = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.event = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.votes.push(reader.string());
                    break;
                case 3:
                    message.accepted = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseEthereumEventVoteRecord);
        message.votes = [];
        if (object.event !== undefined && object.event !== null) {
            message.event = any_1.Any.fromJSON(object.event);
        }
        else {
            message.event = undefined;
        }
        if (object.votes !== undefined && object.votes !== null) {
            for (var _i = 0, _a = object.votes; _i < _a.length; _i++) {
                var e = _a[_i];
                message.votes.push(String(e));
            }
        }
        if (object.accepted !== undefined && object.accepted !== null) {
            message.accepted = Boolean(object.accepted);
        }
        else {
            message.accepted = false;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.event !== undefined &&
            (obj.event = message.event ? any_1.Any.toJSON(message.event) : undefined);
        if (message.votes) {
            obj.votes = message.votes.map(function (e) { return e; });
        }
        else {
            obj.votes = [];
        }
        message.accepted !== undefined && (obj.accepted = message.accepted);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseEthereumEventVoteRecord);
        message.votes = [];
        if (object.event !== undefined && object.event !== null) {
            message.event = any_1.Any.fromPartial(object.event);
        }
        else {
            message.event = undefined;
        }
        if (object.votes !== undefined && object.votes !== null) {
            for (var _i = 0, _a = object.votes; _i < _a.length; _i++) {
                var e = _a[_i];
                message.votes.push(e);
            }
        }
        if (object.accepted !== undefined && object.accepted !== null) {
            message.accepted = object.accepted;
        }
        else {
            message.accepted = false;
        }
        return message;
    }
};
var baseLatestEthereumBlockHeight = {
    ethereumHeight: long_1["default"].UZERO,
    cosmosHeight: long_1["default"].UZERO
};
exports.LatestEthereumBlockHeight = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.ethereumHeight.isZero()) {
            writer.uint32(8).uint64(message.ethereumHeight);
        }
        if (!message.cosmosHeight.isZero()) {
            writer.uint32(16).uint64(message.cosmosHeight);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseLatestEthereumBlockHeight);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.ethereumHeight = reader.uint64();
                    break;
                case 2:
                    message.cosmosHeight = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseLatestEthereumBlockHeight);
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = long_1["default"].fromString(object.ethereumHeight);
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        if (object.cosmosHeight !== undefined && object.cosmosHeight !== null) {
            message.cosmosHeight = long_1["default"].fromString(object.cosmosHeight);
        }
        else {
            message.cosmosHeight = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.ethereumHeight !== undefined &&
            (obj.ethereumHeight = (message.ethereumHeight || long_1["default"].UZERO).toString());
        message.cosmosHeight !== undefined &&
            (obj.cosmosHeight = (message.cosmosHeight || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseLatestEthereumBlockHeight);
        if (object.ethereumHeight !== undefined && object.ethereumHeight !== null) {
            message.ethereumHeight = object.ethereumHeight;
        }
        else {
            message.ethereumHeight = long_1["default"].UZERO;
        }
        if (object.cosmosHeight !== undefined && object.cosmosHeight !== null) {
            message.cosmosHeight = object.cosmosHeight;
        }
        else {
            message.cosmosHeight = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseEthereumSigner = { power: long_1["default"].UZERO, ethereumAddress: "" };
exports.EthereumSigner = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.power.isZero()) {
            writer.uint32(8).uint64(message.power);
        }
        if (message.ethereumAddress !== "") {
            writer.uint32(18).string(message.ethereumAddress);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseEthereumSigner);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.power = reader.uint64();
                    break;
                case 2:
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
        var message = __assign({}, baseEthereumSigner);
        if (object.power !== undefined && object.power !== null) {
            message.power = long_1["default"].fromString(object.power);
        }
        else {
            message.power = long_1["default"].UZERO;
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
        message.power !== undefined &&
            (obj.power = (message.power || long_1["default"].UZERO).toString());
        message.ethereumAddress !== undefined &&
            (obj.ethereumAddress = message.ethereumAddress);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseEthereumSigner);
        if (object.power !== undefined && object.power !== null) {
            message.power = object.power;
        }
        else {
            message.power = long_1["default"].UZERO;
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
var baseSignerSetTx = { nonce: long_1["default"].UZERO, height: long_1["default"].UZERO };
exports.SignerSetTx = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.nonce.isZero()) {
            writer.uint32(8).uint64(message.nonce);
        }
        if (!message.height.isZero()) {
            writer.uint32(16).uint64(message.height);
        }
        for (var _i = 0, _a = message.signers; _i < _a.length; _i++) {
            var v = _a[_i];
            exports.EthereumSigner.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSignerSetTx);
        message.signers = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.nonce = reader.uint64();
                    break;
                case 2:
                    message.height = reader.uint64();
                    break;
                case 3:
                    message.signers.push(exports.EthereumSigner.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSignerSetTx);
        message.signers = [];
        if (object.nonce !== undefined && object.nonce !== null) {
            message.nonce = long_1["default"].fromString(object.nonce);
        }
        else {
            message.nonce = long_1["default"].UZERO;
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = long_1["default"].fromString(object.height);
        }
        else {
            message.height = long_1["default"].UZERO;
        }
        if (object.signers !== undefined && object.signers !== null) {
            for (var _i = 0, _a = object.signers; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signers.push(exports.EthereumSigner.fromJSON(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.nonce !== undefined &&
            (obj.nonce = (message.nonce || long_1["default"].UZERO).toString());
        message.height !== undefined &&
            (obj.height = (message.height || long_1["default"].UZERO).toString());
        if (message.signers) {
            obj.signers = message.signers.map(function (e) {
                return e ? exports.EthereumSigner.toJSON(e) : undefined;
            });
        }
        else {
            obj.signers = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSignerSetTx);
        message.signers = [];
        if (object.nonce !== undefined && object.nonce !== null) {
            message.nonce = object.nonce;
        }
        else {
            message.nonce = long_1["default"].UZERO;
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = object.height;
        }
        else {
            message.height = long_1["default"].UZERO;
        }
        if (object.signers !== undefined && object.signers !== null) {
            for (var _i = 0, _a = object.signers; _i < _a.length; _i++) {
                var e = _a[_i];
                message.signers.push(exports.EthereumSigner.fromPartial(e));
            }
        }
        return message;
    }
};
var baseBatchTx = {
    batchNonce: long_1["default"].UZERO,
    timeout: long_1["default"].UZERO,
    tokenContract: "",
    height: long_1["default"].UZERO
};
exports.BatchTx = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.batchNonce.isZero()) {
            writer.uint32(8).uint64(message.batchNonce);
        }
        if (!message.timeout.isZero()) {
            writer.uint32(16).uint64(message.timeout);
        }
        for (var _i = 0, _a = message.transactions; _i < _a.length; _i++) {
            var v = _a[_i];
            exports.SendToEthereum.encode(v, writer.uint32(26).fork()).ldelim();
        }
        if (message.tokenContract !== "") {
            writer.uint32(34).string(message.tokenContract);
        }
        if (!message.height.isZero()) {
            writer.uint32(40).uint64(message.height);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseBatchTx);
        message.transactions = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.batchNonce = reader.uint64();
                    break;
                case 2:
                    message.timeout = reader.uint64();
                    break;
                case 3:
                    message.transactions.push(exports.SendToEthereum.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.tokenContract = reader.string();
                    break;
                case 5:
                    message.height = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseBatchTx);
        message.transactions = [];
        if (object.batchNonce !== undefined && object.batchNonce !== null) {
            message.batchNonce = long_1["default"].fromString(object.batchNonce);
        }
        else {
            message.batchNonce = long_1["default"].UZERO;
        }
        if (object.timeout !== undefined && object.timeout !== null) {
            message.timeout = long_1["default"].fromString(object.timeout);
        }
        else {
            message.timeout = long_1["default"].UZERO;
        }
        if (object.transactions !== undefined && object.transactions !== null) {
            for (var _i = 0, _a = object.transactions; _i < _a.length; _i++) {
                var e = _a[_i];
                message.transactions.push(exports.SendToEthereum.fromJSON(e));
            }
        }
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = String(object.tokenContract);
        }
        else {
            message.tokenContract = "";
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = long_1["default"].fromString(object.height);
        }
        else {
            message.height = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.batchNonce !== undefined &&
            (obj.batchNonce = (message.batchNonce || long_1["default"].UZERO).toString());
        message.timeout !== undefined &&
            (obj.timeout = (message.timeout || long_1["default"].UZERO).toString());
        if (message.transactions) {
            obj.transactions = message.transactions.map(function (e) {
                return e ? exports.SendToEthereum.toJSON(e) : undefined;
            });
        }
        else {
            obj.transactions = [];
        }
        message.tokenContract !== undefined &&
            (obj.tokenContract = message.tokenContract);
        message.height !== undefined &&
            (obj.height = (message.height || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseBatchTx);
        message.transactions = [];
        if (object.batchNonce !== undefined && object.batchNonce !== null) {
            message.batchNonce = object.batchNonce;
        }
        else {
            message.batchNonce = long_1["default"].UZERO;
        }
        if (object.timeout !== undefined && object.timeout !== null) {
            message.timeout = object.timeout;
        }
        else {
            message.timeout = long_1["default"].UZERO;
        }
        if (object.transactions !== undefined && object.transactions !== null) {
            for (var _i = 0, _a = object.transactions; _i < _a.length; _i++) {
                var e = _a[_i];
                message.transactions.push(exports.SendToEthereum.fromPartial(e));
            }
        }
        if (object.tokenContract !== undefined && object.tokenContract !== null) {
            message.tokenContract = object.tokenContract;
        }
        else {
            message.tokenContract = "";
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = object.height;
        }
        else {
            message.height = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseSendToEthereum = {
    id: long_1["default"].UZERO,
    sender: "",
    ethereumRecipient: ""
};
exports.SendToEthereum = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.id.isZero()) {
            writer.uint32(8).uint64(message.id);
        }
        if (message.sender !== "") {
            writer.uint32(18).string(message.sender);
        }
        if (message.ethereumRecipient !== "") {
            writer.uint32(26).string(message.ethereumRecipient);
        }
        if (message.erc20Token !== undefined) {
            exports.ERC20Token.encode(message.erc20Token, writer.uint32(34).fork()).ldelim();
        }
        if (message.erc20Fee !== undefined) {
            exports.ERC20Token.encode(message.erc20Fee, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseSendToEthereum);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.uint64();
                    break;
                case 2:
                    message.sender = reader.string();
                    break;
                case 3:
                    message.ethereumRecipient = reader.string();
                    break;
                case 4:
                    message.erc20Token = exports.ERC20Token.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.erc20Fee = exports.ERC20Token.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseSendToEthereum);
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
        if (object.ethereumRecipient !== undefined &&
            object.ethereumRecipient !== null) {
            message.ethereumRecipient = String(object.ethereumRecipient);
        }
        else {
            message.ethereumRecipient = "";
        }
        if (object.erc20Token !== undefined && object.erc20Token !== null) {
            message.erc20Token = exports.ERC20Token.fromJSON(object.erc20Token);
        }
        else {
            message.erc20Token = undefined;
        }
        if (object.erc20Fee !== undefined && object.erc20Fee !== null) {
            message.erc20Fee = exports.ERC20Token.fromJSON(object.erc20Fee);
        }
        else {
            message.erc20Fee = undefined;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.id !== undefined &&
            (obj.id = (message.id || long_1["default"].UZERO).toString());
        message.sender !== undefined && (obj.sender = message.sender);
        message.ethereumRecipient !== undefined &&
            (obj.ethereumRecipient = message.ethereumRecipient);
        message.erc20Token !== undefined &&
            (obj.erc20Token = message.erc20Token
                ? exports.ERC20Token.toJSON(message.erc20Token)
                : undefined);
        message.erc20Fee !== undefined &&
            (obj.erc20Fee = message.erc20Fee
                ? exports.ERC20Token.toJSON(message.erc20Fee)
                : undefined);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseSendToEthereum);
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
        if (object.ethereumRecipient !== undefined &&
            object.ethereumRecipient !== null) {
            message.ethereumRecipient = object.ethereumRecipient;
        }
        else {
            message.ethereumRecipient = "";
        }
        if (object.erc20Token !== undefined && object.erc20Token !== null) {
            message.erc20Token = exports.ERC20Token.fromPartial(object.erc20Token);
        }
        else {
            message.erc20Token = undefined;
        }
        if (object.erc20Fee !== undefined && object.erc20Fee !== null) {
            message.erc20Fee = exports.ERC20Token.fromPartial(object.erc20Fee);
        }
        else {
            message.erc20Fee = undefined;
        }
        return message;
    }
};
var baseContractCallTx = {
    invalidationNonce: long_1["default"].UZERO,
    address: "",
    timeout: long_1["default"].UZERO,
    height: long_1["default"].UZERO
};
exports.ContractCallTx = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (!message.invalidationNonce.isZero()) {
            writer.uint32(8).uint64(message.invalidationNonce);
        }
        if (message.invalidationScope.length !== 0) {
            writer.uint32(18).bytes(message.invalidationScope);
        }
        if (message.address !== "") {
            writer.uint32(26).string(message.address);
        }
        if (message.payload.length !== 0) {
            writer.uint32(34).bytes(message.payload);
        }
        if (!message.timeout.isZero()) {
            writer.uint32(40).uint64(message.timeout);
        }
        for (var _i = 0, _a = message.tokens; _i < _a.length; _i++) {
            var v = _a[_i];
            exports.ERC20Token.encode(v, writer.uint32(50).fork()).ldelim();
        }
        for (var _b = 0, _c = message.fees; _b < _c.length; _b++) {
            var v = _c[_b];
            exports.ERC20Token.encode(v, writer.uint32(58).fork()).ldelim();
        }
        if (!message.height.isZero()) {
            writer.uint32(64).uint64(message.height);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseContractCallTx);
        message.tokens = [];
        message.fees = [];
        message.invalidationScope = new Uint8Array();
        message.payload = new Uint8Array();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.invalidationNonce = reader.uint64();
                    break;
                case 2:
                    message.invalidationScope = reader.bytes();
                    break;
                case 3:
                    message.address = reader.string();
                    break;
                case 4:
                    message.payload = reader.bytes();
                    break;
                case 5:
                    message.timeout = reader.uint64();
                    break;
                case 6:
                    message.tokens.push(exports.ERC20Token.decode(reader, reader.uint32()));
                    break;
                case 7:
                    message.fees.push(exports.ERC20Token.decode(reader, reader.uint32()));
                    break;
                case 8:
                    message.height = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseContractCallTx);
        message.tokens = [];
        message.fees = [];
        message.invalidationScope = new Uint8Array();
        message.payload = new Uint8Array();
        if (object.invalidationNonce !== undefined &&
            object.invalidationNonce !== null) {
            message.invalidationNonce = long_1["default"].fromString(object.invalidationNonce);
        }
        else {
            message.invalidationNonce = long_1["default"].UZERO;
        }
        if (object.invalidationScope !== undefined &&
            object.invalidationScope !== null) {
            message.invalidationScope = bytesFromBase64(object.invalidationScope);
        }
        if (object.address !== undefined && object.address !== null) {
            message.address = String(object.address);
        }
        else {
            message.address = "";
        }
        if (object.payload !== undefined && object.payload !== null) {
            message.payload = bytesFromBase64(object.payload);
        }
        if (object.timeout !== undefined && object.timeout !== null) {
            message.timeout = long_1["default"].fromString(object.timeout);
        }
        else {
            message.timeout = long_1["default"].UZERO;
        }
        if (object.tokens !== undefined && object.tokens !== null) {
            for (var _i = 0, _a = object.tokens; _i < _a.length; _i++) {
                var e = _a[_i];
                message.tokens.push(exports.ERC20Token.fromJSON(e));
            }
        }
        if (object.fees !== undefined && object.fees !== null) {
            for (var _b = 0, _c = object.fees; _b < _c.length; _b++) {
                var e = _c[_b];
                message.fees.push(exports.ERC20Token.fromJSON(e));
            }
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = long_1["default"].fromString(object.height);
        }
        else {
            message.height = long_1["default"].UZERO;
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.invalidationNonce !== undefined &&
            (obj.invalidationNonce = (message.invalidationNonce || long_1["default"].UZERO).toString());
        message.invalidationScope !== undefined &&
            (obj.invalidationScope = base64FromBytes(message.invalidationScope !== undefined
                ? message.invalidationScope
                : new Uint8Array()));
        message.address !== undefined && (obj.address = message.address);
        message.payload !== undefined &&
            (obj.payload = base64FromBytes(message.payload !== undefined ? message.payload : new Uint8Array()));
        message.timeout !== undefined &&
            (obj.timeout = (message.timeout || long_1["default"].UZERO).toString());
        if (message.tokens) {
            obj.tokens = message.tokens.map(function (e) {
                return e ? exports.ERC20Token.toJSON(e) : undefined;
            });
        }
        else {
            obj.tokens = [];
        }
        if (message.fees) {
            obj.fees = message.fees.map(function (e) {
                return e ? exports.ERC20Token.toJSON(e) : undefined;
            });
        }
        else {
            obj.fees = [];
        }
        message.height !== undefined &&
            (obj.height = (message.height || long_1["default"].UZERO).toString());
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseContractCallTx);
        message.tokens = [];
        message.fees = [];
        if (object.invalidationNonce !== undefined &&
            object.invalidationNonce !== null) {
            message.invalidationNonce = object.invalidationNonce;
        }
        else {
            message.invalidationNonce = long_1["default"].UZERO;
        }
        if (object.invalidationScope !== undefined &&
            object.invalidationScope !== null) {
            message.invalidationScope = object.invalidationScope;
        }
        else {
            message.invalidationScope = new Uint8Array();
        }
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = "";
        }
        if (object.payload !== undefined && object.payload !== null) {
            message.payload = object.payload;
        }
        else {
            message.payload = new Uint8Array();
        }
        if (object.timeout !== undefined && object.timeout !== null) {
            message.timeout = object.timeout;
        }
        else {
            message.timeout = long_1["default"].UZERO;
        }
        if (object.tokens !== undefined && object.tokens !== null) {
            for (var _i = 0, _a = object.tokens; _i < _a.length; _i++) {
                var e = _a[_i];
                message.tokens.push(exports.ERC20Token.fromPartial(e));
            }
        }
        if (object.fees !== undefined && object.fees !== null) {
            for (var _b = 0, _c = object.fees; _b < _c.length; _b++) {
                var e = _c[_b];
                message.fees.push(exports.ERC20Token.fromPartial(e));
            }
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = object.height;
        }
        else {
            message.height = long_1["default"].UZERO;
        }
        return message;
    }
};
var baseERC20Token = { contract: "", amount: "" };
exports.ERC20Token = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        if (message.contract !== "") {
            writer.uint32(10).string(message.contract);
        }
        if (message.amount !== "") {
            writer.uint32(18).string(message.amount);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseERC20Token);
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.contract = reader.string();
                    break;
                case 2:
                    message.amount = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseERC20Token);
        if (object.contract !== undefined && object.contract !== null) {
            message.contract = String(object.contract);
        }
        else {
            message.contract = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = String(object.amount);
        }
        else {
            message.amount = "";
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        message.contract !== undefined && (obj.contract = message.contract);
        message.amount !== undefined && (obj.amount = message.amount);
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseERC20Token);
        if (object.contract !== undefined && object.contract !== null) {
            message.contract = object.contract;
        }
        else {
            message.contract = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = object.amount;
        }
        else {
            message.amount = "";
        }
        return message;
    }
};
var baseIDSet = { ids: long_1["default"].UZERO };
exports.IDSet = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = minimal_1["default"].Writer.create(); }
        writer.uint32(10).fork();
        for (var _i = 0, _a = message.ids; _i < _a.length; _i++) {
            var v = _a[_i];
            writer.uint64(v);
        }
        writer.ldelim();
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof minimal_1["default"].Reader ? input : new minimal_1["default"].Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = __assign({}, baseIDSet);
        message.ids = [];
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    if ((tag & 7) === 2) {
                        var end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2) {
                            message.ids.push(reader.uint64());
                        }
                    }
                    else {
                        message.ids.push(reader.uint64());
                    }
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (object) {
        var message = __assign({}, baseIDSet);
        message.ids = [];
        if (object.ids !== undefined && object.ids !== null) {
            for (var _i = 0, _a = object.ids; _i < _a.length; _i++) {
                var e = _a[_i];
                message.ids.push(long_1["default"].fromString(e));
            }
        }
        return message;
    },
    toJSON: function (message) {
        var obj = {};
        if (message.ids) {
            obj.ids = message.ids.map(function (e) { return (e || long_1["default"].UZERO).toString(); });
        }
        else {
            obj.ids = [];
        }
        return obj;
    },
    fromPartial: function (object) {
        var message = __assign({}, baseIDSet);
        message.ids = [];
        if (object.ids !== undefined && object.ids !== null) {
            for (var _i = 0, _a = object.ids; _i < _a.length; _i++) {
                var e = _a[_i];
                message.ids.push(e);
            }
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
