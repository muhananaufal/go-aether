package main

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/adapter/cli"
	"github.com/muhananaufal/go-aether/internal/adapter/manifest"
	"github.com/muhananaufal/go-aether/internal/adapter/template"
	"github.com/muhananaufal/go-aether/internal/adapter/writer"
	"github.com/muhananaufal/go-aether/internal/core/service"
	"github.com/muhananaufal/go-aether/templates"
	"github.com/spf13/afero"
)

// Injected by Goreleaser via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Initialize core ports & adapters ("Dogfooding Hexagonal Architecture")
	osFS := afero.NewOsFs()
	fileWriter := writer.NewAferoWriter(osFS)
	resolver := manifest.NewYamlResolver(fileWriter)
	engine := template.NewStdEngine(templates.FS)

	// Wire service orchestration layer
	scaffoldService := service.NewAetherScaffoldService(engine, resolver, fileWriter)

	// Bind CLI root entrypoint and execute
	rootCommand := cli.NewRootCommand(scaffoldService, version, commit, date)
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing go-aether: %v\n", err)
		os.Exit(1)
	}
}
