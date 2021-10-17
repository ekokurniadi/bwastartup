package transaction

import (
	"bwastartup/campaign"
	"bwastartup/datatables"
	"bwastartup/payment"
	"errors"
	"strconv"
)

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransaction) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error

	// web service
	GetAllTransaction(input datatables.DTJson) ([]TransactionOnWeb, error)
}
type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransaction) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("not an owned of the campaign")
	}
	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)

	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	transaction.Code = ""

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transactions{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)

	if err != nil {
		return newTransaction, err
	}
	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)

	if err != nil {
		return newTransaction, err
	}
	return newTransaction, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)
	transaction, err := s.repository.GetByID(transaction_id)

	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "canceled"
	}

	updatedTransaction, err := s.repository.Update(transaction)

	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BakerCount = campaign.BakerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}
	return nil

}

func (s *service) GetAllTransaction(input datatables.DTJson) ([]TransactionOnWeb, error) {
	inputTransaction := datatables.DTJson{}
	inputTransaction.Draw = input.Draw
	inputTransaction.Columms = input.Columms
	inputTransaction.Length = input.Length
	inputTransaction.Orders = input.Orders
	inputTransaction.Start = input.Start
	inputTransaction.Search.Value = input.Search.Value

	data, err := s.repository.GetTransactions(inputTransaction)
	if err != nil {
		return data, err
	}
	return data, nil
}
