package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/adapter/cli/prompt"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdInfraCICD(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "infra:cicd [provider]",
		Short: "Set up the CI/CD pipeline (github, gitlab)",
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

			err = svc.AddCICD(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🚀 Injected [%s] CI/CD pipeline\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing cicd file")
	return cmd
}

func newCmdInfraConfig(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "infra:config [config-type]",
		Short: "Set up a centralized configuration manager (e.g. viper, koanf)",
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

			err = svc.AddConfig(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("⚙️ Injected [%s] centralized configuration manager\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing config file")
	return cmd
}

func newCmdInfraDeploy(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "infra:deploy [target]",
		Short: "Set up the deployment manifests (k8s, helm, lambda)",
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

			err = svc.AddDeploy(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("☁️ Injected [%s] deployment manifests\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing deploy file")
	return cmd
}

func newCmdInfraDrain(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "infra:drain",
		Short: "Set up zero-downtime graceful shutdown and connection draining manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddDrain(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🛑 Generated graceful drain manager in pkg/server/")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing drain files")
	return cmd
}

func newCmdInfraFeatureflags(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "infra:featureflags [flipt]",
		Short: "Set up feature flags & canary release client (flipt)",
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

			err = svc.AddFeatureFlags(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🚩 Injected [%s] feature flags client\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing featureflags files")
	return cmd
}

func newCmdInfraHealthcheck(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "infra:healthcheck",
		Short: "Set up Kubernetes Liveness (/livez) and Readiness (/readyz) probe handlers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddHealthcheck(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🩺 Injected Kubernetes /livez and /readyz probe handlers")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing healthcheck files")
	return cmd
}

func newCmdInfraLint(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "infra:lint",
		Short: "Set up linter rules (golangci-lint) and git pre-commit hooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddLint(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🧹 Scaffolded golangci-lint configuration rules")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing file")
	return cmd
}

func newCmdInfraProfiling(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "infra:profiling [pprof]",
		Short: "Set up protected runtime profiling endpoints (pprof)",
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

			err = svc.AddProfiling(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📈 Injected [%s] protected profiling endpoints\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing profiling files")
	return cmd
}

func newCmdInfraSearch(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "infra:search [meilisearch|elasticsearch]",
		Short: "Set up full-text search engine client",
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

			err = svc.AddSearch(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔍 Injected [%s] full-text search engine client\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing search files")
	return cmd
}
