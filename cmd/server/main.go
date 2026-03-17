package main

import (
	"github.com/atMagicW/go-agent-runtime/api/httpapi"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	// 创建 Gin Router
	router := gin.Default()

	// 初始化 HTTP 路由
	httpapi.RegisterRoutes(router)

	// 启动服务
	logger.Info("Agent Runtime Server starting at :8080")

	err := router.Run(":8080")
	if err != nil {
		logger.Fatal("server start failed", zap.Error(err))
	}
}
