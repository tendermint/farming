package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tendermint/farming/x/farming/keeper"
)

// OpWeightSubmitPublicPlanProposals app params key for public plan proposals
const OpWeightSubmitPublicPlanProposals = "op_weight_submit_public_plan_proposals"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []simtypes.WeightedProposalContent {
	// TODO: not implemented yet
	return nil
	// return []simtypes.WeightedProposalContent{
	// 	simulation.NewWeightedProposalContent(
	// 		OpWeightSubmitPublicPlanProposals,
	// 		simappparams.DefaultWeightPublicPlanProposals,
	// 		SimulateCommunityPoolSpendProposalContent(k),
	// 	),
	// }
}

// SimulatePublicPlanProposalsProposalContent generates random public plan proposals proposal content
func SimulatePublicPlanProposalsProposalContent(k keeper.Keeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		// TODO: not implemented yet
		return nil
	}
}
