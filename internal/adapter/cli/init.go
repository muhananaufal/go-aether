package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newInitCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
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

func newAdoptCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
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

func newDoctorCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
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
