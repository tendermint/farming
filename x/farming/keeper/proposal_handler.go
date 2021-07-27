package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/farming/x/farming/types"
)

// HandlePublicPlanProposal is a handler for executing a public plan creation proposal.
func HandlePublicPlanProposal(ctx sdk.Context, k Keeper, proposal *types.PublicPlanProposal) error {
	if err := proposal.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	switch {
	case proposal.AddRequestProposals != nil:
		if err := k.AddPublicPlanProposal(ctx, proposal.Name, proposal.AddRequestProposals); err != nil {
			return err
		}
	case proposal.UpdateRequestProposals != nil:
		if err := k.UpdatePublicPlanProposal(ctx, proposal.UpdateRequestProposals); err != nil {
			return err
		}
	case proposal.DeleteRequestProposals != nil:
		if err := k.DeletePublicPlanProposal(ctx, proposal.DeleteRequestProposals); err != nil {
			return err
		}
	default:
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unexpected public plan proposal %s", proposal.String())
	}

	return nil
}

// AddPublicPlanProposal adds a new public plan once the governance proposal is passed.
func (k Keeper) AddPublicPlanProposal(ctx sdk.Context, name string, proposals []*types.AddRequestProposal) error {
	for _, p := range proposals {
		farmingPoolAddrAcc, err := sdk.AccAddressFromBech32(p.GetFarmingPoolAddress())
		if err != nil {
			return err
		}

		if !p.EpochAmount.IsZero() && !p.EpochAmount.IsAnyNegative() {
			msg := types.NewMsgCreateFixedAmountPlan(
				name,
				farmingPoolAddrAcc,
				p.GetStakingCoinWeights(),
				p.GetStartTime(),
				p.GetEndTime(),
				p.EpochAmount,
			)

			plan, err := k.CreateFixedAmountPlan(ctx, msg, types.PlanTypePublic)
			if err != nil {
				return err
			}

			logger := k.Logger(ctx)
			logger.Info("created public fixed amount plan", "fixed_amount_plan", plan)
		}

		if !p.EpochRatio.IsZero() && !p.EpochRatio.IsNil() && !p.EpochRatio.IsNegative() {
			msg := types.NewMsgCreateRatioPlan(
				name,
				farmingPoolAddrAcc,
				p.GetStakingCoinWeights(),
				p.GetStartTime(),
				p.GetEndTime(),
				p.EpochRatio,
			)

			plan, err := k.CreateRatioPlan(ctx, msg, types.PlanTypePublic)
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
	for _, proposal := range proposals {
		plan, found := k.GetPlan(ctx, proposal.GetPlanId())
		if !found {
			return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "plan %d is not found", proposal.GetPlanId())
		}

		farmingPoolAddrAcc, err := sdk.AccAddressFromBech32(proposal.GetFarmingPoolAddress())
		if err != nil {
			return err
		}

		terminationAddrAcc, err := sdk.AccAddressFromBech32(proposal.GetTerminationAddress())
		if err != nil {
			return err
		}

		switch p := plan.(type) {
		case *types.FixedAmountPlan:
			if err := p.SetFarmingPoolAddress(farmingPoolAddrAcc); err != nil {
				return err
			}
			if err := p.SetTerminationAddress(terminationAddrAcc); err != nil {
				return err
			}
			if err := p.SetStakingCoinWeights(proposal.GetStakingCoinWeights()); err != nil {
				return err
			}
			if err := p.SetStartTime(proposal.GetStartTime()); err != nil {
				return err
			}
			if err := p.SetEndTime(proposal.GetStartTime()); err != nil {
				return err
			}
			p.EpochAmount = proposal.GetEpochAmount()

			k.SetPlan(ctx, p)

			logger := k.Logger(ctx)
			logger.Info("updated public fixed amount plan", "fixed_amount_plan", plan)

		case *types.RatioPlan:
			if err := p.SetFarmingPoolAddress(farmingPoolAddrAcc); err != nil {
				return err
			}
			if err := p.SetTerminationAddress(terminationAddrAcc); err != nil {
				return err
			}
			if err := p.SetStakingCoinWeights(proposal.GetStakingCoinWeights()); err != nil {
				return err
			}
			if err := p.SetStartTime(proposal.GetStartTime()); err != nil {
				return err
			}
			if err := p.SetEndTime(proposal.GetStartTime()); err != nil {
				return err
			}
			p.EpochRatio = proposal.EpochRatio

			k.SetPlan(ctx, p)

			logger := k.Logger(ctx)
			logger.Info("updated public ratio plan", "ratio_plan", plan)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized plan type: %T", p)
		}
	}

	return nil
}

// DeletePublicPlanProposal delets public plan proposal once the governance proposal is passed.
func (k Keeper) DeletePublicPlanProposal(ctx sdk.Context, proposals []*types.DeleteRequestProposal) error {
	for _, p := range proposals {
		plan, found := k.GetPlan(ctx, p.GetPlanId())
		if !found {
			return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "plan %d is not found", p.GetPlanId())
		}

		k.RemovePlan(ctx, plan)

		logger := k.Logger(ctx)
		logger.Info("removed public ratio plan", "plan_id", plan.GetId())
	}

	return nil
}
