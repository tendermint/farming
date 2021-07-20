package cli_test

import (
	"io/ioutil"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	farmingapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming/keeper"
	"github.com/tendermint/farming/x/farming/types"
)

func createTestInput() (*farmingapp.FarmingApp, sdk.Context) {
	app := farmingapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.FarmingKeeper = keeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.GetSubspace(types.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		map[string]bool{},
	)

	return app, ctx
}

func TestParseJSONFile(t *testing.T) {
	app, _ := createTestInput()

	proposalFile := "./proposal.json"

	proposal := types.PublicPlanProposal{}

	contents, err := ioutil.ReadFile(proposalFile)
	require.NoError(t, err)

	err = app.AppCodec().UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)
}

// func TestMarshalPublic(t *testing.T) {
// 	app, _ := createTestInput()

// 	farmingPoolAddr := sdk.AccAddress([]byte("farmingPoolAddr"))
// 	terminationAddr := sdk.AccAddress([]byte("terminationAddr"))
// 	coinWeights := sdk.NewDecCoins(sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")})
// 	startTime := time.Now().UTC()
// 	endTime := startTime.AddDate(1, 0, 0)
// 	name := ""

// 	proposal := types.AddPublicPlanProposal{}
// 	proposal.Title = "Public Plan Test"
// 	proposal.Description = "TEST..."

// 	basePlan := types.NewBasePlan(
// 		1,
// 		name,
// 		types.PlanTypePublic,
// 		farmingPoolAddr.String(),
// 		terminationAddr.String(),
// 		coinWeights,
// 		startTime,
// 		endTime,
// 	)
// 	epochRatio := sdk.NewDec(1.0)

// 	ratioPlan := types.NewRatioPlan(basePlan, epochRatio)

// 	plans, err := types.PackPlans([]types.PlanI{ratioPlan})
// 	require.NoError(t, err)

// 	proposal.Plans = plans

// 	bz, err := app.AppCodec().MarshalJSON(&proposal)
// 	require.NoError(t, err)

// 	fmt.Println("bz: ", string(bz))
// }
