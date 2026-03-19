package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListCapabilitiesHandler 返回已注册能力列表
func (h *Handler) ListCapabilitiesHandler(c *gin.Context) {
	if h.capabilityService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "capability service is not configured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"capabilities": h.capabilityService.ListCapabilities(),
	})
}
