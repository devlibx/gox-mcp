package mcp

import (
	"errors"
	"github.com/devlibx/gox-mcp/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTool[T any](server mcpApi.Server, tool mcp.Tool, handlerFunc server.ToolHandlerFunc) error {
	if s, ok := server.(*serverImpl); ok {
		s.McpServer.AddTool(tool, handlerFunc)
		return nil
	} else {
		return errors.New("invalid server type")
	}
}
