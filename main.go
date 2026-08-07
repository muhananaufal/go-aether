package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/adapter/cli"
	"github.com/muhananaufal/go-aether/internal/adapter/manifest"
	"github.com/muhananaufal/go-aether/internal/adapter/scanner"
	"github.com/muhananaufal/go-aether/internal/adapter/template"
	"github.com/muhananaufal/go-aether/internal/adapter/writer"
	"github.com/muhananaufal/go-aether/internal/core/domain"
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
	layoutScanner := scanner.NewGoLayoutScanner()
	byoDetector := scanner.NewGoBYODetector()

	// Wire service orchestration layer
	scaffoldService := service.NewAetherScaffoldService(engine, resolver, fileWriter, layoutScanner, byoDetector)

	// Bind CLI root entrypoint and execute
	if err := cli.Execute(scaffoldService, version, commit, date); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitCodeFor(err))
	}
}

// Exit codes. Scripts and CI pipelines need to tell "you asked for something
// impossible" apart from "something broke while doing it"; collapsing both into
// 1 means a typo in a flag is indistinguishable from a full disk.
const (
	exitOperational = 1
	exitUsage       = 2
)

// exitCodeFor maps a failure to its exit code.
//
// Everything listed here is the user's input being wrong in a way they can fix
// by editing the command line or the manifest. Anything else is the tool or the
// environment failing, and stays at 1.
func exitCodeFor(err error) int {
	usageErrors := []error{
		domain.ErrUnsupportedStack,
		domain.ErrInvalidIdentifier,
		domain.ErrStdlibCollision,
		domain.ErrReservedName,
		domain.ErrGoModMissing,
		domain.ErrManifestNotFound,
		domain.ErrPathEscape,
		domain.ErrModuleAlreadyExists,
		domain.ErrFileConflict,
	}
	for _, target := range usageErrors {
		if errors.Is(err, target) {
			return exitUsage
		}
	}
	return exitOperational
}
