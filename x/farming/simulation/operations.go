package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/tendermint/farming/app/params"
	"github.com/tendermint/farming/x/farming/keeper"
	"github.com/tendermint/farming/x/farming/types"
)

// Simulation operation weights constants.
const (
	OpWeightMsgCreateFixedAmountPlan = "op_weight_msg_create_fixed_amount_plan"
	OpWeightMsgCreateRatioPlan       = "op_weight_msg_create_ratio_plan"
	OpWeightMsgStake                 = "op_weight_msg_stake"
	OpWeightMsgUnstake               = "op_weight_msg_unstake"
	OpWeightMsgHarvest               = "op_weight_msg_harvest"
)

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper,
) simulation.WeightedOperations {

	var weightMsgCreateFixedAmountPlan int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateFixedAmountPlan, &weightMsgCreateFixedAmountPlan, nil,
		func(_ *rand.Rand) {
			weightMsgCreateFixedAmountPlan = params.DefaultWeightMsgCreateFixedAmountPlan
		},
	)

	var weightMsgCreateRatioPlan int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateRatioPlan, &weightMsgCreateRatioPlan, nil,
		func(_ *rand.Rand) {
			weightMsgCreateRatioPlan = params.DefaultWeightMsgCreateRatioPlan
		},
	)

	var weightMsgStake int
	appParams.GetOrGenerate(cdc, OpWeightMsgStake, &weightMsgStake, nil,
		func(_ *rand.Rand) {
			weightMsgStake = params.DefaultWeightMsgStake
		},
	)

	var weightMsgUnstake int
	appParams.GetOrGenerate(cdc, OpWeightMsgUnstake, &weightMsgUnstake, nil,
		func(_ *rand.Rand) {
			weightMsgUnstake = params.DefaultWeightMsgUnstake
		},
	)

	var weightMsgHarvest int
	appParams.GetOrGenerate(cdc, OpWeightMsgHarvest, &weightMsgHarvest, nil,
		func(_ *rand.Rand) {
			weightMsgHarvest = params.DefaultWeightMsgHarvest
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateFixedAmountPlan,
			SimulateMsgCreateFixedAmountPlan(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgCreateRatioPlan,
			SimulateMsgCreateRatioPlan(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgStake,
			SimulateMsgStake(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUnstake,
			SimulateMsgUnstake(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgHarvest,
			SimulateMsgHarvest(ak, bk, k),
		),
	}
}

// SimulateMsgCreateFixedAmountPlan generates a MsgCreateFixedAmountPlan with random values
// nolint: interfacer
func SimulateMsgCreateFixedAmountPlan(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params := k.GetParams(ctx)
		coins, hasNeg := spendable.SafeSub(params.PrivatePlanCreationFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedAmountPlan, "lower balance"),
				nil, fmt.Errorf("spendable %s is lower than plan creation fee %s", spendable.String(), coins.String())
		}

		name := "simulation"
		creatorAcc := account.GetAddress()
		stakingCoinWeights := sdk.NewDecCoins(sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 1))
		startTime := time.Now().UTC()
		endTime := startTime.AddDate(0, 0, 1)
		epochAmount := sdk.NewCoins(
			sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simtypes.RandIntBetween(r, 1_000_000, 1_000_000_000))),
		)

		msg := types.NewMsgCreateFixedAmountPlan(name, creatorAcc, stakingCoinWeights, startTime, endTime, epochAmount)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgCreateRatioPlan generates a MsgCreateRatioPlan with random values
// nolint: interfacer
func SimulateMsgCreateRatioPlan(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params := k.GetParams(ctx)
		coins, hasNeg := spendable.SafeSub(params.PrivatePlanCreationFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateRatioPlan, "lower balance"),
				nil, fmt.Errorf("spendable %s is lower than plan creation fee %s", spendable.String(), coins.String())
		}

		name := "simulation"
		creatorAcc := account.GetAddress()
		stakingCoinWeights := sdk.NewDecCoins(sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 1))
		startTime := time.Now().UTC()
		endTime := startTime.AddDate(0, 0, 1)
		epochRatio := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 10)), 1)

		msg := types.NewMsgCreateRatioPlan(name, creatorAcc, stakingCoinWeights, startTime, endTime, epochRatio)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgStake generates a MsgCreateFixedAmountPlan with random values
// nolint: interfacer
func SimulateMsgStake(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params := k.GetParams(ctx)
		coins, hasNeg := spendable.SafeSub(params.StakingCreationFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStake, "lower balance"),
				nil, fmt.Errorf("spendable %s is lower than staking creation fee %s", spendable.String(), coins.String())
		}

		farmer := account.GetAddress()
		stakingCoins := sdk.NewCoins(
			sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simtypes.RandIntBetween(r, 1_000_000, 100_000_000))),
		)

		msg := types.NewMsgStake(farmer, stakingCoins)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgUnstake generates a SimulateMsgUnstake with random values
// nolint: interfacer
func SimulateMsgUnstake(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		farmer := account.GetAddress()
		unstakingCoins := sdk.NewCoins(
			sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simtypes.RandIntBetween(r, 1_000_000, 100_000_000))),
		)

		// staking must exist in order to unharvest
		staking, found := k.GetStakingByFarmer(ctx, farmer)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstake, "unable to find staking"),
				nil, fmt.Errorf("staking by %s not found", farmer)
		}

		if !staking.StakedCoins.Add(staking.QueuedCoins...).IsAllGTE(unstakingCoins) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstake, "insufficient funds"),
				nil, fmt.Errorf("%s is smaller than %s", staking.StakedCoins.Add(staking.QueuedCoins...).String(), unstakingCoins.String())
		}

		msg := types.NewMsgUnstake(farmer, unstakingCoins)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgHarvest generates a MsgHarvest with random values
// nolint: interfacer
func SimulateMsgHarvest(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		farmer := account.GetAddress()
		stakingCoinDenoms := []string{sdk.DefaultBondDenom}

		// add a day to increase epoch if there is no harvest rewards
		rewards := k.GetRewardsByFarmer(ctx, farmer)
		if len(rewards) == 0 {
			ctx = ctx.WithBlockTime(ctx.BlockTime().AddDate(0, 0, 1))
			k.ProcessQueuedCoins(ctx)
			k.DistributeRewards(ctx)
			k.SetLastEpochTime(ctx, ctx.BlockTime())

			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgHarvest, "no rewards to harvest"), nil, nil
		}

		msg := types.NewMsgHarvest(farmer, stakingCoinDenoms)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
