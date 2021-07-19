package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/farming/x/farming/types"
)

// GetReward returns a specific reward.
func (k Keeper) GetReward(ctx sdk.Context, farmerAcc sdk.AccAddress, stakingCoinDenom string) (reward types.Reward, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRewardKey(stakingCoinDenom, farmerAcc))
	if bz == nil {
		return reward, false
	}
	k.cdc.MustUnmarshal(bz, &reward)
	return reward, true
}

// GetRewardsByFarmer reads from kvstore and return a specific Reward indexed by given farmer's address
func (k Keeper) GetRewardsByFarmer(ctx sdk.Context, farmer sdk.AccAddress) (rewards []types.Reward) {
	k.IterateRewardsByFarmer(ctx, farmer, func(_ sdk.AccAddress, _ string, reward types.Reward) bool {
		rewards = append(rewards, reward)
		return false
	})

	return rewards
}

// SetReward implements Reward.
func (k Keeper) SetReward(ctx sdk.Context, farmerAcc sdk.AccAddress, stakingCoinDenom string, reward types.Reward) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&reward)
	store.Set(types.GetRewardKey(stakingCoinDenom, farmerAcc), bz)
	store.Set(types.GetRewardByFarmerAddrIndexKey(farmerAcc, stakingCoinDenom), nil)
}

// DeleteReward deletes a reward for the reward mapper store.
func (k Keeper) DeleteReward(ctx sdk.Context, farmerAcc sdk.AccAddress, stakingCoinDenom string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetRewardKey(stakingCoinDenom, farmerAcc))
	store.Delete(types.GetRewardByFarmerAddrIndexKey(farmerAcc, stakingCoinDenom))
}

// IterateAllRewards iterates over all the stored rewards and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAllRewards(ctx sdk.Context, cb func(farmer sdk.AccAddress, stakingCoinDenom string, reward types.Reward) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RewardKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		farmer := types.GetFarmerAddrFromRewardKey(key)
		stakingCoinDenom := types.GetStakingCoinDenomFromRewardKey(key)
		var reward types.Reward
		k.cdc.MustUnmarshal(iterator.Value(), &reward)
		if cb(farmer, stakingCoinDenom, reward) {
			break
		}
	}
}

// IterateRewardsByStakingCoinDenom iterates over all the stored rewards indexed by given staking coin denom and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateRewardsByStakingCoinDenom(ctx sdk.Context, denom string, cb func(farmer sdk.AccAddress, stakingCoinDenom string, reward types.Reward) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetRewardByStakingCoinDenomPrefix(denom))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		farmer := types.GetFarmerAddrFromRewardKey(key)
		stakingCoinDenom := types.GetStakingCoinDenomFromRewardKey(key)
		var reward types.Reward
		k.cdc.MustUnmarshal(iterator.Value(), &reward)
		if cb(farmer, stakingCoinDenom, reward) {
			break
		}
	}
}

// IterateRewardsByFarmer iterates over all the stored rewards indexed by given farmer's address and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateRewardsByFarmer(ctx sdk.Context, farmer sdk.AccAddress, cb func(farmer sdk.AccAddress, stakingCoinDenom string, reward types.Reward) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetRewardByFarmerAddrIndexPrefix(farmer))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		stakingCoinDenom := types.GetStakingCoinDenomFromRewardByFarmerAddrIndexKey(iterator.Key())
		reward, found := k.GetReward(ctx, farmer, stakingCoinDenom)
		if !found { // TODO: remove this check
			panic("reward not found")
		}
		if cb(farmer, stakingCoinDenom, reward) {
			break
		}
	}
}

// UnmarshalReward unmarshals a Reward from bytes.
func (k Keeper) UnmarshalReward(bz []byte) (types.Reward, error) {
	var reward types.Reward
	return reward, k.cdc.Unmarshal(bz, &reward)
}

// Harvest claims farming rewards from the reward pool account.
func (k Keeper) Harvest(ctx sdk.Context, farmer sdk.AccAddress, stakingCoinDenoms []string) error {
	amount := sdk.NewCoins()
	for _, denom := range stakingCoinDenoms {
		reward, found := k.GetReward(ctx, farmer, denom)
		if !found {
			return sdkerrors.Wrapf(types.ErrRewardNotExists, "no reward for staking coin denom %s", denom)
		}
		amount = amount.Add(reward.RewardCoins...)
	}

	if err := k.ReleaseStakingCoins(ctx, farmer, amount); err != nil {
		return err
	}

	for _, denom := range stakingCoinDenoms {
		k.DeleteReward(ctx, farmer, denom)
	}

	if len(k.GetRewardsByFarmer(ctx, farmer)) == 0 {
		k.GetStakingIDByFarmer(ctx, farmer)
		staking, found := k.GetStakingByFarmer(ctx, farmer)
		if !found { // TODO: remove this check
			return fmt.Errorf("staking not found")
		}
		if staking.StakedCoins.IsZero() && staking.QueuedCoins.IsZero() {
			k.DeleteStaking(ctx, staking.Id)
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeHarvest,
			sdk.NewAttribute(types.AttributeKeyFarmer, farmer.String()),
			sdk.NewAttribute(types.AttributeKeyRewardCoins, amount.String()),
		),
	})

	return nil
}
