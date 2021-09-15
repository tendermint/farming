package farming_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming"
	"github.com/tendermint/farming/x/farming/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *ModuleTestSuite) TestEndBlocker() {
	// set NextEpochDays and GlobalCurrentEpochDays to 7 days
	params := suite.keeper.GetParams(suite.ctx)
	params.NextEpochDays = 7
	suite.keeper.SetParams(suite.ctx, params)
	suite.keeper.SetGlobalCurrentEpochDays(suite.ctx, params.NextEpochDays)

	// set fixed amount plan
	suite.keeper.SetPlan(suite.ctx, suite.sampleFixedAmtPlans[0])

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom2, 10_000_000)))
	suite.keeper.ProcessQueuedCoins(suite.ctx)

	currEpochDays := suite.keeper.GetGlobalCurrentEpochDays(suite.ctx)
	fmt.Println("currEpochDays: ", currEpochDays)

	balancesBefore := suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.addrs[0])
	fmt.Println("balancesBefore: ", balancesBefore)

	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-08-05T00:00:00Z"))
	farming.EndBlocker(suite.ctx, suite.keeper)

	balancesAfter := suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.addrs[0])
	fmt.Println("balancesAfter: ", balancesAfter)

	farming.EndBlocker(suite.ctx, suite.keeper)

	// suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-08-11T23:59:59Z"))
	// farming.EndBlocker(suite.ctx, suite.keeper)

	//
	// params.NextEpochDays = 1
	// suite.keeper.SetParams(suite.ctx, params)
}
