package keeper_test

import (
	"fmt"
	"testing"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming"
	"github.com/tendermint/farming/x/farming/keeper"
	"github.com/tendermint/farming/x/farming/types"
)

const (
	denom1 = "denom1"
	denom2 = "denom2"
	denom3 = "denom3"
)

var (
	initialBalances = sdk.NewCoins(
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 1_000_000_000),
		sdk.NewInt64Coin(denom1, 1_000_000_000),
		sdk.NewInt64Coin(denom2, 1_000_000_000),
		sdk.NewInt64Coin(denom3, 1_000_000_000))
)

type KeeperTestSuite struct {
	suite.Suite

	app                 *simapp.FarmingApp
	ctx                 sdk.Context
	keeper              keeper.Keeper
	querier             keeper.Querier
	govHandler          govtypes.Handler
	addrs               []sdk.AccAddress
	sampleFixedAmtPlans []types.PlanI
	sampleRatioPlans    []types.PlanI
	samplePlans         []types.PlanI
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app = app
	suite.ctx = ctx
	suite.keeper = suite.app.FarmingKeeper
	suite.querier = keeper.Querier{Keeper: suite.keeper}
	suite.govHandler = farming.NewPublicPlanProposalHandler(suite.keeper)
	suite.addrs = simapp.AddTestAddrs(suite.app, suite.ctx, 6, sdk.ZeroInt())
	for _, addr := range suite.addrs {
		err := simapp.FundAccount(suite.app.BankKeeper, suite.ctx, addr, initialBalances)
		suite.Require().NoError(err)
	}
	suite.sampleFixedAmtPlans = []types.PlanI{
		types.NewFixedAmountPlan(
			types.NewBasePlan(
				1,
				"testPlan1",
				types.PlanTypePrivate,
				suite.addrs[4].String(),
				suite.addrs[4].String(),
				sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(denom1, sdk.NewDecWithPrec(3, 1)), // 30%
					sdk.NewDecCoinFromDec(denom2, sdk.NewDecWithPrec(7, 1)), // 70%
				),
				types.ParseTime("2021-08-02T00:00:00Z"),
				types.ParseTime("2021-08-10T00:00:00Z"),
			),
			sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000000)),
		),
		types.NewFixedAmountPlan(
			types.NewBasePlan(
				2,
				"testPlan2",
				types.PlanTypePublic,
				suite.addrs[5].String(),
				suite.addrs[5].String(),
				sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(denom1, sdk.OneDec()), // 100%
				),
				types.ParseTime("2021-08-04T00:00:00Z"),
				types.ParseTime("2021-08-12T00:00:00Z"),
			),
			sdk.NewCoins(sdk.NewInt64Coin(denom3, 2000000)),
		),
	}
	suite.sampleRatioPlans = []types.PlanI{
		types.NewRatioPlan(
			types.NewBasePlan(
				3,
				"testPlan3",
				types.PlanTypePrivate,
				suite.addrs[4].String(),
				suite.addrs[4].String(),
				sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(denom1, sdk.NewDecWithPrec(5, 1)), // 50%
					sdk.NewDecCoinFromDec(denom2, sdk.NewDecWithPrec(5, 1)), // 50%
				),
				types.ParseTime("2021-08-01T00:00:00Z"),
				types.ParseTime("2021-08-09T00:00:00Z"),
			),
			sdk.NewDecWithPrec(4, 2), // 4%
		),
		types.NewRatioPlan(
			types.NewBasePlan(
				4,
				"testPlan4",
				types.PlanTypePublic,
				suite.addrs[5].String(),
				suite.addrs[5].String(),
				sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(denom2, sdk.OneDec()), // 100%
				),
				types.ParseTime("2021-08-03T00:00:00Z"),
				types.ParseTime("2021-08-07T00:00:00Z"),
			),
			sdk.NewDecWithPrec(3, 2), // 3%
		),
	}
	suite.samplePlans = append(suite.sampleFixedAmtPlans, suite.sampleRatioPlans...)
}

func (suite *KeeperTestSuite) AddTestAddrs(num int, coins sdk.Coins) []sdk.AccAddress {
	addrs := simapp.AddTestAddrs(suite.app, suite.ctx, num, sdk.ZeroInt())
	for _, addr := range addrs {
		err := simapp.FundAccount(suite.app.BankKeeper, suite.ctx, addr, coins)
		suite.Require().NoError(err)
	}
	return addrs
}

// Stake is a convenient method to test Keeper.Stake.
func (suite *KeeperTestSuite) Stake(farmerAcc sdk.AccAddress, amt sdk.Coins) {
	err := suite.keeper.Stake(suite.ctx, farmerAcc, amt)
	suite.Require().NoError(err)
}

// Unstake is a convenient method to test Keeper.Unstake.
func (suite *KeeperTestSuite) Unstake(farmerAcc sdk.AccAddress, amt sdk.Coins) {
	err := suite.keeper.Unstake(suite.ctx, farmerAcc, amt)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) Harvest(farmerAcc sdk.AccAddress, stakingCoinDenoms []string) {
	err := suite.keeper.Harvest(suite.ctx, farmerAcc, stakingCoinDenoms)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) AllRewards(farmerAcc sdk.AccAddress) sdk.Coins {
	return suite.keeper.AllRewards(suite.ctx, farmerAcc)
}

func (suite *KeeperTestSuite) AdvanceEpoch() {
	err := suite.keeper.AdvanceEpoch(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) CreateFixedAmountPlan(farmingPoolAcc sdk.AccAddress, stakingCoinWeightsStr, epochAmountStr string) {
	stakingCoinWeights, err := sdk.ParseDecCoins(stakingCoinWeightsStr)
	if err != nil {
		panic(err)
	}

	epochAmount, err := sdk.ParseCoinsNormalized(epochAmountStr)
	if err != nil {
		panic(err)
	}

	msg := types.NewMsgCreateFixedAmountPlan(
		fmt.Sprintf("plan%d", suite.keeper.GetGlobalPlanId(suite.ctx)+1),
		farmingPoolAcc,
		stakingCoinWeights,
		types.ParseTime("0001-01-01T00:00:00Z"),
		types.ParseTime("9999-12-31T00:00:00Z"),
		epochAmount,
	)
	_, err = suite.keeper.CreateFixedAmountPlan(suite.ctx, msg, farmingPoolAcc, farmingPoolAcc, types.PlanTypePublic)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) CreateRatioPlan(farmingPoolAcc sdk.AccAddress, stakingCoinWeightsStr, epochRatioStr string) {
	stakingCoinWeights, err := sdk.ParseDecCoins(stakingCoinWeightsStr)
	if err != nil {
		panic(err)
	}

	epochRatio := sdk.MustNewDecFromStr(epochRatioStr)

	msg := types.NewMsgCreateRatioPlan(
		fmt.Sprintf("plan%d", suite.keeper.GetGlobalPlanId(suite.ctx)+1),
		farmingPoolAcc,
		stakingCoinWeights,
		types.ParseTime("0001-01-01T00:00:00Z"),
		types.ParseTime("9999-12-31T00:00:00Z"),
		epochRatio,
	)
	_, err = suite.keeper.CreateRatioPlan(suite.ctx, msg, farmingPoolAcc, farmingPoolAcc, types.PlanTypePublic)
	suite.Require().NoError(err)
}

func intEq(exp, got sdk.Int) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func decEq(exp, got sdk.Dec) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func coinsEq(exp, got sdk.Coins) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func decCoinsEq(exp, got sdk.DecCoins) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}
