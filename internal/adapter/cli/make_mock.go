package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newMakeMockCommand(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "arch:mock [interface-name]",
		Short: "Scaffold interface mock implementation using Mockery directives for isolated unit tests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaceName := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := svc.MakeMock(cmd.Context(), cwd, interfaceName, globals.DryRun, force); err != nil {
				return err
			}

			fmt.Printf("🎭 Generated Mockery mock contract for [%s] in mocks/\n", interfaceName)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing mock file")

	return cmd
}
