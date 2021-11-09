package usecase

import (
	"funding/model"
	"funding/repository"
	"time"
)

type serviceCampaign struct {
	repositiry repository.CampaignRepository
}

type ServiceCampaign interface {
	Create(input model.CreateCampaignInput, user model.User) (*model.Campaign, error)
	GetAllCampaigns() ([]model.Campaign, error)
	GetCampaignID(id uint) (model.Campaign, error)
}

func NewServiceCampaign(repo repository.CampaignRepository) *serviceCampaign {
	return &serviceCampaign{repositiry: repo}
}

func (s *serviceCampaign) Create(input model.CreateCampaignInput, user model.User) (*model.Campaign, error) {

	var campaign model.Campaign
	campaign.UserID = user.ID
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	campaign.Allowed = model.Wait
	campaign.CurrentAmount = 0

	campaigns, err := s.repositiry.Create(campaign)
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil

}

func (s *serviceCampaign) GetAllCampaigns() ([]model.Campaign, error) {
	campaign, err := s.repositiry.FindAll()
	if err != nil {
		return []model.Campaign{}, err
	}

	return campaign, nil
}

func (s *serviceCampaign) GetCampaignID(id uint) (model.Campaign, error) {
	campaign, err := s.repositiry.FindByID(id)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
