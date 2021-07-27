package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/tendermint/farming/x/farming/client/cli"
)

// ProposalHandler is the public plan creation handler.
var (
	AddProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitPublicPlanProposal, nil)
)
