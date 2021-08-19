package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
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
		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedAmountPlan, "unable to generate fees"), nil, err
		}

		txGen := params.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCreateRatioPlan generates a MsgCreateRatioPlan with random values
// nolint: interfacer
func SimulateMsgCreateRatioPlan(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: not implemented yet
		return simtypes.OperationMsg{}, nil, nil
	}
}

// SimulateMsgStake generates a MsgCreateFixedAmountPlan with random values
// nolint: interfacer
func SimulateMsgStake(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: not implemented yet
		return simtypes.OperationMsg{}, nil, nil
	}
}

// SimulateMsgUnstake generates a SimulateMsgUnstake with random values
// nolint: interfacer
func SimulateMsgUnstake(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: not implemented yet
		return simtypes.OperationMsg{}, nil, nil
	}
}

// SimulateMsgHarvest generates a MsgHarvest with random values
// nolint: interfacer
func SimulateMsgHarvest(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: not implemented yet
		return simtypes.OperationMsg{}, nil, nil
	}
}
