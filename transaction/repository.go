package transaction

import (
	"bwastartup/datatables"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
	GetByUserID(UserID int) ([]Transaction, error)
	GetByID(ID int) (Transaction, error)
	Save(transaction Transaction) (Transaction, error)
	Update(transaction Transaction) (Transaction, error)

	// repository for campaign
	GetTransactions(input datatables.DTJson) ([]TransactionOnWeb, error)
	GetTotalTransactions(input datatables.DTJson) (int, error)
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (r *repository) GetByUserID(UserID int) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary=1").Where("user_id = ?", UserID).Order("id desc").Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil

}

func (r *repository) Save(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Update(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) GetByID(ID int) (Transaction, error) {
	var transaction Transaction
	err := r.db.Where("id = ?", ID).Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) GetTransactions(input datatables.DTJson) ([]TransactionOnWeb, error) {
	var transactions []TransactionOnWeb
	sql := "SELECT a.id,b.name,a.status,a.amount from transactions a join campaigns b on a.campaign_id=b.id WHERE 1=1 "

	if search := input.Search.Value; search != "" {
		sql = fmt.Sprintf("%s AND (b.name LIKE '%%%s%%' OR a.status LIKE '%%%s%%' OR a.amount LIKE '%%%s%%') ", sql, search, search, search)
	}
	start := input.Start
	length := input.Length

	sql = fmt.Sprintf("%s LIMIT %d, %d", sql, start, length)
	err := r.db.Raw(sql).Scan(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
func (r *repository) GetTotalTransactions(input datatables.DTJson) (int, error) {
	var transactions []TransactionOnWeb
	sql := "SELECT a.id,b.name,a.status,a.amount from transactions a join campaigns b on a.campaign_id=b.id WHERE 1=1 "

	if search := input.Search.Value; search != "" {
		sql = fmt.Sprintf("%s AND (b.name LIKE '%%%s%%' OR a.status LIKE '%%%s%%' OR a.amount LIKE '%%%s%%') ", sql, search, search, search)
	}

	err := r.db.Raw(sql).Scan(&transactions).Error
	if err != nil {
		return len(transactions), err
	}
	return len(transactions), nil
}
