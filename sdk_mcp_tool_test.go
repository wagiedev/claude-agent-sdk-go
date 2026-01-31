package claudesdk

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextResult(t *testing.T) {
	result := TextResult("Hello, World!")

	assert.Len(t, result.Content, 1)
	assert.False(t, result.IsError)

	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "Hello, World!", textContent.Text)
}

func TestErrorResult(t *testing.T) {
	result := ErrorResult("Something went wrong")

	assert.Len(t, result.Content, 1)
	assert.True(t, result.IsError)

	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "Something went wrong", textContent.Text)
}

func TestImageResult(t *testing.T) {
	result := ImageResult([]byte("base64data"), "image/png")

	assert.Len(t, result.Content, 1)
	assert.False(t, result.IsError)

	imageContent, ok := result.Content[0].(*mcp.ImageContent)
	require.True(t, ok)
	assert.Equal(t, []byte("base64data"), imageContent.Data)
	assert.Equal(t, "image/png", imageContent.MIMEType)
}

func TestSdkMcpTool(t *testing.T) {
	t.Run("has name and description", func(t *testing.T) {
		tool := NewSdkMcpTool(
			"test_tool",
			"A test tool",
			SimpleSchema(map[string]string{"value": "string"}),
			func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return TextResult("ok"), nil
			},
		)

		assert.Equal(t, "test_tool", tool.Name())
		assert.Equal(t, "A test tool", tool.Description())

		schema := tool.InputSchema()
		assert.NotNil(t, schema)
		assert.Equal(t, "object", schema.Type)
		assert.Contains(t, schema.Properties, "value")
	})

	t.Run("handler executes correctly", func(t *testing.T) {
		tool := NewSdkMcpTool(
			"adder",
			"Adds numbers",
			SimpleSchema(map[string]string{"a": "float64", "b": "float64"}),
			func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				args, err := ParseArguments(req)
				if err != nil {
					return ErrorResult(err.Error()), nil
				}

				_, _ = args["a"].(float64)
				_, _ = args["b"].(float64)

				return TextResult("Result: sum"), nil
			},
		)

		// Create a mock request
		inputJSON, _ := json.Marshal(map[string]any{"a": 1.0, "b": 2.0})
		req := &mcp.CallToolRequest{
			Params: &mcp.CallToolParamsRaw{
				Name:      "adder",
				Arguments: inputJSON,
			},
		}

		result, err := tool.Handler()(context.Background(), req)
		require.NoError(t, err)
		assert.Len(t, result.Content, 1)

		textContent, ok := result.Content[0].(*mcp.TextContent)
		require.True(t, ok)
		assert.Equal(t, "Result: sum", textContent.Text)
	})
}

func TestSimpleSchema(t *testing.T) {
	t.Run("converts simple type map to JSON Schema", func(t *testing.T) {
		schema := SimpleSchema(map[string]string{
			"name":  "string",
			"count": "int",
			"value": "float64",
			"flag":  "bool",
		})

		assert.Equal(t, "object", schema.Type)
		assert.Len(t, schema.Properties, 4)
		assert.Len(t, schema.Required, 4)

		assert.Equal(t, "string", schema.Properties["name"].Type)
		assert.Equal(t, "integer", schema.Properties["count"].Type)
		assert.Equal(t, "number", schema.Properties["value"].Type)
		assert.Equal(t, "boolean", schema.Properties["flag"].Type)
	})

	t.Run("handles array types", func(t *testing.T) {
		schema := SimpleSchema(map[string]string{
			"items": "[]string",
		})

		assert.Equal(t, "array", schema.Properties["items"].Type)
		assert.Equal(t, "string", schema.Properties["items"].Items.Type)
	})
}

func TestParseArguments(t *testing.T) {
	t.Run("parses valid JSON arguments", func(t *testing.T) {
		inputJSON, _ := json.Marshal(map[string]any{"a": 1.0, "b": "hello"})
		req := &mcp.CallToolRequest{
			Params: &mcp.CallToolParamsRaw{
				Name:      "test",
				Arguments: inputJSON,
			},
		}

		args, err := ParseArguments(req)
		require.NoError(t, err)
		assert.Equal(t, 1.0, args["a"])
		assert.Equal(t, "hello", args["b"])
	})

	t.Run("handles nil request", func(t *testing.T) {
		args, err := ParseArguments(nil)
		require.NoError(t, err)
		assert.Empty(t, args)
	})

	t.Run("handles empty arguments", func(t *testing.T) {
		req := &mcp.CallToolRequest{
			Params: &mcp.CallToolParamsRaw{
				Name:      "test",
				Arguments: nil,
			},
		}

		args, err := ParseArguments(req)
		require.NoError(t, err)
		assert.Empty(t, args)
	})
}

func TestCreateSdkMcpServer(t *testing.T) {
	tool := NewSdkMcpTool(
		"test_tool",
		"A test tool",
		SimpleSchema(map[string]string{"value": "string"}),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return TextResult("ok"), nil
		},
	)

	config := CreateSdkMcpServer("test_server", "1.0.0", tool)

	assert.Equal(t, MCPServerTypeSDK, config.Type)
	assert.Equal(t, "test_server", config.Name)
	assert.NotNil(t, config.Instance)
}
