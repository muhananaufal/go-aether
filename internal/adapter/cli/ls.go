package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newLsCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List all active modules, architectural paths, and installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			// Read aether.yaml if present in directory hierarchy
			var manifestPath string
			curr := cwd
			for {
				candidate := filepath.Join(curr, "aether.yaml")
				if _, err := os.Stat(candidate); err == nil {
					manifestPath = candidate
					break
				}
				parent := filepath.Dir(curr)
				if parent == curr {
					break
				}
				curr = parent
			}

			if manifestPath == "" {
				return fmt.Errorf("no active aether project found in current directory tree (missing aether.yaml)")
			}

			projectRoot := filepath.Dir(manifestPath)
			projectName := filepath.Base(projectRoot)

			fmt.Printf("📦 Project: %s (%s)\n", projectName, manifestPath)
			fmt.Println(strings.Repeat("─", 60))

			// Discover modules inside internal/core/domain or internal/core/service
			domainDir := filepath.Join(projectRoot, "internal", "core", "domain")
			var modules []string
			if entries, err := os.ReadDir(domainDir); err == nil {
				for _, entry := range entries {
					if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
						name := strings.TrimSuffix(entry.Name(), ".go")
						modules = append(modules, name)
					}
				}
			}

			fmt.Println("🔷 Registered Domain Modules:")
			if len(modules) == 0 {
				fmt.Println("  (No domain modules generated yet. Use `go-aether make:module <name>`)")
			} else {
				for _, mod := range modules {
					fmt.Printf("  • %-20s [internal/core/domain/%s.go]\n", mod, mod)
				}
			}

			fmt.Println("\n🔌 Installed Plugins & Ecosystem Packages:")
			pkgDir := filepath.Join(projectRoot, "pkg")
			var plugins []string
			if entries, err := os.ReadDir(pkgDir); err == nil {
				for _, entry := range entries {
					if entry.IsDir() {
						plugins = append(plugins, entry.Name())
					}
				}
			}

			if len(plugins) == 0 {
				fmt.Println("  (No extra plugins installed. Use `go-aether add:<plugin>`)")
			} else {
				for _, p := range plugins {
					fmt.Printf("  ✓ pkg/%-18s [Active]\n", p)
				}
			}

			fmt.Println(strings.Repeat("─", 60))
			return nil
		},
	}
	return cmd
}
