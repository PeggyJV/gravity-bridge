package gravity

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/keeper"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Question: what here can be epoched?
	slashing(ctx, k)
	attestationTally(ctx, k)
	cleanupTimedOutBatches(ctx, k)
	cleanupTimedOutLogicCalls(ctx, k)
	createValsets(ctx, k)
	pruneSignerSetTxs(ctx, k)
}

func createValsets(ctx sdk.Context, k keeper.Keeper) {
	// Auto ValsetRequest Creation.
	// 1. If there are no valset requests, create a new one.
	// 2. If there is at least one validator who started unbonding in current block. (we persist last unbonded block height in hooks.go)
	//      This will make sure the unbonding validator has to provide an attestation to a new Valset
	//	    that excludes him before he completely Unbonds.  Otherwise he will be slashed
	// 3. If power change between validators of CurrentValset and latest valset request is > 5%
	latestValset := k.GetLatestSignerSetTx(ctx)
	lastUnbondingHeight := k.GetLastUnBondingBlockHeight(ctx)

	powerDiff := types.EthereumSigners(k.NewSignerSetTx(ctx).Signers).PowerDiff(latestValset.Signers)
	if (latestValset == nil) || (lastUnbondingHeight == uint64(ctx.BlockHeight())) || (powerDiff > 0.05) {
		k.SetOutgoingTx(ctx, k.NewSignerSetTx(ctx))
	}
}

func pruneSignerSetTxs(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	// Validator set pruning
	// prune all validator sets with a nonce less than the
	// last observed nonce, they can't be submitted any longer
	//
	// Only prune valsets after the signed valsets window has passed
	// so that slashing can occur the block before we remove them
	lastObserved := k.GetLastObservedValset(ctx)
	currentBlock := uint64(ctx.BlockHeight())
	tooEarly := currentBlock < params.SignedSignerSetTxsWindow
	if lastObserved != nil && !tooEarly {
		//earliestToPrune := currentBlock - params.SignedSignerSetTxsWindow
		sets := k.GetSignerSetTxs(ctx)
		for _, set := range sets {
			// TODO: do we need height on signersettx?
			if set.Nonce < lastObserved.Nonce { // && set.Height < earliestToPrune {
				k.DeleteOutgoingTx(ctx, set.GetStoreIndex())
			}
		}
	}
}

func slashing(ctx sdk.Context, k keeper.Keeper) {

	params := k.GetParams(ctx)

	// Slash validator for not confirming valset requests, batch requests and not attesting claims rightfully
	ValsetSlashing(ctx, k, params)
	BatchSlashing(ctx, k, params)
	// TODO slashing for arbitrary logic is missing

	// TODO: prune validator sets, older than 6 months, this time is chosen out of an abundance of caution
	// TODO: prune outgoing tx batches while looping over them above, older than 15h and confirmed
	// TODO: prune claims, attestations
}

// Iterate over all attestations currently being voted on in order of nonce and
// "Observe" those who have passed the threshold. Break the loop once we see
// an attestation that has not passed the threshold
func attestationTally(ctx sdk.Context, k keeper.Keeper) {
	attmap := k.GetEthereumEventVoteRecordMapping(ctx)
	// We make a slice with all the event nonces that are in the attestation mapping
	keys := make([]uint64, 0, len(attmap))
	for k := range attmap {
		keys = append(keys, k)
	}
	// Then we sort it
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	// This iterates over all keys (event nonces) in the attestation mapping. Each value contains
	// a slice with one or more attestations at that event nonce. There can be multiple attestations
	// at one event nonce when validators disagree about what event happened at that nonce.
	for _, nonce := range keys {
		// This iterates over all attestations at a particular event nonce.
		// They are ordered by when the first attestation at the event nonce was received.
		// This order is not important.
		for _, att := range attmap[nonce] {
			// We check if the event nonce is exactly 1 higher than the last attestation that was
			// observed. If it is not, we just move on to the next nonce. This will skip over all
			// attestations that have already been observed.
			//
			// Once we hit an event nonce that is one higher than the last observed event, we stop
			// skipping over this conditional and start calling tryAttestation (counting votes)
			// Once an attestation at a given event nonce has enough votes and becomes observed,
			// every other attestation at that nonce will be skipped, since the lastObservedEventNonce
			// will be incremented.
			//
			// Then we go to the next event nonce in the attestation mapping, if there is one. This
			// nonce will once again be one higher than the lastObservedEventNonce.
			// If there is an attestation at this event nonce which has enough votes to be observed,
			// we skip the other attestations and move on to the next nonce again.
			// If no attestation becomes observed, when we get to the next nonce, every attestation in
			// it will be skipped. The same will happen for every nonce after that.
			if nonce == uint64(k.GetLastObservedEventNonce(ctx))+1 {
				k.TryEventVoteRecord(ctx, &att)
			}
		}
	}
}

// cleanupTimedOutBatches deletes batches that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning batch 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single batch at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutBatches(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx).EthereumHeight
	k.IterateOutgoingTxs(ctx, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		btx, _ := otx.(*types.BatchTx)

		if btx.Timeout < ethereumHeight {
			k.CancelBatchTx(ctx, common.HexToAddress(btx.TokenContract), btx.Nonce)
		}

		return false
	})
}

// cleanupTimedOutBatches deletes logic calls that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning call 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single call at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutLogicCalls(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx).EthereumHeight
	var calls []*types.ContractCallTx
	k.IterateOutgoingTxs(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		cctx, _ := otx.(*types.ContractCallTx)
		calls = append(calls, cctx)
		return true
	})
	for _, call := range calls {
		if call.Timeout < ethereumHeight {
			k.CancelContractCallTx(ctx, call.InvalidationScope, call.InvalidationNonce)
		}
	}
}

func ValsetSlashing(ctx sdk.Context, k keeper.Keeper, params types.Params) {

	maxHeight := uint64(0)

	// don't slash in the beginning before there aren't even SignedValsetsWindow blocks yet
	if uint64(ctx.BlockHeight()) > params.SignedValsetsWindow {
		maxHeight = uint64(ctx.BlockHeight()) - params.SignedValsetsWindow
	}

	unslashedValsets := k.GetUnSlashedValsets(ctx, maxHeight)

	// unslashedValsets are sorted by nonce in ASC order
	// Question: do we need to sort each time? See if this can be epoched
	for _, vs := range unslashedValsets {
		confirms := k.GetEthereumSignatures(ctx, vs.Nonce)

		// SLASH BONDED VALIDTORS who didn't attest valset request
		currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
		for _, val := range currentBondedSet {
			consAddr, _ := val.GetConsAddr()
			valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)

			//  Slash validator ONLY if he joined after valset is created
			if exist && valSigningInfo.StartHeight < int64(vs.Nonce) {
				// Check if validator has confirmed valset or not
				found := false
				for _, conf := range confirms {
					if conf.EthereumSigner == k.GetEthAddress(ctx, val.GetOperator()) {
						found = true
						break
					}
				}
				// slash validators for not confirming valsets
				if !found {
					cons, _ := val.GetConsAddr()
					k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionValset)
					if !val.IsJailed() {
						k.StakingKeeper.Jail(ctx, cons)
					}

				}
			}
		}

		// SLASH UNBONDING VALIDATORS who didn't attest valset request
		blockTime := ctx.BlockTime().Add(k.StakingKeeper.GetParams(ctx).UnbondingTime)
		blockHeight := ctx.BlockHeight()
		unbondingValIterator := k.StakingKeeper.ValidatorQueueIterator(ctx, blockTime, blockHeight)
		defer unbondingValIterator.Close()

		// All unbonding validators
		for ; unbondingValIterator.Valid(); unbondingValIterator.Next() {
			unbondingValidators := k.GetUnbondingvalidators(unbondingValIterator.Value())

			for _, valAddr := range unbondingValidators.Addresses {
				addr, err := sdk.ValAddressFromBech32(valAddr)
				if err != nil {
					panic(err)
				}
				validator, _ := k.StakingKeeper.GetValidator(ctx, sdk.ValAddress(addr))
				valConsAddr, _ := validator.GetConsAddr()
				valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, valConsAddr)

				// Only slash validators who joined after valset is created and they are unbonding and UNBOND_SLASHING_WINDOW didn't passed
				if exist && valSigningInfo.StartHeight < int64(vs.Nonce) && validator.IsUnbonding() && vs.Nonce < uint64(validator.UnbondingHeight)+params.UnbondSlashingValsetsWindow {
					// Check if validator has confirmed valset or not
					found := false
					for _, conf := range confirms {
						if conf.EthereumSigner == k.GetEthAddress(ctx, validator.GetOperator()) {
							found = true
							break
						}
					}

					// slash validators for not confirming valsets
					if !found {
						k.StakingKeeper.Slash(ctx, valConsAddr, ctx.BlockHeight(), validator.ConsensusPower(), params.SlashFractionValset)
						if !validator.IsJailed() {
							k.StakingKeeper.Jail(ctx, valConsAddr)
						}
					}
				}
			}
		}
		// then we set the latest slashed valset  nonce
		k.SetLastSlashedValsetNonce(ctx, vs.Nonce)
	}
}

func BatchSlashing(ctx sdk.Context, k keeper.Keeper, params types.Params) {

	// #2 condition
	// We look through the full bonded set (not just the active set, include unbonding validators)
	// and we slash users who haven't signed a batch confirmation that is >15hrs in blocks old
	maxHeight := uint64(0)

	// don't slash in the beginning before there aren't even SignedBatchesWindow blocks yet
	if uint64(ctx.BlockHeight()) > params.SignedBatchesWindow {
		maxHeight = uint64(ctx.BlockHeight()) - params.SignedBatchesWindow
	}

	unslashedBatches := k.GetUnSlashedBatches(ctx, maxHeight)
	for _, batch := range unslashedBatches {

		// SLASH BONDED VALIDTORS who didn't attest batch requests
		currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
		confirms := k.GetBatchTxSignatureByNonceAndTokenContract(ctx, batch.Nonce, batch.TokenContract)
		for _, val := range currentBondedSet {
			// Don't slash validators who joined after batch is created
			consAddr, _ := val.GetConsAddr()
			valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
			if exist && valSigningInfo.StartHeight > int64(batch.EthereumBlock) {
				continue
			}

			found := false
			for _, conf := range confirms {
				// TODO: review this thoroughly
				// TODO: This is currently WRONG! We need to use the EthereumSigner here to
				// get the validator address.
				confVal, _ := sdk.AccAddressFromBech32(conf.EthereumSigner)
				if k.GetOrchestratorValidator(ctx, confVal).Equals(val.GetOperator()) {
					found = true
					break
				}
			}
			if !found {
				cons, _ := val.GetConsAddr()
				k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionBatch)
				if !val.IsJailed() {
					k.StakingKeeper.Jail(ctx, cons)
				}
			}
		}
		// then we set the latest slashed batch block
		k.SetLastSlashedBatchBlock(ctx, batch.EthereumBlock)

	}
}

// TestingEndBlocker is a second endblocker function only imported in the Gravity codebase itself
// if you are a consuming Cosmos chain DO NOT IMPORT THIS, it simulates a chain using the arbitrary
// logic API to request logic calls
func TestingEndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// if this is nil we have not set our test outgoing logic call yet
	if k.GetContractCallTx(ctx, []byte("GravityTesting"), 0).Payload == nil {
		// TODO this call isn't actually very useful for testing, since it always
		// throws, being just junk data that's expected. But it prevents us from checking
		// the full lifecycle of the call. We need to find some way for this to read data
		// and encode a simple testing call, probably to one of the already deployed ERC20
		// contracts so that we can get the full lifecycle.
		token := []types.ERC20Token{{
			Contract: "0x7580bfe88dd3d07947908fae12d95872a260f2d8",
			Amount:   sdk.NewIntFromUint64(5000),
		}}
		_ = types.ContractCallTx{
			Tokens:            token,
			Fees:              token,
			Address:           "0x510ab76899430424d209a6c9a5b9951fb8a6f47d",
			Payload:           []byte("fake bytes"),
			Timeout:           10000,
			InvalidationScope: []byte("GravityTesting"),
			InvalidationNonce: 1,
		}
		//k.SetContractCallTx(ctx, &call)
	}
}
