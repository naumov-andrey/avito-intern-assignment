package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) GetBalance(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "user id must be an integer")
		return
	}

	balance, err := h.accountService.GetBalance(userId)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(200, gin.H{"balance": balance})
}
