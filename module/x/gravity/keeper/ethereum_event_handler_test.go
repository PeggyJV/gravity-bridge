package keeper

import (
	"math/big"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
)

func TestDetectMaliciousSupply(t *testing.T) {
	input := CreateTestEnv(t)

	// set supply to maximum value
	var testBigInt big.Int
	testBigInt.SetBit(new(big.Int), 256, 1).Sub(&testBigInt, big.NewInt(1))
	bigCoinAmount := sdktypes.NewIntFromBigInt(&testBigInt)

	err := input.GravityKeeper.DetectMaliciousSupply(input.Context, "stake", bigCoinAmount)
	require.Error(t, err, "didn't error out on too much added supply")
}

func TestERC20DeployedEventRejectsExistingTokenContractMapping(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	existingContract := common.HexToAddress("0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e")

	input.GravityKeeper.setCosmosOriginatedDenomToERC20(ctx, "uatom", existingContract)
	input.BankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "uosmo", Exponent: 0},
			{Denom: "osmo", Exponent: 6},
		},
		Base:    "uosmo",
		Display: "osmo",
	})

	err := input.GravityKeeper.Handle(ctx, &types.ERC20DeployedEvent{
		CosmosDenom:   "uosmo",
		TokenContract: existingContract.Hex(),
		Erc20Name:     "osmo",
		Erc20Symbol:   "osmo",
		Erc20Decimals: 6,
		EventNonce:    1,
	})

	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrInvalidERC20Event)
	require.Contains(t, err.Error(), "already exists for token contract")

	isCosmosOriginated, denom := input.GravityKeeper.ERC20ToDenomLookup(ctx, existingContract)
	require.True(t, isCosmosOriginated)
	require.Equal(t, "uatom", denom)

	_, _, err = input.GravityKeeper.DenomToERC20Lookup(ctx, "uosmo")
	require.Error(t, err)
}

func TestERC20DeployedEventRejectsEthereumOriginatedVoucherDenom(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	ethOriginatedContract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	deployedWrapperContract := common.HexToAddress("0x2222222222222222222222222222222222222222")
	voucherDenom := types.GravityDenom(ethOriginatedContract)

	require.NoError(t, input.BankKeeper.MintCoins(ctx, types.ModuleName, sdktypes.NewCoins(
		sdktypes.NewCoin(voucherDenom, sdktypes.NewInt(1)),
	)))

	err := input.GravityKeeper.Handle(ctx, &types.ERC20DeployedEvent{
		CosmosDenom:   voucherDenom,
		TokenContract: deployedWrapperContract.Hex(),
		Erc20Name:     voucherDenom,
		Erc20Symbol:   "",
		Erc20Decimals: 0,
		EventNonce:    1,
	})

	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrInvalidERC20Event)
	require.Contains(t, err.Error(), "ethereum-originated voucher denom")

	_, _, err = input.GravityKeeper.DenomToERC20Lookup(ctx, voucherDenom)
	require.NoError(t, err)

	isCosmosOriginated, denom := input.GravityKeeper.ERC20ToDenomLookup(ctx, deployedWrapperContract)
	require.False(t, isCosmosOriginated)
	require.Equal(t, types.GravityDenom(deployedWrapperContract), denom)
}
