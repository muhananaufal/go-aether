package cli

import (
	"fmt"
	"os"

	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

func newCmdTestBenchmark(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "test:benchmark",
		Short: "Scaffold micro-benchmark suite and memory allocation profiler (testing.B)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := svc.TestBenchmark(cmd.Context(), cwd, globals.DryRun, force); err != nil {
				return err
			}

			fmt.Println("⏱️ Generated micro-benchmark suite in tests/bench/alloc_benchmark_test.go")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing benchmark file")

	return cmd
}

func newCmdTestChaos(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "test:chaos",
		Short: "Scaffold chaos engineering fault injection middleware (latency, abort, error rate)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := svc.TestChaos(cmd.Context(), cwd, globals.DryRun, force); err != nil {
				return err
			}

			fmt.Println("🌪️ Generated chaos engineering fault injection middleware in pkg/middleware/chaos.go")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing middleware file")

	return cmd
}

func newCmdTestContainer(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "test:container",
		Short: "Scaffold Testcontainers integration testing harness for PostgreSQL and Redis in Docker",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := svc.TestContainer(cmd.Context(), cwd, "postgres-redis", globals.DryRun, force); err != nil {
				return err
			}

			fmt.Println("🐳 Generated Testcontainers integration testing harness in tests/integration/postgres_redis_test.go")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing integration test file")

	return cmd
}

func newCmdTestFuzz(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "test:fuzz",
		Short: "Scaffold Go native continuous fuzz testing harness",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := svc.TestFuzz(cmd.Context(), cwd, "json", globals.DryRun, force); err != nil {
				return err
			}

			fmt.Println("🧬 Generated continuous fuzz testing harness in tests/fuzz/fuzz_test.go")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing fuzz test file")

	return cmd
}

func newCmdTestIntegration(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "test:integration",
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

func newCmdTestMutation(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "test:mutation",
		Short: "Scaffold mutation testing verification runner script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := svc.TestMutation(cmd.Context(), cwd, globals.DryRun, force); err != nil {
				return err
			}

			fmt.Println("🧪 Generated mutation testing runner script in scripts/mutation_test.ps1")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing script file")

	return cmd
}

func newCmdTestStress(svc port.ScaffoldService, globals *globalFlags) *cobra.Command {
	var engine string
	var force bool

	cmd := &cobra.Command{
		Use:   "test:stress",
		Short: "Scaffold high-concurrency load and stress testing suite (k6 / vegeta)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := svc.TestStress(cmd.Context(), cwd, engine, globals.DryRun, force); err != nil {
				return err
			}

			fmt.Printf("🚀 Generated high-throughput stress testing suite using [%s] in tests/stress/\n", engine)
			return nil
		},
	}

	cmd.Flags().StringVarP(&engine, "engine", "e", "k6", "Load testing engine (k6 or vegeta)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing test files")

	return cmd
}
