package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisStateValidate(t *testing.T) {
	var nilByteSlice []byte
	specs := map[string]struct {
		src    *GenesisState
		expErr bool
	}{
		"default params": {src: DefaultGenesisState(), expErr: false},
		"empty params":   {src: &GenesisState{Params: &Params{}}, expErr: true},
		"invalid params": {src: &GenesisState{
			Params: &Params{
				GravityId:             "foo",
				ContractSourceHash:    "laksdjflasdkfja",
				BridgeEthereumAddress: "invalid-eth-address",
				BridgeChainId:         3279089,
			},
		}, expErr: true},
		"valid delegate": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					// note: the "valid" EthSignature here is simply the correct format, it is not a signature from a private key represented
					// by the above cosmos address, but genesis isn't currently validating that anyway
					ValidatorAddress:    "cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4",
					OrchestratorAddress: "cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        []byte("0x2471d20201d38a6d8f5301be45560161d770f8ca8642ac45a1eb5c82fd853a8670fb6f18aaedffd7fcf27bba62d29b782a1c12059ab43b79d7d2d727f596d7701c"),
				},
			},
		}, expErr: false},
		"valid delegate with placeholder signature": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4",
					OrchestratorAddress: "cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        []byte("unused"), // this will marshal into "dW51c2Vk" as []byte will be encoded as base64
				},
			},
		}, expErr: false},
		"valid delegate with nil signature": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4",
					OrchestratorAddress: "cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        nilByteSlice,
				},
			},
		}, expErr: true},
		"delegate with bad validator address": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "cosmosvaloper1wrong",
					OrchestratorAddress: "cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        []byte("0x2471d20201d38a6d8f5301be45560161d770f8ca8642ac45a1eb5c82fd853a8670fb6f18aaedffd7fcf27bba62d29b782a1c12059ab43b79d7d2d727f596d7701c"),
				},
			},
		}, expErr: true},
		"delegate with bad orchestrator address": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4",
					OrchestratorAddress: "cosmos1wrong",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        []byte("0x2471d20201d38a6d8f5301be45560161d770f8ca8642ac45a1eb5c82fd853a8670fb6f18aaedffd7fcf27bba62d29b782a1c12059ab43b79d7d2d727f596d7701c"),
				},
			},
		}, expErr: true},
		"delegate with bad eth address": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4",
					OrchestratorAddress: "cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7",
					EthereumAddress:     "0xdeadbeef",
					EthSignature:        []byte("0x2471d20201d38a6d8f5301be45560161d770f8ca8642ac45a1eb5c82fd853a8670fb6f18aaedffd7fcf27bba62d29b782a1c12059ab43b79d7d2d727f596d7701c"),
				},
			},
		}, expErr: true},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestStringToByteArray(t *testing.T) {
	specs := map[string]struct {
		testString string
		expErr     bool
	}{
		"16 bytes": {"lakjsdflaksdjfds", false},
		"32 bytes": {"lakjsdflaksdjfdslakjsdflaksdjfds", false},
		"33 bytes": {"€€€€€€€€€€€", true},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			_, err := strToFixByteArray(spec.testString)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
