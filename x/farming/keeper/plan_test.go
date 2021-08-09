package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming/types"
)

func (suite *KeeperTestSuite) TestGetSetNewPlan() {
	name := ""
	farmingPoolAddr := sdk.AccAddress("farmingPoolAddr")
	terminationAddr := sdk.AccAddress("terminationAddr")

	stakingCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000)))
	coinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)

	addrs := app.AddTestAddrs(suite.app, suite.ctx, 2, sdk.NewInt(2000000))
	farmerAddr := addrs[0]

	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)
	basePlan := types.NewBasePlan(1, name, 1, farmingPoolAddr.String(), terminationAddr.String(), coinWeights, startTime, endTime)
	fixedPlan := types.NewFixedAmountPlan(basePlan, sdk.NewCoins(sdk.NewCoin("testFarmCoinDenom", sdk.NewInt(1000000))))
	suite.keeper.SetPlan(suite.ctx, fixedPlan)

	planGet, found := suite.keeper.GetPlan(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(fixedPlan.Id, planGet.GetId())
	suite.Require().Equal(fixedPlan.FarmingPoolAddress, planGet.GetFarmingPoolAddress().String())

	plans := suite.keeper.GetAllPlans(suite.ctx)
	suite.Require().Len(plans, 1)
	suite.Require().Equal(fixedPlan.Id, plans[0].GetId())
	suite.Require().Equal(fixedPlan.FarmingPoolAddress, plans[0].GetFarmingPoolAddress().String())

	_, err := suite.keeper.Stake(suite.ctx, farmerAddr, stakingCoins)
	suite.Require().NoError(err)

	stakings := suite.keeper.GetAllStakings(suite.ctx)
	stakingByFarmer, found := suite.keeper.GetStakingByFarmer(suite.ctx, farmerAddr)
	stakingsByDenom := suite.keeper.GetStakingsByStakingCoinDenom(suite.ctx, sdk.DefaultBondDenom)

	suite.Require().True(found)
	suite.Require().Equal(stakings[0], stakingByFarmer)
	suite.Require().Equal(stakings, stakingsByDenom)
}
