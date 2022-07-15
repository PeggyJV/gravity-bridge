module github.com/peggyjv/gravity-bridge/module/v2

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.46.0-rc2
	github.com/cosmos/ibc-go/v4 v4.0.0-rc0
	github.com/ethereum/go-ethereum v1.10.17
	github.com/gogo/protobuf v1.3.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.12.0
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/tendermint v0.34.20-rc0
	github.com/tendermint/tm-db v0.6.6
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd
	google.golang.org/grpc v1.47.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76

replace github.com/cosmos/ibc-go/v4 => github.com/notional-labs/ibc-go/v4 v4.0.0-20220705111408-32fa3063ed41
