package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/tendermint/farming/x/farming/client/cli"
	"github.com/tendermint/farming/x/farming/client/rest"
)

// ProposalHandler is the public plan creation handler.
var (
	AddProposalHandler    = govclient.NewProposalHandler(cli.GetCmdSubmitAddPublicPlanProposal, rest.ProposalRESTHandler)
	UpdateProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitUpdatePublicPlanProposal, rest.ProposalRESTHandler)
	DeleteProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitDeletePublicPlanProposal, rest.ProposalRESTHandler)
)
