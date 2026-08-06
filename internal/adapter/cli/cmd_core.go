package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdAdopt(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var scan bool

	cmd := &cobra.Command{
		Use:   "adopt",
		Short: "Scan and adopt a legacy brownfield Go repository into go-aether",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			if err := svc.AdoptProject(cmd.Context(), cwd, scan, globals.DryRun); err != nil {
				return err
			}
			fmt.Println("🛡️ Successfully adopted existing repository under go-aether management!")
			return nil
		},
	}

	cmd.Flags().BoolVar(&scan, "scan", false, "perform interactive anomaly directory scanning")

	return cmd
}

func newCmdDoctor(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var fix bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run structural health diagnostics against aether.yaml and project layout",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			return svc.RunDoctor(cmd.Context(), cwd, fix, cmd.OutOrStdout())
		},
	}

	cmd.Flags().BoolVar(&fix, "fix", false, "attempt automatic remediation of manifest inconsistencies")

	return cmd
}

func newCmdInit(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var moduleName, arch, dbDriver, router string

	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Bootstrap a greenfield Go backend project with clean Hexagonal Architecture",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed resolving current working directory: %w", err)
			}

			if err := svc.InitProject(cmd.Context(), cwd, projectName, moduleName, arch, dbDriver, router, globals.DryRun); err != nil {
				return err
			}
			fmt.Printf("🚀 Successfully initialized project [%s] with [%s] architecture!\n", projectName, arch)
			return nil
		},
	}

	cmd.Flags().StringVar(&moduleName, "module", "github.com/example/app", "Go module identifier for go.mod")
	cmd.Flags().StringVar(&arch, "arch", "hexagonal", "Architecture blueprint (hexagonal, clean, ddd)")
	cmd.Flags().StringVar(&dbDriver, "db", "postgres", "Database engine driver (postgres, mysql, sqlite)")
	cmd.Flags().StringVar(&router, "router", "chi", "HTTP routing framework (chi, gin, echo)")

	return cmd
}

func newCmdLs(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
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
