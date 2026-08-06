package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdAsyncCron(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "async:cron [job-name]",
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

func newCmdAsyncEventing(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var broker string
	var force bool

	cmd := &cobra.Command{
		Use:   "async:eventing",
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

func newCmdAsyncOutbox(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "async:outbox",
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

func newCmdAsyncSaga(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "async:saga [workflow-name]",
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

func newCmdAsyncWebhook(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "async:webhook",
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

func newCmdAsyncWorker(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var broker, pattern string
	var force bool

	cmd := &cobra.Command{
		Use:   "async:worker [worker-name]",
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
