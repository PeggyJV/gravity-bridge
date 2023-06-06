package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func TestSignerSetTxExecuted(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context.WithBlockHeight(100)
	// storeKey := input.GravityStoreKey
	// cdc := input.Marshaler

	// latestEthereumBlockHeight := &types.LatestEthereumBlockHeight{
	// 	CosmosHeight:   100,
	// 	EthereumHeight: 1000,
	// }

	// ctx.KVStore(storeKey).Set([]byte{types.LastEthereumBlockHeightKey}, cdc.MustMarshal(latestEthereumBlockHeight))

	input.GravityKeeper.CreateSignerSetTx(
		ctx,
	)

	input.GravityKeeper.CreateSignerSetTx(
		ctx,
	)

	sstx1 := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(1)).(*types.SignerSetTx)
	sstx2 := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(2)).(*types.SignerSetTx)

	assert.NotNil(t, sstx1)
	assert.NotNil(t, sstx2)

	input.GravityKeeper.signerSetExecuted(ctx, 2)

	otx1 := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(1))
	otx2 := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(2))

	assert.Nil(t, otx1)
	assert.Nil(t, otx2)

	otx2 = input.GravityKeeper.GetCompletedOutgoingTx(ctx, sstx2.GetStoreIndex())

	assert.NotNil(t, otx2)
	assert.Equal(t, uint64(2), input.GravityKeeper.GetLatestSignerSetTxNonce(ctx))
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

	require.Empty(t, gk.GetUnconfirmedSignerSetTxs(ctx, val1))
	require.Equal(t, 1, len(gk.GetUnconfirmedSignerSetTxs(ctx, val2)))
}
