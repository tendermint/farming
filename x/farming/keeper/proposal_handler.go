package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming/types"
)

// HandleAddPublicPlanProposal is a handler for executing a public plan creation proposal.
func HandleAddPublicPlanProposal(ctx sdk.Context, k Keeper, p *types.AddPublicPlanProposal) error {
	return nil
}

// HandleUpdatePublicPlanProposal is a handler for executing an update to the public plan.
func HandleUpdatePublicPlanProposal(ctx sdk.Context, k Keeper, p *types.UpdatePublicPlanProposal) error {
	return nil
}

// HandleDeletePublicPlanProposal is a handler for executing a removal of the public plan.
func HandleDeletePublicPlanProposal(ctx sdk.Context, k Keeper, p *types.DeletePublicPlanProposal) error {
	return nil
}
