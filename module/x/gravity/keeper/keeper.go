package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	StakingKeeper          types.StakingKeeper
	StoreKey               sdk.StoreKey
	ParamSpace             paramtypes.Subspace
	Cdc                    codec.Codec
	AccountKeeper          types.AccountKeeper
	BankKeeper             types.BankKeeper
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
		ParamSpace:             paramSpace,
		StoreKey:               storeKey,
		AccountKeeper:          accKeeper,
		StakingKeeper:          stakingKeeper,
		BankKeeper:             bankKeeper,
		SlashingKeeper:         slashingKeeper,
		DistributionKeeper:     distributionKeeper,
		PowerReduction:         powerReduction,
		ReceiverModuleAccounts: receiverModuleAccounts,
		SenderModuleAccounts:   senderModuleAccounts,
	}

	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

/////////////////////////////
//     SignerSetTxNonce    //
/////////////////////////////

// incrementLatestSignerSetTxNonce sets the latest valset nonce
func (k Keeper) incrementLatestSignerSetTxNonce(ctx sdk.Context, chainID uint32) uint64 {
	current := k.GetLatestSignerSetTxNonce(ctx, chainID)
	next := current + 1
	ctx.KVStore(k.StoreKey).Set(types.MakeLatestSignerSetTxNonceKey(chainID), sdk.Uint64ToBigEndian(next))
	return next
}

// GetLatestSignerSetTxNonce returns the latest valset nonce
func (k Keeper) GetLatestSignerSetTxNonce(ctx sdk.Context, chainID uint32) uint64 {
	if bz := ctx.KVStore(k.StoreKey).Get(types.MakeLatestSignerSetTxNonceKey(chainID)); bz != nil {
		return binary.BigEndian.Uint64(bz)
	}
	return 0
}

//////////////////////////////
// LastUnbondingBlockHeight //
//////////////////////////////

// setLastUnbondingBlockHeight sets the last unbonding block height
func (k Keeper) setLastUnbondingBlockHeight(ctx sdk.Context, unbondingBlockHeight uint64) {
	ctx.KVStore(k.StoreKey).Set(types.MakeLastUnBondingBlockHeightKey(), sdk.Uint64ToBigEndian(unbondingBlockHeight))
}

// GetLastUnbondingBlockHeight returns the last unbonding block height
func (k Keeper) GetLastUnbondingBlockHeight(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.StoreKey).Get(types.MakeLastUnBondingBlockHeightKey()); len(bz) == 0 {
		return 0
	} else {
		return binary.BigEndian.Uint64(bz)
	}
}

///////////////////////////////
//     ETHEREUM SIGNATURES   //
///////////////////////////////

// getEVMSignature returns a valset confirmation by a nonce and validator address
func (k Keeper) getEVMSignature(ctx sdk.Context, chainID uint32, storeIndex []byte, validator sdk.ValAddress) []byte {
	return ctx.KVStore(k.StoreKey).Get(types.MakeEVMSignatureKeyForValidator(chainID, storeIndex, validator))
}

// SetEVMSignature sets a valset confirmation
func (k Keeper) SetEVMSignature(ctx sdk.Context, chainID uint32, sig types.EVMTxConfirmation, val sdk.ValAddress) []byte {
	key := types.MakeEVMSignatureKeyForValidator(chainID, sig.GetStoreIndex(chainID), val)
	ctx.KVStore(k.StoreKey).Set(key, sig.GetSignature())
	return key
}

// GetEVMSignatures returns all evm signatures for a given outgoing tx by store index
func (k Keeper) GetEVMSignatures(ctx sdk.Context, chainID uint32, storeIndex []byte) map[string][]byte {
	var signatures = make(map[string][]byte)
	k.iterateEVMSignaturesByStoreIndex(ctx, chainID, storeIndex, func(val sdk.ValAddress, h []byte) bool {
		signatures[val.String()] = h
		return false
	})
	return signatures
}

// iterateEVMSignaturesByStoreIndex iterates through all valset confirms by nonce in ASC order
func (k Keeper) iterateEVMSignaturesByStoreIndex(ctx sdk.Context, chainID uint32, storeIndex []byte, cb func(sdk.ValAddress, []byte) bool) {
	prefixKey := types.EVMSignatureKeyStoreIndexPrefix(chainID, storeIndex)
	prefixStore := prefix.NewStore(ctx.KVStore(k.StoreKey), prefixKey)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		// cb returns true to stop early
		if cb(iter.Key(), iter.Value()) {
			break
		}
	}
}

/////////////////////////
//  ORC -> VAL ADDRESS //
/////////////////////////

// SetOrchestratorValidatorAddress sets the Orchestrator key for a given validator.
func (k Keeper) SetOrchestratorValidatorAddress(ctx sdk.Context, val sdk.ValAddress, orchAddr sdk.AccAddress) {
	store := ctx.KVStore(k.StoreKey)
	key := types.MakeOrchestratorValidatorAddressKey(orchAddr)

	store.Set(key, val.Bytes())
}

// GetOrchestratorValidatorAddress returns the validator key associated with an
// orchestrator key.
func (k Keeper) GetOrchestratorValidatorAddress(ctx sdk.Context, orchAddr sdk.AccAddress) sdk.ValAddress {
	store := ctx.KVStore(k.StoreKey)
	key := types.MakeOrchestratorValidatorAddressKey(orchAddr)

	return store.Get(key)
}

////////////////////////
// VAL -> ETH ADDRESS //
////////////////////////

// setValidatorEVMAddress sets the ethereum address for a given validator
func (k Keeper) setValidatorEVMAddress(ctx sdk.Context, valAddr sdk.ValAddress, ethAddr common.Address) {
	store := ctx.KVStore(k.StoreKey)
	key := types.MakeValidatorEVMAddressKeyForValidator(valAddr)

	store.Set(key, ethAddr.Bytes())
}

// GetValidatorEVMAddress returns the eth address for a given gravity validator.
func (k Keeper) GetValidatorEVMAddress(ctx sdk.Context, valAddr sdk.ValAddress) common.Address {
	store := ctx.KVStore(k.StoreKey)
	key := types.MakeValidatorEVMAddressKeyForValidator(valAddr)

	return common.BytesToAddress(store.Get(key))
}

func (k Keeper) getValidatorsByEVMAddress(ctx sdk.Context, ethAddr common.Address) (vals []sdk.ValAddress) {
	iter := ctx.KVStore(k.StoreKey).Iterator(nil, nil)

	for ; iter.Valid(); iter.Next() {
		if common.BytesToAddress(iter.Value()) == ethAddr {
			valBs := bytes.TrimPrefix(iter.Key(), types.ValidatorToEVMAddressKeyPrefix())
			val := sdk.ValAddress(valBs)
			vals = append(vals, val)
		}
	}

	return
}

////////////////////////
// ETH -> ORC ADDRESS //
////////////////////////

// setEVMOrchestratorAddress sets the eth orch addr mapping
func (k Keeper) setEVMOrchestratorAddress(ctx sdk.Context, ethAddr common.Address, orch sdk.AccAddress) {
	store := ctx.KVStore(k.StoreKey)
	key := types.MakeEVMOrchestratorAddressKey(ethAddr)

	store.Set(key, orch.Bytes())
}

// GetEVMOrchestratorAddress gets the orch address for a given eth address
func (k Keeper) GetEVMOrchestratorAddress(ctx sdk.Context, ethAddr common.Address) sdk.AccAddress {
	store := ctx.KVStore(k.StoreKey)
	key := types.MakeEVMOrchestratorAddressKey(ethAddr)

	return store.Get(key)
}

func (k Keeper) getEVMAddressesByOrchestrator(ctx sdk.Context, orch sdk.AccAddress) (ethAddrs []common.Address) {
	iter := ctx.KVStore(k.StoreKey).Iterator(nil, nil)

	for ; iter.Valid(); iter.Next() {
		if sdk.AccAddress(iter.Value()).String() == orch.String() {
			ethBs := bytes.TrimPrefix(iter.Key(), types.EVMToOrchestratorAddressKeyPrefix())
			ethAddr := common.BytesToAddress(ethBs)
			ethAddrs = append(ethAddrs, ethAddr)
		}
	}

	return
}

/////////////////
// SignerSetTx //
/////////////////

// GetLatestSignerSetTx returns the latest validator set in state
func (k Keeper) GetLatestSignerSetTx(ctx sdk.Context, chainID uint32) *types.SignerSetTx {
	key := types.MakeSignerSetTxKey(chainID, k.GetLatestSignerSetTxNonce(ctx, chainID))
	otx := k.GetOutgoingTx(ctx, chainID, key)
	out, _ := otx.(*types.SignerSetTx)
	return out
}

// CreateSignerSetTx gets the current signer set from the staking keeper, increments the nonce,
// creates the signer set tx object, emits an event and sets the signer set in state
func (k Keeper) CreateSignerSetTx(ctx sdk.Context, chainID uint32) *types.SignerSetTx {
	nonce := k.incrementLatestSignerSetTxNonce(ctx, chainID)
	currSignerSet := k.CurrentSignerSet(ctx)
	newSignerSetTx := types.NewSignerSetTx(nonce, uint64(ctx.BlockHeight()), currSignerSet)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMultisigUpdateRequest,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyChainID, fmt.Sprint(chainID)),
			sdk.NewAttribute(types.AttributeKeySignerSetNonce, fmt.Sprint(nonce)),
		),
	)
	k.SetOutgoingTx(ctx, chainID, newSignerSetTx)
	k.Logger(ctx).Info(
		"SignerSetTx created",
		"nonce", newSignerSetTx.Nonce,
		"height", newSignerSetTx.Height,
		"signers", len(newSignerSetTx.Signers),
	)
	return newSignerSetTx
}

// CurrentSignerSet gets powers from the store and normalizes them
// into an integer percentage with a resolution of uint32 Max meaning
// a given validators 'gravity power' is computed as
// Cosmos power / total cosmos power = x / uint32 Max
// where x is the voting power on the gravity contract. This allows us
// to only use integer division which produces a known rounding error
// from truncation equal to the ratio of the validators
// Cosmos power / total cosmos power ratio, leaving us at uint32 Max - 1
// total voting power. This is an acceptable rounding error since floating
// point may cause consensus problems if different floating point unit
// implementations are involved.
func (k Keeper) CurrentSignerSet(ctx sdk.Context) types.EVMSigners {
	validators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	ethereumSigners := make([]*types.EVMSigner, 0)
	var totalPower uint64
	for _, validator := range validators {
		val := validator.GetOperator()

		p := uint64(k.StakingKeeper.GetLastValidatorPower(ctx, val))

		if ethAddr := k.GetValidatorEVMAddress(ctx, val); ethAddr.Hex() != "0x0000000000000000000000000000000000000000" {
			es := &types.EVMSigner{Power: p, EVMAddress: ethAddr.Hex()}
			ethereumSigners = append(ethereumSigners, es)
			totalPower += p
		}
	}
	// normalize power values
	for i := range ethereumSigners {
		ethereumSigners[i].Power = sdk.NewUint(ethereumSigners[i].Power).MulUint64(math.MaxUint32).QuoUint64(totalPower).Uint64()
	}

	return ethereumSigners
}

// GetSignerSetTxs returns all the signer set txs from the store
func (k Keeper) GetSignerSetTxs(ctx sdk.Context, chainID uint32) (out []*types.SignerSetTx) {
	k.IterateOutgoingTxsByType(ctx, chainID, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sstx, _ := otx.(*types.SignerSetTx)
		out = append(out, sstx)
		return false
	})
	return
}

/////////////////////////////
//       PARAMETERS        //
/////////////////////////////

// GetParams returns the parameters from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.ParamSpace.GetParamSet(ctx, &params)
	return
}

// SetParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, ps types.Params) {
	k.ParamSpace.SetParamSet(ctx, &ps)
}

func (k Keeper) chainIDsContains(ctx sdk.Context, chainID uint32) bool {
	params := k.GetParams(ctx)
	for id, _ := range params.ChainParams {
		if chainID == id {
			return true
		}
	}
	return false
}

// getGravityID returns the GravityID the GravityID is essentially a salt value
// for bridge signatures, provided each chain running Gravity has a unique ID
// it won't be possible to play back signatures from one bridge onto another
// even if they share a validator set.
//
// The lifecycle of the GravityID is that it is set in the Genesis file
// read from the live chain for the contract deployment, once a Gravity contract
// is deployed the GravityID CAN NOT BE CHANGED. Meaning that it can't just be the
// same as the chain id since the chain id may be changed many times with each
// successive chain in charge of the same bridge
func (k Keeper) getGravityID(ctx sdk.Context, chainID uint32) string {
	params := k.GetParams(ctx)
	return params.ChainParams[chainID].GravityId
}

// getDelegateKeys iterates both the EthAddress and Orchestrator address indexes to produce
// a vector of MsgDelegateKeys entries containing all the delegate keys for state
// export / import. This may seem at first glance to be excessively complicated, why not combine
// the EthAddress and Orchestrator address indexes and simply iterate one thing? The answer is that
// even though we set the Eth and Orchestrator address in the same place we use them differently we
// always go from Orchestrator address to Validator address and from validator address to EVM address
// we want to keep looking up the validator address for various reasons, so a direct Orchestrator to EVM
// address mapping will mean having to keep two of the same data around just to provide lookups.
//
// For the time being this will serve
func (k Keeper) getDelegateKeys(ctx sdk.Context) (out []*types.MsgDelegateKeys) {
	store := ctx.KVStore(k.StoreKey)
	iter := prefix.NewStore(store, types.ValidatorToEVMAddressKeyPrefix()).Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		out = append(out, &types.MsgDelegateKeys{
			ValidatorAddress: sdk.ValAddress(iter.Key()).String(),
			EVMAddress:       common.BytesToAddress(iter.Value()).Hex(),
		})
	}
	iter.Close()

	for _, msg := range out {
		msg.OrchestratorAddress = k.GetEVMOrchestratorAddress(ctx, common.HexToAddress(msg.EVMAddress)).String()
	}

	// we iterated over a map, so now we have to sort to ensure the
	// output here is deterministic, eth address chosen for no particular
	// reason
	sort.Slice(out[:], func(i, j int) bool {
		return out[i].EVMAddress < out[j].EVMAddress
	})

	return out
}

// GetUnbondingValidators returns UnbondingValidators.
// Adding here in gravity keeper as cdc is available inside endblocker.
func (k Keeper) GetUnbondingValidators(unbondingVals []byte) stakingtypes.ValAddresses {
	unbondingValidators := stakingtypes.ValAddresses{}
	k.Cdc.MustUnmarshal(unbondingVals, &unbondingValidators)
	return unbondingValidators
}

/////////////////
// OUTGOING TX //
/////////////////

// GetOutgoingTx todo: outgoingTx prefix byte
func (k Keeper) GetOutgoingTx(ctx sdk.Context, chainID uint32, storeIndex []byte) (out types.OutgoingTx) {
	if err := k.Cdc.UnmarshalInterface(ctx.KVStore(k.StoreKey).Get(types.MakeOutgoingTxKey(chainID, storeIndex)), &out); err != nil {
		panic(err)
	}
	return out
}

func (k Keeper) SetOutgoingTx(ctx sdk.Context, chainID uint32, outgoing types.OutgoingTx) {
	outgoingTx, err := types.PackOutgoingTx(outgoing)
	if err != nil {
		panic(err)
	}
	ctx.KVStore(k.StoreKey).Set(
		types.MakeOutgoingTxKey(chainID, outgoing.GetStoreIndex()),
		k.Cdc.MustMarshal(outgoingTx),
	)
}

// DeleteOutgoingTx deletes a given outgoingtx
func (k Keeper) DeleteOutgoingTx(ctx sdk.Context, chainID uint32, storeIndex []byte) {
	ctx.KVStore(k.StoreKey).Delete(types.MakeOutgoingTxKey(chainID, storeIndex))
}

func (k Keeper) PaginateOutgoingTxsByType(ctx sdk.Context, chainID uint32, pageReq *query.PageRequest, prefixByte byte, cb func(key []byte, outgoing types.OutgoingTx) bool) (*query.PageResponse, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.StoreKey), types.MakeOutgoingTxKey(chainID, []byte{prefixByte}))

	return query.FilteredPaginate(prefixStore, pageReq, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if !accumulate {
			return false, nil
		}

		var anyOTx cdctypes.Any
		k.Cdc.MustUnmarshal(value, &anyOTx)
		var otx types.OutgoingTx
		if err := k.Cdc.UnpackAny(&anyOTx, &otx); err != nil {
			panic(err)
		}
		if accumulate {
			return cb(key, otx), nil
		}

		return false, nil
	})
}

// IterateOutgoingTxsByType iterates over a specific type of outgoing transaction denoted by the chosen prefix byte
func (k Keeper) IterateOutgoingTxsByType(ctx sdk.Context, chainID uint32, prefixByte byte, cb func(key []byte, outgoing types.OutgoingTx) (stop bool)) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.StoreKey), types.MakeOutgoingTxKey(chainID, []byte{prefixByte}))
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var anyOTx cdctypes.Any
		k.Cdc.MustUnmarshal(iter.Value(), &anyOTx)
		var otx types.OutgoingTx
		if err := k.Cdc.UnpackAny(&anyOTx, &otx); err != nil {
			panic(err)
		}
		if cb(iter.Key(), otx) {
			break
		}
	}
}

// iterateOutgoingTxs iterates over a specific type of outgoing transaction denoted by the chosen prefix byte
func (k Keeper) iterateOutgoingTxs(ctx sdk.Context, chainID uint32, cb func(key []byte, outgoing types.OutgoingTx) bool) {
	prefixKey := types.OutgoingTxKeyPrefix(chainID)
	prefixStore := prefix.NewStore(ctx.KVStore(k.StoreKey), prefixKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var anyOTx cdctypes.Any
		k.Cdc.MustUnmarshal(iter.Value(), &anyOTx)
		var otx types.OutgoingTx
		if err := k.Cdc.UnpackAny(&anyOTx, &otx); err != nil {
			panic(err)
		}
		if cb(iter.Key(), otx) {
			break
		}
	}
}

// GetLastObservedSignerSetTx retrieves the last observed validator set from the store
func (k Keeper) GetLastObservedSignerSetTx(ctx sdk.Context, chainID uint32) *types.SignerSetTx {
	key := types.MakeLastObservedSignerSetKey(chainID)
	if val := ctx.KVStore(k.StoreKey).Get(key); val != nil {
		var out types.SignerSetTx
		k.Cdc.MustUnmarshal(val, &out)
		return &out
	}
	return nil
}

// SetLastObservedSignerSetTx updates the last observed validator set in the store
func (k Keeper) SetLastObservedSignerSetTx(ctx sdk.Context, chainID uint32, signerSet types.SignerSetTx) {
	key := types.MakeLastObservedSignerSetKey(chainID)
	ctx.KVStore(k.StoreKey).Set(key, k.Cdc.MustMarshal(&signerSet))
}

// CreateContractCallTx xxx
func (k Keeper) CreateContractCallTx(ctx sdk.Context, chainID uint32, invalidationNonce uint64, invalidationScope tmbytes.HexBytes,
	address common.Address, payload []byte, tokens []types.ERC20Token, fees []types.ERC20Token) *types.ContractCallTx {
	params := k.GetParams(ctx)

	newContractCallTx := &types.ContractCallTx{
		InvalidationNonce: invalidationNonce,
		InvalidationScope: invalidationScope,
		Address:           address.String(),
		Payload:           payload,
		Timeout:           params.ChainParams[chainID].TargetEvmTxTimeout,
		Tokens:            tokens,
		Fees:              fees,
		Height:            uint64(ctx.BlockHeight()),
	}

	var tokenString []string
	for _, token := range tokens {
		tokenString = append(tokenString, token.String())
	}

	var feeString []string
	for _, fee := range fees {
		feeString = append(feeString, fee.String())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMultisigUpdateRequest,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyChainID, fmt.Sprint(chainID)),
			sdk.NewAttribute(types.AttributeKeyContractCallInvalidationNonce, fmt.Sprint(invalidationNonce)),
			sdk.NewAttribute(types.AttributeKeyContractCallInvalidationScope, fmt.Sprint(invalidationScope)),
			sdk.NewAttribute(types.AttributeKeyContractCallAddress, fmt.Sprint(address.String())),
			sdk.NewAttribute(types.AttributeKeyContractCallPayload, string(payload)),
			sdk.NewAttribute(types.AttributeKeyContractCallTokens, strings.Join(tokenString, "|")),
			sdk.NewAttribute(types.AttributeKeyContractCallFees, strings.Join(feeString, "|")),
			sdk.NewAttribute(types.AttributeKeyEvmTxTimeout, strconv.FormatUint(params.ChainParams[chainID].TargetEvmTxTimeout, 10)),
		),
	)
	k.SetOutgoingTx(ctx, chainID, newContractCallTx)
	k.Logger(ctx).Info(
		"ContractCallTx created",
		"bridge_chain_id", fmt.Sprint(chainID),
		"invalidation_nonce", newContractCallTx.InvalidationNonce,
		"invalidation_scope", newContractCallTx.InvalidationScope,
		"address", address.String(),
		"payload", string(payload),
		"tokens", strings.Join(tokenString, "|"),
		"fees", strings.Join(feeString, "|"),
		"eth_tx_timeout", strconv.FormatUint(params.ChainParams[chainID].TargetEvmTxTimeout, 10),
	)
	return newContractCallTx
}

//////////////////////////////////////
// Observed Ethereum/Cosmos heights //
//////////////////////////////////////

// GetEVMHeightVote gets the height vote for a validator on an EVM chain
func (k Keeper) GetEVMHeightVote(ctx sdk.Context, chainID uint32, valAddress sdk.ValAddress) types.LatestEVMBlockHeight {
	store := ctx.KVStore(k.StoreKey)
	key := types.MakeEVMHeightVoteKey(chainID, valAddress)
	bz := store.Get(key)

	if len(bz) == 0 {
		return types.LatestEVMBlockHeight{
			CosmosHeight: 0,
			EVMHeight:    0,
		}
	}

	height := types.LatestEVMBlockHeight{}
	k.Cdc.MustUnmarshal(bz, &height)
	return height
}

// SetEthereumHeightVoteRecord sets the latest observed heights per validator
func (k Keeper) SetEVMHeightVote(ctx sdk.Context, chainID uint32, valAddress sdk.ValAddress, evmHeight uint64) {
	store := ctx.KVStore(k.StoreKey)
	height := types.LatestEVMBlockHeight{
		EVMHeight:    evmHeight,
		CosmosHeight: uint64(ctx.BlockHeight()),
	}
	key := types.MakeEVMHeightVoteKey(chainID, valAddress)
	store.Set(key, k.Cdc.MustMarshal(&height))
}

func (k Keeper) IterateEVMHeightVotes(ctx sdk.Context, cb func(val sdk.ValAddress, height types.LatestEVMBlockHeight) (stop bool)) {
	store := ctx.KVStore(k.StoreKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.EVMHeightVoteKey})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var height types.LatestEVMBlockHeight
		key := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.EVMHeightVoteKey}))
		val := sdk.ValAddress(key.Next(20))

		k.Cdc.MustUnmarshal(iter.Value(), &height)
		if cb(val, height) {
			break
		}
	}
}

/////////////////
// MIGRATE     //
/////////////////

// MigrateGravityContract Cleans up all state associated a previous gravity contract and set a new contract. This is intended to run in the upgrade handler.
// This implementation is partial at best. It does not contain necessary functionality to freeze the bridge.
// We will have yet to implement functionality to Migrate the Cosmos ERC20 tokens or any other ERC20 tokens bridged to the gravity contracts.
// This just does keeper state cleanup if a new gravity contract has been deployed
func (k Keeper) MigrateGravityContract(ctx sdk.Context, chainID uint32, newBridgeAddress string, bridgeDeploymentHeight uint64) {
	// Delete Any Outgoing TXs.

	prefixStoreOtx := prefix.NewStore(ctx.KVStore(k.StoreKey), types.OutgoingTxKeyPrefix(chainID))
	iterOtx := prefixStoreOtx.ReverseIterator(nil, nil)
	defer iterOtx.Close()
	for ; iterOtx.Valid(); iterOtx.Next() {

		var any cdctypes.Any
		k.Cdc.MustUnmarshal(iterOtx.Value(), &any)
		var otx types.OutgoingTx
		if err := k.Cdc.UnpackAny(&any, &otx); err != nil {
			panic(err)
		}
		// Delete any partial Eth Signatures handging around
		prefixStoreSig := prefix.NewStore(ctx.KVStore(k.StoreKey), types.EVMSignatureKeyStoreIndexPrefix(chainID, otx.GetStoreIndex()))
		iterSig := prefixStoreSig.Iterator(nil, nil)
		defer iterSig.Close()

		for ; iterSig.Valid(); iterSig.Next() {
			prefixStoreSig.Delete(iterSig.Key())
		}

		prefixStoreOtx.Delete(iterOtx.Key())
	}

	// Reset the last observed signer set nonce
	store := ctx.KVStore(k.StoreKey)
	store.Set(types.MakeLatestSignerSetTxNonceKey(chainID), sdk.Uint64ToBigEndian(0))

	// Reset all ethereum event nonces to zero
	k.SetLastObservedEventNonce(ctx, chainID, 0)
	k.iterateEVMEventVoteRecords(ctx, chainID, func(_ []byte, voteRecord *types.EVMEventVoteRecord) bool {
		for _, vote := range voteRecord.Votes {
			val, err := sdk.ValAddressFromBech32(vote)

			if err != nil {
				panic(err)
			}

			k.setLastEventNonceByValidator(ctx, chainID, val, 0)
		}

		return false
	})

	// Delete all EVM Events
	prefixStoreEVMEvent := prefix.NewStore(ctx.KVStore(k.StoreKey), types.EVMEventVoteRecordPrefix(chainID))
	iterEvent := prefixStoreEVMEvent.Iterator(nil, nil)
	defer iterEvent.Close()
	for ; iterEvent.Valid(); iterEvent.Next() {
		prefixStoreEVMEvent.Delete(iterEvent.Key())
	}

	// Set the Last oberved EVM Blockheight to zero
	height := types.LatestEVMBlockHeight{
		EVMHeight:    (bridgeDeploymentHeight - 1),
		CosmosHeight: uint64(ctx.BlockHeight()),
	}

	store.Set(types.MakeLastEVMBlockHeightKey(chainID), k.Cdc.MustMarshal(&height))

	k.SetLastObservedSignerSetTx(ctx, chainID, types.SignerSetTx{
		Nonce:   0,
		Height:  0,
		Signers: nil,
	})

	// Set the batch Nonce to zero
	store.Set(types.MakeLastOutgoingBatchNonceKey(chainID), sdk.Uint64ToBigEndian(0))

	// Update the bridge contract address
	params := k.GetParams(ctx)
	k.SetParams(ctx, params)
}
