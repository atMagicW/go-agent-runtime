package app

import (
	"context"
	"fmt"
	"path/filepath"

	deepseekadapter "github.com/atMagicW/go-agent-runtime/internal/adapters/llm/deepseek"
	openaiadapter "github.com/atMagicW/go-agent-runtime/internal/adapters/llm/openai"
	filerepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/file"
	memrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/memory"
	pgrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/postgres"
	promptrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/prompt"
	memrag "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/memory"
	mockembedding "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/mock_embedding"
	openaiembedding "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/openai_embedding"
	pgrag "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/pgvector"
	rerankadapter "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/rerank"
	"github.com/atMagicW/go-agent-runtime/internal/adapters/skillloader"
	"github.com/atMagicW/go-agent-runtime/internal/domain/model"
	"github.com/atMagicW/go-agent-runtime/internal/pkg/config"
	"github.com/atMagicW/go-agent-runtime/internal/pkg/textsplitter"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	agentgov "github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
	agentrouter "github.com/atMagicW/go-agent-runtime/internal/usecase/router"
)

type BootstrapResult struct {
	AgentService      *AgentService
	SessionService    *SessionService
	CapabilityService *CapabilityService
	IngestService     *IngestService
	PromptService     *PromptService
	MCPService        *MCPService
	SkillService      *SkillService
	CloseFn           func()
}

func Bootstrap(ctx context.Context, appCfg *config.Config) (*BootstrapResult, error) {
	modelCfg, err := config.LoadModels("configs/models.yaml")
	if err != nil {
		return nil, fmt.Errorf("load models config failed: %w", err)
	}

	capCfg, err := config.LoadCapabilities("configs/capabilities.yaml")
	if err != nil {
		return nil, fmt.Errorf("load capabilities config failed: %w", err)
	}

	kbCfg, err := config.LoadKnowledgeBases("configs/knowledge_bases.yaml")
	if err != nil {
		return nil, fmt.Errorf("load knowledge bases config failed: %w", err)
	}

	fallbackCfg, err := config.LoadFallback("configs/fallback.yaml")
	if err != nil {
		return nil, fmt.Errorf("load fallback config failed: %w", err)
	}

	pricingCfg, err := config.LoadPricing("configs/pricing.yaml")
	if err != nil {
		return nil, fmt.Errorf("load pricing config failed: %w", err)
	}

	pricingService := NewPricingService(pricingCfg)
	breakers := agentgov.NewBreakerRegistry()
	fallbacks := agentgov.NewFallbackPolicyFromConfig(fallbackCfg)

	sessionRepo, ragRepo, closeFn, err := buildStorage(ctx, appCfg)
	if err != nil {
		return nil, err
	}

	sessionService := NewSessionService(sessionRepo)

	modelRegistry := NewModelRegistry(modelCfg)

	openAIClient := openaiadapter.NewClient(appCfg.LLM.OpenAIAPIKey, pricingService)
	deepSeekClient := deepseekadapter.NewClient(
		appCfg.LLM.DeepSeekAPIKey,
		appCfg.LLM.DeepSeekBaseURL,
		pricingService,
	)

	llmClients := map[string]ports.LLMClient{
		string(model.ProviderOpenAI):   openAIClient,
		string(model.ProviderDeepSeek): deepSeekClient,
	}

	modelRouter := agentrouter.NewModelRouter(
		llmClients,
		modelRegistry,
		breakers,
		fallbacks,
	)

	embeddingProvider := buildEmbeddingProvider(appCfg, pricingService)
	var reranker ports.Reranker
	if appCfg.RAG.RerankEnabled {
		reranker = rerankadapter.NewSimpleReranker()
	}

	ragService := NewRAGService(ragRepo, embeddingProvider, reranker)
	splitter := textsplitter.NewSplitter(appCfg.TextSplitter.ChunkSize, appCfg.TextSplitter.Overlap)
	ingestService := NewIngestService(ragRepo, embeddingProvider, splitter)

	// skill
	skillLoader := skillloader.NewFileLoader("skills")
	skillDefs, err := skillLoader.Load()
	if err != nil {
		skillDefs = nil
	}
	skillRegistry := NewSkillRegistry(skillDefs)

	// mcp
	mcpClient := buildMCPClient()

	// capability registry
	registry := BuildCapabilityRegistry(capCfg, skillRegistry, mcpClient)

	// knowledge bases
	if err := InitKnowledgeBases(ctx, ragRepo, embeddingProvider, kbCfg, appCfg.RAG.SeedOnBootstrap); err != nil {
		return nil, fmt.Errorf("init knowledge bases failed: %w", err)
	}

	// prompt repo
	var pRepo ports.PromptRepository
	filePromptRepo, err := promptrepo.NewFileRepository("prompts")
	if err != nil {
		pRepo = promptrepo.NewInMemoryRepository()
	} else {
		pRepo = filePromptRepo
	}

	agentService := NewAgentService(
		sessionService,
		modelRouter,
		registry,
		ragService,
		breakers,
		fallbacks,
	)

	capabilityService := NewCapabilityService(registry)
	promptService := NewPromptService(pRepo)
	mcpService := NewMCPService(capCfg)
	skillService := NewSkillService(skillRegistry)

	return &BootstrapResult{
		AgentService:      agentService,
		SessionService:    sessionService,
		CapabilityService: capabilityService,
		IngestService:     ingestService,
		PromptService:     promptService,
		MCPService:        mcpService,
		SkillService:      skillService,
		CloseFn:           closeFn,
	}, nil
}

func buildStorage(
	ctx context.Context,
	appCfg *config.Config,
) (ports.SessionRepository, ports.RAGRepository, func(), error) {
	switch appCfg.Storage.Mode {
	case "postgres":
		if appCfg.Database.PostgresDSN == "" {
			return nil, nil, nil, fmt.Errorf("postgres dsn is empty")
		}

		db, err := pgrepo.NewDB(ctx, appCfg.Database.PostgresDSN)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("init postgres failed: %w", err)
		}

		sessionRepo := pgrepo.NewSessionRepository(db)
		ragRepo := pgrag.NewRepository(db)

		return sessionRepo, ragRepo, func() {
			db.Close()
		}, nil

	case "file":
		sessionRepo, err := filerepo.NewSessionRepository(filepath.Join(appCfg.Storage.DataDir, "session"))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("init file session repo failed: %w", err)
		}

		// file 模式第一版：session 持久化到文件，RAG 走内存
		ragRepo := memrag.NewRepository()

		return sessionRepo, ragRepo, func() {}, nil

	case "memory":
		sessionRepo := memrepo.NewSessionRepository()
		ragRepo := memrag.NewRepository()

		return sessionRepo, ragRepo, func() {}, nil

	default:
		return nil, nil, nil, fmt.Errorf("unsupported storage mode: %s", appCfg.Storage.Mode)
	}
}

func buildEmbeddingProvider(appCfg *config.Config, pricingService *PricingService) ports.EmbeddingProvider {
	switch appCfg.RAG.EmbeddingProvider {
	case string(model.ProviderOpenAI):
		return openaiembedding.NewProvider(
			appCfg.LLM.OpenAIAPIKey,
			appCfg.RAG.EmbeddingModel,
			appCfg.RAG.EmbeddingDim,
			pricingService,
		)
	default:
		return mockembedding.NewProvider(appCfg.RAG.EmbeddingDim)
	}
}

func buildMCPClient() ports.MCPClient {
	// 先复用你现有的 mock / adapter 实现
	return &noopMCPClient{}
}

type noopMCPClient struct{}

func (n *noopMCPClient) CallTool(ctx context.Context, req ports.MCPCallRequest) (ports.MCPCallResponse, error) {
	_ = ctx
	return ports.MCPCallResponse{
		Output: map[string]any{
			"server_name": req.ServerName,
			"tool_name":   req.ToolName,
			"result":      "mock mcp call success",
		},
	}, nil
}
