package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
)

func (m *mockCapabilityService) ListCapabilities() []capability.Descriptor {
	return []capability.Descriptor{
		{
			Name:    "resume_analyzer",
			Kind:    capability.KindSkill,
			Enabled: true,
		},
	}
}

func TestListCapabilitiesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := &Handler{
		capabilityService: &mockCapabilityService{},
	}
	r.GET("/v1/capabilities", h.ListCapabilitiesHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/capabilities", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
