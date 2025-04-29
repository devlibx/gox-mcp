package mcp

import (
	"errors"
	mcpGo "github.com/metoro-io/mcp-golang"
	"goc-mcp/api"
)

func RegisterTool[T any](server mcpApi.Server, name string, description string, handler func(input T) (*mcpGo.ToolResponse, error)) error {
	if s, ok := server.(*serverImpl); ok {
		return s.server.RegisterTool(name, description, handler)
	} else {
		return errors.New("invalid server type")
	}
}
