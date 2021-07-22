package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	suite.Require().True(balance.Amount.Equal(sdk.NewInt(30000000)))

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
	fmt.Println(err)
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
			err := suite.keeper.Stake(suite.ctx, suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000)))
			suite.Require().NoError(err)
			err = suite.keeper.Unstake(suite.ctx, suite.addrs[tc.addrIdx], tc.amt)
			if tc.expectErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}
