package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
)

func TestBatches(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
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

	// CREATE FIRST BATCH
	// ==================

	// add some TX to the pool
	input.AddSendToEthTxsToPool(t, ctx, myTokenContractAddr, mySender, myReceiver, 2, 3, 2, 1)

	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	firstBatch := input.GravityKeeper.CreateBatchTx(ctx, myTokenContractAddr, 2)

	// then batch is persisted
	gotFirstBatch := input.GravityKeeper.GetOutgoingTx(ctx, firstBatch.GetStoreIndex())
	require.NotNil(t, gotFirstBatch)

	gfb := gotFirstBatch.(*types.BatchTx)
	expFirstBatch := &types.BatchTx{
		BatchNonce: 1,
		Transactions: []*types.SendToEthereum{
			types.NewSendToEthereumTx(2, myTokenContractAddr, mySender, myReceiver, 101, 3),
			types.NewSendToEthereumTx(3, myTokenContractAddr, mySender, myReceiver, 102, 2),
		},
		TokenContract: myTokenContractAddr.Hex(),
		Height:        1234567,
	}

	assert.Equal(t, expFirstBatch.Transactions, gfb.Transactions)

	// and verify remaining available Tx in the pool
	var gotUnbatchedTx []*types.SendToEthereum
	input.GravityKeeper.IterateUnbatchedSendToEthereums(ctx, func(tx *types.SendToEthereum) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx := []*types.SendToEthereum{
		types.NewSendToEthereumTx(1, myTokenContractAddr, mySender, myReceiver, 100, 2),
		types.NewSendToEthereumTx(4, myTokenContractAddr, mySender, myReceiver, 103, 1),
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)

	// CREATE SECOND, MORE PROFITABLE BATCH
	// ====================================

	// add some more TX to the pool to create a more profitable batch
	input.AddSendToEthTxsToPool(t, ctx, myTokenContractAddr, mySender, myReceiver, 4, 5)

	// create the more profitable batch
	ctx = ctx.WithBlockTime(now)
	// tx batch size is 2, so that some of them stay behind
	secondBatch := input.GravityKeeper.CreateBatchTx(ctx, myTokenContractAddr, 2)

	// check that the more profitable batch has the right txs in it
	expSecondBatch := &types.BatchTx{
		BatchNonce: 2,
		Transactions: []*types.SendToEthereum{
			types.NewSendToEthereumTx(6, myTokenContractAddr, mySender, myReceiver, 101, 5),
			types.NewSendToEthereumTx(5, myTokenContractAddr, mySender, myReceiver, 100, 4),
		},
		TokenContract: myTokenContractAddr.Hex(),
		Height:        1234567,
	}

	assert.Equal(t, expSecondBatch, secondBatch)

	// EXECUTE THE MORE PROFITABLE BATCH
	// =================================

	// Execute the batch
	input.GravityKeeper.batchTxExecuted(ctx, common.HexToAddress(secondBatch.TokenContract), secondBatch.BatchNonce)

	// check batch has been deleted
	gotSecondBatch := input.GravityKeeper.GetOutgoingTx(ctx, secondBatch.GetStoreIndex())
	require.Nil(t, gotSecondBatch)

	// check that txs from first batch have been freed
	gotUnbatchedTx = nil
	input.GravityKeeper.IterateUnbatchedSendToEthereums(ctx, func(tx *types.SendToEthereum) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx = []*types.SendToEthereum{
		types.NewSendToEthereumTx(2, myTokenContractAddr, mySender, myReceiver, 101, 3),
		types.NewSendToEthereumTx(3, myTokenContractAddr, mySender, myReceiver, 102, 2),
		types.NewSendToEthereumTx(1, myTokenContractAddr, mySender, myReceiver, 100, 2),
		types.NewSendToEthereumTx(4, myTokenContractAddr, mySender, myReceiver, 103, 1),
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)
}

// tests that batches work with large token amounts, mostly a duplicate of the above
// tests but using much bigger numbers
func TestBatchesFullCoins(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = common.HexToAddress("0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7")
		myTokenContractAddr = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5") // Pickle
		totalCoins, _       = sdk.NewIntFromString("1500000000000000000000")                    // 1,500 ETH worth
		oneEth, _           = sdk.NewIntFromString("1000000000000000000")
		allVouchers         = sdk.NewCoins(
			types.NewSDKIntERC20Token(totalCoins, myTokenContractAddr).GravityCoin(),
		)
	)

	// mint some voucher first
	require.NoError(t, input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers))
	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	require.NoError(t, fundAccount(ctx, input.BankKeeper, mySender, allVouchers))

	// CREATE FIRST BATCH
	// ==================

	// add some TX to the pool
	for _, v := range []uint64{20, 300, 25, 10} {
		vAsSDKInt := sdk.NewIntFromUint64(v)
		amount := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).GravityCoin()
		fee := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).GravityCoin()
		_, err := input.GravityKeeper.createSendToEthereum(ctx, mySender, myReceiver.Hex(), amount, fee)
		require.NoError(t, err)
	}

	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	firstBatch := input.GravityKeeper.CreateBatchTx(ctx, myTokenContractAddr, 2)

	// then batch is persisted
	gotFirstBatch := input.GravityKeeper.GetOutgoingTx(ctx, firstBatch.GetStoreIndex())
	require.NotNil(t, gotFirstBatch)

	expFirstBatch := &types.BatchTx{
		BatchNonce: 1,
		Transactions: []*types.SendToEthereum{
			{
				Id:                2,
				Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(300)), myTokenContractAddr),
				Sender:            mySender.String(),
				EthereumRecipient: myReceiver.Hex(),
				Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(300)), myTokenContractAddr),
			},
			{
				Id:                3,
				Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(25)), myTokenContractAddr),
				Sender:            mySender.String(),
				EthereumRecipient: myReceiver.Hex(),
				Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(25)), myTokenContractAddr),
			},
		},
		TokenContract: myTokenContractAddr.Hex(),
		Height:        1234567,
	}
	assert.Equal(t, expFirstBatch, gotFirstBatch)

	// and verify remaining available Tx in the pool
	var gotUnbatchedTx []*types.SendToEthereum
	input.GravityKeeper.IterateUnbatchedSendToEthereums(ctx, func(tx *types.SendToEthereum) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx := []*types.SendToEthereum{
		{
			Id:                1,
			Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(20)), myTokenContractAddr),
			Sender:            mySender.String(),
			EthereumRecipient: myReceiver.Hex(),
			Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(20)), myTokenContractAddr),
		},
		{
			Id:                4,
			Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(10)), myTokenContractAddr),
			Sender:            mySender.String(),
			EthereumRecipient: myReceiver.Hex(),
			Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(10)), myTokenContractAddr),
		},
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)

	// CREATE SECOND, MORE PROFITABLE BATCH
	// ====================================

	// add some more TX to the pool to create a more profitable batch
	for _, v := range []uint64{4, 5} {
		vAsSDKInt := sdk.NewIntFromUint64(v)
		amount := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).GravityCoin()
		fee := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).GravityCoin()
		_, err := input.GravityKeeper.createSendToEthereum(ctx, mySender, myReceiver.Hex(), amount, fee)
		require.NoError(t, err)
	}

	// create the more profitable batch
	ctx = ctx.WithBlockTime(now)
	// tx batch size is 2, so that some of them stay behind
	secondBatch := input.GravityKeeper.CreateBatchTx(ctx, myTokenContractAddr, 2)

	// check that the more profitable batch has the right txs in it
	expSecondBatch := &types.BatchTx{
		BatchNonce: 2,
		Transactions: []*types.SendToEthereum{
			{
				Id:                1,
				Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(20)), myTokenContractAddr),
				Sender:            mySender.String(),
				EthereumRecipient: myReceiver.Hex(),
				Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(20)), myTokenContractAddr),
			},
			{
				Id:                4,
				Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(10)), myTokenContractAddr),
				Sender:            mySender.String(),
				EthereumRecipient: myReceiver.Hex(),
				Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(10)), myTokenContractAddr),
			},
		},
		TokenContract: myTokenContractAddr.Hex(),
		Height:        1234567,
	}

	assert.Equal(t, expSecondBatch, secondBatch)

	// EXECUTE THE MORE PROFITABLE BATCH
	// =================================

	// Execute the batch
	input.GravityKeeper.batchTxExecuted(ctx, common.HexToAddress(secondBatch.TokenContract), secondBatch.BatchNonce)

	// check batch has been deleted
	gotSecondBatch := input.GravityKeeper.GetOutgoingTx(ctx, secondBatch.GetStoreIndex())
	require.Nil(t, gotSecondBatch)

	// check that txs from first batch have been freed
	gotUnbatchedTx = nil
	input.GravityKeeper.IterateUnbatchedSendToEthereums(ctx, func(tx *types.SendToEthereum) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx = []*types.SendToEthereum{
		{
			Id:                2,
			Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(300)), myTokenContractAddr),
			Sender:            mySender.String(),
			EthereumRecipient: myReceiver.Hex(),
			Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(300)), myTokenContractAddr),
		},
		{
			Id:                3,
			Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(25)), myTokenContractAddr),
			Sender:            mySender.String(),
			EthereumRecipient: myReceiver.Hex(),
			Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(25)), myTokenContractAddr),
		},
		{
			Id:                6,
			Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(5)), myTokenContractAddr),
			Sender:            mySender.String(),
			EthereumRecipient: myReceiver.Hex(),
			Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(5)), myTokenContractAddr),
		},
		{
			Id:                5,
			Erc20Fee:          types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(4)), myTokenContractAddr),
			Sender:            mySender.String(),
			EthereumRecipient: myReceiver.Hex(),
			Erc20Token:        types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(4)), myTokenContractAddr),
		},
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)
}

func TestPoolTxRefund(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = common.HexToAddress("0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7")
		myTokenContractAddr = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5") // Pickle
		allVouchers         = sdk.NewCoins(
			types.NewERC20Token(414, myTokenContractAddr).GravityCoin(),
		)
		myDenom = types.NewERC20Token(1, myTokenContractAddr).GravityCoin().Denom
	)

	// mint some voucher first
	require.NoError(t, input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers))
	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	require.NoError(t, fundAccount(ctx, input.BankKeeper, mySender, allVouchers))

	// CREATE FIRST BATCH
	// ==================

	// add some TX to the pool
	// for i, v := range []uint64{2, 3, 2, 1} {
	// 	amount := types.NewERC20Token(uint64(i+100), myTokenContractAddr).GravityCoin()
	// 	fee := types.NewERC20Token(v, myTokenContractAddr).GravityCoin()
	// 	_, err := input.GravityKeeper.CreateSendToEthereum(ctx, mySender, myReceiver, amount, fee)
	// 	require.NoError(t, err)
	// }
	input.AddSendToEthTxsToPool(t, ctx, myTokenContractAddr, mySender, myReceiver, 2, 3, 2, 1)

	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	input.GravityKeeper.CreateBatchTx(ctx, myTokenContractAddr, 2)

	// try to refund a tx that's in a batch
	err := input.GravityKeeper.cancelSendToEthereum(ctx, 2, mySender.String())
	require.Error(t, err)

	// try to refund a tx that's in the pool
	err = input.GravityKeeper.cancelSendToEthereum(ctx, 4, mySender.String())
	require.NoError(t, err)

	// make sure refund was issued
	balances := input.BankKeeper.GetAllBalances(ctx, mySender)
	require.Equal(t, sdk.NewInt(104), balances.AmountOf(myDenom))
}

func TestEmptyBatch(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
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

	// no transactions should be included in this batch
	ctx = ctx.WithBlockTime(now)
	batchTx := input.GravityKeeper.CreateBatchTx(ctx, myTokenContractAddr, 2)

	require.Nil(t, batchTx)
}

func TestGetUnconfirmedBatchTxs(t *testing.T) {
	input, ctx := SetupFiveValChain(t)
	gk := input.GravityKeeper
	vals := input.StakingKeeper.GetAllValidators(ctx)
	val1, err := sdk.ValAddressFromBech32(vals[0].OperatorAddress)
	require.NoError(t, err)
	val2, err := sdk.ValAddressFromBech32(vals[1].OperatorAddress)
	require.NoError(t, err)

	blockheight := uint64(ctx.BlockHeight())
	sig := []byte("dummysig")
	gk.SetCompletedOutgoingTx(ctx, &types.BatchTx{
		BatchNonce: 1,
		Height:     uint64(ctx.BlockHeight()),
	})
	gk.SetOutgoingTx(ctx, &types.BatchTx{
		BatchNonce: 2,
		Height:     uint64(ctx.BlockHeight()),
	})

	// val1 signs both
	// val2 signs one
	gk.SetEthereumSignature(
		ctx,
		&types.BatchTxConfirmation{
			BatchNonce: 1,
			Signature:  sig,
		},
		val1,
	)
	gk.SetEthereumSignature(
		ctx,
		&types.BatchTxConfirmation{
			BatchNonce: 2,
			Signature:  sig,
		},
		val1,
	)
	gk.SetEthereumSignature(
		ctx,
		&types.BatchTxConfirmation{
			BatchNonce: 1,
			Signature:  sig,
		},
		val2,
	)

	require.Empty(t, gk.GetUnsignedBatchTxs(ctx, val1))
	require.Equal(t, 1, len(gk.GetUnsignedBatchTxs(ctx, val2)))

	// Confirm ordering
	val3, err := sdk.ValAddressFromBech32(vals[2].OperatorAddress)
	require.NoError(t, err)

	addressA := "0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	addressB := "0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
	gk.SetCompletedOutgoingTx(ctx, &types.BatchTx{
		TokenContract: addressB,
		BatchNonce:    3,
		Height:        blockheight,
	})
	gk.SetCompletedOutgoingTx(ctx, &types.BatchTx{
		TokenContract: addressA,
		BatchNonce:    4,
		Height:        blockheight,
	})
	gk.SetOutgoingTx(ctx, &types.BatchTx{
		TokenContract: addressB,
		BatchNonce:    5,
		Height:        blockheight,
	})
	gk.SetOutgoingTx(ctx, &types.BatchTx{
		TokenContract: addressA,
		BatchNonce:    6,
		Height:        blockheight,
	})
	gk.SetOutgoingTx(ctx, &types.BatchTx{
		TokenContract: addressB,
		BatchNonce:    7,
		Height:        blockheight,
	})

	unconfirmed := gk.GetUnsignedBatchTxs(ctx, val3)

	require.EqualValues(t, 7, len(unconfirmed))
	require.EqualValues(t, unconfirmed[0].BatchNonce, 1)
	require.EqualValues(t, unconfirmed[1].BatchNonce, 2)
	require.EqualValues(t, unconfirmed[2].BatchNonce, 3)
	require.EqualValues(t, unconfirmed[3].BatchNonce, 4)
	require.EqualValues(t, unconfirmed[4].BatchNonce, 5)
	require.EqualValues(t, unconfirmed[5].BatchNonce, 6)
	require.EqualValues(t, unconfirmed[6].BatchNonce, 7)
}
