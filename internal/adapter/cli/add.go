package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newAddMiddlewareCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var module string
	var force bool

	cmd := &cobra.Command{
		Use:   "add:middleware [middleware-type]",
		Short: "Inject a middleware (e.g. jwt-auth, rate-limit) into a module's transport handler",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mwType := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if module == "" {
				return fmt.Errorf("--module flag is required to target a specific module's handler")
			}

			err = svc.AddMiddleware(cmd.Context(), cwd, module, mwType, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🛡️ Injected [%s] middleware securely into module [%s]\n", mwType, module)
			return nil
		},
	}
	
	cmd.Flags().StringVarP(&module, "module", "m", "", "Target module name to inject middleware into")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing middleware package file")
	cmd.MarkFlagRequired("module")

	return cmd
}

func newAddCacheCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:cache [cache-type]",
		Short: "Set up the global cache layer configuration and generate the cache provider infrastructure",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cacheType := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddCache(cmd.Context(), cwd, cacheType, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("⚡ Injected [%s] cache infrastructure\n", cacheType)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing cache package file")
	return cmd
}

func newAddTransportCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add:transport [transport-type]",
		Short: "Register a new global transport protocol in aether.yaml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			transport := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddTransport(cmd.Context(), cwd, transport, globals.DryRun, false)
			if err != nil {
				return err
			}

			fmt.Printf("📡 Registered [%s] transport to global stack\n", transport)
			return nil
		},
	}
	return cmd
}
