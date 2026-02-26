// Package claudesdk provides a Go SDK for interacting with the Claude CLI agent.
//
// This SDK enables Go applications to programmatically communicate with Claude
// through the official Claude CLI tool. It supports both one-shot queries and
// interactive multi-turn conversations.
//
// # Basic Usage
//
// For simple, one-shot queries, use the Query function:
//
//	ctx := context.Background()
//	messages, err := claudesdk.Query(ctx, "What is 2+2?",
//	    claudesdk.WithPermissionMode("acceptEdits"),
//	    claudesdk.WithMaxTurns(1),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for msg := range messages {
//	    switch m := msg.(type) {
//	    case *claudesdk.AssistantMessage:
//	        for _, block := range m.Content {
//	            if text, ok := block.(*claudesdk.TextBlock); ok {
//	                fmt.Println(text.Text)
//	            }
//	        }
//	    case *claudesdk.ResultMessage:
//	        fmt.Printf("Completed in %dms\n", m.DurationMs)
//	    }
//	}
//
// # Interactive Sessions
//
// For multi-turn conversations, use NewClient or the WithClient helper:
//
//	// Using WithClient for automatic lifecycle management
//	err := claudesdk.WithClient(ctx, func(c claudesdk.Client) error {
//	    if err := c.Query(ctx, "Hello Claude"); err != nil {
//	        return err
//	    }
//	    for msg, err := range c.ReceiveResponse(ctx) {
//	        if err != nil {
//	            return err
//	        }
//	        // process message...
//	    }
//	    return nil
//	},
//	    claudesdk.WithLogger(slog.Default()),
//	    claudesdk.WithPermissionMode("acceptEdits"),
//	)
//
//	// Or using NewClient directly for more control
//	client := claudesdk.NewClient()
//	defer client.Close(ctx)
//
//	err := client.Start(ctx,
//	    claudesdk.WithLogger(slog.Default()),
//	    claudesdk.WithPermissionMode("acceptEdits"),
//	)
//
// # Logging
//
// For detailed operation tracking, use WithLogger:
//
//	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
//	messages, err := claudesdk.Query(ctx, "Hello Claude",
//	    claudesdk.WithLogger(logger),
//	)
//
// # Error Handling
//
// The SDK provides typed errors for different failure scenarios:
//
//	messages, err := claudesdk.Query(ctx, prompt, claudesdk.WithPermissionMode("acceptEdits"))
//	if err != nil {
//	    if cliErr, ok := errors.AsType[*claudesdk.CLINotFoundError](err); ok {
//	        log.Fatalf("Claude CLI not installed, searched: %v", cliErr.SearchedPaths)
//	    }
//	    if procErr, ok := errors.AsType[*claudesdk.ProcessError](err); ok {
//	        log.Fatalf("CLI process failed with exit code %d: %s", procErr.ExitCode, procErr.Stderr)
//	    }
//	    log.Fatal(err)
//	}
//
// # Requirements
//
// This SDK requires the Claude CLI to be installed and available in your system PATH.
// You can specify a custom CLI path using the WithCliPath option.
package claudesdk
