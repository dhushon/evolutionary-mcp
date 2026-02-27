package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
			mcp.WithDescription("Create a new memory"),
			mcp.WithString("content", mcp.Required(), mcp.Description("The content of the memory")),
		),
		s.handleRemember,
	)

	s.mcpServer.AddTool(
		mcp.NewTool(
			"recall",
			mcp.WithDescription("Recall memories based on a query"),
			mcp.WithString("query", mcp.Required(), mcp.Description("The query to search for")),
		),
		s.handleRecall,
	)

	s.mcpServer.AddTool(
		mcp.NewTool(
			"give_feedback",
			mcp.WithDescription("Give feedback on a memory"),
			mcp.WithString("id", mcp.Required(), mcp.Description("The ID of the memory")),
			mcp.WithNumber("confidence", mcp.Required(), mcp.Description("The new confidence score for the memory")),
		),
		s.handleGiveFeedback,
	)
}

func (s *Server) handleRemember(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	return mcp.NewToolResultText("Feedback received"), nil
}


func MountHTTPHandlers(mux *http.ServeMux, mcpServer *server.MCPServer) {
	// Use SSE server for /mcp/sse and /mcp/message endpoints
	sseServer := server.NewSSEServer(mcpServer, server.WithStaticBasePath("/mcp"))
	
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		// Direct POST for tool calls
		if r.Method == http.MethodPost {
			sseServer.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
	
	// SSE endpoints
	mux.HandleFunc("/mcp/sse", sseServer.ServeHTTP)
	mux.HandleFunc("/mcp/message", sseServer.ServeHTTP)
}
