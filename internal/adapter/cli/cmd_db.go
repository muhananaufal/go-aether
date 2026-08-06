package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/adapter/cli/prompt"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdDBMigration(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "db:migration [name]",
		Short: "Generate a SQL migration file pair (up/down)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg0, err := prompt.GetArgOrPrompt(args, 0, "Argument", "Please provide the required argument", true)
			if err != nil {
				return err
			}
			name := arg0
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

func newCmdDBPaginator(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "db:paginator",
		Short: "Generate cursor-based opaque base64 pagination helper",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.MakeCursorPaginator(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("📄 Generated cursor-based pagination helper")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing file")
	return cmd
}

func newCmdDBReadreplica(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "db:readreplica",
		Short: "Set up DB connection pool splitter for Primary-Write and Replica-Read",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddReadReplica(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("⚖️ Generated Read/Write Replica connection pool manager")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing file")
	return cmd
}

func newCmdDBSQLC(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "db:sqlc",
		Short: "Set up SQLC configuration, base schema, and type-safe query templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddSQLC(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🗄️ Successfully configured SQLC type-safe code generator in sqlc.yaml and db/")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing SQLC files")
	return cmd
}

func newCmdDBSeeder(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "db:seeder [name]",
		Short: "Generate a database seeder file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg0, err := prompt.GetArgOrPrompt(args, 0, "Argument", "Please provide the required argument", true)
			if err != nil {
				return err
			}
			name := arg0
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

func newCmdDBUOW(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "db:uow",
		Short: "Set up Unit of Work transactional orchestrator pattern",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddUOW(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("📦 Generated Unit of Work orchestrator interface")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing file")
	return cmd
}
