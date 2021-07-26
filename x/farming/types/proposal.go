package types

import (
	"fmt"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypePublicPlan string = "PublicPlan"
)

// Implements Proposal Interface
var _ gov.Content = &PublicPlanProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypePublicPlan)
	gov.RegisterProposalTypeCodec(&PublicPlanProposal{}, "cosmos-sdk/PublicPlanProposal")
}

func NewPublicPlanProposal(title, description, name string, addReq []*AddRequestProposal,
	updateReq []*UpdateRequestProposal, deleteReq []*DeleteRequestProposal) (gov.Content, error) {
	return &PublicPlanProposal{
		Title:                  title,
		Description:            description,
		Name:                   name,
		AddRequestProposals:    addReq,
		UpdateRequestProposals: updateReq,
		DeleteRequestProposals: deleteReq,
	}, nil
}

func (p *PublicPlanProposal) GetTitle() string { return p.Title }

func (p *PublicPlanProposal) GetDescription() string { return p.Description }

func (p *PublicPlanProposal) ProposalRoute() string { return RouterKey }

func (p *PublicPlanProposal) ProposalType() string { return ProposalTypePublicPlan }

func (p *PublicPlanProposal) ValidateBasic() error {
	// if p.AddRequestProposals == nil && p.UpdateRequestProposals == nil && p.DeleteRequestProposals == nil {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "proposal must not be empty")
	// }
	return gov.ValidateAbstract(p)
}

func (p PublicPlanProposal) String() string {
	return fmt.Sprintf(`Public Plan Proposal:
  Title:       			  %s
  Description: 		      %s
  AddRequestProposals: 	  %s
  UpdateRequestProposals: %s
  DeleteRequestProposals: %s
`, p.Title, p.Description, p.AddRequestProposals, p.UpdateRequestProposals, p.DeleteRequestProposals)
}
