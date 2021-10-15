package campaign

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

type CampaignRepository interface {
	FindAll() ([]Campaign, error)
	Create(campaign Campaign) error
	FindByID(campaign_id int) (Campaigns, error)
}

func (r *repository) FindAll() ([]Campaign, error) {
	querry := `SELECT c.*, ci.filename as "campaign_image.filename" FROM campaign c LEFT JOIN campaign_image ci ON c.id = ci.campaign_id`

	var campaigns []Campaign
	err := r.db.Select(&campaigns, querry)
	if err != sql.ErrNoRows {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) Create(campaign Campaign) error {
	querry := `INSERT INTO campaign (id, user_id, name, short_desc, description, goal_amount, current_amount, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(querry, campaign.ID, campaign.UserID, campaign.Name, campaign.ShortDescription, campaign.Description, campaign.GoalAmount, campaign.CurrentAmount, campaign.CreatedAt, campaign.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByID(campaign_id int) (Campaigns, error) {
	querry := `SELECT campaign_id, filename, is_primary FROM campaign_image WHERE campaign_id = $1`

	var images []CampaignImage
	err := r.db.Select(&images, querry, campaign_id)

	if err != nil {
		return Campaigns{}, err
	}

	var campaign Campaigns
	querry1 := `SELECT * FROM campaign WHERE id = $1`
	err = r.db.Get(&campaign, querry1, campaign_id)
	if err != sql.ErrNoRows {
		return Campaigns{}, err
	}

	campaign.CampaignImage = append(campaign.CampaignImage, images...)
	return campaign, nil

}
