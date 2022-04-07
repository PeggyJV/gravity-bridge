package v1_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	v1 "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/migrations/v1"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/migrations/v1/keeper"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/migrations/v1/types"
	"github.com/stretchr/testify/assert"
)

const denom string = "cosmos"
const tokenContractString string = "0x2a24af0501a534fca004ee1bd667b783f205a546"

func TestMigrateCosmosOriginatedERC20ToDenom(t *testing.T) {
	input := keeper.CreateTestEnv(t)
	ctx := input.Context
	storeKey := input.GravityStoreKey

	ctx.KVStore(storeKey).Set(types.MakeERC20ToDenomKey(tokenContractString), []byte(denom))

	err := v1.MigrateStore(ctx, storeKey, input.Marshaler)
	assert.NoError(t, err)

	tokenContract := common.HexToAddress(tokenContractString)
	storedDenom := ctx.KVStore(storeKey).Get(types.MakeNewERC20ToDenomKey(tokenContract))
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

	ctx.KVStore(storeKey).Set([]byte{types.LastEthereumBlockHeightKey}, cdc.MustMarshal(latestEthereumBlockHeight))

	erc20Token := types.ERC20Token{
		Contract: tokenContractString,
		Amount:   sdk.NewInt(1),
	}

	otx := &types.ContractCallTx{
		InvalidationNonce: uint64(1),
		InvalidationScope: []byte("test-scope"),
		Address:           "0x2a24af0501a534fca004ee1bd667b783f205a546",
		Payload:           []byte("payload"),
		Timeout:           uint64(1000),
		Tokens:            []types.ERC20Token{erc20Token},
		Fees:              []types.ERC20Token{erc20Token},
		Height:            uint64(100),
	}

	any, err := types.PackOutgoingTx(otx)
	if err != nil {
		panic(err)
	}

	ctx.KVStore(storeKey).Set(
		types.MakeOutgoingTxKey(otx.GetStoreIndex()),
		cdc.MustMarshal(any),
	)

	err = v1.MigrateStore(ctx, storeKey, input.Marshaler)
	assert.NoError(t, err)

	var updatedOtx types.OutgoingTx
	err = cdc.UnmarshalInterface(ctx.KVStore(storeKey).Get(types.MakeOutgoingTxKey(otx.GetStoreIndex())), &updatedOtx)
	assert.NoError(t, err)

	cctx, _ := updatedOtx.(*types.ContractCallTx)
	assert.Equal(t, cctx.InvalidationNonce, uint64(1))
	assert.Equal(t, cctx.InvalidationScope, []byte("test-scope"))
	assert.Equal(t, cctx.Address, "0x2a24af0501a534fca004ee1bd667b783f205a546")
	assert.Equal(t, cctx.Payload, []byte("payload"))
	assert.Equal(t, cctx.Timeout, uint64(1800)) // migration added 800 to the latest height of 1000
	assert.Equal(t, cctx.Tokens, []types.ERC20Token{erc20Token})
	assert.Equal(t, cctx.Fees, []types.ERC20Token{erc20Token})
	assert.Equal(t, cctx.Height, uint64(100))
}
