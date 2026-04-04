package openai

import (
	"strings"
	"testing"

	"google.golang.org/genai"
)

func TestConvertTools_MCPServer(t *testing.T) {
	cfg := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				MCPServers: []*genai.MCPServer{
					{
						Name: "github",
						StreamableHTTPTransport: &genai.StreamableHTTPTransport{
							URL: "https://api.githubcopilot.com/mcp/x/repos/readonly",
							Headers: map[string]string{
								"Authorization": "Bearer token",
							},
						},
					},
				},
			},
		},
	}

	tools, err := convertTools(cfg)
	if err != nil {
		t.Fatalf("convertTools() err = %v", err)
	}
	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}
	if tools[0].OfMcp == nil {
		t.Fatalf("expected MCP tool, got %+v", tools[0])
	}
	if got, want := tools[0].OfMcp.ServerLabel, "github"; got != want {
		t.Fatalf("server label mismatch got=%q want=%q", got, want)
	}
	if !tools[0].OfMcp.ServerURL.Valid() {
		t.Fatalf("expected server URL to be set")
	}
	if got, want := tools[0].OfMcp.ServerURL.Value, "https://api.githubcopilot.com/mcp/x/repos/readonly"; got != want {
		t.Fatalf("server URL mismatch got=%q want=%q", got, want)
	}
	if got, want := tools[0].OfMcp.Headers["Authorization"], "Bearer token"; got != want {
		t.Fatalf("header mismatch got=%q want=%q", got, want)
	}
}

func TestConvertTools_MCPServerMissingTransportURL(t *testing.T) {
	cfg := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				MCPServers: []*genai.MCPServer{
					{
						Name:                    "github",
						StreamableHTTPTransport: &genai.StreamableHTTPTransport{},
					},
				},
			},
		},
	}

	_, err := convertTools(cfg)
	if err == nil || !strings.Contains(err.Error(), "missing transport url") {
		t.Fatalf("expected missing transport url error, got %v", err)
	}
}

func TestConvertTools_RejectUnsupportedNonFunctionNonMCP(t *testing.T) {
	cfg := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{{GoogleSearch: &genai.GoogleSearch{}}},
	}

	_, err := convertTools(cfg)
	if err == nil || !strings.Contains(err.Error(), "only function and mcp tools are supported") {
		t.Fatalf("expected unsupported tool error, got %v", err)
	}
}
