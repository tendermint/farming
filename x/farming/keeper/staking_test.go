package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming"
)

func (suite *KeeperTestSuite) TestStake() {
	for _, tc := range []struct {
		name      string
		amt       sdk.Coins
		expectErr bool
	}{
		{
			"normal",
			sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)),
			false,
		},
		{
			"more than balance",
			sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 50000000)),
			true,
		},
	} {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			err := suite.keeper.Stake(suite.ctx, suite.addrs[0], tc.amt)
			if tc.expectErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
				_, found := suite.keeper.GetStakingIDByFarmer(suite.ctx, suite.addrs[0])
				suite.True(found, "staking should be present")
			}
		})
	}
}

func (suite *KeeperTestSuite) TestStakingCreationFee() {
	params := suite.keeper.GetParams(suite.ctx)
	params.StakingCreationFee = sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000))
	suite.keeper.SetParams(suite.ctx, params)

	// Test accounts have 30000000 coins by default.
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.addrs[0], sdk.DefaultBondDenom)
	suite.Require().True(IntEq(sdk.NewInt(30000000), balance.Amount))

	// Stake 29000000 coins and pay 1000000 coins as staking creation fee because
	// it's the first time staking.
	err := suite.keeper.Stake(suite.ctx, suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 29000000)))
	suite.Require().NoError(err)

	// Balance should be zero now.
	balance = suite.app.BankKeeper.GetBalance(suite.ctx, suite.addrs[0], sdk.DefaultBondDenom)
	suite.Require().True(balance.Amount.IsZero())

	// Taking a new account, staking 30000000 coins should fail because
	// there is no sufficient balance for staking creation fee.
	err = suite.keeper.Stake(suite.ctx, suite.addrs[1], sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 30000000)))
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestUnstake() {
	for _, tc := range []struct {
		name      string
		addrIdx   int
		amt       sdk.Coins
		expectErr bool
	}{
		{
			"normal",
			0,
			sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000)),
			false,
		},
		{
			"more than staked",
			0,
			sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1100000)),
			true,
		},
		{
			"no staking",
			1,
			sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)),
			true,
		},
	} {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000)))
			err := suite.keeper.Unstake(suite.ctx, suite.addrs[tc.addrIdx], tc.amt)
			if tc.expectErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestProcessQueuedCoins() {
	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)))

	staking, _ := suite.keeper.GetStakingByFarmer(suite.ctx, suite.addrs[0])

	suite.Require().True(staking.StakedCoins.IsZero())
	suite.Require().True(IntEq(sdk.NewInt(1000), staking.QueuedCoins.AmountOf(sdk.DefaultBondDenom)))

	suite.keeper.ProcessQueuedCoins(suite.ctx)

	staking, _ = suite.keeper.GetStakingByFarmer(suite.ctx, suite.addrs[0])

	suite.Require().True(IntEq(sdk.NewInt(1000), staking.StakedCoins.AmountOf(sdk.DefaultBondDenom)))
	suite.Require().True(staking.QueuedCoins.IsZero())
}

func (suite *KeeperTestSuite) TestEndBlockerProcessQueuedCoins() {
	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)))

	ctx := suite.ctx.WithBlockTime(time.Date(2021, 7, 23, 5, 0, 0, 0, time.UTC))
	farming.EndBlocker(ctx, suite.keeper)

	staking, _ := suite.keeper.GetStakingByFarmer(ctx, suite.addrs[0])
	suite.Require().True(IntEq(sdk.NewInt(1000), staking.QueuedCoins.AmountOf(sdk.DefaultBondDenom)))
	suite.Require().True(staking.StakedCoins.IsZero(), "staked coins must be empty")

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 500)))

	ctx = ctx.WithBlockTime(time.Date(2021, 7, 3, 23, 59, 59, 0, time.UTC))
	farming.EndBlocker(ctx, suite.keeper)

	staking, _ = suite.keeper.GetStakingByFarmer(ctx, suite.addrs[0])
	suite.Require().True(IntEq(sdk.NewInt(1500), staking.QueuedCoins.AmountOf(sdk.DefaultBondDenom)))
	suite.Require().True(staking.StakedCoins.IsZero(), "staked coins must be empty")

	ctx = ctx.WithBlockTime(time.Date(2021, 7, 24, 0, 0, 1, 0, time.UTC))
	farming.EndBlocker(ctx, suite.keeper)

	staking, _ = suite.keeper.GetStakingByFarmer(ctx, suite.addrs[0])
	suite.Require().True(staking.QueuedCoins.IsZero(), "queued coins must be empty")
	suite.Require().True(IntEq(sdk.NewInt(1500), staking.StakedCoins.AmountOf(sdk.DefaultBondDenom)))
}
