package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func (m *mockCapabilityService) ListCapabilities() []CapabilityView {
	return []CapabilityView{
		{
			Name:    "resume_analyzer",
			Kind:    "skill",
			Source:  "local",
			Enabled: true,
		},
		{
			Name:       "mcp_web_search",
			Kind:       "mcp_tool",
			Source:     "remote",
			ServerName: "search-server",
			RemoteTool: "web_search",
			Enabled:    true,
		},
	}
}

func TestListCapabilitiesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		path      string
		wantCount int
		wantName  string
	}{
		{name: "all", path: "/v1/capabilities", wantCount: 2, wantName: "mcp_web_search"},
		{name: "filter by source", path: "/v1/capabilities?source=remote", wantCount: 1, wantName: "mcp_web_search"},
		{name: "filter by kind", path: "/v1/capabilities?kind=skill", wantCount: 1, wantName: "resume_analyzer"},
		{name: "filter by kind and source", path: "/v1/capabilities?kind=mcp_tool&source=remote", wantCount: 1, wantName: "mcp_web_search"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			h := &Handler{
				capabilityService: &mockCapabilityService{},
			}
			r.GET("/v1/capabilities", h.ListCapabilitiesHandler)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
			}

			var body struct {
				Capabilities []CapabilityView `json:"capabilities"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("unmarshal response error = %v", err)
			}

			if len(body.Capabilities) != tt.wantCount {
				t.Fatalf("capabilities count = %d, want %d", len(body.Capabilities), tt.wantCount)
			}

			if len(body.Capabilities) > 0 && body.Capabilities[len(body.Capabilities)-1].Name != tt.wantName {
				t.Fatalf("last capability name = %s, want %s", body.Capabilities[len(body.Capabilities)-1].Name, tt.wantName)
			}
		})
	}
}
