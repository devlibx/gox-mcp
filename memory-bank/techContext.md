# Technical Context

## Technology Stack

### Core Technologies
- **Go**: The project is written in Go (version 1.24.2), leveraging its strong typing, concurrency features, and performance characteristics.
- **JSON-RPC 2.0**: The communication protocol used for tool discovery and invocation.
- **HTTP**: The transport layer for exposing MCP servers.

### Key Dependencies
- **github.com/metoro-io/mcp-golang v0.11.0**: The core MCP implementation for Go.
- **go.uber.org/fx v1.23.0**: Dependency injection framework for managing component lifecycle.
- **github.com/devlibx/gox-base/v2 v2.0.27**: Utility library providing cross-cutting concerns.
- **github.com/go-resty/resty/v2 v2.7.0**: HTTP client for testing and interacting with MCP servers.
- **github.com/stretchr/testify v1.9.0**: Testing framework.

## Development Environment

### Requirements
- Go 1.24.2 or later
- Git for version control
- Standard Go development tools (go mod, go test, etc.)

### Setup
1. Clone the repository
2. Run `go mod download` to fetch dependencies
3. Use `go test ./...` to run tests

## Technical Constraints

### Performance Considerations
- MCP servers should handle requests with low latency
- Tool execution should be efficient to avoid timeouts
- Concurrent tool invocations should be supported

### Security Considerations
- Input validation using JSON Schema
- Proper error handling to avoid information leakage
- Transport security (HTTPS) for production deployments

### Compatibility
- Follows the MCP specification for interoperability
- Compatible with various Go application frameworks
- Works with different deployment environments (containers, VMs, etc.)

## Code Organization

### Package Structure
- **api/**: Core interfaces and types
- **mcp-metoro-io/**: Implementation of MCP server
- **examples/**: Example usage and test cases

### Design Principles
- Clear separation of concerns
- Interface-based design for testability
- Minimal dependencies
- Type safety through generics
- Idiomatic Go code

## Testing Strategy

### Unit Tests
- Test individual components in isolation
- Mock dependencies using interfaces
- Focus on edge cases and error handling

### Integration Tests
- Test the complete flow from client to server
- Verify tool registration and invocation
- Test JSON-RPC protocol compliance

### Example-based Tests
- Provide working examples that serve as both documentation and tests
- Cover common use cases and patterns

## Deployment Considerations

### Configuration
- Server port and other settings via configuration
- Environment-based configuration for different deployments

### Monitoring
- Log important events and errors
- Expose metrics for monitoring (future enhancement)

### Scaling
- Horizontal scaling through multiple server instances
- Load balancing for high-availability deployments