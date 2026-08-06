package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/muhananaufal/go-aether/internal/adapter/cli/prompt"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdArchAggregate(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:aggregate [name]",
		Short: "Generate a DDD Aggregate Root entity with event recording",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeAggregate(cmd.Context(), cwd, name, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🌳 Generated Aggregate Root [%s]\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}

func newCmdArchCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:command [name]",
		Short: "Generate a CQRS Command DTO and execution handler",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeCommand(cmd.Context(), cwd, name, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("⚡ Generated CQRS Command [%s]\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}

func newCmdArchDI(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:di [di-type]",
		Short: "Set up a dependency injection container (e.g. fx, wire)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddDI(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🧩 Injected [%s] dependency injection container\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing DI file")
	return cmd
}

func newCmdArchDomain(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:domain [module-name]",
		Short: "Generate only the domain layer entity for a specific module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeDomain(cmd.Context(), cwd, moduleName, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✨ Generated domain entity for module [%s]\n", moduleName)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}

func newCmdArchError(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:error",
		Short: "Set up a standardized centralized error handler",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddError(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🚨 Injected standardized error handler")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing error file")
	return cmd
}

func newCmdArchEvent(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:event [name]",
		Short: "Generate a Domain Event struct and serializer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeEvent(cmd.Context(), cwd, name, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📣 Generated Domain Event [%s]\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}

func newCmdArchHandler(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var transport string
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:handler [module-name]",
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

func newCmdArchMock(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:mock [interface-name]",
		Short: "Scaffold interface mock implementation using Mockery directives for isolated unit tests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaceName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := svc.MakeMock(cmd.Context(), cwd, interfaceName, globals.DryRun, force); err != nil {
				return err
			}

			fmt.Printf("🎭 Generated Mockery mock contract for [%s] in mocks/\n", interfaceName)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing mock file")

	return cmd
}

func newCmdArchModule(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var transports []string
	var hasCache, hasWorker, force bool

	cmd := &cobra.Command{
		Use:   "arch:module [module-name]",
		Short: "Generate a complete vertical feature slice (Domain, Port, Service, Handler, Repo)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var moduleName string
			if len(args) > 0 {
				moduleName = args[0]
			} else if prompt.IsInteractive() {
				var err error
				moduleName, err = prompt.AskString("Module Name", "Name of the new vertical slice (e.g. users, products)", true)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("accepts 1 arg(s), received 0. (interactive prompt unavailable in CI)")
			}

			if !cmd.Flags().Changed("transports") && prompt.IsInteractive() {
				t, err := prompt.AskString("Transports", "Comma-separated transport targets (e.g. http, grpc)", false)
				if err == nil && t != "" {
					transports = strings.Split(t, ",")
				}
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if globals.DryRun {
				prompt.PrintDryRunDiff(fmt.Sprintf("Generate vertical slice for module %s with transports %v", moduleName, transports))
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

func newCmdArchPipeline(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:pipeline [name]",
		Short: "Generate Fan-Out / Fan-In bounded concurrency pipeline helper",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakePipeline(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🌊 Generated %s pipeline in pkg/concurrency/\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing pipeline file")
	return cmd
}

func newCmdArchPort(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:port [module-name]",
		Short: "Generate only the port interface contract for a specific module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakePort(cmd.Context(), cwd, moduleName, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✨ Generated port interface for module [%s]\n", moduleName)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}

func newCmdArchQuery(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:query [name]",
		Short: "Generate a CQRS Query DTO and read-model handler",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeQuery(cmd.Context(), cwd, name, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔍 Generated CQRS Query [%s]\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}

func newCmdArchRepository(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:repository [module-name]",
		Short: "Generate only the infrastructure repository for a specific module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeRepository(cmd.Context(), cwd, moduleName, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✨ Generated infrastructure repository for module [%s]\n", moduleName)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}

func newCmdArchService(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:service [module-name]",
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

func newCmdArchSpecification(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:specification [name]",
		Short: "Generate reusable DDD Specification pattern for dynamic query rules",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeSpecification(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📜 Generated %s specification in internal/core/domain/\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing specification file")
	return cmd
}

func newCmdArchValueobject(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:valueobject [name]",
		Short: "Generate an immutable DDD Value Object",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeValueObject(cmd.Context(), cwd, name, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("💎 Generated Value Object [%s]\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}
