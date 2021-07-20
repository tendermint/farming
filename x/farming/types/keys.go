package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the farming module
	ModuleName = "farming"

	// RouterKey is the message router key for the farming module
	RouterKey = ModuleName

	// StoreKey is the default store key for the farming module
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the farming module
	QuerierRoute = ModuleName

	// RewardPoolAccKeyPrefix is prefix for generating deterministic reward pool module account of the each plan
	RewardPoolAccKeyPrefix = "RewardPoolAcc"

	// StakingReserveAccKeyPrefix is prefix for generating deterministic staking reserve module account of the each plan
	StakingReserveAccKeyPrefix = "StakingReserveAcc"
)

var (
	// param key for global farming plan IDs
	GlobalPlanIdKey           = []byte("globalPlanId")
	GlobalLastEpochTimePrefix = []byte("globalLastEpochTime")
	GlobalStakingIdKey        = []byte("globalStakingId")

	PlanKeyPrefix                     = []byte{0x11}
	PlanByFarmerAddrIndexKeyPrefix    = []byte{0x12}
	LastDistributedTimeKeyPrefix      = []byte{0x13}
	TotalDistributedRewardCoinsPrefix = []byte{0x14}

	StakingKeyPrefix                          = []byte{0x21}
	StakingByFarmerAddrIndexKeyPrefix         = []byte{0x22}
	StakingByStakingCoinDenomIdIndexKeyPrefix = []byte{0x23}

	RewardKeyPrefix                  = []byte{0x31}
	RewardByFarmerAddrIndexKeyPrefix = []byte{0x32}

	StakingReserveAcc = sdk.AccAddress(address.Module(ModuleName, []byte("StakingReserveAcc")))
)

// GetPlanKey returns kv indexing key of the plan
func GetPlanKey(planID uint64) []byte {
	key := make([]byte, 9)
	key[0] = PlanKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(planID))
	return key
}

// GetPlansByFarmerAddrIndexKey returns kv indexing key of the plan indexed by reserve account
func GetPlansByFarmerAddrIndexKey(farmerAcc sdk.AccAddress) []byte {
	return append(PlanByFarmerAddrIndexKeyPrefix, address.MustLengthPrefix(farmerAcc.Bytes())...)
}

// GetPlanByFarmerAddrIndexKey returns kv indexing key of the plan indexed by reserve account
func GetPlanByFarmerAddrIndexKey(farmerAcc sdk.AccAddress, planID uint64) []byte {
	return append(append(PlanByFarmerAddrIndexKeyPrefix, address.MustLengthPrefix(farmerAcc.Bytes())...), sdk.Uint64ToBigEndian(planID)...)
}

// GetStakingPrefix returns prefix of staking records in the plan
func GetStakingPrefix(planID uint64) []byte {
	key := make([]byte, 9)
	key[0] = StakingKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(planID))
	return key
}

// GetStakingKey returns a key for staking of corresponding the id
func GetStakingKey(id uint64) []byte {
	return append(StakingKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

// GetStakingByFarmerAddrIndexKey returns key for the farmer's staking of corresponding
func GetStakingByFarmerAddrIndexKey(farmerAcc sdk.AccAddress) []byte {
	return append(StakingByFarmerAddrIndexKeyPrefix, address.MustLengthPrefix(farmerAcc.Bytes())...)
}

// GetStakingByStakingCoinDenomIdIndexPrefix returns prefix for the iterable staking list by the staking coin denomination
func GetStakingByStakingCoinDenomIdIndexPrefix(denom string) []byte {
	return append(StakingByStakingCoinDenomIdIndexKeyPrefix, MustLengthPrefixString(denom)...)
}

// GetStakingByStakingCoinDenomIdIndexKey returns key for the staking index by the staking coin denomination
func GetStakingByStakingCoinDenomIdIndexKey(denom string, id uint64) []byte {
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, id)
	return append(append(StakingByStakingCoinDenomIdIndexKeyPrefix, MustLengthPrefixString(denom)...), idBytes...)
}

// MustLengthPrefixString is LengthPrefix with panic on error.
// TODO: update documentation
func MustLengthPrefixString(str string) []byte {
	bz := []byte(str)
	bzLen := len(bz)
	if bzLen == 0 {
		return bz
	}
	return append([]byte{byte(bzLen)}, bz...)
}

// GetRewardKey returns key for staking coin denomination's reward of corresponding the farmer
func GetRewardKey(stakingCoinDenom string, farmerAcc sdk.AccAddress) []byte {
	return append(append(RewardKeyPrefix, MustLengthPrefixString(stakingCoinDenom)...), address.MustLengthPrefix(farmerAcc.Bytes())...)
}

// GetRewardByFarmerAddrIndexKey returns key for farmer's reward of corresponding the staking coin denomination
func GetRewardByFarmerAddrIndexKey(farmerAcc sdk.AccAddress, stakingCoinDenom string) []byte {
	return append(append(RewardByFarmerAddrIndexKeyPrefix, address.MustLengthPrefix(farmerAcc.Bytes())...), MustLengthPrefixString(stakingCoinDenom)...)
}

// GetRewardByStakingCoinDenomPrefix returns prefix for staking coin denomination's reward list
func GetRewardByStakingCoinDenomPrefix(stakingCoinDenom string) []byte {
	return append(RewardKeyPrefix, MustLengthPrefixString(stakingCoinDenom)...)
}

// GetRewardByFarmerAddrIndexPrefix returns prefix for farmer's reward list
func GetRewardByFarmerAddrIndexPrefix(farmerAcc sdk.AccAddress) []byte {
	return append(RewardByFarmerAddrIndexKeyPrefix, address.MustLengthPrefix(farmerAcc.Bytes())...)
}

// ParseRewardKey parses a RewardKey.
func ParseRewardKey(key []byte) (stakingCoinDenom string, farmer sdk.AccAddress) {
	denomLen := key[1]
	stakingCoinDenom = string(key[2 : 2+denomLen])
	farmer = key[2+denomLen+1:]
	return
}

// ParseRewardByFarmerAddrIndexKey parses a key of RewardByFarmerAddrIndex from bytes.
func ParseRewardByFarmerAddrIndexKey(key []byte) (farmer sdk.AccAddress, stakingCoinDenom string) {
	addrLen := key[1]
	farmer = key[2 : 2+addrLen]
	stakingCoinDenom = string(key[2+addrLen+1:])
	return
}

// ParseStakingByStakingCoinDenomIdIndexKey parses a key of StakingByStakingCoinDenomIdIndex from bytes.
func ParseStakingByStakingCoinDenomIdIndexKey(bz []byte) (stakingCoinDenom string, stakingID uint64) {
	denomLen := bz[1]
	stakingCoinDenom = string(bz[2 : 2+denomLen])
	stakingID = binary.BigEndian.Uint64(bz[2+denomLen:])
	return
}
