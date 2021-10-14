package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// farming module sentinel errors
var (
	ErrInvalidPlanType                = sdkerrors.Register(ModuleName, 2, "invalid plan type")
	ErrInvalidPlanName                = sdkerrors.Register(ModuleName, 3, "invalid plan name")
	ErrInvalidPlanNameLength          = sdkerrors.Register(ModuleName, 4, "invalid plan name length")
	ErrInvalidPlanEndTime             = sdkerrors.Register(ModuleName, 5, "invalid plan end time")
	ErrInvalidStakingCoinWeights      = sdkerrors.Register(ModuleName, 6, "invalid staking coin weights")
	ErrInvalidTotalEpochRatio         = sdkerrors.Register(ModuleName, 7, "invalid total epoch ratio")
	ErrStakingNotExists               = sdkerrors.Register(ModuleName, 8, "staking not exists")
	ErrConflictPrivatePlanFarmingPool = sdkerrors.Register(ModuleName, 9, "the address is already in use, please use a different plan name")
	ErrInvalidStakingReservedAmount   = sdkerrors.Register(ModuleName, 10, "staking reserved amount invariant broken")
	ErrInvalidRemainingRewardsAmount  = sdkerrors.Register(ModuleName, 11, "remaining rewards amount invariant broken")
)
