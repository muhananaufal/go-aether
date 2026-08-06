package cli

import (
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

type globalFlags struct {
	Verbose bool
	DryRun  bool
}

// NewRootCommand constructs the parent cobra command and attaches all available subcommands.
func NewRootCommand(svc port.ScaffoldService) *cobra.Command {
	var flags globalFlags

	rootCmd := &cobra.Command{
		Use:   "go-aether",
		Short: "Opinionated Architecture Scaffold CLI Engine for Go Backend Engineers",
		Long: `go-aether is a lightning-fast dev-time CLI scaffolding engine built to standardize
Hexagonal / Ports & Adapters architectures, enforce zero runtime overhead, and eliminate legacy boilerplate.`,
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "enable verbose structured logging output")
	rootCmd.PersistentFlags().BoolVar(&flags.DryRun, "dry-run", false, "simulate execution without mutating physical filesystem")

	// Attach subcommands
	rootCmd.AddCommand(newLsCommand(svc, &flags))
	rootCmd.AddCommand(newInitCommand(svc, &flags))
	rootCmd.AddCommand(newAdoptCommand(svc, &flags))
	rootCmd.AddCommand(newDoctorCommand(svc, &flags))
	rootCmd.AddCommand(newMakeModuleCommand(svc, &flags))
	rootCmd.AddCommand(newMakeServiceCommand(svc, &flags))
	rootCmd.AddCommand(newMakeHandlerCommand(svc, &flags))
	rootCmd.AddCommand(newMakeDomainCommand(svc, &flags))
	rootCmd.AddCommand(newMakePortCommand(svc, &flags))
	rootCmd.AddCommand(newMakeRepositoryCommand(svc, &flags))
	rootCmd.AddCommand(newAddMiddlewareCommand(svc, &flags))
	rootCmd.AddCommand(newAddCacheCommand(svc, &flags))
	rootCmd.AddCommand(newAddTransportCommand(svc, &flags))
	rootCmd.AddCommand(newAddWorkerCommand(svc, &flags))
	rootCmd.AddCommand(newAddEventingCommand(svc, &flags))
	rootCmd.AddCommand(newAddMetricsCommand(svc, &flags))
	rootCmd.AddCommand(newAddTracingCommand(svc, &flags))
	rootCmd.AddCommand(newAddDeployCommand(svc, &flags))
	rootCmd.AddCommand(newAddCICDCommand(svc, &flags))
	rootCmd.AddCommand(newAddAICommand(svc, &flags))
	rootCmd.AddCommand(newAddDICommand(svc, &flags))
	rootCmd.AddCommand(newAddConfigCommand(svc, &flags))
	rootCmd.AddCommand(newAddErrorCommand(svc, &flags))
	rootCmd.AddCommand(newAddValidatorCommand(svc, &flags))
	rootCmd.AddCommand(newMakeMigrationCommand(svc, &flags))
	rootCmd.AddCommand(newMakeSeederCommand(svc, &flags))
	rootCmd.AddCommand(newMakeValueObjectCommand(svc, &flags))
	rootCmd.AddCommand(newMakeAggregateCommand(svc, &flags))
	rootCmd.AddCommand(newMakeEventCommand(svc, &flags))
	rootCmd.AddCommand(newMakeCommandCommand(svc, &flags))
	rootCmd.AddCommand(newMakeQueryCommand(svc, &flags))
	rootCmd.AddCommand(newAddTestCommand(svc, &flags))
	rootCmd.AddCommand(newAddMultitenancyCommand(svc, &flags))
	rootCmd.AddCommand(newAddCQRSCommand(svc, &flags))
	rootCmd.AddCommand(newAddOutboxCommand(svc, &flags))
	rootCmd.AddCommand(newAddSagaCommand(svc, &flags))
	rootCmd.AddCommand(newAddWebhookCommand(svc, &flags))
	rootCmd.AddCommand(newAddDiscoveryCommand(svc, &flags))
	rootCmd.AddCommand(newAddAuthCommand(svc, &flags))
	rootCmd.AddCommand(newAddStorageCommand(svc, &flags))
	rootCmd.AddCommand(newAddCronCommand(svc, &flags))
	rootCmd.AddCommand(newAddMailerCommand(svc, &flags))
	rootCmd.AddCommand(newAddFirebaseCommand(svc, &flags))
	rootCmd.AddCommand(newAddLoggerCommand(svc, &flags))
	rootCmd.AddCommand(newAddHealthcheckCommand(svc, &flags))
	rootCmd.AddCommand(newAddSecretsCommand(svc, &flags))
	rootCmd.AddCommand(newAddLockCommand(svc, &flags))
	rootCmd.AddCommand(newAddAuthzCommand(svc, &flags))
	rootCmd.AddCommand(newAddCryptoCommand(svc, &flags))
	rootCmd.AddCommand(newAddProfilingCommand(svc, &flags))
	rootCmd.AddCommand(newAddFeatureFlagsCommand(svc, &flags))
	rootCmd.AddCommand(newAddIdempotencyCommand(svc, &flags))
	rootCmd.AddCommand(newAddLedgerCommand(svc, &flags))
	rootCmd.AddCommand(newAddDecimalCommand(svc, &flags))
	rootCmd.AddCommand(newAddReconciliationCommand(svc, &flags))
	rootCmd.AddCommand(newAddPricingEngineCommand(svc, &flags))
	rootCmd.AddCommand(newAddWebSocketCommand(svc, &flags))
	rootCmd.AddCommand(newAddSSECommand(svc, &flags))
	rootCmd.AddCommand(newAddWebRTCCommand(svc, &flags))
	rootCmd.AddCommand(newAddMQTTCommand(svc, &flags))
	rootCmd.AddCommand(newAddTwilioCommand(svc, &flags))
	rootCmd.AddCommand(newAddMultiLevelCacheCommand(svc, &flags))
	rootCmd.AddCommand(newAddBloomFilterCommand(svc, &flags))
	rootCmd.AddCommand(newAddS3Command(svc, &flags))
	rootCmd.AddCommand(newAddResilienceCommand(svc, &flags))
	rootCmd.AddCommand(newAddSearchCommand(svc, &flags))
	rootCmd.AddCommand(newTestStressCommand(svc, &flags))
	rootCmd.AddCommand(newTestChaosCommand(svc, &flags))
	rootCmd.AddCommand(newTestFuzzCommand(svc, &flags))
	rootCmd.AddCommand(newTestBenchmarkCommand(svc, &flags))
	rootCmd.AddCommand(newTestContainerCommand(svc, &flags))
	rootCmd.AddCommand(newTestMutationCommand(svc, &flags))
	rootCmd.AddCommand(newMakeMockCommand(svc, &flags))
	rootCmd.AddCommand(newAddSQLCCommand(svc, &flags))
	rootCmd.AddCommand(newAddGRPCStreamCommand(svc, &flags))
	rootCmd.AddCommand(newAddGRPCGatewayCommand(svc, &flags))
	rootCmd.AddCommand(newAddTenantContextCommand(svc, &flags))
	rootCmd.AddCommand(newMakePipelineCommand(svc, &flags))
	rootCmd.AddCommand(newMakeSpecificationCommand(svc, &flags))
	rootCmd.AddCommand(newAddSingleflightCommand(svc, &flags))
	rootCmd.AddCommand(newAddDrainCommand(svc, &flags))
	rootCmd.AddCommand(newAddOAuth2Command(svc, &flags))
	rootCmd.AddCommand(newAddAuditLogCommand(svc, &flags))
	rootCmd.AddCommand(newAddArgon2Command(svc, &flags))

	// v0.9.0 Final Gap commands
	rootCmd.AddCommand(newAddLintCommand(svc, &flags))
	rootCmd.AddCommand(newAddUOWCommand(svc, &flags))
	rootCmd.AddCommand(newAddGraphQLCommand(svc, &flags))
	rootCmd.AddCommand(newAddReadReplicaCommand(svc, &flags))
	rootCmd.AddCommand(newAddOpenAPICommand(svc, &flags))
	rootCmd.AddCommand(newMakeCursorPaginatorCommand(svc, &flags))

	return rootCmd
}
