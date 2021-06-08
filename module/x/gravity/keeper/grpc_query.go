package keeper

import (
	"context"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	var params types.Params
	k.paramSpace.GetParamSet(sdk.UnwrapSDKContext(c), &params)
	return &types.ParamsResponse{Params: params}, nil
}

func (k Keeper) LatestSignerSetTx(c context.Context, req *types.LatestSignerSetTxRequest) (*types.SignerSetTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// TODO: audit once we finalize storage
	var otx types.OutgoingTx

	store := prefix.NewStore(ctx.KVStore(k.storeKey), append([]byte{types.OutgoingTxKey}, types.SignerSetTxPrefixByte))
	iter := store.ReverseIterator(nil, nil)
	defer iter.Close()

	var any cdctypes.Any
	k.cdc.MustUnmarshalBinaryBare(iter.Value(), &any)

	if err := k.cdc.UnpackAny(&any, &otx); err != nil {
		return nil, err
	}

	ss, ok := otx.(*types.SignerSetTx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't cast to signer set for latest")
	}

	return &types.SignerSetTxResponse{SignerSet: ss}, nil
}

func (k Keeper) SignerSetTx(c context.Context, req *types.SignerSetTxRequest) (*types.SignerSetTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	key := types.MakeSignerSetTxKey(req.SignerSetNonce)
	otx := k.GetOutgoingTx(ctx, key)
	if otx == nil {
		return &types.SignerSetTxResponse{}, nil
	}

	ss, ok := otx.(*types.SignerSetTx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't cast to signer set for %d", req.SignerSetNonce)
	}

	return &types.SignerSetTxResponse{SignerSet: ss}, nil
}

func (k Keeper) BatchTx(c context.Context, req *types.BatchTxRequest) (*types.BatchTxResponse, error) {
	if !common.IsHexAddress(req.TokenContract) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid hex address %s", req.TokenContract)
	}

	res := &types.BatchTxResponse{}

	key := types.MakeBatchTxKey(common.HexToAddress(req.TokenContract), req.BatchNonce)
	otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), key)
	if otx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no batch tx found for %d %s", req.BatchNonce, req.TokenContract)
	}
	batch, ok := otx.(*types.BatchTx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't cast to batch tx for %d %s", req.BatchNonce, req.TokenContract)
	}
	res.Batch = batch

	return res, nil
}

func (k Keeper) ContractCallTx(c context.Context, req *types.ContractCallTxRequest) (*types.ContractCallTxResponse, error) {
	key := types.MakeContractCallTxKey(req.InvalidationScope, req.InvalidationNonce)
	otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), key)
	if otx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no contract call found for %d %s", req.InvalidationNonce, req.InvalidationScope)
	}

	cctx, ok := otx.(*types.ContractCallTx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't cast to contract call for %d %s", req.InvalidationNonce, req.InvalidationScope)
	}

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

func (k Keeper) EthereumTxConfirmations(c context.Context, req *types.EthereumTxConfirmationsRequest) (*types.EthereumTxConfirmationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var out []*types.EthereumTxConfirmation
	k.EthereumSignatureStore.Iterate(ctx, req.Checkpoint, func(ethereumSigner string, sig []byte) bool {
		out = append(out, &types.EthereumTxConfirmation{
			EthereumSigner:    ethereumSigner,
			EthereumSignature: sig,
		})
		return false
	})

	return &types.EthereumTxConfirmationsResponse{Confirmations: out}, nil
}

func (k Keeper) UnsignedSignerSetTxs(c context.Context, req *types.UnsignedSignerSetTxsRequest) (*types.UnsignedSignerSetTxsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	ethereumSigner := k.GetValidatorEthereumAddress(ctx, val).Hex()

	var sstxs []*types.SignerSetTx
	k.IterateOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.EthereumSignatureStore.Get(ctx, otx.GetCheckpoint([]byte(k.getGravityID(ctx))), ethereumSigner)
		if sig == nil { // it's pending
			sstx, ok := otx.(*types.SignerSetTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to signer set tx for %s", otx))
			}
			sstxs = append(sstxs, sstx)
		}
		return false
	})
	return &types.UnsignedSignerSetTxsResponse{SignerSets: sstxs}, nil
}

func (k Keeper) UnsignedBatchTxs(c context.Context, req *types.UnsignedBatchTxsRequest) (*types.UnsignedBatchTxsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	ethereumSigner := k.GetValidatorEthereumAddress(ctx, val).Hex()

	var batches []*types.BatchTx
	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.EthereumSignatureStore.Get(ctx, otx.GetCheckpoint([]byte(k.getGravityID(ctx))), ethereumSigner)
		if sig == nil { // it's pending
			batch, ok := otx.(*types.BatchTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to batch tx for %s", otx))
			}
			batches = append(batches, batch)
		}
		return false
	})
	return &types.UnsignedBatchTxsResponse{Batches: batches}, nil
}

func (k Keeper) UnsignedContractCallTxs(c context.Context, req *types.UnsignedContractCallTxsRequest) (*types.UnsignedContractCallTxsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	ethereumSigner := k.GetValidatorEthereumAddress(ctx, val).Hex()

	var cctxs []*types.ContractCallTx
	k.IterateOutgoingTxsByType(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.EthereumSignatureStore.Get(ctx, otx.GetCheckpoint([]byte(k.getGravityID(ctx))), ethereumSigner)
		if sig == nil { // it's pending
			cctx, ok := otx.(*types.ContractCallTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to contract call tx for %s", otx))
			}
			cctxs = append(cctxs, cctx)
		}
		return false
	})
	return &types.UnsignedContractCallTxsResponse{Calls: cctxs}, nil
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

// TODO(levi) this should be removed; it doesn't have a clear use case
func (k Keeper) UnbatchedSendToEthereums(c context.Context, req *types.UnbatchedSendToEthereumsRequest) (*types.UnbatchedSendToEthereumsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sendToEthereums, pageRes, err := k.SendToEthereumStore.PaginateBySender(ctx, req.SenderAddress, req.Pagination)
	res := &types.UnbatchedSendToEthereumsResponse{
		SendToEthereums: sendToEthereums,
		Pagination:      pageRes,
	}
	return res, err
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
