package farming_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	simapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming"
	"github.com/tendermint/farming/x/farming/keeper"
	"github.com/tendermint/farming/x/farming/types"
)

const (
	denom1 = "denom1"
	denom2 = "denom2"
	denom3 = "denom3"
)

var (
	initialBalances = sdk.NewCoins(
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 1_000_000_000),
		sdk.NewInt64Coin(denom1, 1_000_000_000),
		sdk.NewInt64Coin(denom2, 1_000_000_000),
		sdk.NewInt64Coin(denom3, 1_000_000_000))
)

// createTestInput returns a simapp with custom FarmingKeeper
// to avoid messing with the hooks.
func createTestInput() (*simapp.FarmingApp, sdk.Context, []sdk.AccAddress) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.FarmingKeeper = keeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.GetSubspace(types.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		map[string]bool{},
	)

	addrs := simapp.AddTestAddrs(app, ctx, 6, sdk.ZeroInt())
	for _, addr := range addrs {
		if err := simapp.FundAccount(app.BankKeeper, ctx, addr, initialBalances); err != nil {
			panic(err)
		}
	}

	return app, ctx, addrs
}

func TestMsgCreateFixedAmountPlan(t *testing.T) {
	app, ctx, addrs := createTestInput()

	msg := types.NewMsgCreateFixedAmountPlan(
		"handler-test",
		addrs[0],
		sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(denom1, sdk.NewDecWithPrec(3, 1)), // 30%
			sdk.NewDecCoinFromDec(denom2, sdk.NewDecWithPrec(7, 1)), // 70%
		),
		mustParseRFC3339("2021-08-02T00:00:00Z"),
		mustParseRFC3339("2021-08-10T00:00:00Z"),
		sdk.NewCoins(sdk.NewInt64Coin(denom3, 10_000_000)),
	)

	handler := farming.NewHandler(app.FarmingKeeper)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	plan, found := app.FarmingKeeper.GetPlan(ctx, 1)
	require.Equal(t, true, found)

	require.Equal(t, msg.Name, plan.GetName())
	require.Equal(t, msg.Creator, plan.GetTerminationAddress().String())
	require.Equal(t, msg.StakingCoinWeights, plan.GetStakingCoinWeights())
	require.Equal(t, types.PrivatePlanFarmingPoolAddress(msg.Name, 1), plan.GetFarmingPoolAddress())
	require.Equal(t, mustParseRFC3339("2021-08-02T00:00:00Z"), plan.GetStartTime())
	require.Equal(t, mustParseRFC3339("2021-08-10T00:00:00Z"), plan.GetEndTime())
	require.Equal(t, msg.EpochAmount, plan.(*types.FixedAmountPlan).EpochAmount)
}

func TestMsgCreateRatioPlan(t *testing.T) {
	app, ctx, addrs := createTestInput()

	msg := types.NewMsgCreateRatioPlan(
		"handler-test",
		addrs[0],
		sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(denom1, sdk.NewDecWithPrec(3, 1)), // 30%
			sdk.NewDecCoinFromDec(denom2, sdk.NewDecWithPrec(7, 1)), // 70%
		),
		mustParseRFC3339("2021-08-02T00:00:00Z"),
		mustParseRFC3339("2021-08-10T00:00:00Z"),
		sdk.NewDecWithPrec(4, 2), // 4%,
	)

	handler := farming.NewHandler(app.FarmingKeeper)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	plan, found := app.FarmingKeeper.GetPlan(ctx, 1)
	require.Equal(t, true, found)

	require.Equal(t, msg.Name, plan.GetName())
	require.Equal(t, msg.Creator, plan.GetTerminationAddress().String())
	require.Equal(t, msg.StakingCoinWeights, plan.GetStakingCoinWeights())
	require.Equal(t, types.PrivatePlanFarmingPoolAddress(msg.Name, 1), plan.GetFarmingPoolAddress())
	require.Equal(t, mustParseRFC3339("2021-08-02T00:00:00Z"), plan.GetStartTime())
	require.Equal(t, mustParseRFC3339("2021-08-10T00:00:00Z"), plan.GetEndTime())
	require.Equal(t, msg.EpochRatio, plan.(*types.RatioPlan).EpochRatio)
}

func TestMsgStake(t *testing.T) {
	app, ctx, addrs := createTestInput()

	msg := types.NewMsgStake(
		addrs[0],
		sdk.NewCoins(sdk.NewInt64Coin(denom1, 10_000_000)),
	)

	handler := farming.NewHandler(app.FarmingKeeper)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	_, found := app.FarmingKeeper.GetQueuedStaking(ctx, denom1, addrs[0])
	require.Equal(t, true, found)

	queuedCoins := sdk.NewCoins()
	app.FarmingKeeper.IterateQueuedStakingsByFarmer(ctx, addrs[0], func(stakingCoinDenom string, queuedStaking types.QueuedStaking) (stop bool) {
		queuedCoins = queuedCoins.Add(sdk.NewCoin(stakingCoinDenom, queuedStaking.Amount))
		return false
	})
	require.Equal(t, msg.StakingCoins, queuedCoins)
}

func TestMsgUnstake(t *testing.T) {
	app, ctx, addrs := createTestInput()

	// stake some amount
	err := app.FarmingKeeper.Stake(ctx, addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 10_000_000)))
	require.NoError(t, err)

	_, found := app.FarmingKeeper.GetQueuedStaking(ctx, denom1, addrs[0])
	require.Equal(t, true, found)

	// check balance before unstake
	balanceBefore := app.BankKeeper.GetBalance(ctx, addrs[0], denom1)
	require.Equal(t, sdk.NewInt(990_000_000), balanceBefore.Amount)

	msg := types.NewMsgUnstake(
		addrs[0],
		sdk.NewCoins(sdk.NewInt64Coin(denom1, 5_000_000)),
	)

	handler := farming.NewHandler(app.FarmingKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	// check balance after unstake
	balanceAfter := app.BankKeeper.GetBalance(ctx, addrs[0], denom1)
	require.Equal(t, sdk.NewInt(995_000_000), balanceAfter.Amount)
}

func TestMsgHarvest(t *testing.T) {
	app, ctx, addrs := createTestInput()
	creator := addrs[0] // use addrs[0] to create a fixed amount plan
	staker := addrs[1]  // use addrs[1] to stake some amount with staking coin denom

	planMsg := types.NewMsgCreateFixedAmountPlan(
		"handler-test",
		creator,
		sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(denom1, sdk.NewDecWithPrec(10, 1)), // 100%
		),
		mustParseRFC3339("2021-08-02T00:00:00Z"),
		mustParseRFC3339("2021-08-10T00:00:00Z"),
		sdk.NewCoins(sdk.NewInt64Coin(denom3, 77_000_000)),
	)

	// create a fixed amount plan
	plan, err := app.FarmingKeeper.CreateFixedAmountPlan(
		ctx,
		planMsg,
		creator,
		creator,
		types.PlanTypePrivate,
	)
	require.NoError(t, err)

	_, found := app.FarmingKeeper.GetPlan(ctx, plan.GetId())
	require.Equal(t, true, found)

	// stake some amount
	err = app.FarmingKeeper.Stake(
		ctx,
		staker,
		sdk.NewCoins(sdk.NewInt64Coin(denom1, 10_000_000)),
	)
	require.NoError(t, err)

	_, found = app.FarmingKeeper.GetQueuedStaking(ctx, denom1, staker)
	require.Equal(t, true, found)

	// move queued coins into staked coins
	app.FarmingKeeper.ProcessQueuedCoins(ctx)

	_, found = app.FarmingKeeper.GetStaking(ctx, denom1, staker)
	require.Equal(t, true, found)

	// check balances before unstake
	balanceBefore := app.BankKeeper.GetBalance(ctx, staker, denom3)
	require.Equal(t, sdk.NewInt(1_000_000_000), balanceBefore.Amount)

	// allocate rewards
	ctx = ctx.WithBlockTime(mustParseRFC3339("2021-08-05T00:00:00Z"))
	err = app.FarmingKeeper.AllocateRewards(ctx)
	require.NoError(t, err)

	// harvest
	msg := types.NewMsgHarvest(staker, []string{denom1})
	handler := farming.NewHandler(app.FarmingKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	// check balances after unstake
	balanceAfter := app.BankKeeper.GetBalance(ctx, staker, denom3)
	require.Equal(t, sdk.NewInt(1_077_000_000), balanceAfter.Amount)
}

func mustParseRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
