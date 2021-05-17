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

func TestPrefixRange(t *testing.T) {
	cases := map[string]struct {
		src      []byte
		expStart []byte
		expEnd   []byte
		expPanic bool
	}{
		"normal":              {src: []byte{1, 3, 4}, expStart: []byte{1, 3, 4}, expEnd: []byte{1, 3, 5}},
		"normal short":        {src: []byte{79}, expStart: []byte{79}, expEnd: []byte{80}},
		"empty case":          {src: []byte{}},
		"roll-over example 1": {src: []byte{17, 28, 255}, expStart: []byte{17, 28, 255}, expEnd: []byte{17, 29, 0}},
		"roll-over example 2": {src: []byte{15, 42, 255, 255},
			expStart: []byte{15, 42, 255, 255}, expEnd: []byte{15, 43, 0, 0}},
		"pathological roll-over": {src: []byte{255, 255, 255, 255}, expStart: []byte{255, 255, 255, 255}},
		"nil prohibited":         {expPanic: true},
	}

	for testName, tc := range cases {
		tc := tc
		t.Run(testName, func(t *testing.T) {
			if tc.expPanic {
				require.Panics(t, func() {
					prefixRange(tc.src)
				})
				return
			}
			start, end := prefixRange(tc.src)
			assert.Equal(t, tc.expStart, start)
			assert.Equal(t, tc.expEnd, end)
		})
	}
}

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
			r := input.GravityKeeper.NewSignerSetTx(ctx)
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
	input.GravityKeeper.SetEthereumEventVoteRecord(ctx, dep1.EventNonce, dep1.Hash(), att1)
	input.GravityKeeper.SetEthereumEventVoteRecord(ctx, dep2.EventNonce, dep2.Hash(), att2)

	atts := []types.EthereumEventVoteRecord{}
	input.GravityKeeper.IterateEthereumEventVoteRecords(ctx, func(_ []byte, att types.EthereumEventVoteRecord) bool {
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
		k.SetValidatorEthereumAddress(ctx, val, ethAddrs[i])
		k.SetEthereumOrchestratorAddress(ctx, ethAddrs[i], orch)
	}
	addresses := k.GetDelegateKeys(ctx)
	for i := range addresses {
		res := addresses[i]
		assert.Equal(t, valAddrs[i], res.ValidatorAddress)
		assert.Equal(t, orchAddrs[i], res.OrchestratorAddress)
		assert.Equal(t, ethAddrs[i], res.EthereumAddress)
	}
}

func TestStoreEventVoteRecord(t *testing.T) {
	input := CreateTestEnv(t)
	event := &types.SendToCosmosEvent{
		EventNonce:     1,
		TokenContract:  EthAddrs[0].Hex(),
		EthereumSender: EthAddrs[0].Hex(),
		CosmosReceiver: AccAddrs[0].String(),
		EthereumHeight: 10,
		Amount:         sdk.NewInt(1000000),
	}
	eva, err := types.PackEvent(event)
	require.NoError(t, err)

	evr := &types.EthereumEventVoteRecord{
		Event: eva,
		Votes: []string{
			ValAddrs[0].String(),
			ValAddrs[1].String(),
			ValAddrs[2].String(),
		},
		Accepted: false,
	}

	input.GravityKeeper.SetEthereumEventVoteRecord(input.Context, event.GetEventNonce(), event.Hash(), evr)

	stored := input.GravityKeeper.GetEthereumEventVoteRecord(input.Context, event.GetEventNonce(), event.Hash())
	require.NotNil(t, stored)

	var storedEvent types.EthereumEvent
	require.NoError(t, input.GravityKeeper.cdc.UnpackAny(stored.Event, &storedEvent))

	require.EqualValues(t, storedEvent.GetNonce(), 1)
	require.EqualValues(t, storedEvent.Hash(), event.Hash())
}

// TODO: uncomment
// func TestLastSlashedValsetNonce(t *testing.T) {
// 	input := CreateTestEnv(t)
// 	k := input.GravityKeeper
// 	ctx := input.Context

// 	vs := k.NewSignerSetTx(ctx)

// 	i := 1
// 	for ; i < 10; i++ {
// 		vs.Nonce = uint64(i)
// 		k.SetOutgoingTx(ctx, vs)
// 	}

// 	latestValsetNonce := k.GetLatestSignerSetTxNonce(ctx)
// 	assert.Equal(t, latestValsetNonce, uint64(i-1))

// 	//  lastSlashedValsetNonce should be zero initially.
// 	lastSlashedValsetNonce := k.GetLastSlashedValsetNonce(ctx)
// 	assert.Equal(t, lastSlashedValsetNonce, uint64(0))
// 	unslashedValsets := k.GetUnSlashedSignerSetTxs(ctx, uint64(12))
// 	assert.Equal(t, len(unslashedValsets), 9)

// 	// check if last Slashed Valset nonce is set properly or not
// 	k.SetLastSlashedValsetNonce(ctx, uint64(3))
// 	lastSlashedValsetNonce = k.GetLastSlashedValsetNonce(ctx)
// 	assert.Equal(t, lastSlashedValsetNonce, uint64(3))

// 	// when maxHeight < lastSlashedValsetNonce, len(unslashedValsets) should be zero
// 	unslashedValsets = k.GetUnSlashedValsets(ctx, uint64(2))
// 	assert.Equal(t, len(unslashedValsets), 0)

// 	// when maxHeight == lastSlashedValsetNonce, len(unslashedValsets) should be zero
// 	unslashedValsets = k.GetUnSlashedValsets(ctx, uint64(3))
// 	assert.Equal(t, len(unslashedValsets), 0)

// 	// when maxHeight > lastSlashedValsetNonce && maxHeight <= latestValsetNonce
// 	unslashedValsets = k.GetUnSlashedValsets(ctx, uint64(6))
// 	assert.Equal(t, len(unslashedValsets), 2)

// 	// when maxHeight > latestValsetNonce
// 	unslashedValsets = k.GetUnSlashedValsets(ctx, uint64(15))
// 	assert.Equal(t, len(unslashedValsets), 6)
// 	fmt.Println("unslashedValsetsRange", unslashedValsets)
// }
