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
					ValidatorAddress:    "cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4",
					OrchestratorAddress: "cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        []byte("xIm9M32dXl7II8c72qXS/x9xVQMbX9sHVkuigxSmspNhSnEc2dLe6YDn/+Yi3VAOwXKh8zUBSHMplQSEuOY1zhw="),
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
					EthSignature:        []byte("xIm9M32dXl7II8c72qXS/x9xVQMbX9sHVkuigxSmspNhSnEc2dLe6YDn/+Yi3VAOwXKh8zUBSHMplQSEuOY1zhw="),
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
					EthSignature:        []byte("xIm9M32dXl7II8c72qXS/x9xVQMbX9sHVkuigxSmspNhSnEc2dLe6YDn/+Yi3VAOwXKh8zUBSHMplQSEuOY1zhw="),
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
					EthSignature:        []byte("xIm9M32dXl7II8c72qXS/x9xVQMbX9sHVkuigxSmspNhSnEc2dLe6YDn/+Yi3VAOwXKh8zUBSHMplQSEuOY1zhw="),
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
