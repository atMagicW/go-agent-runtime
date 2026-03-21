package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockAgentService struct{}
type mockSessionService struct{}
type mockCapabilityService struct{}
type mockIngestService struct{}

func (m *mockAgentService) Run(reqCtx any, message string) (any, error)                { return nil, nil }
func (m *mockAgentService) RunStream(reqCtx any, message string, writer *StreamWriter) {}

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := &Handler{}
	r.GET("/v1/health", h.HealthHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
