package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"evolutionary-mcp/backend/internal/contextutil"
	"evolutionary-mcp/backend/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	mcpServer    *server.MCPServer
	memoryService *services.MemoryService
}

func NewServer(memoryService *services.MemoryService) *Server {
	s := &Server{
		mcpServer: server.NewMCPServer(
			"Evolutionary Memory",
			"1.0.0",
			server.WithToolCapabilities(true),
		),
		memoryService: memoryService,
	}

	s.registerTools()
	return s
}

func (s *Server) GetMCPServer() *server.MCPServer {
	return s.mcpServer
}

func (s *Server) registerTools() {
	s.mcpServer.AddTool(
		mcp.NewTool(
			"remember",
			mcp.WithDescription("Create a new semantic memory"),
			mcp.WithString("content", mcp.Required(), mcp.Description("The content of the memory")),
		),
		s.handleRemember,
	)

	s.mcpServer.AddTool(
		mcp.NewTool(
			"recall",
			mcp.WithDescription("Recall semantic memories based on a natural language query"),
			mcp.WithString("query", mcp.Required(), mcp.Description("The query to search for")),
		),
		s.handleRecall,
	)

	s.mcpServer.AddTool(
		mcp.NewTool(
			"give_feedback",
			mcp.WithDescription("Evolve a memory by providing a new confidence score"),
			mcp.WithString("id", mcp.Required(), mcp.Description("The ID of the memory")),
			mcp.WithNumber("confidence", mcp.Required(), mcp.Description("The new confidence score (0.0 to 1.0)")),
		),
		s.handleGiveFeedback,
	)

	s.mcpServer.AddTool(
		mcp.NewTool(
			"list_grounding_rules",
			mcp.WithDescription("Retrieve foundational grounding rules and reasoning constraints"),
		),
		s.handleListGroundingRules,
	)
}

// withAmbientContext ensures the correct tenant identity is present in the context.
// In a real implementation, this would extract tenant information from the MCP session
// or connection metadata. For now, it defaults to a 'mcp-user' tenant if none is provided.
func (s *Server) withAmbientContext(ctx context.Context) context.Context {
	tenantID := contextutil.GetTenant(ctx)
	if tenantID == "" {
		// Placeholder: In production, map the MCP connection to a Tenant
		tenantID = "default" 
	}
	return contextutil.WithTenant(ctx, tenantID)
}

func (s *Server) handleRemember(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = s.withAmbientContext(ctx)
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments type"), nil
	}

	content, ok := args["content"].(string)
	if !ok || content == "" {
		return mcp.NewToolResultError("Missing required parameter: content"), nil
	}

	memory, err := s.memoryService.Remember(ctx, content)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remember: %v", err)), nil
	}

	jsonBytes, _ := json.Marshal(memory)
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (s *Server) handleRecall(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = s.withAmbientContext(ctx)
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments type"), nil
	}

	query, ok := args["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("Missing required parameter: query"), nil
	}

	memories, err := s.memoryService.Recall(ctx, query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to recall: %v", err)), nil
	}

	jsonBytes, _ := json.Marshal(memories)
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (s *Server) handleGiveFeedback(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = s.withAmbientContext(ctx)
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments type"), nil
	}

	id, ok := args["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("Missing required parameter: id"), nil
	}

	confidence, ok := args["confidence"].(float64)
	if !ok {
		return mcp.NewToolResultError("Missing required parameter: confidence"), nil
	}

	err := s.memoryService.GiveFeedback(ctx, id, confidence)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to give feedback: %v", err)), nil
	}

	return mcp.NewToolResultText("Feedback received and memory evolved"), nil
}

func (s *Server) handleListGroundingRules(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = s.withAmbientContext(ctx)
	
	rules, err := s.memoryService.GetGroundingRules(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list grounding rules: %v", err)), nil
	}

	jsonBytes, _ := json.Marshal(rules)
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func MountHTTPHandlers(mux *http.ServeMux, mcpServer *server.MCPServer) {
	sseServer := server.NewSSEServer(mcpServer, server.WithStaticBasePath("/mcp"))
	
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sseServer.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
	
	mux.HandleFunc("/mcp/sse", sseServer.ServeHTTP)
	mux.HandleFunc("/mcp/message", sseServer.ServeHTTP)
}
