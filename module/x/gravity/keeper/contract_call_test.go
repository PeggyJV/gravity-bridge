package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
	"github.com/stretchr/testify/assert"
)

func TestContractCallTxExecuted(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context.WithBlockHeight(100)
	storeKey := input.GravityStoreKey
	cdc := input.Marshaler

	latestEVMBlockHeight := &types.LatestEVMBlockHeight{
		CosmosHeight: 100,
		EVMHeight:    1000,
	}

	ctx.KVStore(storeKey).Set([]byte{types.LastEVMBlockHeightKey}, cdc.MustMarshal(latestEVMBlockHeight))

	scope := []byte("test-scope")
	contract := common.HexToAddress("0x2a24af0501a534fca004ee1bd667b783f205a546")
	nonce1 := uint64(1)
	nonce2 := uint64(2)
	payload := []byte("payload")
	erc20Tokens := []types.ERC20Token{
		{
			Contract: "0x2a24af0501a534fca004ee1bd667b783f205a546",
			Amount:   sdk.NewInt(1),
		},
	}

	input.GravityKeeper.CreateContractCallTx(
		ctx,
		types.EthereumChainID,
		nonce1,
		scope,
		contract,
		payload,
		erc20Tokens,
		erc20Tokens,
	)

	input.GravityKeeper.CreateContractCallTx(
		ctx,
		types.EthereumChainID,
		nonce2,
		scope,
		contract,
		payload,
		erc20Tokens,
		erc20Tokens,
	)

	cctx1 := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeContractCallTxStoreIndex(types.EthereumChainID, scope, nonce1)).(*types.ContractCallTx)
	assert.Equal(t, scope, cctx1.InvalidationScope)
	assert.Equal(t, nonce1, cctx1.InvalidationNonce)
	assert.Equal(t, contract.Hex(), cctx1.Address)
	assert.Equal(t, payload, cctx1.Payload)
	assert.Equal(t, erc20Tokens, cctx1.Tokens)
	assert.Equal(t, erc20Tokens, cctx1.Fees)
	assert.Equal(t, uint32(types.EthereumChainID), cctx1.ChainId)

	cctx2 := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeContractCallTxStoreIndex(types.EthereumChainID, scope, nonce2)).(*types.ContractCallTx)
	assert.Equal(t, scope, cctx2.InvalidationScope)
	assert.Equal(t, nonce2, cctx2.InvalidationNonce)
	assert.Equal(t, contract.Hex(), cctx2.Address)
	assert.Equal(t, payload, cctx2.Payload)
	assert.Equal(t, erc20Tokens, cctx2.Tokens)
	assert.Equal(t, erc20Tokens, cctx2.Fees)
	assert.Equal(t, uint32(types.EthereumChainID), cctx2.ChainId)

	input.GravityKeeper.contractCallExecuted(ctx, types.EthereumChainID, scope, nonce2)

	otx1 := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeContractCallTxStoreIndex(types.EthereumChainID, scope, nonce1))
	otx2 := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeContractCallTxStoreIndex(types.EthereumChainID, scope, nonce2))

	assert.Nil(t, otx1)
	assert.Nil(t, otx2)
}
