package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newMakeModuleCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var transports []string
	var hasCache, hasWorker, force bool

	cmd := &cobra.Command{
		Use:   "make:module [module-name]",
		Short: "Generate a complete vertical feature slice (Domain, Port, Service, Handler, Repo)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeModule(cmd.Context(), cwd, moduleName, transports, hasCache, hasWorker, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✨ Generated vertical slice for module [%s] cleanly across Hexagonal boundaries.\n", moduleName)
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&transports, "transports", []string{"http"}, "Comma-separated transport targets (http, grpc, nats)")
	cmd.Flags().BoolVar(&hasCache, "cache", false, "Inject L1/L2 Redis caching decorators into repository layer")
	cmd.Flags().BoolVar(&hasWorker, "worker", false, "Scaffold async worker job processor for domain events")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")

	return cmd
}

func newMakeServiceCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:service [module-name]",
		Short: "Generate only the service layer component for a specific module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeService(cmd.Context(), cwd, moduleName, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✨ Generated service layer for module [%s]\n", moduleName)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}

func newMakeHandlerCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var transport string
	var force bool

	cmd := &cobra.Command{
		Use:   "make:handler [module-name]",
		Short: "Generate only the transport handler component for a specific module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeHandler(cmd.Context(), cwd, moduleName, transport, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✨ Generated [%s] handler for module [%s]\n", transport, moduleName)
			return nil
		},
	}
	cmd.Flags().StringVar(&transport, "transport", "http", "Transport protocol target (e.g. http, grpc)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}
