package types

import (
	"bytes"
	"fmt"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParamspace defines the default auth module parameter subspace
const (
	DefaultParamspace = ModuleName
)

var (
	// ParamsStoreKeyParamsForChain stores a map of chain ID to chain specific params
	ParamsStoreKeyParamsForChain   = []byte("ParamsForChain")
	ParamsStoreKeyAverageBlockTime = []byte("AverageBlockTime")

	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{}
)

func (gs *EVMSpecificGenesisState) UnpackInterfaces(unpacker cdctypes.AnyUnpacker) error {
	for _, otx := range gs.OutgoingTxs {
		var outgoing OutgoingTx
		if err := unpacker.UnpackAny(otx, &outgoing); err != nil {
			return err
		}
	}
	for _, sig := range gs.Confirmations {
		var signature EVMTxConfirmation
		if err := unpacker.UnpackAny(sig, &signature); err != nil {
			return err
		}
	}
	for _, evr := range gs.EvmEventVoteRecords {
		if err := evr.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}
	return nil
}

func EventVoteRecordPowerThreshold(totalPower sdk.Int) sdk.Int {
	return sdk.NewInt(66).Mul(totalPower).Quo(sdk.NewInt(100))
}

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (s GenesisState) ValidateBasic() error {
	if err := s.Params.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}
	if len(s.DelegateKeys) != 0 {
		for _, delegateKey := range s.DelegateKeys {
			if err := delegateKey.ValidateBasic(); err != nil {
				return sdkerrors.Wrap(err, "delegates")
			}
		}
	}
	return nil
}

// DefaultGenesisState returns empty genesis state
// TODO: set some better defaults here
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// DefaultParams returns a copy of the default params
func DefaultParams() *Params {
	cp := DefaultParamsForChain()

	return &Params{
		AverageBlockTime: 5000,
		ParamsForChains:  []*ParamsForChain{cp},
	}
}

func DefaultParamsForChain() *ParamsForChain {
	return &ParamsForChain{
		ChainId:                              EthereumChainID,
		GravityId:                            "defaultgravityid",
		SignedSignerSetTxsWindow:             10000,
		SignedBatchesWindow:                  10000,
		EvmSignaturesWindow:                  10000,
		TargetEvmTxTimeout:                   43200000,
		AverageEvmBlockTime:                  15000,
		SlashFractionSignerSetTx:             sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionBatch:                   sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionEvmSignature:            sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionConflictingEvmSignature: sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		UnbondSlashingSignerSetTxsWindow:     10000}
}

// ValidateBasic checks that the parameters have valid values.
func (p Params) ValidateBasic() error {
	gravityIDs := make(map[string]bool)
	chainIDs := make(map[uint32]bool)

	if err := validateAverageBlockTime(p.AverageBlockTime); err != nil {
		return sdkerrors.Wrap(err, "Block time")
	}

	for _, cp := range p.ParamsForChains {
		if err := cp.ValidateBasic(); err != nil {
			return err
		}

		if exists := gravityIDs[cp.GravityId]; exists {
			return sdkerrors.Wrap(ErrDuplicateGravityID, "gravity id")
		}
		gravityIDs[cp.GravityId] = true

		if exists := chainIDs[cp.ChainId]; exists {
			return sdkerrors.Wrap(ErrDuplicateChainID, "chain id")
		}
		chainIDs[cp.ChainId] = true
	}

	return nil
}

func (cp ParamsForChain) ValidateBasic() error {
	if err := validateChainID(cp.ChainId); err != nil {
		return sdkerrors.Wrap(err, "chain id")
	}
	if err := validateGravityID(cp.GravityId); err != nil {
		return sdkerrors.Wrap(err, "gravity id")
	}
	if err := validateTargetEvmTxTimeout(cp.TargetEvmTxTimeout); err != nil {
		return sdkerrors.Wrap(err, "Batch timeout")
	}
	if err := validateAverageEVMBlockTime(cp.AverageEvmBlockTime); err != nil {
		return sdkerrors.Wrap(err, "EVM block time")
	}
	if err := validateSignedSignerSetTxsWindow(cp.SignedSignerSetTxsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed signer set txs window")
	}
	if err := validateSignedBatchesWindow(cp.SignedBatchesWindow); err != nil {
		return sdkerrors.Wrap(err, "signed batches window")
	}
	if err := validateEVMSignaturesWindow(cp.EvmSignaturesWindow); err != nil {
		return sdkerrors.Wrap(err, "signatures window")
	}
	if err := validateSlashFractionSignerSetTx(cp.SlashFractionSignerSetTx); err != nil {
		return sdkerrors.Wrap(err, "slash fraction signersettx")
	}
	if err := validateSlashFractionBatch(cp.SlashFractionBatch); err != nil {
		return sdkerrors.Wrap(err, "slash fraction batch tx")
	}
	if err := validateSlashFractionEVMSignature(cp.SlashFractionEvmSignature); err != nil {
		return sdkerrors.Wrap(err, "slash fraction EVM signature")
	}
	if err := validateSlashFractionConflictingEVMSignature(cp.SlashFractionConflictingEvmSignature); err != nil {
		return sdkerrors.Wrap(err, "slash fraction conflicting EVM signature")
	}
	if err := validateUnbondSlashingSignerSetTxsWindow(cp.UnbondSlashingSignerSetTxsWindow); err != nil {
		return sdkerrors.Wrap(err, "unbond slashing signersettx window")
	}

	return nil
}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value
// pairs of auth module's parameters.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageBlockTime, &p.AverageBlockTime, validateAverageBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreKeyParamsForChain, &p.ParamsForChains, validateParamsForChain),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	pb, err := p.Marshal()
	if err != nil {
		panic(err)
	}
	p2b, err := p2.Marshal()
	if err != nil {
		panic(err)
	}
	return bytes.Equal(pb, p2b)
}

func validateParamsForChain(i interface{}) error {
	v, ok := i.([]*ParamsForChain)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, cp := range v {
		if err := cp.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

func validateChainID(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateGravityID(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if _, err := strToFixByteArray(v); err != nil {
		return err
	}
	return nil
}

func validateTargetEvmTxTimeout(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 60000 {
		return fmt.Errorf("invalid target batch timeout, less than 60 seconds is too short")
	}
	return nil
}

func validateAverageBlockTime(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 100 {
		return fmt.Errorf("invalid average Cosmos block time, too short for latency limitations")
	}
	return nil
}

func validateAverageEVMBlockTime(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 100 {
		return fmt.Errorf("invalid average EVM block time, too short for latency limitations")
	}
	return nil
}

func validateSignedSignerSetTxsWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateUnbondSlashingSignerSetTxsWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionSignerSetTx(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSignedBatchesWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateEVMSignaturesWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionBatch(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionEVMSignature(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionConflictingEVMSignature(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func strToFixByteArray(s string) ([32]byte, error) {
	var out [32]byte
	if len([]byte(s)) > 32 {
		return out, fmt.Errorf("string too long")
	}
	copy(out[:], s)
	return out, nil
}

func byteArrayToFixByteArray(b []byte) (out [32]byte, err error) {
	if len(b) > 32 {
		return out, fmt.Errorf("array too long")
	}
	copy(out[:], b)
	return out, nil
}
