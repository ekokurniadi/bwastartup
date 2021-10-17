package campaign

import "fmt"

type CampaignWebFormatter struct {
	Number           int    `json:"no"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	GoalAmount       int    `json:"goal_amount"`
	Link             string `json:"link"`
	CampaignImages   string `json:"campaign_images"`
}

func FormatWebCampaign(campaign CampaignOnWeb) CampaignWebFormatter {
	campaignFormatter := CampaignWebFormatter{}
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CampaignImages = ""

	if (campaign.CampaignImages) != "" {
		campaignFormatter.CampaignImages = fmt.Sprintf("<img src='%s' width='90' class='img-fluid img-circle'>", campaign.CampaignImages)
	}

	campaignFormatter.Link = fmt.Sprintf("<a href='/campaigns/view/%d' class='btn btn-light'><i class='fa fa-camera'></i></a> <a href='/campaigns/edit/%d' class='btn btn-warning'><i class='fa fa-pencil'></i></a> <a href='/campaigns/delete/%d' class='btn btn-danger' onclick='javascript: return confirm(\"Are You Sure ?\")''><i class='fa fa-trash'></i></a>", campaign.ID, campaign.ID, campaign.ID)

	return campaignFormatter
}

func FormatWebCampaigns(campaigns []CampaignOnWeb) []CampaignWebFormatter {

	campaignsFormatters := []CampaignWebFormatter{}
	for _, campaign := range campaigns {
		campaignsFormatter := FormatWebCampaign(campaign)
		campaignsFormatters = append(campaignsFormatters, campaignsFormatter)
	}

	return campaignsFormatters
}
