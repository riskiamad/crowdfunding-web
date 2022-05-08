package handler

import (
	"crowdfunding-web/helper"
	"crowdfunding-web/transaction"
	"crowdfunding-web/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// catch parameter from URI
// mapping parameter to input struct
// call service, input struct as parameter
// service, get campaign ID
// repo search data which match with the campaign ID

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	if err := c.ShouldBindUri(&input); err != nil {
		response := helper.APIResponse("Failed to get campaigns's transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaigns's transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign's Transaction", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}
