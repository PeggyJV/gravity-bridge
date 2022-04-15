package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
	"github.com/stretchr/testify/require"
	"testing"
)

// Tests that the gravity module's balance is accounted for with unbatched txs, including tx cancellation
func TestModuleBalanceUnbatchedTxs(t *testing.T) {
	////////////////// SETUP //////////////////
	input := CreateTestEnv(t)
	defer func() { input.Context.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	ctx := input.Context
	var (
		mySender, _         = sdk.AccAddressFromBech32("gravity1ahx7f8wyertuus9r20284ej0asrs085ceqtfnm")
		myReceiver          = "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7"
		myTokenContractAddr = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5")
		myTokenDenom        = "gravity0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5"
	)

	input.GravityKeeper.setCosmosOriginatedDenomToERC20(ctx, myTokenDenom, myTokenContractAddr)
	// mint some voucher first
	voucher := sdk.NewCoin(myTokenDenom, sdk.NewInt(99999))
	allVouchers := sdk.Coins{voucher}
	err := input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers)
	require.NoError(t, err)
	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	err = input.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mySender, allVouchers)
	require.NoError(t, err)

	////////////////// EXECUTE //////////////////
	// Check the invariant without any transactions
	checkInvariant(t, ctx, input.GravityKeeper, true)

	// Create some unbatched transactions
	for i, v := range []uint64{2, 3, 2, 1} {
		input.GravityKeeper.createSendToEthereum(
			ctx,
			mySender,
			myReceiver,
			sdk.NewCoin(myTokenDenom, sdk.NewInt(int64(i+100))),
			sdk.NewCoin(myTokenDenom, sdk.NewIntFromUint64(v)))
		// Should create:
		// 1: amount 100, fee 2
		// 2: amount 101, fee 3
		// 3: amount 102, fee 2
		// 4: amount 103, fee 1
	}
	checkInvariant(t, ctx, input.GravityKeeper, true)

	// Remove one of the transactions
	err = input.GravityKeeper.cancelSendToEthereum(ctx, 1, mySender.String())
	require.NoError(t, err)
	checkInvariant(t, ctx, input.GravityKeeper, true)

	// Ensure an error is returned for a mismatched balance
	oneVoucher := sdk.NewCoin(myTokenDenom, sdk.NewInt(2))
	checkImbalancedModule(t, ctx, input.GravityKeeper, input.BankKeeper, mySender, sdk.NewCoins(oneVoucher))
}

func checkInvariant(t *testing.T, ctx sdk.Context, k Keeper, succeed bool) {
	res, ok := ModuleBalanceInvariant(k)(ctx)
	if succeed {
		require.False(t, ok, "Invariant should have returned false")
		require.Empty(t, res, "Invariant should have returned no message")
	} else {
		require.True(t, ok, "Invariant should have returned true")
		require.NotEmpty(t, res, "Invariant should have returned a message")
	}
}

func checkImbalancedModule(t *testing.T, ctx sdk.Context, gravityKeeper Keeper, bankKeeper bankkeeper.BaseKeeper, sender sdk.AccAddress, coins sdk.Coins) {
	// Imbalance the module
	bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins)
	checkInvariant(t, ctx, gravityKeeper, false)
	// Rebalance the module
	bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coins)
}
