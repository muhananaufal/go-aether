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

func newAddIdempotencyCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:idempotency [redis|memory]",
		Short: "Set up Idempotency-Key validation middleware",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddIdempotency(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔁 Injected [%s] idempotency key engine\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing idempotency files")
	return cmd
}

func newAddLedgerCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:ledger",
		Short: "Set up Double-Entry bookkeeping ledger engine",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddLedger(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("📒 Injected Double-Entry bookkeeping ledger engine")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing ledger files")
	return cmd
}

func newAddDecimalCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:decimal",
		Short: "Set up high-precision financial decimal money helpers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddDecimal(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("💰 Injected high-precision decimal money arithmetic helpers")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing money files")
	return cmd
}

func newAddReconciliationCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:reconciliation",
		Short: "Set up automated settlement & transaction reconciliation engine",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddReconciliation(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("⚖️ Injected automated transaction reconciliation matching engine")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing reconciliation files")
	return cmd
}

func newAddPricingEngineCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:pricing-engine",
		Short: "Set up rule-based tiered pricing and fee calculation engine",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddPricingEngine(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🏷️ Injected rule-based tiered pricing engine")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing pricing files")
	return cmd
}

func newAddWebSocketCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:websocket [gorilla|nhooyr]",
		Short: "Set up WebSocket hub and connection pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddWebSocket(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔌 Injected [%s] WebSocket hub engine\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing websocket files")
	return cmd
}

func newAddSSECommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:sse",
		Short: "Set up Server-Sent Events (SSE) live streaming broker",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddSSE(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("📡 Injected Server-Sent Events (SSE) streaming broker")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing sse files")
	return cmd
}

func newAddWebRTCCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:webrtc [pion]",
		Short: "Set up Pion WebRTC peer-to-peer data channel hub",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddWebRTC(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📹 Injected [%s] WebRTC signaling and data channel engine\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing webrtc files")
	return cmd
}

func newAddMQTTCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:mqtt [paho]",
		Short: "Set up Paho MQTT client for IoT telemetry pub/sub",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddMQTT(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📶 Injected [%s] MQTT IoT pub/sub client\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing mqtt files")
	return cmd
}

func newAddTwilioCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:twilio",
		Short: "Set up Twilio SMS & WhatsApp omni-channel messaging client",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddTwilio(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("💬 Injected Twilio SMS & WhatsApp messaging client")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing twilio files")
	return cmd
}

func newAddMultiLevelCacheCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:multilevelcache",
		Short: "Set up synchronized L1 memory + L2 Redis cache",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddMultiLevelCache(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🚀 Injected synchronized Multi-Level (L1 Memory + L2 Redis) cache engine")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing cache files")
	return cmd
}

func newAddBloomFilterCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:bloomfilter",
		Short: "Set up probabilistic Bloom Filter cache penetration guard",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddBloomFilter(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🌸 Injected probabilistic Bloom Filter cache guard")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing bloom files")
	return cmd
}

func newAddS3Command(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:s3 [minio|aws]",
		Short: "Set up S3 object storage client with pre-signed URL generator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddS3(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🪣 Injected [%s] S3 object storage client\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing storage files")
	return cmd
}

func newAddResilienceCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:resilience [hystrix|resilience4go]",
		Short: "Set up Circuit Breaker and Bulkhead resilience engine",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddResilience(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🛡️ Injected [%s] Circuit Breaker and Bulkhead resilience engine\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing resilience files")
	return cmd
}

func newAddSearchCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:search [meilisearch|elasticsearch]",
		Short: "Set up full-text search engine client",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddSearch(cmd.Context(), cwd, args[0], globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔍 Injected [%s] full-text search engine client\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing search files")
	return cmd
}

func newAddSQLCCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:sqlc",
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

func newAddGRPCStreamCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:grpc-stream",
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

func newAddGRPCGatewayCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:grpc-gateway",
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

func newAddTenantContextCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:tenant-context",
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
