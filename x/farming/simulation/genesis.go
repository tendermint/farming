package simulation

// DONTCOVER

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tendermint/farming/x/farming/types"
)

// Simulation parameter constants.
const (
	PrivatePlanCreationFee  = "private_plan_creation_fee"
	NextEpochDays           = "next_epoch_days"
	FarmingFeeCollector     = "farming_fee_collector"
	CurrentEpochDays        = "current_epoch_days"
	MaxNumPrivatePlans      = "max_num_private_plans"
	PrivatePlanMaxNumDenoms = "private_plan_max_num_denoms"
	PublicPlanMaxNumDenoms  = "public_plan_max_num_denoms"
)

// GenPrivatePlanCreationFee return randomized private plan creation fee.
func GenPrivatePlanCreationFee(r *rand.Rand) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simulation.RandIntBetween(r, 0, 100_000_000))))
}

// GenNextEpochDays return default next epoch days.
func GenNextEpochDays(r *rand.Rand) uint32 {
	return uint32(simulation.RandIntBetween(r, int(types.DefaultNextEpochDays), 10))
}

// GenCurrentEpochDays returns current epoch days.
func GenCurrentEpochDays(r *rand.Rand) uint32 {
	return uint32(simulation.RandIntBetween(r, int(types.DefaultCurrentEpochDays), 10))
}

// GenFarmingFeeCollector returns default farming fee collector.
func GenFarmingFeeCollector(r *rand.Rand) string {
	return types.DefaultFarmingFeeCollector.String()
}

// GenMaxNumPrivatePlans returns a randomized value for MaxNumPrivatePlans param.
func GenMaxNumPrivatePlans(r *rand.Rand) uint32 {
	return uint32(simulation.RandIntBetween(r, 1, 10000))
}

// GenPrivatePlanMaxNumDenoms returns a randomized value for PrivatePlanMaxNumDenoms param.
func GenPrivatePlanMaxNumDenoms(r *rand.Rand) uint32 {
	return uint32(simulation.RandIntBetween(r, 1, 100))
}

// GenPublicPlanMaxNumDenoms returns a randomized value for PublicPlanMaxNumDenoms param.
func GenPublicPlanMaxNumDenoms(r *rand.Rand) uint32 {
	return uint32(simulation.RandIntBetween(r, 1, 1000))
}

// RandomizedGenState generates a random GenesisState for farming.
func RandomizedGenState(simState *module.SimulationState) {
	var privatePlanCreationFee sdk.Coins
	simState.AppParams.GetOrGenerate(
		simState.Cdc, PrivatePlanCreationFee, &privatePlanCreationFee, simState.Rand,
		func(r *rand.Rand) { privatePlanCreationFee = GenPrivatePlanCreationFee(r) },
	)

	var nextEpochDays uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, NextEpochDays, &nextEpochDays, simState.Rand,
		func(r *rand.Rand) { nextEpochDays = GenNextEpochDays(r) },
	)

	var feeCollector string
	simState.AppParams.GetOrGenerate(
		simState.Cdc, FarmingFeeCollector, &feeCollector, simState.Rand,
		func(r *rand.Rand) { feeCollector = GenFarmingFeeCollector(r) },
	)

	var currentEpochDays uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, CurrentEpochDays, &currentEpochDays, simState.Rand,
		func(r *rand.Rand) { currentEpochDays = GenCurrentEpochDays(r) },
	)

	var maxNumPrivatePlans uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxNumPrivatePlans, &maxNumPrivatePlans, simState.Rand,
		func(r *rand.Rand) { maxNumPrivatePlans = GenMaxNumPrivatePlans(r) },
	)

	var privatePlanMaxNumDenoms uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, PrivatePlanMaxNumDenoms, &privatePlanMaxNumDenoms, simState.Rand,
		func(r *rand.Rand) { privatePlanMaxNumDenoms = GenPrivatePlanMaxNumDenoms(r) },
	)

	var publicPlanMaxNumDenoms uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, PublicPlanMaxNumDenoms, &publicPlanMaxNumDenoms, simState.Rand,
		func(r *rand.Rand) { publicPlanMaxNumDenoms = GenPublicPlanMaxNumDenoms(r) },
	)

	farmingGenesis := types.GenesisState{
		Params: types.Params{
			PrivatePlanCreationFee:  privatePlanCreationFee,
			NextEpochDays:           nextEpochDays,
			FarmingFeeCollector:     feeCollector,
			MaxNumPrivatePlans:      maxNumPrivatePlans,
			PrivatePlanMaxNumDenoms: privatePlanMaxNumDenoms,
			PublicPlanMaxNumDenoms:  publicPlanMaxNumDenoms,
		},
		CurrentEpochDays: currentEpochDays,
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&farmingGenesis)
}
