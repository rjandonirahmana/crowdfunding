package campaign

import (
	"funding/account"
	"time"
)

type Service struct {
	repositiry CampaignRepository
}

type ServiceCampaign interface {
	Create(input CreateCampaignInput, user account.User) error
	GetAllCampaigns() ([]Campaign, error)
}

func NewServiceCampaign(repo CampaignRepository) *Service {
	return &Service{repositiry: repo}
}

func (s *Service) Create(input CreateCampaignInput, user account.User) error {

	var campaign Campaign
	campaign.UserID = user.ID
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	campaign.CurrentAmount = 0

	err := s.repositiry.Create(campaign)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetAllCampaigns() ([]Campaign, error) {
	campaign, err := s.repositiry.FindAll()
	if err != nil {
		return []Campaign{}, err
	}

	return campaign, nil
}
