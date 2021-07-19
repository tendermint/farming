package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"

	"github.com/tendermint/farming/x/farming/types"
)

// GetTxCmd returns a root CLI command handler for all x/farming transaction commands.
func GetTxCmd() *cobra.Command {
	farmingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Farming transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	farmingTxCmd.AddCommand(
		NewCreateFixedAmountPlanCmd(),
		NewCreateRatioPlanCmd(),
		NewStakeCmd(),
		NewUnstakeCmd(),
		NewHarvestCmd(),
	)

	return farmingTxCmd
}

func NewCreateFixedAmountPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-fixed-plan",
		Aliases: []string{"cf"},
		Args:    cobra.ExactArgs(2),
		Short:   "create fixed amount farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create fixed amount farming plan.
Example:
$ %s tx %s create-fixed-plan --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			farmingPoolAddr := sdk.AccAddress{}
			stakingCoinWeights := sdk.DecCoins{}
			startTime := time.Time{}
			endTime := time.Time{}
			epochAmount := sdk.Coins{}

			msg := types.NewMsgCreateFixedAmountPlan(
				farmingPoolAddr,
				stakingCoinWeights,
				startTime,
				endTime,
				epochAmount,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func NewCreateRatioPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-ratio-plan",
		Aliases: []string{"cr"},
		Args:    cobra.ExactArgs(2),
		Short:   "create ratio farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create ratio farming plan.
Example:
$ %s tx %s create-ratio-plan --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			farmingPoolAddr := sdk.AccAddress{}
			stakingCoinWeights := sdk.DecCoins{}
			startTime := time.Time{}
			endTime := time.Time{}
			epochRatio := sdk.Dec{}

			msg := types.NewMsgCreateRatioPlan(
				farmingPoolAddr,
				stakingCoinWeights,
				startTime,
				endTime,
				epochRatio,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func NewStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake",
		Args:  cobra.ExactArgs(2),
		Short: "stake coins to the farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`stake coins to the farming plan.
Example:
$ %s tx %s stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			farmer := sdk.AccAddress{}
			stakingCoins := sdk.Coins{}

			msg := types.NewMsgStake(
				farmer,
				stakingCoins,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func NewUnstakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unstake",
		Args:  cobra.ExactArgs(2),
		Short: "unstake coins from the farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`unstake coins from the farming plan.
Example:
$ %s tx %s unstake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			farmer := sdk.AccAddress{}
			stakingCoins := sdk.Coins{}

			msg := types.NewMsgUnstake(
				farmer,
				stakingCoins,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func NewHarvestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "harvest",
		Args:  cobra.ExactArgs(2),
		Short: "harvest farming rewards from the farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`claim farming rewards from the farming plan.
Example:
$ %s tx %s claim --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			farmer := sdk.AccAddress{}
			stakingCoinDenoms := []string{"test"}

			msg := types.NewMsgHarvest(farmer, stakingCoinDenoms)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

// TODO: not implemented yet
// GetCmdSubmitAddPublicPlanProposal implements a command handler for submitting a public farming plan creation transaction.
func GetCmdSubmitAddPublicPlanProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-public-farming-plan [proposal-file] [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a public farming plan creation",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a a public farming plan creation along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal public-farming-plan <path/to/proposal.json> --from=<key_or_address> --deposit=<deposit_amount>

`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			title := "title"
			desc := "description"
			coinWeights := sdk.NewDecCoins(sdk.NewDecCoin("test", sdk.NewInt(1)))
			startTime := time.Now().UTC()
			endTime := startTime.AddDate(1, 0, 0)

			content, err := types.NewAddPublicPlanProposal(title, desc, coinWeights, startTime, endTime)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			msg, err := gov.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	return cmd
}

// TODO: not implemented yet
// GetCmdSubmitUpdatePublicPlanProposal implements a command handler for submitting a public farming plan creation transaction.
func GetCmdSubmitUpdatePublicPlanProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-public-farming-plan [proposal-file] [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a public farming plan creation",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a a public farming plan creation along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal public-farming-plan <path/to/proposal.json> --from=<key_or_address> --deposit=<deposit_amount>

`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			title := "title"
			desc := "description"
			coinWeights := sdk.NewDecCoins(sdk.NewDecCoin("test", sdk.NewInt(1)))
			startTime := time.Now().UTC()
			endTime := startTime.AddDate(1, 0, 0)

			content, err := types.NewAddPublicPlanProposal(title, desc, coinWeights, startTime, endTime)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			msg, err := gov.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	return cmd
}

// TODO: not implemented yet
// GetCmdSubmitDeletePublicPlanProposal implements a command handler for submitting a public farming plan creation transaction.
func GetCmdSubmitDeletePublicPlanProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-public-farming-plan [proposal-file] [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a public farming plan creation",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a a public farming plan creation along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal public-farming-plan <path/to/proposal.json> --from=<key_or_address> --deposit=<deposit_amount>

`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			title := "title"
			desc := "description"
			coinWeights := sdk.NewDecCoins(sdk.NewDecCoin("test", sdk.NewInt(1)))
			startTime := time.Now().UTC()
			endTime := startTime.AddDate(1, 0, 0)

			content, err := types.NewAddPublicPlanProposal(title, desc, coinWeights, startTime, endTime)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			msg, err := gov.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	return cmd
}
