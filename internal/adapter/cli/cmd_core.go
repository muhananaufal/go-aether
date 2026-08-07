package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/muhananaufal/go-aether/internal/adapter/cli/prompt"
	"github.com/muhananaufal/go-aether/internal/adapter/config"
	"github.com/muhananaufal/go-aether/internal/core/domain"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

// maxManifestWalkUp bounds how far a command searches upward for aether.yaml.
// Chosen to match the resolver in internal/adapter/manifest so the CLI and the
// core never disagree about which project the user is standing in.
const maxManifestWalkUp = 12

// findManifestUpwards locates aether.yaml from dir, stopping at a repository
// boundary or after maxManifestWalkUp levels. It returns "" when none is found.
func findManifestUpwards(dir string) string {
	curr := dir
	for i := 0; i < maxManifestWalkUp; i++ {
		candidate := filepath.Join(curr, "aether.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		// A .git directory marks the top of this repository. Crossing it would
		// start describing a different project entirely.
		if _, err := os.Stat(filepath.Join(curr, ".git")); err == nil && i > 0 {
			return ""
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			return ""
		}
		curr = parent
	}
	return ""
}

func newCmdAdopt(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var scan, apply bool

	cmd := &cobra.Command{
		Use:   "adopt",
		Short: "Scan and adopt a legacy brownfield Go repository into go-aether",
		Long: `Inspect an existing Go repository and propose an aether.yaml that matches
the structure that is actually there.

adopt previews by default and writes nothing. Review the proposed layer
mapping, then re-run with --apply to save it.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			// Preview is the default. This command is pointed at repositories
			// somebody already depends on, and a wrong path mapping sends every
			// later generator into a directory that means something else. Writing
			// only on an explicit --apply keeps the mistake recoverable.
			dryRun := globals.DryRun || !apply

			if err := svc.AdoptProject(cmd.Context(), cwd, scan, dryRun); err != nil {
				return err
			}

			if !dryRun {
				fmt.Println("\n🛡️ Repository is now under go-aether management.")
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&scan, "scan", true, "inspect the repository layout instead of assuming the default one")
	cmd.Flags().BoolVar(&apply, "apply", false, "write aether.yaml; without this the plan is only printed")

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
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			if len(args) > 0 {
				projectName = args[0]
			} else if prompt.IsInteractive() {
				var err error
				projectName, err = prompt.AskString("Project Name", "What is the name of your project? (e.g. invoice-api)", true)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("accepts 1 arg(s), received 0. (interactive prompt unavailable in CI)")
			}

			if !cmd.Flags().Changed("module") && prompt.IsInteractive() {
				mod, err := prompt.AskString("Go Module Name", "e.g. github.com/username/"+projectName, true)
				if err != nil {
					return err
				}
				moduleName = mod
			}

			if !cmd.Flags().Changed("arch") && prompt.IsInteractive() {
				// Option lists are read from the domain rather than typed out here.
				// Hardcoding them is how the prompt came to offer "clean", "ddd" and
				// "mux": choices with no template behind them, which a user could
				// select and only then be rejected.
				selectedArch, err := prompt.AskSelect("Architecture Pattern", domain.SupportedArchitectures())
				if err != nil {
					return err
				}
				arch = selectedArch
			}

			if !cmd.Flags().Changed("db") && prompt.IsInteractive() {
				selectedDB, err := prompt.AskSelect("Database Engine", domain.SupportedDBDrivers())
				if err != nil {
					return err
				}
				dbDriver = selectedDB
			}

			if !cmd.Flags().Changed("router") && prompt.IsInteractive() {
				selectedRouter, err := prompt.AskSelect("HTTP Router", domain.SupportedRouters())
				if err != nil {
					return err
				}
				router = selectedRouter
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed resolving current working directory: %w", err)
			}

			if globals.DryRun {
				prompt.PrintDryRunDiff(fmt.Sprintf("Init project %s with %s / %s / %s", projectName, arch, dbDriver, router))
			}

			if err := svc.InitProject(cmd.Context(), cwd, projectName, moduleName, arch, dbDriver, router, globals.DryRun); err != nil {
				return err
			}

			// Save Context Memory
			if !globals.DryRun {
				// Initialize context memory
				cfg := &config.AetherConfig{
					Version: 1,
					Preferences: config.ProjectPreferences{
						Architecture: arch,
						ORM:          dbDriver,
						Engine:       router,
					},
				}

				// Ensure config is saved in the project root directory
				projectRoot := filepath.Join(cwd, projectName)

				// Change working directory temporarily to save config there
				oldWd, _ := os.Getwd()
				_ = os.Chdir(projectRoot)
				_ = config.SaveConfig(cfg)
				_ = os.Chdir(oldWd)
			}

			fmt.Printf("🚀 Successfully initialized project [%s] with [%s] architecture!\n", projectName, arch)
			return nil
		},
	}

	cmd.Flags().StringVar(&moduleName, "module", "github.com/example/app", "Go module identifier for go.mod")
	// Help text is generated from the same sets the validator enforces, so it
	// cannot drift into advertising a value that would then be rejected.
	cmd.Flags().StringVar(&arch, "arch", domain.DefaultArchitecture,
		"Architecture blueprint ("+strings.Join(domain.SupportedArchitectures(), ", ")+")")
	cmd.Flags().StringVar(&dbDriver, "db", domain.DefaultDBDriver,
		"Database engine driver ("+strings.Join(domain.SupportedDBDrivers(), ", ")+")")
	cmd.Flags().StringVar(&router, "router", domain.DefaultRouter,
		"HTTP routing framework ("+strings.Join(domain.SupportedRouters(), ", ")+")")

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

			// Bounded exactly like the manifest resolver. This loop used to walk to
			// the filesystem root, so running `ls` in a directory with no project
			// could find an unrelated aether.yaml several levels up and describe
			// somebody else's project as if it were this one.
			manifestPath := findManifestUpwards(cwd)

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
