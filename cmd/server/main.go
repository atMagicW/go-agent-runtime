package main

import (
	"context"
	"os"
	"time"

	mockembedding "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/mock_embedding"
	pgrag "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/pgvector"
	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/atMagicW/go-agent-runtime/api/httpapi"
	capregistry "github.com/atMagicW/go-agent-runtime/internal/adapters/capability"
	mcpcap "github.com/atMagicW/go-agent-runtime/internal/adapters/capability/mcp"
	"github.com/atMagicW/go-agent-runtime/internal/adapters/capability/skills"
	"github.com/atMagicW/go-agent-runtime/internal/adapters/capability/tools"
	openaiadapter "github.com/atMagicW/go-agent-runtime/internal/adapters/llm/openai"
	pgrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/postgres"
	"github.com/atMagicW/go-agent-runtime/internal/app"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	agentgov "github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
	agentrouter "github.com/atMagicW/go-agent-runtime/internal/usecase/router"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	pgDSN := os.Getenv("POSTGRES_DSN")
	if pgDSN == "" {
		logger.Fatal("POSTGRES_DSN is not set")
	}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		logger.Fatal("OPENAI_API_KEY is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgrepo.NewDB(ctx, pgDSN)
	if err != nil {
		logger.Fatal("init db failed", zap.Error(err))
	}
	defer db.Close()

	sessionRepo := pgrepo.NewSessionRepository(db)
	sessionService := app.NewSessionService(sessionRepo)

	// 统一治理组件
	breakers := agentgov.NewBreakerRegistry()
	fallbacks := agentgov.NewDefaultFallbackPolicy()

	// 初始化 LLM clients
	openAIClient := openaiadapter.NewClient(openAIKey)
	llmClients := map[string]ports.LLMClient{
		"openai": openAIClient,
	}
	modelRouter := agentrouter.NewModelRouter(llmClients, breakers, fallbacks)

	// 初始化统一能力注册表
	registry := capregistry.NewRegistry()

	// 注册本地 Skill / Tool
	registry.MustRegister(skills.NewResumeAnalyzerSkill())
	registry.MustRegister(tools.NewKeywordExtractTool())

	// 注册 MCP Tool
	mcpClient := mcpcap.NewClient()
	for _, spec := range mcpcap.DefaultToolSpecs() {
		registry.MustRegister(mcpcap.NewToolCapability(mcpClient, spec))
	}

	// 初始化 RAG
	ragRepo := pgrag.NewRepository(db)
	embeddingProvider := mockembedding.NewProvider(1536)
	ragService := app.NewRAGService(ragRepo, embeddingProvider)

	// 确保演示知识库存在
	if err := ragRepo.EnsureKnowledgeBase(ctx, rag.KnowledgeBase{
		KBID:        "default",
		TenantID:    "default",
		Name:        "Default Knowledge Base",
		Description: "默认演示知识库",
		Enabled:     true,
	}); err != nil {
		logger.Fatal("init default kb failed", zap.Error(err))
	}

	if err := ragRepo.EnsureKnowledgeBase(ctx, rag.KnowledgeBase{
		KBID:        "knowledge_a",
		TenantID:    "default",
		Name:        "Knowledge A",
		Description: "演示知识库 A",
		Enabled:     true,
	}); err != nil {
		logger.Fatal("init knowledge_a failed", zap.Error(err))
	}

	if err := ragRepo.EnsureKnowledgeBase(ctx, rag.KnowledgeBase{
		KBID:        "knowledge_b",
		TenantID:    "default",
		Name:        "Knowledge B",
		Description: "演示知识库 B",
		Enabled:     true,
	}); err != nil {
		logger.Fatal("init knowledge_b failed", zap.Error(err))
	}

	agentService := app.NewAgentService(
		sessionService,
		modelRouter,
		registry,
		ragService,
		breakers,
		fallbacks,
	)

	capabilityService := app.NewCapabilityService(registry)

	handler := httpapi.NewHandler(
		agentService,
		sessionService,
		capabilityService,
	)

	router := gin.Default()
	httpapi.RegisterRoutes(router, handler)

	logger.Info("Agent Runtime Server starting at :8080")

	// 测试写入数据
	if err := pgrag.SeedDemoKnowledgeBase(ctx, ragRepo, embeddingProvider, "default"); err != nil {
		logger.Fatal("seed default kb failed", zap.Error(err))
	}

	if err := pgrag.SeedDemoKnowledgeBase(ctx, ragRepo, embeddingProvider, "knowledge_a"); err != nil {
		logger.Fatal("seed knowledge_a failed", zap.Error(err))
	}

	if err := pgrag.SeedDemoKnowledgeBase(ctx, ragRepo, embeddingProvider, "knowledge_b"); err != nil {
		logger.Fatal("seed knowledge_b failed", zap.Error(err))
	}

	if err := router.Run(":8080"); err != nil {
		logger.Fatal("server start failed", zap.Error(err))
	}
}
