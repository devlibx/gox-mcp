// Package mcp_v1 provides a simplified interface for Model Context Protocol (MCP) servers.
// This package defines the core interfaces and implementations for MCP server functionality.
package mcpApi

type Config struct {
	Disabled bool `yaml:"disabled" json:"disabled"`
	Port     int  `yaml:"port" json:"port"`
}

// Server defines the interface for an MCP server.
// Implementations of this interface are responsible for managing the lifecycle
// of an MCP server, including starting and stopping server operations.
type Server interface {
	// Start initializes and begins the server operations.
	// It should set up any necessary resources, open connections,
	// and begin listening for requests.
	// Returns an error if the server fails to start properly.
	Start() error

	// Stop gracefully terminates the server operations.
	// It should clean up resources, close connections,
	// and ensure all pending operations are completed or properly cancelled.
	// Returns an error if the server fails to stop properly.
	Stop() error
}
