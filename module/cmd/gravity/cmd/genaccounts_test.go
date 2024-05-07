package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	tmjson "github.com/cometbft/cometbft/libs/json"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltest "github.com/cosmos/cosmos-sdk/x/genutil/client/testutil"
)

var testMbm = module.NewBasicManager(genutil.AppModuleBasic{})

func TestAddGenesisAccountCmd(t *testing.T) {
	_, _, addr1 := testdata.KeyTestPubAddr()
	tests := []struct {
		name      string
		addr      string
		denom     string
		expectErr bool
	}{
		{
			name:      "invalid address",
			addr:      "",
			denom:     "1000atom",
			expectErr: true,
		},
		{
			name:      "valid address",
			addr:      addr1.String(),
			denom:     "1000atom",
			expectErr: false,
		},
		{
			name:      "multiple denoms",
			addr:      addr1.String(),
			denom:     "1000atom,2000stake",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			logger := log.NewNopLogger()
			cfg, err := genutiltest.CreateDefaultTendermintConfig(home)
			require.NoError(t, err)

			appCodec := moduletestutil.MakeTestEncodingConfig(auth.AppModuleBasic{}).Codec
			err = genutiltest.ExecInitCmd(testMbm, home, appCodec)
			require.NoError(t, err)

			serverCtx := server.NewContext(viper.New(), cfg, logger)
			clientCtx := client.Context{}.WithCodec(appCodec).WithHomeDir(home)

			ctx := context.Background()
			ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
			ctx = context.WithValue(ctx, server.ServerContextKey, serverCtx)

			cmd := AddGenesisAccountCmd(home)
			cmd.SetArgs([]string{
				tc.addr,
				tc.denom,
				fmt.Sprintf("--%s=home", flags.FlagHome),
			})

			if tc.expectErr {
				require.Error(t, cmd.ExecuteContext(ctx))
			} else {
				require.NoError(t, cmd.ExecuteContext(ctx))

				genFile := serverCtx.Config.GenesisFile()
				bytes, err := os.ReadFile(genFile)
				require.NoError(t, err)

				var genDoc tmtypes.GenesisDoc
				tmjson.Unmarshal(bytes, &genDoc)

				var appState map[string]json.RawMessage
				err = json.Unmarshal(genDoc.AppState, &appState)
				require.NoError(t, err)

				var bankGenState banktypes.GenesisState
				bankGenStateBz := appState[banktypes.ModuleName]
				clientCtx.Codec.MustUnmarshalJSON(bankGenStateBz, &bankGenState)

				require.Equal(t, bankGenState.Supply.String(), tc.denom)
			}
		})
	}
}
