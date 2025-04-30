# GoX-MCP: Go Extension for Model Context Protocol

[![Go Reference](https://pkg.go.dev/badge/github.com/devlibx/gox-mcp.svg)](https://pkg.go.dev/github.com/devlibx/gox-mcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/devlibx/gox-mcp)](https://goreportcard.com/report/github.com/devlibx/gox-mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

GoX-MCP is a robust Go implementation of the [Model Context Protocol (MCP)](https://github.com/mark3labs/mcp-go). It provides a simple, type-safe, and flexible way to create MCP servers that expose tools and resources to AI assistants, enabling seamless integration between your Go applications and AI models.

## Features

- 🚀 **Simple API**: Clean, intuitive interfaces for creating MCP servers
- 🔧 **Type-Safe Tool Registration**: Leverage Go's type system for safer tool implementations
- 🔌 **HTTP/SSE Transport**: Built-in support for HTTP and Server-Sent Events
- 📦 **Minimal Dependencies**: Core functionality with minimal external dependencies
- 🧩 **Uber FX Integration**: First-class support for [Uber FX](https://github.com/uber-go/fx) dependency injection
- 🛡️ **Production Ready**: Designed for reliability and performance in production environments

## Installation

```bash
go get github.com/devlibx/gox-mcp
```

## Quick Start

### Creating an MCP Server with Uber FX

The recommended way to create an MCP server is using Uber FX for dependency injection:

```go
package main

import (
	"context"
	"errors"
	mcpApi "github.com/devlibx/gox-mcp/api"
	goxMcp "github.com/devlibx/gox-mcp/mcp-metoro-io"
	"github.com/devlibx/gox-base/v2"
	goxJsonUtils "github.com/devlibx/gox-base/v2/serialization/utils/json"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/fx"
	"time"
)

func main() {
	// Create an FX application with MCP server
	var server mcpApi.Server
	app := fx.New(
		// Provide the cross function for logging and utilities
		fx.Provide(gox.NewNoOpCrossFunction),
		
		// Provide the server configuration
		fx.Supply(&mcpApi.Config{
			Disabled: false,
			Port:     8089,
		}),
		
		// Register the MCP server provider
		goxMcp.Provider,
		
		// Register the server lifecycle hooks
		fx.Invoke(goxMcp.NewServerLifecycleInvoker),
		
		// Populate the server variable for use outside FX
		fx.Populate(&server),
	)
	
	// Start the application
	err := app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	// Define a content type for your tool
	type content struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
	}

	// Create a new tool definition
	tool := mcp.NewTool("hello_tool",
		mcp.WithDescription("Say hello to user"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the person to greet")),
	)

	// Register the tool with a handler function
	err = goxMcp.RegisterTool[content](server, tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if paramsStr, err := goxJsonUtils.ObjectToString(request.Params); err == nil {
			if obj, err := goxJsonUtils.StringToObject[content](paramsStr); err == nil {
				return mcp.NewToolResultText("Hello, " + obj.Title + "!"), nil
			}
		}
		return nil, errors.New("name must be a string")
	})
	
	if err != nil {
		panic(err)
	}

	// Keep the server running
	time.Sleep(1 * time.Hour)
}
```

### Creating a Simple MCP Server (Without FX)

If you prefer not to use FX, you can create an MCP server directly:

```go
package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/devlibx/gox-base/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	mcpApi "github.com/devlibx/gox-mcp/api"
	goxMcp "github.com/devlibx/gox-mcp/mcp-metoro-io"
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
	
	// Create a new tool definition
	tool := mcp.NewTool("hello_tool",
		mcp.WithDescription("Say hello to user"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the person to greet")),
	)
	
	// Register the tool with a handler function
	err = goxMcp.RegisterTool[map[string]interface{}](
		server, 
		tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Extract the name parameter
			name, ok := request.Params["name"].(string)
			if !ok {
				return nil, fmt.Errorf("name must be a string")
			}
			
			// Create a greeting
			greeting := fmt.Sprintf("Hello, %s!", name)
			
			// Return the result
			return mcp.NewToolResultText(greeting), nil
		},
	)
	if err != nil {
		log.Fatalf("Failed to register tool: %v", err)
	}
	
	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
	
	// Keep the server running
	select {}
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
	"github.com/devlibx/gox-mcp/api"
)

func main() {
	// Create a REST client
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("http://localhost:%d/mcp", 8089))
	
	// Create a tool list request
	request := api.NewMcpToolListRequest()
	
	// Send the request
	resp, err := client.R().SetBody(request).Post("")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	
	// Parse the response
	var toolListResponse api.McpToolListResponse
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
	server mcpApi.Server,
	tool mcp.Tool,
	handlerFunc server.ToolHandlerFunc,
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
type McpToolListResponse struct {
	ID      interface{}    `json:"id"`
	JSONRPC string         `json:"jsonrpc"`
	Result  *ToolListResult `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}
```

## Advanced Usage

### Creating a Weather Forecast Tool

Here's an example of creating a more complex tool that provides weather forecasts:

```go
package main

import (
	"context"
	"fmt"
	mcpApi "github.com/devlibx/gox-mcp/api"
	goxMcp "github.com/devlibx/gox-mcp/mcp-metoro-io"
	"github.com/devlibx/gox-base/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/fx"
	"time"
)

// WeatherService is a mock service for getting weather data
type WeatherService struct{}

func (s *WeatherService) GetForecast(city string, days int) (string, error) {
	// In a real application, this would call an external API
	return fmt.Sprintf("Weather forecast for %s for the next %d days: Sunny with a chance of clouds", city, days), nil
}

func main() {
	// Create a weather service
	weatherService := &WeatherService{}
	
	// Create an FX application with MCP server
	var server mcpApi.Server
	app := fx.New(
		fx.Provide(gox.NewNoOpCrossFunction),
		fx.Supply(&mcpApi.Config{
			Disabled: false,
			Port:     8089,
		}),
		goxMcp.Provider,
		fx.Invoke(goxMcp.NewServerLifecycleInvoker),
		fx.Populate(&server),
	)
	
	err := app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	// Define the input type for the weather tool
	type WeatherRequest struct {
		City string `json:"city"`
		Days int    `json:"days"`
	}

	// Create a new tool definition
	weatherTool := mcp.NewTool("get_forecast",
		mcp.WithDescription("Get weather forecast for a city"),
		mcp.WithString("city", mcp.Required(), mcp.Description("City name")),
		mcp.WithInteger("days", mcp.Required(), mcp.Description("Number of days for forecast")),
	)

	// Register the weather tool
	err = goxMcp.RegisterTool[WeatherRequest](server, weatherTool, 
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Parse the request parameters
			params, ok := request.Params.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid parameters format")
			}
			
			city, ok := params["city"].(string)
			if !ok {
				return nil, fmt.Errorf("city must be a string")
			}
			
			daysFloat, ok := params["days"].(float64)
			if !ok {
				return nil, fmt.Errorf("days must be a number")
			}
			days := int(daysFloat)
			
			// Get the forecast from the weather service
			forecast, err := weatherService.GetForecast(city, days)
			if err != nil {
				return nil, err
			}
			
			// Return the forecast as a text result
			return mcp.NewToolResultText(forecast), nil
		},
	)
	
	if err != nil {
		panic(err)
	}

	// Keep the server running
	time.Sleep(1 * time.Hour)
}
```

### Using with Dependency Injection

GoX-MCP is designed to work seamlessly with Uber FX for dependency injection. Here's a more complex example showing how to integrate with FX and inject dependencies:

```go
package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base/v2"
	"github.com/devlibx/gox-base/v2/metrics"
	mcpApi "github.com/devlibx/gox-mcp/api"
	goxMcp "github.com/devlibx/gox-mcp/mcp-metoro-io"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

// AppConfig holds the application configuration
type AppConfig struct {
	MCP mcpApi.Config `yaml:"mcp" json:"mcp"`
}

// WeatherService provides weather data
type WeatherService struct {
	logger *zap.Logger
	metrics metrics.Metrics
}

// NewWeatherService creates a new weather service
func NewWeatherService(cf gox.CrossFunction) *WeatherService {
	return &WeatherService{
		logger: cf.Logger(),
		metrics: cf.Metrics(),
	}
}

// GetForecast returns a weather forecast for a city
func (s *WeatherService) GetForecast(city string, days int) (string, error) {
	// Log the request
	s.logger.Info("Getting forecast", 
		zap.String("city", city), 
		zap.Int("days", days))
	
	// Record a metric
	s.metrics.Counter("weather_forecast_requests").Inc(1)
	
	// In a real application, this would call an external API
	return fmt.Sprintf("Weather forecast for %s for the next %d days: Sunny with a chance of clouds", city, days), nil
}

// ToolRegistrar registers MCP tools
type ToolRegistrar struct {
	server mcpApi.Server
	weatherService *WeatherService
	logger *zap.Logger
}

// NewToolRegistrar creates a new tool registrar
func NewToolRegistrar(
	server mcpApi.Server,
	weatherService *WeatherService,
	cf gox.CrossFunction,
) *ToolRegistrar {
	return &ToolRegistrar{
		server: server,
		weatherService: weatherService,
		logger: cf.Logger(),
	}
}

// Register registers all tools
func (r *ToolRegistrar) Register() error {
	// Register the weather tool
	weatherTool := mcp.NewTool("get_forecast",
		mcp.WithDescription("Get weather forecast for a city"),
		mcp.WithString("city", mcp.Required(), mcp.Description("City name")),
		mcp.WithInteger("days", mcp.Required(), mcp.Description("Number of days for forecast")),
	)
	
	err := goxMcp.RegisterTool[map[string]interface{}](r.server, weatherTool, 
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Parse the request parameters
			params, ok := request.Params.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid parameters format")
			}
			
			city, ok := params["city"].(string)
			if !ok {
				return nil, fmt.Errorf("city must be a string")
			}
			
			daysFloat, ok := params["days"].(float64)
			if !ok {
				return nil, fmt.Errorf("days must be a number")
			}
			days := int(daysFloat)
			
			// Get the forecast from the weather service
			forecast, err := r.weatherService.GetForecast(city, days)
			if err != nil {
				r.logger.Error("Failed to get forecast", 
					zap.String("city", city), 
					zap.Int("days", days),
					zap.Error(err))
				return nil, err
			}
			
			// Return the forecast as a text result
			return mcp.NewToolResultText(forecast), nil
		},
	)
	
	if err != nil {
		return err
	}
	
	r.logger.Info("Registered weather tool")
	return nil
}

func main() {
	// Create an FX application
	app := fx.New(
		// Provide the cross function
		fx.Provide(gox.NewCrossFunction),
		
		// Provide the application config
		fx.Supply(&AppConfig{
			MCP: mcpApi.Config{
				Disabled: false,
				Port:     8089,
			},
		}),
		
		// Extract the MCP config from the app config
		fx.Provide(
			func(config *AppConfig) *mcpApi.Config {
				return &config.MCP
			},
		),
		
		// Provide the MCP server
		goxMcp.Provider,
		
		// Register the server lifecycle hooks
		fx.Invoke(goxMcp.NewServerLifecycleInvoker),
		
		// Provide the weather service
		fx.Provide(NewWeatherService),
		
		// Provide the tool registrar
		fx.Provide(NewToolRegistrar),
		
		// Register the tools
		fx.Invoke(func(registrar *ToolRegistrar) error {
			return registrar.Register()
		}),
	)
	
	// Start the application
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	
	// Wait for signals to shut down
	<-ctx.Done()
	
	// Stop the application
	if err := app.Stop(ctx); err != nil {
		panic(err)
	}
}
```

## Project Structure

The project is organized into the following packages:

- `api`: Core interfaces and types used throughout the project
  - `Server`: The main interface for MCP servers
  - `Config`: Configuration for MCP servers
  - JSON-RPC request and response types

- `mcp-metoro-io`: Implementation of the MCP server using the mark3labs/mcp-go library
  - Server implementation
  - Tool registration
  - Integration with Uber fx

- `examples`: Example applications demonstrating how to use the library
  - `hello`: A simple example printing "Hello, World!"
  - `test_server`: An example MCP server with a hello tool

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/devlibx/gox-mcp.git
   cd gox-mcp
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
