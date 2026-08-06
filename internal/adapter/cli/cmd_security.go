package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdSecurityArgon2(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "security:argon2",
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

func newCmdSecurityAuditlog(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "security:auditlog",
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

func newCmdSecurityAuth(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "security:auth [oauth2|apikey]",
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

func newCmdSecurityAuthz(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "security:authz [casbin]",
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

func newCmdSecurityCrypto(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "security:crypto [aes-gcm]",
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

func newCmdSecurityOauth2(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool
	var provider string

	cmd := &cobra.Command{
		Use:   "security:oauth2 [provider]",
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

func newCmdSecuritySecrets(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "security:secrets [vault|aws]",
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
