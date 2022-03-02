package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/naumov-andrey/avito-intern-assignment/internal/service"
	"gorm.io/gorm"
)

type Handler struct {
	accountService *service.AccountService
}

func NewHandler(accountService *service.AccountService) *Handler {
	return &Handler{accountService}
}

func (h *Handler) InitRoutes(db *gorm.DB) *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1")
	{
		accounts := api.Group("/account")
		{
			accounts.GET("/:userId/balance", h.GetBalance)
			accounts.PUT("/:userId/balance", transaction(db), h.UpdateBalance)
			accounts.GET("/:userId/history", h.GetHistory)
			accounts.POST("/transfer", transaction(db), h.CreateTransfer)
		}
	}

	return router
}

func newErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"message": message})
}
