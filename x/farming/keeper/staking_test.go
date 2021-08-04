package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestGetNextStakingID() {
	for id := uint64(1); id <= 100; id++ {
		suite.Require().Equal(id, suite.keeper.GetNextStakingIDWithUpdate(suite.ctx))
	}
}

func (suite *KeeperTestSuite) TestGetStaking() {
	_, found := suite.keeper.GetStaking(suite.ctx, 1)
	suite.False(found, "staking should not be present")

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000)))

	_, found = suite.keeper.GetStaking(suite.ctx, 1)
	suite.True(found, "staking should be present")
}

func (suite *KeeperTestSuite) TestStake() {
	for _, tc := range []struct {
		name            string
		amt             int64
		remainingStaked int64
		remainingQueued int64
		expectErr       bool
	}{
		{
			"normal",
			1000,
			0,
			1000,
			false,
		},
		{
			"more than balance",
			10_000_000_000,
			0,
			0,
			true,
		},
	} {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			_, found := suite.keeper.GetStakingByFarmer(suite.ctx, suite.addrs[0])
			suite.Require().False(found, "staking should not be present")

			staking, err := suite.keeper.Stake(suite.ctx, suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, tc.amt)))
			if tc.expectErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
				staking2, found := suite.keeper.GetStakingByFarmer(suite.ctx, suite.addrs[0])
				suite.True(found, "staking should be present")
				suite.True(staking2.StakedCoins.IsEqual(staking.StakedCoins))
				suite.True(staking2.QueuedCoins.IsEqual(staking2.QueuedCoins))

				suite.True(intEq(sdk.NewInt(tc.remainingStaked), staking.StakedCoins.AmountOf(denom1)))
				suite.True(intEq(sdk.NewInt(tc.remainingQueued), staking.QueuedCoins.AmountOf(denom1)))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestStakingCreationFee() {
	params := suite.keeper.GetParams(suite.ctx)
	params.StakingCreationFee = sdk.NewCoins(sdk.NewInt64Coin(denom1, 1_000_000))
	suite.keeper.SetParams(suite.ctx, params)

	// Test accounts have 1,000,000,000 coins by default.
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.addrs[0], denom1)
	suite.Require().True(intEq(sdk.NewInt(1_000_000_000), balance.Amount))

	// Stake 999,000,000 coins and pay 1,000,000 coins as staking creation fee because
	// it's the first time staking.
	_, err := suite.keeper.Stake(suite.ctx, suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 999_000_000)))
	suite.Require().NoError(err)

	// Balance should be zero now.
	balance = suite.app.BankKeeper.GetBalance(suite.ctx, suite.addrs[0], denom1)
	suite.Require().True(balance.Amount.IsZero())

	// Taking a new account, staking 1_000_000_000 coins should fail because
	// there is no sufficient balance for staking creation fee.
	_, err = suite.keeper.Stake(suite.ctx, suite.addrs[1], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1_000_000_000)))
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestProcessQueuedCoins() {
	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000)))

	staking, _ := suite.keeper.GetStakingByFarmer(suite.ctx, suite.addrs[0])

	suite.Require().True(staking.StakedCoins.IsZero())
	suite.Require().True(intEq(sdk.NewInt(1000), staking.QueuedCoins.AmountOf(denom1)))

	suite.keeper.ProcessQueuedCoins(suite.ctx)

	staking, _ = suite.keeper.GetStakingByFarmer(suite.ctx, suite.addrs[0])

	suite.Require().True(intEq(sdk.NewInt(1000), staking.StakedCoins.AmountOf(denom1)))
	suite.Require().True(staking.QueuedCoins.IsZero())
}
