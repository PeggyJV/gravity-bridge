package keeper

import (
	"bytes"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
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
				cAddr := bytes.Repeat([]byte{byte(i)}, 20)
				operators[i] = MockStakingValidatorData{
					// any unique addr
					Operator: cAddr,
					Power:    int64(v),
				}
				input.GravityKeeper.setValidatorEthereumAddress(ctx, cAddr, common.HexToAddress("0xf71402f886b45c134743F4c00750823Bbf5Fd045"))
			}
			input.GravityKeeper.StakingKeeper = NewStakingKeeperWeightedMock(operators...)
			r := input.GravityKeeper.CreateSignerSetTx(ctx)
			assert.Equal(t, spec.expPowers, r.Signers.GetPowers())
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

	var atts []*types.EthereumEventVoteRecord
	input.GravityKeeper.IterateEthereumEventVoteRecords(ctx, func(_ []byte, att *types.EthereumEventVoteRecord) bool {
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

		assert.False(t, k.validatorForEthAddressExists(ctx, ethAddrs[i]))
		assert.False(t, k.ethAddressForOrchestratorExists(ctx, orch))

		k.SetOrchestratorValidatorAddress(ctx, val, orch)
		k.setValidatorEthereumAddress(ctx, val, ethAddrs[i])
		k.setEthereumOrchestratorAddress(ctx, ethAddrs[i], orch)

		assert.True(t, k.validatorForEthAddressExists(ctx, ethAddrs[i]))
		assert.True(t, k.ethAddressForOrchestratorExists(ctx, orch))
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
		InvalidationScope: []byte{0x1, 0x2},
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
	unslashedValsets := k.GetUnslashedOutgoingTxs(ctx, uint64(12))
	assert.Equal(t, 9, len(unslashedValsets))

	// check if last Slashed Valset nonce is set properly or not
	k.SetLastSlashedOutgoingTxBlockHeight(ctx, uint64(3))
	lastSlashedValsetNonce = k.GetLastSlashedOutgoingTxBlockHeight(ctx)
	assert.Equal(t, uint64(3), lastSlashedValsetNonce)

	// when maxHeight < lastSlashedValsetNonce, len(unslashedValsets) should be zero
	unslashedValsets = k.GetUnslashedOutgoingTxs(ctx, uint64(2))
	assert.Equal(t, 0, len(unslashedValsets))

	// when maxHeight == lastSlashedValsetNonce, len(unslashedValsets) should be zero
	unslashedValsets = k.GetUnslashedOutgoingTxs(ctx, uint64(3))
	assert.Equal(t, 0, len(unslashedValsets))

	// when maxHeight > lastSlashedValsetNonce && maxHeight <= latestValsetNonce
	unslashedValsets = k.GetUnslashedOutgoingTxs(ctx, uint64(6))
	assert.Equal(t, 2, len(unslashedValsets))

	// when maxHeight > latestValsetNonce
	unslashedValsets = k.GetUnslashedOutgoingTxs(ctx, uint64(15))
	assert.Equal(t, 6, len(unslashedValsets))
}

// ---

func TestKeeper_GetLatestSignerSetTx(t *testing.T) {
	t.Run("read before there's one in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		got := gk.GetLatestSignerSetTx(ctx)
		require.Nil(t, got)
	})

	t.Run("read after there's one in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		{ // setup
			gk.SetOutgoingTx(ctx, &types.SignerSetTx{
				Nonce:   gk.incrementLatestSignerSetTxNonce(ctx),
				Height:  1,
				Signers: nil,
			})
		}

		{ // validate
			got := gk.GetLatestSignerSetTx(env.Context)
			require.NotNil(t, got)
			require.EqualValues(t, got.Height, gk.GetLatestSignerSetTxNonce(ctx))
		}
	})
}

func TestKeeper_GetSignerSetTxs(t *testing.T) {
	t.Run("read before there's any in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		got := gk.GetSignerSetTxs(ctx)
		require.Nil(t, got)
	})

	t.Run("read after there's one in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		{ // setup
			gk.SetOutgoingTx(ctx, &types.SignerSetTx{
				Nonce:   gk.incrementLatestSignerSetTxNonce(ctx),
				Height:  1,
				Signers: nil,
			})
		}

		{ // validate
			got := gk.GetSignerSetTxs(ctx)
			require.NotNil(t, got)
			require.Len(t, got, 1)
		}
	})
}

func TestKeeper_GetLastObservedSignerSetTx(t *testing.T) {
	t.Run("read before there's any in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		got := gk.GetLastObservedSignerSetTx(ctx)
		require.Nil(t, got)
	})

	t.Run("read after there's one in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		{ // setup
			gk.setLastObservedSignerSetTx(ctx, types.SignerSetTx{
				Nonce:   1,
				Height:  1,
				Signers: nil,
			})
		}

		{ // validate
			got := gk.GetLastObservedSignerSetTx(ctx)
			require.NotNil(t, got)
		}
	})
}

func TestKeeper_GetLastUnbondingBlockHeight(t *testing.T) {
	t.Run("read before there's any in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		got := gk.GetLastUnbondingBlockHeight(ctx)
		require.Zero(t, got)
	})

	t.Run("read after there's one in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		{ // setup
			gk.setLastUnbondingBlockHeight(ctx, 10)
		}

		{ // validate
			got := gk.GetLastUnbondingBlockHeight(ctx)
			require.EqualValues(t, 10, got)
		}
	})
}

func TestKeeper_GetEthereumSignatures(t *testing.T) {
	t.Run("read before there's anything in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		storeIndexes := [][]byte{
			types.MakeSignerSetTxKey(1),
			types.MakeBatchTxKey(common.HexToAddress(""), 1), // weird that empty address is okay
			types.MakeContractCallTxKey(nil, 0),
		}
		for _, storeIndex := range storeIndexes {
			got := gk.GetEthereumSignatures(ctx, storeIndex)
			require.Empty(t, got)
		}
	})

	t.Run("read after there's one signer-set-tx-confirmation in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		ethAddr := common.HexToAddress("0x3146D2d6Eed46Afa423969f5dDC3152DfC359b09")
		valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4")
		require.NoError(t, err)

		const signerSetNonce = 10

		{ // setup
			signerSetTxConfirmation := &types.SignerSetTxConfirmation{
				SignerSetNonce: signerSetNonce,
				EthereumSigner: ethAddr.Hex(),
				Signature:      []byte("fake-signature"),
			}
			key := gk.SetEthereumSignature(ctx, signerSetTxConfirmation, valAddr)
			require.NotEmpty(t, key)
		}

		{ // validate
			storeIndex := types.MakeSignerSetTxKey(signerSetNonce)

			{ // getEthereumSignature
				got := gk.getEthereumSignature(ctx, storeIndex, valAddr)
				require.Equal(t, []byte("fake-signature"), got)
			}
			{ // GetEthereumSignatures
				got := gk.GetEthereumSignatures(ctx, storeIndex)
				require.Len(t, got, 1)
			}
		}
	})

	t.Run("read after there's one batch-tx-confirmation in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		ethAddr := common.HexToAddress("0x3146D2d6Eed46Afa423969f5dDC3152DfC359b09")
		valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4")
		require.NoError(t, err)

		const (
			batchNonce    = 10
			tokenContract = "0x1111111111111111111111111111111111111111"
		)

		{ // setup
			batchTxConfirmation := &types.BatchTxConfirmation{
				TokenContract:  tokenContract,
				BatchNonce:     batchNonce,
				EthereumSigner: ethAddr.Hex(),
				Signature:      []byte("fake-signature"),
			}
			key := gk.SetEthereumSignature(ctx, batchTxConfirmation, valAddr)
			require.NotEmpty(t, key)
		}

		{ // validate
			storeIndex := types.MakeBatchTxKey(common.HexToAddress(tokenContract), batchNonce)

			{ // getEthereumSignature
				got := gk.getEthereumSignature(ctx, storeIndex, valAddr)
				require.Equal(t, []byte("fake-signature"), got)
			}
			{ // GetEthereumSignatures
				got := gk.GetEthereumSignatures(ctx, storeIndex)
				require.Len(t, got, 1)
			}
		}
	})

	t.Run("read after there's one contract-call-tx-confirmation in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		ethAddr := common.HexToAddress("0x3146D2d6Eed46Afa423969f5dDC3152DfC359b09")

		valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4")
		require.NoError(t, err)

		const (
			invalidationScope = "some-invalidation-scope"
			invalidationNonce = 10
		)

		{ // setup
			contractCallConfirmation := &types.ContractCallTxConfirmation{
				InvalidationScope: []byte(invalidationScope),
				InvalidationNonce: invalidationNonce,
				EthereumSigner:    ethAddr.Hex(),
				Signature:         []byte("fake-signature"),
			}
			key := gk.SetEthereumSignature(ctx, contractCallConfirmation, valAddr)
			require.NotEmpty(t, key)
		}

		{ // validate
			storeIndex := types.MakeContractCallTxKey([]byte(invalidationScope), invalidationNonce)

			{ // getEthereumSignature
				got := gk.getEthereumSignature(ctx, storeIndex, valAddr)
				require.Equal(t, []byte("fake-signature"), got)
			}
			{ // GetEthereumSignatures
				got := gk.GetEthereumSignatures(ctx, storeIndex)
				require.Len(t, got, 1)
			}
		}
	})
}

func TestKeeper_Migration(t *testing.T) {

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
		InvalidationScope: []byte{0x1, 0x2},
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

	//Put an outgoing transaction into the system

	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = common.HexToAddress("0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7")
		myTokenContractAddr = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5") // Pickle
		allVouchers         = sdk.NewCoins(
			types.NewERC20Token(99999, myTokenContractAddr).GravityCoin(),
		)
	)

	// mint some voucher first
	require.NoError(t, input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers))
	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	require.NoError(t, fundAccount(ctx, input.BankKeeper, mySender, allVouchers))

	// add some TX to the pool
	input.AddSendToEthTxsToPool(t, ctx, myTokenContractAddr, mySender, myReceiver, 2, 3, 2, 1)

	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	firstBatch := input.GravityKeeper.CreateBatchTx(ctx, myTokenContractAddr, 2)

	// then batch is persisted
	gotFirstBatch := input.GravityKeeper.GetOutgoingTx(ctx, firstBatch.GetStoreIndex())
	require.NotNil(t, gotFirstBatch)

	gk.setEthereumEventVoteRecord(ctx, stce.GetEventNonce(), stce.Hash(), evr)
	gk.setLastObservedEventNonce(ctx, stce.GetEventNonce())
	gk.setEthereumEventVoteRecord(ctx, cctxe.GetEventNonce(), cctxe.Hash(), evr2)
	gk.setLastObservedEventNonce(ctx, cctxe.GetEventNonce())

	stored := gk.GetEthereumEventVoteRecord(ctx, stce.GetEventNonce(), stce.Hash())
	require.NotNil(t, stored)

	stored2 := gk.GetEthereumEventVoteRecord(ctx, cctxe.GetEventNonce(), cctxe.Hash())
	require.NotNil(t, stored2)

	ethAddr := common.HexToAddress("0x3146D2d6Eed46Afa423969f5dDC3152DfC359b09")

	valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4")
	require.NoError(t, err)

	{ // setup
		batchTxConfirmation := &types.BatchTxConfirmation{
			TokenContract:  myTokenContractAddr.Hex(),
			BatchNonce:     firstBatch.BatchNonce,
			EthereumSigner: ethAddr.Hex(),
			Signature:      []byte("fake-signature"),
		}
		key := gk.SetEthereumSignature(ctx, batchTxConfirmation, valAddr)
		require.NotEmpty(t, key)
	}

	{ // validate
		storeIndex := gotFirstBatch.GetStoreIndex()

		{ // getEthereumSignature
			got := gk.getEthereumSignature(ctx, storeIndex, valAddr)
			require.Equal(t, []byte("fake-signature"), got)
		}
		{ // GetEthereumSignatures
			got := gk.GetEthereumSignatures(ctx, storeIndex)
			require.Len(t, got, 1)
		}
	}

	nonce := gk.GetLastObservedEventNonce(ctx)
	require.Equal(t, cctxe.GetEventNonce(), nonce)

	gk.setLastObservedSignerSetTx(ctx, types.SignerSetTx{
		Nonce:   1,
		Height:  1,
		Signers: nil,
	})

	for _, val := range ValAddrs {
		gk.setLastEventNonceByValidator(ctx, val, nonce)
	}

	gk.MigrateGravityContract(ctx, "0x5e175bE4d23Fa25604CE7848F60FB340894D5CDA", 1000)

	storedAfterMigrate := gk.GetEthereumEventVoteRecord(ctx, stce.GetEventNonce(), stce.Hash())
	require.Nil(t, storedAfterMigrate)

	stored2AfterMigrate := gk.GetEthereumEventVoteRecord(ctx, cctxe.GetEventNonce(), cctxe.Hash())
	require.Nil(t, stored2AfterMigrate)

	nonce2 := gk.GetLastObservedEventNonce(ctx)
	require.Equal(t, uint64(0), nonce2)

	for _, val := range ValAddrs {
		require.Equal(t, uint64(0), gk.getLastEventNonceByValidator(ctx, val))
	}

	got := gk.GetLastObservedSignerSetTx(ctx)
	require.Equal(t, got, &types.SignerSetTx{Nonce: 0x0, Height: 0x0, Signers: types.EthereumSigners(nil)})

	{ // GetEthereumSignatures
		storeIndex := gotFirstBatch.GetStoreIndex()
		got := gk.GetEthereumSignatures(ctx, storeIndex)
		require.Len(t, got, 0)
	}

}

func TestEthereumSignatureIterators(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	k := input.GravityKeeper

	// add some signatures to the store
	valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4")
	require.NoError(t, err)
	signer := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	k.setValidatorEthereumAddress(ctx, valAddr, signer)
	b1 := &types.BatchTxConfirmation{
		TokenContract:  "0x1111111111111111111111111111111111111111",
		BatchNonce:     1,
		EthereumSigner: signer.Hex(),
		Signature:      []byte("batch-signature-1"),
	}
	b2 := &types.BatchTxConfirmation{
		TokenContract:  "0x2222222222222222222222222222222222222222",
		BatchNonce:     2,
		EthereumSigner: signer.Hex(),
		Signature:      []byte("batch-signature-2"),
	}
	k.SetEthereumSignature(ctx, b1, valAddr)
	k.SetEthereumSignature(ctx, b2, valAddr)

	iterationCount := 0
	var batchSigs []*types.BatchTxConfirmation
	k.IterateBatchTxEthereumSignatures(ctx, func(contractAddr common.Address, nonce uint64, val sdk.ValAddress, sig []byte) bool {
		iterationCount++
		batchSigs = append(batchSigs, &types.BatchTxConfirmation{
			TokenContract:  contractAddr.Hex(),
			BatchNonce:     nonce,
			EthereumSigner: k.GetValidatorEthereumAddress(ctx, valAddr).Hex(),
			Signature:      sig,
		})
		return false
	})

	require.Len(t, batchSigs, 2)
	require.Equal(t, iterationCount, len(batchSigs))
	require.Equal(t, batchSigs[0], b1)
	require.Equal(t, batchSigs[1], b2)

	// ContractCallTxConfirmations

	scope := crypto.Keccak256Hash([]byte("test-scope")).Bytes()
	cc1 := &types.ContractCallTxConfirmation{
		InvalidationScope: scope,
		InvalidationNonce: 1,
		EthereumSigner:    signer.Hex(),
		Signature:         []byte("contract-call-signature-1"),
	}
	cc2 := &types.ContractCallTxConfirmation{
		InvalidationScope: scope,
		InvalidationNonce: 2,
		EthereumSigner:    signer.Hex(),
		Signature:         []byte("contract-call-signature-2"),
	}
	k.SetEthereumSignature(ctx, cc1, valAddr)
	k.SetEthereumSignature(ctx, cc2, valAddr)

	iterationCount = 0
	var ccSigs []*types.ContractCallTxConfirmation
	k.IterateContractCallTxEthereumSignatures(ctx, func(invalidationScope []byte, invalidationNonce uint64, val sdk.ValAddress, sig []byte) bool {
		iterationCount++
		ccSigs = append(ccSigs, &types.ContractCallTxConfirmation{
			InvalidationScope: invalidationScope,
			InvalidationNonce: invalidationNonce,
			EthereumSigner:    k.GetValidatorEthereumAddress(ctx, valAddr).Hex(),
			Signature:         sig,
		})
		return false
	})

	require.Len(t, ccSigs, 2)
	require.Equal(t, iterationCount, len(ccSigs))
	require.Equal(t, ccSigs[0], cc1)
	require.Equal(t, ccSigs[1], cc2)

	// SignerSetTxConfirmations

	ss1 := &types.SignerSetTxConfirmation{
		SignerSetNonce: 1,
		EthereumSigner: signer.Hex(),
		Signature:      []byte("signer-set-signature-1"),
	}
	ss2 := &types.SignerSetTxConfirmation{
		SignerSetNonce: 2,
		EthereumSigner: signer.Hex(),
		Signature:      []byte("signer-set-signature-2"),
	}
	k.SetEthereumSignature(ctx, ss1, valAddr)
	k.SetEthereumSignature(ctx, ss2, valAddr)

	iterationCount = 0
	var ssSigs []*types.SignerSetTxConfirmation
	k.IterateSignerSetTxEthereumSignatures(ctx, func(nonce uint64, val sdk.ValAddress, sig []byte) bool {
		iterationCount++
		ssSigs = append(ssSigs, &types.SignerSetTxConfirmation{
			SignerSetNonce: nonce,
			EthereumSigner: k.GetValidatorEthereumAddress(ctx, valAddr).Hex(),
			Signature:      sig,
		})
		return false
	})

	require.Len(t, ssSigs, 2)
	require.Equal(t, iterationCount, len(ssSigs))
	require.Equal(t, ssSigs[0], ss1)
	require.Equal(t, ssSigs[1], ss2)
}

// TODO(levi) review/ensure coverage for:
// PaginateOutgoingTxsByType
// GetUnbondingvalidators(unbondingVals []byte) stakingtypes.ValAddresses
