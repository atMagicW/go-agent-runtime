package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSessionHandler 第一版先返回占位结果
func GetSessionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"session_id": c.Param("id"),
		"message":    "not implemented yet",
	})
}
