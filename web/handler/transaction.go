package handler

import (
	"bwastartup/transaction"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler(transactionService transaction.Service) *transactionHandler {
	return &transactionHandler{transactionService}
}

func (h *transactionHandler) Index(c *gin.Context) {
	session := sessions.Default(c)
	data := session.Get("message")
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "transaction_index.html", gin.H{"data": data})
}
