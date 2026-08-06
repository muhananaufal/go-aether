package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/adapter/cli/prompt"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdCloudS3(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "cloud:s3 [minio|aws]",
		Short: "Set up S3 object storage client with pre-signed URL generator",
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

			err = svc.AddS3(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🪣 Injected [%s] S3 object storage client\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing storage files")
	return cmd
}

func newCmdCloudStorage(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "cloud:storage [provider]",
		Short: "Set up cloud blob storage adapter (s3, gcs, local)",
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

			err = svc.AddStorage(cmd.Context(), cwd, arg0, globals.DryRun, force)
			if err != nil {
				return err
			}

			fmt.Printf("🗄️ Injected [%s] cloud storage adapter\n", arg0)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing storage files")
	return cmd
}
