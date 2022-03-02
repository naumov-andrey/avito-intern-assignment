package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func transaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		c.Set("tx", tx)
		c.Next()

		if c.Writer.Status() == http.StatusOK {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
}

func getUserId(c *gin.Context) (int, error) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "User id must be an integer")
		return 0, err
	}
	return userId, nil
}
