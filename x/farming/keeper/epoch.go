package keeper

import (
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming/types"
)

// GetEpochEndTime returns the last time the epoch ended in UTC format.
func (k Keeper) GetEpochEndTime(ctx sdk.Context) (t time.Time, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EpochEndTimeKey)
	if bz == nil {
		return
	}
	var ts gogotypes.Timestamp
	k.cdc.MustUnmarshal(bz, &ts)
	var err error
	t, err = gogotypes.TimestampFromProto(&ts)
	if err != nil {
		panic(err)
	}
	found = true
	return
}

// SetEpochEndTime sets the last time the epoch ended.
func (k Keeper) SetEpochEndTime(ctx sdk.Context, t time.Time) {
	store := ctx.KVStore(k.storeKey)
	ts, err := gogotypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	bz := k.cdc.MustMarshal(ts)
	store.Set(types.EpochEndTimeKey, bz)
}


// GetNextEpochDuration returns the current epoch days(period).
func (k Keeper) GetNextEpochDuration(ctx sdk.Context) uint32 {
	var epochDays uint32
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextEpochDaysKey)
	if bz == nil {
		// initialize with next epoch days
		epochDays = k.GetParams(ctx).NextEpochDays
	} else {
		val := gogotypes.UInt32Value{}
		if err := k.cdc.Unmarshal(bz, &val); err != nil {
			panic(err)
		}
		epochDays = val.GetValue()
	}
	return epochDays
}

// SetNextEpochDuration sets the current epoch days(period).
func (k Keeper) SetNextEpochDuration(ctx sdk.Context, days uint32) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt32Value{Value: days})
	store.Set(types.NextEpochDaysKey, bz)
}
