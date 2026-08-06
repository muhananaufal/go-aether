package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newAddSingleflightCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:singleflight",
		Short: "Set up request deduplication helper to prevent cache stampede",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddSingleflight(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🛡️ Generated Singleflight request deduplication helper in pkg/concurrency/")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing singleflight files")
	return cmd
}

func newAddDrainCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:drain",
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

func newAddOAuth2Command(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool
	var provider string

	cmd := &cobra.Command{
		Use:   "add:oauth2 [provider]",
		Short: "Set up OIDC/OAuth2 login client with PKCE state verification",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if len(args) > 0 {
				provider = args[0]
			} else {
				provider = "google"
			}

			err = svc.AddOAuth2(cmd.Context(), cwd, provider, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🔐 Generated OAuth2 SSO handler in internal/adapter/handler/http/")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing oauth2 files")
	return cmd
}

func newAddAuditLogCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:auditlog",
		Short: "Set up tamper-evident immutable audit log with PII scrubbing",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddAuditLog(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("📝 Generated AuditLog tamper-evident middleware in pkg/middleware/")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing auditlog files")
	return cmd
}

func newAddArgon2Command(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add:argon2",
		Short: "Set up GPU-resistant Argon2id password security hasher",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			err = svc.AddArgon2(cmd.Context(), cwd, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Println("🔑 Generated Argon2id password hasher in pkg/security/")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing argon2 files")
	return cmd
}
