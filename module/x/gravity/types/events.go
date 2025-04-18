package types

const (
	EventTypeObservation              = "observation"
	EventTypeOutgoingBatch            = "outgoing_batch"
	EventTypeMultisigUpdateRequest    = "multisig_update_request"
	EventTypeOutgoingBatchCanceled    = "outgoing_batch_canceled"
	EventTypeContractCallTxCanceled   = "outgoing_logic_call_canceled"
	EventTypeBridgeWithdrawalReceived = "withdrawal_received"
	EventTypeBridgeDepositReceived    = "deposit_received"
	EventTypeBridgeWithdrawCanceled   = "withdraw_canceled"
	EventTypeContractCallTxCompleted  = "contract_call_tx_completed"

	AttributeKeyEthereumEventVoteRecordID     = "ethereum_event_vote_record_id"
	AttributeKeyBatchConfirmKey               = "batch_confirm_key"
	AttributeKeyEthereumSignatureKey          = "ethereum_signature_key"
	AttributeKeyOutgoingBatchID               = "batch_id"
	AttributeKeyOutgoingTXID                  = "outgoing_tx_id"
	AttributeKeyEthereumEventType             = "ethereum_event_type"
	AttributeKeyContract                      = "bridge_contract"
	AttributeKeyNonce                         = "nonce"
	AttributeKeySignerSetNonce                = "signerset_nonce"
	AttributeKeyBatchNonce                    = "batch_nonce"
	AttributeKeyBridgeChainID                 = "bridge_chain_id"
	AttributeKeySetOrchestratorAddr           = "set_orchestrator_address"
	AttributeKeySetEthereumAddr               = "set_ethereum_address"
	AttributeKeyValidatorAddr                 = "validator_address"
	AttributeKeyContractCallInvalidationScope = "contract_call_invalidation_scope"
	AttributeKeyContractCallInvalidationNonce = "contract_call_invalidation_nonce"
	AttributeKeyContractCallPayload           = "contract_call_payload"
	AttributeKeyContractCallTokens            = "contract_call_tokens"
	AttributeKeyContractCallFees              = "contract_call_fees"
	AttributeKeyContractCallAddress           = "contract_call_address"
	AttributeKeyEthTxTimeout                  = "eth_tx_timeout"

	// slashing reasons
	AttributeMissingSignerSetSignature = "missing_signer_set_signature"
	AttributeMissingSignature          = "missing_signature"
)
