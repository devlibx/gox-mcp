package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer("test", "1.0.0")
	sseServer := server.NewSSEServer(mcpServer,
		server.WithBaseURL("http://localhost:8089"),
	)

	// Add tool
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	// Add tool handler
	// s.AddTool(tool, helloHandler)
	mcpServer.AddTool(tool, helloHandler)

	if err := sseServer.Start(":8089"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
	fmt.Println("Server started at http://localhost:8089/mcp")
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
