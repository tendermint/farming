package keeper

import (
	"strconv"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/farming/x/farming/types"
)

// GetPlan returns a plan for a given plan id.
func (k Keeper) GetPlan(ctx sdk.Context, id uint64) (plan types.PlanI, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPlanKey(id))
	if bz == nil {
		return plan, false
	}

	return k.decodePlan(bz), true
}

// GetPlans returns all plans in the store.
func (k Keeper) GetPlans(ctx sdk.Context) (plans []types.PlanI) {
	k.IteratePlans(ctx, func(plan types.PlanI) (stop bool) {
		plans = append(plans, plan)
		return false
	})

	return plans
}

// SetPlan sets a plan for a given plan id.
func (k Keeper) SetPlan(ctx sdk.Context, plan types.PlanI) {
	id := plan.GetId()
	store := ctx.KVStore(k.storeKey)

	bz, err := k.MarshalPlan(plan)
	if err != nil {
		panic(err)
	}

	store.Set(types.GetPlanKey(id), bz)
}

// DeletePlan deletes a plan from the store.
// NOTE: this will cause supply invariant violation if called
func (k Keeper) DeletePlan(ctx sdk.Context, plan types.PlanI) {
	id := plan.GetId()
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPlanKey(id))
}

// IteratePlans iterates over all the stored plans and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IteratePlans(ctx sdk.Context, cb func(plan types.PlanI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PlanKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		plan := k.decodePlan(iterator.Value())

		if cb(plan) {
			break
		}
	}
}

// GetNextPlanIdWithUpdate returns and increments the global Plan ID counter.
// If the global plan number is not set, it initializes it with value 0.
func (k Keeper) GetNextPlanIdWithUpdate(ctx sdk.Context) uint64 {
	id := k.GetGlobalPlanId(ctx) + 1
	k.SetGlobalPlanId(ctx, id)
	return id
}

// SetGlobalPlanId sets the global Plan ID counter.
func (k Keeper) SetGlobalPlanId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.GlobalPlanIdKey, bz)
}

// GetGlobalPlanId returns the global Plan ID counter.
func (k Keeper) GetGlobalPlanId(ctx sdk.Context) uint64 {
	var id uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GlobalPlanIdKey)
	if bz == nil {
		// initialize the PlanId
		id = 0
	} else {
		val := gogotypes.UInt64Value{}
		k.cdc.MustUnmarshal(bz, &val)
		id = val.GetValue()
	}
	return id
}

func (k Keeper) decodePlan(bz []byte) types.PlanI {
	acc, err := k.UnmarshalPlan(bz)
	if err != nil {
		panic(err)
	}

	return acc
}

// GetNumPrivatePlans returns the current number of private plans.
func (k Keeper) GetNumPrivatePlans(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NumPrivatePlansKey)
	var val gogotypes.UInt32Value
	k.cdc.MustUnmarshal(bz, &val)
	return val.GetValue()
}

// SetNumPrivatePlans sets the current number of private plans.
func (k Keeper) SetNumPrivatePlans(ctx sdk.Context, num uint32) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt32Value{Value: num})
	store.Set(types.NumPrivatePlansKey, bz)
}

// IncrementNumPrivatePlans increments the current number of private plans.
func (k Keeper) IncrementNumPrivatePlans(ctx sdk.Context, amt uint32) {
	num := k.GetNumPrivatePlans(ctx)
	num += amt
	k.SetNumPrivatePlans(ctx, num)
}

// DecrementNumPrivatePlans decrements the current number of private plans.
func (k Keeper) DecrementNumPrivatePlans(ctx sdk.Context, amt uint32) {
	num := k.GetNumPrivatePlans(ctx)
	if num < amt { // sanity check
		panic("cannot set negative NumPrivatePlans")
	}
	num -= amt
	k.SetNumPrivatePlans(ctx, num)
}

// MarshalPlan serializes a plan.
func (k Keeper) MarshalPlan(plan types.PlanI) ([]byte, error) { // nolint:interfacer
	return k.cdc.MarshalInterface(plan)
}

// UnmarshalPlan returns a plan from raw serialized
// bytes of a Proto-based Plan type.
func (k Keeper) UnmarshalPlan(bz []byte) (plan types.PlanI, err error) {
	return plan, k.cdc.UnmarshalInterface(bz, &plan)
}

// CreateFixedAmountPlan sets fixed amount plan.
func (k Keeper) CreateFixedAmountPlan(ctx sdk.Context, msg *types.MsgCreateFixedAmountPlan, farmingPoolAcc, terminationAcc sdk.AccAddress, typ types.PlanType) (types.PlanI, error) {
	if typ == types.PlanTypePrivate {
		params := k.GetParams(ctx)

		if k.GetNumPrivatePlans(ctx) >= params.MaxNumPrivatePlans {
			return nil, types.ErrNumPrivatePlansLimit
		}

		feeCollectorAcc, _ := sdk.AccAddressFromBech32(params.FarmingFeeCollector) // Already validated
		if err := k.bankKeeper.SendCoins(ctx, msg.GetCreator(), feeCollectorAcc, params.PrivatePlanCreationFee); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to pay private plan creation fee")
		}

		k.IncrementNumPrivatePlans(ctx, 1)
	}

	nextId := k.GetNextPlanIdWithUpdate(ctx)
	basePlan := types.NewBasePlan(
		nextId,
		msg.Name,
		typ,
		farmingPoolAcc.String(),
		terminationAcc.String(),
		msg.StakingCoinWeights,
		msg.StartTime,
		msg.EndTime,
	)

	fixedPlan := types.NewFixedAmountPlan(basePlan, msg.EpochAmount)

	k.SetPlan(ctx, fixedPlan)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateFixedAmountPlan,
			sdk.NewAttribute(types.AttributeKeyPlanId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyPlanName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyFarmingPoolAddress, farmingPoolAcc.String()),
			sdk.NewAttribute(types.AttributeKeyStartTime, msg.StartTime.String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyEpochAmount, msg.EpochAmount.String()),
		),
	})

	return fixedPlan, nil
}

// CreateRatioPlan sets ratio plan.
func (k Keeper) CreateRatioPlan(ctx sdk.Context, msg *types.MsgCreateRatioPlan, farmingPoolAcc, terminationAcc sdk.AccAddress, typ types.PlanType) (types.PlanI, error) {
	nextId := k.GetNextPlanIdWithUpdate(ctx)
	if typ == types.PlanTypePrivate {
		params := k.GetParams(ctx)

		if k.GetNumPrivatePlans(ctx) >= params.MaxNumPrivatePlans {
			return nil, types.ErrNumPrivatePlansLimit
		}

		feeCollectorAcc, _ := sdk.AccAddressFromBech32(params.FarmingFeeCollector) // Already validated
		if err := k.bankKeeper.SendCoins(ctx, msg.GetCreator(), feeCollectorAcc, params.PrivatePlanCreationFee); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to pay private plan creation fee")
		}

		k.IncrementNumPrivatePlans(ctx, 1)
	}

	basePlan := types.NewBasePlan(
		nextId,
		msg.Name,
		typ,
		farmingPoolAcc.String(),
		terminationAcc.String(),
		msg.StakingCoinWeights,
		msg.StartTime,
		msg.EndTime,
	)

	ratioPlan := types.NewRatioPlan(basePlan, msg.EpochRatio)

	k.SetPlan(ctx, ratioPlan)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateRatioPlan,
			sdk.NewAttribute(types.AttributeKeyPlanId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyPlanName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyFarmingPoolAddress, farmingPoolAcc.String()),
			sdk.NewAttribute(types.AttributeKeyStartTime, msg.StartTime.String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyEpochRatio, msg.EpochRatio.String()),
		),
	})

	return ratioPlan, nil
}

// TerminatePlan marks the plan as terminated.
// It moves the plan under different store key, which is for terminated plans.
func (k Keeper) TerminatePlan(ctx sdk.Context, plan types.PlanI) error {
	_ = plan.SetTerminated(true)
	k.SetPlan(ctx, plan)
	if plan.GetType() == types.PlanTypePrivate {
		k.DecrementNumPrivatePlans(ctx, 1)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlanTerminated,
			sdk.NewAttribute(types.AttributeKeyPlanId, strconv.FormatUint(plan.GetId(), 10)),
		),
	})

	return nil
}

// RemovePlan removes a terminated plan and sends all remaining coins in the
// farming pool address to the termination address.
func (k Keeper) RemovePlan(ctx sdk.Context, creator sdk.AccAddress, planId uint64) error {
	plan, found := k.GetPlan(ctx, planId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "plan %d not found", planId)
	}

	if !plan.IsTerminated() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "plan %d is not terminated yet", planId)
	}

	if plan.GetType() != types.PlanTypePrivate {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "plan %d is not a private plan", planId)
	}

	if !plan.GetTerminationAddress().Equals(creator) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only the plan creator can remove the plan")
	}

	if plan.GetFarmingPoolAddress().String() != plan.GetTerminationAddress().String() {
		balances := k.bankKeeper.GetAllBalances(ctx, plan.GetFarmingPoolAddress())
		if !balances.IsZero() {
			if err := k.bankKeeper.SendCoins(ctx, plan.GetFarmingPoolAddress(), plan.GetTerminationAddress(), balances); err != nil {
				return err
			}
		}
	}

	// Refund private plan creation fee.
	params := k.GetParams(ctx)
	feeCollectorAcc, _ := sdk.AccAddressFromBech32(params.FarmingFeeCollector) // Already validated
	if err := k.bankKeeper.SendCoins(ctx, feeCollectorAcc, creator, params.PrivatePlanCreationFee); err != nil {
		return sdkerrors.Wrap(err, "failed to refund private plan creation fee")
	}

	k.DeletePlan(ctx, plan)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlanRemoved,
			sdk.NewAttribute(types.AttributeKeyPlanId, strconv.FormatUint(plan.GetId(), 10)),
			sdk.NewAttribute(types.AttributeKeyFarmingPoolAddress, plan.GetFarmingPoolAddress().String()),
			sdk.NewAttribute(types.AttributeKeyTerminationAddress, plan.GetTerminationAddress().String()),
		),
	})

	return nil
}

// DerivePrivatePlanFarmingPoolAcc returns a unique account address
// of a farming pool for a private plan.
func (k Keeper) DerivePrivatePlanFarmingPoolAcc(ctx sdk.Context, name string) (sdk.AccAddress, error) {
	nextPlanId := k.GetGlobalPlanId(ctx) + 1
	poolAcc := types.PrivatePlanFarmingPoolAcc(name, nextPlanId)
	if !k.bankKeeper.GetAllBalances(ctx, poolAcc).Empty() {
		return nil, types.ErrConflictPrivatePlanFarmingPool
	}
	return poolAcc, nil
}
