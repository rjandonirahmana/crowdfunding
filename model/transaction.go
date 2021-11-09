package model

type transaction struct {
	ID         int `db:"id" json:"transaction_id,omitempty"`
	CampaignID int `db:"campaign_id" json:",omitempty"`
}
