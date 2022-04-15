package keeper

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	accType "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

// TODO: Add any future invariants here
// TODO: (see the sdk docs for more info https://docs.cosmos.network/master/building-modules/invariants.html)
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return ModuleBalanceInvariant(k)(ctx)

		// Example additional invariants
		//  res, stop := FutureInvariant(k)(ctx)
		//	if stop {
		//		return res, stop
		//	}
		//
		//	return AnotherFutureInvariant(k)(ctx)
	}
}

// ModuleBalanceInvariant checks that the module account's balance is equal to the balance of unbatched transactions and unobserved batches
// Note that the returned bool should be true if there is an error, e.g. an unexpected module balance
func ModuleBalanceInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		modAcc := accType.NewModuleAddress(types.ModuleName)
		actualBals := k.bankKeeper.GetAllBalances(ctx, modAcc)
		expectedBals := make(map[string]*sdk.Int, len(actualBals)) // Collect balances by contract
		for _, v := range actualBals {
			newInt := sdk.NewInt(0)
			expectedBals[v.Denom] = &newInt
		}
		expectedBals = sumUnconfirmedBatchModuleBalances(ctx, k, expectedBals)
		expectedBals = sumUnbatchedSendToEthereumsModuleBalances(ctx, k, expectedBals)

		// Compare actual vs expected balances
		for _, actual := range actualBals {
			denom := actual.GetDenom()
			cosmosOriginated, _, err := k.DenomToERC20Lookup(ctx, denom)
			if err != nil {
				// Here we do not return because a user could halt the chain by gifting gravity a cosmos asset with no erc20 repr
				ctx.Logger().Error("Unexpected gravity module balance of cosmos-originated asset with no erc20 representation", "asset", denom)
				continue
			}
			expected, ok := expectedBals[denom]
			if !ok {
				return fmt.Sprint("Could not find expected balance for actual module balance of ", actual), true
			}

			if cosmosOriginated { // Cosmos originated mismatched balance
				// We cannot make any assertions about cosmosOriginated assets because we do not have enough information.
				// There is no index of denom => amount bridged, which would force us to parse all logs in existence
			} else if !actual.Amount.Equal(*expected) { // Eth originated mismatched balance
				return fmt.Sprint("Mismatched balance of eth-originated ", denom, ": actual balance ", actual.Amount, " != expected balance ", expected), true
			}
		}
		return "", false
	}
}

// sumUnconfirmedBatchModuleBalances calculate the value the module should have stored due to unconfirmed batches
func sumUnconfirmedBatchModuleBalances(ctx sdk.Context, k Keeper, expectedBals map[string]*sdk.Int) map[string]*sdk.Int {
	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		batchTotal := sdk.NewInt(0)
		// Collect the send amount + fee amount for each tx
		batch, _ := otx.(*types.BatchTx)
		for _, tx := range batch.Transactions {
			newTotal := batchTotal.Add(tx.Erc20Token.Amount.Add(tx.Erc20Fee.Amount))
			batchTotal = newTotal
		}
		contract := batch.TokenContract
		_, denom := k.ERC20ToDenomLookup(ctx, common.HexToAddress(contract))
		// Add the batch total to the contract counter
		_, ok := expectedBals[denom]
		if !ok {
			zero := sdk.ZeroInt()
			expectedBals[denom] = &zero
		}

		*expectedBals[denom] = expectedBals[denom].Add(batchTotal)

		return false // continue iterating
	})

	return expectedBals
}

// sumUnbatchedSendToEthereumsModuleBalances calculates the value the module should have stored due to unbatched txs
func sumUnbatchedSendToEthereumsModuleBalances(ctx sdk.Context, k Keeper, expectedBals map[string]*sdk.Int) map[string]*sdk.Int {
	// It is also given the balance of all unbatched txs in the pool
	k.IterateUnbatchedSendToEthereums(ctx, func(ste *types.SendToEthereum) bool {
		contract := ste.Erc20Token.Contract
		_, denom := k.ERC20ToDenomLookup(ctx, common.HexToAddress(contract))

		// Collect the send amount + fee amount for each tx
		txTotal := ste.Erc20Token.Amount.Add(ste.Erc20Fee.Amount)
		_, ok := expectedBals[denom]
		if !ok {
			zero := sdk.ZeroInt()
			expectedBals[denom] = &zero
		}
		*expectedBals[denom] = expectedBals[denom].Add(txTotal)

		return false // continue iterating
	})

	return expectedBals
}
