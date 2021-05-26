package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	var params types.Params
	k.paramSpace.GetParamSet(sdk.UnwrapSDKContext(c), &params)
	return &types.ParamsResponse{Params: params}, nil
}

func (k Keeper) SignerSetTx(c context.Context, req *types.SignerSetTxRequest) (*types.SignerSetTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// TODO: audit once we finalize storage
	var otx types.OutgoingTx

	// given a 0 nonce, we will retrieve the latest by retrieving the last value off of
	// the reverse iterator for signer sets. As nonces only increase, this should result in the latest signer set.
	if req.Nonce == 0 {
		store := prefix.NewStore(ctx.KVStore(k.storeKey), append([]byte{types.OutgoingTxKey}, types.SignerSetTxPrefixByte))
		iter := store.ReverseIterator(nil, nil)
		defer iter.Close()

		var any cdctypes.Any
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &any)

		if err := k.cdc.UnpackAny(&any, &otx); err != nil {
			return nil, err
		}
	} else {
		storeIndex := sdk.Uint64ToBigEndian(req.Nonce)
		otx = k.GetOutgoingTx(sdk.UnwrapSDKContext(c), types.MakeOutgoingTxKey(storeIndex))
		if otx == nil {
			return nil, status.Errorf(codes.InvalidArgument, "no signer set found for %d", req.Nonce)
		}
	}

	ss, ok := otx.(*types.SignerSetTx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't cast to signer set for %d", req.Nonce)
	}

	// TODO: ensure that latest signer set tx nonce index is set properly

	return &types.SignerSetTxResponse{SignerSet: ss}, nil
}

func (k Keeper) BatchTx(c context.Context, req *types.BatchTxRequest) (*types.BatchTxResponse, error) {
	if !common.IsHexAddress(req.ContractAddress) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid hex address %s", req.ContractAddress)
	}
	ctx := sdk.UnwrapSDKContext(c)
	res := &types.BatchTxResponse{}

	if req.Nonce == 0 {
		// given a 0 nonce, we will retrieve the latest by iterating through batch txs in reverse,
		// as nonces should be increasing over time. We must iterator through potentially multiple
		// batch txs because we must compare the contract address
		store := prefix.NewStore(ctx.KVStore(k.storeKey), append([]byte{types.OutgoingTxKey}, types.BatchTxPrefixByte))
		iter := store.ReverseIterator(nil, nil)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			var any cdctypes.Any
			k.cdc.MustUnmarshalBinaryBare(iter.Value(), &any)
			var otx types.OutgoingTx
			if err := k.cdc.UnpackAny(&any, &otx); err != nil {
				panic(err)
			}
			batch, ok := otx.(*types.BatchTx)
			if !ok {
				return nil, status.Errorf(codes.InvalidArgument, "couldn't cast to batch tx for %d %s", req.Nonce, req.ContractAddress)
			}
			if batch.TokenContract == req.ContractAddress {
				res.Batch = batch
				break
			}
		}
	} else {
		// TODO: audit once we finalize storage
		storeIndex := append(sdk.Uint64ToBigEndian(req.Nonce), common.Hex2Bytes(req.ContractAddress)...)
		otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), types.MakeOutgoingTxKey(storeIndex))
		if otx == nil {
			return nil, status.Errorf(codes.InvalidArgument, "no batch tx found for %d %s", req.Nonce, req.ContractAddress)
		}
		batch, ok := otx.(*types.BatchTx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "couldn't cast to batch tx for %d %s", req.Nonce, req.ContractAddress)
		}
		res.Batch = batch
	}

	return res, nil
}

func (k Keeper) ContractCallTx(c context.Context, req *types.ContractCallTxRequest) (*types.ContractCallTxResponse, error) {
	storeIndex := append(sdk.Uint64ToBigEndian(req.InvalidationNonce), req.InvalidationScope...)
	otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), types.MakeOutgoingTxKey(storeIndex))
	if otx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no contract call found for %d %s", req.InvalidationNonce, req.InvalidationScope)
	}

	cctx, ok := otx.(*types.ContractCallTx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't cast to contract call for %d %s", req.InvalidationNonce, req.InvalidationScope)
	}

	// TODO: figure out how to call latest

	return &types.ContractCallTxResponse{LogicCall: cctx}, nil
}

func (k Keeper) SignerSetTxs(c context.Context, req *types.SignerSetTxsRequest) (*types.SignerSetTxsResponse, error) {
	var signers []*types.SignerSetTx
	pageRes, err := k.PaginateOutgoingTxsByType(sdk.UnwrapSDKContext(c), req.Pagination, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) (hit bool) {
		signer, ok := otx.(*types.SignerSetTx)
		if !ok {
			panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for %s", otx))
		}
		signers = append(signers, signer)

		return true
	})
	if err != nil {
		return nil, err
	}

	return &types.SignerSetTxsResponse{SignerSets: signers, Pagination: pageRes}, nil
}

func (k Keeper) BatchTxs(c context.Context, req *types.BatchTxsRequest) (*types.BatchTxsResponse, error) {
	var batches []*types.BatchTx
	pageRes, err := k.PaginateOutgoingTxsByType(sdk.UnwrapSDKContext(c), req.Pagination, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) (hit bool) {
		batch, ok := otx.(*types.BatchTx)
		if !ok {
			panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to batch tx for %s", otx))
		}
		batches = append(batches, batch)
		return true
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchTxsResponse{Batches: batches, Pagination: pageRes}, nil
}

func (k Keeper) ContractCallTxs(c context.Context, req *types.ContractCallTxsRequest) (*types.ContractCallTxsResponse, error) {
	var calls []*types.ContractCallTx
	pageRes, err := k.PaginateOutgoingTxsByType(sdk.UnwrapSDKContext(c), req.Pagination, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) (hit bool) {
		call, ok := otx.(*types.ContractCallTx)
		if !ok {
			panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to contract call for %s", otx))
		}
		calls = append(calls, call)
		return true
	})
	if err != nil {
		return nil, err
	}

	return &types.ContractCallTxsResponse{Calls: calls, Pagination: pageRes}, nil
}

func (k Keeper) SignerSetTxEthereumSignatures(c context.Context, req *types.SignerSetTxEthereumSignaturesRequest) (*types.SignerSetTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeSignerSetTxKey(req.Nonce)
	if req.Address != "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.SignerSetTxEthereumSignaturesResponse{Signature: []hexutil.Bytes{k.getEthereumSignature(ctx, key, val)}}, nil
	}

	var out []hexutil.Bytes
	k.iterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.SignerSetTxEthereumSignaturesResponse{Signature: out}, nil
}

func (k Keeper) BatchTxEthereumSignatures(c context.Context, req *types.BatchTxEthereumSignaturesRequest) (*types.BatchTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeBatchTxKey(common.HexToAddress(req.ContractAddress), req.Nonce)
	if req.Address != "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.BatchTxEthereumSignaturesResponse{Signature: [][]byte{k.getEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.iterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.BatchTxEthereumSignaturesResponse{Signature: out}, nil
}

func (k Keeper) ContractCallTxEthereumSignatures(c context.Context, req *types.ContractCallTxEthereumSignaturesRequest) (*types.ContractCallTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeContractCallTxKey(req.InvalidationScope, req.InvalidationNonce)
	if req.Address != "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.ContractCallTxEthereumSignaturesResponse{Signature: [][]byte{k.getEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.iterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.ContractCallTxEthereumSignaturesResponse{Signature: out}, nil
}

func (k Keeper) PendingSignerSetTxEthereumSignatures(c context.Context, req *types.PendingSignerSetTxEthereumSignaturesRequest) (*types.PendingSignerSetTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	var signerSets []*types.SignerSetTx
	k.IterateOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, otx.GetStoreIndex(), val)
		if len(sig) == 0 { // it's pending
			signerSet, ok := otx.(*types.SignerSetTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for %s", otx))
			}
			signerSets = append(signerSets, signerSet)
		}
		return false
	})
	return &types.PendingSignerSetTxEthereumSignaturesResponse{SignerSets: signerSets}, nil
}

func (k Keeper) PendingBatchTxEthereumSignatures(c context.Context, req *types.PendingBatchTxEthereumSignaturesRequest) (*types.PendingBatchTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	var batches []*types.BatchTx
	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, otx.GetStoreIndex(), val)
		if len(sig) == 0 { // it's pending
			batch, ok := otx.(*types.BatchTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to batch tx for %s", otx))
			}
			batches = append(batches, batch)
		}
		return false
	})
	return &types.PendingBatchTxEthereumSignaturesResponse{Batches: batches}, nil
}

func (k Keeper) PendingContractCallTxEthereumSignatures(c context.Context, req *types.PendingContractCallTxEthereumSignaturesRequest) (*types.PendingContractCallTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	var calls []*types.ContractCallTx
	k.IterateOutgoingTxsByType(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, otx.GetStoreIndex(), val)
		if len(sig) == 0 { // it's pending
			call, ok := otx.(*types.ContractCallTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to contract call for %s", otx))
			}
			calls = append(calls, call)
		}
		return false
	})
	return &types.PendingContractCallTxEthereumSignaturesResponse{Calls: calls}, nil
}

func (k Keeper) LastSubmittedEthereumEvent(c context.Context, req *types.LastSubmittedEthereumEventRequest) (*types.LastSubmittedEthereumEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	res := &types.LastSubmittedEthereumEventResponse{}

	valAddr, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	res.EventNonce = k.getLastEventNonceByValidator(ctx, valAddr)

	return &types.LastSubmittedEthereumEventResponse{}, nil
}

func (k Keeper) BatchTxFees(c context.Context, req *types.BatchTxFeesRequest) (*types.BatchTxFeesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	res := &types.BatchTxFeesResponse{}

	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		btx, _ := otx.(*types.BatchTx)
		for _, tx := range btx.Transactions {
			res.Fees = append(res.Fees, tx.Erc20Fee.GravityCoin())
		}
		return false
	})

	return res, nil
}

func (k Keeper) ERC20ToDenom(c context.Context, req *types.ERC20ToDenomRequest) (*types.ERC20ToDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	cosmosOriginated, denom := k.ERC20ToDenomLookup(ctx, req.Erc20)
	res := &types.ERC20ToDenomResponse{
		Denom:            denom,
		CosmosOriginated: cosmosOriginated,
	}
	return res, nil
}

func (k Keeper) DenomToERC20(c context.Context, req *types.DenomToERC20Request) (*types.DenomToERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(c)
	cosmosOriginated, erc20, err := k.DenomToERC20Lookup(ctx, req.Denom)
	if err != nil {
		return nil, err
	}
	res := &types.DenomToERC20Response{
		Erc20:            erc20.Hex(),
		CosmosOriginated: cosmosOriginated,
	}
	return res, nil
}

func (k Keeper) BatchedSendToEthereums(c context.Context, req *types.BatchedSendToEthereumsRequest) (*types.BatchedSendToEthereumsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	res := &types.BatchedSendToEthereumsResponse{}

	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(_ []byte, outgoing types.OutgoingTx) bool {
		batchTx := outgoing.(*types.BatchTx)
		for _, ste := range batchTx.Transactions {
			if ste.Sender == req.SenderAddress {
				res.SendToEthereums = append(res.SendToEthereums, ste)
			}
		}

		return false
	})

	return res, nil
}

func (k Keeper) UnbatchedSendToEthereums(c context.Context, req *types.UnbatchedSendToEthereumsRequest) (*types.UnbatchedSendToEthereumsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	res := &types.UnbatchedSendToEthereumsResponse{}

	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.SendToEthereumKey})
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var ste types.SendToEthereum
		k.cdc.MustUnmarshalBinaryBare(value, &ste)
		if ste.Sender == req.SenderAddress {
			res.SendToEthereums = append(res.SendToEthereums, &ste)
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	res.Pagination = pageRes

	return res, nil
}

func (k Keeper) DelegateKeysByValidator(c context.Context, req *types.DelegateKeysByValidatorRequest) (*types.DelegateKeysByValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	ethAddr := k.GetValidatorEthereumAddress(ctx, valAddr)
	orchAddr := k.GetEthereumOrchestratorAddress(ctx, ethAddr)
	res := &types.DelegateKeysByValidatorResponse{
		EthAddress:          ethAddr.Hex(),
		OrchestratorAddress: orchAddr.String(),
	}
	return res, nil
}

func (k Keeper) DelegateKeysByEthereumSigner(c context.Context, req *types.DelegateKeysByEthereumSignerRequest) (*types.DelegateKeysByEthereumSignerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if !common.IsHexAddress(req.EthereumSigner) {
		return nil, nil // TODO(levi) make and return an error
	}
	ethAddr := common.HexToAddress(req.EthereumSigner)
	orchAddr := k.GetEthereumOrchestratorAddress(ctx, ethAddr)
	valAddr := k.GetOrchestratorValidatorAddress(ctx, orchAddr)
	res := &types.DelegateKeysByEthereumSignerResponse{
		ValidatorAddress:    valAddr.String(),
		OrchestratorAddress: orchAddr.String(),
	}
	return res, nil
}

func (k Keeper) DelegateKeysByOrchestrator(c context.Context, req *types.DelegateKeysByOrchestratorRequest) (*types.DelegateKeysByOrchestratorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	orchAddr, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, err
	}
	valAddr := k.GetOrchestratorValidatorAddress(ctx, orchAddr)
	ethAddr := k.GetValidatorEthereumAddress(ctx, valAddr)
	res := &types.DelegateKeysByOrchestratorResponse{
		ValidatorAddress: valAddr.String(),
		EthereumSigner:   ethAddr.Hex(),
	}
	return res, nil
}
