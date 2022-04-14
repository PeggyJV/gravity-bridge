/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";
import { Any } from "../../google/protobuf/any";
import { EthereumSigner } from "../../gravity/v1/gravity";

export const protobufPackage = "gravity.v1";

/**
 * MsgSendToEthereum submits a SendToEthereum attempt to bridge an asset over to
 * Ethereum. The SendToEthereum will be stored and then included in a batch and
 * then submitted to Ethereum.
 */
export interface MsgSendToEthereum {
  sender: string;
  ethereumRecipient: string;
  amount?: Coin;
  bridgeFee?: Coin;
}

/**
 * MsgSendToEthereumResponse returns the SendToEthereum transaction ID which
 * will be included in the batch tx.
 */
export interface MsgSendToEthereumResponse {
  id: Long;
}

/**
 * MsgCancelSendToEthereum allows the sender to cancel its own outgoing
 * SendToEthereum tx and recieve a refund of the tokens and bridge fees. This tx
 * will only succeed if the SendToEthereum tx hasn't been batched to be
 * processed and relayed to Ethereum.
 */
export interface MsgCancelSendToEthereum {
  id: Long;
  sender: string;
}

export interface MsgCancelSendToEthereumResponse {}

/**
 * MsgRequestBatchTx requests a batch of transactions with a given coin
 * denomination to send across the bridge to Ethereum.
 */
export interface MsgRequestBatchTx {
  denom: string;
  signer: string;
}

export interface MsgRequestBatchTxResponse {}

/**
 * MsgSubmitEthereumTxConfirmation submits an ethereum signature for a given
 * validator
 */
export interface MsgSubmitEthereumTxConfirmation {
  /** TODO: can we make this take an array? */
  confirmation?: Any;
  signer: string;
}

/**
 * ContractCallTxConfirmation is a signature on behalf of a validator for a
 * ContractCallTx.
 */
export interface ContractCallTxConfirmation {
  invalidationScope: Uint8Array;
  invalidationNonce: Long;
  ethereumSigner: string;
  signature: Uint8Array;
}

/** BatchTxConfirmation is a signature on behalf of a validator for a BatchTx. */
export interface BatchTxConfirmation {
  tokenContract: string;
  batchNonce: Long;
  ethereumSigner: string;
  signature: Uint8Array;
}

/**
 * SignerSetTxConfirmation is a signature on behalf of a validator for a
 * SignerSetTx
 */
export interface SignerSetTxConfirmation {
  signerSetNonce: Long;
  ethereumSigner: string;
  signature: Uint8Array;
}

export interface MsgSubmitEthereumTxConfirmationResponse {}

/** MsgSubmitEthereumEvent */
export interface MsgSubmitEthereumEvent {
  event?: Any;
  signer: string;
}

export interface MsgSubmitEthereumEventResponse {}

/**
 * MsgDelegateKey allows validators to delegate their voting responsibilities
 * to a given orchestrator address. This key is then used as an optional
 * authentication method for attesting events from Ethereum.
 */
export interface MsgDelegateKeys {
  validatorAddress: string;
  orchestratorAddress: string;
  ethereumAddress: string;
  ethSignature: Uint8Array;
}

export interface MsgDelegateKeysResponse {}

/**
 * DelegateKeysSignMsg defines the message structure an operator is expected to
 * sign when submitting a MsgDelegateKeys message. The resulting signature
 * should populate the eth_signature field.
 */
export interface DelegateKeysSignMsg {
  validatorAddress: string;
  nonce: Long;
}

/**
 * SendToCosmosEvent is submitted when the SendToCosmosEvent is emitted by they
 * gravity contract. ERC20 representation coins are minted to the cosmosreceiver
 * address.
 */
export interface SendToCosmosEvent {
  eventNonce: Long;
  tokenContract: string;
  amount: string;
  ethereumSender: string;
  cosmosReceiver: string;
  ethereumHeight: Long;
}

/**
 * BatchExecutedEvent claims that a batch of BatchTxExecutedal operations on the
 * bridge contract was executed successfully on ETH
 */
export interface BatchExecutedEvent {
  tokenContract: string;
  eventNonce: Long;
  ethereumHeight: Long;
  batchNonce: Long;
}

/**
 * NOTE: bytes.HexBytes is supposed to "help" with json encoding/decoding
 * investigate?
 */
export interface ContractCallExecutedEvent {
  eventNonce: Long;
  invalidationScope: Uint8Array;
  invalidationNonce: Long;
  ethereumHeight: Long;
}

/**
 * ERC20DeployedEvent is submitted when an ERC20 contract
 * for a Cosmos SDK coin has been deployed on Ethereum.
 */
export interface ERC20DeployedEvent {
  eventNonce: Long;
  cosmosDenom: string;
  tokenContract: string;
  erc20Name: string;
  erc20Symbol: string;
  erc20Decimals: Long;
  ethereumHeight: Long;
}

/**
 * This informs the Cosmos module that a validator
 * set has been updated.
 */
export interface SignerSetTxExecutedEvent {
  eventNonce: Long;
  signerSetTxNonce: Long;
  ethereumHeight: Long;
  members: EthereumSigner[];
}

const baseMsgSendToEthereum: object = { sender: "", ethereumRecipient: "" };

export const MsgSendToEthereum = {
  encode(
    message: MsgSendToEthereum,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.sender !== "") {
      writer.uint32(10).string(message.sender);
    }
    if (message.ethereumRecipient !== "") {
      writer.uint32(18).string(message.ethereumRecipient);
    }
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(26).fork()).ldelim();
    }
    if (message.bridgeFee !== undefined) {
      Coin.encode(message.bridgeFee, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendToEthereum {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgSendToEthereum } as MsgSendToEthereum;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sender = reader.string();
          break;
        case 2:
          message.ethereumRecipient = reader.string();
          break;
        case 3:
          message.amount = Coin.decode(reader, reader.uint32());
          break;
        case 4:
          message.bridgeFee = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSendToEthereum {
    const message = { ...baseMsgSendToEthereum } as MsgSendToEthereum;
    message.sender =
      object.sender !== undefined && object.sender !== null
        ? String(object.sender)
        : "";
    message.ethereumRecipient =
      object.ethereumRecipient !== undefined &&
      object.ethereumRecipient !== null
        ? String(object.ethereumRecipient)
        : "";
    message.amount =
      object.amount !== undefined && object.amount !== null
        ? Coin.fromJSON(object.amount)
        : undefined;
    message.bridgeFee =
      object.bridgeFee !== undefined && object.bridgeFee !== null
        ? Coin.fromJSON(object.bridgeFee)
        : undefined;
    return message;
  },

  toJSON(message: MsgSendToEthereum): unknown {
    const obj: any = {};
    message.sender !== undefined && (obj.sender = message.sender);
    message.ethereumRecipient !== undefined &&
      (obj.ethereumRecipient = message.ethereumRecipient);
    message.amount !== undefined &&
      (obj.amount = message.amount ? Coin.toJSON(message.amount) : undefined);
    message.bridgeFee !== undefined &&
      (obj.bridgeFee = message.bridgeFee
        ? Coin.toJSON(message.bridgeFee)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgSendToEthereum>): MsgSendToEthereum {
    const message = { ...baseMsgSendToEthereum } as MsgSendToEthereum;
    message.sender = object.sender ?? "";
    message.ethereumRecipient = object.ethereumRecipient ?? "";
    message.amount =
      object.amount !== undefined && object.amount !== null
        ? Coin.fromPartial(object.amount)
        : undefined;
    message.bridgeFee =
      object.bridgeFee !== undefined && object.bridgeFee !== null
        ? Coin.fromPartial(object.bridgeFee)
        : undefined;
    return message;
  },
};

const baseMsgSendToEthereumResponse: object = { id: Long.UZERO };

export const MsgSendToEthereumResponse = {
  encode(
    message: MsgSendToEthereumResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.id.isZero()) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): MsgSendToEthereumResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgSendToEthereumResponse,
    } as MsgSendToEthereumResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSendToEthereumResponse {
    const message = {
      ...baseMsgSendToEthereumResponse,
    } as MsgSendToEthereumResponse;
    message.id =
      object.id !== undefined && object.id !== null
        ? Long.fromString(object.id)
        : Long.UZERO;
    return message;
  },

  toJSON(message: MsgSendToEthereumResponse): unknown {
    const obj: any = {};
    message.id !== undefined &&
      (obj.id = (message.id || Long.UZERO).toString());
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgSendToEthereumResponse>
  ): MsgSendToEthereumResponse {
    const message = {
      ...baseMsgSendToEthereumResponse,
    } as MsgSendToEthereumResponse;
    message.id =
      object.id !== undefined && object.id !== null
        ? Long.fromValue(object.id)
        : Long.UZERO;
    return message;
  },
};

const baseMsgCancelSendToEthereum: object = { id: Long.UZERO, sender: "" };

export const MsgCancelSendToEthereum = {
  encode(
    message: MsgCancelSendToEthereum,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.id.isZero()) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.sender !== "") {
      writer.uint32(18).string(message.sender);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): MsgCancelSendToEthereum {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgCancelSendToEthereum,
    } as MsgCancelSendToEthereum;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint64() as Long;
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

  fromJSON(object: any): MsgCancelSendToEthereum {
    const message = {
      ...baseMsgCancelSendToEthereum,
    } as MsgCancelSendToEthereum;
    message.id =
      object.id !== undefined && object.id !== null
        ? Long.fromString(object.id)
        : Long.UZERO;
    message.sender =
      object.sender !== undefined && object.sender !== null
        ? String(object.sender)
        : "";
    return message;
  },

  toJSON(message: MsgCancelSendToEthereum): unknown {
    const obj: any = {};
    message.id !== undefined &&
      (obj.id = (message.id || Long.UZERO).toString());
    message.sender !== undefined && (obj.sender = message.sender);
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgCancelSendToEthereum>
  ): MsgCancelSendToEthereum {
    const message = {
      ...baseMsgCancelSendToEthereum,
    } as MsgCancelSendToEthereum;
    message.id =
      object.id !== undefined && object.id !== null
        ? Long.fromValue(object.id)
        : Long.UZERO;
    message.sender = object.sender ?? "";
    return message;
  },
};

const baseMsgCancelSendToEthereumResponse: object = {};

export const MsgCancelSendToEthereumResponse = {
  encode(
    _: MsgCancelSendToEthereumResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): MsgCancelSendToEthereumResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgCancelSendToEthereumResponse,
    } as MsgCancelSendToEthereumResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgCancelSendToEthereumResponse {
    const message = {
      ...baseMsgCancelSendToEthereumResponse,
    } as MsgCancelSendToEthereumResponse;
    return message;
  },

  toJSON(_: MsgCancelSendToEthereumResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgCancelSendToEthereumResponse>
  ): MsgCancelSendToEthereumResponse {
    const message = {
      ...baseMsgCancelSendToEthereumResponse,
    } as MsgCancelSendToEthereumResponse;
    return message;
  },
};

const baseMsgRequestBatchTx: object = { denom: "", signer: "" };

export const MsgRequestBatchTx = {
  encode(
    message: MsgRequestBatchTx,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    if (message.signer !== "") {
      writer.uint32(18).string(message.signer);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRequestBatchTx {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgRequestBatchTx } as MsgRequestBatchTx;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): MsgRequestBatchTx {
    const message = { ...baseMsgRequestBatchTx } as MsgRequestBatchTx;
    message.denom =
      object.denom !== undefined && object.denom !== null
        ? String(object.denom)
        : "";
    message.signer =
      object.signer !== undefined && object.signer !== null
        ? String(object.signer)
        : "";
    return message;
  },

  toJSON(message: MsgRequestBatchTx): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    message.signer !== undefined && (obj.signer = message.signer);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgRequestBatchTx>): MsgRequestBatchTx {
    const message = { ...baseMsgRequestBatchTx } as MsgRequestBatchTx;
    message.denom = object.denom ?? "";
    message.signer = object.signer ?? "";
    return message;
  },
};

const baseMsgRequestBatchTxResponse: object = {};

export const MsgRequestBatchTxResponse = {
  encode(
    _: MsgRequestBatchTxResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): MsgRequestBatchTxResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgRequestBatchTxResponse,
    } as MsgRequestBatchTxResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgRequestBatchTxResponse {
    const message = {
      ...baseMsgRequestBatchTxResponse,
    } as MsgRequestBatchTxResponse;
    return message;
  },

  toJSON(_: MsgRequestBatchTxResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgRequestBatchTxResponse>
  ): MsgRequestBatchTxResponse {
    const message = {
      ...baseMsgRequestBatchTxResponse,
    } as MsgRequestBatchTxResponse;
    return message;
  },
};

const baseMsgSubmitEthereumTxConfirmation: object = { signer: "" };

export const MsgSubmitEthereumTxConfirmation = {
  encode(
    message: MsgSubmitEthereumTxConfirmation,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.confirmation !== undefined) {
      Any.encode(message.confirmation, writer.uint32(10).fork()).ldelim();
    }
    if (message.signer !== "") {
      writer.uint32(18).string(message.signer);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): MsgSubmitEthereumTxConfirmation {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgSubmitEthereumTxConfirmation,
    } as MsgSubmitEthereumTxConfirmation;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.confirmation = Any.decode(reader, reader.uint32());
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

  fromJSON(object: any): MsgSubmitEthereumTxConfirmation {
    const message = {
      ...baseMsgSubmitEthereumTxConfirmation,
    } as MsgSubmitEthereumTxConfirmation;
    message.confirmation =
      object.confirmation !== undefined && object.confirmation !== null
        ? Any.fromJSON(object.confirmation)
        : undefined;
    message.signer =
      object.signer !== undefined && object.signer !== null
        ? String(object.signer)
        : "";
    return message;
  },

  toJSON(message: MsgSubmitEthereumTxConfirmation): unknown {
    const obj: any = {};
    message.confirmation !== undefined &&
      (obj.confirmation = message.confirmation
        ? Any.toJSON(message.confirmation)
        : undefined);
    message.signer !== undefined && (obj.signer = message.signer);
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgSubmitEthereumTxConfirmation>
  ): MsgSubmitEthereumTxConfirmation {
    const message = {
      ...baseMsgSubmitEthereumTxConfirmation,
    } as MsgSubmitEthereumTxConfirmation;
    message.confirmation =
      object.confirmation !== undefined && object.confirmation !== null
        ? Any.fromPartial(object.confirmation)
        : undefined;
    message.signer = object.signer ?? "";
    return message;
  },
};

const baseContractCallTxConfirmation: object = {
  invalidationNonce: Long.UZERO,
  ethereumSigner: "",
};

export const ContractCallTxConfirmation = {
  encode(
    message: ContractCallTxConfirmation,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
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

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ContractCallTxConfirmation {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseContractCallTxConfirmation,
    } as ContractCallTxConfirmation;
    message.invalidationScope = new Uint8Array();
    message.signature = new Uint8Array();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.invalidationScope = reader.bytes();
          break;
        case 2:
          message.invalidationNonce = reader.uint64() as Long;
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

  fromJSON(object: any): ContractCallTxConfirmation {
    const message = {
      ...baseContractCallTxConfirmation,
    } as ContractCallTxConfirmation;
    message.invalidationScope =
      object.invalidationScope !== undefined &&
      object.invalidationScope !== null
        ? bytesFromBase64(object.invalidationScope)
        : new Uint8Array();
    message.invalidationNonce =
      object.invalidationNonce !== undefined &&
      object.invalidationNonce !== null
        ? Long.fromString(object.invalidationNonce)
        : Long.UZERO;
    message.ethereumSigner =
      object.ethereumSigner !== undefined && object.ethereumSigner !== null
        ? String(object.ethereumSigner)
        : "";
    message.signature =
      object.signature !== undefined && object.signature !== null
        ? bytesFromBase64(object.signature)
        : new Uint8Array();
    return message;
  },

  toJSON(message: ContractCallTxConfirmation): unknown {
    const obj: any = {};
    message.invalidationScope !== undefined &&
      (obj.invalidationScope = base64FromBytes(
        message.invalidationScope !== undefined
          ? message.invalidationScope
          : new Uint8Array()
      ));
    message.invalidationNonce !== undefined &&
      (obj.invalidationNonce = (
        message.invalidationNonce || Long.UZERO
      ).toString());
    message.ethereumSigner !== undefined &&
      (obj.ethereumSigner = message.ethereumSigner);
    message.signature !== undefined &&
      (obj.signature = base64FromBytes(
        message.signature !== undefined ? message.signature : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<ContractCallTxConfirmation>
  ): ContractCallTxConfirmation {
    const message = {
      ...baseContractCallTxConfirmation,
    } as ContractCallTxConfirmation;
    message.invalidationScope = object.invalidationScope ?? new Uint8Array();
    message.invalidationNonce =
      object.invalidationNonce !== undefined &&
      object.invalidationNonce !== null
        ? Long.fromValue(object.invalidationNonce)
        : Long.UZERO;
    message.ethereumSigner = object.ethereumSigner ?? "";
    message.signature = object.signature ?? new Uint8Array();
    return message;
  },
};

const baseBatchTxConfirmation: object = {
  tokenContract: "",
  batchNonce: Long.UZERO,
  ethereumSigner: "",
};

export const BatchTxConfirmation = {
  encode(
    message: BatchTxConfirmation,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): BatchTxConfirmation {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBatchTxConfirmation } as BatchTxConfirmation;
    message.signature = new Uint8Array();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tokenContract = reader.string();
          break;
        case 2:
          message.batchNonce = reader.uint64() as Long;
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

  fromJSON(object: any): BatchTxConfirmation {
    const message = { ...baseBatchTxConfirmation } as BatchTxConfirmation;
    message.tokenContract =
      object.tokenContract !== undefined && object.tokenContract !== null
        ? String(object.tokenContract)
        : "";
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromString(object.batchNonce)
        : Long.UZERO;
    message.ethereumSigner =
      object.ethereumSigner !== undefined && object.ethereumSigner !== null
        ? String(object.ethereumSigner)
        : "";
    message.signature =
      object.signature !== undefined && object.signature !== null
        ? bytesFromBase64(object.signature)
        : new Uint8Array();
    return message;
  },

  toJSON(message: BatchTxConfirmation): unknown {
    const obj: any = {};
    message.tokenContract !== undefined &&
      (obj.tokenContract = message.tokenContract);
    message.batchNonce !== undefined &&
      (obj.batchNonce = (message.batchNonce || Long.UZERO).toString());
    message.ethereumSigner !== undefined &&
      (obj.ethereumSigner = message.ethereumSigner);
    message.signature !== undefined &&
      (obj.signature = base64FromBytes(
        message.signature !== undefined ? message.signature : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(object: DeepPartial<BatchTxConfirmation>): BatchTxConfirmation {
    const message = { ...baseBatchTxConfirmation } as BatchTxConfirmation;
    message.tokenContract = object.tokenContract ?? "";
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromValue(object.batchNonce)
        : Long.UZERO;
    message.ethereumSigner = object.ethereumSigner ?? "";
    message.signature = object.signature ?? new Uint8Array();
    return message;
  },
};

const baseSignerSetTxConfirmation: object = {
  signerSetNonce: Long.UZERO,
  ethereumSigner: "",
};

export const SignerSetTxConfirmation = {
  encode(
    message: SignerSetTxConfirmation,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
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

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): SignerSetTxConfirmation {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseSignerSetTxConfirmation,
    } as SignerSetTxConfirmation;
    message.signature = new Uint8Array();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signerSetNonce = reader.uint64() as Long;
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

  fromJSON(object: any): SignerSetTxConfirmation {
    const message = {
      ...baseSignerSetTxConfirmation,
    } as SignerSetTxConfirmation;
    message.signerSetNonce =
      object.signerSetNonce !== undefined && object.signerSetNonce !== null
        ? Long.fromString(object.signerSetNonce)
        : Long.UZERO;
    message.ethereumSigner =
      object.ethereumSigner !== undefined && object.ethereumSigner !== null
        ? String(object.ethereumSigner)
        : "";
    message.signature =
      object.signature !== undefined && object.signature !== null
        ? bytesFromBase64(object.signature)
        : new Uint8Array();
    return message;
  },

  toJSON(message: SignerSetTxConfirmation): unknown {
    const obj: any = {};
    message.signerSetNonce !== undefined &&
      (obj.signerSetNonce = (message.signerSetNonce || Long.UZERO).toString());
    message.ethereumSigner !== undefined &&
      (obj.ethereumSigner = message.ethereumSigner);
    message.signature !== undefined &&
      (obj.signature = base64FromBytes(
        message.signature !== undefined ? message.signature : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<SignerSetTxConfirmation>
  ): SignerSetTxConfirmation {
    const message = {
      ...baseSignerSetTxConfirmation,
    } as SignerSetTxConfirmation;
    message.signerSetNonce =
      object.signerSetNonce !== undefined && object.signerSetNonce !== null
        ? Long.fromValue(object.signerSetNonce)
        : Long.UZERO;
    message.ethereumSigner = object.ethereumSigner ?? "";
    message.signature = object.signature ?? new Uint8Array();
    return message;
  },
};

const baseMsgSubmitEthereumTxConfirmationResponse: object = {};

export const MsgSubmitEthereumTxConfirmationResponse = {
  encode(
    _: MsgSubmitEthereumTxConfirmationResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): MsgSubmitEthereumTxConfirmationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgSubmitEthereumTxConfirmationResponse,
    } as MsgSubmitEthereumTxConfirmationResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgSubmitEthereumTxConfirmationResponse {
    const message = {
      ...baseMsgSubmitEthereumTxConfirmationResponse,
    } as MsgSubmitEthereumTxConfirmationResponse;
    return message;
  },

  toJSON(_: MsgSubmitEthereumTxConfirmationResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgSubmitEthereumTxConfirmationResponse>
  ): MsgSubmitEthereumTxConfirmationResponse {
    const message = {
      ...baseMsgSubmitEthereumTxConfirmationResponse,
    } as MsgSubmitEthereumTxConfirmationResponse;
    return message;
  },
};

const baseMsgSubmitEthereumEvent: object = { signer: "" };

export const MsgSubmitEthereumEvent = {
  encode(
    message: MsgSubmitEthereumEvent,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.event !== undefined) {
      Any.encode(message.event, writer.uint32(10).fork()).ldelim();
    }
    if (message.signer !== "") {
      writer.uint32(18).string(message.signer);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): MsgSubmitEthereumEvent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgSubmitEthereumEvent } as MsgSubmitEthereumEvent;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.event = Any.decode(reader, reader.uint32());
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

  fromJSON(object: any): MsgSubmitEthereumEvent {
    const message = { ...baseMsgSubmitEthereumEvent } as MsgSubmitEthereumEvent;
    message.event =
      object.event !== undefined && object.event !== null
        ? Any.fromJSON(object.event)
        : undefined;
    message.signer =
      object.signer !== undefined && object.signer !== null
        ? String(object.signer)
        : "";
    return message;
  },

  toJSON(message: MsgSubmitEthereumEvent): unknown {
    const obj: any = {};
    message.event !== undefined &&
      (obj.event = message.event ? Any.toJSON(message.event) : undefined);
    message.signer !== undefined && (obj.signer = message.signer);
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgSubmitEthereumEvent>
  ): MsgSubmitEthereumEvent {
    const message = { ...baseMsgSubmitEthereumEvent } as MsgSubmitEthereumEvent;
    message.event =
      object.event !== undefined && object.event !== null
        ? Any.fromPartial(object.event)
        : undefined;
    message.signer = object.signer ?? "";
    return message;
  },
};

const baseMsgSubmitEthereumEventResponse: object = {};

export const MsgSubmitEthereumEventResponse = {
  encode(
    _: MsgSubmitEthereumEventResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): MsgSubmitEthereumEventResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgSubmitEthereumEventResponse,
    } as MsgSubmitEthereumEventResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgSubmitEthereumEventResponse {
    const message = {
      ...baseMsgSubmitEthereumEventResponse,
    } as MsgSubmitEthereumEventResponse;
    return message;
  },

  toJSON(_: MsgSubmitEthereumEventResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgSubmitEthereumEventResponse>
  ): MsgSubmitEthereumEventResponse {
    const message = {
      ...baseMsgSubmitEthereumEventResponse,
    } as MsgSubmitEthereumEventResponse;
    return message;
  },
};

const baseMsgDelegateKeys: object = {
  validatorAddress: "",
  orchestratorAddress: "",
  ethereumAddress: "",
};

export const MsgDelegateKeys = {
  encode(
    message: MsgDelegateKeys,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.validatorAddress !== "") {
      writer.uint32(10).string(message.validatorAddress);
    }
    if (message.orchestratorAddress !== "") {
      writer.uint32(18).string(message.orchestratorAddress);
    }
    if (message.ethereumAddress !== "") {
      writer.uint32(26).string(message.ethereumAddress);
    }
    if (message.ethSignature.length !== 0) {
      writer.uint32(34).bytes(message.ethSignature);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDelegateKeys {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgDelegateKeys } as MsgDelegateKeys;
    message.ethSignature = new Uint8Array();
    while (reader.pos < end) {
      const tag = reader.uint32();
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
        case 4:
          message.ethSignature = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgDelegateKeys {
    const message = { ...baseMsgDelegateKeys } as MsgDelegateKeys;
    message.validatorAddress =
      object.validatorAddress !== undefined && object.validatorAddress !== null
        ? String(object.validatorAddress)
        : "";
    message.orchestratorAddress =
      object.orchestratorAddress !== undefined &&
      object.orchestratorAddress !== null
        ? String(object.orchestratorAddress)
        : "";
    message.ethereumAddress =
      object.ethereumAddress !== undefined && object.ethereumAddress !== null
        ? String(object.ethereumAddress)
        : "";
    message.ethSignature =
      object.ethSignature !== undefined && object.ethSignature !== null
        ? bytesFromBase64(object.ethSignature)
        : new Uint8Array();
    return message;
  },

  toJSON(message: MsgDelegateKeys): unknown {
    const obj: any = {};
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    message.orchestratorAddress !== undefined &&
      (obj.orchestratorAddress = message.orchestratorAddress);
    message.ethereumAddress !== undefined &&
      (obj.ethereumAddress = message.ethereumAddress);
    message.ethSignature !== undefined &&
      (obj.ethSignature = base64FromBytes(
        message.ethSignature !== undefined
          ? message.ethSignature
          : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(object: DeepPartial<MsgDelegateKeys>): MsgDelegateKeys {
    const message = { ...baseMsgDelegateKeys } as MsgDelegateKeys;
    message.validatorAddress = object.validatorAddress ?? "";
    message.orchestratorAddress = object.orchestratorAddress ?? "";
    message.ethereumAddress = object.ethereumAddress ?? "";
    message.ethSignature = object.ethSignature ?? new Uint8Array();
    return message;
  },
};

const baseMsgDelegateKeysResponse: object = {};

export const MsgDelegateKeysResponse = {
  encode(
    _: MsgDelegateKeysResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): MsgDelegateKeysResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgDelegateKeysResponse,
    } as MsgDelegateKeysResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgDelegateKeysResponse {
    const message = {
      ...baseMsgDelegateKeysResponse,
    } as MsgDelegateKeysResponse;
    return message;
  },

  toJSON(_: MsgDelegateKeysResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgDelegateKeysResponse>
  ): MsgDelegateKeysResponse {
    const message = {
      ...baseMsgDelegateKeysResponse,
    } as MsgDelegateKeysResponse;
    return message;
  },
};

const baseDelegateKeysSignMsg: object = {
  validatorAddress: "",
  nonce: Long.UZERO,
};

export const DelegateKeysSignMsg = {
  encode(
    message: DelegateKeysSignMsg,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.validatorAddress !== "") {
      writer.uint32(10).string(message.validatorAddress);
    }
    if (!message.nonce.isZero()) {
      writer.uint32(16).uint64(message.nonce);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DelegateKeysSignMsg {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDelegateKeysSignMsg } as DelegateKeysSignMsg;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.validatorAddress = reader.string();
          break;
        case 2:
          message.nonce = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DelegateKeysSignMsg {
    const message = { ...baseDelegateKeysSignMsg } as DelegateKeysSignMsg;
    message.validatorAddress =
      object.validatorAddress !== undefined && object.validatorAddress !== null
        ? String(object.validatorAddress)
        : "";
    message.nonce =
      object.nonce !== undefined && object.nonce !== null
        ? Long.fromString(object.nonce)
        : Long.UZERO;
    return message;
  },

  toJSON(message: DelegateKeysSignMsg): unknown {
    const obj: any = {};
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    message.nonce !== undefined &&
      (obj.nonce = (message.nonce || Long.UZERO).toString());
    return obj;
  },

  fromPartial(object: DeepPartial<DelegateKeysSignMsg>): DelegateKeysSignMsg {
    const message = { ...baseDelegateKeysSignMsg } as DelegateKeysSignMsg;
    message.validatorAddress = object.validatorAddress ?? "";
    message.nonce =
      object.nonce !== undefined && object.nonce !== null
        ? Long.fromValue(object.nonce)
        : Long.UZERO;
    return message;
  },
};

const baseSendToCosmosEvent: object = {
  eventNonce: Long.UZERO,
  tokenContract: "",
  amount: "",
  ethereumSender: "",
  cosmosReceiver: "",
  ethereumHeight: Long.UZERO,
};

export const SendToCosmosEvent = {
  encode(
    message: SendToCosmosEvent,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): SendToCosmosEvent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseSendToCosmosEvent } as SendToCosmosEvent;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.eventNonce = reader.uint64() as Long;
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
          message.ethereumHeight = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SendToCosmosEvent {
    const message = { ...baseSendToCosmosEvent } as SendToCosmosEvent;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromString(object.eventNonce)
        : Long.UZERO;
    message.tokenContract =
      object.tokenContract !== undefined && object.tokenContract !== null
        ? String(object.tokenContract)
        : "";
    message.amount =
      object.amount !== undefined && object.amount !== null
        ? String(object.amount)
        : "";
    message.ethereumSender =
      object.ethereumSender !== undefined && object.ethereumSender !== null
        ? String(object.ethereumSender)
        : "";
    message.cosmosReceiver =
      object.cosmosReceiver !== undefined && object.cosmosReceiver !== null
        ? String(object.cosmosReceiver)
        : "";
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromString(object.ethereumHeight)
        : Long.UZERO;
    return message;
  },

  toJSON(message: SendToCosmosEvent): unknown {
    const obj: any = {};
    message.eventNonce !== undefined &&
      (obj.eventNonce = (message.eventNonce || Long.UZERO).toString());
    message.tokenContract !== undefined &&
      (obj.tokenContract = message.tokenContract);
    message.amount !== undefined && (obj.amount = message.amount);
    message.ethereumSender !== undefined &&
      (obj.ethereumSender = message.ethereumSender);
    message.cosmosReceiver !== undefined &&
      (obj.cosmosReceiver = message.cosmosReceiver);
    message.ethereumHeight !== undefined &&
      (obj.ethereumHeight = (message.ethereumHeight || Long.UZERO).toString());
    return obj;
  },

  fromPartial(object: DeepPartial<SendToCosmosEvent>): SendToCosmosEvent {
    const message = { ...baseSendToCosmosEvent } as SendToCosmosEvent;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromValue(object.eventNonce)
        : Long.UZERO;
    message.tokenContract = object.tokenContract ?? "";
    message.amount = object.amount ?? "";
    message.ethereumSender = object.ethereumSender ?? "";
    message.cosmosReceiver = object.cosmosReceiver ?? "";
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromValue(object.ethereumHeight)
        : Long.UZERO;
    return message;
  },
};

const baseBatchExecutedEvent: object = {
  tokenContract: "",
  eventNonce: Long.UZERO,
  ethereumHeight: Long.UZERO,
  batchNonce: Long.UZERO,
};

export const BatchExecutedEvent = {
  encode(
    message: BatchExecutedEvent,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): BatchExecutedEvent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBatchExecutedEvent } as BatchExecutedEvent;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tokenContract = reader.string();
          break;
        case 2:
          message.eventNonce = reader.uint64() as Long;
          break;
        case 3:
          message.ethereumHeight = reader.uint64() as Long;
          break;
        case 4:
          message.batchNonce = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BatchExecutedEvent {
    const message = { ...baseBatchExecutedEvent } as BatchExecutedEvent;
    message.tokenContract =
      object.tokenContract !== undefined && object.tokenContract !== null
        ? String(object.tokenContract)
        : "";
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromString(object.eventNonce)
        : Long.UZERO;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromString(object.ethereumHeight)
        : Long.UZERO;
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromString(object.batchNonce)
        : Long.UZERO;
    return message;
  },

  toJSON(message: BatchExecutedEvent): unknown {
    const obj: any = {};
    message.tokenContract !== undefined &&
      (obj.tokenContract = message.tokenContract);
    message.eventNonce !== undefined &&
      (obj.eventNonce = (message.eventNonce || Long.UZERO).toString());
    message.ethereumHeight !== undefined &&
      (obj.ethereumHeight = (message.ethereumHeight || Long.UZERO).toString());
    message.batchNonce !== undefined &&
      (obj.batchNonce = (message.batchNonce || Long.UZERO).toString());
    return obj;
  },

  fromPartial(object: DeepPartial<BatchExecutedEvent>): BatchExecutedEvent {
    const message = { ...baseBatchExecutedEvent } as BatchExecutedEvent;
    message.tokenContract = object.tokenContract ?? "";
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromValue(object.eventNonce)
        : Long.UZERO;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromValue(object.ethereumHeight)
        : Long.UZERO;
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromValue(object.batchNonce)
        : Long.UZERO;
    return message;
  },
};

const baseContractCallExecutedEvent: object = {
  eventNonce: Long.UZERO,
  invalidationNonce: Long.UZERO,
  ethereumHeight: Long.UZERO,
};

export const ContractCallExecutedEvent = {
  encode(
    message: ContractCallExecutedEvent,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.eventNonce.isZero()) {
      writer.uint32(8).uint64(message.eventNonce);
    }
    if (message.invalidationScope.length !== 0) {
      writer.uint32(18).bytes(message.invalidationScope);
    }
    if (!message.invalidationNonce.isZero()) {
      writer.uint32(24).uint64(message.invalidationNonce);
    }
    if (!message.ethereumHeight.isZero()) {
      writer.uint32(32).uint64(message.ethereumHeight);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ContractCallExecutedEvent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseContractCallExecutedEvent,
    } as ContractCallExecutedEvent;
    message.invalidationScope = new Uint8Array();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.eventNonce = reader.uint64() as Long;
          break;
        case 2:
          message.invalidationScope = reader.bytes();
          break;
        case 3:
          message.invalidationNonce = reader.uint64() as Long;
          break;
        case 4:
          message.ethereumHeight = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractCallExecutedEvent {
    const message = {
      ...baseContractCallExecutedEvent,
    } as ContractCallExecutedEvent;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromString(object.eventNonce)
        : Long.UZERO;
    message.invalidationScope =
      object.invalidationScope !== undefined &&
      object.invalidationScope !== null
        ? bytesFromBase64(object.invalidationScope)
        : new Uint8Array();
    message.invalidationNonce =
      object.invalidationNonce !== undefined &&
      object.invalidationNonce !== null
        ? Long.fromString(object.invalidationNonce)
        : Long.UZERO;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromString(object.ethereumHeight)
        : Long.UZERO;
    return message;
  },

  toJSON(message: ContractCallExecutedEvent): unknown {
    const obj: any = {};
    message.eventNonce !== undefined &&
      (obj.eventNonce = (message.eventNonce || Long.UZERO).toString());
    message.invalidationScope !== undefined &&
      (obj.invalidationScope = base64FromBytes(
        message.invalidationScope !== undefined
          ? message.invalidationScope
          : new Uint8Array()
      ));
    message.invalidationNonce !== undefined &&
      (obj.invalidationNonce = (
        message.invalidationNonce || Long.UZERO
      ).toString());
    message.ethereumHeight !== undefined &&
      (obj.ethereumHeight = (message.ethereumHeight || Long.UZERO).toString());
    return obj;
  },

  fromPartial(
    object: DeepPartial<ContractCallExecutedEvent>
  ): ContractCallExecutedEvent {
    const message = {
      ...baseContractCallExecutedEvent,
    } as ContractCallExecutedEvent;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromValue(object.eventNonce)
        : Long.UZERO;
    message.invalidationScope = object.invalidationScope ?? new Uint8Array();
    message.invalidationNonce =
      object.invalidationNonce !== undefined &&
      object.invalidationNonce !== null
        ? Long.fromValue(object.invalidationNonce)
        : Long.UZERO;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromValue(object.ethereumHeight)
        : Long.UZERO;
    return message;
  },
};

const baseERC20DeployedEvent: object = {
  eventNonce: Long.UZERO,
  cosmosDenom: "",
  tokenContract: "",
  erc20Name: "",
  erc20Symbol: "",
  erc20Decimals: Long.UZERO,
  ethereumHeight: Long.UZERO,
};

export const ERC20DeployedEvent = {
  encode(
    message: ERC20DeployedEvent,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ERC20DeployedEvent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseERC20DeployedEvent } as ERC20DeployedEvent;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.eventNonce = reader.uint64() as Long;
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
          message.erc20Decimals = reader.uint64() as Long;
          break;
        case 7:
          message.ethereumHeight = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ERC20DeployedEvent {
    const message = { ...baseERC20DeployedEvent } as ERC20DeployedEvent;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromString(object.eventNonce)
        : Long.UZERO;
    message.cosmosDenom =
      object.cosmosDenom !== undefined && object.cosmosDenom !== null
        ? String(object.cosmosDenom)
        : "";
    message.tokenContract =
      object.tokenContract !== undefined && object.tokenContract !== null
        ? String(object.tokenContract)
        : "";
    message.erc20Name =
      object.erc20Name !== undefined && object.erc20Name !== null
        ? String(object.erc20Name)
        : "";
    message.erc20Symbol =
      object.erc20Symbol !== undefined && object.erc20Symbol !== null
        ? String(object.erc20Symbol)
        : "";
    message.erc20Decimals =
      object.erc20Decimals !== undefined && object.erc20Decimals !== null
        ? Long.fromString(object.erc20Decimals)
        : Long.UZERO;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromString(object.ethereumHeight)
        : Long.UZERO;
    return message;
  },

  toJSON(message: ERC20DeployedEvent): unknown {
    const obj: any = {};
    message.eventNonce !== undefined &&
      (obj.eventNonce = (message.eventNonce || Long.UZERO).toString());
    message.cosmosDenom !== undefined &&
      (obj.cosmosDenom = message.cosmosDenom);
    message.tokenContract !== undefined &&
      (obj.tokenContract = message.tokenContract);
    message.erc20Name !== undefined && (obj.erc20Name = message.erc20Name);
    message.erc20Symbol !== undefined &&
      (obj.erc20Symbol = message.erc20Symbol);
    message.erc20Decimals !== undefined &&
      (obj.erc20Decimals = (message.erc20Decimals || Long.UZERO).toString());
    message.ethereumHeight !== undefined &&
      (obj.ethereumHeight = (message.ethereumHeight || Long.UZERO).toString());
    return obj;
  },

  fromPartial(object: DeepPartial<ERC20DeployedEvent>): ERC20DeployedEvent {
    const message = { ...baseERC20DeployedEvent } as ERC20DeployedEvent;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromValue(object.eventNonce)
        : Long.UZERO;
    message.cosmosDenom = object.cosmosDenom ?? "";
    message.tokenContract = object.tokenContract ?? "";
    message.erc20Name = object.erc20Name ?? "";
    message.erc20Symbol = object.erc20Symbol ?? "";
    message.erc20Decimals =
      object.erc20Decimals !== undefined && object.erc20Decimals !== null
        ? Long.fromValue(object.erc20Decimals)
        : Long.UZERO;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromValue(object.ethereumHeight)
        : Long.UZERO;
    return message;
  },
};

const baseSignerSetTxExecutedEvent: object = {
  eventNonce: Long.UZERO,
  signerSetTxNonce: Long.UZERO,
  ethereumHeight: Long.UZERO,
};

export const SignerSetTxExecutedEvent = {
  encode(
    message: SignerSetTxExecutedEvent,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.eventNonce.isZero()) {
      writer.uint32(8).uint64(message.eventNonce);
    }
    if (!message.signerSetTxNonce.isZero()) {
      writer.uint32(16).uint64(message.signerSetTxNonce);
    }
    if (!message.ethereumHeight.isZero()) {
      writer.uint32(24).uint64(message.ethereumHeight);
    }
    for (const v of message.members) {
      EthereumSigner.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): SignerSetTxExecutedEvent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseSignerSetTxExecutedEvent,
    } as SignerSetTxExecutedEvent;
    message.members = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.eventNonce = reader.uint64() as Long;
          break;
        case 2:
          message.signerSetTxNonce = reader.uint64() as Long;
          break;
        case 3:
          message.ethereumHeight = reader.uint64() as Long;
          break;
        case 4:
          message.members.push(EthereumSigner.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignerSetTxExecutedEvent {
    const message = {
      ...baseSignerSetTxExecutedEvent,
    } as SignerSetTxExecutedEvent;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromString(object.eventNonce)
        : Long.UZERO;
    message.signerSetTxNonce =
      object.signerSetTxNonce !== undefined && object.signerSetTxNonce !== null
        ? Long.fromString(object.signerSetTxNonce)
        : Long.UZERO;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromString(object.ethereumHeight)
        : Long.UZERO;
    message.members = (object.members ?? []).map((e: any) =>
      EthereumSigner.fromJSON(e)
    );
    return message;
  },

  toJSON(message: SignerSetTxExecutedEvent): unknown {
    const obj: any = {};
    message.eventNonce !== undefined &&
      (obj.eventNonce = (message.eventNonce || Long.UZERO).toString());
    message.signerSetTxNonce !== undefined &&
      (obj.signerSetTxNonce = (
        message.signerSetTxNonce || Long.UZERO
      ).toString());
    message.ethereumHeight !== undefined &&
      (obj.ethereumHeight = (message.ethereumHeight || Long.UZERO).toString());
    if (message.members) {
      obj.members = message.members.map((e) =>
        e ? EthereumSigner.toJSON(e) : undefined
      );
    } else {
      obj.members = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<SignerSetTxExecutedEvent>
  ): SignerSetTxExecutedEvent {
    const message = {
      ...baseSignerSetTxExecutedEvent,
    } as SignerSetTxExecutedEvent;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromValue(object.eventNonce)
        : Long.UZERO;
    message.signerSetTxNonce =
      object.signerSetTxNonce !== undefined && object.signerSetTxNonce !== null
        ? Long.fromValue(object.signerSetTxNonce)
        : Long.UZERO;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromValue(object.ethereumHeight)
        : Long.UZERO;
    message.members = (object.members ?? []).map((e) =>
      EthereumSigner.fromPartial(e)
    );
    return message;
  },
};

/** Msg defines the state transitions possible within gravity */
export interface Msg {
  SendToEthereum(
    request: MsgSendToEthereum
  ): Promise<MsgSendToEthereumResponse>;
  CancelSendToEthereum(
    request: MsgCancelSendToEthereum
  ): Promise<MsgCancelSendToEthereumResponse>;
  RequestBatchTx(
    request: MsgRequestBatchTx
  ): Promise<MsgRequestBatchTxResponse>;
  SubmitEthereumTxConfirmation(
    request: MsgSubmitEthereumTxConfirmation
  ): Promise<MsgSubmitEthereumTxConfirmationResponse>;
  SubmitEthereumEvent(
    request: MsgSubmitEthereumEvent
  ): Promise<MsgSubmitEthereumEventResponse>;
  SetDelegateKeys(request: MsgDelegateKeys): Promise<MsgDelegateKeysResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.SendToEthereum = this.SendToEthereum.bind(this);
    this.CancelSendToEthereum = this.CancelSendToEthereum.bind(this);
    this.RequestBatchTx = this.RequestBatchTx.bind(this);
    this.SubmitEthereumTxConfirmation =
      this.SubmitEthereumTxConfirmation.bind(this);
    this.SubmitEthereumEvent = this.SubmitEthereumEvent.bind(this);
    this.SetDelegateKeys = this.SetDelegateKeys.bind(this);
  }
  SendToEthereum(
    request: MsgSendToEthereum
  ): Promise<MsgSendToEthereumResponse> {
    const data = MsgSendToEthereum.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Msg", "SendToEthereum", data);
    return promise.then((data) =>
      MsgSendToEthereumResponse.decode(new _m0.Reader(data))
    );
  }

  CancelSendToEthereum(
    request: MsgCancelSendToEthereum
  ): Promise<MsgCancelSendToEthereumResponse> {
    const data = MsgCancelSendToEthereum.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Msg",
      "CancelSendToEthereum",
      data
    );
    return promise.then((data) =>
      MsgCancelSendToEthereumResponse.decode(new _m0.Reader(data))
    );
  }

  RequestBatchTx(
    request: MsgRequestBatchTx
  ): Promise<MsgRequestBatchTxResponse> {
    const data = MsgRequestBatchTx.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Msg", "RequestBatchTx", data);
    return promise.then((data) =>
      MsgRequestBatchTxResponse.decode(new _m0.Reader(data))
    );
  }

  SubmitEthereumTxConfirmation(
    request: MsgSubmitEthereumTxConfirmation
  ): Promise<MsgSubmitEthereumTxConfirmationResponse> {
    const data = MsgSubmitEthereumTxConfirmation.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Msg",
      "SubmitEthereumTxConfirmation",
      data
    );
    return promise.then((data) =>
      MsgSubmitEthereumTxConfirmationResponse.decode(new _m0.Reader(data))
    );
  }

  SubmitEthereumEvent(
    request: MsgSubmitEthereumEvent
  ): Promise<MsgSubmitEthereumEventResponse> {
    const data = MsgSubmitEthereumEvent.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Msg",
      "SubmitEthereumEvent",
      data
    );
    return promise.then((data) =>
      MsgSubmitEthereumEventResponse.decode(new _m0.Reader(data))
    );
  }

  SetDelegateKeys(request: MsgDelegateKeys): Promise<MsgDelegateKeysResponse> {
    const data = MsgDelegateKeys.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Msg", "SetDelegateKeys", data);
    return promise.then((data) =>
      MsgDelegateKeysResponse.decode(new _m0.Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

const atob: (b64: string) => string =
  globalThis.atob ||
  ((b64) => globalThis.Buffer.from(b64, "base64").toString("binary"));
function bytesFromBase64(b64: string): Uint8Array {
  const bin = atob(b64);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i);
  }
  return arr;
}

const btoa: (bin: string) => string =
  globalThis.btoa ||
  ((bin) => globalThis.Buffer.from(bin, "binary").toString("base64"));
function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = [];
  for (const byte of arr) {
    bin.push(String.fromCharCode(byte));
  }
  return btoa(bin.join(""));
}

type Builtin =
  | Date
  | Function
  | Uint8Array
  | string
  | number
  | boolean
  | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Long
  ? string | number | Long
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}
