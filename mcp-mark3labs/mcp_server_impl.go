// Package mcp_v1 provides a simplified interface for Model Context Protocol (MCP) servers.
// This package defines the core interfaces and implementations for MCP server functionality.
package mcp

import (
	"fmt"
	"github.com/devlibx/gox-base/v2"
	"github.com/devlibx/gox-mcp/api"
	"github.com/mark3labs/mcp-go/server"
)

// serverImpl is the concrete implementation of the Server interface.
// It provides the actual functionality for starting and stopping an MCP server.
// This implementation is a minimal placeholder that can be extended with
// actual server logic as needed.
type serverImpl struct {
	gox.CrossFunction
	config *mcpApi.Config

	McpServer *server.MCPServer
	sseServer *server.SSEServer
}

func (s *serverImpl) internalStart() error {
	s.McpServer = server.NewMCPServer("mcp", "1.0.0")
	s.sseServer = server.NewSSEServer(s.McpServer,
		server.WithBaseURL(fmt.Sprintf(":%d", s.config.Port)),
	)
	return nil
}

func (s *serverImpl) Start() error {

	go func() {
		if err := s.sseServer.Start(fmt.Sprintf(":%d", s.config.Port)); err != nil {
			panic(err)
		}
	}()
	return nil
}

func (s *serverImpl) Stop() error {
	return nil
}
