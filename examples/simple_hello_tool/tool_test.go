package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base/v2"
	goxJsonUtils "github.com/devlibx/gox-base/v2/serialization/utils/json"
	"github.com/go-resty/resty/v2"
	mcpGo "github.com/metoro-io/mcp-golang"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"goc-mcp/api"
	goxMcp "goc-mcp/mcp-metoro-io"
	"testing"
	"time"
)

func TestTool(t *testing.T) {

	// Create a app
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
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)

	restyClient := resty.New()
	restyClient.SetBaseURL(fmt.Sprintf("http://localhost:%d/mcp", 1234))
	restyClient.Debug = true
	resp, err := restyClient.R().SetBody(mcpApi.NewMcpToolListRequest()).Post("")
	assert.NoError(t, err)
	mapTollListRep, err := goxJsonUtils.BytesToObject[*mcpApi.McpTollListResponse](resp.Body())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mapTollListRep.Result.Tools))
	assert.Equal(t, "prompt_tool", mapTollListRep.Result.Tools[0].Name)
}
