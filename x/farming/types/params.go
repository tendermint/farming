package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyPrivatePlanCreationFee = []byte("PrivatePlanCreationFee")
	KeyNextEpochDays          = []byte("NextEpochDays")
	KeyFarmingFeeCollector    = []byte("FarmingFeeCollector")
	KeyDelayedStakingGasFee   = []byte("DelayedStakingGasFee")

	DefaultPrivatePlanCreationFee = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000)))
	DefaultCurrentEpochDays       = uint32(1)
	DefaultNextEpochDays          = uint32(1)
	DefaultFarmingFeeCollector    = sdk.AccAddress(address.Module(ModuleName, []byte("FarmingFeeCollectorAcc"))).String()
	DefaultDelayedStakingGasFee   = sdk.Gas(60000) // See https://github.com/tendermint/farming/issues/102 for details.

	// TODO: remove global reserve account due to split the reserveAcc by staking coin denom
	StakingReserveAcc = sdk.AccAddress(address.Module(ModuleName, []byte("StakingReserveAcc")))
	RewardsReserveAcc = sdk.AccAddress(address.Module(ModuleName, []byte("RewardsReserveAcc")))

	// ReserveAccAddressType is Address type of reserve accounts for staking, rewards, It could be changed,
	// details on https://github.com/tendermint/farming/issues/200
	ReserveAccAddressType = AddressType32Bytes
)

var _ paramstypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns the default farming module parameters.
func DefaultParams() Params {
	return Params{
		PrivatePlanCreationFee: DefaultPrivatePlanCreationFee,
		NextEpochDays:          DefaultNextEpochDays,
		FarmingFeeCollector:    DefaultFarmingFeeCollector,
		DelayedStakingGasFee:   DefaultDelayedStakingGasFee,
	}
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyPrivatePlanCreationFee, &p.PrivatePlanCreationFee, validatePrivatePlanCreationFee),
		paramstypes.NewParamSetPair(KeyNextEpochDays, &p.NextEpochDays, validateNextEpochDays),
		paramstypes.NewParamSetPair(KeyFarmingFeeCollector, &p.FarmingFeeCollector, validateFarmingFeeCollector),
		paramstypes.NewParamSetPair(KeyDelayedStakingGasFee, &p.DelayedStakingGasFee, validateDelayedStakingGas),
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates parameters.
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.PrivatePlanCreationFee, validatePrivatePlanCreationFee},
		{p.NextEpochDays, validateNextEpochDays},
		{p.FarmingFeeCollector, validateFarmingFeeCollector},
		{p.DelayedStakingGasFee, validateDelayedStakingGas},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

func validatePrivatePlanCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := v.Validate(); err != nil {
		return err
	}

	return nil
}

func validateNextEpochDays(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("next epoch days must be positive: %d", v)
	}

	return nil
}

func validateFarmingFeeCollector(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("farming fee collector address must not be empty")
	}

	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("invalid account address: %v", v)
	}

	return nil
}

func validateDelayedStakingGas(i interface{}) error {
	_, ok := i.(sdk.Gas)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
