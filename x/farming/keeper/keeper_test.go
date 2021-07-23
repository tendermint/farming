package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	app    *simapp.FarmingApp
	ctx    sdk.Context
	keeper keeper.Keeper
	addrs  []sdk.AccAddress
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
	suite.addrs = simapp.AddTestAddrs(suite.app, suite.ctx, 4, sdk.NewInt(30000000))
}

// Stake is a convenient method to test Keeper.Stake.
func (suite *KeeperTestSuite) Stake(addr sdk.AccAddress, amt sdk.Coins) {
	staking, found := suite.keeper.GetStakingByFarmer(suite.ctx, addr)
	if !found {
		staking.QueuedCoins = sdk.NewCoins()
	}

	err := suite.keeper.Stake(suite.ctx, addr, amt)
	suite.Require().NoError(err)

	staking2, found := suite.keeper.GetStakingByFarmer(suite.ctx, addr)
	suite.Require().True(found, "staking should be present")

	suite.Require().True(staking2.QueuedCoins.IsEqual(staking.QueuedCoins.Add(amt...)), "inconsistent queued coins amount")
}

func IntEq(exp, got sdk.Int) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}