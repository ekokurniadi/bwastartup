package campaign

type GetCampaignDetail struct {
	ID int `uri:"id" binding:"required"`
}
