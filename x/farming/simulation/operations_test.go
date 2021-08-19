package simulation_test

import (
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	farmingapp "github.com/tendermint/farming/app"
	farmingparams "github.com/tendermint/farming/app/params"
	"github.com/tendermint/farming/x/farming/simulation"
	"github.com/tendermint/farming/x/farming/types"
)

// TestWeightedOperations tests the weights of the operations.
func TestWeightedOperations(t *testing.T) {
	app, ctx := createTestApp(false)

	ctx.WithChainID("test-chain")

	cdc := app.AppCodec()
	appParams := make(simtypes.AppParams)

	weightedOps := simulation.WeightedOperations(appParams, cdc, app.AccountKeeper, app.BankKeeper, app.FarmingKeeper)

	s := rand.NewSource(1)
	r := rand.New(s)
	accs := simtypes.RandomAccounts(r, 3)

	expected := []struct {
		weight     int
		opMsgRoute string
		opMsgName  string
	}{
		{farmingparams.DefaultWeightMsgCreateFixedAmountPlan, types.ModuleName, types.TypeMsgCreateFixedAmountPlan},
		{farmingparams.DefaultWeightMsgCreateRatioPlan, types.ModuleName, types.TypeMsgCreateRatioPlan},
		{farmingparams.DefaultWeightMsgStake, types.ModuleName, types.TypeMsgStake},
		{farmingparams.DefaultWeightMsgUnstake, types.ModuleName, types.TypeMsgUnstake},
		{farmingparams.DefaultWeightMsgHarvest, types.ModuleName, types.TypeMsgHarvest},
	}

	for i, w := range weightedOps {
		operationMsg, _, _ := w.Op()(r, app.BaseApp, ctx, accs, ctx.ChainID())
		// the following checks are very much dependent from the ordering of the output given
		// by WeightedOperations. if the ordering in WeightedOperations changes some tests
		// will fail
		require.Equal(t, expected[i].weight, w.Weight(), "weight should be the same")
		require.Equal(t, expected[i].opMsgRoute, operationMsg.Route, "route should be the same")
		require.Equal(t, expected[i].opMsgName, operationMsg.Name, "operation Msg name should be the same")
	}
}

// TestSimulateMsgCreateFixedAmountPlan tests the normal scenario of a valid message of type TypeMsgCreateFixedAmountPlan.
// Abnormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgCreateFixedAmountPlan(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// setup randomly generated liquidity pool creation fees
	feeCoins := simulation.GenPrivatePlanCreationFee(r)
	params := app.FarmingKeeper.GetParams(ctx)
	params.PrivatePlanCreationFee = feeCoins
	app.FarmingKeeper.SetParams(ctx, params)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgCreateFixedAmountPlan(app.AccountKeeper, app.BankKeeper, app.FarmingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgCreateFixedAmountPlan
	require.NoError(t, types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg))

	require.True(t, operationMsg.OK)
	require.Equal(t, types.TypeMsgCreateFixedAmountPlan, msg.Type())
	require.Len(t, futureOperations, 0)
}

func createTestApp(isCheckTx bool) (*farmingapp.FarmingApp, sdk.Context) {
	app := farmingapp.Setup(false)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, minttypes.DefaultParams())
	app.MintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	return app, ctx
}

func getTestingAccounts(t *testing.T, r *rand.Rand, app *farmingapp.FarmingApp, ctx sdk.Context, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := app.StakingKeeper.TokensFromConsensusPower(ctx, 200)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, account.Address)
		app.AccountKeeper.SetAccount(ctx, acc)
		err := simapp.FundAccount(app.BankKeeper, ctx, account.Address, initCoins)
		require.NoError(t, err)
	}

	return accounts
}
