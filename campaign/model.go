package campaign

import (
	"time"
)

type Campaign struct {
	ID               int             `db:"id" json:"id"`
	UserID           int             `db:"user_id" json:"user_id"`
	Name             string          `db:"name" json:"campaign_name"`
	ShortDescription string          `db:"short_desc" json:"short_description"`
	Description      string          `db:"description" json:"description"`
	GoalAmount       uint32          `db:"goal_amount" json:"goal_amount"`
	CurrentAmount    uint32          `db:"current_amount" json:"current_amount"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
	CampaignImages   []CampaignImage `db:"campaign_image" json:"campaign_image"`
}

type CampaignImage struct {
	ID         int    `db:"id" json:"images_id,omitempty"`
	CampaignID int    `db:"campaign_id" json:"campaign_id"`
	FileName   string `db:"filename" json:"filename"`
	IsPrimary  int    `db:"is_primary" json:"is_primary"`
}

type CreateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       uint32 `json:"goal_amount" binding:"required"`
}
