package keeper

import (
	"fmt"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
)

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, *data.Params)

	// reset pool transactions in state
	for _, tx := range data.UnbatchedSendToEthereumTxs {
		k.setUnbatchedSendToEthereum(ctx, tx)
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
		ctx.Logger().Info("gravity: event hash", "hash", event.Hash(), "nonce", event.GetEventNonce())
		k.setEthereumEventVoteRecord(ctx, event.GetEventNonce(), event.Hash(), evr)
	}

	// reset last observed event nonce
	k.setLastObservedEventNonce(ctx, data.LastObservedEventNonce)

	// reset attestation state of all validators
	for _, eventVoteRecord := range data.EthereumEventVoteRecords {
		event, _ := types.UnpackEvent(eventVoteRecord.Event)
		for _, vote := range eventVoteRecord.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := k.getLastEventNonceByValidator(ctx, val)
			if event.GetEventNonce() > last {
				k.setLastEventNonceByValidator(ctx, val, event.GetEventNonce())
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
		k.setCosmosOriginatedDenomToERC20(ctx, item.Denom, common.HexToAddress(item.Erc20))
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
			panic(fmt.Sprintf("invalid etheruem signature in genesis: %s", err))
		}
		// TODO: not currently an easy way to get the validator address from the
		// etherum address here. once we implement the third index for keys
		// this will be easy.
		k.SetEthereumSignature(ctx, conf, sdk.ValAddress{})
	}
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var (
		p                        = k.GetParams(ctx)
		outgoingTxs              []*cdctypes.Any
		ethereumTxConfirmations  []*cdctypes.Any
		attmap                   = k.GetEthereumEventVoteRecordMapping(ctx)
		ethereumEventVoteRecords []*types.EthereumEventVoteRecord
		delegates                = k.getDelegateKeys(ctx)
		lastobserved             = k.GetLastObservedEventNonce(ctx)
		erc20ToDenoms            []*types.ERC20ToDenom
		unbatchedTransfers       = k.getUnbatchedSendToEthereums(ctx)
	)

	// export ethereumEventVoteRecords from state
	for _, atts := range attmap {
		// TODO: set height = 0?
		ethereumEventVoteRecords = append(ethereumEventVoteRecords, atts...)
	}

	// export erc20 to denom relations
	k.iterateERC20ToDenom(ctx, func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
		erc20ToDenoms = append(erc20ToDenoms, erc20ToDenom)
		return false
	})

	// export signer set txs and sigs
	k.IterateOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		ota, _ := types.PackOutgoingTx(otx)
		outgoingTxs = append(outgoingTxs, ota)
		sstx, _ := otx.(*types.SignerSetTx)
		k.iterateEthereumSignatures(ctx, sstx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
			siga, _ := types.PackConfirmation(&types.SignerSetTxConfirmation{
				SignerSetNonce: sstx.Nonce,
				EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(),
				Signature:      sig,
			})
			ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
			return false
		})
		return false
	})

	// export batch txs and sigs
	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		ota, _ := types.PackOutgoingTx(otx)
		outgoingTxs = append(outgoingTxs, ota)
		btx, _ := otx.(*types.BatchTx)
		k.iterateEthereumSignatures(ctx, btx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
			siga, _ := types.PackConfirmation(&types.BatchTxConfirmation{
				TokenContract:  btx.TokenContract,
				BatchNonce:     btx.BatchNonce,
				EthereumSigner: k.GetValidatorEthereumAddress(ctx, val).Hex(),
				Signature:      sig,
			})
			ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
			return false
		})
		return false
	})

	// export contract call txs and sigs
	k.IterateOutgoingTxsByType(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		ota, _ := types.PackOutgoingTx(otx)
		outgoingTxs = append(outgoingTxs, ota)
		btx, _ := otx.(*types.ContractCallTx)
		k.iterateEthereumSignatures(ctx, btx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
			siga, _ := types.PackConfirmation(&types.ContractCallTxConfirmation{
				InvalidationScope: btx.InvalidationScope,
				InvalidationNonce: btx.InvalidationNonce,
				EthereumSigner:    k.GetValidatorEthereumAddress(ctx, val).Hex(),
				Signature:         sig,
			})
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
