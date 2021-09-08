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
func ProposalContents(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightSimulatePublicPlanProposal,
			params.DefaultWeightPublicPlanProposal,
			SimulatePublicPlanProposal(ak, bk, k),
		),
	}
}

// SimulatePublicPlanProposal generates random public plan proposal content
func SimulatePublicPlanProposal(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params := k.GetParams(ctx)
		_, hasNeg := spendable.SafeSub(params.PrivatePlanCreationFee)
		if hasNeg {
			return nil
		}

		poolCoins, err := mintPoolCoins(ctx, r, bk, simAccount)
		if err != nil {
			return nil
		}

		// create add request proposal
		// TODO: randomized values of the fields
		req := &types.AddRequestProposal{
			Name:               "simulation-test-" + simtypes.RandStringOfLength(r, 5),
			FarmingPoolAddress: simAccount.Address.String(),
			TerminationAddress: simAccount.Address.String(),
			StakingCoinWeights: sdk.NewDecCoins(sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 1)),
			StartTime:          ctx.BlockTime(),
			EndTime:            ctx.BlockTime().AddDate(0, 1, 0),
			EpochAmount:        sdk.NewCoins(sdk.NewInt64Coin(poolCoins[r.Intn(3)].Denom, int64(simtypes.RandIntBetween(r, 10_000_000, 1_000_000_000)))),
			EpochRatio:         sdk.ZeroDec(),
		}
		addRequests := []*types.AddRequestProposal{req}

		// TODO
		// create update request proposal
		// updating plan can only be allowed (owner)

		// TODO
		// create delete request proposal
		// deleting plan can only be allowed

		return types.NewPublicPlanProposal(
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 100),
			addRequests,
			[]*types.UpdateRequestProposal{},
			[]*types.DeleteRequestProposal{},
		)
	}
}
