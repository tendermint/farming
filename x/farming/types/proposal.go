package types

import (
	"fmt"
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeAddPublicPlan    string = "AddPublicPlan"
	ProposalTypeUpdatePublicPlan string = "UpdatePublicPlan"
	ProposalTypeDeletePublicPlan string = "DeletePublicPlan"
)

// Implements Proposal Interface
var _ gov.Content = &AddPublicPlanProposal{}
var _ gov.Content = &UpdatePublicPlanProposal{}
var _ gov.Content = &DeletePublicPlanProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeAddPublicPlan)
	gov.RegisterProposalTypeCodec(&AddPublicPlanProposal{}, "cosmos-sdk/AddPublicPlanProposal")
	gov.RegisterProposalType(ProposalTypeUpdatePublicPlan)
	gov.RegisterProposalTypeCodec(&UpdatePublicPlanProposal{}, "cosmos-sdk/UpdatePublicPlanProposal")
	gov.RegisterProposalType(ProposalTypeDeletePublicPlan)
	gov.RegisterProposalTypeCodec(&DeletePublicPlanProposal{}, "cosmos-sdk/DeletePublicPlanProposal")
}

func NewAddPublicPlanProposal(title, description string, coinWeights sdk.DecCoins, startTime, endTime time.Time) (gov.Content, error) {
	return &AddPublicPlanProposal{
		Title:              title,
		Description:        description,
		StakingCoinWeights: coinWeights,
		StartTime:          startTime,
		EndTime:            endTime,
	}, nil
}

func (p *AddPublicPlanProposal) GetTitle() string { return p.Title }

func (p *AddPublicPlanProposal) GetDescription() string { return p.Description }

func (p *AddPublicPlanProposal) ProposalRoute() string { return RouterKey }

func (p *AddPublicPlanProposal) ProposalType() string { return ProposalTypeAddPublicPlan }

func (p *AddPublicPlanProposal) ValidateBasic() error {
	// TODO: not implemented yet
	return gov.ValidateAbstract(p)
}

func (p AddPublicPlanProposal) String() string {
	return fmt.Sprintf(`Add Public Plan Proposal:
  Title:       		  %s
  Description: 		  %s
  StakingCoinWeights: %s
  StartTime: 	      %s
  EndTime: 	   	  	  %s
`, p.Title, p.Description, p.StakingCoinWeights, p.StartTime, p.EndTime)
}

func NewUpdatePublicPlanProposal(title, description string, id uint64, coinWeights sdk.DecCoins, startTime, endTime time.Time) (gov.Content, error) {
	return &UpdatePublicPlanProposal{
		Title:              title,
		Description:        description,
		PlanId:             id,
		StakingCoinWeights: coinWeights,
		StartTime:          startTime,
		EndTime:            endTime,
	}, nil
}

func (p *UpdatePublicPlanProposal) GetTitle() string { return p.Title }

func (p *UpdatePublicPlanProposal) GetDescription() string { return p.Description }

func (p *UpdatePublicPlanProposal) ProposalRoute() string { return RouterKey }

func (p *UpdatePublicPlanProposal) ProposalType() string { return ProposalTypeUpdatePublicPlan }

func (p *UpdatePublicPlanProposal) ValidateBasic() error {
	// TODO: not implemented yet
	return gov.ValidateAbstract(p)
}

func (p UpdatePublicPlanProposal) String() string {
	return fmt.Sprintf(`Update Public Plan Proposal:
  Title:       		  %s
  Description: 		  %s
  PlanId: 		  	  %s
  StakingCoinWeights: %s
  StartTime: 	      %s
  EndTime: 	   	  	  %s
`, p.Title, p.Description, p.PlanId, p.StakingCoinWeights, p.StartTime, p.EndTime)
}

func NewDeletePublicPlanProposal(title, description string, id uint64) (gov.Content, error) {
	return &DeletePublicPlanProposal{
		Title:       title,
		Description: description,
		PlanId:      id,
	}, nil
}

func (p *DeletePublicPlanProposal) GetTitle() string { return p.Title }

func (p *DeletePublicPlanProposal) GetDescription() string { return p.Description }

func (p *DeletePublicPlanProposal) ProposalRoute() string { return RouterKey }

func (p *DeletePublicPlanProposal) ProposalType() string { return ProposalTypeUpdatePublicPlan }

func (p *DeletePublicPlanProposal) ValidateBasic() error {
	// TODO: not implemented yet
	return gov.ValidateAbstract(p)
}

func (p DeletePublicPlanProposal) String() string {
	return fmt.Sprintf(`Delete Public Plan Proposal:
  Title:       		  %s
  Description: 		  %s
  PlanId: 		  	  %s
`, p.Title, p.Description, p.PlanId)
}
