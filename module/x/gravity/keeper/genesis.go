package keeper

import (
	"fmt"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
)

// InitGenesis starts a chain from a genesis state
//
// Deprecated: please use InitGenesisMultiChain
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.setParams(ctx, *data.Params)

	// reset pool transactions in state
	for _, tx := range data.UnbatchedSendToEthereumTxs {
		k.setUnbatchedSendToEthereum(ctx, uint32(data.Params.BridgeChainId), tx)
	}

	// reset ethereum event vote records in state
	for _, evr := range data.EthereumEventVoteRecords {
		event, err := types.UnpackEvent(evr.Event)
		if err != nil {
			panic(fmt.Sprintf("couldn't cast to event: %s", err))
		}
		if err := event.Validate(); err != nil {
			panic(fmt.Sprintf("invalid event in genesis: %s", err))
		}
		k.setEthereumEventVoteRecord(ctx, uint32(data.Params.BridgeChainId), event.GetEventNonce(), event.Hash(), evr)
	}

	// reset last observed event nonce
	k.setLastObservedEventNonce(ctx, uint32(data.Params.BridgeChainId), data.LastObservedEventNonce)

	// reset attestation state of all validators
	for _, eventVoteRecord := range data.EthereumEventVoteRecords {
		event, _ := types.UnpackEvent(eventVoteRecord.Event)
		for _, vote := range eventVoteRecord.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := k.getLastEventNonceByValidator(ctx, uint32(data.Params.BridgeChainId), val)
			if event.GetEventNonce() > last {
				k.setLastEventNonceByValidator(ctx, uint32(data.Params.BridgeChainId), val, event.GetEventNonce())
			}
		}
	}

	// reset delegate keys in state
	for _, keys := range data.DelegateKeys {
		if err := keys.ValidateBasic(); err != nil {
			panic(fmt.Sprintf("Invalid delegate key in Genesis: %s", err))
		}

		val, _ := sdk.ValAddressFromBech32(keys.ValidatorAddress)
		orch, _ := sdk.AccAddressFromBech32(keys.OrchestratorAddress)
		eth := common.HexToAddress(keys.EthereumAddress)

		// set the orchestrator address
		k.SetOrchestratorValidatorAddress(ctx, val, orch)
		// set the ethereum address
		k.setValidatorEthereumAddress(ctx, val, common.HexToAddress(keys.EthereumAddress))
		k.setEthereumOrchestratorAddress(ctx, eth, orch)
	}

	// populate state with cosmos originated denom-erc20 mapping
	for _, item := range data.Erc20ToDenoms {
		k.setCosmosOriginatedDenomToERC20(ctx, uint32(data.Params.BridgeChainId), item.Denom, item.Erc20)
	}

	// reset outgoing txs in state
	for _, ota := range data.OutgoingTxs {
		otx, err := types.UnpackOutgoingTx(ota)
		if err != nil {
			panic(fmt.Sprintf("invalid outgoing tx any in genesis file: %s", err))
		}
		k.SetOutgoingTx(ctx, uint32(data.Params.BridgeChainId), otx)
	}

	// reset signatures in state
	for _, confa := range data.Confirmations {
		conf, err := types.UnpackConfirmation(confa)
		if err != nil {
			panic(fmt.Sprintf("invalid etheruem signature in genesis: %s", err))
		}
		// TODO: not currently an easy way to get the validator address from the
		// etherum address here. once we implement the third index for keys
		// this will be easy.
		k.SetEthereumSignature(ctx, uint32(data.Params.BridgeChainId), conf, sdk.ValAddress{})
	}
}

func chainIDinParamsChainIDs(chainID uint32, chainIDs []uint32) bool {
	for _, id := range chainIDs {
		if id == chainID {
			return true
		}
	}

	return false
}

func InitGenesisMultiChain(ctx sdk.Context, k Keeper, data types.GenesisStateMultiChain) {
	k.setParams(ctx, *data.Params)

	// reset delegate keys in state
	for _, keys := range data.DelegateKeys {
		if err := keys.ValidateBasic(); err != nil {
			panic(fmt.Sprintf("Invalid delegate key in Genesis: %s", err))
		}

		val, _ := sdk.ValAddressFromBech32(keys.ValidatorAddress)
		orch, _ := sdk.AccAddressFromBech32(keys.OrchestratorAddress)
		eth := common.HexToAddress(keys.EthereumAddress)

		// set the orchestrator address
		k.SetOrchestratorValidatorAddress(ctx, val, orch)
		// set the ethereum address
		k.setValidatorEthereumAddress(ctx, val, common.HexToAddress(keys.EthereumAddress))
		k.setEthereumOrchestratorAddress(ctx, eth, orch)
	}

	for _, chainGS := range data.ChainGenesisStates {
		if !chainIDinParamsChainIDs(chainGS.ChainID, data.Params.ChainIds) {
			panic(fmt.Sprintf("chain ID %d presented in state, but not in params", chainGS.ChainID))
		}

		initGenesisForChain(ctx, k, *chainGS)
	}
}

func initGenesisForChain(ctx sdk.Context, k Keeper, data types.ChainGenesisState) {
	// reset pool transactions in state
	for _, tx := range data.UnbatchedSendToEvmTxs {
		k.setUnbatchedSendToEthereum(ctx, data.ChainID, tx)
	}

	// reset ethereum event vote records in state
	for _, evr := range data.EvmEventVoteRecords {
		event, err := types.UnpackEvent(evr.Event)
		if err != nil {
			panic(fmt.Sprintf("couldn't cast to event: %s", err))
		}
		if err := event.Validate(); err != nil {
			panic(fmt.Sprintf("invalid event in genesis: %s", err))
		}
		k.setEthereumEventVoteRecord(ctx, data.ChainID, event.GetEventNonce(), event.Hash(), evr)
	}

	// reset last observed event nonce
	k.setLastObservedEventNonce(ctx, data.ChainID, data.LastObservedEventNonce)

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
		k.setCosmosOriginatedDenomToERC20(ctx, data.ChainID, item.Denom, item.Erc20)
	}

	// reset outgoing txs in state
	for _, ota := range data.OutgoingTxs {
		otx, err := types.UnpackOutgoingTx(ota)
		if err != nil {
			panic(fmt.Sprintf("invalid outgoing tx any in genesis file: %s", err))
		}
		k.SetOutgoingTx(ctx, data.ChainID, otx)
	}

	// reset signatures in state
	for _, confa := range data.Confirmations {
		conf, err := types.UnpackConfirmation(confa)
		if err != nil {
			panic(fmt.Sprintf("invalid etheruem signature in genesis: %s", err))
		}
		// TODO: not currently an easy way to get the validator address from the
		// etherum address here. once we implement the third index for keys
		// this will be easy.
		k.SetEthereumSignature(ctx, data.ChainID, conf, sdk.ValAddress{})
	}
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
//
// Deprecated: use ExportGenesisMultiChain. This should only be used
// if the deprecated InitGenesisState must be used
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var (
		p                        = k.GetParams(ctx)
		outgoingTxs              []*cdctypes.Any
		ethereumTxConfirmations  []*cdctypes.Any
		attmap                   = k.GetEthereumEventVoteRecordMapping(ctx, uint32(p.BridgeChainId))
		ethereumEventVoteRecords []*types.EthereumEventVoteRecord
		delegates                = k.getDelegateKeys(ctx)
		lastobserved             = k.GetLastObservedEventNonce(ctx, uint32(p.BridgeChainId))
		erc20ToDenoms            []*types.ERC20ToDenom
		unbatchedTransfers       = k.getUnbatchedSendToEthereums(ctx, uint32(p.BridgeChainId))
	)

	// export ethereumEventVoteRecords from state
	for _, atts := range attmap {
		// TODO: set height = 0?
		ethereumEventVoteRecords = append(ethereumEventVoteRecords, atts...)
	}

	// export erc20 to denom relations
	k.iterateERC20ToDenom(ctx, uint32(p.BridgeChainId), func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
		erc20ToDenoms = append(erc20ToDenoms, erc20ToDenom)
		return false
	})

	// export signer set txs and sigs
	k.IterateOutgoingTxsByType(ctx, uint32(p.BridgeChainId), types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		ota, _ := types.PackOutgoingTx(otx)
		outgoingTxs = append(outgoingTxs, ota)
		sstx, _ := otx.(*types.SignerSetTx)
		k.iterateEthereumSignatures(ctx, uint32(p.BridgeChainId), sstx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
			sstxc := &types.SignerSetTxConfirmation{SignerSetNonce: sstx.Nonce, EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(), Signature: sig}
			sstxcChainID := &types.SignerSetTxConfirmation_ChainId{ChainId: uint32(p.BridgeChainId)}
			sstxc.XChainId = sstxcChainID
			siga, _ := types.PackConfirmation(sstxc)
			ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
			return false
		})
		return false
	})

	// export batch txs and sigs
	k.IterateOutgoingTxsByType(ctx, uint32(p.BridgeChainId), types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		ota, _ := types.PackOutgoingTx(otx)
		outgoingTxs = append(outgoingTxs, ota)
		btx, _ := otx.(*types.BatchTx)
		k.iterateEthereumSignatures(ctx, uint32(p.BridgeChainId), btx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
			btc := &types.BatchTxConfirmation{TokenContract: btx.TokenContract, BatchNonce: btx.BatchNonce, EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(), Signature: sig}
			btcChainID := &types.BatchTxConfirmation_ChainId{ChainId: uint32(p.BridgeChainId)}
			btc.XChainId = btcChainID
			siga, _ := types.PackConfirmation(btc)
			ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
			return false
		})
		return false
	})

	// export contract call txs and sigs
	k.IterateOutgoingTxsByType(ctx, uint32(p.BridgeChainId), types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		ota, _ := types.PackOutgoingTx(otx)
		outgoingTxs = append(outgoingTxs, ota)
		btx, _ := otx.(*types.ContractCallTx)
		k.iterateEthereumSignatures(ctx, uint32(p.BridgeChainId), btx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
			cctc := &types.ContractCallTxConfirmation{InvalidationScope: btx.InvalidationScope, InvalidationNonce: btx.InvalidationNonce, EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(), Signature: sig}
			cctc.XChainId = &types.ContractCallTxConfirmation_ChainId{ChainId: uint32(p.BridgeChainId)}
			siga, _ := types.PackConfirmation(cctc)
			ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
			return false
		})
		return false
	})

	// this will marshal into "dW51c2Vk" as []byte will be encoded as base64
	for _, delegate := range delegates {
		delegate.EthSignature = []byte("unused")
	}

	return types.GenesisState{
		Params:                     &p,
		LastObservedEventNonce:     lastobserved,
		OutgoingTxs:                outgoingTxs,
		Confirmations:              ethereumTxConfirmations,
		EthereumEventVoteRecords:   ethereumEventVoteRecords,
		DelegateKeys:               delegates,
		Erc20ToDenoms:              erc20ToDenoms,
		UnbatchedSendToEthereumTxs: unbatchedTransfers,
	}
}

// ExportGenesisMultiChain exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesisMultiChain(ctx sdk.Context, k Keeper) types.GenesisStateMultiChain {
	var (
		p         = k.GetParams(ctx)
		delegates = k.getDelegateKeys(ctx)
	)

	var chainGenesisStates []*types.ChainGenesisState

	for _, chainID := range p.ChainIds {
		var (
			outgoingTxs              []*cdctypes.Any
			ethereumTxConfirmations  []*cdctypes.Any
			attmap                   = k.GetEthereumEventVoteRecordMapping(ctx, chainID)
			ethereumEventVoteRecords []*types.EthereumEventVoteRecord
			lastobserved             = k.GetLastObservedEventNonce(ctx, chainID)
			erc20ToDenoms            []*types.ERC20ToDenom
			unbatchedSendToEVMTxs    = k.getUnbatchedSendToEthereums(ctx, chainID)
		)

		// export ethereumEventVoteRecords from state
		for _, atts := range attmap {
			// TODO: set height = 0?
			ethereumEventVoteRecords = append(ethereumEventVoteRecords, atts...)
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
			k.iterateEthereumSignatures(ctx, chainID, sstx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
				siga, _ := types.PackConfirmation(&types.SignerSetTxConfirmation{
					SignerSetNonce: sstx.Nonce,
					EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(),
					Signature:      sig,
					XChainId:       &types.SignerSetTxConfirmation_ChainId{ChainId: chainID},
				})
				ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
				return false
			})
			return false
		})

		// export batch txs and sigs
		k.IterateOutgoingTxsByType(ctx, chainID, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
			ota, _ := types.PackOutgoingTx(otx)
			outgoingTxs = append(outgoingTxs, ota)
			btx, _ := otx.(*types.BatchTx)
			k.iterateEthereumSignatures(ctx, chainID, btx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
				siga, _ := types.PackConfirmation(&types.BatchTxConfirmation{
					TokenContract:  btx.TokenContract,
					BatchNonce:     btx.BatchNonce,
					EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(),
					Signature:      sig,
					XChainId:       &types.BatchTxConfirmation_ChainId{ChainId: chainID},
				})
				ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
				return false
			})
			return false
		})

		// export contract call txs and sigs
		k.IterateOutgoingTxsByType(ctx, chainID, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
			ota, _ := types.PackOutgoingTx(otx)
			outgoingTxs = append(outgoingTxs, ota)
			btx, _ := otx.(*types.ContractCallTx)
			k.iterateEthereumSignatures(ctx, chainID, btx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
				siga, _ := types.PackConfirmation(&types.ContractCallTxConfirmation{
					InvalidationScope: btx.InvalidationScope,
					InvalidationNonce: btx.InvalidationNonce,
					EthereumSigner:    k.GetValidatorEthereumAddress(ctx, val).Hex(),
					Signature:         sig,
					XChainId:          &types.ContractCallTxConfirmation_ChainId{ChainId: chainID},
				})
				ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
				return false
			})
			return false
		})

		cgs := &types.ChainGenesisState{
			ChainID:                chainID,
			LastObservedEventNonce: lastobserved,
			OutgoingTxs:            outgoingTxs,
			Confirmations:          ethereumTxConfirmations,
			EvmEventVoteRecords:    ethereumEventVoteRecords,
			Erc20ToDenoms:          erc20ToDenoms,
			UnbatchedSendToEvmTxs:  unbatchedSendToEVMTxs,
		}

		chainGenesisStates = append(chainGenesisStates, cgs)
	}

	// this will marshal into "dW51c2Vk" as []byte will be encoded as base64
	for _, delegate := range delegates {
		delegate.EthSignature = []byte("unused")
	}

	return types.GenesisStateMultiChain{
		Params:             &p,
		DelegateKeys:       delegates,
		ChainGenesisStates: chainGenesisStates,
	}
}
