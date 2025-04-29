package mcp

import (
	"errors"
	"github.com/devlibx/goc-mcp/api"
	mcpGo "github.com/metoro-io/mcp-golang"
)

func RegisterTool[T any](server mcpApi.Server, name string, description string, handler func(input T) (*mcpGo.ToolResponse, error)) error {
	if s, ok := server.(*serverImpl); ok {
		return s.server.RegisterTool(name, description, handler)
	} else {
		return errors.New("invalid server type")
	}
}
