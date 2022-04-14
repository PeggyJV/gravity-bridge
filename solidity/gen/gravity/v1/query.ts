/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Params } from "../../gravity/v1/genesis";
import {
  SignerSetTx,
  BatchTx,
  ContractCallTx,
  SendToEthereum,
} from "../../gravity/v1/gravity";
import {
  PageRequest,
  PageResponse,
} from "../../cosmos/base/query/v1beta1/pagination";
import {
  SignerSetTxConfirmation,
  ContractCallTxConfirmation,
  BatchTxConfirmation,
  MsgDelegateKeys,
} from "../../gravity/v1/msgs";
import { Coin } from "../../cosmos/base/v1beta1/coin";

export const protobufPackage = "gravity.v1";

/** rpc Params */
export interface ParamsRequest {}

export interface ParamsResponse {
  params?: Params;
}

/** rpc SignerSetTx */
export interface SignerSetTxRequest {
  signerSetNonce: Long;
}

export interface LatestSignerSetTxRequest {}

export interface SignerSetTxResponse {
  signerSet?: SignerSetTx;
}

/** rpc BatchTx */
export interface BatchTxRequest {
  tokenContract: string;
  batchNonce: Long;
}

export interface BatchTxResponse {
  batch?: BatchTx;
}

/** rpc ContractCallTx */
export interface ContractCallTxRequest {
  invalidationScope: Uint8Array;
  invalidationNonce: Long;
}

export interface ContractCallTxResponse {
  logicCall?: ContractCallTx;
}

/** rpc SignerSetTxConfirmations */
export interface SignerSetTxConfirmationsRequest {
  signerSetNonce: Long;
}

export interface SignerSetTxConfirmationsResponse {
  signatures: SignerSetTxConfirmation[];
}

/** rpc SignerSetTxs */
export interface SignerSetTxsRequest {
  pagination?: PageRequest;
}

export interface SignerSetTxsResponse {
  signerSets: SignerSetTx[];
  pagination?: PageResponse;
}

/** rpc BatchTxs */
export interface BatchTxsRequest {
  pagination?: PageRequest;
}

export interface BatchTxsResponse {
  batches: BatchTx[];
  pagination?: PageResponse;
}

/** rpc ContractCallTxs */
export interface ContractCallTxsRequest {
  pagination?: PageRequest;
}

export interface ContractCallTxsResponse {
  calls: ContractCallTx[];
  pagination?: PageResponse;
}

/** rpc UnsignedSignerSetTxs */
export interface UnsignedSignerSetTxsRequest {
  /**
   * NOTE: this is an sdk.AccAddress and can represent either the
   * orchestrator address or the corresponding validator address
   */
  address: string;
}

export interface UnsignedSignerSetTxsResponse {
  signerSets: SignerSetTx[];
}

export interface UnsignedBatchTxsRequest {
  /**
   * NOTE: this is an sdk.AccAddress and can represent either the
   * orchestrator address or the corresponding validator address
   */
  address: string;
}

export interface UnsignedBatchTxsResponse {
  /** Note these are returned with the signature empty */
  batches: BatchTx[];
}

/** rpc UnsignedContractCallTxs */
export interface UnsignedContractCallTxsRequest {
  address: string;
}

export interface UnsignedContractCallTxsResponse {
  calls: ContractCallTx[];
}

export interface BatchTxFeesRequest {}

export interface BatchTxFeesResponse {
  fees: Coin[];
}

export interface ContractCallTxConfirmationsRequest {
  invalidationScope: Uint8Array;
  invalidationNonce: Long;
}

export interface ContractCallTxConfirmationsResponse {
  signatures: ContractCallTxConfirmation[];
}

export interface BatchTxConfirmationsRequest {
  batchNonce: Long;
  tokenContract: string;
}

export interface BatchTxConfirmationsResponse {
  signatures: BatchTxConfirmation[];
}

export interface LastSubmittedEthereumEventRequest {
  address: string;
}

export interface LastSubmittedEthereumEventResponse {
  eventNonce: Long;
}

export interface ERC20ToDenomRequest {
  erc20: string;
}

export interface ERC20ToDenomResponse {
  denom: string;
  cosmosOriginated: boolean;
}

export interface DenomToERC20ParamsRequest {
  denom: string;
}

export interface DenomToERC20ParamsResponse {
  baseDenom: string;
  erc20Name: string;
  erc20Symbol: string;
  erc20Decimals: Long;
}

export interface DenomToERC20Request {
  denom: string;
}

export interface DenomToERC20Response {
  erc20: string;
  cosmosOriginated: boolean;
}

export interface DelegateKeysByValidatorRequest {
  validatorAddress: string;
}

export interface DelegateKeysByValidatorResponse {
  ethAddress: string;
  orchestratorAddress: string;
}

export interface DelegateKeysByEthereumSignerRequest {
  ethereumSigner: string;
}

export interface DelegateKeysByEthereumSignerResponse {
  validatorAddress: string;
  orchestratorAddress: string;
}

export interface DelegateKeysByOrchestratorRequest {
  orchestratorAddress: string;
}

export interface DelegateKeysByOrchestratorResponse {
  validatorAddress: string;
  ethereumSigner: string;
}

export interface DelegateKeysRequest {}

export interface DelegateKeysResponse {
  delegateKeys: MsgDelegateKeys[];
}

/** NOTE: if there is no sender address, return all */
export interface BatchedSendToEthereumsRequest {
  /**
   * todo: figure out how to paginate given n Batches with m Send To Ethereums
   *  cosmos.base.query.v1beta1.PageRequest pagination = 2;
   */
  senderAddress: string;
}

export interface BatchedSendToEthereumsResponse {
  /** cosmos.base.query.v1beta1.PageResponse pagination = 2; */
  sendToEthereums: SendToEthereum[];
}

export interface UnbatchedSendToEthereumsRequest {
  senderAddress: string;
  pagination?: PageRequest;
}

export interface UnbatchedSendToEthereumsResponse {
  sendToEthereums: SendToEthereum[];
  pagination?: PageResponse;
}

const baseParamsRequest: object = {};

export const ParamsRequest = {
  encode(
    _: ParamsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseParamsRequest } as ParamsRequest;
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

  fromJSON(_: any): ParamsRequest {
    const message = { ...baseParamsRequest } as ParamsRequest;
    return message;
  },

  toJSON(_: ParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<ParamsRequest>): ParamsRequest {
    const message = { ...baseParamsRequest } as ParamsRequest;
    return message;
  },
};

const baseParamsResponse: object = {};

export const ParamsResponse = {
  encode(
    message: ParamsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseParamsResponse } as ParamsResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ParamsResponse {
    const message = { ...baseParamsResponse } as ParamsResponse;
    message.params =
      object.params !== undefined && object.params !== null
        ? Params.fromJSON(object.params)
        : undefined;
    return message;
  },

  toJSON(message: ParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined &&
      (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<ParamsResponse>): ParamsResponse {
    const message = { ...baseParamsResponse } as ParamsResponse;
    message.params =
      object.params !== undefined && object.params !== null
        ? Params.fromPartial(object.params)
        : undefined;
    return message;
  },
};

const baseSignerSetTxRequest: object = { signerSetNonce: Long.UZERO };

export const SignerSetTxRequest = {
  encode(
    message: SignerSetTxRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.signerSetNonce.isZero()) {
      writer.uint32(8).uint64(message.signerSetNonce);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SignerSetTxRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseSignerSetTxRequest } as SignerSetTxRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signerSetNonce = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignerSetTxRequest {
    const message = { ...baseSignerSetTxRequest } as SignerSetTxRequest;
    message.signerSetNonce =
      object.signerSetNonce !== undefined && object.signerSetNonce !== null
        ? Long.fromString(object.signerSetNonce)
        : Long.UZERO;
    return message;
  },

  toJSON(message: SignerSetTxRequest): unknown {
    const obj: any = {};
    message.signerSetNonce !== undefined &&
      (obj.signerSetNonce = (message.signerSetNonce || Long.UZERO).toString());
    return obj;
  },

  fromPartial(object: DeepPartial<SignerSetTxRequest>): SignerSetTxRequest {
    const message = { ...baseSignerSetTxRequest } as SignerSetTxRequest;
    message.signerSetNonce =
      object.signerSetNonce !== undefined && object.signerSetNonce !== null
        ? Long.fromValue(object.signerSetNonce)
        : Long.UZERO;
    return message;
  },
};

const baseLatestSignerSetTxRequest: object = {};

export const LatestSignerSetTxRequest = {
  encode(
    _: LatestSignerSetTxRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): LatestSignerSetTxRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseLatestSignerSetTxRequest,
    } as LatestSignerSetTxRequest;
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

  fromJSON(_: any): LatestSignerSetTxRequest {
    const message = {
      ...baseLatestSignerSetTxRequest,
    } as LatestSignerSetTxRequest;
    return message;
  },

  toJSON(_: LatestSignerSetTxRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<LatestSignerSetTxRequest>
  ): LatestSignerSetTxRequest {
    const message = {
      ...baseLatestSignerSetTxRequest,
    } as LatestSignerSetTxRequest;
    return message;
  },
};

const baseSignerSetTxResponse: object = {};

export const SignerSetTxResponse = {
  encode(
    message: SignerSetTxResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.signerSet !== undefined) {
      SignerSetTx.encode(message.signerSet, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SignerSetTxResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseSignerSetTxResponse } as SignerSetTxResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signerSet = SignerSetTx.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignerSetTxResponse {
    const message = { ...baseSignerSetTxResponse } as SignerSetTxResponse;
    message.signerSet =
      object.signerSet !== undefined && object.signerSet !== null
        ? SignerSetTx.fromJSON(object.signerSet)
        : undefined;
    return message;
  },

  toJSON(message: SignerSetTxResponse): unknown {
    const obj: any = {};
    message.signerSet !== undefined &&
      (obj.signerSet = message.signerSet
        ? SignerSetTx.toJSON(message.signerSet)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<SignerSetTxResponse>): SignerSetTxResponse {
    const message = { ...baseSignerSetTxResponse } as SignerSetTxResponse;
    message.signerSet =
      object.signerSet !== undefined && object.signerSet !== null
        ? SignerSetTx.fromPartial(object.signerSet)
        : undefined;
    return message;
  },
};

const baseBatchTxRequest: object = {
  tokenContract: "",
  batchNonce: Long.UZERO,
};

export const BatchTxRequest = {
  encode(
    message: BatchTxRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.tokenContract !== "") {
      writer.uint32(10).string(message.tokenContract);
    }
    if (!message.batchNonce.isZero()) {
      writer.uint32(16).uint64(message.batchNonce);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BatchTxRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBatchTxRequest } as BatchTxRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tokenContract = reader.string();
          break;
        case 2:
          message.batchNonce = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BatchTxRequest {
    const message = { ...baseBatchTxRequest } as BatchTxRequest;
    message.tokenContract =
      object.tokenContract !== undefined && object.tokenContract !== null
        ? String(object.tokenContract)
        : "";
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromString(object.batchNonce)
        : Long.UZERO;
    return message;
  },

  toJSON(message: BatchTxRequest): unknown {
    const obj: any = {};
    message.tokenContract !== undefined &&
      (obj.tokenContract = message.tokenContract);
    message.batchNonce !== undefined &&
      (obj.batchNonce = (message.batchNonce || Long.UZERO).toString());
    return obj;
  },

  fromPartial(object: DeepPartial<BatchTxRequest>): BatchTxRequest {
    const message = { ...baseBatchTxRequest } as BatchTxRequest;
    message.tokenContract = object.tokenContract ?? "";
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromValue(object.batchNonce)
        : Long.UZERO;
    return message;
  },
};

const baseBatchTxResponse: object = {};

export const BatchTxResponse = {
  encode(
    message: BatchTxResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.batch !== undefined) {
      BatchTx.encode(message.batch, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BatchTxResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBatchTxResponse } as BatchTxResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.batch = BatchTx.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BatchTxResponse {
    const message = { ...baseBatchTxResponse } as BatchTxResponse;
    message.batch =
      object.batch !== undefined && object.batch !== null
        ? BatchTx.fromJSON(object.batch)
        : undefined;
    return message;
  },

  toJSON(message: BatchTxResponse): unknown {
    const obj: any = {};
    message.batch !== undefined &&
      (obj.batch = message.batch ? BatchTx.toJSON(message.batch) : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<BatchTxResponse>): BatchTxResponse {
    const message = { ...baseBatchTxResponse } as BatchTxResponse;
    message.batch =
      object.batch !== undefined && object.batch !== null
        ? BatchTx.fromPartial(object.batch)
        : undefined;
    return message;
  },
};

const baseContractCallTxRequest: object = { invalidationNonce: Long.UZERO };

export const ContractCallTxRequest = {
  encode(
    message: ContractCallTxRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.invalidationScope.length !== 0) {
      writer.uint32(10).bytes(message.invalidationScope);
    }
    if (!message.invalidationNonce.isZero()) {
      writer.uint32(16).uint64(message.invalidationNonce);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ContractCallTxRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseContractCallTxRequest } as ContractCallTxRequest;
    message.invalidationScope = new Uint8Array();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.invalidationScope = reader.bytes();
          break;
        case 2:
          message.invalidationNonce = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractCallTxRequest {
    const message = { ...baseContractCallTxRequest } as ContractCallTxRequest;
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
    return message;
  },

  toJSON(message: ContractCallTxRequest): unknown {
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
    return obj;
  },

  fromPartial(
    object: DeepPartial<ContractCallTxRequest>
  ): ContractCallTxRequest {
    const message = { ...baseContractCallTxRequest } as ContractCallTxRequest;
    message.invalidationScope = object.invalidationScope ?? new Uint8Array();
    message.invalidationNonce =
      object.invalidationNonce !== undefined &&
      object.invalidationNonce !== null
        ? Long.fromValue(object.invalidationNonce)
        : Long.UZERO;
    return message;
  },
};

const baseContractCallTxResponse: object = {};

export const ContractCallTxResponse = {
  encode(
    message: ContractCallTxResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.logicCall !== undefined) {
      ContractCallTx.encode(
        message.logicCall,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ContractCallTxResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseContractCallTxResponse } as ContractCallTxResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.logicCall = ContractCallTx.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractCallTxResponse {
    const message = { ...baseContractCallTxResponse } as ContractCallTxResponse;
    message.logicCall =
      object.logicCall !== undefined && object.logicCall !== null
        ? ContractCallTx.fromJSON(object.logicCall)
        : undefined;
    return message;
  },

  toJSON(message: ContractCallTxResponse): unknown {
    const obj: any = {};
    message.logicCall !== undefined &&
      (obj.logicCall = message.logicCall
        ? ContractCallTx.toJSON(message.logicCall)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<ContractCallTxResponse>
  ): ContractCallTxResponse {
    const message = { ...baseContractCallTxResponse } as ContractCallTxResponse;
    message.logicCall =
      object.logicCall !== undefined && object.logicCall !== null
        ? ContractCallTx.fromPartial(object.logicCall)
        : undefined;
    return message;
  },
};

const baseSignerSetTxConfirmationsRequest: object = {
  signerSetNonce: Long.UZERO,
};

export const SignerSetTxConfirmationsRequest = {
  encode(
    message: SignerSetTxConfirmationsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.signerSetNonce.isZero()) {
      writer.uint32(8).uint64(message.signerSetNonce);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): SignerSetTxConfirmationsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseSignerSetTxConfirmationsRequest,
    } as SignerSetTxConfirmationsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signerSetNonce = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignerSetTxConfirmationsRequest {
    const message = {
      ...baseSignerSetTxConfirmationsRequest,
    } as SignerSetTxConfirmationsRequest;
    message.signerSetNonce =
      object.signerSetNonce !== undefined && object.signerSetNonce !== null
        ? Long.fromString(object.signerSetNonce)
        : Long.UZERO;
    return message;
  },

  toJSON(message: SignerSetTxConfirmationsRequest): unknown {
    const obj: any = {};
    message.signerSetNonce !== undefined &&
      (obj.signerSetNonce = (message.signerSetNonce || Long.UZERO).toString());
    return obj;
  },

  fromPartial(
    object: DeepPartial<SignerSetTxConfirmationsRequest>
  ): SignerSetTxConfirmationsRequest {
    const message = {
      ...baseSignerSetTxConfirmationsRequest,
    } as SignerSetTxConfirmationsRequest;
    message.signerSetNonce =
      object.signerSetNonce !== undefined && object.signerSetNonce !== null
        ? Long.fromValue(object.signerSetNonce)
        : Long.UZERO;
    return message;
  },
};

const baseSignerSetTxConfirmationsResponse: object = {};

export const SignerSetTxConfirmationsResponse = {
  encode(
    message: SignerSetTxConfirmationsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.signatures) {
      SignerSetTxConfirmation.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): SignerSetTxConfirmationsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseSignerSetTxConfirmationsResponse,
    } as SignerSetTxConfirmationsResponse;
    message.signatures = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signatures.push(
            SignerSetTxConfirmation.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignerSetTxConfirmationsResponse {
    const message = {
      ...baseSignerSetTxConfirmationsResponse,
    } as SignerSetTxConfirmationsResponse;
    message.signatures = (object.signatures ?? []).map((e: any) =>
      SignerSetTxConfirmation.fromJSON(e)
    );
    return message;
  },

  toJSON(message: SignerSetTxConfirmationsResponse): unknown {
    const obj: any = {};
    if (message.signatures) {
      obj.signatures = message.signatures.map((e) =>
        e ? SignerSetTxConfirmation.toJSON(e) : undefined
      );
    } else {
      obj.signatures = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<SignerSetTxConfirmationsResponse>
  ): SignerSetTxConfirmationsResponse {
    const message = {
      ...baseSignerSetTxConfirmationsResponse,
    } as SignerSetTxConfirmationsResponse;
    message.signatures = (object.signatures ?? []).map((e) =>
      SignerSetTxConfirmation.fromPartial(e)
    );
    return message;
  },
};

const baseSignerSetTxsRequest: object = {};

export const SignerSetTxsRequest = {
  encode(
    message: SignerSetTxsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SignerSetTxsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseSignerSetTxsRequest } as SignerSetTxsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignerSetTxsRequest {
    const message = { ...baseSignerSetTxsRequest } as SignerSetTxsRequest;
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageRequest.fromJSON(object.pagination)
        : undefined;
    return message;
  },

  toJSON(message: SignerSetTxsRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<SignerSetTxsRequest>): SignerSetTxsRequest {
    const message = { ...baseSignerSetTxsRequest } as SignerSetTxsRequest;
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageRequest.fromPartial(object.pagination)
        : undefined;
    return message;
  },
};

const baseSignerSetTxsResponse: object = {};

export const SignerSetTxsResponse = {
  encode(
    message: SignerSetTxsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.signerSets) {
      SignerSetTx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): SignerSetTxsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseSignerSetTxsResponse } as SignerSetTxsResponse;
    message.signerSets = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signerSets.push(SignerSetTx.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignerSetTxsResponse {
    const message = { ...baseSignerSetTxsResponse } as SignerSetTxsResponse;
    message.signerSets = (object.signerSets ?? []).map((e: any) =>
      SignerSetTx.fromJSON(e)
    );
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageResponse.fromJSON(object.pagination)
        : undefined;
    return message;
  },

  toJSON(message: SignerSetTxsResponse): unknown {
    const obj: any = {};
    if (message.signerSets) {
      obj.signerSets = message.signerSets.map((e) =>
        e ? SignerSetTx.toJSON(e) : undefined
      );
    } else {
      obj.signerSets = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<SignerSetTxsResponse>): SignerSetTxsResponse {
    const message = { ...baseSignerSetTxsResponse } as SignerSetTxsResponse;
    message.signerSets = (object.signerSets ?? []).map((e) =>
      SignerSetTx.fromPartial(e)
    );
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageResponse.fromPartial(object.pagination)
        : undefined;
    return message;
  },
};

const baseBatchTxsRequest: object = {};

export const BatchTxsRequest = {
  encode(
    message: BatchTxsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BatchTxsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBatchTxsRequest } as BatchTxsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BatchTxsRequest {
    const message = { ...baseBatchTxsRequest } as BatchTxsRequest;
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageRequest.fromJSON(object.pagination)
        : undefined;
    return message;
  },

  toJSON(message: BatchTxsRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<BatchTxsRequest>): BatchTxsRequest {
    const message = { ...baseBatchTxsRequest } as BatchTxsRequest;
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageRequest.fromPartial(object.pagination)
        : undefined;
    return message;
  },
};

const baseBatchTxsResponse: object = {};

export const BatchTxsResponse = {
  encode(
    message: BatchTxsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.batches) {
      BatchTx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BatchTxsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBatchTxsResponse } as BatchTxsResponse;
    message.batches = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.batches.push(BatchTx.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BatchTxsResponse {
    const message = { ...baseBatchTxsResponse } as BatchTxsResponse;
    message.batches = (object.batches ?? []).map((e: any) =>
      BatchTx.fromJSON(e)
    );
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageResponse.fromJSON(object.pagination)
        : undefined;
    return message;
  },

  toJSON(message: BatchTxsResponse): unknown {
    const obj: any = {};
    if (message.batches) {
      obj.batches = message.batches.map((e) =>
        e ? BatchTx.toJSON(e) : undefined
      );
    } else {
      obj.batches = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<BatchTxsResponse>): BatchTxsResponse {
    const message = { ...baseBatchTxsResponse } as BatchTxsResponse;
    message.batches = (object.batches ?? []).map((e) => BatchTx.fromPartial(e));
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageResponse.fromPartial(object.pagination)
        : undefined;
    return message;
  },
};

const baseContractCallTxsRequest: object = {};

export const ContractCallTxsRequest = {
  encode(
    message: ContractCallTxsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ContractCallTxsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseContractCallTxsRequest } as ContractCallTxsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractCallTxsRequest {
    const message = { ...baseContractCallTxsRequest } as ContractCallTxsRequest;
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageRequest.fromJSON(object.pagination)
        : undefined;
    return message;
  },

  toJSON(message: ContractCallTxsRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<ContractCallTxsRequest>
  ): ContractCallTxsRequest {
    const message = { ...baseContractCallTxsRequest } as ContractCallTxsRequest;
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageRequest.fromPartial(object.pagination)
        : undefined;
    return message;
  },
};

const baseContractCallTxsResponse: object = {};

export const ContractCallTxsResponse = {
  encode(
    message: ContractCallTxsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.calls) {
      ContractCallTx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ContractCallTxsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseContractCallTxsResponse,
    } as ContractCallTxsResponse;
    message.calls = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.calls.push(ContractCallTx.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractCallTxsResponse {
    const message = {
      ...baseContractCallTxsResponse,
    } as ContractCallTxsResponse;
    message.calls = (object.calls ?? []).map((e: any) =>
      ContractCallTx.fromJSON(e)
    );
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageResponse.fromJSON(object.pagination)
        : undefined;
    return message;
  },

  toJSON(message: ContractCallTxsResponse): unknown {
    const obj: any = {};
    if (message.calls) {
      obj.calls = message.calls.map((e) =>
        e ? ContractCallTx.toJSON(e) : undefined
      );
    } else {
      obj.calls = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<ContractCallTxsResponse>
  ): ContractCallTxsResponse {
    const message = {
      ...baseContractCallTxsResponse,
    } as ContractCallTxsResponse;
    message.calls = (object.calls ?? []).map((e) =>
      ContractCallTx.fromPartial(e)
    );
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageResponse.fromPartial(object.pagination)
        : undefined;
    return message;
  },
};

const baseUnsignedSignerSetTxsRequest: object = { address: "" };

export const UnsignedSignerSetTxsRequest = {
  encode(
    message: UnsignedSignerSetTxsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): UnsignedSignerSetTxsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseUnsignedSignerSetTxsRequest,
    } as UnsignedSignerSetTxsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): UnsignedSignerSetTxsRequest {
    const message = {
      ...baseUnsignedSignerSetTxsRequest,
    } as UnsignedSignerSetTxsRequest;
    message.address =
      object.address !== undefined && object.address !== null
        ? String(object.address)
        : "";
    return message;
  },

  toJSON(message: UnsignedSignerSetTxsRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnsignedSignerSetTxsRequest>
  ): UnsignedSignerSetTxsRequest {
    const message = {
      ...baseUnsignedSignerSetTxsRequest,
    } as UnsignedSignerSetTxsRequest;
    message.address = object.address ?? "";
    return message;
  },
};

const baseUnsignedSignerSetTxsResponse: object = {};

export const UnsignedSignerSetTxsResponse = {
  encode(
    message: UnsignedSignerSetTxsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.signerSets) {
      SignerSetTx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): UnsignedSignerSetTxsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseUnsignedSignerSetTxsResponse,
    } as UnsignedSignerSetTxsResponse;
    message.signerSets = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signerSets.push(SignerSetTx.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnsignedSignerSetTxsResponse {
    const message = {
      ...baseUnsignedSignerSetTxsResponse,
    } as UnsignedSignerSetTxsResponse;
    message.signerSets = (object.signerSets ?? []).map((e: any) =>
      SignerSetTx.fromJSON(e)
    );
    return message;
  },

  toJSON(message: UnsignedSignerSetTxsResponse): unknown {
    const obj: any = {};
    if (message.signerSets) {
      obj.signerSets = message.signerSets.map((e) =>
        e ? SignerSetTx.toJSON(e) : undefined
      );
    } else {
      obj.signerSets = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnsignedSignerSetTxsResponse>
  ): UnsignedSignerSetTxsResponse {
    const message = {
      ...baseUnsignedSignerSetTxsResponse,
    } as UnsignedSignerSetTxsResponse;
    message.signerSets = (object.signerSets ?? []).map((e) =>
      SignerSetTx.fromPartial(e)
    );
    return message;
  },
};

const baseUnsignedBatchTxsRequest: object = { address: "" };

export const UnsignedBatchTxsRequest = {
  encode(
    message: UnsignedBatchTxsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): UnsignedBatchTxsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseUnsignedBatchTxsRequest,
    } as UnsignedBatchTxsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): UnsignedBatchTxsRequest {
    const message = {
      ...baseUnsignedBatchTxsRequest,
    } as UnsignedBatchTxsRequest;
    message.address =
      object.address !== undefined && object.address !== null
        ? String(object.address)
        : "";
    return message;
  },

  toJSON(message: UnsignedBatchTxsRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnsignedBatchTxsRequest>
  ): UnsignedBatchTxsRequest {
    const message = {
      ...baseUnsignedBatchTxsRequest,
    } as UnsignedBatchTxsRequest;
    message.address = object.address ?? "";
    return message;
  },
};

const baseUnsignedBatchTxsResponse: object = {};

export const UnsignedBatchTxsResponse = {
  encode(
    message: UnsignedBatchTxsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.batches) {
      BatchTx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): UnsignedBatchTxsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseUnsignedBatchTxsResponse,
    } as UnsignedBatchTxsResponse;
    message.batches = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.batches.push(BatchTx.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnsignedBatchTxsResponse {
    const message = {
      ...baseUnsignedBatchTxsResponse,
    } as UnsignedBatchTxsResponse;
    message.batches = (object.batches ?? []).map((e: any) =>
      BatchTx.fromJSON(e)
    );
    return message;
  },

  toJSON(message: UnsignedBatchTxsResponse): unknown {
    const obj: any = {};
    if (message.batches) {
      obj.batches = message.batches.map((e) =>
        e ? BatchTx.toJSON(e) : undefined
      );
    } else {
      obj.batches = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnsignedBatchTxsResponse>
  ): UnsignedBatchTxsResponse {
    const message = {
      ...baseUnsignedBatchTxsResponse,
    } as UnsignedBatchTxsResponse;
    message.batches = (object.batches ?? []).map((e) => BatchTx.fromPartial(e));
    return message;
  },
};

const baseUnsignedContractCallTxsRequest: object = { address: "" };

export const UnsignedContractCallTxsRequest = {
  encode(
    message: UnsignedContractCallTxsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): UnsignedContractCallTxsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseUnsignedContractCallTxsRequest,
    } as UnsignedContractCallTxsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): UnsignedContractCallTxsRequest {
    const message = {
      ...baseUnsignedContractCallTxsRequest,
    } as UnsignedContractCallTxsRequest;
    message.address =
      object.address !== undefined && object.address !== null
        ? String(object.address)
        : "";
    return message;
  },

  toJSON(message: UnsignedContractCallTxsRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnsignedContractCallTxsRequest>
  ): UnsignedContractCallTxsRequest {
    const message = {
      ...baseUnsignedContractCallTxsRequest,
    } as UnsignedContractCallTxsRequest;
    message.address = object.address ?? "";
    return message;
  },
};

const baseUnsignedContractCallTxsResponse: object = {};

export const UnsignedContractCallTxsResponse = {
  encode(
    message: UnsignedContractCallTxsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.calls) {
      ContractCallTx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): UnsignedContractCallTxsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseUnsignedContractCallTxsResponse,
    } as UnsignedContractCallTxsResponse;
    message.calls = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.calls.push(ContractCallTx.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnsignedContractCallTxsResponse {
    const message = {
      ...baseUnsignedContractCallTxsResponse,
    } as UnsignedContractCallTxsResponse;
    message.calls = (object.calls ?? []).map((e: any) =>
      ContractCallTx.fromJSON(e)
    );
    return message;
  },

  toJSON(message: UnsignedContractCallTxsResponse): unknown {
    const obj: any = {};
    if (message.calls) {
      obj.calls = message.calls.map((e) =>
        e ? ContractCallTx.toJSON(e) : undefined
      );
    } else {
      obj.calls = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnsignedContractCallTxsResponse>
  ): UnsignedContractCallTxsResponse {
    const message = {
      ...baseUnsignedContractCallTxsResponse,
    } as UnsignedContractCallTxsResponse;
    message.calls = (object.calls ?? []).map((e) =>
      ContractCallTx.fromPartial(e)
    );
    return message;
  },
};

const baseBatchTxFeesRequest: object = {};

export const BatchTxFeesRequest = {
  encode(
    _: BatchTxFeesRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BatchTxFeesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBatchTxFeesRequest } as BatchTxFeesRequest;
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

  fromJSON(_: any): BatchTxFeesRequest {
    const message = { ...baseBatchTxFeesRequest } as BatchTxFeesRequest;
    return message;
  },

  toJSON(_: BatchTxFeesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<BatchTxFeesRequest>): BatchTxFeesRequest {
    const message = { ...baseBatchTxFeesRequest } as BatchTxFeesRequest;
    return message;
  },
};

const baseBatchTxFeesResponse: object = {};

export const BatchTxFeesResponse = {
  encode(
    message: BatchTxFeesResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.fees) {
      Coin.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BatchTxFeesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBatchTxFeesResponse } as BatchTxFeesResponse;
    message.fees = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fees.push(Coin.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BatchTxFeesResponse {
    const message = { ...baseBatchTxFeesResponse } as BatchTxFeesResponse;
    message.fees = (object.fees ?? []).map((e: any) => Coin.fromJSON(e));
    return message;
  },

  toJSON(message: BatchTxFeesResponse): unknown {
    const obj: any = {};
    if (message.fees) {
      obj.fees = message.fees.map((e) => (e ? Coin.toJSON(e) : undefined));
    } else {
      obj.fees = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<BatchTxFeesResponse>): BatchTxFeesResponse {
    const message = { ...baseBatchTxFeesResponse } as BatchTxFeesResponse;
    message.fees = (object.fees ?? []).map((e) => Coin.fromPartial(e));
    return message;
  },
};

const baseContractCallTxConfirmationsRequest: object = {
  invalidationNonce: Long.UZERO,
};

export const ContractCallTxConfirmationsRequest = {
  encode(
    message: ContractCallTxConfirmationsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.invalidationScope.length !== 0) {
      writer.uint32(10).bytes(message.invalidationScope);
    }
    if (!message.invalidationNonce.isZero()) {
      writer.uint32(16).uint64(message.invalidationNonce);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ContractCallTxConfirmationsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseContractCallTxConfirmationsRequest,
    } as ContractCallTxConfirmationsRequest;
    message.invalidationScope = new Uint8Array();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.invalidationScope = reader.bytes();
          break;
        case 2:
          message.invalidationNonce = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractCallTxConfirmationsRequest {
    const message = {
      ...baseContractCallTxConfirmationsRequest,
    } as ContractCallTxConfirmationsRequest;
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
    return message;
  },

  toJSON(message: ContractCallTxConfirmationsRequest): unknown {
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
    return obj;
  },

  fromPartial(
    object: DeepPartial<ContractCallTxConfirmationsRequest>
  ): ContractCallTxConfirmationsRequest {
    const message = {
      ...baseContractCallTxConfirmationsRequest,
    } as ContractCallTxConfirmationsRequest;
    message.invalidationScope = object.invalidationScope ?? new Uint8Array();
    message.invalidationNonce =
      object.invalidationNonce !== undefined &&
      object.invalidationNonce !== null
        ? Long.fromValue(object.invalidationNonce)
        : Long.UZERO;
    return message;
  },
};

const baseContractCallTxConfirmationsResponse: object = {};

export const ContractCallTxConfirmationsResponse = {
  encode(
    message: ContractCallTxConfirmationsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.signatures) {
      ContractCallTxConfirmation.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ContractCallTxConfirmationsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseContractCallTxConfirmationsResponse,
    } as ContractCallTxConfirmationsResponse;
    message.signatures = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signatures.push(
            ContractCallTxConfirmation.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractCallTxConfirmationsResponse {
    const message = {
      ...baseContractCallTxConfirmationsResponse,
    } as ContractCallTxConfirmationsResponse;
    message.signatures = (object.signatures ?? []).map((e: any) =>
      ContractCallTxConfirmation.fromJSON(e)
    );
    return message;
  },

  toJSON(message: ContractCallTxConfirmationsResponse): unknown {
    const obj: any = {};
    if (message.signatures) {
      obj.signatures = message.signatures.map((e) =>
        e ? ContractCallTxConfirmation.toJSON(e) : undefined
      );
    } else {
      obj.signatures = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<ContractCallTxConfirmationsResponse>
  ): ContractCallTxConfirmationsResponse {
    const message = {
      ...baseContractCallTxConfirmationsResponse,
    } as ContractCallTxConfirmationsResponse;
    message.signatures = (object.signatures ?? []).map((e) =>
      ContractCallTxConfirmation.fromPartial(e)
    );
    return message;
  },
};

const baseBatchTxConfirmationsRequest: object = {
  batchNonce: Long.UZERO,
  tokenContract: "",
};

export const BatchTxConfirmationsRequest = {
  encode(
    message: BatchTxConfirmationsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.batchNonce.isZero()) {
      writer.uint32(8).uint64(message.batchNonce);
    }
    if (message.tokenContract !== "") {
      writer.uint32(18).string(message.tokenContract);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): BatchTxConfirmationsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseBatchTxConfirmationsRequest,
    } as BatchTxConfirmationsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.batchNonce = reader.uint64() as Long;
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

  fromJSON(object: any): BatchTxConfirmationsRequest {
    const message = {
      ...baseBatchTxConfirmationsRequest,
    } as BatchTxConfirmationsRequest;
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromString(object.batchNonce)
        : Long.UZERO;
    message.tokenContract =
      object.tokenContract !== undefined && object.tokenContract !== null
        ? String(object.tokenContract)
        : "";
    return message;
  },

  toJSON(message: BatchTxConfirmationsRequest): unknown {
    const obj: any = {};
    message.batchNonce !== undefined &&
      (obj.batchNonce = (message.batchNonce || Long.UZERO).toString());
    message.tokenContract !== undefined &&
      (obj.tokenContract = message.tokenContract);
    return obj;
  },

  fromPartial(
    object: DeepPartial<BatchTxConfirmationsRequest>
  ): BatchTxConfirmationsRequest {
    const message = {
      ...baseBatchTxConfirmationsRequest,
    } as BatchTxConfirmationsRequest;
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromValue(object.batchNonce)
        : Long.UZERO;
    message.tokenContract = object.tokenContract ?? "";
    return message;
  },
};

const baseBatchTxConfirmationsResponse: object = {};

export const BatchTxConfirmationsResponse = {
  encode(
    message: BatchTxConfirmationsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.signatures) {
      BatchTxConfirmation.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): BatchTxConfirmationsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseBatchTxConfirmationsResponse,
    } as BatchTxConfirmationsResponse;
    message.signatures = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signatures.push(
            BatchTxConfirmation.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BatchTxConfirmationsResponse {
    const message = {
      ...baseBatchTxConfirmationsResponse,
    } as BatchTxConfirmationsResponse;
    message.signatures = (object.signatures ?? []).map((e: any) =>
      BatchTxConfirmation.fromJSON(e)
    );
    return message;
  },

  toJSON(message: BatchTxConfirmationsResponse): unknown {
    const obj: any = {};
    if (message.signatures) {
      obj.signatures = message.signatures.map((e) =>
        e ? BatchTxConfirmation.toJSON(e) : undefined
      );
    } else {
      obj.signatures = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<BatchTxConfirmationsResponse>
  ): BatchTxConfirmationsResponse {
    const message = {
      ...baseBatchTxConfirmationsResponse,
    } as BatchTxConfirmationsResponse;
    message.signatures = (object.signatures ?? []).map((e) =>
      BatchTxConfirmation.fromPartial(e)
    );
    return message;
  },
};

const baseLastSubmittedEthereumEventRequest: object = { address: "" };

export const LastSubmittedEthereumEventRequest = {
  encode(
    message: LastSubmittedEthereumEventRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): LastSubmittedEthereumEventRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseLastSubmittedEthereumEventRequest,
    } as LastSubmittedEthereumEventRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): LastSubmittedEthereumEventRequest {
    const message = {
      ...baseLastSubmittedEthereumEventRequest,
    } as LastSubmittedEthereumEventRequest;
    message.address =
      object.address !== undefined && object.address !== null
        ? String(object.address)
        : "";
    return message;
  },

  toJSON(message: LastSubmittedEthereumEventRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial(
    object: DeepPartial<LastSubmittedEthereumEventRequest>
  ): LastSubmittedEthereumEventRequest {
    const message = {
      ...baseLastSubmittedEthereumEventRequest,
    } as LastSubmittedEthereumEventRequest;
    message.address = object.address ?? "";
    return message;
  },
};

const baseLastSubmittedEthereumEventResponse: object = {
  eventNonce: Long.UZERO,
};

export const LastSubmittedEthereumEventResponse = {
  encode(
    message: LastSubmittedEthereumEventResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.eventNonce.isZero()) {
      writer.uint32(8).uint64(message.eventNonce);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): LastSubmittedEthereumEventResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseLastSubmittedEthereumEventResponse,
    } as LastSubmittedEthereumEventResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.eventNonce = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LastSubmittedEthereumEventResponse {
    const message = {
      ...baseLastSubmittedEthereumEventResponse,
    } as LastSubmittedEthereumEventResponse;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromString(object.eventNonce)
        : Long.UZERO;
    return message;
  },

  toJSON(message: LastSubmittedEthereumEventResponse): unknown {
    const obj: any = {};
    message.eventNonce !== undefined &&
      (obj.eventNonce = (message.eventNonce || Long.UZERO).toString());
    return obj;
  },

  fromPartial(
    object: DeepPartial<LastSubmittedEthereumEventResponse>
  ): LastSubmittedEthereumEventResponse {
    const message = {
      ...baseLastSubmittedEthereumEventResponse,
    } as LastSubmittedEthereumEventResponse;
    message.eventNonce =
      object.eventNonce !== undefined && object.eventNonce !== null
        ? Long.fromValue(object.eventNonce)
        : Long.UZERO;
    return message;
  },
};

const baseERC20ToDenomRequest: object = { erc20: "" };

export const ERC20ToDenomRequest = {
  encode(
    message: ERC20ToDenomRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.erc20 !== "") {
      writer.uint32(10).string(message.erc20);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ERC20ToDenomRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseERC20ToDenomRequest } as ERC20ToDenomRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): ERC20ToDenomRequest {
    const message = { ...baseERC20ToDenomRequest } as ERC20ToDenomRequest;
    message.erc20 =
      object.erc20 !== undefined && object.erc20 !== null
        ? String(object.erc20)
        : "";
    return message;
  },

  toJSON(message: ERC20ToDenomRequest): unknown {
    const obj: any = {};
    message.erc20 !== undefined && (obj.erc20 = message.erc20);
    return obj;
  },

  fromPartial(object: DeepPartial<ERC20ToDenomRequest>): ERC20ToDenomRequest {
    const message = { ...baseERC20ToDenomRequest } as ERC20ToDenomRequest;
    message.erc20 = object.erc20 ?? "";
    return message;
  },
};

const baseERC20ToDenomResponse: object = { denom: "", cosmosOriginated: false };

export const ERC20ToDenomResponse = {
  encode(
    message: ERC20ToDenomResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    if (message.cosmosOriginated === true) {
      writer.uint32(16).bool(message.cosmosOriginated);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ERC20ToDenomResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseERC20ToDenomResponse } as ERC20ToDenomResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): ERC20ToDenomResponse {
    const message = { ...baseERC20ToDenomResponse } as ERC20ToDenomResponse;
    message.denom =
      object.denom !== undefined && object.denom !== null
        ? String(object.denom)
        : "";
    message.cosmosOriginated =
      object.cosmosOriginated !== undefined && object.cosmosOriginated !== null
        ? Boolean(object.cosmosOriginated)
        : false;
    return message;
  },

  toJSON(message: ERC20ToDenomResponse): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    message.cosmosOriginated !== undefined &&
      (obj.cosmosOriginated = message.cosmosOriginated);
    return obj;
  },

  fromPartial(object: DeepPartial<ERC20ToDenomResponse>): ERC20ToDenomResponse {
    const message = { ...baseERC20ToDenomResponse } as ERC20ToDenomResponse;
    message.denom = object.denom ?? "";
    message.cosmosOriginated = object.cosmosOriginated ?? false;
    return message;
  },
};

const baseDenomToERC20ParamsRequest: object = { denom: "" };

export const DenomToERC20ParamsRequest = {
  encode(
    message: DenomToERC20ParamsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DenomToERC20ParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseDenomToERC20ParamsRequest,
    } as DenomToERC20ParamsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): DenomToERC20ParamsRequest {
    const message = {
      ...baseDenomToERC20ParamsRequest,
    } as DenomToERC20ParamsRequest;
    message.denom =
      object.denom !== undefined && object.denom !== null
        ? String(object.denom)
        : "";
    return message;
  },

  toJSON(message: DenomToERC20ParamsRequest): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    return obj;
  },

  fromPartial(
    object: DeepPartial<DenomToERC20ParamsRequest>
  ): DenomToERC20ParamsRequest {
    const message = {
      ...baseDenomToERC20ParamsRequest,
    } as DenomToERC20ParamsRequest;
    message.denom = object.denom ?? "";
    return message;
  },
};

const baseDenomToERC20ParamsResponse: object = {
  baseDenom: "",
  erc20Name: "",
  erc20Symbol: "",
  erc20Decimals: Long.UZERO,
};

export const DenomToERC20ParamsResponse = {
  encode(
    message: DenomToERC20ParamsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.baseDenom !== "") {
      writer.uint32(10).string(message.baseDenom);
    }
    if (message.erc20Name !== "") {
      writer.uint32(18).string(message.erc20Name);
    }
    if (message.erc20Symbol !== "") {
      writer.uint32(26).string(message.erc20Symbol);
    }
    if (!message.erc20Decimals.isZero()) {
      writer.uint32(32).uint64(message.erc20Decimals);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DenomToERC20ParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseDenomToERC20ParamsResponse,
    } as DenomToERC20ParamsResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseDenom = reader.string();
          break;
        case 2:
          message.erc20Name = reader.string();
          break;
        case 3:
          message.erc20Symbol = reader.string();
          break;
        case 4:
          message.erc20Decimals = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DenomToERC20ParamsResponse {
    const message = {
      ...baseDenomToERC20ParamsResponse,
    } as DenomToERC20ParamsResponse;
    message.baseDenom =
      object.baseDenom !== undefined && object.baseDenom !== null
        ? String(object.baseDenom)
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
    return message;
  },

  toJSON(message: DenomToERC20ParamsResponse): unknown {
    const obj: any = {};
    message.baseDenom !== undefined && (obj.baseDenom = message.baseDenom);
    message.erc20Name !== undefined && (obj.erc20Name = message.erc20Name);
    message.erc20Symbol !== undefined &&
      (obj.erc20Symbol = message.erc20Symbol);
    message.erc20Decimals !== undefined &&
      (obj.erc20Decimals = (message.erc20Decimals || Long.UZERO).toString());
    return obj;
  },

  fromPartial(
    object: DeepPartial<DenomToERC20ParamsResponse>
  ): DenomToERC20ParamsResponse {
    const message = {
      ...baseDenomToERC20ParamsResponse,
    } as DenomToERC20ParamsResponse;
    message.baseDenom = object.baseDenom ?? "";
    message.erc20Name = object.erc20Name ?? "";
    message.erc20Symbol = object.erc20Symbol ?? "";
    message.erc20Decimals =
      object.erc20Decimals !== undefined && object.erc20Decimals !== null
        ? Long.fromValue(object.erc20Decimals)
        : Long.UZERO;
    return message;
  },
};

const baseDenomToERC20Request: object = { denom: "" };

export const DenomToERC20Request = {
  encode(
    message: DenomToERC20Request,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DenomToERC20Request {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDenomToERC20Request } as DenomToERC20Request;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): DenomToERC20Request {
    const message = { ...baseDenomToERC20Request } as DenomToERC20Request;
    message.denom =
      object.denom !== undefined && object.denom !== null
        ? String(object.denom)
        : "";
    return message;
  },

  toJSON(message: DenomToERC20Request): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    return obj;
  },

  fromPartial(object: DeepPartial<DenomToERC20Request>): DenomToERC20Request {
    const message = { ...baseDenomToERC20Request } as DenomToERC20Request;
    message.denom = object.denom ?? "";
    return message;
  },
};

const baseDenomToERC20Response: object = { erc20: "", cosmosOriginated: false };

export const DenomToERC20Response = {
  encode(
    message: DenomToERC20Response,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.erc20 !== "") {
      writer.uint32(10).string(message.erc20);
    }
    if (message.cosmosOriginated === true) {
      writer.uint32(16).bool(message.cosmosOriginated);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DenomToERC20Response {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDenomToERC20Response } as DenomToERC20Response;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): DenomToERC20Response {
    const message = { ...baseDenomToERC20Response } as DenomToERC20Response;
    message.erc20 =
      object.erc20 !== undefined && object.erc20 !== null
        ? String(object.erc20)
        : "";
    message.cosmosOriginated =
      object.cosmosOriginated !== undefined && object.cosmosOriginated !== null
        ? Boolean(object.cosmosOriginated)
        : false;
    return message;
  },

  toJSON(message: DenomToERC20Response): unknown {
    const obj: any = {};
    message.erc20 !== undefined && (obj.erc20 = message.erc20);
    message.cosmosOriginated !== undefined &&
      (obj.cosmosOriginated = message.cosmosOriginated);
    return obj;
  },

  fromPartial(object: DeepPartial<DenomToERC20Response>): DenomToERC20Response {
    const message = { ...baseDenomToERC20Response } as DenomToERC20Response;
    message.erc20 = object.erc20 ?? "";
    message.cosmosOriginated = object.cosmosOriginated ?? false;
    return message;
  },
};

const baseDelegateKeysByValidatorRequest: object = { validatorAddress: "" };

export const DelegateKeysByValidatorRequest = {
  encode(
    message: DelegateKeysByValidatorRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.validatorAddress !== "") {
      writer.uint32(10).string(message.validatorAddress);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DelegateKeysByValidatorRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseDelegateKeysByValidatorRequest,
    } as DelegateKeysByValidatorRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): DelegateKeysByValidatorRequest {
    const message = {
      ...baseDelegateKeysByValidatorRequest,
    } as DelegateKeysByValidatorRequest;
    message.validatorAddress =
      object.validatorAddress !== undefined && object.validatorAddress !== null
        ? String(object.validatorAddress)
        : "";
    return message;
  },

  toJSON(message: DelegateKeysByValidatorRequest): unknown {
    const obj: any = {};
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    return obj;
  },

  fromPartial(
    object: DeepPartial<DelegateKeysByValidatorRequest>
  ): DelegateKeysByValidatorRequest {
    const message = {
      ...baseDelegateKeysByValidatorRequest,
    } as DelegateKeysByValidatorRequest;
    message.validatorAddress = object.validatorAddress ?? "";
    return message;
  },
};

const baseDelegateKeysByValidatorResponse: object = {
  ethAddress: "",
  orchestratorAddress: "",
};

export const DelegateKeysByValidatorResponse = {
  encode(
    message: DelegateKeysByValidatorResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.ethAddress !== "") {
      writer.uint32(10).string(message.ethAddress);
    }
    if (message.orchestratorAddress !== "") {
      writer.uint32(18).string(message.orchestratorAddress);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DelegateKeysByValidatorResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseDelegateKeysByValidatorResponse,
    } as DelegateKeysByValidatorResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): DelegateKeysByValidatorResponse {
    const message = {
      ...baseDelegateKeysByValidatorResponse,
    } as DelegateKeysByValidatorResponse;
    message.ethAddress =
      object.ethAddress !== undefined && object.ethAddress !== null
        ? String(object.ethAddress)
        : "";
    message.orchestratorAddress =
      object.orchestratorAddress !== undefined &&
      object.orchestratorAddress !== null
        ? String(object.orchestratorAddress)
        : "";
    return message;
  },

  toJSON(message: DelegateKeysByValidatorResponse): unknown {
    const obj: any = {};
    message.ethAddress !== undefined && (obj.ethAddress = message.ethAddress);
    message.orchestratorAddress !== undefined &&
      (obj.orchestratorAddress = message.orchestratorAddress);
    return obj;
  },

  fromPartial(
    object: DeepPartial<DelegateKeysByValidatorResponse>
  ): DelegateKeysByValidatorResponse {
    const message = {
      ...baseDelegateKeysByValidatorResponse,
    } as DelegateKeysByValidatorResponse;
    message.ethAddress = object.ethAddress ?? "";
    message.orchestratorAddress = object.orchestratorAddress ?? "";
    return message;
  },
};

const baseDelegateKeysByEthereumSignerRequest: object = { ethereumSigner: "" };

export const DelegateKeysByEthereumSignerRequest = {
  encode(
    message: DelegateKeysByEthereumSignerRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.ethereumSigner !== "") {
      writer.uint32(10).string(message.ethereumSigner);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DelegateKeysByEthereumSignerRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseDelegateKeysByEthereumSignerRequest,
    } as DelegateKeysByEthereumSignerRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): DelegateKeysByEthereumSignerRequest {
    const message = {
      ...baseDelegateKeysByEthereumSignerRequest,
    } as DelegateKeysByEthereumSignerRequest;
    message.ethereumSigner =
      object.ethereumSigner !== undefined && object.ethereumSigner !== null
        ? String(object.ethereumSigner)
        : "";
    return message;
  },

  toJSON(message: DelegateKeysByEthereumSignerRequest): unknown {
    const obj: any = {};
    message.ethereumSigner !== undefined &&
      (obj.ethereumSigner = message.ethereumSigner);
    return obj;
  },

  fromPartial(
    object: DeepPartial<DelegateKeysByEthereumSignerRequest>
  ): DelegateKeysByEthereumSignerRequest {
    const message = {
      ...baseDelegateKeysByEthereumSignerRequest,
    } as DelegateKeysByEthereumSignerRequest;
    message.ethereumSigner = object.ethereumSigner ?? "";
    return message;
  },
};

const baseDelegateKeysByEthereumSignerResponse: object = {
  validatorAddress: "",
  orchestratorAddress: "",
};

export const DelegateKeysByEthereumSignerResponse = {
  encode(
    message: DelegateKeysByEthereumSignerResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.validatorAddress !== "") {
      writer.uint32(10).string(message.validatorAddress);
    }
    if (message.orchestratorAddress !== "") {
      writer.uint32(18).string(message.orchestratorAddress);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DelegateKeysByEthereumSignerResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseDelegateKeysByEthereumSignerResponse,
    } as DelegateKeysByEthereumSignerResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): DelegateKeysByEthereumSignerResponse {
    const message = {
      ...baseDelegateKeysByEthereumSignerResponse,
    } as DelegateKeysByEthereumSignerResponse;
    message.validatorAddress =
      object.validatorAddress !== undefined && object.validatorAddress !== null
        ? String(object.validatorAddress)
        : "";
    message.orchestratorAddress =
      object.orchestratorAddress !== undefined &&
      object.orchestratorAddress !== null
        ? String(object.orchestratorAddress)
        : "";
    return message;
  },

  toJSON(message: DelegateKeysByEthereumSignerResponse): unknown {
    const obj: any = {};
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    message.orchestratorAddress !== undefined &&
      (obj.orchestratorAddress = message.orchestratorAddress);
    return obj;
  },

  fromPartial(
    object: DeepPartial<DelegateKeysByEthereumSignerResponse>
  ): DelegateKeysByEthereumSignerResponse {
    const message = {
      ...baseDelegateKeysByEthereumSignerResponse,
    } as DelegateKeysByEthereumSignerResponse;
    message.validatorAddress = object.validatorAddress ?? "";
    message.orchestratorAddress = object.orchestratorAddress ?? "";
    return message;
  },
};

const baseDelegateKeysByOrchestratorRequest: object = {
  orchestratorAddress: "",
};

export const DelegateKeysByOrchestratorRequest = {
  encode(
    message: DelegateKeysByOrchestratorRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.orchestratorAddress !== "") {
      writer.uint32(10).string(message.orchestratorAddress);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DelegateKeysByOrchestratorRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseDelegateKeysByOrchestratorRequest,
    } as DelegateKeysByOrchestratorRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): DelegateKeysByOrchestratorRequest {
    const message = {
      ...baseDelegateKeysByOrchestratorRequest,
    } as DelegateKeysByOrchestratorRequest;
    message.orchestratorAddress =
      object.orchestratorAddress !== undefined &&
      object.orchestratorAddress !== null
        ? String(object.orchestratorAddress)
        : "";
    return message;
  },

  toJSON(message: DelegateKeysByOrchestratorRequest): unknown {
    const obj: any = {};
    message.orchestratorAddress !== undefined &&
      (obj.orchestratorAddress = message.orchestratorAddress);
    return obj;
  },

  fromPartial(
    object: DeepPartial<DelegateKeysByOrchestratorRequest>
  ): DelegateKeysByOrchestratorRequest {
    const message = {
      ...baseDelegateKeysByOrchestratorRequest,
    } as DelegateKeysByOrchestratorRequest;
    message.orchestratorAddress = object.orchestratorAddress ?? "";
    return message;
  },
};

const baseDelegateKeysByOrchestratorResponse: object = {
  validatorAddress: "",
  ethereumSigner: "",
};

export const DelegateKeysByOrchestratorResponse = {
  encode(
    message: DelegateKeysByOrchestratorResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.validatorAddress !== "") {
      writer.uint32(10).string(message.validatorAddress);
    }
    if (message.ethereumSigner !== "") {
      writer.uint32(18).string(message.ethereumSigner);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DelegateKeysByOrchestratorResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseDelegateKeysByOrchestratorResponse,
    } as DelegateKeysByOrchestratorResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): DelegateKeysByOrchestratorResponse {
    const message = {
      ...baseDelegateKeysByOrchestratorResponse,
    } as DelegateKeysByOrchestratorResponse;
    message.validatorAddress =
      object.validatorAddress !== undefined && object.validatorAddress !== null
        ? String(object.validatorAddress)
        : "";
    message.ethereumSigner =
      object.ethereumSigner !== undefined && object.ethereumSigner !== null
        ? String(object.ethereumSigner)
        : "";
    return message;
  },

  toJSON(message: DelegateKeysByOrchestratorResponse): unknown {
    const obj: any = {};
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    message.ethereumSigner !== undefined &&
      (obj.ethereumSigner = message.ethereumSigner);
    return obj;
  },

  fromPartial(
    object: DeepPartial<DelegateKeysByOrchestratorResponse>
  ): DelegateKeysByOrchestratorResponse {
    const message = {
      ...baseDelegateKeysByOrchestratorResponse,
    } as DelegateKeysByOrchestratorResponse;
    message.validatorAddress = object.validatorAddress ?? "";
    message.ethereumSigner = object.ethereumSigner ?? "";
    return message;
  },
};

const baseDelegateKeysRequest: object = {};

export const DelegateKeysRequest = {
  encode(
    _: DelegateKeysRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DelegateKeysRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDelegateKeysRequest } as DelegateKeysRequest;
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

  fromJSON(_: any): DelegateKeysRequest {
    const message = { ...baseDelegateKeysRequest } as DelegateKeysRequest;
    return message;
  },

  toJSON(_: DelegateKeysRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<DelegateKeysRequest>): DelegateKeysRequest {
    const message = { ...baseDelegateKeysRequest } as DelegateKeysRequest;
    return message;
  },
};

const baseDelegateKeysResponse: object = {};

export const DelegateKeysResponse = {
  encode(
    message: DelegateKeysResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.delegateKeys) {
      MsgDelegateKeys.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DelegateKeysResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDelegateKeysResponse } as DelegateKeysResponse;
    message.delegateKeys = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegateKeys.push(
            MsgDelegateKeys.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DelegateKeysResponse {
    const message = { ...baseDelegateKeysResponse } as DelegateKeysResponse;
    message.delegateKeys = (object.delegateKeys ?? []).map((e: any) =>
      MsgDelegateKeys.fromJSON(e)
    );
    return message;
  },

  toJSON(message: DelegateKeysResponse): unknown {
    const obj: any = {};
    if (message.delegateKeys) {
      obj.delegateKeys = message.delegateKeys.map((e) =>
        e ? MsgDelegateKeys.toJSON(e) : undefined
      );
    } else {
      obj.delegateKeys = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<DelegateKeysResponse>): DelegateKeysResponse {
    const message = { ...baseDelegateKeysResponse } as DelegateKeysResponse;
    message.delegateKeys = (object.delegateKeys ?? []).map((e) =>
      MsgDelegateKeys.fromPartial(e)
    );
    return message;
  },
};

const baseBatchedSendToEthereumsRequest: object = { senderAddress: "" };

export const BatchedSendToEthereumsRequest = {
  encode(
    message: BatchedSendToEthereumsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.senderAddress !== "") {
      writer.uint32(10).string(message.senderAddress);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): BatchedSendToEthereumsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseBatchedSendToEthereumsRequest,
    } as BatchedSendToEthereumsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): BatchedSendToEthereumsRequest {
    const message = {
      ...baseBatchedSendToEthereumsRequest,
    } as BatchedSendToEthereumsRequest;
    message.senderAddress =
      object.senderAddress !== undefined && object.senderAddress !== null
        ? String(object.senderAddress)
        : "";
    return message;
  },

  toJSON(message: BatchedSendToEthereumsRequest): unknown {
    const obj: any = {};
    message.senderAddress !== undefined &&
      (obj.senderAddress = message.senderAddress);
    return obj;
  },

  fromPartial(
    object: DeepPartial<BatchedSendToEthereumsRequest>
  ): BatchedSendToEthereumsRequest {
    const message = {
      ...baseBatchedSendToEthereumsRequest,
    } as BatchedSendToEthereumsRequest;
    message.senderAddress = object.senderAddress ?? "";
    return message;
  },
};

const baseBatchedSendToEthereumsResponse: object = {};

export const BatchedSendToEthereumsResponse = {
  encode(
    message: BatchedSendToEthereumsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.sendToEthereums) {
      SendToEthereum.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): BatchedSendToEthereumsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseBatchedSendToEthereumsResponse,
    } as BatchedSendToEthereumsResponse;
    message.sendToEthereums = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sendToEthereums.push(
            SendToEthereum.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BatchedSendToEthereumsResponse {
    const message = {
      ...baseBatchedSendToEthereumsResponse,
    } as BatchedSendToEthereumsResponse;
    message.sendToEthereums = (object.sendToEthereums ?? []).map((e: any) =>
      SendToEthereum.fromJSON(e)
    );
    return message;
  },

  toJSON(message: BatchedSendToEthereumsResponse): unknown {
    const obj: any = {};
    if (message.sendToEthereums) {
      obj.sendToEthereums = message.sendToEthereums.map((e) =>
        e ? SendToEthereum.toJSON(e) : undefined
      );
    } else {
      obj.sendToEthereums = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<BatchedSendToEthereumsResponse>
  ): BatchedSendToEthereumsResponse {
    const message = {
      ...baseBatchedSendToEthereumsResponse,
    } as BatchedSendToEthereumsResponse;
    message.sendToEthereums = (object.sendToEthereums ?? []).map((e) =>
      SendToEthereum.fromPartial(e)
    );
    return message;
  },
};

const baseUnbatchedSendToEthereumsRequest: object = { senderAddress: "" };

export const UnbatchedSendToEthereumsRequest = {
  encode(
    message: UnbatchedSendToEthereumsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.senderAddress !== "") {
      writer.uint32(10).string(message.senderAddress);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): UnbatchedSendToEthereumsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseUnbatchedSendToEthereumsRequest,
    } as UnbatchedSendToEthereumsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.senderAddress = reader.string();
          break;
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnbatchedSendToEthereumsRequest {
    const message = {
      ...baseUnbatchedSendToEthereumsRequest,
    } as UnbatchedSendToEthereumsRequest;
    message.senderAddress =
      object.senderAddress !== undefined && object.senderAddress !== null
        ? String(object.senderAddress)
        : "";
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageRequest.fromJSON(object.pagination)
        : undefined;
    return message;
  },

  toJSON(message: UnbatchedSendToEthereumsRequest): unknown {
    const obj: any = {};
    message.senderAddress !== undefined &&
      (obj.senderAddress = message.senderAddress);
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnbatchedSendToEthereumsRequest>
  ): UnbatchedSendToEthereumsRequest {
    const message = {
      ...baseUnbatchedSendToEthereumsRequest,
    } as UnbatchedSendToEthereumsRequest;
    message.senderAddress = object.senderAddress ?? "";
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageRequest.fromPartial(object.pagination)
        : undefined;
    return message;
  },
};

const baseUnbatchedSendToEthereumsResponse: object = {};

export const UnbatchedSendToEthereumsResponse = {
  encode(
    message: UnbatchedSendToEthereumsResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.sendToEthereums) {
      SendToEthereum.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): UnbatchedSendToEthereumsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseUnbatchedSendToEthereumsResponse,
    } as UnbatchedSendToEthereumsResponse;
    message.sendToEthereums = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sendToEthereums.push(
            SendToEthereum.decode(reader, reader.uint32())
          );
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnbatchedSendToEthereumsResponse {
    const message = {
      ...baseUnbatchedSendToEthereumsResponse,
    } as UnbatchedSendToEthereumsResponse;
    message.sendToEthereums = (object.sendToEthereums ?? []).map((e: any) =>
      SendToEthereum.fromJSON(e)
    );
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageResponse.fromJSON(object.pagination)
        : undefined;
    return message;
  },

  toJSON(message: UnbatchedSendToEthereumsResponse): unknown {
    const obj: any = {};
    if (message.sendToEthereums) {
      obj.sendToEthereums = message.sendToEthereums.map((e) =>
        e ? SendToEthereum.toJSON(e) : undefined
      );
    } else {
      obj.sendToEthereums = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnbatchedSendToEthereumsResponse>
  ): UnbatchedSendToEthereumsResponse {
    const message = {
      ...baseUnbatchedSendToEthereumsResponse,
    } as UnbatchedSendToEthereumsResponse;
    message.sendToEthereums = (object.sendToEthereums ?? []).map((e) =>
      SendToEthereum.fromPartial(e)
    );
    message.pagination =
      object.pagination !== undefined && object.pagination !== null
        ? PageResponse.fromPartial(object.pagination)
        : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service */
export interface Query {
  /** Module parameters query */
  Params(request: ParamsRequest): Promise<ParamsResponse>;
  /** get info on individual outgoing data */
  SignerSetTx(request: SignerSetTxRequest): Promise<SignerSetTxResponse>;
  LatestSignerSetTx(
    request: LatestSignerSetTxRequest
  ): Promise<SignerSetTxResponse>;
  BatchTx(request: BatchTxRequest): Promise<BatchTxResponse>;
  ContractCallTx(
    request: ContractCallTxRequest
  ): Promise<ContractCallTxResponse>;
  /** get collections of outgoing traffic from the bridge */
  SignerSetTxs(request: SignerSetTxsRequest): Promise<SignerSetTxsResponse>;
  BatchTxs(request: BatchTxsRequest): Promise<BatchTxsResponse>;
  ContractCallTxs(
    request: ContractCallTxsRequest
  ): Promise<ContractCallTxsResponse>;
  /** TODO: can/should we group these into one endpoint? */
  SignerSetTxConfirmations(
    request: SignerSetTxConfirmationsRequest
  ): Promise<SignerSetTxConfirmationsResponse>;
  BatchTxConfirmations(
    request: BatchTxConfirmationsRequest
  ): Promise<BatchTxConfirmationsResponse>;
  ContractCallTxConfirmations(
    request: ContractCallTxConfirmationsRequest
  ): Promise<ContractCallTxConfirmationsResponse>;
  /**
   * pending ethereum signature queries for orchestrators to figure out which
   * signatures they are missing
   * TODO: can/should we group this into one endpoint?
   */
  UnsignedSignerSetTxs(
    request: UnsignedSignerSetTxsRequest
  ): Promise<UnsignedSignerSetTxsResponse>;
  UnsignedBatchTxs(
    request: UnsignedBatchTxsRequest
  ): Promise<UnsignedBatchTxsResponse>;
  UnsignedContractCallTxs(
    request: UnsignedContractCallTxsRequest
  ): Promise<UnsignedContractCallTxsResponse>;
  LastSubmittedEthereumEvent(
    request: LastSubmittedEthereumEventRequest
  ): Promise<LastSubmittedEthereumEventResponse>;
  /**
   * Queries the fees for all pending batches, results are returned in sdk.Coin
   * (fee_amount_int)(contract_address) style
   */
  BatchTxFees(request: BatchTxFeesRequest): Promise<BatchTxFeesResponse>;
  /** Query for info about denoms tracked by gravity */
  ERC20ToDenom(request: ERC20ToDenomRequest): Promise<ERC20ToDenomResponse>;
  /**
   * DenomToERC20Params implements a query that allows ERC-20 parameter
   * information to be retrieved by a Cosmos base denomination.
   */
  DenomToERC20Params(
    request: DenomToERC20ParamsRequest
  ): Promise<DenomToERC20ParamsResponse>;
  /** Query for info about denoms tracked by gravity */
  DenomToERC20(request: DenomToERC20Request): Promise<DenomToERC20Response>;
  /** Query for batch send to ethereums */
  BatchedSendToEthereums(
    request: BatchedSendToEthereumsRequest
  ): Promise<BatchedSendToEthereumsResponse>;
  /** Query for unbatched send to ethereums */
  UnbatchedSendToEthereums(
    request: UnbatchedSendToEthereumsRequest
  ): Promise<UnbatchedSendToEthereumsResponse>;
  /** delegate keys */
  DelegateKeysByValidator(
    request: DelegateKeysByValidatorRequest
  ): Promise<DelegateKeysByValidatorResponse>;
  DelegateKeysByEthereumSigner(
    request: DelegateKeysByEthereumSignerRequest
  ): Promise<DelegateKeysByEthereumSignerResponse>;
  DelegateKeysByOrchestrator(
    request: DelegateKeysByOrchestratorRequest
  ): Promise<DelegateKeysByOrchestratorResponse>;
  DelegateKeys(request: DelegateKeysRequest): Promise<DelegateKeysResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
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
    this.ContractCallTxConfirmations =
      this.ContractCallTxConfirmations.bind(this);
    this.UnsignedSignerSetTxs = this.UnsignedSignerSetTxs.bind(this);
    this.UnsignedBatchTxs = this.UnsignedBatchTxs.bind(this);
    this.UnsignedContractCallTxs = this.UnsignedContractCallTxs.bind(this);
    this.LastSubmittedEthereumEvent =
      this.LastSubmittedEthereumEvent.bind(this);
    this.BatchTxFees = this.BatchTxFees.bind(this);
    this.ERC20ToDenom = this.ERC20ToDenom.bind(this);
    this.DenomToERC20Params = this.DenomToERC20Params.bind(this);
    this.DenomToERC20 = this.DenomToERC20.bind(this);
    this.BatchedSendToEthereums = this.BatchedSendToEthereums.bind(this);
    this.UnbatchedSendToEthereums = this.UnbatchedSendToEthereums.bind(this);
    this.DelegateKeysByValidator = this.DelegateKeysByValidator.bind(this);
    this.DelegateKeysByEthereumSigner =
      this.DelegateKeysByEthereumSigner.bind(this);
    this.DelegateKeysByOrchestrator =
      this.DelegateKeysByOrchestrator.bind(this);
    this.DelegateKeys = this.DelegateKeys.bind(this);
  }
  Params(request: ParamsRequest): Promise<ParamsResponse> {
    const data = ParamsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "Params", data);
    return promise.then((data) => ParamsResponse.decode(new _m0.Reader(data)));
  }

  SignerSetTx(request: SignerSetTxRequest): Promise<SignerSetTxResponse> {
    const data = SignerSetTxRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "SignerSetTx", data);
    return promise.then((data) =>
      SignerSetTxResponse.decode(new _m0.Reader(data))
    );
  }

  LatestSignerSetTx(
    request: LatestSignerSetTxRequest
  ): Promise<SignerSetTxResponse> {
    const data = LatestSignerSetTxRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "LatestSignerSetTx",
      data
    );
    return promise.then((data) =>
      SignerSetTxResponse.decode(new _m0.Reader(data))
    );
  }

  BatchTx(request: BatchTxRequest): Promise<BatchTxResponse> {
    const data = BatchTxRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "BatchTx", data);
    return promise.then((data) => BatchTxResponse.decode(new _m0.Reader(data)));
  }

  ContractCallTx(
    request: ContractCallTxRequest
  ): Promise<ContractCallTxResponse> {
    const data = ContractCallTxRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "ContractCallTx",
      data
    );
    return promise.then((data) =>
      ContractCallTxResponse.decode(new _m0.Reader(data))
    );
  }

  SignerSetTxs(request: SignerSetTxsRequest): Promise<SignerSetTxsResponse> {
    const data = SignerSetTxsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "SignerSetTxs", data);
    return promise.then((data) =>
      SignerSetTxsResponse.decode(new _m0.Reader(data))
    );
  }

  BatchTxs(request: BatchTxsRequest): Promise<BatchTxsResponse> {
    const data = BatchTxsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "BatchTxs", data);
    return promise.then((data) =>
      BatchTxsResponse.decode(new _m0.Reader(data))
    );
  }

  ContractCallTxs(
    request: ContractCallTxsRequest
  ): Promise<ContractCallTxsResponse> {
    const data = ContractCallTxsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "ContractCallTxs",
      data
    );
    return promise.then((data) =>
      ContractCallTxsResponse.decode(new _m0.Reader(data))
    );
  }

  SignerSetTxConfirmations(
    request: SignerSetTxConfirmationsRequest
  ): Promise<SignerSetTxConfirmationsResponse> {
    const data = SignerSetTxConfirmationsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "SignerSetTxConfirmations",
      data
    );
    return promise.then((data) =>
      SignerSetTxConfirmationsResponse.decode(new _m0.Reader(data))
    );
  }

  BatchTxConfirmations(
    request: BatchTxConfirmationsRequest
  ): Promise<BatchTxConfirmationsResponse> {
    const data = BatchTxConfirmationsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "BatchTxConfirmations",
      data
    );
    return promise.then((data) =>
      BatchTxConfirmationsResponse.decode(new _m0.Reader(data))
    );
  }

  ContractCallTxConfirmations(
    request: ContractCallTxConfirmationsRequest
  ): Promise<ContractCallTxConfirmationsResponse> {
    const data = ContractCallTxConfirmationsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "ContractCallTxConfirmations",
      data
    );
    return promise.then((data) =>
      ContractCallTxConfirmationsResponse.decode(new _m0.Reader(data))
    );
  }

  UnsignedSignerSetTxs(
    request: UnsignedSignerSetTxsRequest
  ): Promise<UnsignedSignerSetTxsResponse> {
    const data = UnsignedSignerSetTxsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "UnsignedSignerSetTxs",
      data
    );
    return promise.then((data) =>
      UnsignedSignerSetTxsResponse.decode(new _m0.Reader(data))
    );
  }

  UnsignedBatchTxs(
    request: UnsignedBatchTxsRequest
  ): Promise<UnsignedBatchTxsResponse> {
    const data = UnsignedBatchTxsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "UnsignedBatchTxs",
      data
    );
    return promise.then((data) =>
      UnsignedBatchTxsResponse.decode(new _m0.Reader(data))
    );
  }

  UnsignedContractCallTxs(
    request: UnsignedContractCallTxsRequest
  ): Promise<UnsignedContractCallTxsResponse> {
    const data = UnsignedContractCallTxsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "UnsignedContractCallTxs",
      data
    );
    return promise.then((data) =>
      UnsignedContractCallTxsResponse.decode(new _m0.Reader(data))
    );
  }

  LastSubmittedEthereumEvent(
    request: LastSubmittedEthereumEventRequest
  ): Promise<LastSubmittedEthereumEventResponse> {
    const data = LastSubmittedEthereumEventRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "LastSubmittedEthereumEvent",
      data
    );
    return promise.then((data) =>
      LastSubmittedEthereumEventResponse.decode(new _m0.Reader(data))
    );
  }

  BatchTxFees(request: BatchTxFeesRequest): Promise<BatchTxFeesResponse> {
    const data = BatchTxFeesRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "BatchTxFees", data);
    return promise.then((data) =>
      BatchTxFeesResponse.decode(new _m0.Reader(data))
    );
  }

  ERC20ToDenom(request: ERC20ToDenomRequest): Promise<ERC20ToDenomResponse> {
    const data = ERC20ToDenomRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "ERC20ToDenom", data);
    return promise.then((data) =>
      ERC20ToDenomResponse.decode(new _m0.Reader(data))
    );
  }

  DenomToERC20Params(
    request: DenomToERC20ParamsRequest
  ): Promise<DenomToERC20ParamsResponse> {
    const data = DenomToERC20ParamsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "DenomToERC20Params",
      data
    );
    return promise.then((data) =>
      DenomToERC20ParamsResponse.decode(new _m0.Reader(data))
    );
  }

  DenomToERC20(request: DenomToERC20Request): Promise<DenomToERC20Response> {
    const data = DenomToERC20Request.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "DenomToERC20", data);
    return promise.then((data) =>
      DenomToERC20Response.decode(new _m0.Reader(data))
    );
  }

  BatchedSendToEthereums(
    request: BatchedSendToEthereumsRequest
  ): Promise<BatchedSendToEthereumsResponse> {
    const data = BatchedSendToEthereumsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "BatchedSendToEthereums",
      data
    );
    return promise.then((data) =>
      BatchedSendToEthereumsResponse.decode(new _m0.Reader(data))
    );
  }

  UnbatchedSendToEthereums(
    request: UnbatchedSendToEthereumsRequest
  ): Promise<UnbatchedSendToEthereumsResponse> {
    const data = UnbatchedSendToEthereumsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "UnbatchedSendToEthereums",
      data
    );
    return promise.then((data) =>
      UnbatchedSendToEthereumsResponse.decode(new _m0.Reader(data))
    );
  }

  DelegateKeysByValidator(
    request: DelegateKeysByValidatorRequest
  ): Promise<DelegateKeysByValidatorResponse> {
    const data = DelegateKeysByValidatorRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "DelegateKeysByValidator",
      data
    );
    return promise.then((data) =>
      DelegateKeysByValidatorResponse.decode(new _m0.Reader(data))
    );
  }

  DelegateKeysByEthereumSigner(
    request: DelegateKeysByEthereumSignerRequest
  ): Promise<DelegateKeysByEthereumSignerResponse> {
    const data = DelegateKeysByEthereumSignerRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "DelegateKeysByEthereumSigner",
      data
    );
    return promise.then((data) =>
      DelegateKeysByEthereumSignerResponse.decode(new _m0.Reader(data))
    );
  }

  DelegateKeysByOrchestrator(
    request: DelegateKeysByOrchestratorRequest
  ): Promise<DelegateKeysByOrchestratorResponse> {
    const data = DelegateKeysByOrchestratorRequest.encode(request).finish();
    const promise = this.rpc.request(
      "gravity.v1.Query",
      "DelegateKeysByOrchestrator",
      data
    );
    return promise.then((data) =>
      DelegateKeysByOrchestratorResponse.decode(new _m0.Reader(data))
    );
  }

  DelegateKeys(request: DelegateKeysRequest): Promise<DelegateKeysResponse> {
    const data = DelegateKeysRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "DelegateKeys", data);
    return promise.then((data) =>
      DelegateKeysResponse.decode(new _m0.Reader(data))
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
