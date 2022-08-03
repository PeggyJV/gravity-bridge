package v2

import (
	"bytes"
	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	oldKeeper "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v2/keeper"
	oldTypes "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v2/types"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

type NewKeeper struct {
	StakingKeeper          types.StakingKeeper
	StoreKey               sdk.StoreKey
	ParamSpace             paramtypes.Subspace
	Cdc                    codec.Codec
	AccountKeeper          types.AccountKeeper
	BankKeeper             types.BankKeeper
	SlashingKeeper         types.SlashingKeeper
	DistributionKeeper     types.DistributionKeeper
	PowerReduction         sdk.Int
	ReceiverModuleAccounts map[string]string
	SenderModuleAccounts   map[string]string
}

func MigrateStore(ctx sdk.Context, newK *NewKeeper) error {
	ctx.Logger().Info("Gravity v2 to v3: Beginning store migration")

	oldK := oldKeeper.NewKeeper(
		newK.Cdc,
		newK.StoreKey,
		newK.ParamSpace,
		newK.AccountKeeper,
		newK.StakingKeeper,
		newK.BankKeeper,
		newK.SlashingKeeper,
		newK.DistributionKeeper,
		newK.PowerReduction,
		newK.ReceiverModuleAccounts,
		newK.SenderModuleAccounts,
	)

	// params
	migrateParams(ctx, newK, &oldK)

	// nonces
	migrateSignerSetTxNonce(ctx, newK, &oldK)
	migrateLastEventNonceByValidators(ctx, newK)
	migrateLastObservedEventNonce(ctx, newK, &oldK)
	migrateLastSlashedOutgoingTxBlockHeight(ctx, newK, &oldK)
	migrateLastOutgoingBatchNonce(ctx, newK)

	// evm signatures
	migrateEVMSignatures(ctx, newK)
	migrateSendToEVMID(ctx, newK)
	migrateEVMBlockHeight(ctx, newK, &oldK)
	migrateSendToEVMs(ctx, newK, &oldK)
	migrateEVMEventVoteRecords(ctx, newK, &oldK)
	migrateEVMEventHeightVotes(ctx, newK, &oldK)

	// delegate keys (not chain specific but need the new proto encoding)
	migrateDelegateKeys(ctx, newK)

	// outgoing txs
	migrateSignerSetTxs(ctx, newK, &oldK)
	migrateLastObservedSignerSet(ctx, newK, &oldK)
	migrateBatchTXs(ctx, newK, &oldK)
	migrateContractCallTxs(ctx, newK, &oldK)

	// denoms
	migrateERC20ToDenom(ctx, newK, &oldK)
	migrateDenomToERC20(ctx, newK)

	ctx.Logger().Info("Gravity v2 to v3: Store migration complete")

	return nil
}

func migrateLastObservedSignerSet(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	sstx := oldK.GetLastObservedSignerSetTx(ctx)
	store := ctx.KVStore(newK.StoreKey)
	store.Delete([]byte{oldTypes.LastObservedSignerSetKey})

	var signers types.EVMSigners
	for _, s := range sstx.Signers {
		signers = append(signers, &types.EVMSigner{
			Power:      s.Power,
			EVMAddress: s.EthereumAddress,
		})
	}

	newSstx := types.SignerSetTx{
		Nonce:   sstx.Nonce,
		Height:  sstx.Height,
		Signers: signers,
		ChainId: types.EthereumChainID,
	}

	store.Set(types.MakeLastObservedSignerSetKey(types.EthereumChainID), newK.Cdc.MustMarshal(&newSstx))
}

func migrateSignerSetTxNonce(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	nonce := oldK.GetLatestSignerSetTxNonce(ctx)
	store := ctx.KVStore(newK.StoreKey)
	store.Delete([]byte{oldTypes.LatestSignerSetTxNonceKey})

	store.Set(types.MakeLatestSignerSetTxNonceKey(types.EthereumChainID), sdk.Uint64ToBigEndian(nonce))
}

func migrateLastObservedEventNonce(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	nonce := oldK.GetLastObservedEventNonce(ctx)
	store := ctx.KVStore(newK.StoreKey)
	store.Delete([]byte{oldTypes.LastObservedEventNonceKey})
	store.Set(types.MakeLastObservedEventNonceKey(types.EthereumChainID), sdk.Uint64ToBigEndian(nonce))
}

func (k NewKeeper) setLastSlashedOutgoingTxBlockHeight(ctx sdk.Context, chainID uint32, blockHeight uint64) {
	ctx.KVStore(k.StoreKey).Set(types.MakeLastSlashedOutgoingTxBlockKey(chainID), sdk.Uint64ToBigEndian(blockHeight))
}

func migrateLastSlashedOutgoingTxBlockHeight(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	height := oldK.GetLastSlashedOutgoingTxBlockHeight(ctx)
	store := ctx.KVStore(newK.StoreKey)
	store.Delete([]byte{oldTypes.LastSlashedOutgoingTxBlockKey})
	newK.setLastSlashedOutgoingTxBlockHeight(ctx, types.EthereumChainID, height)
}

func migrateLastOutgoingBatchNonce(ctx sdk.Context, newK *NewKeeper) {
	store := ctx.KVStore(newK.StoreKey)
	nonce := store.Get([]byte{oldTypes.LastOutgoingBatchNonceKey})
	store.Delete([]byte{oldTypes.LastOutgoingBatchNonceKey})
	store.Set(types.MakeLastOutgoingBatchNonceKey(types.EthereumChainID), nonce)
}

func migrateLastEventNonceByValidators(ctx sdk.Context, newK *NewKeeper) {
	store := ctx.KVStore(newK.StoreKey)
	prefixStore := prefix.NewStore(store, []byte{oldTypes.LastEventNonceByValidatorKey})
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	var nonceByValidatorKeys [][]byte
	var nonceByValidatorValues [][]byte

	for ; iter.Valid(); iter.Next() {
		nonceByValidatorKeys = append(nonceByValidatorKeys, iter.Key())
		nonceByValidatorValues = append(nonceByValidatorValues, iter.Value())
	}

	for i, key := range nonceByValidatorKeys {
		store.Delete(key)
		newKey := bytes.Join([][]byte{{types.LastEventNonceByValidatorKey}, types.Uint32ToBigEndian(types.EthereumChainID), key[1:]}, []byte{})
		store.Set(newKey, nonceByValidatorValues[i])
	}
}

func migrateEVMSignatures(ctx sdk.Context, newK *NewKeeper) {
	store := ctx.KVStore(newK.StoreKey)
	prefixStore := prefix.NewStore(store, []byte{oldTypes.EthereumSignatureKey})
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	var evmSignatureKeys [][]byte
	var evmSignatureValues [][]byte

	for ; iter.Valid(); iter.Next() {
		evmSignatureKeys = append(evmSignatureKeys, iter.Key())
		evmSignatureValues = append(evmSignatureKeys, iter.Value())
	}

	for i, key := range evmSignatureKeys {
		store.Delete(key)
		newKey := bytes.Join([][]byte{types.EVMSignatureKeyPrefix(), key[:1]}, []byte{})
		store.Set(newKey, evmSignatureValues[i])
	}
}

func migrateSendToEVMID(ctx sdk.Context, newK *NewKeeper) {
	store := ctx.KVStore(newK.StoreKey)
	id := store.Get([]byte{oldTypes.LastSendToEthereumIDKey})
	store.Delete([]byte{oldTypes.LastSendToEthereumIDKey})
	store.Set(types.MakeLastSendToEVMIDKey(types.EthereumChainID), id)
}

func migrateEVMBlockHeight(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	oldTypeHeight := oldK.GetLastObservedEthereumBlockHeight(ctx)
	store := ctx.KVStore(newK.StoreKey)
	store.Delete([]byte{oldTypes.LastEthereumBlockHeightKey})

	newTypeHeight := types.LatestEVMBlockHeight{
		EVMHeight:    oldTypeHeight.EthereumHeight,
		CosmosHeight: oldTypeHeight.CosmosHeight,
		ChainId:      types.EthereumChainID,
	}

	store.Set(types.MakeLastEVMBlockHeightKey(types.EthereumChainID), newK.Cdc.MustMarshal(&newTypeHeight))
}

func migrateParams(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	oldParams := oldK.GetParams(ctx)

	newParams := types.Params{
		AverageBlockTime: oldParams.AverageBlockTime,
		ParamsByChain: map[string]*types.ParamsForChain{
			strconv.Itoa(types.EthereumChainID): {
				GravityId:                            oldParams.GravityId,
				SignedSignerSetTxsWindow:             oldParams.SignedSignerSetTxsWindow,
				SignedBatchesWindow:                  oldParams.SignedBatchesWindow,
				EvmSignaturesWindow:                  oldParams.EthereumSignaturesWindow,
				TargetEvmTxTimeout:                   oldParams.TargetEthTxTimeout,
				AverageEvmBlockTime:                  oldParams.AverageEthereumBlockTime,
				SlashFractionSignerSetTx:             oldParams.SlashFractionSignerSetTx,
				SlashFractionBatch:                   oldParams.SlashFractionBatch,
				SlashFractionEvmSignature:            oldParams.SlashFractionEthereumSignature,
				SlashFractionConflictingEvmSignature: oldParams.SlashFractionConflictingEthereumSignature,
				UnbondSlashingSignerSetTxsWindow:     oldParams.UnbondSlashingSignerSetTxsWindow,
			},
		},
	}

	// todo: is there a way to delete the current paramset? will this clash?
	newK.ParamSpace.SetParamSet(ctx, &newParams)
}

func migrateDelegateKeys(ctx sdk.Context, newK *NewKeeper) {
	store := ctx.KVStore(newK.StoreKey)
	iter := prefix.NewStore(store, []byte{oldTypes.ValidatorEthereumAddressKey}).Iterator(nil, nil)

	for ; iter.Valid(); iter.Next() {
		ethAddr := common.BytesToAddress(iter.Value())
		key := types.MakeEVMOrchestratorAddressKey(ethAddr)
		msg := types.MsgDelegateKeys{
			ValidatorAddress:    sdk.ValAddress(iter.Key()).String(),
			EVMAddress:          ethAddr.Hex(),
			OrchestratorAddress: sdk.AccAddress(store.Get(key)).String(),
		}

		store.Set(iter.Key(), newK.Cdc.MustMarshal(&msg))
	}
	iter.Close()
}

func (k NewKeeper) setOutgoingTx(ctx sdk.Context, chainID uint32, outgoing types.OutgoingTx) {
	outgoingTx, err := types.PackOutgoingTx(outgoing)
	if err != nil {
		panic(err)
	}
	ctx.KVStore(k.StoreKey).Set(
		types.MakeOutgoingTxKey(outgoing.GetStoreIndex()),
		k.Cdc.MustMarshal(outgoingTx),
	)
}

func migrateSignerSetTxs(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	var oldSignerSetTxs []*oldTypes.SignerSetTx

	oldLastObserved := oldK.GetLastObservedSignerSetTx(ctx)

	oldK.IterateOutgoingTxsByType(ctx, oldTypes.SignerSetTxPrefixByte, func(key []byte, outgoing oldTypes.OutgoingTx) (stop bool) {
		oldSignerSetTxs = append(oldSignerSetTxs, outgoing.(*oldTypes.SignerSetTx))
		return false
	})

	for _, otx := range oldSignerSetTxs {
		oldK.DeleteOutgoingTx(ctx, otx.GetStoreIndex())

		var evmSigners types.EVMSigners

		for _, signer := range otx.Signers {
			evmSigners = append(evmSigners, &types.EVMSigner{
				Power:      signer.Power,
				EVMAddress: signer.EthereumAddress,
			})
		}

		newOtx := types.SignerSetTx{
			Nonce:   otx.Nonce,
			Height:  otx.Height,
			Signers: evmSigners,
			ChainId: types.EthereumChainID,
		}

		newK.setOutgoingTx(ctx, types.EthereumChainID, &newOtx)

		if otx == oldLastObserved {
			ctx.KVStore(newK.StoreKey).Delete([]byte{oldTypes.LastObservedSignerSetKey})
			key := types.MakeLastObservedSignerSetKey(types.EthereumChainID)
			ctx.KVStore(newK.StoreKey).Set(key, newK.Cdc.MustMarshal(&newOtx))
		}
	}

}

func migrateBatchTXs(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	var oldBatchTxs []*oldTypes.BatchTx

	oldK.IterateOutgoingTxsByType(ctx, oldTypes.BatchTxPrefixByte, func(key []byte, outgoing oldTypes.OutgoingTx) (stop bool) {
		oldBatchTxs = append(oldBatchTxs, outgoing.(*oldTypes.BatchTx))
		return false
	})

	for _, otx := range oldBatchTxs {
		oldK.DeleteOutgoingTx(ctx, otx.GetStoreIndex())
	}

	for _, otx := range oldBatchTxs {
		var transactions []*types.SendToEVM

		for _, t := range otx.Transactions {
			transactions = append(transactions, &types.SendToEVM{
				Id:           t.Id,
				Sender:       t.Sender,
				EVMRecipient: t.EthereumRecipient,
				Erc20Token: types.ERC20Token{
					Contract: t.Erc20Token.Contract,
					Amount:   t.Erc20Token.Amount,
				},
				Erc20Fee: types.ERC20Token{
					Contract: t.Erc20Fee.Contract,
					Amount:   t.Erc20Fee.Amount,
				},
				ChainId: types.EthereumChainID,
			})
		}

		newOtx := types.BatchTx{
			BatchNonce:    otx.BatchNonce,
			Timeout:       otx.Timeout,
			Transactions:  transactions,
			TokenContract: otx.TokenContract,
			Height:        otx.Height,
			ChainId:       types.EthereumChainID,
		}

		newK.setOutgoingTx(ctx, types.EthereumChainID, &newOtx)
	}
}

func migrateContractCallTxs(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	var oldContractCallTxs []*oldTypes.ContractCallTx

	oldK.IterateOutgoingTxsByType(ctx, oldTypes.ContractCallTxPrefixByte, func(key []byte, outgoing oldTypes.OutgoingTx) (stop bool) {
		oldContractCallTxs = append(oldContractCallTxs, outgoing.(*oldTypes.ContractCallTx))
		return false
	})

	for _, otx := range oldContractCallTxs {
		oldK.DeleteOutgoingTx(ctx, otx.GetStoreIndex())
	}

	for _, otx := range oldContractCallTxs {
		var tokens []types.ERC20Token
		var fees []types.ERC20Token

		for _, t := range otx.Tokens {
			tokens = append(tokens, types.ERC20Token{
				Contract: t.Contract,
				Amount:   t.Amount,
			})
		}

		for _, f := range otx.Fees {
			fees = append(fees, types.ERC20Token{
				Contract: f.Contract,
				Amount:   f.Amount,
			})
		}

		newOtx := types.ContractCallTx{
			InvalidationNonce: otx.InvalidationNonce,
			InvalidationScope: otx.InvalidationScope,
			Address:           otx.Address,
			Payload:           otx.Payload,
			Timeout:           otx.Timeout,
			Tokens:            tokens,
			Fees:              fees,
			Height:            otx.Height,
			ChainId:           types.EthereumChainID,
		}

		newK.setOutgoingTx(ctx, types.EthereumChainID, &newOtx)
	}
}

func migrateSendToEVMs(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	store := ctx.KVStore(newK.StoreKey)
	iter := prefix.NewStore(store, []byte{oldTypes.SendToEthereumKey}).Iterator(nil, nil)
	defer iter.Close()

	var oldSendKeys [][]byte
	var oldSends []*oldTypes.SendToEthereum

	for ; iter.Valid(); iter.Next() {
		var oldSend oldTypes.SendToEthereum
		oldK.Cdc.MustUnmarshal(iter.Value(), &oldSend)
		oldSendKeys = append(oldSendKeys, iter.Key())
		oldSends = append(oldSends, &oldSend)
	}

	for i, key := range oldSendKeys {
		store.Delete(key)
		newKey := bytes.Join([][]byte{{types.SendToEVMKey}, types.Uint32ToBigEndian(types.EthereumChainID), key[:1]}, []byte{})

		oldSend := oldSends[i]
		newSend := types.SendToEVM{
			Id:           oldSend.Id,
			Sender:       oldSend.Sender,
			EVMRecipient: oldSend.EthereumRecipient,
			Erc20Token: types.ERC20Token{
				Contract: oldSend.Erc20Token.Contract,
				Amount:   oldSend.Erc20Token.Amount,
			},
			Erc20Fee: types.ERC20Token{
				Contract: oldSend.Erc20Fee.Contract,
				Amount:   oldSend.Erc20Fee.Amount,
			},
			ChainId: types.EthereumChainID,
		}

		store.Set(newKey, newK.Cdc.MustMarshal(&newSend))
	}
}

func migrateEVMEventVoteRecords(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	store := ctx.KVStore(newK.StoreKey)
	iter := prefix.NewStore(store, []byte{oldTypes.EthereumEventVoteRecordKey}).Iterator(nil, nil)
	defer iter.Close()

	var oldKeys [][]byte
	var oldRecords []*oldTypes.EthereumEventVoteRecord

	for ; iter.Valid(); iter.Next() {
		var oldSend oldTypes.EthereumEventVoteRecord
		oldK.Cdc.MustUnmarshal(iter.Value(), &oldSend)
		oldKeys = append(oldKeys, iter.Key())
		oldRecords = append(oldRecords, &oldSend)
	}

	for i, key := range oldKeys {
		store.Delete(key)
		newKey := bytes.Join([][]byte{{types.EVMEventVoteRecordKey}, types.Uint32ToBigEndian(types.EthereumChainID), key[:1]}, []byte{})

		oldRecord := oldRecords[i]
		newRecord := types.EVMEventVoteRecord{
			Event:    oldRecord.Event, // TODO: does this any type need to be re-serialized?
			Votes:    oldRecord.Votes,
			Accepted: oldRecord.Accepted,
			ChainId:  types.EthereumChainID,
		}

		store.Set(newKey, newK.Cdc.MustMarshal(&newRecord))
	}
}

func migrateEVMEventHeightVotes(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	store := ctx.KVStore(newK.StoreKey)
	iter := prefix.NewStore(store, []byte{oldTypes.EthereumHeightVoteKey}).Iterator(nil, nil)
	defer iter.Close()

	var oldKeys [][]byte
	var oldVotes []*oldTypes.LatestEthereumBlockHeight

	for ; iter.Valid(); iter.Next() {
		var oldVote oldTypes.LatestEthereumBlockHeight
		oldK.Cdc.MustUnmarshal(iter.Value(), &oldVote)
		oldKeys = append(oldKeys, iter.Key())
		oldVotes = append(oldVotes, &oldVote)
	}

	for i, key := range oldKeys {
		store.Delete(key)
		newKey := bytes.Join([][]byte{{types.EVMHeightVoteKey}, types.Uint32ToBigEndian(types.EthereumChainID), key[:1]}, []byte{})

		oldVote := oldVotes[i]
		newRecord := types.LatestEVMBlockHeight{
			EVMHeight:    oldVote.EthereumHeight,
			CosmosHeight: oldVote.CosmosHeight,
			ChainId:      types.EthereumChainID,
		}

		store.Set(newKey, newK.Cdc.MustMarshal(&newRecord))
	}
}

func migrateERC20ToDenom(ctx sdk.Context, newK *NewKeeper, oldK *oldKeeper.Keeper) {
	store := ctx.KVStore(newK.StoreKey)
	iter := prefix.NewStore(store, []byte{oldTypes.ERC20ToDenomKey}).Iterator(nil, nil)
	defer iter.Close()

	var oldKeys [][]byte
	var oldE2Ds []*oldTypes.ERC20ToDenom

	for ; iter.Valid(); iter.Next() {
		var oldE2D oldTypes.ERC20ToDenom
		oldK.Cdc.MustUnmarshal(iter.Value(), &oldE2D)
		oldKeys = append(oldKeys, iter.Key())
		oldE2Ds = append(oldE2Ds, &oldE2D)
	}

	for i, key := range oldKeys {
		store.Delete(key)
		newKey := bytes.Join([][]byte{{types.ERC20ToDenomKey}, types.Uint32ToBigEndian(types.EthereumChainID), key[:1]}, []byte{})

		oldE2D := oldE2Ds[i]
		newE2D := types.ERC20ToDenom{
			Erc20: oldE2D.Erc20,
			Denom: oldE2D.Denom,
		}
		store.Set(newKey, newK.Cdc.MustMarshal(&newE2D))
	}
}

func migrateDenomToERC20(ctx sdk.Context, newK *NewKeeper) {
	store := ctx.KVStore(newK.StoreKey)
	iter := prefix.NewStore(store, []byte{oldTypes.DenomToERC20Key}).Iterator(nil, nil)
	defer iter.Close()

	var oldKeys [][]byte
	var oldD2Es [][]byte

	for ; iter.Valid(); iter.Next() {
		oldKeys = append(oldKeys, iter.Key())
		oldD2Es = append(oldD2Es, iter.Value())
	}

	for i, key := range oldKeys {
		store.Delete(key)
		newKey := bytes.Join([][]byte{{types.DenomToERC20Key}, types.Uint32ToBigEndian(types.EthereumChainID), key[:1]}, []byte{})

		store.Set(newKey, oldD2Es[i])
	}
}
