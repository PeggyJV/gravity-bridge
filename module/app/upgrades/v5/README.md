# v5 upgrade

This upgrade simply fixes the bug in the `SubmitEthereumTxConfirmation` message handler where `CompletedOutgoingTx`s were not checked and therefore could not be signed by orchestrators that failed to sign before the transaction was completed, resulting in jailing and slashing. Because this is a state breaking change, we bump the consensus version of the module.

