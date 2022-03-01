package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
)

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

	out, err := h.service.Transaction.GetHistory(userId, limit, cursor, sortBy, orderBy)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *Handler) Transfer(c *gin.Context) {
	var input model.TransferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Request body must contain transfer data")
		return
	}

	amount, err := decimal.NewFromString(input.Amount)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Amount data must be a decimal in format of string")
		return
	}

	transfer := model.TransferData{
		DebitUserId:  input.DebitUserId,
		CreditUserId: input.CreditUserId,
		Amount:       amount,
		Description:  input.Description,
	}

	out, err := h.service.Transaction.Transfer(transfer)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
}
