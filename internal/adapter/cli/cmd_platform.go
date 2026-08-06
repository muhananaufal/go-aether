package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/adapter/cli/prompt"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdPlatformAI(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "platform:ai [provider]",
		Short: "Set up the LLM proxy infrastructure (e.g. llm-proxy)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg0, err := prompt.GetArgOrPrompt(args, 0, "Argument", "Please provide the required argument", true)
			if err != nil {
				return err
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddAI(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🤖 Injected [%s] AI LLM-Proxy infrastructure\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing ai file")
	return cmd
}

func newCmdPlatformCQRS(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "platform:cqrs [module-name]",
		Short: "Set up CQRS Command and Query handlers within module scope",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg0, err := prompt.GetArgOrPrompt(args, 0, "Argument", "Please provide the required argument", true)
			if err != nil {
				return err
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddCQRS(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("⚡ Injected CQRS handlers for [%s]\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing CQRS files")
	return cmd
}

func newCmdPlatformDiscovery(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "platform:discovery [provider]",
		Short: "Set up service discovery client (consul, etcd)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg0, err := prompt.GetArgOrPrompt(args, 0, "Argument", "Please provide the required argument", true)
			if err != nil {
				return err
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddDiscovery(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🌐 Injected [%s] service discovery client\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing discovery files")
	return cmd
}

func newCmdPlatformLock(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "platform:lock [redis]",
		Short: "Set up distributed lock (Redlock) engine",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg0, err := prompt.GetArgOrPrompt(args, 0, "Argument", "Please provide the required argument", true)
			if err != nil {
				return err
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddLock(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔒 Injected [%s] distributed mutex lock engine\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing lock files")
	return cmd
}

func newCmdPlatformMultitenancy(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "platform:multitenancy [module-name]",
		Short: "Set up Row Level Security (RLS) SQL policies for tenant isolation",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg0, err := prompt.GetArgOrPrompt(args, 0, "Argument", "Please provide the required argument", true)
			if err != nil {
				return err
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddMultitenancy(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🏢 Injected multitenancy RLS script for [%s]\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing multitenancy file")
	return cmd
}

func newCmdPlatformResilience(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "platform:resilience [hystrix|resilience4go]",
		Short: "Set up Circuit Breaker and Bulkhead resilience engine",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg0, err := prompt.GetArgOrPrompt(args, 0, "Argument", "Please provide the required argument", true)
			if err != nil {
				return err
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddResilience(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🛡️ Injected [%s] Circuit Breaker and Bulkhead resilience engine\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing resilience files")
	return cmd
}

func newCmdPlatformTenant(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "platform:tenant",
		Short: "Set up multi-tenancy middleware and tenant context isolation helper",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddTenantContext(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🏢 Generated multi-tenancy isolation context and middleware in pkg/tenant/")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing tenant context files")
	return cmd
}
