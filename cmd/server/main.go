package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/atMagicW/go-agent-runtime/api/httpapi"
	"github.com/atMagicW/go-agent-runtime/internal/app"
	"github.com/atMagicW/go-agent-runtime/internal/pkg/config"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	appCfg, err := config.Load("configs/app.yaml")
	if err != nil {
		logger.Fatal("load app config failed", zap.Error(err))
	}

	bootstrapCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	boot, err := app.Bootstrap(bootstrapCtx, appCfg)
	if err != nil {
		logger.Fatal("bootstrap failed", zap.Error(err))
	}
	defer boot.CloseFn()

	handler := httpapi.NewHandler(
		boot.AgentService,
		boot.SessionService,
		boot.CapabilityService,
		boot.IngestService,
		boot.PromptService,
		boot.MCPService,
		boot.SkillService,
	)

	router := gin.Default()
	httpapi.RegisterRoutes(router, handler)

	addr := fmt.Sprintf(":%d", appCfg.App.Port)
	logger.Info("Agent Runtime Server starting",
		zap.String("addr", addr),
		zap.String("env", appCfg.App.Env),
		zap.String("storage_mode", appCfg.Storage.Mode),
	)

	if err := router.Run(addr); err != nil {
		logger.Fatal("server start failed", zap.Error(err))
	}
}
