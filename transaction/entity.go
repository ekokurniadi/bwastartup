package transaction

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       user.User
	Campaign   campaign.Campaign
	PaymentURL string
}

type TransactionOnWeb struct {
	ID     int
	Name   string
	Amount int
	Status string
}
