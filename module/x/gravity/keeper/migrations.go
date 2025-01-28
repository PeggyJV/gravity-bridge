package keeper

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) MigrateStore(ctx sdk.Context) error {
	ctx.Logger().Info("gravity: Migrating store")

	m.DeletePendingEventVoteRecords(ctx)

	return nil
}

// DeletePendingEventVoteRecords deletes pending event vote records and adjusts the last observed nonce for validators
// who voted on unapproved events. This upgrade includes changes to how event hashes are calculated, so we delete
// pending event vote records that were created with the old hash calculation method to prevent inconsistent hashes.
func (m Migrator) DeletePendingEventVoteRecords(ctx sdk.Context) error {
	ctx.Logger().Info("gravity: Deleting pending event vote records")
	lastObservedEventNonce := m.keeper.GetLastObservedEventNonce(ctx)
	ctx.Logger().Info("gravity: Last observed event nonce", "nonce", lastObservedEventNonce)

	// delete event vote records with nonce greater than lastObservedEventNonce
	type deleteInfo struct {
		nonce uint64
		hash  []byte
	}

	var recordsToDelete []deleteInfo
	var lowestNonce uint64
	var unapprovedEventSigners []string
	m.keeper.IterateEthereumEventVoteRecords(ctx, func(key []byte, eventVoteRecord *types.EthereumEventVoteRecord) bool {
		// Skip approved events
		if eventVoteRecord.Accepted {
			return false
		}

		event, err := types.UnpackEvent(eventVoteRecord.Event)
		if err != nil {
			panic(err)
		}

		if event.GetEventNonce() > lastObservedEventNonce {
			v5Hash := OldHash(event)
			ctx.Logger().Info("gravity: v5 event hash: ", v5Hash)
			ctx.Logger().Info("gravity: v6 event hash: ", event.Hash())
			ctx.Logger().Info("gravity: Deleting pending event vote record", "nonce", event.GetEventNonce(), "hash", v5Hash)
			unapprovedEventSigners = append(unapprovedEventSigners, eventVoteRecord.Votes...)
			recordsToDelete = append(recordsToDelete, deleteInfo{
				nonce: event.GetEventNonce(),
				hash:  v5Hash,
			})
			if lowestNonce == 0 || event.GetEventNonce() < lowestNonce {
				lowestNonce = event.GetEventNonce()
			}
		}
		return false
	})

	if len(recordsToDelete) == 0 {
		ctx.Logger().Info("gravity: No pending event vote records to delete")
		return nil
	}

	for _, record := range recordsToDelete {
		m.keeper.DeleteEthereumEventVoteRecord(ctx, record.nonce, record.hash)
	}

	// Dedup the list of validators
	var dedupedValidators []string
	dedupedValidatorsMap := make(map[string]struct{})
	for _, signer := range unapprovedEventSigners {
		if _, exists := dedupedValidatorsMap[signer]; !exists {
			dedupedValidatorsMap[signer] = struct{}{}
			dedupedValidators = append(dedupedValidators, signer)
		}
	}

	// For each validator who voted on unapproved events, set their last observed
	// nonce to the global last observed nonce
	for _, signer := range dedupedValidators {
		validator, err := sdk.ValAddressFromBech32(signer)
		if err != nil {
			ctx.Logger().Error("gravity: Error converting validator address from bech32", "error", err)
			panic(err)
		}

		currentNonce := m.keeper.getLastEventNonceByValidator(ctx, validator)
		if currentNonce < lastObservedEventNonce {
			// This shouldn't happen unless there's a bug but we just leave it be since the skipping event is getting deleted.
			ctx.Logger().Error("gravity: Validator that signed an unapproved event has a last observed nonce less than the last accepted nonce. This is a bug in the module code.", "validator", signer, "current_nonce", currentNonce, "last_observed_nonce", lastObservedEventNonce)
			continue
		}

		m.keeper.setLastEventNonceByValidator(ctx, validator, lastObservedEventNonce)
	}

	return nil
}

func OldHash(event types.EthereumEvent) tmbytes.HexBytes {
	switch e := event.(type) {
	case *types.SendToCosmosEvent:
		return HashSendToCosmosEvent(event.(*types.SendToCosmosEvent))
	case *types.BatchExecutedEvent:
		return HashBatchExecutedEvent(event.(*types.BatchExecutedEvent))
	case *types.ContractCallExecutedEvent:
		return HashContractCallExecutedEvent(event.(*types.ContractCallExecutedEvent))
	case *types.SignerSetTxExecutedEvent:
		return HashSignerSetTxExecutedEvent(event.(*types.SignerSetTxExecutedEvent))
	case *types.ERC20DeployedEvent:
		return HashERC20DeployedEvent(event.(*types.ERC20DeployedEvent))
	default:
		panic(fmt.Sprintf("unknown event type: %T", e))
	}
}

func HashSendToCosmosEvent(event *types.SendToCosmosEvent) tmbytes.HexBytes {
	rcv, _ := sdk.AccAddressFromBech32(event.CosmosReceiver)
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(event.EventNonce),
			common.HexToAddress(event.TokenContract).Bytes(),
			event.Amount.BigInt().Bytes(),
			common.Hex2Bytes(event.EthereumSender),
			rcv.Bytes(),
			sdk.Uint64ToBigEndian(event.EthereumHeight),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func HashBatchExecutedEvent(event *types.BatchExecutedEvent) tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			common.HexToAddress(event.TokenContract).Bytes(),
			sdk.Uint64ToBigEndian(event.EventNonce),
			sdk.Uint64ToBigEndian(event.BatchNonce),
			sdk.Uint64ToBigEndian(event.EthereumHeight),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func HashContractCallExecutedEvent(event *types.ContractCallExecutedEvent) tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(event.EventNonce),
			event.InvalidationScope,
			sdk.Uint64ToBigEndian(event.InvalidationNonce),
			sdk.Uint64ToBigEndian(event.EthereumHeight),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func HashERC20DeployedEvent(event *types.ERC20DeployedEvent) tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(event.EventNonce),
			[]byte(event.CosmosDenom),
			common.HexToAddress(event.TokenContract).Bytes(),
			[]byte(event.Erc20Name),
			[]byte(event.Erc20Symbol),
			sdk.Uint64ToBigEndian(event.Erc20Decimals),
			sdk.Uint64ToBigEndian(event.EthereumHeight),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func HashSignerSetTxExecutedEvent(event *types.SignerSetTxExecutedEvent) tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(event.EventNonce),
			sdk.Uint64ToBigEndian(event.SignerSetTxNonce),
			sdk.Uint64ToBigEndian(event.EthereumHeight),
			types.EthereumSigners(event.Members).Hash(),
		},
		[]byte{},
	)
	hash := sha256.Sum256(([]byte(path)))
	return hash[:]
}
