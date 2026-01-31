package claudesdk

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	internalmcp "github.com/Savid/claude-agent-sdk-go/internal/mcp"
)

// Re-export MCP SDK types for public API.
// These are the official MCP protocol types.
type (
	// CallToolResult is the server's response to a tool call.
	// Use TextResult, ErrorResult, or ImageResult helpers to create results.
	CallToolResult = mcp.CallToolResult

	// CallToolRequest is the request passed to tool handlers.
	CallToolRequest = mcp.CallToolRequest

	// McpContent is the interface for content types in tool results.
	McpContent = mcp.Content

	// McpTextContent represents text content in a tool result.
	McpTextContent = mcp.TextContent

	// McpImageContent represents image content in a tool result.
	McpImageContent = mcp.ImageContent

	// McpAudioContent represents audio content in a tool result.
	McpAudioContent = mcp.AudioContent

	// McpTool represents an MCP tool definition from the official SDK.
	McpTool = mcp.Tool

	// McpToolHandler is the function signature for low-level tool handlers.
	McpToolHandler = mcp.ToolHandler

	// Schema is a JSON Schema object for tool input validation.
	Schema = jsonschema.Schema
)

// SdkMcpToolHandler is the function signature for SdkMcpTool handlers.
// It receives the context and request, and returns the result.
//
// Use ParseArguments to extract input as map[string]any from the request.
// Use TextResult, ErrorResult, or ImageResult helpers to create results.
//
// Example:
//
//	func(ctx context.Context, req *claudesdk.CallToolRequest) (*claudesdk.CallToolResult, error) {
//	    args, err := claudesdk.ParseArguments(req)
//	    if err != nil {
//	        return claudesdk.ErrorResult(err.Error()), nil
//	    }
//	    a := args["a"].(float64)
//	    return claudesdk.TextResult(fmt.Sprintf("Result: %v", a)), nil
//	}
type SdkMcpToolHandler = mcp.ToolHandler

// SdkMcpTool represents a tool created with NewSdkMcpTool.
type SdkMcpTool struct {
	ToolName        string
	ToolDescription string
	ToolSchema      *jsonschema.Schema
	ToolHandler     SdkMcpToolHandler
}

// Name returns the tool name.
func (t *SdkMcpTool) Name() string {
	return t.ToolName
}

// Description returns the tool description.
func (t *SdkMcpTool) Description() string {
	return t.ToolDescription
}

// InputSchema returns the JSON Schema for the tool input.
func (t *SdkMcpTool) InputSchema() *jsonschema.Schema {
	return t.ToolSchema
}

// Handler returns the tool handler function.
func (t *SdkMcpTool) Handler() SdkMcpToolHandler {
	return t.ToolHandler
}

// NewSdkMcpTool creates an SdkMcpTool.
//
// The inputSchema should be a *jsonschema.Schema. Use SimpleSchema for convenience
// or create a full Schema struct for more control.
//
// Example with SimpleSchema:
//
//	addTool := claudesdk.NewSdkMcpTool("add", "Add two numbers",
//	    claudesdk.SimpleSchema(map[string]string{"a": "float64", "b": "float64"}),
//	    func(ctx context.Context, req *claudesdk.CallToolRequest) (*claudesdk.CallToolResult, error) {
//	        args, _ := claudesdk.ParseArguments(req)
//	        a, b := args["a"].(float64), args["b"].(float64)
//	        return claudesdk.TextResult(fmt.Sprintf("Result: %v", a+b)), nil
//	    },
//	)
func NewSdkMcpTool(
	name, description string,
	inputSchema *jsonschema.Schema,
	handler SdkMcpToolHandler,
) *SdkMcpTool {
	return &SdkMcpTool{
		ToolName:        name,
		ToolDescription: description,
		ToolSchema:      inputSchema,
		ToolHandler:     handler,
	}
}

// SimpleSchema creates a jsonschema.Schema from a simple type map.
//
// Input format: {"a": "float64", "b": "string"}
//
// Type mappings:
//   - "string"           → {"type": "string"}
//   - "int", "int64"     → {"type": "integer"}
//   - "float64", "float" → {"type": "number"}
//   - "bool"             → {"type": "boolean"}
//   - "[]string"         → {"type": "array", "items": {"type": "string"}}
//   - "any", "object"    → {"type": "object"}
func SimpleSchema(props map[string]string) *jsonschema.Schema {
	return internalmcp.SimpleSchema(props)
}

// TextResult creates a CallToolResult with text content.
func TextResult(text string) *mcp.CallToolResult {
	return internalmcp.TextResult(text)
}

// ErrorResult creates a CallToolResult indicating an error.
func ErrorResult(message string) *mcp.CallToolResult {
	return internalmcp.ErrorResult(message)
}

// ImageResult creates a CallToolResult with image content.
func ImageResult(data []byte, mimeType string) *mcp.CallToolResult {
	return internalmcp.ImageResult(data, mimeType)
}

// ParseArguments unmarshals CallToolRequest arguments into a map.
// This is a convenience function for extracting tool input.
func ParseArguments(req *mcp.CallToolRequest) (map[string]any, error) {
	return internalmcp.ParseArguments(req)
}

// NewMcpTool creates an mcp.Tool with the given parameters.
// This is useful when you need direct access to the MCP Tool type.
func NewMcpTool(name, description string, inputSchema *jsonschema.Schema) *mcp.Tool {
	return internalmcp.NewTool(name, description, inputSchema)
}
