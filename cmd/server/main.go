package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/atMagicW/go-agent-runtime/api/httpapi"
	pgrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/postgres"
	"github.com/atMagicW/go-agent-runtime/internal/app"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	pgDSN := os.Getenv("POSTGRES_DSN")
	if pgDSN == "" {
		logger.Fatal("POSTGRES_DSN is not set")
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
	agentService := app.NewAgentService(sessionService)

	handler := httpapi.NewHandler(agentService, sessionService)

	router := gin.Default()
	httpapi.RegisterRoutes(router, handler)

	logger.Info("Agent Runtime Server starting at :8080")

	if err := router.Run(":8080"); err != nil {
		logger.Fatal("server start failed", zap.Error(err))
	}
}
