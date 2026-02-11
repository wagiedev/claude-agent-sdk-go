module github.com/wagiedev/claude-agent-sdk-go/examples/custom_logger

go 1.25

require (
	github.com/wagiedev/claude-agent-sdk-go v0.1.0
	github.com/sirupsen/logrus v1.9.4
)

require (
	github.com/google/jsonschema-go v0.4.2 // indirect
	github.com/modelcontextprotocol/go-sdk v1.2.0 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/wagiedev/claude-agent-sdk-go => ../..
