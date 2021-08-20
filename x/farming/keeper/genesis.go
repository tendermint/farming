package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming/types"
)

// InitGenesis initializes the farming module's state from a given genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := k.ValidateGenesis(ctx, genState); err != nil {
		panic(err)
	}

	k.SetParams(ctx, genState.Params)
	moduleAcc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)

	_, err := sdk.AccAddressFromBech32(genState.Params.FarmingFeeCollector)
	if err != nil {
		panic(err)
	}

	for _, record := range genState.PlanRecords {
		plan, err := types.UnpackPlan(&record.Plan)
		if err != nil {
			panic(err)
		}
		k.SetPlan(ctx, plan)
		k.SetGlobalPlanId(ctx, plan.GetId())
	}
	for _, staking := range genState.Stakings {
		k.SetStaking(ctx, staking)
	}
	for _, reward := range genState.Rewards {
		k.SetReward(ctx, reward.StakingCoinDenom, reward.GetFarmer(), reward.RewardCoins)
	}
}

// ExportGenesis returns the farming module's genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	var planRecords []types.PlanRecord

	plans := k.GetAllPlans(ctx)
	stakings := k.GetAllStakings(ctx)
	rewards := k.GetAllRewards(ctx)

	for _, plan := range plans {
		any, err := types.PackPlan(plan)
		if err != nil {
			panic(err)
		}
		planRecords = append(planRecords, types.PlanRecord{
			Plan:             *any,
			FarmingPoolCoins: k.bankKeeper.GetAllBalances(ctx, plan.GetFarmingPoolAddress()),
		})
	}

	epochTime, _ := k.GetLastEpochTime(ctx)
	return types.NewGenesisState(params, planRecords, stakings, rewards, k.bankKeeper.GetAllBalances(ctx, types.StakingReserveAcc), k.bankKeeper.GetAllBalances(ctx, types.RewardsReserveAcc), epochTime)
}

// ValidateGenesis validates the farming module's genesis state.
func (k Keeper) ValidateGenesis(ctx sdk.Context, genState types.GenesisState) error {
	err := types.ValidateRatioPlans(genState)
	if err != nil {
		return nil
	}

	// TODO: Make it more efficient by eliminating duplicate init
	cc, _ := ctx.CacheContext()
	k.InitGenesis(cc, genState)
	if err := k.ValidateRemainingRewardsAmount(cc); err != nil {
		return err
	}
	if err := k.ValidateStakingReservedAmount(cc); err != nil {
		return err
	}
	return nil
}
