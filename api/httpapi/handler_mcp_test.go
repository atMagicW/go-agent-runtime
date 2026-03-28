package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockMCPService struct{}

func (m *mockMCPService) ListServers() []MCPServerView {
	return []MCPServerView{
		{
			Name:      "search-server",
			Mode:      "mock",
			Enabled:   true,
			ToolPath:  "/tools/call",
			TimeoutMS: 5000,
			Tools: []MCPToolView{
				{
					Name:       "mcp_web_search",
					RemoteTool: "web_search",
					Enabled:    true,
				},
			},
		},
		{
			Name:      "docs-server",
			Mode:      "http",
			BaseURL:   "https://mcp.example.com",
			ToolPath:  "/tools/call",
			TimeoutMS: 3000,
			Enabled:   true,
		},
	}
}

func TestListMCPServersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := &Handler{
		mcpService: &mockMCPService{},
	}
	r.GET("/v1/mcp/servers", h.ListMCPServersHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/mcp/servers", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body struct {
		Servers []MCPServerView `json:"servers"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response error = %v", err)
	}

	if len(body.Servers) != 2 {
		t.Fatalf("servers count = %d, want 2", len(body.Servers))
	}

	if body.Servers[0].Mode != "mock" {
		t.Fatalf("first server mode = %s, want mock", body.Servers[0].Mode)
	}

	if body.Servers[1].BaseURL != "https://mcp.example.com" {
		t.Fatalf("second server base_url = %s, want https://mcp.example.com", body.Servers[1].BaseURL)
	}
}
