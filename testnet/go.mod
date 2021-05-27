module github.com/cosmos/gravity-bridge/testnet

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.42.5
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/cosmos/gravity-bridge/module v0.0.0-20210512001249-3c4c5d1dbe2e
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tendermint/tendermint v0.34.10
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/text v0.3.4 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
