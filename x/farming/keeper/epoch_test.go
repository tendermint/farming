package keeper_test

import (
	"time"
)

func (suite *KeeperTestSuite) TestLastEpochTime() {
	_, found := suite.keeper.GetLastEpochTime(suite.ctx)
	suite.Require().False(found)

	t := time.Date(2021, 7, 23, 5, 1, 2, 3, time.UTC)
	suite.keeper.SetLastEpochTime(suite.ctx, t)

	t2, found := suite.keeper.GetLastEpochTime(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(t, t2)
}
