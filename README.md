# First Resonance MCP Server

> **Note:** The First Resonance MCP Server is currently a simple prototype and may be subject to significant changes in future releases. It is provided as a proof of concept for integrating with the First Resonance API.

The First Resonance MCP Server is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction)
server that provides seamless integration with First Resonance APIs, enabling advanced
automation and interaction capabilities for developers and tools.

[![Install with Docker in VS Code](https://img.shields.io/badge/VS_Code-Install_Server-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=firstresonance&inputs=%5B%7B%22id%22%3A%22firstresonance_token%22%2C%22type%22%3A%22promptString%22%2C%22description%22%3A%22First%20Resonance%20API%20Token%22%2C%22password%22%3Atrue%7D%5D&config=%7B%22command%22%3A%22docker%22%2C%22args%22%3A%5B%22run%22%2C%22-i%22%2C%22--rm%22%2C%22-e%22%2C%22FIRSTRESONANCE_API_TOKEN%22%2C%22ghcr.io%2Ffirstresonance%2Ffirstresonance-mcp-server%22%5D%2C%22env%22%3A%7B%22FIRSTRESONANCE_API_TOKEN%22%3A%22%24%7Binput%3Afirstresonance_token%7D%22%7D%7D) [![Install with Docker in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install_Server-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=firstresonance&inputs=%5B%7B%22id%22%3A%22firstresonance_token%22%2C%22type%22%3A%22promptString%22%2C%22description%22%3A%22First%20Resonance%20API%20Token%22%2C%22password%22%3Atrue%7D%5D&config=%7B%22command%22%3A%22docker%22%2C%22args%22%3A%5B%22run%22%2C%22-i%22%2C%22--rm%22%2C%22-e%22%2C%22FIRSTRESONANCE_API_TOKEN%22%2C%22ghcr.io%2Ffirstresonance%2Ffirstresonance-mcp-server%22%5D%2C%22env%22%3A%7B%22FIRSTRESONANCE_API_TOKEN%22%3A%22%24%7Binput%3Afirstresonance_token%7D%22%7D%7D&quality=insiders)

## Use Cases

- Automating First Resonance workflows and processes.
- Extracting and analyzing data from First Resonance.
- Building AI powered tools and applications that interact with First Resonance's ecosystem.

## Prerequisites

1. To run the server in a container, you will need to have [Docker](https://www.docker.com/) installed.
2. Once Docker is installed, you will also need to ensure Docker is running.
3. Lastly you will need to [Create a First Resonance API Token](https://app.firstresonance.io/settings/api-tokens).
The MCP server can use many of the First Resonance APIs, so enable the permissions that you feel comfortable granting your AI tools.

## Installation

### Usage with VS Code

For quick installation, use one of the one-click install buttons at the top of this README.

For manual installation, add the following JSON block to your User Settings (JSON) file in VS Code. You can do this by pressing `Ctrl + Shift + P` and typing `Preferences: Open User Settings (JSON)`.

Optionally, you can add it to a file called `.vscode/mcp.json` in your workspace. This will allow you to share the configuration with others.

> Note that the `mcp` key is not needed in the `.vscode/mcp.json` file.

```json
{
  "mcp": {
    "inputs": [
      {
        "type": "promptString",
        "id": "firstresonance_token",
        "description": "First Resonance API Token",
        "password": true
      }
    ],
    "servers": {
      "firstresonance": {
        "command": "docker",
        "args": [
          "run",
          "-i",
          "--rm",
          "-e",
          "FIRSTRESONANCE_API_TOKEN",
          "ghcr.io/firstresonance/firstresonance-mcp-server"
        ],
        "env": {
          "FIRSTRESONANCE_API_TOKEN": "${input:firstresonance_token}"
        }
      }
    }
  }
}
```

More about using MCP server tools in VS Code's [agent mode documentation](https://code.visualstudio.com/docs/copilot/chat/mcp-servers).

### Usage with Claude Desktop

```json
{
  "mcpServers": {
    "firstresonance": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "FIRSTRESONANCE_API_TOKEN",
        "ghcr.io/firstresonance/firstresonance-mcp-server"
      ],
      "env": {
        "FIRSTRESONANCE_API_TOKEN": "<YOUR_TOKEN>"
      }
    }
  }
}
```

### Build from source

If you don't have Docker, you can use `go` to build the binary in the
`cmd/firstresonance-mcp-server` directory, and use the `firstresonance-mcp-server stdio`
command with the `FIRSTRESONANCE_API_TOKEN` environment variable set to
your token.

## First Resonance Enterprise Server

The flag `--fr-host` and the environment variable `FR_HOST` can be used to set
the First Resonance Enterprise Server hostname.

## i18n / Overriding Descriptions

The descriptions of the tools can be overridden by creating a
`firstresonance-mcp-server-config.json` file in the same directory as the binary.

The file should contain a JSON object with the tool names as keys and the new
descriptions as values. For example:

```json
{
  "TOOL_GET_PART_DESCRIPTION": "an alternative description",
  "TOOL_CREATE_PART_DESCRIPTION": "Create a new part in First Resonance"
}
```

You can create an export of the current translations by running the binary with
the `--export-translations` flag.

This flag will preserve any translations/overrides you have made, while adding
any new translations that have been added to the binary since the last time you
exported.

```sh
./firstresonance-mcp-server --export-translations
cat firstresonance-mcp-server-config.json
```

You can also use ENV vars to override the descriptions. The environment
variable names are the same as the keys in the JSON file, prefixed with
`FIRSTRESONANCE_MCP_` and all uppercase.

For example, to override the `TOOL_GET_PART_DESCRIPTION` tool, you can
set the following environment variable:

```sh
export FIRSTRESONANCE_MCP_TOOL_GET_PART_DESCRIPTION="an alternative description"
```

## Tools

### Users

- **get_me** - Get details of the authenticated user
  - No parameters required

### Parts

- **get_part** - Gets the details of a part within First Resonance

  - `part_id`: Part ID (string, required)

- **list_parts** - List and filter parts

  - `status`: Filter by status (string, optional)
  - `type`: Filter by type (string, optional)
  - `page`: Page number (number, optional)
  - `perPage`: Results per page (number, optional)

- **create_part** - Create a new part in First Resonance

  - `name`: Part name (string, required)
  - `description`: Part description (string, optional)
  - `type`: Part type (string, required)
  - `status`: Part status (string, optional)

- **update_part** - Update an existing part in First Resonance

  - `part_id`: Part ID to update (string, required)
  - `name`: New name (string, optional)
  - `description`: New description (string, optional)
  - `type`: New type (string, optional)
  - `status`: New status (string, optional)

### Orders

- **get_order** - Get details of a specific order

  - `order_id`: Order ID (string, required)

- **list_orders** - List and filter orders

  - `status`: Order status (string, optional)
  - `sort`: Sort field (string, optional)
  - `direction`: Sort direction (string, optional)
  - `perPage`: Results per page (number, optional)
  - `page`: Page number (number, optional)

- **create_order** - Create a new order

  - `customer_id`: Customer ID (string, required)
  - `items`: Order items (array, required)
  - `priority`: Order priority (string, optional)
  - `due_date`: Due date (string, optional)

- **update_order** - Update an existing order

  - `order_id`: Order ID to update (string, required)
  - `status`: New status (string, optional)
  - `priority`: New priority (string, optional)
  - `due_date`: New due date (string, optional)

### Suppliers

- **get_supplier** - Get details of a specific supplier

  - `supplier_id`: Supplier ID (string, required)

- **list_suppliers** - List and filter suppliers

  - `status`: Supplier status (string, optional)
  - `sort`: Sort field (string, optional)
  - `direction`: Sort direction (string, optional)
  - `perPage`: Results per page (number, optional)
  - `page`: Page number (number, optional)

- **create_supplier** - Create a new supplier

  - `name`: Supplier name (string, required)
  - `contact_info`: Contact information (object, optional)
  - `status`: Supplier status (string, optional)

- **update_supplier** - Update an existing supplier

  - `supplier_id`: Supplier ID to update (string, required)
  - `name`: New name (string, optional)
  - `contact_info`: New contact information (object, optional)
  - `status`: New status (string, optional)

### Inventory

- **get_inventory_item** - Get details of a specific inventory item

  - `item_id`: Inventory item ID (string, required)

- **list_inventory_items** - List and filter inventory items

  - `location`: Filter by location (string, optional)
  - `status`: Filter by status (string, optional)
  - `sort`: Sort field (string, optional)
  - `direction`: Sort direction (string, optional)
  - `perPage`: Results per page (number, optional)
  - `page`: Page number (number, optional)

- **update_inventory_item** - Update an existing inventory item

  - `item_id`: Inventory item ID to update (string, required)
  - `quantity`: New quantity (number, optional)
  - `location`: New location (string, optional)
  - `status`: New status (string, optional)

### Search

- **search_parts** - Search for parts across First Resonance

  - `query`: Search query (string, required)
  - `sort`: Sort field (string, optional)
  - `order`: Sort order (string, optional)
  - `page`: Page number (number, optional)
  - `perPage`: Results per page (number, optional)

- **search_orders** - Search for orders
  - `query`: Search query (string, required)
  - `sort`: Sort field (string, optional)
  - `order`: Sort order (string, optional)
  - `page`: Page number (number, optional)
  - `perPage`: Results per page (number, optional)

## Resources

### First Resonance Content

- **Get Part Content**
  Retrieves the content of a part.

  - **Template**: `part://{part_id}`
  - **Parameters**:
    - `part_id`: Part ID (string, required)

- **Get Order Content**
  Retrieves the content of an order.

  - **Template**: `order://{order_id}`
  - **Parameters**:
    - `order_id`: Order ID (string, required)

- **Get Supplier Content**
  Retrieves the content of a supplier.

  - **Template**: `supplier://{supplier_id}`
  - **Parameters**:
    - `supplier_id`: Supplier ID (string, required)

- **Get Inventory Item Content**
  Retrieves the content of an inventory item.

  - **Template**: `inventory://{item_id}`
  - **Parameters**:
    - `item_id`: Inventory item ID (string, required)

## Library Usage

The exported Go API of this module should currently be considered unstable, and subject to breaking changes. In the future, we may offer stability; please file an issue if there is a use case where this would be valuable.

## Microservice Integration

The First Resonance MCP Server can be deployed as a service within your microservice ecosystem. This section explains how to interact with the service in this context.

### Service Configuration

When deploying the MCP Server as a service, you'll need to configure it with the following environment variables:

```yaml
environment:
  - FIRSTRESONANCE_API_TOKEN=${FIRSTRESONANCE_API_TOKEN}
  - FR_HOST=${FR_HOST}  # Optional: For First Resonance Enterprise Server
  - PORT=8080  # Optional: Default port for the service
```

### API Endpoints

The MCP Server exposes a REST API that can be used by other services in your ecosystem:

#### Authentication

All requests to the MCP Server must include an API key in the `Authorization` header:

```
Authorization: Bearer ${MCP_SERVER_API_KEY}
```

#### Tool Execution

To execute a tool, send a POST request to the `/tools/{tool_name}` endpoint:

```
POST /tools/get_part
Content-Type: application/json

{
  "part_id": "part-123"
}
```

The response will be in the following format:

```json
{
  "result": {
    "id": "part-123",
    "name": "Example Part",
    "description": "This is an example part",
    "type": "component",
    "status": "active"
  },
  "error": null
}
```

#### Batch Tool Execution

For efficiency, you can execute multiple tools in a single request:

```
POST /tools/batch
Content-Type: application/json

{
  "tools": [
    {
      "name": "get_part",
      "params": {
        "part_id": "part-123"
      }
    },
    {
      "name": "list_inventory_items",
      "params": {
        "location": "warehouse-1"
      }
    }
  ]
}
```

The response will contain results for all tools in the same order:

```json
{
  "results": [
    {
      "result": {
        "id": "part-123",
        "name": "Example Part",
        "description": "This is an example part",
        "type": "component",
        "status": "active"
      },
      "error": null
    },
    {
      "result": {
        "items": [
          {
            "id": "inv-456",
            "part_id": "part-123",
            "quantity": 10,
            "location": "warehouse-1"
          }
        ]
      },
      "error": null
    }
  ]
}
```

### Service Discovery

The MCP Server registers itself with your service discovery system (e.g., Consul, etcd) using the following metadata:

```json
{
  "name": "firstresonance-mcp-server",
  "version": "1.0.0",
  "tags": ["mcp", "firstresonance", "api"],
  "health": "/health"
}
```

### Health Checks

The MCP Server provides a health check endpoint at `/health` that returns:

```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": "24h30m15s"
}
```

### Rate Limiting

The MCP Server implements rate limiting to prevent abuse. By default, it allows:

- 100 requests per minute per API key
- 1000 requests per hour per API key

Rate limit headers are included in all responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1620000000
```

### Error Handling

All errors follow a consistent format:

```json
{
  "error": {
    "code": "TOOL_EXECUTION_ERROR",
    "message": "Failed to execute tool: get_part",
    "details": {
      "reason": "Part not found",
      "part_id": "part-123"
    }
  }
}
```

Common error codes:
- `INVALID_REQUEST`: The request format is invalid
- `UNAUTHORIZED`: The API key is invalid or missing
- `RATE_LIMIT_EXCEEDED`: The rate limit has been exceeded
- `TOOL_EXECUTION_ERROR`: The tool execution failed
- `INTERNAL_ERROR`: An internal server error occurred

### Client Libraries

For easier integration, client libraries are available in the following languages:

- Go: `github.com/firstresonance/firstresonance-mcp-client-go`
- Python: `firstresonance-mcp-client-python`
- JavaScript: `@firstresonance/mcp-client-js`

Example usage with the Go client:

```go
import (
  "github.com/firstresonance/firstresonance-mcp-client-go"
)

client := mcp.NewClient("http://firstresonance-mcp-server:8080", "your-api-key")

// Execute a single tool
result, err := client.ExecuteTool("get_part", map[string]interface{}{
  "part_id": "part-123",
})

// Execute multiple tools in batch
results, err := client.ExecuteBatch([]mcp.ToolRequest{
  {
    Name: "get_part",
    Params: map[string]interface{}{
      "part_id": "part-123",
    },
  },
  {
    Name: "list_inventory_items",
    Params: map[string]interface{}{
      "location": "warehouse-1",
    },
  },
})
```

## License

This project is licensed under the terms of the MIT open source license. Please refer to [MIT](./LICENSE) for the full terms.

## AI Model Handler

The First Resonance MCP Server includes an AI Model handler that manages requests from AI models. This handler provides a secure and controlled way for AI models to interact with the First Resonance API.

> **Note:** The AI Model handler is currently a simple prototype and may be subject to significant changes in future releases. It is provided as a proof of concept for integrating AI models with the First Resonance API.

### Features

- **Model Registration**: AI models must be registered before they can make requests
- **Context Management**: Supports context inheritance for maintaining state across requests
- **Action-Based Interface**: Handles specific actions like `get_part` and `create_order`
- **Type-Safe Parameter Handling**: Validates and converts parameters to appropriate types
- **Error Handling**: Provides detailed error messages for invalid requests

### Available Actions

#### get_part
Retrieves details of a specific part.

Parameters:
- `part_id` (string, required): The ID of the part to retrieve

#### create_order
Creates a new order in First Resonance.

Parameters:
- `customer_id` (string, required): The ID of the customer
- `items` (array, required): Array of order items
- `priority` (string, optional): Order priority
- `due_date` (string, optional): Due date for the order
- `status` (string, optional): Initial status of the order

### Example Usage

```go
// Create a new model handler
handler := aimodel.NewModelHandler(contextManager, frClient)

// Register a model
handler.RegisterModel("model-123")

// Handle a request
req := &aimodel.ModelRequest{
    ModelID:   "model-123",
    Action:    "get_part",
    ContextID: "ctx-456",
    Params: map[string]interface{}{
        "part_id": "part-789",
    },
}

resp, err := handler.HandleRequest(ctx, req)
```

### Security Considerations

- Models must be explicitly registered before they can make requests
- Each request is validated against the registered model list
- Context inheritance is controlled and audited
- All parameters are type-checked before processing
- Error messages are sanitized to prevent information leakage
