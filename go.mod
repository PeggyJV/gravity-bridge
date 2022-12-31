module github.com/peggyjv/gravity-bridge

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.44.5-patch
	github.com/cosmos/go-bip39 v1.0.0
	github.com/ethereum/go-ethereum v1.10.17
	github.com/icza/dyno v0.0.0-20220812133438-f0b6f8a18845 // indirect
	github.com/miguelmota/go-ethereum-hdwallet v0.1.1
	github.com/ory/dockertest/v3 v3.9.1
	github.com/peggyjv/gravity-bridge/module/v3 v3.0.0-20221121235217-bc73418ac79b // indirect
	github.com/spf13/viper v1.13.0
	github.com/strangelove-ventures/ibctest/v6 v6.0.0 // indirect
	github.com/stretchr/testify v1.8.1
	github.com/tendermint/tendermint v0.34.21
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0

replace github.com/ChainSafe/go-schnorrkel => github.com/ChainSafe/go-schnorrkel v0.0.0-20200405005733-88cbf1b4c40d

replace github.com/ChainSafe/go-schnorrkel/1 => github.com/ChainSafe/go-schnorrkel v1.0.0

//replace github.com/strangelove-ventures/ibctest/v6 => github.com/PeggyJV/ibctest/v6 v6.0.0-20221227111146-1e7c3e8e79b0
replace github.com/strangelove-ventures/ibctest/v6 => /Users/mv/Development/somm/ibctest