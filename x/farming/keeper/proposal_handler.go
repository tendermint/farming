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

	// TODO: ctx로 전체 플랜 가져온 후 동일한 farmer 주소의 epoch ratio 합이 1이 넘을경우 리턴

	return nil
}

// AddPublicPlanProposal adds a new public plan once the governance proposal is passed.
func (k Keeper) AddPublicPlanProposal(ctx sdk.Context, proposals []*types.AddRequestProposal) error {
	for _, p := range proposals {
		plans := k.GetAllPlans(ctx)
		for _, plan := range plans {
			if plan.(*types.BasePlan).Name == p.Name {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "plan name '%s' already exists", p.Name)
			}
		}

		farmingPoolAddrAcc, err := sdk.AccAddressFromBech32(p.GetFarmingPoolAddress())
		if err != nil {
			return err
		}

		if !p.EpochAmount.IsZero() && !p.EpochAmount.IsAnyNegative() {
			msg := types.NewMsgCreateFixedAmountPlan(
				p.GetName(),
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

		} else if !p.EpochRatio.IsZero() && !p.EpochRatio.IsNil() && !p.EpochRatio.IsNegative() {
			msg := types.NewMsgCreateRatioPlan(
				p.GetName(),
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
		plans := k.GetAllPlans(ctx)
		for _, plan := range plans {
			if plan.(*types.BasePlan).Name == proposal.Name {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "plan name '%s' already exists", proposal.Name)
			}
		}

		plan, found := k.GetPlan(ctx, proposal.GetPlanId())
		if !found {
			return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "plan %d is not found", proposal.GetPlanId())
		}

		switch p := plan.(type) {
		case *types.FixedAmountPlan:
			if proposal.GetFarmingPoolAddress() != "" {
				farmingPoolAddrAcc, err := sdk.AccAddressFromBech32(proposal.GetFarmingPoolAddress())
				if err != nil {
					return err
				}
				if err := p.SetFarmingPoolAddress(farmingPoolAddrAcc); err != nil {
					return err
				}
			}
			if proposal.GetTerminationAddress() != "" {
				terminationAddrAcc, err := sdk.AccAddressFromBech32(proposal.GetTerminationAddress())
				if err != nil {
					return err
				}
				if err := p.SetTerminationAddress(terminationAddrAcc); err != nil {
					return err
				}
			}
			if proposal.GetStakingCoinWeights() != nil {
				if err := p.SetStakingCoinWeights(proposal.GetStakingCoinWeights()); err != nil {
					return err
				}
			}
			if proposal.GetStartTime() != nil {
				if err := p.SetStartTime(*proposal.GetStartTime()); err != nil {
					return err
				}
			}
			if proposal.GetEndTime() != nil {
				if err := p.SetEndTime(*proposal.GetEndTime()); err != nil {
					return err
				}
			}
			if proposal.GetName() != "" {
				p.Name = proposal.GetName()
			}
			if proposal.GetEpochAmount() != nil {
				p.EpochAmount = proposal.GetEpochAmount()
			}

			k.SetPlan(ctx, p)

			logger := k.Logger(ctx)
			logger.Info("updated public fixed amount plan", "fixed_amount_plan", plan)

		case *types.RatioPlan:
			if proposal.GetFarmingPoolAddress() != "" {
				farmingPoolAddrAcc, err := sdk.AccAddressFromBech32(proposal.GetFarmingPoolAddress())
				if err != nil {
					return err
				}
				if err := p.SetFarmingPoolAddress(farmingPoolAddrAcc); err != nil {
					return err
				}
			}
			if proposal.GetTerminationAddress() != "" {
				terminationAddrAcc, err := sdk.AccAddressFromBech32(proposal.GetTerminationAddress())
				if err != nil {
					return err
				}
				if err := p.SetTerminationAddress(terminationAddrAcc); err != nil {
					return err
				}
			}
			if proposal.GetStakingCoinWeights() != nil {
				if err := p.SetStakingCoinWeights(proposal.GetStakingCoinWeights()); err != nil {
					return err
				}
			}
			if proposal.GetStartTime() != nil {
				if err := p.SetStartTime(*proposal.GetStartTime()); err != nil {
					return err
				}
			}
			if proposal.GetEndTime() != nil {
				if err := p.SetEndTime(*proposal.GetEndTime()); err != nil {
					return err
				}
			}
			if proposal.GetName() != "" {
				p.Name = proposal.GetName()
			}
			if !proposal.EpochRatio.IsNil() {
				p.EpochRatio = proposal.EpochRatio
			}

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
