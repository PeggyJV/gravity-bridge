package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
	"github.com/stretchr/testify/require"
)

func TestDeletePendingEventVoteRecords(t *testing.T) {
	env := CreateTestEnv(t)
	gk := env.GravityKeeper

	valAddr1 := sdk.ValAddress("cosmosvaloper1")
	valAddr2 := sdk.ValAddress("cosmosvaloper2")
	valAddr3 := sdk.ValAddress("cosmosvaloper3")

	gk.setLastObservedEventNonce(env.Context, 10)
	gk.setLastEventNonceByValidator(env.Context, valAddr1, 15)
	gk.setLastEventNonceByValidator(env.Context, valAddr2, 15)
	gk.setLastEventNonceByValidator(env.Context, valAddr3, 15)

	approvedEvent := types.ContractCallExecutedEvent{
		EventNonce:        10,
		EthereumHeight:    99,
		InvalidationScope: []byte("scope"),
		InvalidationNonce: 1,
	}
	unapprovedEvent1 := types.ContractCallExecutedEvent{
		EventNonce:        11,
		EthereumHeight:    100,
		InvalidationScope: []byte("scope1"),
		InvalidationNonce: 1,
	}
	unapprovedEvent2 := types.SignerSetTxExecutedEvent{
		EventNonce:     12,
		EthereumHeight: 100,
		Members: []*types.EthereumSigner{{
			Power:           100,
			EthereumAddress: "0x1234567890",
		}},
		SignerSetTxNonce: 5,
	}
	unapprovedEvent3 := types.BatchExecutedEvent{
		TokenContract:  "0x1234567890",
		EventNonce:     13,
		EthereumHeight: 100,
		BatchNonce:     20,
	}
	unapprovedEvent4 := types.ERC20DeployedEvent{
		EventNonce:     14,
		EthereumHeight: 100,
		TokenContract:  "0x1234567890",
	}
	unapprovedEvent5 := types.SendToCosmosEvent{
		EventNonce:     15,
		EthereumHeight: 100,
		TokenContract:  "0x1234567890",
		Amount:         sdk.NewInt(100),
	}

	event, err := types.PackEvent(&approvedEvent)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 10, []byte("scope"), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: true,
	})

	event, err = types.PackEvent(&unapprovedEvent1)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 11, OldHash(&unapprovedEvent1), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: false,
	})

	event, err = types.PackEvent(&unapprovedEvent2)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 12, OldHash(&unapprovedEvent2), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: false,
	})

	event, err = types.PackEvent(&unapprovedEvent3)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 13, OldHash(&unapprovedEvent3), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: false,
	})

	event, err = types.PackEvent(&unapprovedEvent4)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 14, OldHash(&unapprovedEvent4), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: false,
	})

	event, err = types.PackEvent(&unapprovedEvent5)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 15, OldHash(&unapprovedEvent5), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: false,
	})

	require.Equal(t, uint64(15), gk.getLastEventNonceByValidator(env.Context, valAddr1))
	require.Equal(t, uint64(15), gk.getLastEventNonceByValidator(env.Context, valAddr2))
	require.Equal(t, uint64(15), gk.getLastEventNonceByValidator(env.Context, valAddr3))
	require.Equal(t, uint64(10), gk.GetLastObservedEventNonce(env.Context))

	migrator := NewMigrator(gk)
	err = migrator.DeletePendingEventVoteRecords(env.Context)
	require.NoError(t, err)

	require.Equal(t, uint64(10), gk.GetLastObservedEventNonce(env.Context))
	require.NotNil(t, gk.GetEthereumEventVoteRecordMapping(env.Context)[10])
	require.Nil(t, gk.GetEthereumEventVoteRecordMapping(env.Context)[11])
	require.Nil(t, gk.GetEthereumEventVoteRecordMapping(env.Context)[12])
	require.Nil(t, gk.GetEthereumEventVoteRecordMapping(env.Context)[13])
	require.Nil(t, gk.GetEthereumEventVoteRecordMapping(env.Context)[14])
	require.Nil(t, gk.GetEthereumEventVoteRecordMapping(env.Context)[15])
	require.Equal(t, uint64(10), gk.getLastEventNonceByValidator(env.Context, valAddr1))
	require.Equal(t, uint64(10), gk.getLastEventNonceByValidator(env.Context, valAddr2))
	require.Equal(t, uint64(10), gk.getLastEventNonceByValidator(env.Context, valAddr3))

}

func TestDeletePendingEventVoteRecordsEmptyState(t *testing.T) {
	env := CreateTestEnv(t)
	gk := env.GravityKeeper
	migrator := NewMigrator(gk)

	// Test migration with no events
	err := migrator.DeletePendingEventVoteRecords(env.Context)
	require.NoError(t, err)
	require.Equal(t, uint64(0), gk.GetLastObservedEventNonce(env.Context))
	require.Empty(t, gk.GetEthereumEventVoteRecordMapping(env.Context))
}

func TestDeletePendingEventVoteRecordsPartialVotes(t *testing.T) {
	env := CreateTestEnv(t)
	gk := env.GravityKeeper

	valAddr1 := sdk.ValAddress("cosmosvaloper1")
	valAddr2 := sdk.ValAddress("cosmosvaloper2")
	valAddr3 := sdk.ValAddress("cosmosvaloper3")
	valAddr4 := sdk.ValAddress("cosmosvaloper4")
	valAddr5 := sdk.ValAddress("cosmosvaloper5")

	gk.setLastObservedEventNonce(env.Context, 5)
	gk.setLastEventNonceByValidator(env.Context, valAddr1, 7)
	gk.setLastEventNonceByValidator(env.Context, valAddr2, 6)
	gk.setLastEventNonceByValidator(env.Context, valAddr3, 7)
	gk.setLastEventNonceByValidator(env.Context, valAddr4, 3)
	gk.setLastEventNonceByValidator(env.Context, valAddr5, 5)

	// Create events with different voting patterns
	approvedEvent := types.ContractCallExecutedEvent{
		EventNonce:        5,
		EthereumHeight:    90,
		InvalidationScope: []byte("scope"),
		InvalidationNonce: 1,
	}
	partialEvent := types.ContractCallExecutedEvent{
		EventNonce:        6,
		EthereumHeight:    91,
		InvalidationScope: []byte("scope2"),
		InvalidationNonce: 2,
	}
	unapprovedEvent := types.ContractCallExecutedEvent{
		EventNonce:        7,
		EthereumHeight:    92,
		InvalidationScope: []byte("scope3"),
		InvalidationNonce: 3,
	}

	// Set up approved event
	event, err := types.PackEvent(&approvedEvent)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 5, OldHash(&approvedEvent), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: true,
	})

	// Set up partial votes event
	event, err = types.PackEvent(&partialEvent)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 6, OldHash(&partialEvent), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String()},
		Accepted: false,
	})

	// Set up unapproved event with single vote
	event, err = types.PackEvent(&unapprovedEvent)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 7, OldHash(&unapprovedEvent), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr3.String()},
		Accepted: false,
	})

	require.Equal(t, uint64(7), gk.getLastEventNonceByValidator(env.Context, valAddr1))
	require.Equal(t, uint64(6), gk.getLastEventNonceByValidator(env.Context, valAddr2))
	require.Equal(t, uint64(7), gk.getLastEventNonceByValidator(env.Context, valAddr3))
	require.Equal(t, uint64(3), gk.getLastEventNonceByValidator(env.Context, valAddr4))
	require.Equal(t, uint64(5), gk.getLastEventNonceByValidator(env.Context, valAddr5))
	require.Equal(t, 3, len(gk.GetEthereumEventVoteRecordMapping(env.Context)))

	migrator := NewMigrator(gk)
	err = migrator.DeletePendingEventVoteRecords(env.Context)
	require.NoError(t, err)

	// Verify state after migration
	require.Equal(t, uint64(5), gk.GetLastObservedEventNonce(env.Context))
	require.NotNil(t, gk.GetEthereumEventVoteRecordMapping(env.Context)[5])
	require.Nil(t, gk.GetEthereumEventVoteRecordMapping(env.Context)[6], "expected 6 to be deleted: %v", gk.GetEthereumEventVoteRecordMapping(env.Context)[6])
	require.Nil(t, gk.GetEthereumEventVoteRecordMapping(env.Context)[7], "expected 7 to be deleted: %v", gk.GetEthereumEventVoteRecordMapping(env.Context)[7])

	// Verify validator nonces are reset correctly
	require.Equal(t, uint64(5), gk.getLastEventNonceByValidator(env.Context, valAddr1))
	require.Equal(t, uint64(5), gk.getLastEventNonceByValidator(env.Context, valAddr2))
	require.Equal(t, uint64(5), gk.getLastEventNonceByValidator(env.Context, valAddr3))
}

func TestDeletePendingEventVoteRecordsInactiveValidators(t *testing.T) {
	env := CreateTestEnv(t)
	gk := env.GravityKeeper

	// Active validators who participate in voting
	valAddr1 := sdk.ValAddress("cosmosvaloper1")
	valAddr2 := sdk.ValAddress("cosmosvaloper2")
	valAddr3 := sdk.ValAddress("cosmosvaloper3")

	// Inactive validators with lower nonces
	inactiveVal1 := sdk.ValAddress("cosmosvaloper4")
	inactiveVal2 := sdk.ValAddress("cosmosvaloper5")

	// Set global last observed nonce
	gk.setLastObservedEventNonce(env.Context, 10)

	// Set active validators' nonces
	gk.setLastEventNonceByValidator(env.Context, valAddr1, 11)
	gk.setLastEventNonceByValidator(env.Context, valAddr2, 12)
	gk.setLastEventNonceByValidator(env.Context, valAddr3, 13)

	// Set inactive validators with lower nonces
	gk.setLastEventNonceByValidator(env.Context, inactiveVal1, 5)
	gk.setLastEventNonceByValidator(env.Context, inactiveVal2, 7)

	// Create and record some events
	approvedEvent := types.ContractCallExecutedEvent{
		EventNonce:        10,
		EthereumHeight:    99,
		InvalidationScope: []byte("scope"),
		InvalidationNonce: 1,
	}
	unapprovedEvent1 := types.ContractCallExecutedEvent{
		EventNonce:        11,
		EthereumHeight:    100,
		InvalidationScope: []byte("scope1"),
		InvalidationNonce: 1,
	}
	unapprovedEvent2 := types.ContractCallExecutedEvent{
		EventNonce:        12,
		EthereumHeight:    100,
		InvalidationScope: []byte("scope2"),
		InvalidationNonce: 2,
	}
	unapprovedEvent3 := types.ContractCallExecutedEvent{
		EventNonce:        13,
		EthereumHeight:    100,
		InvalidationScope: []byte("scope3"),
		InvalidationNonce: 3,
	}

	event, err := types.PackEvent(&approvedEvent)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 10, OldHash(&approvedEvent), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: true,
	})

	event, err = types.PackEvent(&unapprovedEvent1)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 11, OldHash(&unapprovedEvent1), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: false,
	})

	event, err = types.PackEvent(&unapprovedEvent2)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 12, OldHash(&unapprovedEvent2), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: false,
	})

	event, err = types.PackEvent(&unapprovedEvent3)
	require.NoError(t, err)
	gk.setEthereumEventVoteRecord(env.Context, 13, OldHash(&unapprovedEvent3), &types.EthereumEventVoteRecord{
		Event:    event,
		Votes:    []string{valAddr1.String(), valAddr2.String(), valAddr3.String()},
		Accepted: false,
	})

	// Verify initial state
	require.Equal(t, uint64(11), gk.getLastEventNonceByValidator(env.Context, valAddr1))
	require.Equal(t, uint64(12), gk.getLastEventNonceByValidator(env.Context, valAddr2))
	require.Equal(t, uint64(13), gk.getLastEventNonceByValidator(env.Context, valAddr3))
	require.Equal(t, uint64(5), gk.getLastEventNonceByValidator(env.Context, inactiveVal1))
	require.Equal(t, uint64(7), gk.getLastEventNonceByValidator(env.Context, inactiveVal2))

	// Run migration
	migrator := NewMigrator(gk)
	err = migrator.DeletePendingEventVoteRecords(env.Context)
	require.NoError(t, err)

	// Verify active validators were reset to last observed nonce
	require.Equal(t, uint64(10), gk.getLastEventNonceByValidator(env.Context, valAddr1))
	require.Equal(t, uint64(10), gk.getLastEventNonceByValidator(env.Context, valAddr2))
	require.Equal(t, uint64(10), gk.getLastEventNonceByValidator(env.Context, valAddr3))

	// Verify inactive validators' nonces were not changed
	require.Equal(t, uint64(5), gk.getLastEventNonceByValidator(env.Context, inactiveVal1))
	require.Equal(t, uint64(7), gk.getLastEventNonceByValidator(env.Context, inactiveVal2))
}
