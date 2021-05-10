# Orchestrator folder

### client/

This folder builds a binary that is a client application for the gravity system. It contains the following commands:
- `cosmos-to-eth`
- `eth-to-cosmos`
- `deploy-erc20-representation`

### cosmos_gravity/

This is a library for interacting with the cosmos chain both queries and transactions. It substantally wraps `gravity_proto`.

### ethereum_gravity/

This is a library that contains code for the interactions with the counterparty ethereum chain.

### gravity_proto/

`prost` generated bindings for working with the gravity protobuf objects.

### gravity_utils/

Various utilities for working with the `gravity` code.

### orchestrator/

The package to build the orchestartor binary.

### proto_build/

Run `cargo run` in this folder to build `gravity_proto` also note, this will generate too many files. Only `gravity.v1.rs` is required.

### register_delegate_keys/

This is a sepreate binary for running a command to register delegate keys for a validator. NOTE: this needs to be done in `gentx` now so this is likely no longer needed.

### relayer/

This is to build the relayer logic (i.e. cosmos to ethereum) as a seperate binary. It also contains the library for the relayer.

### scripts/

Supporting bash scripts for this library

### test_runner/

A binary which runs tests against a cosmos chain