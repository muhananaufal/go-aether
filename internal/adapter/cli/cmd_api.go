package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/adapter/cli/prompt"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdAPIGRPC(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "api:grpc",
		Short: "Set up gRPC bi-directional duplex streaming RPC handler and protobuf contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddGRPCStream(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("⚡ Generated gRPC bi-directional duplex streaming handler and protobuf contract")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing gRPC stream files")
	return cmd
}

func newCmdAPIGRPCGateway(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "api:grpc-gateway",
		Short: "Set up gRPC-Gateway reverse-proxy HTTP REST JSON bridge",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddGRPCGateway(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🌐 Generated gRPC-Gateway HTTP/REST JSON reverse-proxy server")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing gRPC-Gateway files")
	return cmd
}

func newCmdAPIGraphql(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "api:graphql",
		Short: "Set up gqlgen GraphQL server with dataloader boilerplate",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddGraphQL(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🕸️ Generated GraphQL server bootstrap")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing file")
	return cmd
}

func newCmdAPIIdempotency(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "api:idempotency [redis|memory]",
		Short: "Set up Idempotency-Key validation middleware",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var engine string
			if len(args) > 0 {
				engine = args[0]
			} else if prompt.IsInteractive() {
				var err error
				engine, err = prompt.AskSelect("Idempotency Engine", []string{"redis", "memory"})
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("accepts 1 arg(s), received 0. (interactive prompt unavailable in CI)")
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddIdempotency(cmd.Context(), cwd, engine, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔁 Injected [%s] idempotency key engine\n", engine)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing idempotency files")
	return cmd
}

func newCmdAPIMiddleware(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var module string
	var force bool

	cmd := &cobra.Command{
		Use:   "api:middleware [middleware-type]",
		Short: "Inject a middleware (e.g. jwt-auth, rate-limit) into a module's transport handler",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var mwType string
			if len(args) > 0 {
				mwType = args[0]
			} else if prompt.IsInteractive() {
				var err error
				mwType, err = prompt.AskString("Middleware Type", "e.g. jwt-auth, rate-limit", true)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("accepts 1 arg(s), received 0. (interactive prompt unavailable in CI)")
			}

			if !cmd.Flags().Changed("module") && prompt.IsInteractive() {
				mod, err := prompt.AskString("Target Module", "Module name to inject middleware into (e.g. users)", true)
				if err != nil {
					return err
				}
				module = mod
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if module == "" {
				return fmt.Errorf("--module flag is required to target a specific module's handler")
			}

			if globals.DryRun {
				prompt.PrintDryRunDiff(fmt.Sprintf("Inject middleware %s into module %s", mwType, module))
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

	return cmd
}

func newCmdAPIOpenapi(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "api:openapi",
		Short: "Set up Swagger OpenAPI docs generator middleware",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddOpenAPI(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("📖 Generated Swagger OpenAPI UI middleware route")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing file")
	return cmd
}

func newCmdAPITransport(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api:transport [transport-type]",
		Short: "Register a new global transport protocol in aether.yaml",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var transport string
			if len(args) > 0 {
				transport = args[0]
			} else if prompt.IsInteractive() {
				var err error
				transport, err = prompt.AskSelect("Transport Protocol", []string{"http", "grpc", "graphql"})
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("accepts 1 arg(s), received 0. (interactive prompt unavailable in CI)")
			}

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

func newCmdAPIValidator(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "api:validator [validator-type]",
		Short: "Set up the struct validation wrapper (e.g. playground)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var validatorType string
			if len(args) > 0 {
				validatorType = args[0]
			} else if prompt.IsInteractive() {
				var err error
				validatorType, err = prompt.AskSelect("Validator Wrapper", []string{"playground"})
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("accepts 1 arg(s), received 0. (interactive prompt unavailable in CI)")
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddValidator(cmd.Context(), cwd, validatorType, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✅ Injected [%s] struct validation wrapper\n", validatorType)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing validator file")
	return cmd
}
