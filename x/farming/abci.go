package farming

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming/keeper"
	"github.com/tendermint/farming/x/farming/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	for _, plan := range k.GetPlans(ctx) {
		if !plan.GetTerminated() && ctx.BlockTime().After(plan.GetEndTime()) {
			if err := k.TerminatePlan(ctx, plan); err != nil {
				panic(err)
			}
		}
	}

	// alternative impl that stores the end date of an epoch instead of
	// calculating it each time.
	// XXX: it may have collateral effects in other parts of the code (unchecked)
	blockDate := ctx.BlockTime()
	epochEndTime, isSet := k.GetEpochEndTime(ctx) // GetEpochEndDate replaces GetLastEpochTime
	if !isSet || blockDate.After(epochEndTime) {
		nextEpochDays := k.GetNextEpochDuration(ctx) // GetNextEpochDays replaces GetCurrentEpochDays
		nextEpochEndTime := blockDate.AddDate(0, 0, int(nextEpochDays))
		// advance epoch
		if err := k.AllocateRewards(ctx); err != nil {
			panic(err)
		}
		k.ProcessQueuedCoins(ctx)
		k.SetEpochEndTime(ctx, nextEpochEndTime)
	}


	// CurrentEpochDays is initialized with the value of NextEpochDays in genesis, and
	// it is used here to prevent from affecting the epoch days for farming rewards allocation.
	// Suppose NextEpochDays is 7 days, and it is proposed to change the value to 1 day through governance proposal.
	// Although the proposal is passed, farming rewards allocation should continue to proceed with 7 days,
	// and then it gets updated.

	//lastEpochTime, found := k.GetLastEpochTime(ctx)
	//if !found {
	//	k.SetLastEpochTime(ctx, ctx.BlockTime())
	//} else {
	//	currentEpochDays := k.GetCurrentEpochDays(ctx)
	//	y, m, d := lastEpochTime.AddDate(0, 0, int(currentEpochDays)).Date()
	//	y2, m2, d2 := ctx.BlockTime().Date()
	//	//
	//	epochExpectedEnd := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	//	blockDate := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	//
	//	if !blockDate.Before(epochExpectedEnd) {
	//		if err := k.AdvanceEpoch(ctx); err != nil {
	//			panic(err)
	//		}
	//		if params := k.GetParams(ctx); params.NextEpochDays != currentEpochDays {
	//			k.SetCurrentEpochDays(ctx, params.NextEpochDays)
	//		}
	//	}
	//}
}
