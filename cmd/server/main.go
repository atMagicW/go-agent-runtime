package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/atMagicW/go-agent-runtime/api/httpapi"
	capregistry "github.com/atMagicW/go-agent-runtime/internal/adapters/capability"
	mcpcap "github.com/atMagicW/go-agent-runtime/internal/adapters/capability/mcp"
	"github.com/atMagicW/go-agent-runtime/internal/adapters/capability/skills"
	"github.com/atMagicW/go-agent-runtime/internal/adapters/capability/tools"
	openaiadapter "github.com/atMagicW/go-agent-runtime/internal/adapters/llm/openai"
	pgrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/postgres"
	mockembedding "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/mock_embedding"
	pgrag "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/pgvector"
	"github.com/atMagicW/go-agent-runtime/internal/app"
	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
	"github.com/atMagicW/go-agent-runtime/internal/pkg/config"
	"github.com/atMagicW/go-agent-runtime/internal/pkg/textsplitter"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	agentgov "github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
	agentrouter "github.com/atMagicW/go-agent-runtime/internal/usecase/router"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load("configs/app.yaml")
	if err != nil {
		logger.Fatal("load config failed", zap.Error(err))
	}

	if cfg.Database.PostgresDSN == "" {
		logger.Fatal("postgres dsn is empty")
	}
	if cfg.LLM.OpenAIAPIKey == "" {
		logger.Warn("OPENAI_API_KEY is empty, OpenAI adapter will run in placeholder mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgrepo.NewDB(ctx, cfg.Database.PostgresDSN)
	if err != nil {
		logger.Fatal("init db failed", zap.Error(err))
	}
	defer db.Close()

	sessionRepo := pgrepo.NewSessionRepository(db)
	sessionService := app.NewSessionService(sessionRepo)

	// 统一治理组件
	breakers := agentgov.NewBreakerRegistry()
	fallbacks := agentgov.NewDefaultFallbackPolicy()

	openAIClient := openaiadapter.NewClient(cfg.LLM.OpenAIAPIKey)
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
	embeddingProvider := mockembedding.NewProvider(cfg.RAG.EmbeddingDim)
	ragService := app.NewRAGService(ragRepo, embeddingProvider)

	splitter := textsplitter.NewSplitter(cfg.TextSplitter.ChunkSize, cfg.TextSplitter.Overlap)
	ingestService := app.NewIngestService(ragRepo, embeddingProvider, splitter)
	// 确保演示知识库存在
	if cfg.App.Env == "local" {
		seedKnowledgeBases(ctx, logger, ragRepo)
		seedDemoData(ctx, logger, ragRepo, embeddingProvider)
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
		ingestService,
	)

	router := gin.Default()
	httpapi.RegisterRoutes(router, handler)

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Info("Agent Runtime Server starting", zap.String("addr", addr), zap.String("env", cfg.App.Env))

	if err := router.Run(addr); err != nil {
		logger.Fatal("server start failed", zap.Error(err))
	}
}

func seedKnowledgeBases(ctx context.Context, logger *zap.Logger, ragRepo *pgrag.Repository) {
	kbs := []rag.KnowledgeBase{
		{
			KBID:        "default",
			TenantID:    "default",
			Name:        "Default Knowledge Base",
			Description: "默认演示知识库",
			Enabled:     true,
		},
		{
			KBID:        "knowledge_a",
			TenantID:    "default",
			Name:        "Knowledge A",
			Description: "演示知识库 A",
			Enabled:     true,
		},
		{
			KBID:        "knowledge_b",
			TenantID:    "default",
			Name:        "Knowledge B",
			Description: "演示知识库 B",
			Enabled:     true,
		},
	}

	for _, kb := range kbs {
		if err := ragRepo.EnsureKnowledgeBase(ctx, kb); err != nil {
			logger.Fatal("ensure knowledge base failed", zap.String("kb_id", kb.KBID), zap.Error(err))
		}
	}
}

func seedDemoData(ctx context.Context, logger *zap.Logger, ragRepo *pgrag.Repository, embeddingProvider ports.EmbeddingProvider) {
	for _, kbID := range []string{"default", "knowledge_a", "knowledge_b"} {
		if err := pgrag.SeedDemoKnowledgeBase(ctx, ragRepo, embeddingProvider, kbID); err != nil {
			logger.Warn("seed knowledge base failed", zap.String("kb_id", kbID), zap.Error(err))
		}
	}
}
