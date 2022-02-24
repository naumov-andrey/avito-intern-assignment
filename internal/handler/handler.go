package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/naumov-andrey/avito-intern-assignment/internal/service"
)

type Handler struct {
	accountService *service.AccountService
}

func NewHandler(accountService *service.AccountService) *Handler {
	return &Handler{accountService}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1")
	{
		balances := api.Group("/balance")
		{
			balances.GET("/:userId", h.GetBalance)
		}
	}

	return router
}

func newErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}
