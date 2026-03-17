package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查接口
func (h *Handler) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
