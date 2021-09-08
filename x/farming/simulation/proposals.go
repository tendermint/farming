package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/tendermint/farming/app/params"
	"github.com/tendermint/farming/x/farming/keeper"
	"github.com/tendermint/farming/x/farming/types"
)

// Simulation operation weights constants.
const OpWeightSimulatePublicPlanProposal = "op_weight_public_plan_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightSimulatePublicPlanProposal,
			params.DefaultWeightPublicPlanProposal,
			SimulatePublicPlanProposal(k),
		),
	}
}

// SimulatePublicPlanProposal generates random public plan proposal content
func SimulatePublicPlanProposal(k keeper.Keeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		// simAccount, _ := simtypes.RandomAcc(r, accs)

		return types.NewPublicPlanProposal(
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 100),
			[]*types.AddRequestProposal{},
			[]*types.UpdateRequestProposal{},
			[]*types.DeleteRequestProposal{},
		)
	}
}
