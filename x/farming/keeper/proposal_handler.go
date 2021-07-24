package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming/types"
)

// HandlePublicPlanProposal is a handler for executing a public plan creation proposal.
func HandlePublicPlanProposal(ctx sdk.Context, k Keeper, publicPlanProposal *types.PublicPlanProposal) error {
	// TODO: not implemented yet
	return nil
}
