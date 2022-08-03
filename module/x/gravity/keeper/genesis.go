package keeper

import (
	"fmt"
	"strconv"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, *data.Params)

	// reset delegate keys in state
	for _, keys := range data.DelegateKeys {
		if err := keys.ValidateBasic(); err != nil {
			panic(fmt.Sprintf("Invalid delegate key in Genesis: %s", err))
		}

		val, _ := sdk.ValAddressFromBech32(keys.ValidatorAddress)
		orch, _ := sdk.AccAddressFromBech32(keys.OrchestratorAddress)
		evm := common.HexToAddress(keys.EVMAddress)

		// set the orchestrator address
		k.SetOrchestratorValidatorAddress(ctx, val, orch)
		// set the evm address
		k.setValidatorEVMAddress(ctx, val, common.HexToAddress(keys.EVMAddress))
		k.setEVMOrchestratorAddress(ctx, evm, orch)
	}

	for _, chainGS := range data.EvmGenesisStates {
		if _, ok := data.Params.ParamsByChain[strconv.Itoa(int(chainGS.ChainID))]; !ok {
			panic(fmt.Sprintf("chain ID %d presented in state, but not in params", chainGS.ChainID))
		}

		initGenesisForChain(ctx, k, *chainGS)
	}
}

func initGenesisForChain(ctx sdk.Context, k Keeper, data types.EVMSpecificGenesisState) {
	// reset pool transactions in state
	for _, tx := range data.UnbatchedSendToEvmTxs {
		k.setUnbatchedSendToEVM(ctx, data.ChainID, tx)
	}

	// reset evm event vote records in state
	for _, evr := range data.EvmEventVoteRecords {
		event, err := types.UnpackEvent(evr.Event)
		if err != nil {
			panic(fmt.Sprintf("couldn't cast to event: %s", err))
		}
		if err := event.Validate(); err != nil {
			panic(fmt.Sprintf("invalid event in genesis: %s", err))
		}
		k.setEVMEventVoteRecord(ctx, data.ChainID, event.GetEventNonce(), event.Hash(), evr)
	}

	// reset last observed event nonce
	k.SetLastObservedEventNonce(ctx, data.ChainID, data.LastObservedEventNonce)

	// reset attestation state of all validators
	for _, eventVoteRecord := range data.EvmEventVoteRecords {
		event, _ := types.UnpackEvent(eventVoteRecord.Event)
		for _, vote := range eventVoteRecord.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := k.getLastEventNonceByValidator(ctx, data.ChainID, val)
			if event.GetEventNonce() > last {
				k.setLastEventNonceByValidator(ctx, data.ChainID, val, event.GetEventNonce())
			}
		}
	}

	// populate state with cosmos originated denom-erc20 mapping
	for _, item := range data.Erc20ToDenoms {
		k.setCosmosOriginatedDenomToERC20(ctx, data.ChainID, item.Denom, common.HexToAddress(item.Erc20))
	}

	// reset outgoing txs in state
	for _, ota := range data.OutgoingTxs {
		otx, err := types.UnpackOutgoingTx(ota)
		if err != nil {
			panic(fmt.Sprintf("invalid outgoing tx any in genesis file: %s", err))
		}
		k.SetOutgoingTx(ctx, otx)
	}

	// reset signatures in state
	for _, confa := range data.Confirmations {
		conf, err := types.UnpackConfirmation(confa)
		if err != nil {
			panic(fmt.Sprintf("invalid evm signature in genesis: %s", err))
		}
		// TODO: not currently an easy way to get the validator address from the
		// evm address here. once we implement the third index for keys
		// this will be easy.
		k.SetEVMSignature(ctx, conf, sdk.ValAddress{})
	}
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var (
		p         = k.GetParams(ctx)
		delegates = k.getDelegateKeys(ctx)
	)

	var EVMSpecificGenesisStates []*types.EVMSpecificGenesisState

	for _, chainID := range k.GetChainIDs(ctx) {
		var (
			outgoingTxs           []*cdctypes.Any
			evmTxConfirmations    []*cdctypes.Any
			attmap                = k.GetEVMEventVoteRecordMapping(ctx, chainID)
			evmEventVoteRecords   []*types.EVMEventVoteRecord
			lastobserved          = k.GetLastObservedEventNonce(ctx, chainID)
			erc20ToDenoms         []*types.ERC20ToDenom
			unbatchedSendToEVMTxs = k.getUnbatchedSendToEVMs(ctx, chainID)
		)

		// export evmEventVoteRecords from state
		for _, atts := range attmap {
			// TODO: set height = 0?
			evmEventVoteRecords = append(evmEventVoteRecords, atts...)
		}

		// export erc20 to denom relations
		k.iterateERC20ToDenom(ctx, chainID, func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
			erc20ToDenoms = append(erc20ToDenoms, erc20ToDenom)
			return false
		})

		// export signer set txs and sigs
		k.IterateOutgoingTxsByType(ctx, chainID, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
			ota, _ := types.PackOutgoingTx(otx)
			outgoingTxs = append(outgoingTxs, ota)
			sstx, _ := otx.(*types.SignerSetTx)
			k.iterateEVMSignaturesByStoreIndex(ctx, sstx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
				siga, _ := types.PackConfirmation(&types.SignerSetTxConfirmation{
					SignerSetNonce: sstx.Nonce,
					EVMSigner:      k.GetValidatorEVMAddress(ctx, val).Hex(),
					Signature:      sig,
					ChainId:        chainID,
				})
				evmTxConfirmations = append(evmTxConfirmations, siga)
				return false
			})
			return false
		})

		// export batch txs and sigs
		k.IterateOutgoingTxsByType(ctx, chainID, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
			ota, _ := types.PackOutgoingTx(otx)
			outgoingTxs = append(outgoingTxs, ota)
			btx, _ := otx.(*types.BatchTx)
			k.iterateEVMSignaturesByStoreIndex(ctx, btx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
				siga, _ := types.PackConfirmation(&types.BatchTxConfirmation{
					TokenContract: btx.TokenContract,
					BatchNonce:    btx.BatchNonce,
					EVMSigner:     k.GetValidatorEVMAddress(ctx, val).Hex(),
					Signature:     sig,
					ChainId:       chainID,
				})
				evmTxConfirmations = append(evmTxConfirmations, siga)
				return false
			})
			return false
		})

		// export contract call txs and sigs
		k.IterateOutgoingTxsByType(ctx, chainID, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
			ota, _ := types.PackOutgoingTx(otx)
			outgoingTxs = append(outgoingTxs, ota)
			cctx, _ := otx.(*types.ContractCallTx)
			k.iterateEVMSignaturesByStoreIndex(ctx, cctx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
				siga, _ := types.PackConfirmation(&types.ContractCallTxConfirmation{
					InvalidationScope: cctx.InvalidationScope,
					InvalidationNonce: cctx.InvalidationNonce,
					EVMSigner:         k.GetValidatorEVMAddress(ctx, val).Hex(),
					Signature:         sig,
					ChainId:           chainID,
				})
				evmTxConfirmations = append(evmTxConfirmations, siga)
				return false
			})
			return false
		})

		cgs := &types.EVMSpecificGenesisState{
			ChainID:                chainID,
			LastObservedEventNonce: lastobserved,
			OutgoingTxs:            outgoingTxs,
			Confirmations:          evmTxConfirmations,
			EvmEventVoteRecords:    evmEventVoteRecords,
			Erc20ToDenoms:          erc20ToDenoms,
			UnbatchedSendToEvmTxs:  unbatchedSendToEVMTxs,
		}

		EVMSpecificGenesisStates = append(EVMSpecificGenesisStates, cgs)
	}

	// this will marshal into "dW51c2Vk" as []byte will be encoded as base64
	for _, delegate := range delegates {
		delegate.EVMSignature = []byte("unused")
	}

	return types.GenesisState{
		Params:           &p,
		DelegateKeys:     delegates,
		EvmGenesisStates: EVMSpecificGenesisStates,
	}
}
