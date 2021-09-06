package simulation

import (
	"math/rand"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/tendermint/farming/x/farming/keeper"
)

// Simulation operation weights constants.
const (
	OpWeightSubmitPublicFixedAmountPlanProposal = "op_weight_submit_public_fixed_amount_plan_proposal"
	OpWeightSubmitPublicRatioPlanProposal       = "op_weight_submit_public_ratio_plan_proposal"
)

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightSubmitPublicFixedAmountPlanProposal,
			simappparams.DefaultWeightTextProposal,
			SimulateTextProposalContent,
		),
		// simulation.NewWeightedProposalContent(
		// 	OpWeightSubmitPublicRatioPlanProposal,
		// 	simappparams.DefaultWeightTextProposal,
		// 	SimulateTextProposalContent(k),
		// ),
	}
}

// SimulateTextProposalContent returns a random text proposal content.
func SimulateTextProposalContent(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) simtypes.Content {
	return govtypes.NewTextProposal(
		simtypes.RandStringOfLength(r, 140),
		simtypes.RandStringOfLength(r, 5000),
	)
}
