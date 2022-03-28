package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
)

func (k Keeper) HandleCommunityPoolEthereumSpendProposal(ctx sdk.Context, p *types.CommunityPoolEthereumSpendProposal) error {
	// TODO(bolten): is this implicitly called elsewhere?
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	feePool := k.DistributionKeeper.GetFeePool(ctx)

	// NOTE the community pool isn't a module account, however its coins
	// are held in the distribution module account. Thus the community pool
	// must be reduced separately from the createSendToEthereum calls
	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(p.Amount...))
	if negative {
		return distributiontypes.ErrBadDistribution
	}

	feePool.CommunityPool = newPool
	sender := authtypes.NewModuleAddress(distributiontypes.ModuleName)

	for _, coin := range p.Amount {
		// TODO(bolten): currently just setting fees for these to zero. Given that an individual
		// createSendToEthereum could fail  and thus the proposal would end up in a half-complete state,
		// should we restrict proposals to only allowing one denom?
		_, err := k.createSendToEthereum(ctx, sender, p.Recipient, coin, sdk.NewCoin(coin.Denom, sdk.NewInt(0)))

		if err != nil {
			return err
		}
	}

	k.DistributionKeeper.SetFeePool(ctx, feePool)

	return nil
}
