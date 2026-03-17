package httpapi

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	pgrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/postgres"
	"github.com/atMagicW/go-agent-runtime/internal/app"
)

// GetSessionHandler 获取会话详情
func GetSessionHandler(c *gin.Context) {
	sessionID := c.Param("id")

	pgDSN := os.Getenv("POSTGRES_DSN")
	if pgDSN == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "POSTGRES_DSN is not set",
		})
		return
	}

	db, err := pgrepo.NewDB(context.Background(), pgDSN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer db.Close()

	sessionRepo := pgrepo.NewSessionRepository(db)
	sessionService := app.NewSessionService(sessionRepo)

	state, err := sessionService.LoadConversationState(context.Background(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	session, err := sessionService.GetSession(context.Background(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session": session,
		"state":   state,
	})
}
