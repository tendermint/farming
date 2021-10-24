package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/farming/x/farming/types"
)

// HandlePublicPlanProposal is a handler for executing a public plan creation proposal.
func HandlePublicPlanProposal(ctx sdk.Context, k Keeper, proposal *types.PublicPlanProposal) error {
	if proposal.AddRequestProposals != nil {
		if err := k.AddPublicPlanProposal(ctx, proposal.AddRequestProposals); err != nil {
			return err
		}
	}

	if proposal.UpdateRequestProposals != nil {
		if err := k.UpdatePublicPlanProposal(ctx, proposal.UpdateRequestProposals); err != nil {
			return err
		}
	}

	if proposal.DeleteRequestProposals != nil {
		if err := k.DeletePublicPlanProposal(ctx, proposal.DeleteRequestProposals); err != nil {
			return err
		}
	}

	plans := k.GetPlans(ctx)
	if err := types.ValidateTotalEpochRatio(plans); err != nil {
		return err
	}

	return nil
}

// AddPublicPlanProposal adds a new public plan once the governance proposal is passed.
func (k Keeper) AddPublicPlanProposal(ctx sdk.Context, proposals []*types.AddRequestProposal) error {
	for _, p := range proposals {
		farmingPoolAddrAcc, err := sdk.AccAddressFromBech32(p.GetFarmingPoolAddress())
		if err != nil {
			return err
		}
		// ^ Pick either farmingPoolAddr or farmingPoolAcc throughout this file :)

		terminationAcc, err := sdk.AccAddressFromBech32(p.GetTerminationAddress())
		if err != nil {
			return err
		}

		if p.IsForFixedAmountPlan() {
			msg := types.NewMsgCreateFixedAmountPlan(
				p.GetName(),
				farmingPoolAddrAcc,
				p.GetStakingCoinWeights(),
				p.GetStartTime(),
				p.GetEndTime(),
				p.EpochAmount,
			)

			plan, err := k.CreateFixedAmountPlan(ctx, msg, farmingPoolAddrAcc, terminationAcc, types.PlanTypePublic)
			if err != nil {
				return err
			}

			logger := k.Logger(ctx)
			logger.Info("created public fixed amount plan", "fixed_amount_plan", plan)
		} else {
			msg := types.NewMsgCreateRatioPlan(
				p.GetName(),
				farmingPoolAddrAcc,
				p.GetStakingCoinWeights(),
				p.GetStartTime(),
				p.GetEndTime(),
				p.EpochRatio,
			)

			plan, err := k.CreateRatioPlan(ctx, msg, farmingPoolAddrAcc, terminationAcc, types.PlanTypePublic)
			if err != nil {
				return err
			}

			logger := k.Logger(ctx)
			logger.Info("created public ratio amount plan", "ratio_plan", plan)
		}
	}

	return nil
}

// UpdatePublicPlanProposal overwrites the plan with the new plan proposal once the governance proposal is passed.
func (k Keeper) UpdatePublicPlanProposal(ctx sdk.Context, proposals []*types.UpdateRequestProposal) error {
	for _, p := range proposals {
		plan, found := k.GetPlan(ctx, p.GetPlanId())
		if !found {
			return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "plan %d is not found", p.GetPlanId())
		}

		if p.GetName() != "" {
			if err := plan.SetName(p.GetName()); err != nil {
				return err
			}
		}

		if p.GetFarmingPoolAddress() != "" {
			farmingPoolAddrAcc, err := sdk.AccAddressFromBech32(p.GetFarmingPoolAddress())
			if err != nil {
				return err
			}
			if err := plan.SetFarmingPoolAddress(farmingPoolAddrAcc); err != nil {
				return err
			}
			// ^ Pick either farmingPoolAddr or farmingPoolAcc :)
		}

		if p.GetTerminationAddress() != "" {
			terminationAddrAcc, err := sdk.AccAddressFromBech32(p.GetTerminationAddress())
			if err != nil {
				return err
			}
			if err := plan.SetTerminationAddress(terminationAddrAcc); err != nil {
				return err
			}
			// ^ Pick either terminationAddr or terminationAcc :)
		}

		if p.GetStakingCoinWeights() != nil {
			if err := plan.SetStakingCoinWeights(p.GetStakingCoinWeights()); err != nil {
				return err
			}
		}

		if p.GetStartTime() != nil {
			if err := plan.SetStartTime(*p.GetStartTime()); err != nil {
				return err
			}
		}

		if p.GetEndTime() != nil {
			if err := plan.SetEndTime(*p.GetEndTime()); err != nil {
				return err
			}
		}

		if p.IsForFixedAmountPlan() {
			// change the plan to fixed amount plan
			plan = types.NewFixedAmountPlan(plan.GetBasePlan(), p.GetEpochAmount())

			logger := k.Logger(ctx)
			logger.Info("updated public fixed amount plan", "fixed_amount_plan", plan)

		} else if p.IsForRatioPlan() {
			// change the plan to ratio plan
			plan = types.NewRatioPlan(plan.GetBasePlan(), p.EpochRatio)

			logger := k.Logger(ctx)
			logger.Info("updated public ratio plan", "ratio_plan", plan)
		}

		k.SetPlan(ctx, plan)
	}

	return nil
}

// DeletePublicPlanProposal deletes public plan proposal once the governance proposal is passed.
func (k Keeper) DeletePublicPlanProposal(ctx sdk.Context, proposals []*types.DeleteRequestProposal) error {
	for _, p := range proposals {
		plan, found := k.GetPlan(ctx, p.GetPlanId())
		if !found {
			return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "plan %d is not found", p.GetPlanId())
		}

		if err := k.TerminatePlan(ctx, plan); err != nil {
			panic(err)
		}

		k.RemovePlan(ctx, plan)

		logger := k.Logger(ctx)
		logger.Info("removed public ratio plan", "plan_id", plan.GetId())
	}

	return nil
}
