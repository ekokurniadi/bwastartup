package handler

import (
	"bwastartup/campaign"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	session := sessions.Default(c)
	data := session.Get("message")
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"data": data})
}
