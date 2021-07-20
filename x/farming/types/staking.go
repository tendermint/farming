package types

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s Staking) String() string {
	out, _ := s.MarshalYAML()
	return out.(string)
}

func (s Staking) MarshalYAML() (interface{}, error) {
	bz, err := codec.MarshalYAML(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()), &s)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}

func (s Staking) GetFarmerAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(s.Farmer)
	return addr
}

func (s Staking) IdBytes() []byte {
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, s.Id)
	return idBytes
}

func (s Staking) Denoms() (denomList []string) {
	keys := make(map[string]bool)
	for _, coin := range s.QueuedCoins {
		if _, value := keys[coin.Denom]; !value {
			keys[coin.Denom] = true
			denomList = append(denomList, coin.Denom)
		}
	}
	for _, coin := range s.StakedCoins {
		if _, value := keys[coin.Denom]; !value {
			keys[coin.Denom] = true
			denomList = append(denomList, coin.Denom)
		}
	}
	return denomList
}
