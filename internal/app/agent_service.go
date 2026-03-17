package app

import (
	"github.com/atMagicW/go-agent-runtime/api/sse"
	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// AgentService 是Agent运行时入口
type AgentService struct {
}

// NewAgentService 创建AgentService
func NewAgentService() *AgentService {

	return &AgentService{}
}

// Run 非流式执行
func (s *AgentService) Run(
	ctx agent.RequestContext,
	message string,
) (*agent.FinalResponse, error) {

	// TODO:
	// 1. 构建上下文
	// 2. 意图识别
	// 3. 生成计划
	// 4. 执行计划
	// 5. 返回最终回复

	return &agent.FinalResponse{
		Message: "Agent runtime demo running",
	}, nil
}

// RunStream 流式执行
func (s *AgentService) RunStream(
	ctx agent.RequestContext,
	message string,
	writer *httpapi.StreamWriter,
) {

	writer.WriteEvent("plan", "creating execution plan")

	// TODO: 调用Orchestrator

	writer.WriteToken("Hello ")
	writer.WriteToken("this ")
	writer.WriteToken("is ")
	writer.WriteToken("agent ")
	writer.WriteToken("runtime")

	writer.WriteEvent("done", "completed")
}
