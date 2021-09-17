package keeper_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming"
	"github.com/tendermint/farming/x/farming/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestInitGenesis() {
	plans := []types.PlanI{
		types.NewFixedAmountPlan(
			types.NewBasePlan(
				1,
				"name1",
				types.PlanTypePrivate,
				suite.addrs[0].String(),
				suite.addrs[0].String(),
				sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(denom1, sdk.NewDecWithPrec(3, 1)),
					sdk.NewDecCoinFromDec(denom2, sdk.NewDecWithPrec(7, 1))),
				types.ParseTime("2021-07-30T00:00:00Z"),
				types.ParseTime("2021-08-30T00:00:00Z"),
			),
			sdk.NewCoins(sdk.NewInt64Coin(denom3, 1_000_000))),
		types.NewRatioPlan(
			types.NewBasePlan(
				2,
				"name2",
				types.PlanTypePublic,
				suite.addrs[0].String(),
				suite.addrs[0].String(),
				sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(denom1, sdk.NewDecWithPrec(3, 1)),
					sdk.NewDecCoinFromDec(denom2, sdk.NewDecWithPrec(7, 1))),
				types.ParseTime("2021-07-30T00:00:00Z"),
				types.ParseTime("2021-08-30T00:00:00Z"),
			),
			sdk.MustNewDecFromStr("0.01")),
	}
	//for _, plan := range plans {
	//	suite.keeper.SetPlan(suite.ctx, plan)
	//}
	suite.keeper.SetPlan(suite.ctx, plans[1])
	suite.keeper.SetPlan(suite.ctx, plans[0])

	suite.Stake(suite.addrs[1], sdk.NewCoins(
		sdk.NewInt64Coin(denom1, 1_000_000),
		sdk.NewInt64Coin(denom2, 1_000_000)))
	suite.keeper.ProcessQueuedCoins(suite.ctx)

	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-07-31T00:00:00Z"))

	// Advance 2 epochs
	err := suite.keeper.AdvanceEpoch(suite.ctx)
	suite.Require().NoError(err)
	err = suite.keeper.AdvanceEpoch(suite.ctx)
	suite.Require().NoError(err)

	var genState *types.GenesisState
	suite.Require().NotPanics(func() {
		genState = suite.keeper.ExportGenesis(suite.ctx)
	})

	err = types.ValidateGenesis(*genState)
	suite.Require().NoError(err)

	suite.Require().NotPanics(func() {
		suite.keeper.InitGenesis(suite.ctx, *genState)
	})
	suite.Require().Equal(genState, suite.keeper.ExportGenesis(suite.ctx))
}

func (suite *KeeperTestSuite) TestMarshalUnmarshalDefaultGenesis() {
	genState := suite.keeper.ExportGenesis(suite.ctx)
	bz, err := suite.app.AppCodec().MarshalJSON(genState)
	suite.Require().NoError(err)
	var genState2 types.GenesisState
	err = suite.app.AppCodec().UnmarshalJSON(bz, &genState2)
	suite.Require().NoError(err)
	suite.Require().Equal(*genState, genState2)

	app2 := simapp.Setup(false)
	ctx2 := app2.BaseApp.NewContext(false, tmproto.Header{})
	keeper2 := app2.FarmingKeeper
	keeper2.InitGenesis(ctx2, genState2)

	genState3 := keeper2.ExportGenesis(ctx2)
	suite.Require().Equal(genState2, *genState3)
}

func (suite *KeeperTestSuite) TestExportGenesis() {
	for i := range suite.samplePlans {
		plan := suite.samplePlans[len(suite.samplePlans)-i-1]
		suite.keeper.SetPlan(suite.ctx, plan)
	}

	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-08-04T23:00:00Z"))
	farming.EndBlocker(suite.ctx, suite.keeper)
	suite.Stake(suite.addrs[1], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000), sdk.NewInt64Coin(denom2, 800000)))
	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 500000), sdk.NewInt64Coin(denom2, 700000)))
	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-08-05T00:00:00Z"))
	farming.EndBlocker(suite.ctx, suite.keeper) // queued coins => staked coins
	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-08-06T00:00:00Z"))
	farming.EndBlocker(suite.ctx, suite.keeper) // allocate rewards
	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 2000000), sdk.NewInt64Coin(denom2, 1200000)))
	suite.Stake(suite.addrs[1], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1500000), sdk.NewInt64Coin(denom2, 300000)))

	genState := suite.keeper.ExportGenesis(suite.ctx)
	bz, err := suite.app.AppCodec().MarshalJSON(genState)
	suite.Require().NoError(err)
	*genState = types.GenesisState{}
	err = suite.app.AppCodec().UnmarshalJSON(bz, genState)
	suite.Require().NoError(err)

	for _, tc := range []struct {
		name  string
		check func()
	}{
		{
			"Params",
			func() {
				err := genState.Params.Validate()
				suite.Require().NoError(err)
				suite.Require().Equal(suite.keeper.GetParams(suite.ctx), genState.Params)
			},
		},
		{
			"PlanRecords",
			func() {
				suite.Require().Len(genState.PlanRecords, len(suite.samplePlans))
				for _, record := range genState.PlanRecords {
					err := record.Validate()
					suite.Require().NoError(err)
					_, err = types.UnpackPlan(&record.Plan)
					suite.Require().NoError(err)
					// TODO: add more checks
				}
			},
		},
		{
			"StakingRecords",
			func() {
				suite.Require().Len(genState.StakingRecords, 4)
				for _, record := range genState.StakingRecords {
					switch record.Farmer {
					case suite.addrs[0].String():
						switch record.StakingCoinDenom {
						case denom1:
							suite.Require().True(intEq(record.Staking.Amount, sdk.NewInt(500000)))
						case denom2:
							suite.Require().True(intEq(record.Staking.Amount, sdk.NewInt(700000)))
						}
					case suite.addrs[1].String():
						switch record.StakingCoinDenom {
						case denom1:
							suite.Require().True(intEq(record.Staking.Amount, sdk.NewInt(1000000)))
						case denom2:
							suite.Require().True(intEq(record.Staking.Amount, sdk.NewInt(800000)))
						}
					}
				}
			},
		},
		{
			"QueuedStakingRecords",
			func() {
				suite.Require().Len(genState.QueuedStakingRecords, 4)
				for _, record := range genState.QueuedStakingRecords {
					switch record.Farmer {
					case suite.addrs[0].String():
						switch record.StakingCoinDenom {
						case denom1:
							suite.Require().True(intEq(record.QueuedStaking.Amount, sdk.NewInt(2000000)))
						case denom2:
							suite.Require().True(intEq(record.QueuedStaking.Amount, sdk.NewInt(1200000)))
						}
					case suite.addrs[1].String():
						switch record.StakingCoinDenom {
						case denom1:
							suite.Require().True(intEq(record.QueuedStaking.Amount, sdk.NewInt(1500000)))
						case denom2:
							suite.Require().True(intEq(record.QueuedStaking.Amount, sdk.NewInt(300000)))
						}
					}
				}
			},
		},
		{
			"HistoricalRewards",
			func() {},
		},
		{
			"OutstandingRewards",
			func() {},
		},
		{
			"CurrentEpochRecords",
			func() {},
		},
		{
			"StakingReserveCoins",
			func() {},
		},
		{
			"RewardPoolCoins",
			func() {},
		},
		{
			"LastEpochTime",
			func() {},
		},
		{
			"CurrentEpochDays",
			func() {},
		},
	} {
		suite.Run(tc.name, tc.check)
	}
}
