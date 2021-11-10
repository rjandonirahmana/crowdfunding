package repository

import (
	"fmt"

	"funding/model"

	"github.com/jmoiron/sqlx"
)

type campaignRepo struct {
	db *sqlx.DB
}

func NewRepositoryCampaign(db *sqlx.DB) *campaignRepo {
	return &campaignRepo{db: db}
}

type CampaignRepository interface {
	FindAll() ([]*model.Campaign, error)
	Create(campaign *model.Campaign) (*model.Campaign, error)
	FindByID(campaign_id uint) (*model.Campaign, error)
}

// SELECT c.*, ci.filename as "campaign_image.filename" FROM campaign c LEFT JOIN campaign_image ci ON c.id = ci.campaign_id LIMIT 10

func (r *campaignRepo) FindAll() ([]*model.Campaign, error) {
	querry := `SELECT * FROM campaign`

	var campaigns []*model.Campaign
	err := r.db.Select(&campaigns, querry)

	if err != nil {
		return nil, err
	}

	return campaigns, nil
}

func (r *campaignRepo) Create(campaign *model.Campaign) (*model.Campaign, error) {
	querry := `INSERT INTO campaign (user_id, name, short_desc, description, goal_amount, current_amount, created_at, updated_at, allowed) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	err := r.db.QueryRowx(querry, campaign.UserID, campaign.Name, campaign.ShortDescription, campaign.Description, campaign.GoalAmount, campaign.CurrentAmount, campaign.CreatedAt, campaign.UpdatedAt, campaign.Allowed).Scan(&campaign.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *campaignRepo) FindByID(campaign_id uint) (*model.Campaign, error) {

	var campaign model.Campaign
	querry1 := `SELECT * FROM campaign WHERE id = $1`
	err := r.db.Get(&campaign, querry1, campaign_id)
	if err != nil {
		return &campaign, err
	}

	fmt.Println(campaign.Name)
	querry := `SELECT campaign_id, filename, is_primary FROM campaign_image WHERE campaign_id = $1`

	var images []model.CampaignImage
	err = r.db.Select(&images, querry, campaign_id)

	if err != nil {
		return &campaign, err
	}
	campaign.CampaignImage = append(campaign.CampaignImage, images...)

	return &campaign, nil

}
