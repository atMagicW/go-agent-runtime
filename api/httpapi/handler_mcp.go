package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListMCPServersHandler 返回 MCP server 配置
func (h *Handler) ListMCPServersHandler(c *gin.Context) {
	if h.mcpService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "mcp service is not configured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"servers": h.mcpService.ListServers(),
	})
}
