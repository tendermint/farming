package types

import (
	"fmt"
	time "time"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	proto "github.com/gogo/protobuf/proto"
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

// PackPlans converts PlanIs to Any slice.
func PackPlans(plans []PlanI) ([]*types.Any, error) {
	plansAny := make([]*types.Any, len(plans))
	for i, plan := range plans {
		msg, ok := plan.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("cannot proto marshal %T", plan)
		}
		any, err := types.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
		plansAny[i] = any
	}

	return plansAny, nil
}

// UnpackPlans converts Any slice to PlanIs.
func UnpackPlans(plansAny []*types.Any) ([]PlanI, error) {
	plans := make([]PlanI, len(plansAny))
	for i, any := range plansAny {
		p, ok := any.GetCachedValue().(PlanI)
		if !ok {
			return nil, fmt.Errorf("expected planI")
		}
		plans[i] = p
	}

	return plans, nil
}
