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
	// app, ctx, addrs := createTestInput()

	// err := app.FarmingKeeper.Stake(ctx, addrs[0], sdk.NewCoins(sdk.NewInt64Coin(denom1, 10_000_000)))
	// require.NoError(t, err)

	// _, found := app.FarmingKeeper.GetQueuedStaking(ctx, denom1, addrs[0])
	// require.Equal(t, true, found)

	// msg := types.NewMsgUnstake(
	// 	addrs[0],
	// 	sdk.NewCoins(sdk.NewInt64Coin(denom1, 5_000_000)),
	// )

	// handler := farming.NewHandler(app.FarmingKeeper)
	// _, err = handler(ctx, msg)
	// require.NoError(t, err)
}

func TestMsgHarvest(t *testing.T) {
	// TODO: not implemented yet
}

func mustParseRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
