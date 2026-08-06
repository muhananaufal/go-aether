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

func newAddDICommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:di [di-type]",
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

func newAddConfigCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:config [config-type]",
		Short: "Set up a centralized configuration manager (e.g. viper, koanf)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddConfig(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("⚙️ Injected [%s] centralized configuration manager\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing config file")
	return cmd
}

func newAddErrorCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:error",
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

func newAddValidatorCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:validator [validator-type]",
		Short: "Set up the struct validation wrapper (e.g. playground)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddValidator(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("✅ Injected [%s] struct validation wrapper\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing validator file")
	return cmd
}

func newAddTestCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:test",
		Short: "Set up the integration test helpers and mocking base",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddTest(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🧪 Injected integration test setup")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing test file")
	return cmd
}

func newAddMultitenancyCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:multitenancy [module-name]",
		Short: "Set up Row Level Security (RLS) SQL policies for tenant isolation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddMultitenancy(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🏢 Injected multitenancy RLS script for [%s]\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing multitenancy file")
	return cmd
}

func newAddCQRSCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:cqrs [module-name]",
		Short: "Set up CQRS Command and Query handlers within module scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddCQRS(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("⚡ Injected CQRS handlers for [%s]\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing CQRS files")
	return cmd
}

func newAddOutboxCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:outbox",
		Short: "Set up Transactional Outbox pattern infrastructure and SQL migrations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddOutbox(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("📦 Injected Transactional Outbox pattern & migrations")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing outbox files")
	return cmd
}

func newAddSagaCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:saga [workflow-name]",
		Short: "Set up Distributed Saga orchestrator and compensation workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddSaga(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔄 Injected Saga workflow for [%s]\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing saga files")
	return cmd
}

func newAddWebhookCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:webhook",
		Short: "Set up secure HMAC-SHA256 signed webhook dispatcher & receiver",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddWebhook(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🪝 Injected secure webhook dispatcher & receiver")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing webhook files")
	return cmd
}

func newAddDiscoveryCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:discovery [provider]",
		Short: "Set up service discovery client (consul, etcd)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddDiscovery(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🌐 Injected [%s] service discovery client\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing discovery files")
	return cmd
}

func newAddAuthCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:auth [oauth2|apikey]",
		Short: "Set up authentication handlers and middleware (oauth2, apikey)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddAuth(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔑 Injected [%s] authentication provider\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing auth files")
	return cmd
}

func newAddStorageCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:storage [provider]",
		Short: "Set up cloud blob storage adapter (s3, gcs, local)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddStorage(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🗄️ Injected [%s] cloud storage adapter\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing storage files")
	return cmd
}

func newAddCronCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:cron [job-name]",
		Short: "Set up in-process recurring background cron job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddCron(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("⏰ Injected cron job scheduler & [%s] job runner\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing cron files")
	return cmd
}

func newAddMailerCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:mailer [provider]",
		Short: "Set up transactional email delivery client (smtp, resend)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddMailer(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📧 Injected [%s] transactional mailer client\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing mailer files")
	return cmd
}

func newAddFirebaseCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:firebase",
		Short: "Set up Firebase Auth token decoding and FCM push messaging",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddFirebase(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🔥 Injected Firebase Auth & FCM Cloud Messaging client")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing firebase files")
	return cmd
}

func newAddLoggerCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:logger [provider]",
		Short: "Set up structured JSON logger with context correlation tracking (slog, zap)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddLogger(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📝 Injected [%s] structured context logger\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing logger files")
	return cmd
}

func newAddHealthcheckCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:healthcheck",
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

func newAddSecretsCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:secrets [vault|aws]",
		Short: "Set up secret manager client (vault, aws)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddSecrets(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔐 Injected [%s] secrets manager client\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing secrets files")
	return cmd
}

func newAddLockCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:lock [redis]",
		Short: "Set up distributed lock (Redlock) engine",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddLock(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔒 Injected [%s] distributed mutex lock engine\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing lock files")
	return cmd
}

func newAddAuthzCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:authz [casbin]",
		Short: "Set up RBAC / ABAC authorization engine (casbin)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddAuthz(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🛡️ Injected [%s] RBAC/ABAC authorization engine\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing authz files")
	return cmd
}

func newAddCryptoCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:crypto [aes-gcm]",
		Short: "Set up symmetric envelope encryption helper (aes-gcm)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddCrypto(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔑 Injected [%s] envelope encryption helper\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing crypto files")
	return cmd
}

func newAddProfilingCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:profiling [pprof]",
		Short: "Set up protected runtime profiling endpoints (pprof)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddProfiling(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📈 Injected [%s] protected profiling endpoints\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing profiling files")
	return cmd
}

func newAddFeatureFlagsCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:featureflags [flipt]",
		Short: "Set up feature flags & canary release client (flipt)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddFeatureFlags(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🚩 Injected [%s] feature flags client\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing featureflags files")
	return cmd
}
