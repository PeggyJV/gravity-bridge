package gravity

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/keeper"
	oldKeeper "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v2/keeper"
	oldTypes "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v2/types"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func MigrateStore(ctx sdk.Context, newK *keeper.Keeper) error {
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

	migrateSignerSetTxNonce(ctx, newK, &oldK)
	migrateEVMSignatures(ctx, newK, &oldK)
	migrateSignerSetTxs(ctx, newK, &oldK)
	migrateParams(ctx, newK, &oldK)

	// unbonding block height is not chain specific, so doesn't change
	// delegate addresses are not chain specific, and don't change

	ctx.Logger().Info("Gravity v2 to v3: Store migration complete")

	return nil
}

func migrateSignerSetTxNonce(ctx sdk.Context, newK *keeper.Keeper, oldK *oldKeeper.Keeper) {
	nonce := oldK.GetLatestSignerSetTxNonce(ctx)
	ctx.KVStore(newK.StoreKey).Set(types.MakeLatestSignerSetTxNonceKey(types.EthereumChainID), sdk.Uint64ToBigEndian(nonce))
}

func keyInsertChainID(key []byte, chainID uint32) []byte {
	return append(key[:1], append(types.EVMSignatureKeyPrefix(chainID), key[:1]...)...)
}

func migrateEVMSignatures(ctx sdk.Context, newK *keeper.Keeper, oldK *oldKeeper.Keeper) {
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
		newKey := keyInsertChainID(key, types.EthereumChainID)
		store.Set(newKey, evmSignatureValues[i])
	}
}

func migrateSignerSetTxs(ctx sdk.Context, newK *keeper.Keeper, oldK *oldKeeper.Keeper) {
	signerSetTxs := oldK.GetSignerSetTxs(ctx)

	store := ctx.KVStore(newK.StoreKey)

	for _, sstx := range signerSetTxs {
		oldKey := oldTypes.MakeOutgoingTxKey(sstx.GetStoreIndex())
		store.Delete(oldKey)
		newK.SetOutgoingTx(ctx, types.EthereumChainID, sstx)
	}
}

func migrateParams(ctx sdk.Context, newK *keeper.Keeper, oldK *oldKeeper.Keeper) {
	oldParams := oldK.GetParams(ctx)

	newParams := types.Params{
		ChainParams: map[uint32]*types.ChainParams{
			types.EthereumChainID: {
				GravityId:                            oldParams.GravityId,
				ContractSourceHash:                   oldParams.ContractSourceHash,
				SignedSignerSetTxsWindow:             oldParams.SignedSignerSetTxsWindow,
				SignedBatchesWindow:                  oldParams.SignedBatchesWindow,
				EvmSignaturesWindow:                  oldParams.EthereumSignaturesWindow,
				TargetEvmTxTimeout:                   oldParams.TargetEthTxTimeout,
				AverageBlockTime:                     oldParams.AverageBlockTime,
				AverageEvmBlockTime:                  oldParams.AverageEthereumBlockTime,
				SlashFractionSignerSetTx:             oldParams.SlashFractionSignerSetTx,
				SlashFractionBatch:                   oldParams.SlashFractionBatch,
				SlashFractionEvmSignature:            oldParams.SlashFractionEthereumSignature,
				SlashFractionConflictingEvmSignature: oldParams.SlashFractionConflictingEthereumSignature,
				UnbondSlashingSignerSetTxsWindow:     oldParams.UnbondSlashingSignerSetTxsWindow,
			},
		},
	}

	newK.SetParams(ctx, newParams)
}
