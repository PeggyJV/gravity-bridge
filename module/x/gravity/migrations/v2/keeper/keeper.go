package keeper

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v2/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	StakingKeeper          types.StakingKeeper
	storeKey               sdk.StoreKey
	paramSpace             paramtypes.Subspace
	Cdc                    codec.Codec
	accountKeeper          types.AccountKeeper
	bankKeeper             types.BankKeeper
	SlashingKeeper         types.SlashingKeeper
	DistributionKeeper     types.DistributionKeeper
	PowerReduction         sdk.Int
	hooks                  types.GravityHooks
	ReceiverModuleAccounts map[string]string
	SenderModuleAccounts   map[string]string
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(
	cdc codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accKeeper types.AccountKeeper,
	stakingKeeper types.StakingKeeper,
	bankKeeper types.BankKeeper,
	slashingKeeper types.SlashingKeeper,
	distributionKeeper types.DistributionKeeper,
	powerReduction sdk.Int,
	receiverModuleAccounts map[string]string,
	senderModuleAccounts map[string]string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	k := Keeper{
		Cdc:                    cdc,
		paramSpace:             paramSpace,
		storeKey:               storeKey,
		accountKeeper:          accKeeper,
		StakingKeeper:          stakingKeeper,
		bankKeeper:             bankKeeper,
		SlashingKeeper:         slashingKeeper,
		DistributionKeeper:     distributionKeeper,
		PowerReduction:         powerReduction,
		ReceiverModuleAccounts: receiverModuleAccounts,
		SenderModuleAccounts:   senderModuleAccounts,
	}

	return k
}

/////////////////////////////
//     SignerSetTxNonce    //
/////////////////////////////

// GetLatestSignerSetTxNonce returns the latest valset nonce
func (k Keeper) GetLatestSignerSetTxNonce(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.storeKey).Get([]byte{types.LatestSignerSetTxNonceKey}); bz != nil {
		return binary.BigEndian.Uint64(bz)
	}
	return 0
}

//////////////////////////////
// LastUnbondingBlockHeight //
//////////////////////////////

// GetLastUnbondingBlockHeight returns the last unbonding block height
func (k Keeper) GetLastUnbondingBlockHeight(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.storeKey).Get([]byte{types.LastUnBondingBlockHeightKey}); len(bz) == 0 {
		return 0
	} else {
		return binary.BigEndian.Uint64(bz)
	}
}

/////////////////////////////
//       PARAMETERS        //
/////////////////////////////

// GetParams returns the parameters from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

/////////////////
// OUTGOING TX //
/////////////////

// DeleteOutgoingTx deletes a given outgoingtx
func (k Keeper) DeleteOutgoingTx(ctx sdk.Context, storeIndex []byte) {
	ctx.KVStore(k.storeKey).Delete(types.MakeOutgoingTxKey(storeIndex))
}

// IterateOutgoingTxsByType iterates over a specific type of outgoing transaction denoted by the chosen prefix byte
func (k Keeper) IterateOutgoingTxsByType(ctx sdk.Context, prefixByte byte, cb func(key []byte, outgoing types.OutgoingTx) (stop bool)) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.MakeOutgoingTxKey([]byte{prefixByte}))
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var any cdctypes.Any
		k.Cdc.MustUnmarshal(iter.Value(), &any)
		var otx types.OutgoingTx
		if err := k.Cdc.UnpackAny(&any, &otx); err != nil {
			panic(err)
		}
		if cb(iter.Key(), otx) {
			break
		}
	}
}

// GetLastObservedSignerSetTx retrieves the last observed validator set from the store
func (k Keeper) GetLastObservedSignerSetTx(ctx sdk.Context) *types.SignerSetTx {
	key := []byte{types.LastObservedSignerSetKey}
	if val := ctx.KVStore(k.storeKey).Get(key); val != nil {
		var out types.SignerSetTx
		k.Cdc.MustUnmarshal(val, &out)
		return &out
	}
	return nil
}
