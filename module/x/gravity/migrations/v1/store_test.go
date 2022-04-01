package v1

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/keeper"
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
