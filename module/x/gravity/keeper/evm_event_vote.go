package keeper

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

func (k Keeper) recordEventVote(
	ctx sdk.Context,
	event types.EVMEvent,
	val sdk.ValAddress,
) (*types.EVMEventVoteRecord, error) {
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processEVMEvent as well,
	// but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two claims with the same nonce
	lastEventNonce := k.getLastEventNonceByValidator(ctx, event.ChainID(), val)
	expectedNonce := lastEventNonce + 1
	if event.GetEventNonce() != expectedNonce {
		return nil, sdkerrors.Wrapf(types.ErrInvalid,
			"non contiguous event nonce expected %v observed %v for validator %v",
			expectedNonce,
			event.GetEventNonce(),
			val,
		)
	}

	// Tries to get an EVMEventVoteRecord with the same eventNonce and event as the event that was submitted.
	eventVoteRecord := k.GetEVMEventVoteRecord(ctx, event.ChainID(), event.GetEventNonce(), event.Hash())

	// If it does not exist, create a new one.
	if eventVoteRecord == nil {
		any, err := types.PackEvent(event)
		if err != nil {
			return nil, err
		}
		eventVoteRecord = &types.EVMEventVoteRecord{
			Accepted: false,
			Event:    any,
		}
	}

	// Add the validator's vote to this EVMEventVoteRecord
	eventVoteRecord.Votes = append(eventVoteRecord.Votes, val.String())

	k.setEVMEventVoteRecord(ctx, event.ChainID(), event.GetEventNonce(), event.Hash(), eventVoteRecord)
	k.setLastEventNonceByValidator(ctx, event.ChainID(), val, event.GetEventNonce())

	return eventVoteRecord, nil
}

// TryEventVoteRecord checks if an event vote record has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processEVMEvent to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TryEventVoteRecord(ctx sdk.Context, chainID uint32, eventVoteRecord *types.EVMEventVoteRecord) {
	// If the event vote record has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the event vote record from accidentally being applied twice.
	if !eventVoteRecord.Accepted {
		var event types.EVMEvent
		if err := k.cdc.UnpackAny(eventVoteRecord.Event, &event); err != nil {
			panic("unpacking packed any")
		}

		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		// TODO: The different integer types and math here needs a careful review
		requiredPower := types.EventVoteRecordPowerThreshold(k.StakingKeeper.GetLastTotalPower(ctx))
		eventVotePower := sdk.NewInt(0)
		for _, validator := range eventVoteRecord.Votes {
			val, _ := sdk.ValAddressFromBech32(validator)

			validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, val)
			// Add it to the attestation power's sum
			eventVotePower = eventVotePower.Add(sdk.NewInt(validatorPower))
			// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
			// process the attestation, set Observed to true, and break
			if eventVotePower.GTE(requiredPower) {
				lastEventNonce := k.GetLastObservedEventNonce(ctx, chainID)
				// this check is performed at the next level up so this should never panic
				// outside of programmer error.
				if event.GetEventNonce() != lastEventNonce+1 {
					panic("attempting to apply events to state out of order")
				}
				k.setLastObservedEventNonce(ctx, chainID, event.GetEventNonce())
				k.SetLastObservedEVMBlockHeight(ctx, chainID, event.GetEVMHeight())

				eventVoteRecord.Accepted = true
				k.setEVMEventVoteRecord(ctx, chainID, event.GetEventNonce(), event.Hash(), eventVoteRecord)

				k.processEVMEvent(ctx, chainID, event)
				ctx.EventManager().EmitEvent(sdk.NewEvent(
					types.EventTypeObservation,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
					sdk.NewAttribute(types.AttributeKeyEVMEventType, fmt.Sprintf("%T", event)),
					sdk.NewAttribute(types.AttributeKeyContract, k.getBridgeContractAddress(ctx)),
					sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.getBridgeChainID(ctx)))),
					sdk.NewAttribute(types.AttributeKeyEVMEventVoteRecordID,
						string(types.MakeEVMEventVoteRecordKey(chainID, event.GetEventNonce(), event.Hash()))),
					sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(event.GetEventNonce())),
				))

				break
			}
		}
	} else {
		// We panic here because this should never happen
		panic("attempting to process observed ethereum event")
	}
}

// processEVMEvent actually applies the attestation to the consensus state
func (k Keeper) processEVMEvent(ctx sdk.Context, chainID uint32, event types.EVMEvent) {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.Handle(xCtx, event); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.Logger(ctx).Error(
			"ethereum event vote record failed",
			"cause", err.Error(),
			"event type", fmt.Sprintf("%T", event),
			"id", types.MakeEVMEventVoteRecordKey(chainID, event.GetEventNonce(), event.Hash()),
			"nonce", fmt.Sprint(event.GetEventNonce()),
		)
	} else {
		ctx.EventManager().EmitEvents(xCtx.EventManager().Events()) // copy events to original context
		commit()                                                    // persist transient storage
	}
}

// setEVMEventVoteRecord sets the attestation in the store
func (k Keeper) setEVMEventVoteRecord(ctx sdk.Context, chainID uint32, eventNonce uint64, claimHash []byte, eventVoteRecord *types.EVMEventVoteRecord) {
	ctx.KVStore(k.storeKey).Set(types.MakeEVMEventVoteRecordKey(chainID, eventNonce, claimHash), k.cdc.MustMarshal(eventVoteRecord))
}

// GetEVMEventVoteRecord return a vote record given a nonce
func (k Keeper) GetEVMEventVoteRecord(ctx sdk.Context, chainID uint32, eventNonce uint64, claimHash []byte) *types.EVMEventVoteRecord {
	if bz := ctx.KVStore(k.storeKey).Get(types.MakeEVMEventVoteRecordKey(chainID, eventNonce, claimHash)); bz == nil {
		return nil
	} else {
		var out types.EVMEventVoteRecord
		k.cdc.MustUnmarshal(bz, &out)
		return &out
	}
}

// GetEVMEventVoteRecordMapping returns a mapping of eventnonce -> attestations at that nonce
func (k Keeper) GetEVMEventVoteRecordMapping(ctx sdk.Context, chainID uint32) (out map[uint64][]*types.EVMEventVoteRecord) {
	out = make(map[uint64][]*types.EVMEventVoteRecord)
	k.iterateEVMEventVoteRecords(ctx, chainID, func(key []byte, eventVoteRecord *types.EVMEventVoteRecord) bool {
		event, err := types.UnpackEvent(eventVoteRecord.Event)
		if err != nil {
			panic(err)
		}
		if val, ok := out[event.GetEventNonce()]; !ok {
			out[event.GetEventNonce()] = []*types.EVMEventVoteRecord{eventVoteRecord}
		} else {
			out[event.GetEventNonce()] = append(val, eventVoteRecord)
		}
		return false
	})
	return
}

// iterateEVMEventVoteRecords iterates through all attestations
func (k Keeper) iterateEVMEventVoteRecords(ctx sdk.Context, chainID uint32, cb func([]byte, *types.EVMEventVoteRecord) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EVMEventVoteRecordPrefix(chainID))
	iter := store.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		att := &types.EVMEventVoteRecord{}
		k.cdc.MustUnmarshal(iter.Value(), att)
		// cb returns true to stop early
		if cb(iter.Key(), att) {
			return
		}
	}
}

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context, chainID uint32) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.MakeLastObservedEventNonceKey(chainID))

	if len(bytes) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(bytes)
}

// GetLastObservedEVMBlockHeight height gets the block height to of the last observed attestation from
// the store
func (k Keeper) GetLastObservedEVMBlockHeight(ctx sdk.Context, chainID uint32) types.LatestEVMBlockHeight {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.MakeLastEVMBlockHeightKey(chainID))

	if len(bytes) == 0 {
		return types.LatestEVMBlockHeight{
			CosmosHeight: 0,
			EVMHeight:    0,
		}
	}
	height := types.LatestEVMBlockHeight{}
	k.cdc.MustUnmarshal(bytes, &height)
	return height
}

// SetLastObservedEVMBlockHeight sets the block height in the store.
func (k Keeper) SetLastObservedEVMBlockHeight(ctx sdk.Context, chainID uint32, ethereumHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	height := types.LatestEVMBlockHeight{
		EVMHeight:    ethereumHeight,
		CosmosHeight: uint64(ctx.BlockHeight()),
	}
	store.Set(types.MakeLastEVMBlockHeightKey(chainID), k.cdc.MustMarshal(&height))
}

// setLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) setLastObservedEventNonce(ctx sdk.Context, chainID uint32, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MakeLastObservedEventNonceKey(chainID), sdk.Uint64ToBigEndian(nonce))
}

// getLastEventNonceByValidator returns the latest event nonce for a given validator
func (k Keeper) getLastEventNonceByValidator(ctx sdk.Context, chainID uint32, validator sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.MakeLastEventNonceByValidatorKey(chainID, validator))

	if len(bytes) == 0 {
		// in the case that we have no existing value this is the first
		// time a validator is submitting a claim. Since we don't want to force
		// them to replay the entire history of all events ever we can't start
		// at zero
		//
		// We could start at the LastObservedEventNonce but if we do that this
		// validator will be slashed, because they are responsible for making a claim
		// on any attestation that has not yet passed the slashing window.
		//
		// Therefore we need to return to them the lowest attestation that is still within
		// the slashing window. Since we delete attestations after the slashing window that's
		// just the lowest observed event in the store. If no claims have been submitted in for
		// params.SignedClaimsWindow we may have no attestations in our nonce. At which point
		// the last observed which is a persistent and never cleaned counter will suffice.
		lowestObserved := k.GetLastObservedEventNonce(ctx, chainID)
		attmap := k.GetEVMEventVoteRecordMapping(ctx, chainID)
		// no new claims in params.SignedClaimsWindow, we can return the current value
		// because the validator can't be slashed for an event that has already passed.
		// so they only have to worry about the *next* event to occur
		if len(attmap) == 0 {
			return lowestObserved
		}
		for nonce, atts := range attmap {
			for att := range atts {
				if atts[att].Accepted && nonce < lowestObserved {
					lowestObserved = nonce
				}
			}
		}
		// return the latest event minus one so that the validator
		// can submit that event and avoid slashing. special case
		// for zero
		if lowestObserved > 0 {
			return lowestObserved - 1
		}
		return 0
	}
	return binary.BigEndian.Uint64(bytes)
}

// setLastEventNonceByValidator sets the latest event nonce for a give validator
func (k Keeper) setLastEventNonceByValidator(ctx sdk.Context, chainID uint32, validator sdk.ValAddress, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MakeLastEventNonceByValidatorKey(chainID, validator), sdk.Uint64ToBigEndian(nonce))
}
