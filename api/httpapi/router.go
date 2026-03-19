package httpapi

import "github.com/gin-gonic/gin"

// Handler 统一承载 HTTP 层依赖
type Handler struct {
	agentService      AgentService
	sessionService    SessionService
	capabilityService CapabilityService
	ingestService     IngestService
}

// NewHandler 创建 HTTP Handler
func NewHandler(
	agentService AgentService,
	sessionService SessionService,
	capabilityService CapabilityService,
	ingestService IngestService,
) *Handler {
	return &Handler{
		agentService:      agentService,
		sessionService:    sessionService,
		capabilityService: capabilityService,
		ingestService:     ingestService,
	}
}

// RegisterRoutes 注册 HTTP 路由
func RegisterRoutes(r *gin.Engine, h *Handler) {
	v1 := r.Group("/v1")
	{
		v1.GET("/health", h.HealthHandler)
		v1.GET("/capabilities", h.ListCapabilitiesHandler)
		v1.POST("/chat", h.ChatHandler)
		v1.GET("/sessions/:id", h.GetSessionHandler)

		// RAG ingest
		v1.POST("/rag/ingest", h.IngestTextHandler)
	}
}
