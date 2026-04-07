package openai_test

import (
	"context"
	"log"
	"os"

	openai "github.com/amammay/adk-go-openai"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	libopenapi "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"golang.org/x/oauth2"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
)

func githubMCPTransport(ctx context.Context) mcp.Transport {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")},
	)
	return &mcp.StreamableClientTransport{
		Endpoint:   "https://api.githubcopilot.com/mcp/x/repos/readonly",
		HTTPClient: oauth2.NewClient(ctx, ts),
	}
}

func Example() {

	ctx := context.Background()

	client := libopenapi.NewClient(option.WithBaseURL("http://localhost:11434/v1"))

	model, err := openai.NewModel(ctx, "gemma4", client)
	if err != nil {
		log.Fatalf("failed to create model: %v", err)
	}

	ghMCP, err := mcptoolset.New(mcptoolset.Config{
		Transport: githubMCPTransport(ctx),
	})
	if err != nil {
		log.Fatalf("failed to create GitHub MCP toolset: %v", err)
	}

	commitFetcherAgent, err := llmagent.New(llmagent.Config{
		Name:        "github_repo_commit_agent",
		Description: "Specialist in GitHub data extraction. Retrieves the latest 5 commits and outputs only author and PR description list items.",
		Model:       model,
		Instruction: `You are a specialist in extracting GitHub commit data.
Your task is to retrieve the latest 5 commits for the repository specified by the user.
Call list_commits with page set to 1 and perPage set to 5.
Output MUST be non-empty plain text.
If commits are found, output only a bullet list where each bullet is:
- AUTHOR: <author name> | CHANGE: <PR description or commit message>
If there are no commits, output exactly: "No commits found."
If tool results are incomplete, still output at least one bullet using available fields and "unknown" placeholders.
Do not output explanations before or after the list/message.`,
		Toolsets:  []tool.Toolset{ghMCP},
		OutputKey: "raw_commits",
	})
	if err != nil {

	}

	poemMakerAgent, err := llmagent.New(llmagent.Config{
		Name:        "commit_poem_maker_agent",
		Description: "Turns raw GitHub commit data into a poem about the commits and their authors.",
		Model:       model,
		Instruction: `You are a creative poet describing GitHub activity.
You will be given a list of the latest commits:

{raw_commits}

You MUST always return a poem. Never return an empty response.
Write one short poem that:
1. Mentions each author by name.
2. References the key commit themes or changes for each author.
3. Keeps a playful, readable tone.
4. Adds at least one fitting emoji in the poem.
5. Uses 4 to 8 lines.

If the input says "No commits found.", still write a short poem about the quiet repository and include that exact phrase in one line.
If input is partial or malformed, write a best-effort poem from available text and mention uncertainty poetically instead of failing.`,
	})
	if err != nil {
		log.Fatalf("failed to create poem maker agent: %v", err)
	}

	pipelineAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "commit_poem_pipeline",
			Description: "Fetches the latest 5 commits for a GitHub repo and always produces an author-focused poem.",
			SubAgents:   []agent.Agent{commitFetcherAgent, poemMakerAgent},
		},
	})
	if err != nil {
		log.Fatalf("failed to create pipeline agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader:    agent.NewSingleLoader(pipelineAgent),
		SessionService: session.InMemoryService(),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("failed to execute: %v", err)
	}

}
