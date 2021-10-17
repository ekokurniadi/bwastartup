package campaign

import (
	"bwastartup/datatables"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
	FindByID(ID int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Update(campaign Campaign) (Campaign, error)
	CreateImage(campaignImage CampaignImage) (CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)

	// repository for campaign
	GetCampaigns(campaigns datatables.DTJson) ([]CampaignOnWeb, error)
	GetTotalCampaigns(campaigns datatables.DTJson) (int, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary =1").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary= 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil

}

func (r *repository) FindByID(ID int) (Campaign, error) {
	var campaign Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", ID).Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
func (r *repository) Update(campaign Campaign) (Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
func (r *repository) CreateImage(campaignImage CampaignImage) (CampaignImage, error) {
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}
func (r *repository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	err := r.db.Model(&CampaignImage{}).Where("campaign_id= ?", campaignID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *repository) GetCampaigns(campaigns datatables.DTJson) ([]CampaignOnWeb, error) {
	var campaignOnWeb []CampaignOnWeb
	sql := "SELECT a.id,a.name,a.short_description,a.goal_amount,(select b.file_name from campaign_images b where b.campaign_id=a.id and b.is_primary=1) as campaign_images from campaigns a WHERE 1=1 "

	if search := campaigns.Search.Value; search != "" {
		sql = fmt.Sprintf("%s AND (a.name LIKE '%%%s%%' OR a.short_description LIKE '%%%s%%' OR a.goal_amount LIKE '%%%s%%') ", sql, search, search, search)
	}
	start := campaigns.Start
	length := campaigns.Length

	sql = fmt.Sprintf("%s LIMIT %d, %d", sql, start, length)
	err := r.db.Raw(sql).Scan(&campaignOnWeb).Error
	if err != nil {
		return campaignOnWeb, err
	}
	return campaignOnWeb, nil
}
func (r *repository) GetTotalCampaigns(campaigns datatables.DTJson) (int, error) {
	var campaignOnWeb []CampaignOnWeb
	sql := "SELECT a.id,a.name,a.short_description,a.goal_amount,(select b.file_name from campaign_images b where b.campaign_id=a.id and b.is_primary=1) as campaign_images from campaigns a WHERE 1=1 "

	if search := campaigns.Search.Value; search != "" {
		sql = fmt.Sprintf("%s AND (a.name LIKE '%%%s%%' OR a.short_description LIKE '%%%s%%' OR a.goal_amount LIKE '%%%s%%') ", sql, search, search, search)
	}

	err := r.db.Raw(sql).Scan(&campaignOnWeb).Error
	if err != nil {
		return len(campaignOnWeb), err
	}
	return len(campaignOnWeb), nil
}
