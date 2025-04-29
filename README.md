# GoC-MCP: Go Client for Model Context Protocol

[![Go Reference](https://pkg.go.dev/badge/github.com/harishbohara/goc-mcp.svg)](https://pkg.go.dev/github.com/harishbohara/goc-mcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/harishbohara/goc-mcp)](https://goreportcard.com/report/github.com/harishbohara/goc-mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

GoC-MCP is a Go client implementation of the [Model Context Protocol (MCP)](https://github.com/metoro-io/mcp-golang). It provides a simple and flexible way to create MCP servers that can expose tools and resources to AI assistants.

## Features

- 🚀 Simple API for creating MCP servers
- 🔧 Easy tool registration with type safety
- 🔌 HTTP transport support
- 📦 No external dependencies for core functionality
- 🧩 Optional integration with [Uber fx](https://github.com/uber-go/fx) for dependency injection

## Installation

```bash
go get github.com/harishbohara/goc-mcp
```

## Quick Start

### Creating a Simple MCP Server

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/devlibx/gox-base/v2"
	mcpGo "github.com/metoro-io/mcp-golang"
	mcpApi "goc-mcp/api"
	goxMcp "goc-mcp/mcp-metoro-io"
)

func main() {
	// Create a cross function (required for logging and other utilities)
	cf := gox.NewNoOpCrossFunction()
	
	// Create a new server configuration
	config := &mcpApi.Config{
		Disabled: false,
		Port:     1234,
	}
	
	// Create a new server instance
	server, err := goxMcp.NewServer(cf, config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	
	// Define a content type for your tool
	type Content struct {
		Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
		Description *string `json:"description" jsonschema:"description=The description to submit"`
	}
	
	// Register a tool
	err = goxMcp.RegisterTool[Content](
		server,
		"greeting_tool",
		"A tool that generates a greeting",
		func(input Content) (*mcpGo.ToolResponse, error) {
			greeting := fmt.Sprintf("Hello, %s!", input.Title)
			if input.Description != nil {
				greeting += fmt.Sprintf(" %s", *input.Description)
			}
			
			return mcpGo.NewToolResponse(
				&mcpGo.Content{
					Type:        "text",
					TextContent: &mcpGo.TextContent{Text: greeting},
				},
			), nil
		},
	)
	if err != nil {
		log.Fatalf("Failed to register tool: %v", err)
	}
	
	// Start the server (this will block until the server is stopped)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
```

### Testing Your MCP Server

You can test your MCP server by sending a request to list all available tools:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	
	"github.com/go-resty/resty/v2"
	"goc-mcp/api"
)

func main() {
	// Create a REST client
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("http://localhost:%d/mcp", 1234))
	
	// Create a tool list request
	request := api.NewMcpToolListRequest()
	
	// Send the request
	resp, err := client.R().SetBody(request).Post("")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	
	// Parse the response
	var toolListResponse api.McpTollListResponse
	if err := json.Unmarshal(resp.Body(), &toolListResponse); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}
	
	// Print the available tools
	fmt.Printf("Available tools (%d):\n", len(toolListResponse.Result.Tools))
	for _, tool := range toolListResponse.Result.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}
}
```

## API Documentation

### Server Interface

The core of the library is the `Server` interface:

```go
type Server interface {
	// Start initializes and begins the server operations.
	Start() error

	// Stop gracefully terminates the server operations.
	Stop() error
}
```

### Configuration

Configure your MCP server using the `Config` struct:

```go
type Config struct {
	Disabled bool `yaml:"disabled" json:"disabled"`
	Port     int  `yaml:"port" json:"port"`
}
```

### Tool Registration

Register tools with type safety using the `RegisterTool` function:

```go
func RegisterTool[T any](
	server api.Server,
	name string,
	description string,
	handler func(T) (*mcpGo.ToolResponse, error),
) error
```

### JSON-RPC Requests and Responses

The library provides structs for JSON-RPC requests and responses:

```go
// Request to list tools
type McpToolListRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
}

// Response with tool list
type McpTollListResponse struct {
	ID      interface{}    `json:"id"`
	JSONRPC string         `json:"jsonrpc"`
	Result  *ToolListResult `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}
```

## Advanced Usage

### Creating a Simplified API Wrapper

If you prefer not to use the `gox.CrossFunction` dependency directly in your code, you can create a simple wrapper function:

```go
package main

import (
	"github.com/devlibx/gox-base/v2"
	mcpApi "goc-mcp/api"
	goxMcp "goc-mcp/mcp-metoro-io"
)

// SimpleNewServer creates a new MCP server with minimal configuration.
// It handles the creation of the CrossFunction dependency internally.
func SimpleNewServer(port int) (mcpApi.Server, error) {
	// Create a no-op cross function
	cf := gox.NewNoOpCrossFunction()
	
	// Create a new server configuration
	config := &mcpApi.Config{
		Disabled: false,
		Port:     port,
	}
	
	// Create and return the server
	return goxMcp.NewServer(cf, config)
}
```

Then you can use this simplified function in your code:

```go
// Create a new server instance with the simplified API
server, err := SimpleNewServer(1234)
if err != nil {
	log.Fatalf("Failed to create server: %v", err)
}
```

### Custom Tool Implementations

You can create custom tools by implementing the tool handler function:

```go
func myToolHandler(input MyInputType) (*mcpGo.ToolResponse, error) {
	// Process the input
	result := processInput(input)
	
	// Return a response
	return mcpGo.NewToolResponse(
		&mcpGo.Content{
			Type:        "text",
			TextContent: &mcpGo.TextContent{Text: result},
		},
	), nil
}
```

### Error Handling

Handle errors in your tool implementations:

```go
func myToolHandler(input MyInputType) (*mcpGo.ToolResponse, error) {
	// Validate input
	if err := validateInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	
	// Process the input
	result, err := processInput(input)
	if err != nil {
		return nil, fmt.Errorf("processing error: %w", err)
	}
	
	// Return a response
	return mcpGo.NewToolResponse(
		&mcpGo.Content{
			Type:        "text",
			TextContent: &mcpGo.TextContent{Text: result},
		},
	), nil
}
```

### Using with Uber fx

The library provides built-in support for [Uber fx](https://github.com/uber-go/fx), a dependency injection framework. This makes it easy to integrate the MCP server into your existing fx-based applications:

```go
package main

import (
	"context"
	"github.com/devlibx/gox-base/v2"
	"go.uber.org/fx"
	mcpApi "goc-mcp/api"
	goxMcp "goc-mcp/mcp-metoro-io"
)

func main() {
	// Create an fx application
	app := fx.New(
		// Provide the cross function
		fx.Provide(gox.NewNoOpCrossFunction),
		
		// Provide the server configuration
		fx.Supply(&mcpApi.Config{
			Disabled: false,
			Port:     1234,
		}),
		
		// Provide the MCP server
		goxMcp.Provider,
		
		// Register the server lifecycle hooks
		fx.Invoke(goxMcp.NewServerLifecycleInvoker),
		
		// Register your tools
		fx.Invoke(registerTools),
	)
	
	// Start the application
	app.Run()
}

// registerTools registers your tools with the MCP server
func registerTools(server mcpApi.Server) error {
	// Register your tools here
	// ...
	
	return nil
}
```

This approach allows you to leverage fx's dependency injection capabilities to manage the lifecycle of your MCP server and its dependencies.

## Project Structure

The project is organized into the following packages:

- `api`: Core interfaces and types used throughout the project
  - `Server`: The main interface for MCP servers
  - `Config`: Configuration for MCP servers
  - JSON-RPC request and response types

- `mcp-metoro-io`: Implementation of the MCP server using the metoro-io/mcp-golang library
  - Server implementation
  - Tool registration
  - Integration with Uber fx

- `examples`: Example applications demonstrating how to use the library
  - `simple_hello_tool`: A simple example of creating an MCP server with a greeting tool

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
