package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetLatestPromptHandler 获取某个 prompt 的最新版本
func (h *Handler) GetLatestPromptHandler(c *gin.Context) {
	if h.promptService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "prompt service is not configured",
		})
		return
	}

	name := c.Param("name")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item, err := h.promptService.GetLatest(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// ListPromptVersionsHandler 列出某个 prompt 的全部版本
func (h *Handler) ListPromptVersionsHandler(c *gin.Context) {
	if h.promptService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "prompt service is not configured",
		})
		return
	}

	name := c.Param("name")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := h.promptService.ListVersions(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prompt_name": name,
		"versions":    items,
	})
}
