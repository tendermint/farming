package simulation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/tendermint/farming/x/farming/simulation"
	"github.com/tendermint/farming/x/farming/types"
)

// case bytes.Equal(kvA.Key[:1], types.PlanKeyPrefix):
// 	var pA, pB types.BasePlan
// 	cdc.MustUnmarshal(kvA.Value, &pA)
// 	cdc.MustUnmarshal(kvA.Value, &pB)
// 	return fmt.Sprintf("%v\n%v", pA, pB)

// case bytes.Equal(kvA.Key[:1], types.PlansByFarmerIndexKeyPrefix),
// 	bytes.Equal(kvA.Key[:1], types.StakingByFarmerIndexKeyPrefix),
// 	bytes.Equal(kvA.Key[:1], types.RewardsByFarmerIndexKeyPrefix):
// 	return fmt.Sprintf("%v\n%v", sdk.AccAddress(kvA.Value), sdk.AccAddress(kvB.Value))

// case bytes.Equal(kvA.Key[:1], types.StakingKeyPrefix):
// 	var sA, sB types.Staking
// 	cdc.MustUnmarshal(kvA.Value, &sA)
// 	cdc.MustUnmarshal(kvA.Value, &sB)
// 	return fmt.Sprintf("%v\n%v", sA, sB)

// case bytes.Equal(kvA.Key[:1], types.StakingsByStakingCoinDenomIndexKeyPrefix):
// 	return fmt.Sprintf("%s\n%s", kvA.Value, kvB.Value)

// case bytes.Equal(kvA.Key[:1], types.RewardKeyPrefix):
// 	var rA, rB types.Reward
// 	cdc.MustUnmarshal(kvA.Value, &rA)
// 	return fmt.Sprintf("%v\n%v", rA, rB)

var (
	pk1         = ed25519.GenPrivKey().PubKey()
	farmerAddr1 = sdk.AccAddress(pk1.Address())
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
			{Key: types.StakingKeyPrefix, Value: cdc.MustMarshal(&staking)},
			{Key: types.RewardKeyPrefix, Value: cdc.MustMarshal(&reward)},
			{Key: types.PlansByFarmerIndexKeyPrefix, Value: farmerAddr1.Bytes()},
			{Key: types.StakingByFarmerIndexKeyPrefix, Value: farmerAddr1.Bytes()},
			{Key: types.RewardsByFarmerIndexKeyPrefix, Value: farmerAddr1.Bytes()},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Plan", fmt.Sprintf("%v\n%v", basePlan, basePlan)},
		{"Staking", fmt.Sprintf("%v\n%v", staking, staking)},
		{"Reward", fmt.Sprintf("%v\n%v", reward, reward)},
		{"PlansByFarmerIndex", fmt.Sprintf("%v\n%v", farmerAddr1, farmerAddr1)},
		{"StakingByFarmerIndex", fmt.Sprintf("%v\n%v", farmerAddr1, farmerAddr1)},
		{"RewardsByFarmerIndex", fmt.Sprintf("%v\n%v", farmerAddr1, farmerAddr1)},
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
