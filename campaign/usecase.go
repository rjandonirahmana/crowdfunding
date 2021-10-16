package campaign

import (
	"funding/account"
	"time"
)

type Service struct {
	repositiry CampaignRepository
}

type ServiceCampaign interface {
	Create(input CreateCampaignInput, user account.User) (Campaigns, error)
	GetAllCampaigns() ([]Campaign, error)
	GetCampaignID(id uint) (Campaigns, error)
}

func NewServiceCampaign(repo CampaignRepository) *Service {
	return &Service{repositiry: repo}
}

func (s *Service) Create(input CreateCampaignInput, user account.User) (Campaigns, error) {

	var campaign Campaign
	campaign.UserID = user.ID
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	campaign.CurrentAmount = 0

	id, err := s.repositiry.Create(campaign)
	if err != nil {
		return Campaigns{}, err
	}
	campaigns, err := s.repositiry.FindByID(id)
	if err != nil {
		return Campaigns{}, err
	}

	return campaigns, nil

}

func (s *Service) GetAllCampaigns() ([]Campaign, error) {
	campaign, err := s.repositiry.FindAll()
	if err != nil {
		return []Campaign{}, err
	}

	return campaign, nil
}

func (s *Service) GetCampaignID(id uint) (Campaigns, error) {
	campaign, err := s.repositiry.FindByID(id)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
