package campaign

import (
	"bwastartup/user"
	"time"

	"github.com/leekchan/accounting"
)

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BakerCount       int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
	User             user.User
	AmountOnText     string
}

type CampaignImage struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CampaignOnWeb struct {
	ID               int
	Name             string
	ShortDescription string
	GoalAmount       int
	CampaignImages   string
}

func (c CampaignOnWeb) GoalAmountFormatIDR() string {
	format := accounting.Accounting{Symbol: "Rp.", Precision: 2, Thousand: ".", Decimal: ","}
	return format.FormatMoney(c.GoalAmount)
}
