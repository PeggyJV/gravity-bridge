/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Any } from "../../google/protobuf/any";
import { Coin } from "../../cosmos/base/v1beta1/coin";

export const protobufPackage = "gravity.v1";

/**
 * EthereumEventVoteRecord is an event that is pending of confirmation by 2/3 of
 * the signer set. The event is then attested and executed in the state machine
 * once the required threshold is met.
 */
export interface EthereumEventVoteRecord {
  event?: Any;
  votes: string[];
  accepted: boolean;
}

/**
 * LatestEthereumBlockHeight defines the latest observed ethereum block height
 * and the corresponding timestamp value in nanoseconds.
 */
export interface LatestEthereumBlockHeight {
  ethereumHeight: Long;
  cosmosHeight: Long;
}

/**
 * EthereumSigner represents a cosmos validator with its corresponding bridge
 * operator ethereum address and its staking consensus power.
 */
export interface EthereumSigner {
  power: Long;
  ethereumAddress: string;
}

/**
 * SignerSetTx is the Ethereum Bridge multisig set that relays
 * transactions the two chains. The staking validators keep ethereum keys which
 * are used to check signatures on Ethereum in order to get significant gas
 * savings.
 */
export interface SignerSetTx {
  nonce: Long;
  height: Long;
  signers: EthereumSigner[];
}

/**
 * BatchTx represents a batch of transactions going from Cosmos to Ethereum.
 * Batch txs are are identified by a unique hash and the token contract that is
 * shared by all the SendToEthereum
 */
export interface BatchTx {
  batchNonce: Long;
  timeout: Long;
  transactions: SendToEthereum[];
  tokenContract: string;
  height: Long;
}

/**
 * SendToEthereum represents an individual SendToEthereum from Cosmos to
 * Ethereum
 */
export interface SendToEthereum {
  id: Long;
  sender: string;
  ethereumRecipient: string;
  erc20Token?: ERC20Token;
  erc20Fee?: ERC20Token;
}

/**
 * ContractCallTx represents an individual arbitrary logic call transaction
 * from Cosmos to Ethereum.
 */
export interface ContractCallTx {
  invalidationNonce: Long;
  invalidationScope: Uint8Array;
  address: string;
  payload: Uint8Array;
  timeout: Long;
  tokens: ERC20Token[];
  fees: ERC20Token[];
  height: Long;
}

export interface ERC20Token {
  contract: string;
  amount: string;
}

export interface IDSet {
  ids: Long[];
}

export interface CommunityPoolEthereumSpendProposal {
  title: string;
  description: string;
  recipient: string;
  amount?: Coin;
  bridgeFee?: Coin;
}

/**
 * This format of the community spend Ethereum proposal is specifically for
 * the CLI to allow simple text serialization.
 */
export interface CommunityPoolEthereumSpendProposalForCLI {
  title: string;
  description: string;
  recipient: string;
  amount: string;
  bridgeFee: string;
  deposit: string;
}

const baseEthereumEventVoteRecord: object = { votes: "", accepted: false };

export const EthereumEventVoteRecord = {
  encode(
    message: EthereumEventVoteRecord,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.event !== undefined) {
      Any.encode(message.event, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.votes) {
      writer.uint32(18).string(v!);
    }
    if (message.accepted === true) {
      writer.uint32(24).bool(message.accepted);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): EthereumEventVoteRecord {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseEthereumEventVoteRecord,
    } as EthereumEventVoteRecord;
    message.votes = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.event = Any.decode(reader, reader.uint32());
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

  fromJSON(object: any): EthereumEventVoteRecord {
    const message = {
      ...baseEthereumEventVoteRecord,
    } as EthereumEventVoteRecord;
    message.event =
      object.event !== undefined && object.event !== null
        ? Any.fromJSON(object.event)
        : undefined;
    message.votes = (object.votes ?? []).map((e: any) => String(e));
    message.accepted =
      object.accepted !== undefined && object.accepted !== null
        ? Boolean(object.accepted)
        : false;
    return message;
  },

  toJSON(message: EthereumEventVoteRecord): unknown {
    const obj: any = {};
    message.event !== undefined &&
      (obj.event = message.event ? Any.toJSON(message.event) : undefined);
    if (message.votes) {
      obj.votes = message.votes.map((e) => e);
    } else {
      obj.votes = [];
    }
    message.accepted !== undefined && (obj.accepted = message.accepted);
    return obj;
  },

  fromPartial(
    object: DeepPartial<EthereumEventVoteRecord>
  ): EthereumEventVoteRecord {
    const message = {
      ...baseEthereumEventVoteRecord,
    } as EthereumEventVoteRecord;
    message.event =
      object.event !== undefined && object.event !== null
        ? Any.fromPartial(object.event)
        : undefined;
    message.votes = (object.votes ?? []).map((e) => e);
    message.accepted = object.accepted ?? false;
    return message;
  },
};

const baseLatestEthereumBlockHeight: object = {
  ethereumHeight: Long.UZERO,
  cosmosHeight: Long.UZERO,
};

export const LatestEthereumBlockHeight = {
  encode(
    message: LatestEthereumBlockHeight,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.ethereumHeight.isZero()) {
      writer.uint32(8).uint64(message.ethereumHeight);
    }
    if (!message.cosmosHeight.isZero()) {
      writer.uint32(16).uint64(message.cosmosHeight);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): LatestEthereumBlockHeight {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseLatestEthereumBlockHeight,
    } as LatestEthereumBlockHeight;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ethereumHeight = reader.uint64() as Long;
          break;
        case 2:
          message.cosmosHeight = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LatestEthereumBlockHeight {
    const message = {
      ...baseLatestEthereumBlockHeight,
    } as LatestEthereumBlockHeight;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromString(object.ethereumHeight)
        : Long.UZERO;
    message.cosmosHeight =
      object.cosmosHeight !== undefined && object.cosmosHeight !== null
        ? Long.fromString(object.cosmosHeight)
        : Long.UZERO;
    return message;
  },

  toJSON(message: LatestEthereumBlockHeight): unknown {
    const obj: any = {};
    message.ethereumHeight !== undefined &&
      (obj.ethereumHeight = (message.ethereumHeight || Long.UZERO).toString());
    message.cosmosHeight !== undefined &&
      (obj.cosmosHeight = (message.cosmosHeight || Long.UZERO).toString());
    return obj;
  },

  fromPartial(
    object: DeepPartial<LatestEthereumBlockHeight>
  ): LatestEthereumBlockHeight {
    const message = {
      ...baseLatestEthereumBlockHeight,
    } as LatestEthereumBlockHeight;
    message.ethereumHeight =
      object.ethereumHeight !== undefined && object.ethereumHeight !== null
        ? Long.fromValue(object.ethereumHeight)
        : Long.UZERO;
    message.cosmosHeight =
      object.cosmosHeight !== undefined && object.cosmosHeight !== null
        ? Long.fromValue(object.cosmosHeight)
        : Long.UZERO;
    return message;
  },
};

const baseEthereumSigner: object = { power: Long.UZERO, ethereumAddress: "" };

export const EthereumSigner = {
  encode(
    message: EthereumSigner,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.power.isZero()) {
      writer.uint32(8).uint64(message.power);
    }
    if (message.ethereumAddress !== "") {
      writer.uint32(18).string(message.ethereumAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EthereumSigner {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseEthereumSigner } as EthereumSigner;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.power = reader.uint64() as Long;
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

  fromJSON(object: any): EthereumSigner {
    const message = { ...baseEthereumSigner } as EthereumSigner;
    message.power =
      object.power !== undefined && object.power !== null
        ? Long.fromString(object.power)
        : Long.UZERO;
    message.ethereumAddress =
      object.ethereumAddress !== undefined && object.ethereumAddress !== null
        ? String(object.ethereumAddress)
        : "";
    return message;
  },

  toJSON(message: EthereumSigner): unknown {
    const obj: any = {};
    message.power !== undefined &&
      (obj.power = (message.power || Long.UZERO).toString());
    message.ethereumAddress !== undefined &&
      (obj.ethereumAddress = message.ethereumAddress);
    return obj;
  },

  fromPartial(object: DeepPartial<EthereumSigner>): EthereumSigner {
    const message = { ...baseEthereumSigner } as EthereumSigner;
    message.power =
      object.power !== undefined && object.power !== null
        ? Long.fromValue(object.power)
        : Long.UZERO;
    message.ethereumAddress = object.ethereumAddress ?? "";
    return message;
  },
};

const baseSignerSetTx: object = { nonce: Long.UZERO, height: Long.UZERO };

export const SignerSetTx = {
  encode(
    message: SignerSetTx,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.nonce.isZero()) {
      writer.uint32(8).uint64(message.nonce);
    }
    if (!message.height.isZero()) {
      writer.uint32(16).uint64(message.height);
    }
    for (const v of message.signers) {
      EthereumSigner.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SignerSetTx {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseSignerSetTx } as SignerSetTx;
    message.signers = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.nonce = reader.uint64() as Long;
          break;
        case 2:
          message.height = reader.uint64() as Long;
          break;
        case 3:
          message.signers.push(EthereumSigner.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignerSetTx {
    const message = { ...baseSignerSetTx } as SignerSetTx;
    message.nonce =
      object.nonce !== undefined && object.nonce !== null
        ? Long.fromString(object.nonce)
        : Long.UZERO;
    message.height =
      object.height !== undefined && object.height !== null
        ? Long.fromString(object.height)
        : Long.UZERO;
    message.signers = (object.signers ?? []).map((e: any) =>
      EthereumSigner.fromJSON(e)
    );
    return message;
  },

  toJSON(message: SignerSetTx): unknown {
    const obj: any = {};
    message.nonce !== undefined &&
      (obj.nonce = (message.nonce || Long.UZERO).toString());
    message.height !== undefined &&
      (obj.height = (message.height || Long.UZERO).toString());
    if (message.signers) {
      obj.signers = message.signers.map((e) =>
        e ? EthereumSigner.toJSON(e) : undefined
      );
    } else {
      obj.signers = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<SignerSetTx>): SignerSetTx {
    const message = { ...baseSignerSetTx } as SignerSetTx;
    message.nonce =
      object.nonce !== undefined && object.nonce !== null
        ? Long.fromValue(object.nonce)
        : Long.UZERO;
    message.height =
      object.height !== undefined && object.height !== null
        ? Long.fromValue(object.height)
        : Long.UZERO;
    message.signers = (object.signers ?? []).map((e) =>
      EthereumSigner.fromPartial(e)
    );
    return message;
  },
};

const baseBatchTx: object = {
  batchNonce: Long.UZERO,
  timeout: Long.UZERO,
  tokenContract: "",
  height: Long.UZERO,
};

export const BatchTx = {
  encode(
    message: BatchTx,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.batchNonce.isZero()) {
      writer.uint32(8).uint64(message.batchNonce);
    }
    if (!message.timeout.isZero()) {
      writer.uint32(16).uint64(message.timeout);
    }
    for (const v of message.transactions) {
      SendToEthereum.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.tokenContract !== "") {
      writer.uint32(34).string(message.tokenContract);
    }
    if (!message.height.isZero()) {
      writer.uint32(40).uint64(message.height);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BatchTx {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBatchTx } as BatchTx;
    message.transactions = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.batchNonce = reader.uint64() as Long;
          break;
        case 2:
          message.timeout = reader.uint64() as Long;
          break;
        case 3:
          message.transactions.push(
            SendToEthereum.decode(reader, reader.uint32())
          );
          break;
        case 4:
          message.tokenContract = reader.string();
          break;
        case 5:
          message.height = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BatchTx {
    const message = { ...baseBatchTx } as BatchTx;
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromString(object.batchNonce)
        : Long.UZERO;
    message.timeout =
      object.timeout !== undefined && object.timeout !== null
        ? Long.fromString(object.timeout)
        : Long.UZERO;
    message.transactions = (object.transactions ?? []).map((e: any) =>
      SendToEthereum.fromJSON(e)
    );
    message.tokenContract =
      object.tokenContract !== undefined && object.tokenContract !== null
        ? String(object.tokenContract)
        : "";
    message.height =
      object.height !== undefined && object.height !== null
        ? Long.fromString(object.height)
        : Long.UZERO;
    return message;
  },

  toJSON(message: BatchTx): unknown {
    const obj: any = {};
    message.batchNonce !== undefined &&
      (obj.batchNonce = (message.batchNonce || Long.UZERO).toString());
    message.timeout !== undefined &&
      (obj.timeout = (message.timeout || Long.UZERO).toString());
    if (message.transactions) {
      obj.transactions = message.transactions.map((e) =>
        e ? SendToEthereum.toJSON(e) : undefined
      );
    } else {
      obj.transactions = [];
    }
    message.tokenContract !== undefined &&
      (obj.tokenContract = message.tokenContract);
    message.height !== undefined &&
      (obj.height = (message.height || Long.UZERO).toString());
    return obj;
  },

  fromPartial(object: DeepPartial<BatchTx>): BatchTx {
    const message = { ...baseBatchTx } as BatchTx;
    message.batchNonce =
      object.batchNonce !== undefined && object.batchNonce !== null
        ? Long.fromValue(object.batchNonce)
        : Long.UZERO;
    message.timeout =
      object.timeout !== undefined && object.timeout !== null
        ? Long.fromValue(object.timeout)
        : Long.UZERO;
    message.transactions = (object.transactions ?? []).map((e) =>
      SendToEthereum.fromPartial(e)
    );
    message.tokenContract = object.tokenContract ?? "";
    message.height =
      object.height !== undefined && object.height !== null
        ? Long.fromValue(object.height)
        : Long.UZERO;
    return message;
  },
};

const baseSendToEthereum: object = {
  id: Long.UZERO,
  sender: "",
  ethereumRecipient: "",
};

export const SendToEthereum = {
  encode(
    message: SendToEthereum,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
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
      ERC20Token.encode(message.erc20Token, writer.uint32(34).fork()).ldelim();
    }
    if (message.erc20Fee !== undefined) {
      ERC20Token.encode(message.erc20Fee, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendToEthereum {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseSendToEthereum } as SendToEthereum;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint64() as Long;
          break;
        case 2:
          message.sender = reader.string();
          break;
        case 3:
          message.ethereumRecipient = reader.string();
          break;
        case 4:
          message.erc20Token = ERC20Token.decode(reader, reader.uint32());
          break;
        case 5:
          message.erc20Fee = ERC20Token.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SendToEthereum {
    const message = { ...baseSendToEthereum } as SendToEthereum;
    message.id =
      object.id !== undefined && object.id !== null
        ? Long.fromString(object.id)
        : Long.UZERO;
    message.sender =
      object.sender !== undefined && object.sender !== null
        ? String(object.sender)
        : "";
    message.ethereumRecipient =
      object.ethereumRecipient !== undefined &&
      object.ethereumRecipient !== null
        ? String(object.ethereumRecipient)
        : "";
    message.erc20Token =
      object.erc20Token !== undefined && object.erc20Token !== null
        ? ERC20Token.fromJSON(object.erc20Token)
        : undefined;
    message.erc20Fee =
      object.erc20Fee !== undefined && object.erc20Fee !== null
        ? ERC20Token.fromJSON(object.erc20Fee)
        : undefined;
    return message;
  },

  toJSON(message: SendToEthereum): unknown {
    const obj: any = {};
    message.id !== undefined &&
      (obj.id = (message.id || Long.UZERO).toString());
    message.sender !== undefined && (obj.sender = message.sender);
    message.ethereumRecipient !== undefined &&
      (obj.ethereumRecipient = message.ethereumRecipient);
    message.erc20Token !== undefined &&
      (obj.erc20Token = message.erc20Token
        ? ERC20Token.toJSON(message.erc20Token)
        : undefined);
    message.erc20Fee !== undefined &&
      (obj.erc20Fee = message.erc20Fee
        ? ERC20Token.toJSON(message.erc20Fee)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<SendToEthereum>): SendToEthereum {
    const message = { ...baseSendToEthereum } as SendToEthereum;
    message.id =
      object.id !== undefined && object.id !== null
        ? Long.fromValue(object.id)
        : Long.UZERO;
    message.sender = object.sender ?? "";
    message.ethereumRecipient = object.ethereumRecipient ?? "";
    message.erc20Token =
      object.erc20Token !== undefined && object.erc20Token !== null
        ? ERC20Token.fromPartial(object.erc20Token)
        : undefined;
    message.erc20Fee =
      object.erc20Fee !== undefined && object.erc20Fee !== null
        ? ERC20Token.fromPartial(object.erc20Fee)
        : undefined;
    return message;
  },
};

const baseContractCallTx: object = {
  invalidationNonce: Long.UZERO,
  address: "",
  timeout: Long.UZERO,
  height: Long.UZERO,
};

export const ContractCallTx = {
  encode(
    message: ContractCallTx,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
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
    for (const v of message.tokens) {
      ERC20Token.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.fees) {
      ERC20Token.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    if (!message.height.isZero()) {
      writer.uint32(64).uint64(message.height);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ContractCallTx {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseContractCallTx } as ContractCallTx;
    message.tokens = [];
    message.fees = [];
    message.invalidationScope = new Uint8Array();
    message.payload = new Uint8Array();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.invalidationNonce = reader.uint64() as Long;
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
          message.timeout = reader.uint64() as Long;
          break;
        case 6:
          message.tokens.push(ERC20Token.decode(reader, reader.uint32()));
          break;
        case 7:
          message.fees.push(ERC20Token.decode(reader, reader.uint32()));
          break;
        case 8:
          message.height = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractCallTx {
    const message = { ...baseContractCallTx } as ContractCallTx;
    message.invalidationNonce =
      object.invalidationNonce !== undefined &&
      object.invalidationNonce !== null
        ? Long.fromString(object.invalidationNonce)
        : Long.UZERO;
    message.invalidationScope =
      object.invalidationScope !== undefined &&
      object.invalidationScope !== null
        ? bytesFromBase64(object.invalidationScope)
        : new Uint8Array();
    message.address =
      object.address !== undefined && object.address !== null
        ? String(object.address)
        : "";
    message.payload =
      object.payload !== undefined && object.payload !== null
        ? bytesFromBase64(object.payload)
        : new Uint8Array();
    message.timeout =
      object.timeout !== undefined && object.timeout !== null
        ? Long.fromString(object.timeout)
        : Long.UZERO;
    message.tokens = (object.tokens ?? []).map((e: any) =>
      ERC20Token.fromJSON(e)
    );
    message.fees = (object.fees ?? []).map((e: any) => ERC20Token.fromJSON(e));
    message.height =
      object.height !== undefined && object.height !== null
        ? Long.fromString(object.height)
        : Long.UZERO;
    return message;
  },

  toJSON(message: ContractCallTx): unknown {
    const obj: any = {};
    message.invalidationNonce !== undefined &&
      (obj.invalidationNonce = (
        message.invalidationNonce || Long.UZERO
      ).toString());
    message.invalidationScope !== undefined &&
      (obj.invalidationScope = base64FromBytes(
        message.invalidationScope !== undefined
          ? message.invalidationScope
          : new Uint8Array()
      ));
    message.address !== undefined && (obj.address = message.address);
    message.payload !== undefined &&
      (obj.payload = base64FromBytes(
        message.payload !== undefined ? message.payload : new Uint8Array()
      ));
    message.timeout !== undefined &&
      (obj.timeout = (message.timeout || Long.UZERO).toString());
    if (message.tokens) {
      obj.tokens = message.tokens.map((e) =>
        e ? ERC20Token.toJSON(e) : undefined
      );
    } else {
      obj.tokens = [];
    }
    if (message.fees) {
      obj.fees = message.fees.map((e) =>
        e ? ERC20Token.toJSON(e) : undefined
      );
    } else {
      obj.fees = [];
    }
    message.height !== undefined &&
      (obj.height = (message.height || Long.UZERO).toString());
    return obj;
  },

  fromPartial(object: DeepPartial<ContractCallTx>): ContractCallTx {
    const message = { ...baseContractCallTx } as ContractCallTx;
    message.invalidationNonce =
      object.invalidationNonce !== undefined &&
      object.invalidationNonce !== null
        ? Long.fromValue(object.invalidationNonce)
        : Long.UZERO;
    message.invalidationScope = object.invalidationScope ?? new Uint8Array();
    message.address = object.address ?? "";
    message.payload = object.payload ?? new Uint8Array();
    message.timeout =
      object.timeout !== undefined && object.timeout !== null
        ? Long.fromValue(object.timeout)
        : Long.UZERO;
    message.tokens = (object.tokens ?? []).map((e) =>
      ERC20Token.fromPartial(e)
    );
    message.fees = (object.fees ?? []).map((e) => ERC20Token.fromPartial(e));
    message.height =
      object.height !== undefined && object.height !== null
        ? Long.fromValue(object.height)
        : Long.UZERO;
    return message;
  },
};

const baseERC20Token: object = { contract: "", amount: "" };

export const ERC20Token = {
  encode(
    message: ERC20Token,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.contract !== "") {
      writer.uint32(10).string(message.contract);
    }
    if (message.amount !== "") {
      writer.uint32(18).string(message.amount);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ERC20Token {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseERC20Token } as ERC20Token;
    while (reader.pos < end) {
      const tag = reader.uint32();
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

  fromJSON(object: any): ERC20Token {
    const message = { ...baseERC20Token } as ERC20Token;
    message.contract =
      object.contract !== undefined && object.contract !== null
        ? String(object.contract)
        : "";
    message.amount =
      object.amount !== undefined && object.amount !== null
        ? String(object.amount)
        : "";
    return message;
  },

  toJSON(message: ERC20Token): unknown {
    const obj: any = {};
    message.contract !== undefined && (obj.contract = message.contract);
    message.amount !== undefined && (obj.amount = message.amount);
    return obj;
  },

  fromPartial(object: DeepPartial<ERC20Token>): ERC20Token {
    const message = { ...baseERC20Token } as ERC20Token;
    message.contract = object.contract ?? "";
    message.amount = object.amount ?? "";
    return message;
  },
};

const baseIDSet: object = { ids: Long.UZERO };

export const IDSet = {
  encode(message: IDSet, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    writer.uint32(10).fork();
    for (const v of message.ids) {
      writer.uint64(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDSet {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseIDSet } as IDSet;
    message.ids = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.ids.push(reader.uint64() as Long);
            }
          } else {
            message.ids.push(reader.uint64() as Long);
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IDSet {
    const message = { ...baseIDSet } as IDSet;
    message.ids = (object.ids ?? []).map((e: any) => Long.fromString(e));
    return message;
  },

  toJSON(message: IDSet): unknown {
    const obj: any = {};
    if (message.ids) {
      obj.ids = message.ids.map((e) => (e || Long.UZERO).toString());
    } else {
      obj.ids = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<IDSet>): IDSet {
    const message = { ...baseIDSet } as IDSet;
    message.ids = (object.ids ?? []).map((e) => Long.fromValue(e));
    return message;
  },
};

const baseCommunityPoolEthereumSpendProposal: object = {
  title: "",
  description: "",
  recipient: "",
};

export const CommunityPoolEthereumSpendProposal = {
  encode(
    message: CommunityPoolEthereumSpendProposal,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.recipient !== "") {
      writer.uint32(26).string(message.recipient);
    }
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(34).fork()).ldelim();
    }
    if (message.bridgeFee !== undefined) {
      Coin.encode(message.bridgeFee, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CommunityPoolEthereumSpendProposal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseCommunityPoolEthereumSpendProposal,
    } as CommunityPoolEthereumSpendProposal;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.recipient = reader.string();
          break;
        case 4:
          message.amount = Coin.decode(reader, reader.uint32());
          break;
        case 5:
          message.bridgeFee = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CommunityPoolEthereumSpendProposal {
    const message = {
      ...baseCommunityPoolEthereumSpendProposal,
    } as CommunityPoolEthereumSpendProposal;
    message.title =
      object.title !== undefined && object.title !== null
        ? String(object.title)
        : "";
    message.description =
      object.description !== undefined && object.description !== null
        ? String(object.description)
        : "";
    message.recipient =
      object.recipient !== undefined && object.recipient !== null
        ? String(object.recipient)
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

  toJSON(message: CommunityPoolEthereumSpendProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined &&
      (obj.description = message.description);
    message.recipient !== undefined && (obj.recipient = message.recipient);
    message.amount !== undefined &&
      (obj.amount = message.amount ? Coin.toJSON(message.amount) : undefined);
    message.bridgeFee !== undefined &&
      (obj.bridgeFee = message.bridgeFee
        ? Coin.toJSON(message.bridgeFee)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<CommunityPoolEthereumSpendProposal>
  ): CommunityPoolEthereumSpendProposal {
    const message = {
      ...baseCommunityPoolEthereumSpendProposal,
    } as CommunityPoolEthereumSpendProposal;
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.recipient = object.recipient ?? "";
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

const baseCommunityPoolEthereumSpendProposalForCLI: object = {
  title: "",
  description: "",
  recipient: "",
  amount: "",
  bridgeFee: "",
  deposit: "",
};

export const CommunityPoolEthereumSpendProposalForCLI = {
  encode(
    message: CommunityPoolEthereumSpendProposalForCLI,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.recipient !== "") {
      writer.uint32(26).string(message.recipient);
    }
    if (message.amount !== "") {
      writer.uint32(34).string(message.amount);
    }
    if (message.bridgeFee !== "") {
      writer.uint32(42).string(message.bridgeFee);
    }
    if (message.deposit !== "") {
      writer.uint32(50).string(message.deposit);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CommunityPoolEthereumSpendProposalForCLI {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseCommunityPoolEthereumSpendProposalForCLI,
    } as CommunityPoolEthereumSpendProposalForCLI;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.recipient = reader.string();
          break;
        case 4:
          message.amount = reader.string();
          break;
        case 5:
          message.bridgeFee = reader.string();
          break;
        case 6:
          message.deposit = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CommunityPoolEthereumSpendProposalForCLI {
    const message = {
      ...baseCommunityPoolEthereumSpendProposalForCLI,
    } as CommunityPoolEthereumSpendProposalForCLI;
    message.title =
      object.title !== undefined && object.title !== null
        ? String(object.title)
        : "";
    message.description =
      object.description !== undefined && object.description !== null
        ? String(object.description)
        : "";
    message.recipient =
      object.recipient !== undefined && object.recipient !== null
        ? String(object.recipient)
        : "";
    message.amount =
      object.amount !== undefined && object.amount !== null
        ? String(object.amount)
        : "";
    message.bridgeFee =
      object.bridgeFee !== undefined && object.bridgeFee !== null
        ? String(object.bridgeFee)
        : "";
    message.deposit =
      object.deposit !== undefined && object.deposit !== null
        ? String(object.deposit)
        : "";
    return message;
  },

  toJSON(message: CommunityPoolEthereumSpendProposalForCLI): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined &&
      (obj.description = message.description);
    message.recipient !== undefined && (obj.recipient = message.recipient);
    message.amount !== undefined && (obj.amount = message.amount);
    message.bridgeFee !== undefined && (obj.bridgeFee = message.bridgeFee);
    message.deposit !== undefined && (obj.deposit = message.deposit);
    return obj;
  },

  fromPartial(
    object: DeepPartial<CommunityPoolEthereumSpendProposalForCLI>
  ): CommunityPoolEthereumSpendProposalForCLI {
    const message = {
      ...baseCommunityPoolEthereumSpendProposalForCLI,
    } as CommunityPoolEthereumSpendProposalForCLI;
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.recipient = object.recipient ?? "";
    message.amount = object.amount ?? "";
    message.bridgeFee = object.bridgeFee ?? "";
    message.deposit = object.deposit ?? "";
    return message;
  },
};

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
