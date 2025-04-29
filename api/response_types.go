// Package api provides interfaces and types for the Model Context Protocol (MCP) server.
package mcpApi

import "encoding/json"

// McpToolListRequest represents a JSON-RPC 2.0 request for listing available MCP tools.
// This structure follows the JSON-RPC 2.0 specification format for method calls.
type McpToolListRequest struct {
	// JSONRPC specifies the version of the JSON-RPC protocol.
	// Must be exactly "2.0" for JSON-RPC 2.0.
	JSONRPC string `json:"jsonrpc"`

	// ID is the identifier for this request, which will be echoed back in the response.
	// It can be a string or number.
	ID interface{} `json:"id"`

	// Method specifies the RPC method to invoke.
	// For tool listing, this is "tools/list".
	Method string `json:"method"`
}

// NewMcpToolListRequest creates a new McpToolListRequest with default values.
// This function provides a convenient way to create a properly formatted
// JSON-RPC request for listing MCP tools.
//
// Returns:
//   - *McpToolListRequest: A new request object with default values
func NewMcpToolListRequest() *McpToolListRequest {
	return &McpToolListRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	}
}

// ToJSON serializes the McpToolListRequest to a JSON byte array.
// This is useful for sending the request over HTTP or other transport mechanisms.
//
// Returns:
//   - []byte: The JSON representation of the request
//   - error: An error if serialization fails, or nil on success
func (r *McpToolListRequest) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// McpTollListResponse represents a JSON-RPC 2.0 response containing a list of available MCP tools.
// This structure follows the JSON-RPC 2.0 specification format and includes specific fields
// for MCP tool information.
type McpTollListResponse struct {
	// ID is the identifier established by the request.
	// It can be a string or number.
	ID interface{} `json:"id"`

	// JSONRPC specifies the version of the JSON-RPC protocol.
	// Must be exactly "2.0" for JSON-RPC 2.0.
	JSONRPC string `json:"jsonrpc"`

	// Result contains the data returned by the method invocation.
	// For tool listing, this includes an array of available tools.
	Result *ToolListResult `json:"result,omitempty"`

	// Error contains error information if the call was not successful.
	// This field is omitted when the call was successful.
	Error *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents an error in a JSON-RPC 2.0 response.
// This is included in the response when an error occurs.
type JSONRPCError struct {
	// Code is a number indicating the error type.
	Code int `json:"code"`

	// Message is a short description of the error.
	Message string `json:"message"`

	// Data is optional additional information about the error.
	Data interface{} `json:"data,omitempty"`
}

// ToolListResult contains the successful result data of a tool listing request.
// It includes an array of available tools provided by the MCP server.
type ToolListResult struct {
	// Tools is an array of available tools provided by the MCP server.
	Tools []Tool `json:"tools"`
}

// Tool represents a tool available in the MCP server.
// Each tool has a name, description, and an input schema that defines
// the expected parameters.
type Tool struct {
	// Name is the identifier of the tool, used when calling the tool.
	Name string `json:"name"`

	// Description provides information about what the tool does.
	Description string `json:"description"`

	// InputSchema defines the expected structure of input parameters
	// for the tool, using JSON Schema.
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema defines the structure of input parameters for a tool
// using JSON Schema format.
type InputSchema struct {
	// Schema is the JSON Schema version identifier.
	Schema string `json:"$schema"`

	// Properties defines the fields accepted by the tool.
	Properties map[string]PropertyDefinition `json:"properties"`

	// Required is an array of property names that are required.
	Required []string `json:"required"`

	// Type specifies the JSON Schema type, typically "object" for tool inputs.
	Type string `json:"type"`
}

// PropertyDefinition defines a single property in a JSON Schema.
// It describes the type and other constraints for a field.
type PropertyDefinition struct {
	// Description provides information about the property's purpose.
	Description string `json:"description"`

	// Type specifies the JSON data type of the property.
	Type string `json:"type"`
}
