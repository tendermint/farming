package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tendermint/farming/x/farming/types"
)

// Simulation parameter constants.
const (
	PrivatePlanCreationFee = "private_plan_creation_fee"
	StakingCreationFee     = "staking_creation_fee"
	EpochDays              = "epoch_days"
)

// GenPrivatePlanCreationFee return randomized private plan creation fee.
func GenPrivatePlanCreationFee(r *rand.Rand) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simulation.RandIntBetween(r, 0, 1_000_000_000))))
}

// GenStakingCreationFee return randomized staking creation fee.
func GenStakingCreationFee(r *rand.Rand) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simulation.RandIntBetween(r, 0, 1_000_000_000))))
}

// GenEpochDays return default EpochDays.
func GenEpochDays(r *rand.Rand) uint32 {
	return types.DefaultEpochDays
}

// RandomizedGenState generates a random GenesisState for farming.
func RandomizedGenState(simState *module.SimulationState) {
	var privatePlanCreationFee sdk.Coins
	simState.AppParams.GetOrGenerate(
		simState.Cdc, PrivatePlanCreationFee, &privatePlanCreationFee, simState.Rand,
		func(r *rand.Rand) { privatePlanCreationFee = GenPrivatePlanCreationFee(r) },
	)

	var stakingCreationFee sdk.Coins
	simState.AppParams.GetOrGenerate(
		simState.Cdc, StakingCreationFee, &stakingCreationFee, simState.Rand,
		func(r *rand.Rand) { stakingCreationFee = GenStakingCreationFee(r) },
	)

	var epochDays uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, EpochDays, &epochDays, simState.Rand,
		func(r *rand.Rand) { epochDays = GenEpochDays(r) },
	)

	farmingGenesis := types.GenesisState{
		Params: types.Params{
			PrivatePlanCreationFee: privatePlanCreationFee,
			StakingCreationFee:     stakingCreationFee,
			EpochDays:              epochDays,
		},
	}

	bz, _ := json.MarshalIndent(&farmingGenesis, "", " ")
	fmt.Printf("Selected randomly generated farming parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&farmingGenesis)
}
