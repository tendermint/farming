package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/farming/x/farming/types"
)

// ParsePublicPlanProposal reads and parses a PublicPlanProposal from a file.
func ParsePublicPlanProposal(cdc codec.JSONCodec, proposalFile string) (types.AddPublicPlanProposal, error) {
	proposal := types.AddPublicPlanProposal{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
