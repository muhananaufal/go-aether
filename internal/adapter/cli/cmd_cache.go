package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdCacheBloom(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "cache:bloom",
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

func newCmdCacheDedup(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "cache:dedup",
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

func newCmdCacheMultilevel(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "cache:multilevel",
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

func newCmdCacheRedis(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "cache:redis [cache-type]",
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
