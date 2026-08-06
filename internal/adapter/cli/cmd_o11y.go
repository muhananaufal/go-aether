package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/adapter/cli/prompt"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdO11yLogger(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "o11y:logger [provider]",
		Short: "Set up structured JSON logger with context correlation tracking (slog, zap)",
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

			err = svc.AddLogger(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📝 Injected [%s] structured context logger\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing logger files")
	return cmd
}

func newCmdO11yMetrics(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "o11y:metrics [provider]",
		Short: "Set up the metrics middleware and endpoint (e.g. prometheus)",
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

			err = svc.AddMetrics(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("📊 Injected [%s] metrics middleware and endpoint\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing metrics file")
	return cmd
}

func newCmdO11yTracing(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "o11y:tracing [exporter]",
		Short: "Set up the OpenTelemetry tracing infrastructure (e.g. jaeger, stdout)",
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

			err = svc.AddTracing(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🔍 Injected OpenTelemetry tracing infrastructure for [%s]\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing tracing file")
	return cmd
}
