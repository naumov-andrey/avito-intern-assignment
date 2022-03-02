package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (h *Handler) GetBalance(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var balance decimal.Decimal
	currency, ok := c.GetQuery("currency")
	if ok {
		balance, err = h.accountService.GetConvertedBalance(userId, currency)
	} else {
		balance, err = h.accountService.GetBalance(userId)
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

	tx := c.MustGet("tx").(*gorm.DB)

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

	balance, err := h.accountService.WithTx(tx).UpdateBalance(userId, amount, input.Description)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Account{UserId: userId, Balance: balance})
}

func (h *Handler) GetHistory(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	limit := -1
	limitStr, ok := c.GetQuery("limit")
	if ok {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "Limit query parameter must be a positive integer")
			return
		}
	}

	cursor := -1
	cursorStr, ok := c.GetQuery("cursor")
	if ok {
		cursor, err = strconv.Atoi(cursorStr)
		if err != nil || cursor <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "Cursor query parameter must be a positive integer")
			return
		}
	}

	sortBy, ok := c.GetQuery("sort_by")
	if !ok {
		sortBy = "date"
	}

	orderBy, ok := c.GetQuery("order_by")
	if !ok {
		orderBy = "desc"
	}

	out, err := h.accountService.GetHistory(userId, limit, cursor, sortBy, orderBy)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *Handler) CreateTransfer(c *gin.Context) {
	tx := c.MustGet("tx").(*gorm.DB)

	var input model.TransferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Request body must contain transfer data")
		return
	}

	if input.CreditUserId == input.DebitUserId {
		newErrorResponse(c, http.StatusBadRequest, "Debit and credit users must be different")
		return
	}

	amount, err := decimal.NewFromString(input.Amount)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Amount data must be a decimal in format of string")
		return
	}

	if amount.LessThanOrEqual(decimal.Decimal{}) {
		newErrorResponse(c, http.StatusBadRequest, "Amount data must be a positive decimal")
		return
	}

	transfer := model.TransferData{
		DebitUserId:  input.DebitUserId,
		CreditUserId: input.CreditUserId,
		Amount:       amount,
		Description:  input.Description,
	}

	out, err := h.accountService.WithTx(tx).CreateTransfer(transfer)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
}
