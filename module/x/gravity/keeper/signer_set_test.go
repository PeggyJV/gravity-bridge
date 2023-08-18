package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
)

func TestSignerSetTxExecuted(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context.WithBlockHeight(100)

	input.GravityKeeper.CreateSignerSetTx(
		ctx,
	)

	sstx := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(1)).(*types.SignerSetTx)

	assert.NotNil(t, sstx)

	input.GravityKeeper.SignerSetExecuted(ctx, 1)

	otx := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(1))

	assert.NotNil(t, otx)

	cotx := input.GravityKeeper.GetCompletedOutgoingTx(ctx, sstx.GetStoreIndex())

	assert.NotNil(t, cotx)
	assert.Equal(t, uint64(1), input.GravityKeeper.GetLatestSignerSetTxNonce(ctx))
}

func TestGetUnconfirmedSignerSetTxs(t *testing.T) {
	input, ctx := SetupFiveValChain(t)
	gk := input.GravityKeeper
	vals := input.StakingKeeper.GetAllValidators(ctx)
	val1, err := sdk.ValAddressFromBech32(vals[0].OperatorAddress)
	require.NoError(t, err)
	val2, err := sdk.ValAddressFromBech32(vals[1].OperatorAddress)
	require.NoError(t, err)

	gk.CreateSignerSetTx(ctx)
	gk.SetCompletedOutgoingTx(ctx, &types.SignerSetTx{
		Nonce:   2,
		Height:  0,
		Signers: []*types.EthereumSigner{},
	})

	// val1 signs both
	// val2 signs one
	gk.SetEthereumSignature(
		ctx,
		&types.SignerSetTxConfirmation{
			SignerSetNonce: 1,
			EthereumSigner: "",
			Signature:      []byte("dummysig"),
		},
		val1,
	)
	gk.SetEthereumSignature(
		ctx,
		&types.SignerSetTxConfirmation{
			SignerSetNonce: 2,
			EthereumSigner: "",
			Signature:      []byte("dummysig"),
		},
		val1,
	)
	gk.SetEthereumSignature(
		ctx,
		&types.SignerSetTxConfirmation{
			SignerSetNonce: 2,
			EthereumSigner: "",
			Signature:      []byte("dummysig"),
		},
		val2,
	)

	require.Empty(t, gk.GetUnsignedSignerSetTxs(ctx, val1))
	require.Equal(t, 1, len(gk.GetUnsignedSignerSetTxs(ctx, val2)))
}
