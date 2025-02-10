package handler

import (
	"net/http"

	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) createTransaction(c *gin.Context) {
	var input models.TransactionInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	uuid, err := h.services.Transaction.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"uuid": uuid,
	})
}

type getAllTransactionsResponse struct {
	Transactions []models.Transaction `json:"data"`
}

func (h *Handler) getAllTransactions(c *gin.Context) {
	transactions, err := h.services.Transaction.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	c.JSON(http.StatusOK, getAllTransactionsResponse{
		Transactions: transactions,
	})
}

func (h *Handler) getTransactionById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	transaction, err := h.services.Transaction.GetById(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			newErrorResponse(c, http.StatusNotFound, "transaction not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	c.JSON(http.StatusOK, transaction)
}
