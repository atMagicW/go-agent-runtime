package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// ChatRequest 用户请求结构
type ChatRequest struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	Model     string `json:"model"`
	Stream    bool   `json:"stream"`
}

// ChatHandler 是 Agent 主入口
func (h *Handler) ChatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	requestID := uuid.New().String()

	reqCtx := agent.RequestContext{
		RequestID: requestID,
		SessionID: req.SessionID,
		UserID:    req.UserID,
		Model:     req.Model,
		Stream:    req.Stream,
	}

	if req.Stream {
		StreamResponse(c, func(writer *StreamWriter) {
			h.agentService.RunStream(reqCtx, req.Message, writer)
		})
		return
	}

	resp, err := h.agentService.Run(reqCtx, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
