package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
)

// IngestTextHandler 文本入库接口
func (h *Handler) IngestTextHandler(c *gin.Context) {
	if h.ingestService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ingest service is not configured",
		})
		return
	}

	var req rag.IngestTextRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := h.ingestService.IngestText(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
