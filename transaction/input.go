package transaction

import "bwastartup/user"

type GetCampaignTransaction struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
