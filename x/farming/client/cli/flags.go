package cli

// DONTCOVER

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagStakingCoinDenom = "staking-coin-denom"
	FlagFarmerAddr       = "farmer-addr"
)

func flagSetStaking() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagStakingCoinDenom, "", "The staking coin denom")
	fs.String(FlagFarmerAddr, "", "The bech32 address of the farmer account")

	return fs
}

func flagSetRewards() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagStakingCoinDenom, "", "The staking coin denom")
	fs.String(FlagFarmerAddr, "", "The bech32 address of the farmer account")

	return fs
}
