package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming"
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
			50000000,
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

				suite.True(IntEq(sdk.NewInt(tc.remainingStaked), staking.StakedCoins.AmountOf(denom1)))
				suite.True(IntEq(sdk.NewInt(tc.remainingQueued), staking.QueuedCoins.AmountOf(denom1)))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestStakingCreationFee() {
	params := suite.keeper.GetParams(suite.ctx)
	params.StakingCreationFee = sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000))
	suite.keeper.SetParams(suite.ctx, params)

	// Test accounts have 30000000 coins by default.
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.addrs[0], denom1)
	suite.Require().True(IntEq(sdk.NewInt(30000000), balance.Amount))

	// Stake 29000000 coins and pay 1000000 coins as staking creation fee because
	// it's the first time staking.
	_, err := suite.keeper.Stake(suite.ctx, suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 29000000)))
	suite.Require().NoError(err)

	// Balance should be zero now.
	balance = suite.app.BankKeeper.GetBalance(suite.ctx, suite.addrs[0], denom1)
	suite.Require().True(balance.Amount.IsZero())

	// Taking a new account, staking 30000000 coins should fail because
	// there is no sufficient balance for staking creation fee.
	_, err = suite.keeper.Stake(suite.ctx, suite.addrs[1], sdk.NewCoins(sdk.NewInt64Coin(denom1, 30000000)))
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestUnstake() {
	for _, tc := range []struct {
		name            string
		addrIdx         int
		amt             sdk.Coins
		remainingStaked sdk.Coins
		remainingQueued sdk.Coins
		expectErr       bool
	}{
		{
			"from queued coins",
			0,
			sdk.NewCoins(sdk.NewInt64Coin(denom1, 5000)),
			sdk.NewCoins(sdk.NewInt64Coin(denom1, 500000), sdk.NewInt64Coin(denom2, 1000000)),
			sdk.NewCoins(sdk.NewInt64Coin(denom1, 495000)),
			false,
		},
		{
			"from staked coins",
			0,
			sdk.NewCoins(sdk.NewInt64Coin(denom1, 700000), sdk.NewInt64Coin(denom2, 100000)),
			sdk.NewCoins(sdk.NewInt64Coin(denom1, 300000), sdk.NewInt64Coin(denom2, 900000)),
			sdk.NewCoins(),
			false,
		},
		{
			"one coin",
			0,
			sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)),
			sdk.NewCoins(sdk.NewInt64Coin(denom2, 1000000)),
			sdk.NewCoins(),
			false,
		},
		{
			"unstake all",
			0,
			sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000), sdk.NewInt64Coin(denom2, 1000000)),
			sdk.NewCoins(),
			sdk.NewCoins(),
			false,
		},
		{
			"more than staked",
			0,
			sdk.NewCoins(sdk.NewInt64Coin(denom1, 1100000), sdk.NewInt64Coin(denom2, 1100000)),
			// We can use nil since there will be an error and we don't use these fields
			nil,
			nil,
			true,
		},
		{
			"no staking",
			1,
			sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000)),
			nil,
			nil,
			true,
		},
	} {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.Stake(suite.addrs[0], sdk.NewCoins(
				sdk.NewInt64Coin(denom1, 500000),
				sdk.NewInt64Coin(denom2, 1000000)))

			// Make queued coins be staked.
			suite.ctx = suite.ctx.WithBlockTime(time.Date(2021, 7, 23, 5, 0, 0, 0, time.UTC))
			farming.EndBlocker(suite.ctx, suite.keeper)
			suite.ctx = suite.ctx.WithBlockTime(time.Date(2021, 7, 24, 0, 5, 0, 0, time.UTC))
			farming.EndBlocker(suite.ctx, suite.keeper)

			suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 500000)))

			// At this moment, we have 500000denom1,1000000denom2 staked and
			// 500000denom1 queued.

			staking, err := suite.keeper.Unstake(suite.ctx, suite.addrs[tc.addrIdx], tc.amt)
			if tc.expectErr {
				suite.Error(err)
			} else {
				if suite.NoError(err) {
					suite.True(CoinsEq(tc.remainingStaked, staking.StakedCoins))
					suite.True(CoinsEq(tc.remainingQueued, staking.QueuedCoins))
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestProcessQueuedCoins() {
	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000)))

	staking, _ := suite.keeper.GetStakingByFarmer(suite.ctx, suite.addrs[0])

	suite.Require().True(staking.StakedCoins.IsZero())
	suite.Require().True(IntEq(sdk.NewInt(1000), staking.QueuedCoins.AmountOf(denom1)))

	suite.keeper.ProcessQueuedCoins(suite.ctx)

	staking, _ = suite.keeper.GetStakingByFarmer(suite.ctx, suite.addrs[0])

	suite.Require().True(IntEq(sdk.NewInt(1000), staking.StakedCoins.AmountOf(denom1)))
	suite.Require().True(staking.QueuedCoins.IsZero())
}

func (suite *KeeperTestSuite) TestEndBlockerProcessQueuedCoins() {
	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000)))

	ctx := suite.ctx.WithBlockTime(time.Date(2021, 7, 23, 5, 0, 0, 0, time.UTC))
	farming.EndBlocker(ctx, suite.keeper)

	staking, _ := suite.keeper.GetStakingByFarmer(ctx, suite.addrs[0])
	suite.Require().True(IntEq(sdk.NewInt(1000), staking.QueuedCoins.AmountOf(denom1)))
	suite.Require().True(staking.StakedCoins.IsZero(), "staked coins must be empty")

	suite.Stake(suite.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 500)))

	ctx = ctx.WithBlockTime(time.Date(2021, 7, 3, 23, 59, 59, 0, time.UTC))
	farming.EndBlocker(ctx, suite.keeper)

	staking, _ = suite.keeper.GetStakingByFarmer(ctx, suite.addrs[0])
	suite.Require().True(IntEq(sdk.NewInt(1500), staking.QueuedCoins.AmountOf(denom1)))
	suite.Require().True(staking.StakedCoins.IsZero(), "staked coins must be empty")

	ctx = ctx.WithBlockTime(time.Date(2021, 7, 24, 0, 0, 1, 0, time.UTC))
	farming.EndBlocker(ctx, suite.keeper)

	staking, _ = suite.keeper.GetStakingByFarmer(ctx, suite.addrs[0])
	suite.Require().True(staking.QueuedCoins.IsZero(), "queued coins must be empty")
	suite.Require().True(IntEq(sdk.NewInt(1500), staking.StakedCoins.AmountOf(denom1)))
}
