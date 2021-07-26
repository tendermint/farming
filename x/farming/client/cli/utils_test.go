package cli_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/app/params"
	"github.com/tendermint/farming/x/farming/client/cli"
	"github.com/tendermint/farming/x/farming/keeper"
	"github.com/tendermint/farming/x/farming/types"
)

func createTestInput() (*simapp.FarmingApp, sdk.Context) {
	app := simapp.Setup(false)
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

func TestParsePrivateFixedPlan(t *testing.T) {
	okJSON := testutil.WriteToNewTempFile(t, `
{
  "staking_coin_weights": [
	  {
	      "denom": "PoolCoinDenom",
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
`)

	plan, err := cli.ParsePrivateFixedPlan(okJSON.Name())
	require.NoError(t, err)

	require.Equal(t, "1.000000000000000000PoolCoinDenom", plan.StakingCoinWeights.String())
	require.Equal(t, "2021-07-15T08:41:21.662422Z", plan.StartTime.Format(time.RFC3339Nano))
	require.Equal(t, "2022-07-16T08:41:21.662422Z", plan.EndTime.Format(time.RFC3339Nano))
	require.Equal(t, "1uatom", plan.EpochAmount.String())
}

func TestParsePrivateRatioPlan(t *testing.T) {
	okJSON := testutil.WriteToNewTempFile(t, `
{
  "staking_coin_weights": [
	  {
	      "denom": "PoolCoinDenom",
	      "amount": "1.000000000000000000"
	  }
  ],
  "start_time": "2021-07-15T08:41:21.662422Z",
  "end_time": "2022-07-16T08:41:21.662422Z",
  "epoch_ratio":"1.000000000000000000"
}
`)

	plan, err := cli.ParsePrivateRatioPlan(okJSON.Name())
	require.NoError(t, err)

	require.Equal(t, "1.000000000000000000PoolCoinDenom", plan.StakingCoinWeights.String())
	require.Equal(t, "2021-07-15T08:41:21.662422Z", plan.StartTime.Format(time.RFC3339Nano))
	require.Equal(t, "2022-07-16T08:41:21.662422Z", plan.EndTime.Format(time.RFC3339Nano))
	require.Equal(t, "1.000000000000000000", plan.EpochRatio.String())
}

func TestParsePublicPlanProposal(t *testing.T) {
	encodingConfig := params.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Public Farming Plan",
  "description": "Are you ready to farm?",
  "name": "First Public Farming Plan",
  "add_request_proposals": [
    {
      "farming_pool_address": "cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu",
      "termination_address": "cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu",
      "staking_coin_weights": [
        {
          "denom": "PoolCoinDenom",
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
  ]
}
`)

	proposal, err := cli.ParsePublicPlanProposal(encodingConfig.Marshaler, okJSON.Name())
	require.NoError(t, err)

	require.Equal(t, "Public Farming Plan", proposal.Title)
	require.Equal(t, "Are you ready to farm?", proposal.Description)
	require.Equal(t, "First Public Farming Plan", proposal.Name)
}
