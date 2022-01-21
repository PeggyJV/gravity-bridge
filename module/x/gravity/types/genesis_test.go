package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisStateValidate(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForValidator("sommvaloper", "sommvaloperpub")
	sdk.GetConfig().SetBech32PrefixForAccount("somm", "sommpub")
	sdk.GetConfig().SetBech32PrefixForConsensusNode("sommcons", "sommconspub")

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
					ValidatorAddress:    "sommvaloper1w4lcqzsrfgfwrgd386v9n6frj9p66k4xazvsqn",
					OrchestratorAddress: "somm1ht877l3pgqjkalxwtvmydr65jnde3p4g7vxq5t",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        []byte("xIm9M32dXl7II8c72qXS/x9xVQMbX9sHVkuigxSmspNhSnEc2dLe6YDn/+Yi3VAOwXKh8zUBSHMplQSEuOY1zhw="),
				},
			},
		}, expErr: false},
		"valid delegate with placeholder signature": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "sommvaloper1w4lcqzsrfgfwrgd386v9n6frj9p66k4xazvsqn",
					OrchestratorAddress: "somm1ht877l3pgqjkalxwtvmydr65jnde3p4g7vxq5t",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        []byte("unused"),
				},
			},
		}, expErr: false},
		"valid delegate with nil signature": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "sommvaloper1w4lcqzsrfgfwrgd386v9n6frj9p66k4xazvsqn",
					OrchestratorAddress: "somm1ht877l3pgqjkalxwtvmydr65jnde3p4g7vxq5t",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        nilByteSlice,
				},
			},
		}, expErr: true},
		"delegate with bad validator address": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "sommvaloper1wrong",
					OrchestratorAddress: "somm1ht877l3pgqjkalxwtvmydr65jnde3p4g7vxq5t",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        []byte("xIm9M32dXl7II8c72qXS/x9xVQMbX9sHVkuigxSmspNhSnEc2dLe6YDn/+Yi3VAOwXKh8zUBSHMplQSEuOY1zhw="),
				},
			},
		}, expErr: true},
		"delegate with bad orchestrator address": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "sommvaloper1w4lcqzsrfgfwrgd386v9n6frj9p66k4xazvsqn",
					OrchestratorAddress: "somm1wrong",
					EthereumAddress:     "0x494eeff8848254C4fdd5B529FC6E751Ab34597A6",
					EthSignature:        []byte("xIm9M32dXl7II8c72qXS/x9xVQMbX9sHVkuigxSmspNhSnEc2dLe6YDn/+Yi3VAOwXKh8zUBSHMplQSEuOY1zhw="),
				},
			},
		}, expErr: true},
		"delegate with bad eth address": {src: &GenesisState{
			Params: DefaultParams(),
			DelegateKeys: []*MsgDelegateKeys{
				{
					ValidatorAddress:    "sommvaloper1w4lcqzsrfgfwrgd386v9n6frj9p66k4xazvsqn",
					OrchestratorAddress: "somm1ht877l3pgqjkalxwtvmydr65jnde3p4g7vxq5t",
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
