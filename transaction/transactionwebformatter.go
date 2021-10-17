package transaction

import "fmt"

type TransactionWebFormatter struct {
	Number int    `json:"no"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
	Status string `json:"status"`
	Link   string `json:"link"`
}

func FormatWebTransaction(transaction TransactionOnWeb) TransactionWebFormatter {
	transactionFormatter := TransactionWebFormatter{}
	transactionFormatter.Name = transaction.Name

	span := ""
	if (transaction.Status) == "paid" {
		span = "<span class='badge badge-success' style='padding:5px 5px 8px 5px !important'>" + transaction.Status + "</span>"
	} else {
		span = "<span class='badge badge-danger' style='padding:5px 5px 8px 5px !important'>" + transaction.Status + "</span>"
	}
	transactionFormatter.Status = span
	transactionFormatter.Amount = transaction.Amount

	transactionFormatter.Link = fmt.Sprintf("<a href='/transactions/view/%d' class='btn btn-light'><i class='fa fa-camera'></i></a> <a href='/transactions/edit/%d' class='btn btn-warning'><i class='fa fa-pencil'></i></a> <a href='/transactions/delete/%d' class='btn btn-danger' onclick='javascript: return confirm(\"Are You Sure ?\")''><i class='fa fa-trash'></i></a>", transaction.ID, transaction.ID, transaction.ID)

	return transactionFormatter
}

func FormatWebTransactions(transactions []TransactionOnWeb) []TransactionWebFormatter {

	transactionsFormatters := []TransactionWebFormatter{}
	for _, transaction := range transactions {
		transactionFormatter := FormatWebTransaction(transaction)
		transactionsFormatters = append(transactionsFormatters, transactionFormatter)
	}

	return transactionsFormatters
}
