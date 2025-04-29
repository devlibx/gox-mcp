// Package mcp_v1 provides a simplified interface for Model Context Protocol (MCP) servers.
// This package defines the core interfaces and implementations for MCP server functionality.
package mcp

import (
	"context"
	"github.com/devlibx/gox-base/v2"
	"go.uber.org/fx"
	"goc-mcp/api"
)

var Provider = fx.Options(
	fx.Provide(NewServer),
)

// NewServer creates and returns a new instance of the Server interface.
// This function instantiates a serverImpl struct and returns it as a Server.
// It can be used directly or as part of the dependency injection system.
//
// Returns:
//   - Server: A new server instance implementing the Server interface
//   - error: An error if server creation fails, or nil on success
func NewServer(cf gox.CrossFunction, config *mcpApi.Config) (mcpApi.Server, error) {
	ns := &serverImpl{
		CrossFunction: cf,
		config:        config,
	}
	return ns, nil
}

// NewServerLifecycleInvoker registers the server's lifecycle hooks with the fx lifecycle.
// This function is used to ensure that the server's Start and Stop methods are called
// at the appropriate times during the application lifecycle.
//
// Parameters:
//   - lifecycle: The fx.Lifecycle to register hooks with
//   - server: The Server instance to manage
//
// Returns:
//   - error: An error if registration fails, or nil on success
func NewServerLifecycleInvoker(lifecycle fx.Lifecycle, server mcpApi.Server) error {
	lifecycle.Append(
		fx.Hook{
			// OnStart hook calls the server's Start method when the application starts
			OnStart: func(ctx context.Context) error {
				return server.Start()
			},
			// OnStop hook calls the server's Stop method when the application stops
			OnStop: func(ctx context.Context) error {
				return server.Stop()
			},
		},
	)
	return nil
}
