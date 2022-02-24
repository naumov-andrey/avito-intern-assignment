package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
)

func (h *Handler) GetBalance(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "user id must be an integer")
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

	c.JSON(200, gin.H{"balance": balance})
}
