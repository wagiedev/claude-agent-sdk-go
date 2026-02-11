package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	claudesdk "github.com/wagiedev/claude-agent-sdk-go"
)

func ptrString(s string) *string {
	return &s
}

// displayMessage standardizes message display across examples.
func displayMessage(msg claudesdk.Message) {
	switch m := msg.(type) {
	case *claudesdk.UserMessage:
		for _, block := range m.Content.Blocks() {
			if textBlock, ok := block.(*claudesdk.TextBlock); ok {
				fmt.Printf("User: %s\n", textBlock.Text)
			}
		}

	case *claudesdk.AssistantMessage:
		for _, block := range m.Content {
			if textBlock, ok := block.(*claudesdk.TextBlock); ok {
				fmt.Printf("Claude: %s\n", textBlock.Text)
			}
		}

	case *claudesdk.SystemMessage:
		// Ignore system messages in display

	case *claudesdk.ResultMessage:
		fmt.Println("Result ended")

		if m.TotalCostUSD != nil {
			fmt.Printf("Cost: $%.4f\n", *m.TotalCostUSD)
		}
	}
}

func codeReviewerExample() {
	fmt.Println("=== Code Reviewer Agent Example ===")

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	client := claudesdk.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	defer client.Close()

	if err := client.Start(ctx,
		claudesdk.WithLogger(logger),
		claudesdk.WithAgents(map[string]*claudesdk.AgentDefinition{
			"code-reviewer": {
				Description: "Reviews code for best practices and potential issues",
				Prompt: "You are a code reviewer. Analyze code for bugs, performance issues, " +
					"security vulnerabilities, and adherence to best practices. " +
					"Provide constructive feedback.",
				Tools: []string{"Read", "Grep"},
				Model: ptrString("sonnet"),
			},
		}),
		claudesdk.WithPermissionMode("bypassPermissions"),
	); err != nil {
		fmt.Printf("Failed to connect: %v\n", err)

		return
	}

	prompt := "Use the code-reviewer agent to review the code in internal/types/types.go"
	if err := client.Query(ctx, prompt); err != nil {
		fmt.Printf("Failed to send query: %v\n", err)

		return
	}

	for msg, err := range client.ReceiveMessages(ctx) {
		if err != nil {
			break
		}

		displayMessage(msg)

		if _, ok := msg.(*claudesdk.ResultMessage); ok {
			break
		}
	}

	fmt.Println()
}

func documentationWriterExample() {
	fmt.Println("=== Documentation Writer Agent Example ===")

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	client := claudesdk.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	defer client.Close()

	if err := client.Start(ctx,
		claudesdk.WithLogger(logger),
		claudesdk.WithAgents(map[string]*claudesdk.AgentDefinition{
			"doc-writer": {
				Description: "Writes comprehensive documentation",
				Prompt: "You are a technical documentation expert. Write clear, comprehensive " +
					"documentation with examples. Focus on clarity and completeness.",
				Tools: []string{"Read", "Write", "Edit"},
				Model: ptrString("sonnet"),
			},
		}),
		claudesdk.WithPermissionMode("bypassPermissions"),
	); err != nil {
		fmt.Printf("Failed to connect: %v\n", err)

		return
	}

	prompt := "Use the doc-writer agent to explain what AgentDefinition is used for"
	if err := client.Query(ctx, prompt); err != nil {
		fmt.Printf("Failed to send query: %v\n", err)

		return
	}

	for msg, err := range client.ReceiveMessages(ctx) {
		if err != nil {
			break
		}

		displayMessage(msg)

		if _, ok := msg.(*claudesdk.ResultMessage); ok {
			break
		}
	}

	fmt.Println()
}

func multipleAgentsExample() {
	fmt.Println("=== Multiple Agents Example ===")

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	client := claudesdk.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	defer client.Close()

	if err := client.Start(ctx,
		claudesdk.WithLogger(logger),
		claudesdk.WithAgents(map[string]*claudesdk.AgentDefinition{
			"analyzer": {
				Description: "Analyzes code structure and patterns",
				Prompt:      "You are a code analyzer. Examine code structure, patterns, and architecture.",
				Tools:       []string{"Read", "Grep", "Glob"},
			},
			"tester": {
				Description: "Creates and runs tests",
				Prompt:      "You are a testing expert. Write comprehensive tests and ensure code quality.",
				Tools:       []string{"Read", "Write", "Bash"},
				Model:       ptrString("sonnet"),
			},
		}),
		claudesdk.WithSettingSources(
			claudesdk.SettingSourceUser,
			claudesdk.SettingSourceProject,
		),
		claudesdk.WithPermissionMode("bypassPermissions"),
	); err != nil {
		fmt.Printf("Failed to connect: %v\n", err)

		return
	}

	prompt := "Use the analyzer agent to find all Go files in the examples/ directory"
	if err := client.Query(ctx, prompt); err != nil {
		fmt.Printf("Failed to send query: %v\n", err)

		return
	}

	for msg, err := range client.ReceiveMessages(ctx) {
		if err != nil {
			break
		}

		displayMessage(msg)

		if _, ok := msg.(*claudesdk.ResultMessage); ok {
			break
		}
	}

	fmt.Println()
}

func main() {
	fmt.Println("Custom Agents Examples")
	fmt.Println()

	codeReviewerExample()
	documentationWriterExample()
	multipleAgentsExample()
}
