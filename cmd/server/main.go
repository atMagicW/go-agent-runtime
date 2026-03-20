package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/atMagicW/go-agent-runtime/api/httpapi"
	mcpcap "github.com/atMagicW/go-agent-runtime/internal/adapters/capability/mcp"
	openaiadapter "github.com/atMagicW/go-agent-runtime/internal/adapters/llm/openai"
	pgrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/postgres"
	mockembedding "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/mock_embedding"
	pgrag "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/pgvector"
	"github.com/atMagicW/go-agent-runtime/internal/app"
	"github.com/atMagicW/go-agent-runtime/internal/pkg/config"
	"github.com/atMagicW/go-agent-runtime/internal/pkg/textsplitter"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	agentgov "github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
	agentrouter "github.com/atMagicW/go-agent-runtime/internal/usecase/router"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	appCfg, err := config.Load("configs/app.yaml")
	if err != nil {
		logger.Fatal("load app config failed", zap.Error(err))
	}

	modelCfg, err := config.LoadModels("configs/models.yaml")
	if err != nil {
		logger.Fatal("load models config failed", zap.Error(err))
	}

	capCfg, err := config.LoadCapabilities("configs/capabilities.yaml")
	if err != nil {
		logger.Fatal("load capabilities config failed", zap.Error(err))
	}

	kbCfg, err := config.LoadKnowledgeBases("configs/knowledge_bases.yaml")
	if err != nil {
		logger.Fatal("load knowledge bases config failed", zap.Error(err))
	}

	fallbackCfg, err := config.LoadFallback("configs/fallback.yaml")
	if err != nil {
		logger.Fatal("load fallback config failed", zap.Error(err))
	}

	if appCfg.Database.PostgresDSN == "" {
		logger.Fatal("postgres dsn is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgrepo.NewDB(ctx, appCfg.Database.PostgresDSN)
	if err != nil {
		logger.Fatal("init db failed", zap.Error(err))
	}
	defer db.Close()

	sessionRepo := pgrepo.NewSessionRepository(db)
	sessionService := app.NewSessionService(sessionRepo)

	breakers := agentgov.NewBreakerRegistry()
	fallbacks := agentgov.NewFallbackPolicyFromConfig(fallbackCfg)

	modelRegistry := app.NewModelRegistry(modelCfg)

	openAIClient := openaiadapter.NewClient(appCfg.LLM.OpenAIAPIKey)
	llmClients := map[string]ports.LLMClient{
		"openai": openAIClient,
	}

	modelRouter := agentrouter.NewModelRouter(
		llmClients,
		modelRegistry,
		breakers,
		fallbacks,
	)

	mcpClient := mcpcap.NewClient()
	registry := app.BuildCapabilityRegistry(capCfg, mcpClient)

	// 初始化 RAG
	ragRepo := pgrag.NewRepository(db)
	embeddingProvider := mockembedding.NewProvider(appCfg.RAG.EmbeddingDim)
	ragService := app.NewRAGService(ragRepo, embeddingProvider)

	splitter := textsplitter.NewSplitter(appCfg.TextSplitter.ChunkSize, appCfg.TextSplitter.Overlap)
	ingestService := app.NewIngestService(ragRepo, embeddingProvider, splitter)

	if err := app.InitKnowledgeBases(ctx, ragRepo, embeddingProvider, kbCfg); err != nil {
		logger.Fatal("init knowledge bases failed", zap.Error(err))
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

	addr := fmt.Sprintf(":%d", appCfg.App.Port)
	logger.Info("Agent Runtime Server starting", zap.String("addr", addr), zap.String("env", appCfg.App.Env))

	if err := router.Run(addr); err != nil {
		logger.Fatal("server start failed", zap.Error(err))
	}
}
