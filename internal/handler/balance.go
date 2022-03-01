package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/shopspring/decimal"
	"net/http"
)

func (h *Handler) GetBalance(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var balance decimal.Decimal
	currency, ok := c.GetQuery("currency")
	if ok {
		balance, err = h.service.Balance.GetConvertedBalance(userId, currency)
	} else {
		balance, err = h.service.Balance.GetBalance(userId)
	}

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Account{UserId: userId, Balance: balance})
}

func (h *Handler) UpdateBalance(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input model.UpdateBalanceInput
	if err = c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Request body must contain amount data")
		return
	}

	amount, err := decimal.NewFromString(input.Amount)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Amount data must be a decimal in format of string")
		return
	}

	balance, err := h.service.Balance.UpdateBalance(userId, amount, input.Description)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Account{UserId: userId, Balance: balance})
}
