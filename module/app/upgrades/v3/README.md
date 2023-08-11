# v3 upgrade

This upgrade moves the gravity module from consensus version 2 to 3.

## Summary of changes

* Add CompletedOutgoingTx store for marking transactions as executed
* Refactor slashing logic and include CompletedOutgoingTx in unslashed tx getter
* Add tx confirmation pruning 
* Add event vote record pruning
* Fix bug that iterated the entire key store when SetDelegateKeys was called
* Refactor address lookups used in SetDelegateKeys to not require scanning entire list of validators
* Remove MsgRequestBatchTx and handlers
* Add missing Amino registrations
* Fix minor CLI bugs
* Improve and correct terminology in function names and comments

