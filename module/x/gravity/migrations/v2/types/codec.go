package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers the vesting interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDelegateKeys{}, "gravity-bridge/", nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSendToEthereum{},
		&MsgCancelSendToEthereum{},
		&MsgRequestBatchTx{},
		&MsgSubmitEthereumEvent{},
		&MsgSubmitEthereumTxConfirmation{},
		&MsgDelegateKeys{},
		&MsgEthereumHeightVote{},
	)

	registry.RegisterInterface(
		"gravity.v1.EthereumEvent",
		(*EthereumEvent)(nil),
		&SendToCosmosEvent{},
		&BatchExecutedEvent{},
		&ERC20DeployedEvent{},
		&ContractCallExecutedEvent{},
		&SignerSetTxExecutedEvent{},
	)

	registry.RegisterInterface(
		"gravity.v1.EthereumSignature",
		(*EthereumTxConfirmation)(nil),
		&BatchTxConfirmation{},
		&ContractCallTxConfirmation{},
		&SignerSetTxConfirmation{},
	)

	registry.RegisterInterface(
		"gravity.v1.OutgoingTx",
		(*OutgoingTx)(nil),
		&SignerSetTx{},
		&BatchTx{},
		&ContractCallTx{},
	)

	registry.RegisterImplementations((*govtypes.Content)(nil),
		&CommunityPoolEthereumSpendProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
