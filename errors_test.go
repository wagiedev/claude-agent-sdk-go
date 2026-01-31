package claudesdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCLINotFoundError_Creation tests CLINotFoundError creation and formatting.
func TestCLINotFoundError_Creation(t *testing.T) {
	searchedPaths := []string{
		"$PATH",
		"/usr/local/bin/claude",
		"/usr/bin/claude",
	}
	err := &CLINotFoundError{
		SearchedPaths: searchedPaths,
	}

	require.Error(t, err)
	require.Contains(t, err.Error(), "claude CLI not found")
	require.Contains(t, err.Error(), "$PATH")
	require.Contains(t, err.Error(), "/usr/local/bin/claude")
}

// TestCLIConnectionError_Creation tests CLIConnectionError creation and formatting.
func TestCLIConnectionError_Creation(t *testing.T) {
	innerErr := fmt.Errorf("connection refused")
	err := &CLIConnectionError{
		Err: innerErr,
	}

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to connect to CLI")
	require.Contains(t, err.Error(), "connection refused")
}

// TestProcessError_WithExitCodeAndStderr tests ProcessError with exit code and stderr.
func TestProcessError_WithExitCodeAndStderr(t *testing.T) {
	err := &ProcessError{
		ExitCode: 1,
		Stderr:   "Error: authentication failed",
		Err:      nil,
	}

	require.Error(t, err)
	require.Contains(t, err.Error(), "CLI process failed")
	require.Contains(t, err.Error(), "exit 1")
	require.Contains(t, err.Error(), "authentication failed")
}

// TestMessageParseError_Creation tests MessageParseError creation and formatting.
func TestMessageParseError_Creation(t *testing.T) {
	innerErr := fmt.Errorf("invalid JSON")
	err := &MessageParseError{
		Message: `{"incomplete": `,
		Err:     innerErr,
		Data: map[string]any{
			"incomplete": true,
		},
	}

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse message")
	require.Contains(t, err.Error(), "invalid JSON")
}

// TestMessageParseError_PreservesMessage tests that MessageParseError preserves the original message.
func TestMessageParseError_PreservesMessage(t *testing.T) {
	err := &MessageParseError{
		Message: `{"type": "unknown", "data": 123}`,
		Err:     fmt.Errorf("unknown type"),
		Data: map[string]any{
			"type": "unknown",
			"data": 123,
		},
	}

	require.Equal(t, `{"type": "unknown", "data": 123}`, err.Message)
	require.Equal(t, "unknown", err.Data["type"])
	require.Equal(t, 123, err.Data["data"])
}

// TestCLIJSONDecodeError_Creation tests CLIJSONDecodeError creation and formatting.
func TestCLIJSONDecodeError_Creation(t *testing.T) {
	innerErr := fmt.Errorf("unexpected end of JSON input")
	err := &CLIJSONDecodeError{
		RawData: `{"incomplete": `,
		Err:     innerErr,
	}

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode JSON from CLI")
	require.Contains(t, err.Error(), "unexpected end of JSON input")
}

// TestCLIJSONDecodeError_PreservesRawData tests that raw data is preserved.
func TestCLIJSONDecodeError_PreservesRawData(t *testing.T) {
	rawData := `{"type": "user", invalid}`
	err := &CLIJSONDecodeError{
		RawData: rawData,
		Err:     fmt.Errorf("invalid character"),
	}

	require.Equal(t, rawData, err.RawData)
	require.Contains(t, err.Error(), "invalid character")
}

// TestCLIJSONDecodeError_Unwrap tests that the underlying error can be unwrapped.
func TestCLIJSONDecodeError_Unwrap(t *testing.T) {
	innerErr := fmt.Errorf("syntax error")
	err := &CLIJSONDecodeError{
		RawData: `{bad}`,
		Err:     innerErr,
	}

	require.ErrorIs(t, err, innerErr)
}

func TestAsType(t *testing.T) {
	t.Run("CLINotFoundError", func(t *testing.T) {
		paths := []string{"/usr/bin", "/usr/local/bin"}
		err := &CLINotFoundError{SearchedPaths: paths}

		result, ok := AsType[*CLINotFoundError](err)
		require.True(t, ok)
		require.NotNil(t, result)
		require.Equal(t, paths, result.SearchedPaths)

		_, ok = AsType[*CLINotFoundError](fmt.Errorf("other"))
		require.False(t, ok)

		_, ok = AsType[*CLINotFoundError](nil)
		require.False(t, ok)
	})

	t.Run("ProcessError", func(t *testing.T) {
		err := &ProcessError{ExitCode: 1, Stderr: "failed"}

		result, ok := AsType[*ProcessError](err)
		require.True(t, ok)
		require.Equal(t, 1, result.ExitCode)
		require.Equal(t, "failed", result.Stderr)
	})

	t.Run("CLIJSONDecodeError", func(t *testing.T) {
		err := &CLIJSONDecodeError{RawData: "{bad}", Err: fmt.Errorf("syntax")}

		result, ok := AsType[*CLIJSONDecodeError](err)
		require.True(t, ok)
		require.Equal(t, "{bad}", result.RawData)
	})

	t.Run("CLIConnectionError", func(t *testing.T) {
		inner := fmt.Errorf("connection refused")
		err := &CLIConnectionError{Err: inner}

		result, ok := AsType[*CLIConnectionError](err)
		require.True(t, ok)
		require.Equal(t, inner, result.Err)
	})

	t.Run("MessageParseError", func(t *testing.T) {
		err := &MessageParseError{Message: "{}", Err: fmt.Errorf("invalid")}

		result, ok := AsType[*MessageParseError](err)
		require.True(t, ok)
		require.Equal(t, "{}", result.Message)
	})

	t.Run("wrapped errors", func(t *testing.T) {
		inner := &ProcessError{ExitCode: 2, Stderr: "error"}
		wrapped := fmt.Errorf("wrapped: %w", inner)

		result, ok := AsType[*ProcessError](wrapped)
		require.True(t, ok)
		require.Equal(t, 2, result.ExitCode)
	})

	t.Run("type mismatch returns zero value", func(t *testing.T) {
		err := &ProcessError{ExitCode: 1}

		result, ok := AsType[*CLINotFoundError](err)
		require.False(t, ok)
		require.Nil(t, result)
	})
}
