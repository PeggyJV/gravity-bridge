# v2 upgrade

This upgrade moves the gravity module from consensus version 1 to 2.

## Summary of changes

* Switch to use of the upgrade module
* Support sending to specified module accounts over the bridge
* Community spend governance proposal for sending funds over the bridge to an Ethereum address
* Fix a bug affecting the capitalization of ERC20 addresses in denominations
* Fix a bug incorrectly setting the timeouts of ContractCallTxs