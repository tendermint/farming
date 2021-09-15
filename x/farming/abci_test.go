package farming_test

import (
	"time"

	"github.com/tendermint/farming/x/farming"
	"github.com/tendermint/farming/x/farming/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *ModuleTestSuite) TestEndBlockerEdgeCase1() {
	suite.SetupTest()

	nextEpochDays := uint32(7)

	params := suite.keeper.GetParams(suite.ctx)
	params.NextEpochDays = nextEpochDays
	suite.keeper.SetParams(suite.ctx, params)
	suite.keeper.SetCurrentEpochDays(suite.ctx, params.NextEpochDays)

	t := types.ParseTime("2021-08-01T00:00:00Z")
	suite.ctx = suite.ctx.WithBlockTime(t)
	farming.EndBlocker(suite.ctx, suite.keeper)

	lastEpochTime, _ := suite.keeper.GetLastEpochTime(suite.ctx)

	for i := 1; i < 200; i++ {
		t = t.Add(1 * time.Hour)
		suite.ctx = suite.ctx.WithBlockTime(t)
		farming.EndBlocker(suite.ctx, suite.keeper)

		if i == 120 { // 5 days passed
			params := suite.keeper.GetParams(suite.ctx)
			params.NextEpochDays = uint32(1)
			suite.keeper.SetParams(suite.ctx, params)
		}

		currentEpochDays := suite.keeper.GetCurrentEpochDays(suite.ctx)

		t2, _ := suite.keeper.GetLastEpochTime(suite.ctx)
		if t2.After(lastEpochTime) {
			suite.Require().GreaterOrEqual(t2.Sub(lastEpochTime).Hours(), float64(nextEpochDays*24))
			suite.Require().Equal(uint32(1), currentEpochDays)
		}
	}
}

func (suite *ModuleTestSuite) TestEndBlockerEdgeCase2() {
	suite.SetupTest()

	nextEpochDays := uint32(1)

	params := suite.keeper.GetParams(suite.ctx)
	params.NextEpochDays = nextEpochDays
	suite.keeper.SetParams(suite.ctx, params)
	suite.keeper.SetCurrentEpochDays(suite.ctx, params.NextEpochDays)

	t := types.ParseTime("2021-08-01T00:00:00Z")
	suite.ctx = suite.ctx.WithBlockTime(t)
	farming.EndBlocker(suite.ctx, suite.keeper)

	lastEpochTime, _ := suite.keeper.GetLastEpochTime(suite.ctx)

	for i := 1; i < 50; i++ {
		t = t.Add(1 * time.Hour)
		suite.ctx = suite.ctx.WithBlockTime(t)
		farming.EndBlocker(suite.ctx, suite.keeper)

		if i == 10 { // 10 hours passed
			params := suite.keeper.GetParams(suite.ctx)
			params.NextEpochDays = uint32(7)
			suite.keeper.SetParams(suite.ctx, params)
		}

		currentEpochDays := suite.keeper.GetCurrentEpochDays(suite.ctx)

		t2, _ := suite.keeper.GetLastEpochTime(suite.ctx)
		if t2.After(lastEpochTime) {
			suite.Require().GreaterOrEqual(t2.Sub(lastEpochTime).Hours(), float64(nextEpochDays*24))
			suite.Require().Equal(uint32(7), currentEpochDays)
		}
	}
}
