package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming/types"
)

// HandleAddPublicPlanProposal is a handler for executing a public plan creation proposal.
func HandleAddPublicPlanProposal(ctx sdk.Context, k Keeper, publicPlanProposal *types.PublicPlanProposal) error {
	// TODO: not implemented yet
	return nil
}
