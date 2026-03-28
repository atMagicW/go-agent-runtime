package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

type serverConfig struct {
	Name      string
	Mode      string
	BaseURL   string
	ToolPath  string
	TimeoutMS int
}

// Client 支持 mock / http 两种模式的 MCP 客户端。
type Client struct {
	servers     map[string]serverConfig
	httpClient  *http.Client
	mockResults map[string]map[string]any
}

// NewClient 创建 MCP 客户端。
func NewClient(servers []cfg.MCPServerConfig) *Client {
	items := make(map[string]serverConfig, len(servers))
	for _, server := range servers {
		mode := strings.TrimSpace(strings.ToLower(server.Mode))
		if mode == "" {
			mode = "mock"
		}

		toolPath := strings.TrimSpace(server.ToolPath)
		if toolPath == "" {
			toolPath = "/tools/call"
		}

		timeoutMS := server.TimeoutMS
		if timeoutMS <= 0 {
			timeoutMS = 5000
		}

		items[server.Name] = serverConfig{
			Name:      server.Name,
			Mode:      mode,
			BaseURL:   strings.TrimRight(strings.TrimSpace(server.BaseURL), "/"),
			ToolPath:  toolPath,
			TimeoutMS: timeoutMS,
		}
	}

	return &Client{
		servers:    items,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		mockResults: map[string]map[string]any{
			"web_search": {
				"result": "mock web search success",
			},
			"news_search": {
				"result": "mock news search success",
			},
			"doc_lookup": {
				"result": "mock doc lookup success",
			},
		},
	}
}

// CallTool 调用远程 MCP Tool。
func (c *Client) CallTool(ctx context.Context, req ports.MCPCallRequest) (ports.MCPCallResponse, error) {
	server, ok := c.servers[req.ServerName]
	if !ok {
		return ports.MCPCallResponse{}, fmt.Errorf("mcp server not configured: %s", req.ServerName)
	}

	if server.Mode == "http" && server.BaseURL != "" {
		return c.callHTTP(ctx, server, req)
	}

	return c.callMock(req), nil
}

func (c *Client) callMock(req ports.MCPCallRequest) ports.MCPCallResponse {
	payload := map[string]any{
		"server_name": req.ServerName,
		"tool_name":   req.ToolName,
		"transport":   "mock",
		"result":      fmt.Sprintf("mcp tool %s executed on server %s", req.ToolName, req.ServerName),
	}

	if extra, ok := c.mockResults[req.ToolName]; ok {
		for k, v := range extra {
			payload[k] = v
		}
	}

	return ports.MCPCallResponse{Output: payload}
}

func (c *Client) callHTTP(
	ctx context.Context,
	server serverConfig,
	req ports.MCPCallRequest,
) (ports.MCPCallResponse, error) {
	body, err := json.Marshal(map[string]any{
		"tool_name": req.ToolName,
		"input":     req.Input,
	})
	if err != nil {
		return ports.MCPCallResponse{}, fmt.Errorf("marshal mcp request failed: %w", err)
	}

	url := joinURL(server.BaseURL, server.ToolPath)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return ports.MCPCallResponse{}, fmt.Errorf("create mcp http request failed: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := c.httpClient
	if client == nil {
		client = &http.Client{}
	}
	if server.TimeoutMS > 0 {
		client = &http.Client{
			Timeout:   time.Duration(server.TimeoutMS) * time.Millisecond,
			Transport: client.Transport,
		}
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return ports.MCPCallResponse{}, fmt.Errorf("call mcp http server failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return ports.MCPCallResponse{}, fmt.Errorf("read mcp http response failed: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return ports.MCPCallResponse{}, fmt.Errorf("mcp http server returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return ports.MCPCallResponse{}, fmt.Errorf("decode mcp http response failed: %w", err)
	}

	if output, ok := decoded["output"].(map[string]any); ok {
		if _, exists := output["transport"]; !exists {
			output["transport"] = "http"
		}
		return ports.MCPCallResponse{Output: output}, nil
	}

	decoded["transport"] = "http"
	return ports.MCPCallResponse{Output: decoded}, nil
}

func joinURL(baseURL, path string) string {
	if path == "" {
		return baseURL
	}
	if strings.HasPrefix(path, "/") {
		return baseURL + path
	}
	return baseURL + "/" + path
}
