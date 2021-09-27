package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState returns new GenesisState.
func NewGenesisState(
	params Params, plans []PlanRecord, stakings []StakingRecord, queuedStakings []QueuedStakingRecord,
	historicalRewards []HistoricalRewardsRecord, outstandingRewards []OutstandingRewardsRecord,
	currentEpochs []CurrentEpochRecord, stakingReserveCoins, rewardPoolCoins sdk.Coins,
	lastEpochTime *time.Time, currentEpochDays uint32,
) *GenesisState {
	return &GenesisState{
		Params:                    params,
		PlanRecords:               plans,
		StakingRecords:            stakings,
		QueuedStakingRecords:      queuedStakings,
		HistoricalRewardsRecords:  historicalRewards,
		OutstandingRewardsRecords: outstandingRewards,
		CurrentEpochRecords:       currentEpochs,
		StakingReserveCoins:       stakingReserveCoins,
		RewardPoolCoins:           rewardPoolCoins,
		LastEpochTime:             lastEpochTime,
		CurrentEpochDays:          currentEpochDays,
	}
}

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		DefaultParams(),
		[]PlanRecord{},
		[]StakingRecord{},
		[]QueuedStakingRecord{},
		[]HistoricalRewardsRecord{},
		[]OutstandingRewardsRecord{},
		[]CurrentEpochRecord{},
		sdk.Coins{},
		sdk.Coins{},
		nil,
		DefaultCurrentEpochDays,
	)
}

// ValidateGenesis validates GenesisState.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	id := uint64(0)

	var plans []PlanI
	for _, record := range data.PlanRecords {
		plan, err := UnpackPlan(&record.Plan)
		if err != nil {
			return err
		}
		if err := plan.Validate(); err != nil {
			return err
		}
		if plan.GetId() < id {
			return fmt.Errorf("pool records must be sorted")
		}
		plans = append(plans, plan)
		id = plan.GetId() + 1
	}

	if err := ValidatePlanNames(plans); err != nil {
		return err
	}

	if err := ValidateTotalEpochRatio(plans); err != nil {
		return err
	}

	// TODO: validate other fields

	return nil
}
