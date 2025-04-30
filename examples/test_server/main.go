package main

import (
	"context"
	"errors"
	"github.com/devlibx/gox-base/v2"
	goxJsonUtils "github.com/devlibx/gox-base/v2/serialization/utils/json"
	mcpApi "github.com/devlibx/gox-mcp/api"
	goxMcp "github.com/devlibx/gox-mcp/mcp-metoro-io"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/fx"
	"time"
)

func main() {

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

	type content struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
	}

	tool := mcp.NewTool("hello_tool",
		mcp.WithDescription("Say hello to user"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the person to greet")),
	)

	err = goxMcp.RegisterTool[content](server, tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if paramsStr, err := goxJsonUtils.ObjectToString(request.Params); err == nil {
			if obj, err := goxJsonUtils.StringToObject[content](paramsStr); err == nil {
				return mcp.NewToolResultText("this is name " + obj.Title + " "), err
			}
		}
		return nil, errors.New("name must be a string")
	})
	time.Sleep(1 * time.Hour)
}
