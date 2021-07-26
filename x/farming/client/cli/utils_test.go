package cli_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	farmingapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming/client/cli"
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

func TestParsePrivateFixedPlan(t *testing.T) {
	fixedPlanStr := `{
  "staking_coin_weights": [
	  {
	      "denom": "poolCoinDenom",
	      "amount": "1.000000000000000000"
	  }
  ],
  "start_time": "2021-07-15T08:41:21.662422Z",
  "end_time": "2022-07-16T08:41:21.662422Z",
  "epoch_amount": [
    {
      "denom": "uatom",
      "amount": "1"
    }
  ]
}
`
	plan := cli.PrivateFixedPlanRequest{}

	contents := []byte(fixedPlanStr)
	err := json.Unmarshal(contents, &plan)
	require.NoError(t, err)

	require.Equal(t, "1.000000000000000000poolCoinDenom", plan.StakingCoinWeights.String())
	require.Equal(t, "1uatom", plan.EpochAmount.String())
}
