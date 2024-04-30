package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/errors"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	ccodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	icaexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	"github.com/gorilla/mux"
	gravityparams "github.com/peggyjv/gravity-bridge/module/v4/app/params"
	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity"
	gravityclient "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/client"
	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity/keeper"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
)

const (
	appName = "app"

	// MaxAddrLen is the maximum allowed length (in bytes) for an address.
	//
	// NOTE: In the SDK, the default value is 255.
	MaxAddrLen = 20
)

var (
	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome string

	// ModuleBasics The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		consensus.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				upgradeclient.LegacyProposalHandler,
				upgradeclient.LegacyCancelProposalHandler,
				gravityclient.ProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		gravity.AppModuleBasic{},
	)

	// module account permissions
	// NOTE: We believe that this is giving various modules access to functions of the supply module? We will probably need to use this.
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		gravitytypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}

	// verify app interface at compile time
	_ runtime.AppI            = (*Gravity)(nil)
	_ servertypes.Application = (*Gravity)(nil)
)

// MakeCodec creates the application codec. The codec is sealed before it is
// returned.
func MakeCodec() *codec.LegacyAmino {
	var cdc = codec.NewLegacyAmino()
	ModuleBasics.RegisterLegacyAminoCodec(cdc)
	vesting.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)
	ccodec.RegisterCrypto(cdc)
	cdc.Seal()
	return cdc
}

// Gravity extended ABCI application
type Gravity struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tKeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	accountKeeper         authkeeper.AccountKeeper
	bankKeeper            bankkeeper.Keeper
	capabilityKeeper      *capabilitykeeper.Keeper
	stakingKeeper         stakingkeeper.Keeper
	slashingKeeper        slashingkeeper.Keeper
	mintKeeper            mintkeeper.Keeper
	distrKeeper           distrkeeper.Keeper
	govKeeper             govkeeper.Keeper
	consensusParamsKeeper consensusparamkeeper.Keeper
	crisisKeeper          crisiskeeper.Keeper
	upgradeKeeper         upgradekeeper.Keeper
	paramsKeeper          paramskeeper.Keeper
	ibcKeeper             *ibckeeper.Keeper
	evidenceKeeper        evidencekeeper.Keeper
	transferKeeper        ibctransferkeeper.Keeper
	gravityKeeper         keeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	// Module Manager
	mm *module.Manager

	// configurator
	configurator module.Configurator

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".gravity")
}

func NewGravityApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig gravityparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *Gravity {

	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, icaexported.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, consensusparamtypes.StoreKey, capabilitytypes.StoreKey,
		gravitytypes.StoreKey, crisistypes.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	var app = &Gravity{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tKeys:             tKeys,
		memKeys:           memKeys,
	}

	app.paramsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey])

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// set the BaseApp's parameter store
	app.consensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, keys[consensusparamtypes.StoreKey], authority)
	bApp.SetParamStore(&app.consensusParamsKeeper)

	app.capabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)
	scopedIBCKeeper := app.capabilityKeeper.ScopeToModule(icaexported.ModuleName)
	scopedTransferKeeper := app.capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should
	// call `Seal` after creating their scoped modules in the app via
	// `ScopeToModule`.
	app.capabilityKeeper.Seal()

	app.accountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		app.keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		gravityparams.Bech32PrefixAccAddr,
		authority,
	)
	app.bankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		app.keys[banktypes.StoreKey],
		app.accountKeeper,
		app.BlockedAddrs(),
		authority,
	)
	app.stakingKeeper = *stakingkeeper.NewKeeper(
		appCodec,
		app.keys[stakingtypes.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		authority,
	)

	app.mintKeeper = mintkeeper.NewKeeper(
		appCodec,
		app.keys[minttypes.StoreKey],
		app.stakingKeeper,
		app.accountKeeper,
		app.bankKeeper,
		authtypes.FeeCollectorName,
		authority,
	)

	app.distrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		app.stakingKeeper,
		authtypes.FeeCollectorName,
		authority,
	)

	app.slashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		keys[slashingtypes.StoreKey],
		app.stakingKeeper,
		authority,
	)

	app.crisisKeeper = *crisiskeeper.NewKeeper(
		appCodec,
		app.keys[crisistypes.StoreKey],
		invCheckPeriod,
		app.bankKeeper,
		authtypes.FeeCollectorName,
		authority,
	)

	app.upgradeKeeper = *upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authority,
	)

	app.gravityKeeper = keeper.NewKeeper(
		appCodec,
		keys[gravitytypes.StoreKey],
		app.GetSubspace(gravitytypes.ModuleName),
		app.accountKeeper,
		app.stakingKeeper,
		app.bankKeeper,
		app.slashingKeeper,
		app.distrKeeper,
		sdk.DefaultPowerReduction,
		app.ModuleAccountAddressesToNames([]string{}),
		app.ModuleAccountAddressesToNames([]string{distrtypes.ModuleName}),
	)

	app.stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks(),
			app.gravityKeeper.Hooks(),
		),
	)

	app.ibcKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[icaexported.StoreKey],
		app.GetSubspace(icaexported.ModuleName),
		app.stakingKeeper,
		app.upgradeKeeper,
		scopedIBCKeeper,
	)

	app.transferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.ibcKeeper.ChannelKeeper, app.ibcKeeper.ChannelKeeper, &app.ibcKeeper.PortKeeper,
		app.accountKeeper, app.bankKeeper, scopedTransferKeeper,
	)
	transferModule := ibctransfer.NewAppModule(app.transferKeeper)
	transferIBCModule := ibctransfer.NewIBCModule(app.transferKeeper)

	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	app.ibcKeeper.SetRouter(ibcRouter)

	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		app.stakingKeeper,
		app.slashingKeeper,
	)
	app.evidenceKeeper = *evidenceKeeper

	govRouter := govtypesv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypesv1beta1.ProposalHandler).
		AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(&app.upgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.ibcKeeper.ClientKeeper)).
		AddRoute(gravitytypes.RouterKey, gravity.NewCommunityPoolEthereumSpendProposalHandler(app.gravityKeeper))

	// Are defaults ok here?
	govConfig := govtypes.DefaultConfig()
	app.govKeeper = *govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		app.stakingKeeper,
		// TODO: This does not include our custom module proposal handlers anymore since the govRouter parameter was removed from this constructor
		app.MsgServiceRouter(),
		govConfig,
		authority,
	)

	app.govKeeper.SetLegacyRouter(govRouter)

	app.setupUpgradeStoreLoaders()

	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.accountKeeper,
			app.stakingKeeper,
			app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.accountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(
			app.accountKeeper,
			app.bankKeeper,
		),
		bank.NewAppModule(
			appCodec,
			app.bankKeeper,
			app.accountKeeper,
			app.GetSubspace(banktypes.ModuleName),
		),
		capability.NewAppModule(
			appCodec,
			*app.capabilityKeeper,
			false,
		),
		crisis.NewAppModule(
			&app.crisisKeeper,
			skipGenesisInvariants,
			app.GetSubspace(crisistypes.ModuleName),
		),
		gov.NewAppModule(
			appCodec,
			&app.govKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.GetSubspace(govtypes.ModuleName),
		),
		mint.NewAppModule(
			appCodec,
			app.mintKeeper,
			app.accountKeeper,
			nil,
			app.GetSubspace(minttypes.ModuleName),
		),
		slashing.NewAppModule(
			appCodec,
			app.slashingKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.stakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
		),
		distr.NewAppModule(
			appCodec,
			app.distrKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.stakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),
		staking.NewAppModule(appCodec,
			&app.stakingKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		upgrade.NewAppModule(&app.upgradeKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		params.NewAppModule(app.paramsKeeper),
		transferModule,
		gravity.NewAppModule(
			app.gravityKeeper,
			app.bankKeeper,
		),
	)

	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		icaexported.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		gravitytypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		icaexported.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		gravitytypes.ModuleName,
	)
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		icaexported.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		gravitytypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)
	app.setupUpgradeHandlers()

	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.accountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(
			appCodec,
			app.bankKeeper,
			app.accountKeeper,
			app.GetSubspace(banktypes.ModuleName),
		),
		capability.NewAppModule(
			appCodec,
			*app.capabilityKeeper,
			false,
		),
		gov.NewAppModule(
			appCodec,
			&app.govKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.GetSubspace(govtypes.ModuleName),
		),
		mint.NewAppModule(
			appCodec,
			app.mintKeeper,
			app.accountKeeper,
			nil,
			app.GetSubspace(minttypes.ModuleName),
		),
		slashing.NewAppModule(
			appCodec,
			app.slashingKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.stakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
		),
		distr.NewAppModule(
			appCodec,
			app.distrKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.stakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),
		staking.NewAppModule(appCodec,
			&app.stakingKeeper,
			app.accountKeeper,
			app.bankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		evidence.NewAppModule(app.evidenceKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		params.NewAppModule(app.paramsKeeper),
		transferModule,
		gravity.NewAppModule(
			app.gravityKeeper,
			app.bankKeeper,
		),
	)

	app.sm.RegisterStoreDecoders()

	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)
	app.MountMemoryStores(memKeys)

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.accountKeeper,
			BankKeeper:      app.bankKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			FeegrantKeeper:  nil,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create ante handler: %s", err))
	}

	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	return app
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// simapp. It is useful for tests and clients who do not want to construct the
// full simapp
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}

// Name returns the name of the App
func (app *Gravity) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *Gravity) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *Gravity) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *Gravity) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.upgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *Gravity) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *Gravity) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// ModuleAccountNames returns a map of module account address to module name
func (app *Gravity) ModuleAccountAddressesToNames(moduleAccounts []string) map[string]string {
	modAccNames := make(map[string]string)
	for _, acc := range moduleAccounts {
		modAccNames[authtypes.NewModuleAddress(acc).String()] = acc
	}

	return modAccNames
}

func (app *Gravity) ModuleNames() []string {
	return app.mm.ModuleNames()
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *Gravity) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Gravity) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns SimApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Gravity) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns SimApp's InterfaceRegistry
func (app *Gravity) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Gravity) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Gravity) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tKeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *Gravity) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
func (app *Gravity) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *Gravity) SimulationManager() *module.SimulationManager {
	return app.sm
}

// API server.
func (app *Gravity) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// TODO: build the custom gravity swagger files and add here?
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterNodeService registers the node gRPC Query service.
func (app *Gravity) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// RegisterSwaggerAPI registers swagger route with API Server
// TODO: build the custom gravity swagger files and add here?
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *Gravity) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *Gravity) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}

	return modAccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypesv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(gravitytypes.ModuleName)
	paramsKeeper.Subspace(icaexported.ModuleName)

	return paramsKeeper
}

func VerifyAddressFormat(bz []byte) error {
	if len(bz) == 0 {
		return errors.Wrap(sdkerrors.ErrUnknownAddress, "invalid address; cannot be empty")
	}
	if len(bz) != MaxAddrLen {
		return errors.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"invalid address length; got: %d, max: %d", len(bz), MaxAddrLen,
		)
	}

	return nil
}

func (app *Gravity) setupUpgradeStoreLoaders() {
	_, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	// if upgradeInfo.Name matches a plan name with a module being added, renamed, or deleted,
	// create a storetypes.StoreUpgrades struct and
	// app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	// see also:
	// https://github.com/cosmos/cosmos-sdk/blob/master/docs/core/upgrade.md#add-storeupgrades-for-new-modules
}

func (app *Gravity) setupUpgradeHandlers() {

}
