package mcp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestClientCallToolMock(t *testing.T) {
	client := NewClient([]cfg.MCPServerConfig{
		{
			Name:    "search-server",
			Mode:    "mock",
			Enabled: true,
		},
	})

	resp, err := client.CallTool(context.Background(), ports.MCPCallRequest{
		ServerName: "search-server",
		ToolName:   "web_search",
		Input:      map[string]any{"query": "agent runtime"},
	})
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}

	if got := resp.Output["transport"]; got != "mock" {
		t.Fatalf("transport = %v, want mock", got)
	}
	if got := resp.Output["result"]; got != "mock web search success" {
		t.Fatalf("result = %v, want mock web search success", got)
	}
}

func TestClientCallToolHTTP(t *testing.T) {
	client := NewClient([]cfg.MCPServerConfig{
		{
			Name:      "docs-server",
			Mode:      "http",
			BaseURL:   "http://mcp.example.local",
			ToolPath:  "/tools/call",
			TimeoutMS: 1000,
			Enabled:   true,
		},
	})
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/tools/call" {
				t.Fatalf("path = %s, want /tools/call", r.URL.Path)
			}

			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body error = %v", err)
			}

			if got := body["tool_name"]; got != "doc_lookup" {
				t.Fatalf("tool_name = %v, want doc_lookup", got)
			}

			payload, err := json.Marshal(map[string]any{
				"output": map[string]any{
					"result": "real doc lookup success",
				},
			})
			if err != nil {
				t.Fatalf("marshal response error = %v", err)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(string(payload))),
			}, nil
		}),
	}

	resp, err := client.CallTool(context.Background(), ports.MCPCallRequest{
		ServerName: "docs-server",
		ToolName:   "doc_lookup",
		Input:      map[string]any{"query": "planner"},
	})
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}

	if got := resp.Output["transport"]; got != "http" {
		t.Fatalf("transport = %v, want http", got)
	}
	if got := resp.Output["result"]; got != "real doc lookup success" {
		t.Fatalf("result = %v, want real doc lookup success", got)
	}
}
