package main

import (
	"context"
	mcpApi "github.com/devlibx/goc-mcp/api"
	goxMcp "github.com/devlibx/goc-mcp/mcp-metoro-io"
	"github.com/devlibx/gox-base/v2"
	mcpGo "github.com/metoro-io/mcp-golang"
	"go.uber.org/fx"
	"time"
)

func main() {
	var server mcpApi.Server
	app := fx.New(
		fx.Provide(gox.NewNoOpCrossFunction),
		fx.Supply(&mcpApi.Config{
			Disabled: false,
			Port:     1234,
		}),
		goxMcp.Provider,
		fx.Invoke(goxMcp.NewServerLifecycleInvoker),
		fx.Populate(&server),
	)
	err := app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	type Content struct {
		Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
		Description *string `json:"description" jsonschema:"description=The description to submit"`
	}
	err = goxMcp.RegisterTool[Content](
		server,
		"prompt_tool",
		"Prompt Tool",
		func(input Content) (*mcpGo.ToolResponse, error) {
			return mcpGo.NewToolResponse(
				&mcpGo.Content{
					Type:             "",
					TextContent:      nil,
					ImageContent:     nil,
					EmbeddedResource: nil,
					Annotations:      nil,
				},
			), nil
		},
	)
	time.Sleep(1 * time.Hour)
}
