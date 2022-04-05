package v1

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/keeper"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
	"github.com/stretchr/testify/assert"
)

const denom string = "cosmos"
const tokenContractString string = "0x2a24af0501a534fca004ee1bd667b783f205a546"

func TestMigrateCosmosOriginatedERC20ToDenom(t *testing.T) {
	input := keeper.CreateTestEnv(t)
	ctx := input.Context
	storeKey := input.GravityStoreKey

	ctx.KVStore(storeKey).Set(MakeOldERC20ToDenomKey(tokenContractString), []byte(denom))

	err := MigrateStore(ctx, storeKey, input.Marshaler)
	assert.NoError(t, err)

	tokenContract := common.HexToAddress(tokenContractString)
	storedDenom := ctx.KVStore(storeKey).Get(MakeNewERC20ToDenomKey(tokenContract))
	assert.Equal(t, denom, string(storedDenom))
}

func TestMigrateContractCallTxTimeout(t *testing.T) {
	input := keeper.CreateTestEnv(t)
	ctx := input.Context
	storeKey := input.GravityStoreKey
	cdc := input.Marshaler

	latestEthereumBlockHeight := &types.LatestEthereumBlockHeight{
		CosmosHeight:   100,
		EthereumHeight: 1000,
	}

	ctx.KVStore(storeKey).Set([]byte{LastEthereumBlockHeightKey}, cdc.MustMarshal(latestEthereumBlockHeight))

	otx := &types.ContractCallTx{
		InvalidationNonce: uint64(1),
		InvalidationScope: []byte("test-scope"),
		Address:           "0x2a24af0501a534fca004ee1bd667b783f205a546",
		Payload:           []byte("payload"),
		Timeout:           uint64(1000),
		Tokens:            []types.ERC20Token{},
		Fees:              []types.ERC20Token{},
		Height:            uint64(100),
	}

	any, err := types.PackOutgoingTx(otx)
	if err != nil {
		panic(err)
	}

	ctx.KVStore(storeKey).Set(
		MakeOutgoingTxKey(otx.GetStoreIndex()),
		cdc.MustMarshal(any),
	)

	err = MigrateStore(ctx, storeKey, input.Marshaler)
	assert.NoError(t, err)

	var updatedOtx types.OutgoingTx
	err = cdc.UnmarshalInterface(ctx.KVStore(storeKey).Get(types.MakeOutgoingTxKey(otx.GetStoreIndex())), &updatedOtx)
	assert.NoError(t, err)

	cctx, _ := updatedOtx.(*types.ContractCallTx)
	assert.Equal(t, cctx.Timeout, uint64(1800))
}
