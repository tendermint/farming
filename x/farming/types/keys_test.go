package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming/types"
)

type keysTestSuite struct {
	suite.Suite
}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(keysTestSuite))
}

func (s *keysTestSuite) TestGetPlanKey() {
	s.Require().Equal([]byte{0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa}, types.GetPlanKey(10))
	s.Require().Equal([]byte{0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9}, types.GetPlanKey(9))
	s.Require().Equal([]byte{0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, types.GetPlanKey(0))
}

func (s *keysTestSuite) TestGetPlansByFarmerIndexKey() {
	s.Require().Equal([]byte{0x12}, types.GetPlansByFarmerIndexKey(sdk.AccAddress("")))
	s.Require().Equal([]byte{0x12, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x31}, types.GetPlansByFarmerIndexKey(sdk.AccAddress("farmer1")))
	s.Require().Equal([]byte{0x12, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x32}, types.GetPlansByFarmerIndexKey(sdk.AccAddress("farmer2")))
	s.Require().Equal([]byte{0x12, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x33}, types.GetPlansByFarmerIndexKey(sdk.AccAddress("farmer3")))
}

func (s *keysTestSuite) TestGetPlanByFarmerAddrIndexKey() {
	testCases := []struct {
		farmerAcc sdk.AccAddress
		planID    uint64
		expected  []byte
	}{
		{
			nil,
			1,
			[]byte{0x12, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
		},
		{
			sdk.AccAddress("farmer1"),
			1,
			[]byte{0x12, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
		},
		{
			sdk.AccAddress("farmer2"),
			2,
			[]byte{0x12, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x32, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2},
		},
		{
			sdk.AccAddress("farmer3"),
			2,
			[]byte{0x12, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x33, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2},
		},
	}

	for _, tc := range testCases {
		s.Require().Equal(tc.expected, types.GetPlanByFarmerAddrIndexKey(tc.farmerAcc, tc.planID))
	}
}

func (s *keysTestSuite) TestGetStakingKey() {
	testCases := []struct {
		stakingCoinDenom string
		farmerAcc        sdk.AccAddress
		expected         []byte
	}{
		// TODO: see the first case of TestGetStakingIndexKey
		{
			sdk.DefaultBondDenom,
			sdk.AccAddress(""),
			[]byte{0x21, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65},
		},
		{
			sdk.DefaultBondDenom,
			sdk.AccAddress("farmer1"),
			[]byte{0x21, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x31},
		},
		{
			sdk.DefaultBondDenom,
			sdk.AccAddress("farmer2"),
			[]byte{0x21, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x32},
		},
		{
			sdk.DefaultBondDenom,
			sdk.AccAddress("farmer3"),
			[]byte{0x21, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x33},
		},
	}

	for _, tc := range testCases {
		key := types.GetStakingKey(tc.stakingCoinDenom, tc.farmerAcc)
		s.Require().Equal(tc.expected, key)

		stakingCoinDenom, farmerAcc := types.ParseStakingKey(key)
		s.Require().Equal(tc.stakingCoinDenom, stakingCoinDenom)
		s.Require().Equal(tc.farmerAcc, farmerAcc)
	}
}

func (s *keysTestSuite) TestGetStakingIndexKey() {
	testCases := []struct {
		farmerAcc        sdk.AccAddress
		stakingCoinDenom string
		expected         []byte
	}{
		// TODO: should we cover only happy cases? below case returns panic since farmerAcc is empty
		// How about the case1 for TestGetStakingKey()? It allows empty farmerAcc.
		// {
		// 	sdk.AccAddress(""),
		// 	sdk.DefaultBondDenom,
		// 	[]byte{0x22, 0x73, 0x74, 0x61, 0x6b, 0x65},
		// },
		{
			sdk.AccAddress("farmer1"),
			sdk.DefaultBondDenom,
			[]byte{0x22, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x31, 0x73, 0x74, 0x61, 0x6b, 0x65},
		},
		{
			sdk.AccAddress("farmer2"),
			sdk.DefaultBondDenom,
			[]byte{0x22, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x32, 0x73, 0x74, 0x61, 0x6b, 0x65},
		},
		{
			sdk.AccAddress("farmer3"),
			sdk.DefaultBondDenom,
			[]byte{0x22, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x33, 0x73, 0x74, 0x61, 0x6b, 0x65},
		},
	}

	for _, tc := range testCases {
		key := types.GetStakingIndexKey(tc.farmerAcc, tc.stakingCoinDenom)
		s.Require().Equal(tc.expected, key)

		farmerAcc, stakingCoinDenom := types.ParseStakingIndexKey(key)
		s.Require().Equal(tc.farmerAcc, farmerAcc)
		s.Require().Equal(tc.stakingCoinDenom, stakingCoinDenom)
	}
}

func (s *keysTestSuite) TestGetStakingsByFarmerPrefix() {
	s.Require().Equal([]byte{0x22}, types.GetStakingsByFarmerPrefix(sdk.AccAddress("")))
	s.Require().Equal([]byte{0x22, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x31}, types.GetStakingsByFarmerPrefix(sdk.AccAddress("farmer1")))
	s.Require().Equal([]byte{0x22, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x32}, types.GetStakingsByFarmerPrefix(sdk.AccAddress("farmer2")))
	s.Require().Equal([]byte{0x22, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x33}, types.GetStakingsByFarmerPrefix(sdk.AccAddress("farmer3")))
}

func (s *keysTestSuite) TestGetQueuedStakingKey() {
	testCases := []struct {
		stakingCoinDenom string
		farmerAcc        sdk.AccAddress
		expected         []byte
	}{
		{
			sdk.DefaultBondDenom,
			nil,
			[]byte{0x23, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65},
		},
		{
			sdk.DefaultBondDenom,
			sdk.AccAddress("farmer1"),
			[]byte{0x23, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x31},
		},
		{
			sdk.DefaultBondDenom,
			sdk.AccAddress("farmer2"),
			[]byte{0x23, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x32},
		},
		{
			sdk.DefaultBondDenom,
			sdk.AccAddress("farmer3"),
			[]byte{0x23, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x33},
		},
	}

	for _, tc := range testCases {
		s.Require().Equal(tc.expected, types.GetQueuedStakingKey(tc.stakingCoinDenom, tc.farmerAcc))
	}
}

func (s *keysTestSuite) TestGetQueuedStakingIndexKey() {
	testCases := []struct {
		farmerAcc        sdk.AccAddress
		stakingCoinDenom string
		expected         []byte
	}{
		{
			sdk.AccAddress(""),
			sdk.DefaultBondDenom,
			[]byte{0x24, 0x73, 0x74, 0x61, 0x6b, 0x65},
		},
		{
			sdk.AccAddress("farmer1"),
			sdk.DefaultBondDenom,
			[]byte{0x24, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x31, 0x73, 0x74, 0x61, 0x6b, 0x65},
		},
		{
			sdk.AccAddress("farmer2"),
			sdk.DefaultBondDenom,
			[]byte{0x24, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x32, 0x73, 0x74, 0x61, 0x6b, 0x65},
		},
		{
			sdk.AccAddress("farmer3"),
			sdk.DefaultBondDenom,
			[]byte{0x24, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x33, 0x73, 0x74, 0x61, 0x6b, 0x65},
		},
	}

	for _, tc := range testCases {
		key := types.GetQueuedStakingIndexKey(tc.farmerAcc, tc.stakingCoinDenom)
		s.Require().Equal(tc.expected, key)
	}
}

func (s *keysTestSuite) TestGetQueuedStakingByFarmerPrefix() {
	s.Require().Equal([]byte{0x24}, types.GetQueuedStakingByFarmerPrefix(sdk.AccAddress("")))
	s.Require().Equal([]byte{0x24, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x31}, types.GetQueuedStakingByFarmerPrefix(sdk.AccAddress("farmer1")))
	s.Require().Equal([]byte{0x24, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x32}, types.GetQueuedStakingByFarmerPrefix(sdk.AccAddress("farmer2")))
	s.Require().Equal([]byte{0x24, 0x7, 0x66, 0x61, 0x72, 0x6d, 0x65, 0x72, 0x33}, types.GetQueuedStakingByFarmerPrefix(sdk.AccAddress("farmer3")))
}

func (s *keysTestSuite) TestGetTotalStakingKey() {
	s.Require().Equal([]byte{0x25}, types.GetTotalStakingKey(""))
	s.Require().Equal([]byte{0x25, 0x73, 0x74, 0x61, 0x6b, 0x65}, types.GetTotalStakingKey(sdk.DefaultBondDenom))
}

func (s *keysTestSuite) TestGetHistoricalRewardsKey() {
	testCases := []struct {
		stakingCoinDenom string
		epoch            uint64
		expected         []byte
	}{
		{
			"",
			0,
			[]byte{0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
		{
			sdk.DefaultBondDenom,
			1,
			[]byte{0x31, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
		},
		{
			sdk.DefaultBondDenom,
			2,
			[]byte{0x31, 0x5, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2},
		},
	}

	for _, tc := range testCases {
		s.Require().Equal(tc.expected, types.GetHistoricalRewardsKey(tc.stakingCoinDenom, tc.epoch))
	}
}

func (s *keysTestSuite) TestGetCurrentEpochKey() {
	s.Require().Equal([]byte{0x32}, types.GetCurrentEpochKey(""))
	s.Require().Equal([]byte{0x32, 0x73, 0x74, 0x61, 0x6b, 0x65}, types.GetCurrentEpochKey(sdk.DefaultBondDenom))
}
