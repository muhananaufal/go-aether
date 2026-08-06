package cli

import (
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

type globalFlags struct {
	Verbose bool
	DryRun  bool
}

// NewRootCommand constructs the parent cobra command and attaches all available subcommands.
func NewRootCommand(svc port.ScaffoldService) *cobra.Command {
	var flags globalFlags

	rootCmd := &cobra.Command{
		Use:   "go-aether",
		Short: "Opinionated Architecture Scaffold CLI Engine for Go Backend Engineers",
		Long: `go-aether is a lightning-fast dev-time CLI scaffolding engine built to standardize
Hexagonal / Ports & Adapters architectures, enforce zero runtime overhead, and eliminate legacy boilerplate.`,
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "enable verbose structured logging output")
	rootCmd.PersistentFlags().BoolVar(&flags.DryRun, "dry-run", false, "simulate execution without mutating physical filesystem")

	// Attach subcommands
	rootCmd.AddCommand(newInitCommand(svc, &flags))
	rootCmd.AddCommand(newAdoptCommand(svc, &flags))
	rootCmd.AddCommand(newDoctorCommand(svc, &flags))
	rootCmd.AddCommand(newMakeModuleCommand(svc, &flags))
	rootCmd.AddCommand(newMakeServiceCommand(svc, &flags))
	rootCmd.AddCommand(newMakeHandlerCommand(svc, &flags))
	rootCmd.AddCommand(newAddMiddlewareCommand(svc, &flags))
	rootCmd.AddCommand(newAddCacheCommand(svc, &flags))
	rootCmd.AddCommand(newAddTransportCommand(svc, &flags))

	return rootCmd
}
