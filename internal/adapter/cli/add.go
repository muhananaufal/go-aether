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

func newAddWorkerCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var broker, pattern string
	var force bool

	cmd := &cobra.Command{
		Use:   "add:worker [worker-name]",
		Short: "Generate an asynchronous background processor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workerName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddWorker(cmd.Context(), cwd, workerName, broker, pattern, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("⚙️ Injected [%s] background worker using pattern [%s]\n", workerName, pattern)
			return nil
		},
	}
	cmd.Flags().StringVar(&broker, "broker", "redis", "Message broker engine (redis, kafka)")
	cmd.Flags().StringVar(&pattern, "pattern", "asynq", "Worker library pattern (asynq, kafka)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing worker file")
	return cmd
}

func newAddEventingCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var broker string
	var force bool

	cmd := &cobra.Command{
		Use:   "add:eventing",
		Short: "Set up the global Publisher and Subscriber interfaces for event-driven architecture",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddEventing(cmd.Context(), cwd, broker, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📨 Injected Event-Bus Publisher/Subscriber interfaces for [%s]\n", broker)
			return nil
		},
	}
	cmd.Flags().StringVar(&broker, "broker", "nats", "Message broker engine (nats, kafka, rabbitmq)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing event bus file")
	return cmd
}

func newAddMetricsCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:metrics [provider]",
		Short: "Set up the metrics middleware and endpoint (e.g. prometheus)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddMetrics(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📊 Injected [%s] metrics middleware and endpoint\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing metrics file")
	return cmd
}

func newAddTracingCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:tracing [exporter]",
		Short: "Set up the OpenTelemetry tracing infrastructure (e.g. jaeger, stdout)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddTracing(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔍 Injected OpenTelemetry tracing infrastructure for [%s]\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing tracing file")
	return cmd
}

func newAddDeployCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:deploy [target]",
		Short: "Set up the deployment manifests (k8s, helm, lambda)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddDeploy(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("☁️ Injected [%s] deployment manifests\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing deploy file")
	return cmd
}

func newAddCICDCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:cicd [provider]",
		Short: "Set up the CI/CD pipeline (github, gitlab)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddCICD(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🚀 Injected [%s] CI/CD pipeline\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing cicd file")
	return cmd
}

func newAddAICommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:ai [provider]",
		Short: "Set up the LLM proxy infrastructure (e.g. llm-proxy)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddAI(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🤖 Injected [%s] AI LLM-Proxy infrastructure\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing ai file")
	return cmd
}
