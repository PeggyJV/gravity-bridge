package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/v5/x/gravity/types"
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

func TestOrderSignerSetsByNonceAscending(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		// Create test signer sets with different nonces
		signerSets := []*types.SignerSetTx{
			{Nonce: 3},
			{Nonce: 1},
			{Nonce: 4},
			{Nonce: 2},
		}

		// Order the signer sets
		orderedSignerSets := orderSignerSetsByNonceAscending(signerSets)

		// Check if the signer sets are ordered correctly
		assert.Equal(t, uint64(1), orderedSignerSets[0].Nonce)
		assert.Equal(t, uint64(2), orderedSignerSets[1].Nonce)
		assert.Equal(t, uint64(3), orderedSignerSets[2].Nonce)
		assert.Equal(t, uint64(4), orderedSignerSets[3].Nonce)

		// Check if the length of the slice remains the same
		assert.Equal(t, len(signerSets), len(orderedSignerSets))
	})

	t.Run("empty slice", func(t *testing.T) {
		signerSets := []*types.SignerSetTx{}
		orderedSignerSets := orderSignerSetsByNonceAscending(signerSets)
		assert.Empty(t, orderedSignerSets)
	})

	t.Run("nil slice", func(t *testing.T) {
		var signerSets []*types.SignerSetTx
		orderedSignerSets := orderSignerSetsByNonceAscending(signerSets)
		assert.Nil(t, orderedSignerSets)
	})

	t.Run("single element", func(t *testing.T) {
		signerSets := []*types.SignerSetTx{{Nonce: 1}}
		orderedSignerSets := orderSignerSetsByNonceAscending(signerSets)
		assert.Equal(t, 1, len(orderedSignerSets))
		assert.Equal(t, uint64(1), orderedSignerSets[0].Nonce)
	})

	t.Run("duplicate nonces", func(t *testing.T) {
		signerSets := []*types.SignerSetTx{
			{Nonce: 2},
			{Nonce: 1},
			{Nonce: 2},
			{Nonce: 1},
		}
		orderedSignerSets := orderSignerSetsByNonceAscending(signerSets)
		assert.Equal(t, 4, len(orderedSignerSets))
		assert.Equal(t, uint64(1), orderedSignerSets[0].Nonce)
		assert.Equal(t, uint64(1), orderedSignerSets[1].Nonce)
		assert.Equal(t, uint64(2), orderedSignerSets[2].Nonce)
		assert.Equal(t, uint64(2), orderedSignerSets[3].Nonce)
	})
}
