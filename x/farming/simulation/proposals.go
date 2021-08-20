package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tendermint/farming/x/farming/keeper"
)

// OpWeightSubmitCommunitySpendProposal app params key for community spend proposal
const OpWeightSubmitCommunitySpendProposal = "op_weight_submit_community_spend_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []simtypes.WeightedProposalContent {
	// TODO: not implemented yet
	return nil
	// return []simtypes.WeightedProposalContent{
	// 	simulation.NewWeightedProposalContent(
	// 		OpWeightSubmitCommunitySpendProposal,
	// 		simappparams.DefaultWeightCommunitySpendProposal,
	// 		SimulateCommunityPoolSpendProposalContent(k),
	// 	),
	// }
}

// SimulateCommunityPoolSpendProposalContent generates random community-pool-spend proposal content
func SimulateCommunityPoolSpendProposalContent(k keeper.Keeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		// TODO: not implemented yet
		return nil
	}
}
