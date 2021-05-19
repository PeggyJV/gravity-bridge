package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	// TODO: audit once we finalize storage
	storeIndex := sdk.Uint64ToBigEndian(req.Nonce)
	otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), types.GetOutgoingTxKey(storeIndex))
	if otx == nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "no signer set found for %d", req.Nonce)
	}

	ss, ok := otx.(*types.SignerSetTx)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for %d", req.Nonce)
	}

	// TODO: special case nonce = 0 to find latest
	// TODO: ensure that latest signer set tx nonce index is set properly
	// TODO: ensure nonce sequence starts at one

	return &types.SignerSetTxResponse{SignerSet: ss}, nil
}

func (k Keeper) BatchTx(c context.Context, req *types.BatchTxRequest) (*types.BatchTxResponse, error) {
	if !common.IsHexAddress(req.ContractAddress) {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid hex address %s", req.ContractAddress)
	}

	// TODO: audit once we finalize storage
	storeIndex := append(sdk.Uint64ToBigEndian(req.Nonce), common.Hex2Bytes(req.ContractAddress)...)
	otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), types.GetOutgoingTxKey(storeIndex))
	if otx == nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "no batch tx found for %d %s", req.Nonce, req.ContractAddress)
	}

	batch, ok := otx.(*types.BatchTx)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to batch tx for %d %s", req.Nonce, req.ContractAddress)
	}

	// TODO: handle special case nonce = 0 to find latest by contract address

	return &types.BatchTxResponse{Batch: batch}, nil
}

func (k Keeper) ContractCallTx(c context.Context, req *types.ContractCallTxRequest) (*types.ContractCallTxResponse, error) {
	storeIndex := append(sdk.Uint64ToBigEndian(req.InvalidationNonce), req.InvalidationScope...)
	otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), types.GetOutgoingTxKey(storeIndex))
	if otx == nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "no contract call found for %d %s", req.InvalidationNonce, req.InvalidationScope)
	}

	cctx, ok := otx.(*types.ContractCallTx)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to contract call for %d %s", req.InvalidationNonce, req.InvalidationScope)
	}

	// TODO: figure out how to call latest

	return &types.ContractCallTxResponse{LogicCall: cctx}, nil
}

func (k Keeper) SignerSetTxs(c context.Context, req *types.SignerSetTxsRequest) (*types.SignerSetTxsResponse, error) {
	var signers []*types.SignerSetTx
	k.IterateOutgoingTxs(sdk.UnwrapSDKContext(c), types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		signer, ok := otx.(*types.SignerSetTx)
		if !ok {
			panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for %s", otx))
		}
		signers = append(signers, signer)

		return len(signers) < int(req.Count)
	})
	return &types.SignerSetTxsResponse{SignerSets: signers}, nil
}

func (k Keeper) BatchTxs(c context.Context, req *types.BatchTxsRequest) (*types.BatchTxsResponse, error) {
	var batches []*types.BatchTx
	k.IterateOutgoingTxs(sdk.UnwrapSDKContext(c), types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		batch, ok := otx.(*types.BatchTx)
		if !ok {
			panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to batch tx for %s", otx))
		}
		batches = append(batches, batch)
		return false
	})
	return &types.BatchTxsResponse{Batches: batches}, nil
}

func (k Keeper) ContractCallTxs(c context.Context, req *types.ContractCallTxsRequest) (*types.ContractCallTxsResponse, error) {
	var calls []*types.ContractCallTx
	k.IterateOutgoingTxs(sdk.UnwrapSDKContext(c), types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		call, ok := otx.(*types.ContractCallTx)
		if !ok {
			panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to contract call for %s", otx))
		}
		calls = append(calls, call)
		return false
	})
	return &types.ContractCallTxsResponse{Calls: calls}, nil
}

func (k Keeper) SignerSetTxEthereumSignatures(c context.Context, req *types.SignerSetTxEthereumSignaturesRequest) (*types.SignerSetTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeSignerSetTxKey(req.Nonce)
	if req.Address == "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.SignerSetTxEthereumSignaturesResponse{Signature: [][]byte{k.GetEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.IterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.SignerSetTxEthereumSignaturesResponse{Signature: out}, nil
}

func (k Keeper) BatchTxEthereumSignatures(c context.Context, req *types.BatchTxEthereumSignaturesRequest) (*types.BatchTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeBatchTxKey(common.HexToAddress(req.ContractAddress), req.Nonce)
	if req.Address == "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.BatchTxEthereumSignaturesResponse{Signature: [][]byte{k.GetEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.IterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.BatchTxEthereumSignaturesResponse{Signature: out}, nil
}

func (k Keeper) ContractCallTxEthereumSignatures(c context.Context, req *types.ContractCallTxEthereumSignaturesRequest) (*types.ContractCallTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeContractCallTxKey(req.InvalidationScope, req.InvalidationNonce)
	if req.Address == "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.ContractCallTxEthereumSignaturesResponse{Signature: [][]byte{k.GetEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.IterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
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
	k.IterateOutgoingTxs(ctx, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.GetEthereumSignature(ctx, otx.GetStoreIndex(), val)
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
	k.IterateOutgoingTxs(ctx, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.GetEthereumSignature(ctx, otx.GetStoreIndex(), val)
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
	k.IterateOutgoingTxs(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.GetEthereumSignature(ctx, otx.GetStoreIndex(), val)
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

	valAddr, err := sdk.ValAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	res.EventNonce = k.GetLastEventNonceByValidator(ctx, valAddr)

	return &types.LastSubmittedEthereumEventResponse{}, nil
}

func (k Keeper) BatchTxFees(c context.Context, req *types.BatchTxFeesRequest) (*types.BatchTxFeesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	res := &types.BatchTxFeesResponse{}

	k.IterateOutgoingTxs(ctx, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
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

	k.IterateOutgoingTxs(ctx, types.BatchTxPrefixByte, func(_ []byte, outgoing types.OutgoingTx) bool {
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

	k.IterateUnbatchedSendToEthereums(ctx, func(ste *types.SendToEthereum) bool {
		if ste.Sender == req.SenderAddress {
			res.SendToEthereums = append(res.SendToEthereums, ste)
		}
		return false
	})

	return res, nil
}

func (k Keeper) DelegateKeysByValidator(c context.Context, req *types.DelegateKeysByValidatorAddress) (*types.DelegateKeysByValidatorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	ethAddr := k.GetValidatorEthereumAddress(ctx, valAddr)
	orchAddr := k.GetEthereumOrchestratorAddress(ctx, ethAddr)
	res := &types.DelegateKeysByValidatorAddressResponse{
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

func (k Keeper) DelegateKeysByOrchestrator(c context.Context, req *types.DelegateKeysByOrchestratorAddress) (*types.DelegateKeysByOrchestratorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	orchAddr, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, err
	}
	valAddr := k.GetOrchestratorValidatorAddress(ctx, orchAddr)
	ethAddr := k.GetValidatorEthereumAddress(ctx, valAddr)
	res := &types.DelegateKeysByOrchestratorAddressResponse{
		ValidatorAddress: valAddr.String(),
		EthereumSigner:   ethAddr.Hex(),
	}
	return res, nil
}
