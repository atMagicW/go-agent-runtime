package httpapi

import (
	"net/http"
	"strings"

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

	kind := strings.TrimSpace(strings.ToLower(c.Query("kind")))
	source := strings.TrimSpace(strings.ToLower(c.Query("source")))

	capabilities := h.capabilityService.ListCapabilities()
	if kind != "" || source != "" {
		filtered := make([]CapabilityView, 0, len(capabilities))
		for _, item := range capabilities {
			if kind != "" && strings.ToLower(item.Kind) != kind {
				continue
			}
			if source != "" && strings.ToLower(item.Source) != source {
				continue
			}
			filtered = append(filtered, item)
		}
		capabilities = filtered
	}

	c.JSON(http.StatusOK, gin.H{
		"capabilities": capabilities,
	})
}
