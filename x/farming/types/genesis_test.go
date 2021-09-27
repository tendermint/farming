package types_test

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/farming/x/farming/types"
)

func TestValidateGenesis(t *testing.T) {
	addr1 := sdk.AccAddress(crypto.AddressHash([]byte("addr1")))

	testCases := []struct {
		name        string
		configure   func(*types.GenesisState)
		expectedErr string
	}{
		{
			"default case",
			func(genState *types.GenesisState) {
				params := types.DefaultParams()
				genState.Params = params
			},
			"",
		},
		{
			"invalid plan",
			func(genState *types.GenesisState) {
				plan := types.NewRatioPlan(
					types.NewBasePlan(
						1,
						"planA",
						types.PlanTypeNil,
						addr1.String(),
						addr1.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin("denom1", 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.OneDec(),
				)
				planAny, _ := types.PackPlan(plan)
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:                *planAny,
						FarmingPoolCoins:    sdk.NewCoins(),
						StakingReserveCoins: sdk.NewCoins(),
					},
				}
			},
			"unknown plan type: PLAN_TYPE_UNSPECIFIED: invalid plan type",
		},
		{
			"invalid plan record",
			func(genState *types.GenesisState) {
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:                cdctypes.Any{},
						FarmingPoolCoins:    nil,
						StakingReserveCoins: nil,
					},
				}
			},
			"empty type url: invalid type",
		},
		{
			"not sorted plan ids",
			func(genState *types.GenesisState) {
				plan1 := types.NewRatioPlan(
					types.NewBasePlan(
						1,
						"planA",
						types.PlanTypePublic,
						addr1.String(),
						addr1.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin("denom1", 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.OneDec(),
				)
				plan2 := types.NewRatioPlan(
					types.NewBasePlan(
						2,
						"planB",
						types.PlanTypePublic,
						addr1.String(),
						addr1.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin("denom1", 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.OneDec(),
				)
				planAny1, _ := types.PackPlan(plan1)
				planAny2, _ := types.PackPlan(plan2)
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:                *planAny2,
						FarmingPoolCoins:    sdk.NewCoins(),
						StakingReserveCoins: sdk.NewCoins(),
					},
					{
						Plan:                *planAny1,
						FarmingPoolCoins:    sdk.NewCoins(),
						StakingReserveCoins: sdk.NewCoins(),
					},
				}
			},
			"pool records must be sorted",
		},
		{
			"duplicate plan name",
			func(genState *types.GenesisState) {
				plan1 := types.NewRatioPlan(
					types.NewBasePlan(
						1,
						"planA",
						types.PlanTypePublic,
						addr1.String(),
						addr1.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin("denom1", 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.OneDec(),
				)
				plan2 := types.NewRatioPlan(
					types.NewBasePlan(
						2,
						"planA",
						types.PlanTypePublic,
						addr1.String(),
						addr1.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin("denom1", 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.OneDec(),
				)
				planAny1, _ := types.PackPlan(plan1)
				planAny2, _ := types.PackPlan(plan2)
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:                *planAny1,
						FarmingPoolCoins:    sdk.NewCoins(),
						StakingReserveCoins: sdk.NewCoins(),
					},
					{
						Plan:                *planAny2,
						FarmingPoolCoins:    sdk.NewCoins(),
						StakingReserveCoins: sdk.NewCoins(),
					},
				}
			},
			"planA: duplicate plan name",
		},
		{
			"invalid accumulated epoch ratio",
			func(genState *types.GenesisState) {
				plan1 := types.NewRatioPlan(
					types.NewBasePlan(
						1,
						"planA",
						types.PlanTypePublic,
						addr1.String(),
						addr1.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin("denom1", 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.OneDec(),
				)
				plan2 := types.NewRatioPlan(
					types.NewBasePlan(
						2,
						"planB",
						types.PlanTypePublic,
						addr1.String(),
						addr1.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin("denom1", 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.OneDec(),
				)
				planAny1, _ := types.PackPlan(plan1)
				planAny2, _ := types.PackPlan(plan2)
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:                *planAny1,
						FarmingPoolCoins:    sdk.NewCoins(),
						StakingReserveCoins: sdk.NewCoins(),
					},
					{
						Plan:                *planAny2,
						FarmingPoolCoins:    sdk.NewCoins(),
						StakingReserveCoins: sdk.NewCoins(),
					},
				}
			},
			"total epoch ratio must be lower than 1: invalid request",
		},
		{
			"invalid NextEpochDays case",
			func(genState *types.GenesisState) {
				params := types.DefaultParams()
				params.NextEpochDays = 0
				genState.Params = params
			},
			"next epoch days must be positive: 0",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			genState := types.DefaultGenesisState()
			tc.configure(genState)

			err := types.ValidateGenesis(*genState)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
