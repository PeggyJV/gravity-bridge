package keeper

import (
	"context"

	"cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(c))
	return &types.ParamsResponse{Params: params}, nil
}

func (k Keeper) LatestSignerSetTx(c context.Context, req *types.LatestSignerSetTxRequest) (*types.SignerSetTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), append([]byte{types.OutgoingTxKey}, types.SignerSetTxPrefixByte))
	iter := store.ReverseIterator(nil, nil)
	defer iter.Close()

	if !iter.Valid() {
		return nil, status.Errorf(codes.NotFound, "latest signer set not found")
	}

	var any cdctypes.Any
	k.cdc.MustUnmarshal(iter.Value(), &any)

	var otx types.OutgoingTx
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
			panic(errors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for %s", otx))
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
			panic(errors.Wrapf(types.ErrInvalid, "couldn't cast to batch tx for %s", otx))
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
			panic(errors.Wrapf(types.ErrInvalid, "couldn't cast to contract call for %s", otx))
		}
		calls = append(calls, call)
		return true
	})
	if err != nil {
		return nil, err
	}

	return &types.ContractCallTxsResponse{Calls: calls, Pagination: pageRes}, nil
}

func (k Keeper) SignerSetTxConfirmations(c context.Context, req *types.SignerSetTxConfirmationsRequest) (*types.SignerSetTxConfirmationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeSignerSetTxKey(req.SignerSetNonce)

	var out []*types.SignerSetTxConfirmation
	k.iterateEthereumSignatures(ctx, key, func(val sdk.ValAddress, sig []byte) bool {
		out = append(out, &types.SignerSetTxConfirmation{
			SignerSetNonce: req.SignerSetNonce,
			EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(),
			Signature:      sig,
		})
		return false
	})

	return &types.SignerSetTxConfirmationsResponse{Signatures: out}, nil
}

func (k Keeper) BatchTxConfirmations(c context.Context, req *types.BatchTxConfirmationsRequest) (*types.BatchTxConfirmationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeBatchTxKey(common.HexToAddress(req.TokenContract), req.BatchNonce)

	var out []*types.BatchTxConfirmation
	k.iterateEthereumSignatures(ctx, key, func(val sdk.ValAddress, sig []byte) bool {
		out = append(out, &types.BatchTxConfirmation{
			TokenContract:  req.TokenContract,
			BatchNonce:     req.BatchNonce,
			EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(),
			Signature:      sig,
		})
		return false
	})
	return &types.BatchTxConfirmationsResponse{Signatures: out}, nil
}

func (k Keeper) ContractCallTxConfirmations(c context.Context, req *types.ContractCallTxConfirmationsRequest) (*types.ContractCallTxConfirmationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeContractCallTxKey(req.InvalidationScope, req.InvalidationNonce)

	var out []*types.ContractCallTxConfirmation
	k.iterateEthereumSignatures(ctx, key, func(val sdk.ValAddress, sig []byte) bool {
		out = append(out, &types.ContractCallTxConfirmation{
			InvalidationScope: req.InvalidationScope,
			InvalidationNonce: req.InvalidationNonce,
			EthereumSigner:    k.GetValidatorEthereumAddress(ctx, val).Hex(),
			Signature:         sig,
		})
		return false
	})
	return &types.ContractCallTxConfirmationsResponse{Signatures: out}, nil
}

// UnsignedSignerSetTxs returns all signer set txs that have not been signed by the given validator
func (k Keeper) UnsignedSignerSetTxs(c context.Context, req *types.UnsignedSignerSetTxsRequest) (*types.UnsignedSignerSetTxsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	return &types.UnsignedSignerSetTxsResponse{SignerSets: k.GetUnsignedSignerSetTxs(ctx, val)}, nil
}

// UnsignedBatchTxs returns all batch txs that have not been signed by the given validator
func (k Keeper) UnsignedBatchTxs(c context.Context, req *types.UnsignedBatchTxsRequest) (*types.UnsignedBatchTxsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	return &types.UnsignedBatchTxsResponse{Batches: k.GetUnsignedBatchTxs(ctx, val)}, nil
}

// UnsignedContractCallTxs returns all contract call txs that have not been signed by the given validator
func (k Keeper) UnsignedContractCallTxs(c context.Context, req *types.UnsignedContractCallTxsRequest) (*types.UnsignedContractCallTxsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	return &types.UnsignedContractCallTxsResponse{Calls: k.GetUnsignedContractCallTxs(ctx, val)}, nil
}

func (k Keeper) LastSubmittedEthereumEvent(c context.Context, req *types.LastSubmittedEthereumEventRequest) (*types.LastSubmittedEthereumEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	valAddr, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	res := &types.LastSubmittedEthereumEventResponse{
		EventNonce: k.getLastEventNonceByValidator(ctx, valAddr),
	}
	return res, nil
}

func (k Keeper) BatchTxFees(c context.Context, req *types.BatchTxFeesRequest) (*types.BatchTxFeesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	res := &types.BatchTxFeesResponse{}

	// TODO: is this what we want here?
	// Should this calculation return a
	// map[contract_address]fees or something similar?
	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		btx, _ := otx.(*types.BatchTx)
		for _, tx := range btx.Transactions {
			_, denom := k.ERC20ToDenomLookup(ctx, common.HexToAddress(tx.Erc20Fee.Contract))
			res.Fees = append(res.Fees, sdk.NewCoin(denom, tx.Erc20Fee.Amount))
		}
		return false
	})

	return res, nil
}

func (k Keeper) ERC20ToDenom(c context.Context, req *types.ERC20ToDenomRequest) (*types.ERC20ToDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	cosmosOriginated, denom := k.ERC20ToDenomLookup(ctx, common.HexToAddress(req.Erc20))
	res := &types.ERC20ToDenomResponse{
		Denom:            denom,
		CosmosOriginated: cosmosOriginated,
	}
	return res, nil
}

func (k Keeper) DenomToERC20Params(c context.Context, req *types.DenomToERC20ParamsRequest) (*types.DenomToERC20ParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if existingERC20, exists := k.getCosmosOriginatedERC20(ctx, req.Denom); exists {
		return nil, errors.Wrapf(
			types.ErrInvalidERC20Event,
			"ERC20 token %s already exists for denom %s", existingERC20.Hex(), req.Denom,
		)
	}

	// use metadata, if we can find it
	if md, ok := k.bankKeeper.GetDenomMetaData(ctx, req.Denom); ok && md.Base != "" {
		var erc20Decimals uint64
		for _, denomUnit := range md.DenomUnits {
			if denomUnit.Denom == md.Display {
				erc20Decimals = uint64(denomUnit.Exponent)
				break
			}
		}

		return &types.DenomToERC20ParamsResponse{
			BaseDenom:     md.Base,
			Erc20Name:     md.Display,
			Erc20Symbol:   md.Display,
			Erc20Decimals: erc20Decimals,
		}, nil
	}

	if supply := k.bankKeeper.GetSupply(ctx, req.Denom); supply.IsZero() {
		return nil, errors.Wrapf(
			types.ErrInvalidERC20Event,
			"no supply exists for token %s without metadata", req.Denom,
		)
	}

	// no metadata, go with a zero decimal, no symbol erc-20
	res := &types.DenomToERC20ParamsResponse{
		BaseDenom:     req.Denom,
		Erc20Name:     req.Denom,
		Erc20Symbol:   "",
		Erc20Decimals: 0,
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
		k.cdc.MustUnmarshal(value, &ste)
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
		return nil, errors.Wrapf(types.ErrInvalid, "ethereum signer needs to be a hex address")
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

func (k Keeper) DelegateKeys(c context.Context, req *types.DelegateKeysRequest) (*types.DelegateKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	delegateKeys := k.getDelegateKeys(ctx)

	res := &types.DelegateKeysResponse{
		DelegateKeys: delegateKeys,
	}
	return res, nil
}

func (k Keeper) LastObservedEthereumHeight(c context.Context, req *types.LastObservedEthereumHeightRequest) (*types.LastObservedEthereumHeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	lastObservedEthereumHeight := k.GetLastObservedEthereumBlockHeight(ctx)

	res := &types.LastObservedEthereumHeightResponse{
		LastObservedEthereumHeight: &lastObservedEthereumHeight,
	}

	return res, nil
}

func (k Keeper) CompletedBatchTxs(c context.Context, req *types.CompletedBatchTxsRequest) (*types.CompletedBatchTxsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var batches []*types.BatchTx
	k.IterateCompletedOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		batchTx := otx.(*types.BatchTx)
		batches = append(batches, batchTx)
		return false
	})

	res := &types.CompletedBatchTxsResponse{
		CompletedBatchTxs: batches,
	}
	return res, nil
}

func (k Keeper) CompletedContractCallTxs(c context.Context, req *types.CompletedContractCallTxsRequest) (*types.CompletedContractCallTxsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var contractCalls []*types.ContractCallTx
	k.IterateCompletedOutgoingTxsByType(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		contractCallTx := otx.(*types.ContractCallTx)
		contractCalls = append(contractCalls, contractCallTx)
		return false
	})

	res := &types.CompletedContractCallTxsResponse{
		CompletedContractCallTxs: contractCalls,
	}
	return res, nil
}

func (k Keeper) CompletedSignerSetTxs(c context.Context, req *types.CompletedSignerSetTxsRequest) (*types.CompletedSignerSetTxsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var signerSetCalls []*types.SignerSetTx
	k.IterateCompletedOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		signerSetTx := otx.(*types.SignerSetTx)
		signerSetCalls = append(signerSetCalls, signerSetTx)
		return false
	})

	res := &types.CompletedSignerSetTxsResponse{
		CompletedSignerSetTxs: signerSetCalls,
	}
	return res, nil
}

func (k Keeper) BatchTxConfirmationsByValidator(c context.Context, req *types.BatchTxConfirmationsByValidatorRequest) (*types.BatchTxConfirmationsByValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	var confirms []*types.BatchTxConfirmation
	k.IterateBatchTxEthereumSignatures(ctx, func(contractAddress common.Address, nonce uint64, val sdk.ValAddress, sig []byte) bool {
		if !val.Equals(valAddr) {
			return false
		}

		confirms = append(confirms, &types.BatchTxConfirmation{
			TokenContract:  contractAddress.Hex(),
			BatchNonce:     nonce,
			EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(),
			Signature:      sig,
		})
		return false
	})

	res := &types.BatchTxConfirmationsByValidatorResponse{
		BatchTxConfirmations: confirms,
	}
	return res, nil
}

func (k Keeper) ContractCallTxConfirmationsByValidator(c context.Context, req *types.ContractCallTxConfirmationsByValidatorRequest) (*types.ContractCallTxConfirmationsByValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	var confirms []*types.ContractCallTxConfirmation
	k.IterateContractCallTxEthereumSignatures(ctx, func(invalidationScope []byte, invalidationNonce uint64, val sdk.ValAddress, sig []byte) bool {
		if !val.Equals(valAddr) {
			return false
		}

		confirms = append(confirms, &types.ContractCallTxConfirmation{
			InvalidationScope: invalidationScope,
			InvalidationNonce: invalidationNonce,
			EthereumSigner:    k.GetValidatorEthereumAddress(ctx, val).Hex(),
			Signature:         sig,
		})
		return false
	})

	res := &types.ContractCallTxConfirmationsByValidatorResponse{
		ContractCallTxConfirmations: confirms,
	}
	return res, nil
}

func (k Keeper) SignerSetTxConfirmationsByValidator(c context.Context, req *types.SignerSetTxConfirmationsByValidatorRequest) (*types.SignerSetTxConfirmationsByValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	var confirms []*types.SignerSetTxConfirmation
	k.IterateSignerSetTxEthereumSignatures(ctx, func(nonce uint64, val sdk.ValAddress, sig []byte) bool {
		if !val.Equals(valAddr) {
			return false
		}

		confirms = append(confirms, &types.SignerSetTxConfirmation{
			SignerSetNonce: nonce,
			EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(),
			Signature:      sig,
		})
		return false
	})

	res := &types.SignerSetTxConfirmationsByValidatorResponse{
		SignerSetTxConfirmations: confirms,
	}
	return res, nil
}

func (k Keeper) EthereumEventVoteRecords(c context.Context, req *types.EthereumEventVoteRecordsRequest) (*types.EthereumEventVoteRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	res := &types.EthereumEventVoteRecordsResponse{}
	pageRes, err := k.PaginateEthereumEventVoteRecords(ctx, req.Pagination, func(key []byte, eventVoteRecord *types.EthereumEventVoteRecord) bool {
		res.Records = append(res.Records, eventVoteRecord)

		return false
	})
	if err != nil {
		return nil, err
	}

	res.Pagination = pageRes

	return res, nil
}

func (k Keeper) EthereumEventVotes(c context.Context, req *types.EthereumEventVotesRequest) (*types.EthereumEventVotesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	var events []*cdctypes.Any
	k.IterateEthereumEventVoteRecords(ctx, func(key []byte, eventVoteRecord *types.EthereumEventVoteRecord) bool {
		for _, voter := range eventVoteRecord.Votes {
			voterAddr, err := sdk.ValAddressFromBech32(voter)
			if err != nil {
				continue
			}
			if voterAddr.Equals(val) {
				events = append(events, eventVoteRecord.Event)
				return false
			}
		}

		return false
	})

	res := &types.EthereumEventVotesResponse{
		Events: events,
	}
	return res, nil
}
