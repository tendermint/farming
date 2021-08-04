package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/tendermint/farming/x/farming/types"
)

// GetQueryCmd returns a root CLI command handler for all x/farming query commands.
func GetQueryCmd() *cobra.Command {
	farmingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the farming module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	farmingQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryPlans(),
		GetCmdQueryPlan(),
		GetCmdQueryStakings(),
		GetCmdQueryRewards(),
	)

	return farmingQueryCmd
}

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current farming parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as farming parameters.

Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&resp.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPlans() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plans",
		Args:  cobra.NoArgs,
		Short: "Query for all farming plans",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all farming plans on a network.

Example:
$ %s query %s plans
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			resp, err := queryClient.Plans(cmd.Context(), &types.QueryPlansRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "plans")

	return cmd
}

func GetCmdQueryPlan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan [plan-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a specific plan.

Example:
$ %s query %s plan
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			planId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "plan-id %s is not valid", args[0])
			}

			resp, err := queryClient.Plan(cmd.Context(), &types.QueryPlanRequest{
				PlanId: planId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryStakings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stakings [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query for all stakings",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all farming stakings on a network.

Example:
$ %s query %s stakings
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			if err := sdk.ValidateDenom(args[0]); err != nil {
				return err
			}

			resp, err := queryClient.Stakings(cmd.Context(), &types.QueryStakingsRequest{
				StakingCoinDenom: args[0],
				Pagination:       pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "stakings")

	return cmd
}

func GetCmdQueryRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards [farmer-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query for all rewards from a farmer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query rewards that are accumulated on a network from a farmer.

Example:
$ %s query %s rewards %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, types.ModuleName, sdk.Bech32MainPrefix,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			farmerAcc, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if err := sdk.ValidateDenom(args[0]); err != nil {
				return err
			}

			resp, err := queryClient.Rewards(cmd.Context(), &types.QueryRewardsRequest{
				Farmer:           farmerAcc.String(),
				StakingCoinDenom: args[0],
				Pagination:       pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "rewards")

	return cmd
}
