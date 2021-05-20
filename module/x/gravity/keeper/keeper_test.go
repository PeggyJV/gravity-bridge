package keeper

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func TestCurrentValsetNormalization(t *testing.T) {
	specs := map[string]struct {
		srcPowers []uint64
		expPowers []uint64
	}{
		"one": {
			srcPowers: []uint64{100},
			expPowers: []uint64{4294967295},
		},
		"two": {
			srcPowers: []uint64{100, 1},
			expPowers: []uint64{4252442866, 42524428},
		},
	}
	input := CreateTestEnv(t)
	ctx := input.Context
	for msg, spec := range specs {
		spec := spec
		t.Run(msg, func(t *testing.T) {
			operators := make([]MockStakingValidatorData, len(spec.srcPowers))
			for i, v := range spec.srcPowers {
				operators[i] = MockStakingValidatorData{
					// any unique addr
					Operator: bytes.Repeat([]byte{byte(i)}, sdk.AddrLen),
					Power:    int64(v),
				}
			}
			input.GravityKeeper.StakingKeeper = NewStakingKeeperWeightedMock(operators...)
			r := input.GravityKeeper.CreateSignerSetTx(ctx)
			assert.Equal(t, spec.expPowers, types.EthereumSigners(r.Signers).GetPowers())
		})
	}
}

func TestAttestationIterator(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	// add some attestations to the store

	att1 := &types.EthereumEventVoteRecord{
		Accepted: true,
		Votes:    []string{},
	}
	dep1 := &types.SendToCosmosEvent{
		EventNonce:     1,
		TokenContract:  TokenContractAddrs[0],
		Amount:         sdk.NewInt(100),
		EthereumSender: EthAddrs[0].String(),
		CosmosReceiver: AccAddrs[0].String(),
	}
	att2 := &types.EthereumEventVoteRecord{
		Accepted: true,
		Votes:    []string{},
	}
	dep2 := &types.SendToCosmosEvent{
		EventNonce:     2,
		TokenContract:  TokenContractAddrs[0],
		Amount:         sdk.NewInt(100),
		EthereumSender: EthAddrs[0].String(),
		CosmosReceiver: AccAddrs[0].String(),
	}
	input.GravityKeeper.setEthereumEventVoteRecord(ctx, dep1.EventNonce, dep1.Hash(), att1)
	input.GravityKeeper.setEthereumEventVoteRecord(ctx, dep2.EventNonce, dep2.Hash(), att2)

	atts := []*types.EthereumEventVoteRecord{}
	input.GravityKeeper.iterateEthereumEventVoteRecords(ctx, func(_ []byte, att *types.EthereumEventVoteRecord) bool {
		atts = append(atts, att)
		return false
	})

	require.Len(t, atts, 2)
}

func TestDelegateKeys(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	k := input.GravityKeeper
	var (
		ethAddrs = []common.Address{
			common.HexToAddress("0x3146D2d6Eed46Afa423969f5dDC3152DfC359b09"),
			common.HexToAddress("0x610277F0208D342C576b991daFdCb36E36515e76"),
			common.HexToAddress("0x835973768750b3ED2D5c3EF5AdcD5eDb44d12aD4"),
			common.HexToAddress("0xb2A7F3E84F8FdcA1da46c810AEa110dd96BAE6bF"),
		}
		valAddrs = []string{
			"cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4",
			"cosmosvaloper15n79nty2fj37ant3p2gj4wju4ls6eu6tjwmdt0",
			"cosmosvaloper16dnkc6ac6ruuyr6l372fc3p77jgjpet6fka0cq",
			"cosmosvaloper1vrptwhl3ht2txmzy28j9msqkcvmn8gjz507pgu",
		}
		orchAddrs = []string{
			"cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7",
			"cosmos1rhfs24tlw4na04v35tzmjncy785kkw9j27d5kx",
			"cosmos10upq3tmt04zf55f6hw67m0uyrda3mp722q70rw",
			"cosmos1nt2uwjh5peg9vz2wfh2m3jjwqnu9kpjlhgpmen",
		}
	)
	for i := range ethAddrs {
		// set some addresses
		val, err1 := sdk.ValAddressFromBech32(valAddrs[i])
		orch, err2 := sdk.AccAddressFromBech32(orchAddrs[i])
		require.NoError(t, err1)
		require.NoError(t, err2)

		k.SetOrchestratorValidatorAddress(ctx, val, orch)
		k.setValidatorEthereumAddress(ctx, val, ethAddrs[i])
		k.setEthereumOrchestratorAddress(ctx, ethAddrs[i], orch)
	}
	addresses := k.getDelegateKeys(ctx)
	for i := range addresses {
		res := addresses[i]
		assert.Equal(t, valAddrs[i], res.ValidatorAddress)
		assert.Equal(t, orchAddrs[i], res.OrchestratorAddress)
		assert.Equal(t, ethAddrs[i].Hex(), res.EthereumAddress)
	}
}

func TestStoreEventVoteRecord(t *testing.T) {
	input := CreateTestEnv(t)
	gk := input.GravityKeeper
	ctx := input.Context
	stce := &types.SendToCosmosEvent{
		EventNonce:     1,
		TokenContract:  EthAddrs[0].Hex(),
		EthereumSender: EthAddrs[0].Hex(),
		CosmosReceiver: AccAddrs[0].String(),
		EthereumHeight: 10,
		Amount:         sdk.NewInt(1000000),
	}
	stcea, err := types.PackEvent(stce)
	require.NoError(t, err)

	evr := &types.EthereumEventVoteRecord{
		Event: stcea,
		Votes: []string{
			ValAddrs[0].String(),
			ValAddrs[1].String(),
			ValAddrs[2].String(),
		},
		Accepted: false,
	}

	cctxe := &types.ContractCallExecutedEvent{
		EventNonce:        2,
		InvalidationId:    []byte{0x1, 0x2},
		InvalidationNonce: 1,
		EthereumHeight:    11,
	}

	cctxea, err := types.PackEvent(cctxe)
	require.NoError(t, err)

	evr2 := &types.EthereumEventVoteRecord{
		Event: cctxea,
		Votes: []string{
			ValAddrs[2].String(),
			ValAddrs[3].String(),
			ValAddrs[4].String(),
		},
	}

	gk.setEthereumEventVoteRecord(ctx, stce.GetEventNonce(), stce.Hash(), evr)
	gk.setEthereumEventVoteRecord(ctx, cctxe.GetEventNonce(), cctxe.Hash(), evr2)

	stored := gk.GetEthereumEventVoteRecord(ctx, stce.GetEventNonce(), stce.Hash())
	require.NotNil(t, stored)

	stored1 := gk.GetEthereumEventVoteRecord(ctx, cctxe.GetEventNonce(), cctxe.Hash())
	require.NotNil(t, stored1)

	// var storedEvent, storedEvent1 types.EthereumEvent
	storedEvent, err := types.UnpackEvent(stored.Event)
	require.NoError(t, err)
	storedEvent1, err := types.UnpackEvent(stored1.Event)
	require.NoError(t, err)

	require.EqualValues(t, storedEvent.GetEventNonce(), 1)
	require.EqualValues(t, storedEvent.Hash(), stce.Hash())

	require.EqualValues(t, storedEvent1.GetEventNonce(), 2)
	require.EqualValues(t, storedEvent1.Hash(), cctxe.Hash())

	mapping := gk.GetEthereumEventVoteRecordMapping(ctx)
	require.EqualValues(t, 3, len(mapping[1][0].Votes))
	require.EqualValues(t, 3, len(mapping[2][0].Votes))

	eve1, err := types.UnpackEvent(mapping[1][0].Event)
	require.NoError(t, err)
	eve2, err := types.UnpackEvent(mapping[2][0].Event)
	require.NoError(t, err)
	require.EqualValues(t, 1, eve1.GetEventNonce())
	require.EqualValues(t, 2, eve2.GetEventNonce())
	require.EqualValues(t, stce.Hash(), eve1.Hash())
	require.EqualValues(t, cctxe.Hash(), eve2.Hash())
}

// TODO: uncomment
func TestLastSlashedValsetNonce(t *testing.T) {
	input := CreateTestEnv(t)
	k := input.GravityKeeper
	ctx := input.Context

	i := 1
	for ; i < 10; i++ {
		ctx = ctx.WithBlockHeight(int64(i))
		_ = k.CreateSignerSetTx(ctx)
	}

	latestValsetNonce := k.GetLatestSignerSetTxNonce(ctx)
	assert.Equal(t, uint64(i-1), latestValsetNonce)

	//  lastSlashedValsetNonce should be zero initially.
	lastSlashedValsetNonce := k.GetLastSlashedOutgoingTxBlockHeight(ctx)
	assert.Equal(t, uint64(0), lastSlashedValsetNonce)
	unslashedValsets := k.GetUnSlashedOutgoingTxs(ctx, uint64(12))
	assert.Equal(t, 9, len(unslashedValsets))

	// check if last Slashed Valset nonce is set properly or not
	k.SetLastSlashedOutgoingTxBlockHeight(ctx, uint64(3))
	lastSlashedValsetNonce = k.GetLastSlashedOutgoingTxBlockHeight(ctx)
	assert.Equal(t, uint64(3), lastSlashedValsetNonce)

	// when maxHeight < lastSlashedValsetNonce, len(unslashedValsets) should be zero
	unslashedValsets = k.GetUnSlashedOutgoingTxs(ctx, uint64(2))
	assert.Equal(t, 0, len(unslashedValsets))

	// when maxHeight == lastSlashedValsetNonce, len(unslashedValsets) should be zero
	unslashedValsets = k.GetUnSlashedOutgoingTxs(ctx, uint64(3))
	assert.Equal(t, 0, len(unslashedValsets))

	// when maxHeight > lastSlashedValsetNonce && maxHeight <= latestValsetNonce
	unslashedValsets = k.GetUnSlashedOutgoingTxs(ctx, uint64(6))
	assert.Equal(t, 2, len(unslashedValsets))

	// when maxHeight > latestValsetNonce
	unslashedValsets = k.GetUnSlashedOutgoingTxs(ctx, uint64(15))
	assert.Equal(t, 6, len(unslashedValsets))
}
