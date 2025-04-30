# MCP Server Example

This example demonstrates how to create a Model Context Protocol (MCP) server. The server exposes tools, prompts, and resources that can be accessed via HTTP and supports Server-Sent Events (SSE) by default.

## Features

- HTTP server with SSE support
- Two tools:
  - `prompt_tool`: Tool to greet user
  - `weather_tool`: Tool to get weather alerts for US states
- Greeting prompt for generating welcome messages

## Running the Server

To run the server:

```bash
go run main.go
```

The server will start on port 1234 and expose the SSE endpoint at `/sse`.

## Available Endpoints

- **Tools**: `http://localhost:1234/sse`
  - `prompt_tool`: Tool to greet user
  - `weather_tool`: Tool to get weather alerts for US states

- **Prompts**: `http://localhost:1234/sse`
  - `hello_prompt`: Returns a hello message

## Using the Tools

### Using the Greeting Prompt

You can use the hello prompt by sending a POST request to the SSE endpoint:

```bash
curl -X POST http://localhost:1234/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "prompts/use",
    "params": {
      "name": "hello_prompt",
      "arguments": {
        "title": "World",
        "description": "Testing the hello prompt"
      }
    }
  }'
```

### Using the Weather Alerts Tool

You can get weather alerts for a US state by sending a POST request to the SSE endpoint:

```bash
curl -X POST http://localhost:1234/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/use",
    "params": {
      "name": "get_weather_alerts",
      "arguments": {
        "state": "CA"
      }
    }
  }'
```

This will return weather alerts for California. You can replace "CA" with any other two-letter US state code (e.g., "NY", "TX").

## Listing Available Tools

To list all available tools:

```bash
curl -X POST http://localhost:1234/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/list"
  }'
```

## Using with SSE

To connect to the server using Server-Sent Events:

```javascript
// Browser JavaScript example
const eventSource = new EventSource('http://localhost:1234/sse');

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Received SSE message:', data);
};

eventSource.onerror = (error) => {
  console.error('SSE error:', error);
  eventSource.close();
};
```

## Stopping the Server

Press `Ctrl+C` to stop the server.
