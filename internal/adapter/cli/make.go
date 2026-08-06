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

func newMakeDomainCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:domain [module-name]",
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

func newMakePortCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:port [module-name]",
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

func newMakeRepositoryCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:repository [module-name]",
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

func newMakeMigrationCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:migration [name]",
		Short: "Generate a SQL migration file pair (up/down)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeMigration(cmd.Context(), cwd, name, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✨ Generated migration files for [%s]\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}

func newMakeSeederCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:seeder [name]",
		Short: "Generate a database seeder file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeSeeder(cmd.Context(), cwd, name, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✨ Generated database seeder for [%s]\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing generated target files")
	return cmd
}


func newMakePipelineCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:pipeline [name]",
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

func newMakeSpecificationCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:specification [name]",
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

func newMakeValueObjectCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:valueobject [name]",
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

func newMakeAggregateCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:aggregate [name]",
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

func newMakeEventCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:event [name]",
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

func newMakeCommandCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:command [name]",
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

func newMakeQueryCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "make:query [name]",
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
