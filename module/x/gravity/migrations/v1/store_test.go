package v1_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/keeper"
	v1 "github.com/peggyjv/gravity-bridge/module/x/gravity/migrations/v1"
	"github.com/stretchr/testify/assert"
)

const denom string = "cosmos"
const tokenContractString string = "0x2a24af0501a534fca004ee1bd667b783f205a546"

func TestMigrateCosmosOriginatedERC20ToDenom(t *testing.T) {
	input := keeper.CreateTestEnv(t)
	ctx := input.Context
	storeKey := input.GravityStoreKey

	ctx.KVStore(storeKey).Set(v1.MakeOldERC20ToDenomKey(tokenContractString), []byte(denom))

	err := v1.MigrateStore(ctx, storeKey, input.Marshaler)
	assert.NoError(t, err)

	tokenContract := common.HexToAddress(tokenContractString)
	storedDenom := ctx.KVStore(storeKey).Get(v1.MakeNewERC20ToDenomKey(tokenContract))
	assert.Equal(t, denom, string(storedDenom))
}
