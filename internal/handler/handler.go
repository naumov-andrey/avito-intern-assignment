package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/naumov-andrey/avito-intern-assignment/internal/service"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1")
	{
		accounts := api.Group("/account")
		{
			accounts.GET("/:userId/balance", h.GetBalance)
			accounts.PUT("/:userId/balance", h.UpdateBalance)
			accounts.GET("/:userId/history", h.GetHistory)
		}
		api.GET("/transfer", h.Transfer)
	}

	return router
}

func newErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"message": message})
}

func getUserId(c *gin.Context) (int, error) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "User id must be an integer")
		return 0, err
	}
	return userId, nil
}
