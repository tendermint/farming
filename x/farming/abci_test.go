package farming_test

import (
	"github.com/tendermint/farming/x/farming"

	_ "github.com/stretchr/testify/suite"
)

func (suite *ModuleTestSuite) TestEndBlocker() {
	params := suite.keeper.GetParams(suite.ctx)
	suite.Require().Equal(uint32(1), params.NextEpochDays)

	suite.ctx = suite.ctx.WithBlockTime(mustParseRFC3339("2021-08-11T23:59:59Z"))
	farming.EndBlocker(suite.ctx, suite.keeper)

	// WIP
}
