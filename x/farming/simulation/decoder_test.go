package simulation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/tendermint/farming/x/farming/simulation"
	"github.com/tendermint/farming/x/farming/types"
)

func TestDecodeFarmingStore(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	basePlan := types.BasePlan{}
	staking := types.Staking{}
	reward := types.Reward{}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.PlanKeyPrefix, Value: cdc.MustMarshal(&basePlan)},
			{Key: types.PlansByFarmerIndexKeyPrefix, Value: cdc.MustMarshal(&basePlan)},
			{Key: types.StakingKeyPrefix, Value: cdc.MustMarshal(&staking)},
			{Key: types.StakingByFarmerIndexKeyPrefix, Value: cdc.MustMarshal(&staking)},
			{Key: types.StakingsByStakingCoinDenomIndexKeyPrefix, Value: cdc.MustMarshal(&staking)},
			{Key: types.RewardKeyPrefix, Value: cdc.MustMarshal(&reward)},
			{Key: types.RewardsByFarmerIndexKeyPrefix, Value: cdc.MustMarshal(&reward)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Plan", fmt.Sprintf("%v\n%v", basePlan, basePlan)},
		{"Plans", fmt.Sprintf("%v\n%v", basePlan, basePlan)},
		{"Staking", fmt.Sprintf("%v\n%v", staking, staking)},
		{"StakingByFarmer", fmt.Sprintf("%v\n%v", staking, staking)},
		{"Stakings", fmt.Sprintf("%v\n%v", staking, staking)},
		{"Reward", fmt.Sprintf("%v\n%v", reward, reward)},
		{"Rewards", fmt.Sprintf("%v\n%v", reward, reward)},
		{"other", ""},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
