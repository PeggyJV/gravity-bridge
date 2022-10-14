module github.com/peggyjv/gravity-bridge

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.44.5-patch
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/ibc-go v1.0.1 // indirect
	github.com/ethereum/go-ethereum v1.10.11
	github.com/miguelmota/go-ethereum-hdwallet v0.1.1
	github.com/ory/dockertest/v3 v3.8.1
	github.com/peggyjv/gravity-bridge/module/v2 v2.0.2 // indirect
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/tendermint v0.34.14
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0
