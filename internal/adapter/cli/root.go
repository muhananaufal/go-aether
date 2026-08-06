package cli

import (
	"fmt"

	"github.com/muhananaufal/go-aether/internal/adapter/config"
	"github.com/muhananaufal/go-aether/internal/core/port"
	"github.com/spf13/cobra"
)

type globalFlags struct {
	Verbose bool
	DryRun  bool
}

// NewRootCommand constructs the parent cobra command and attaches all available subcommands.
func NewRootCommand(svc port.ScaffoldService, version, commit, date string) *cobra.Command {
	var flags globalFlags

	rootCmd := &cobra.Command{
		Use:   "go-aether",
		Short: "Opinionated Architecture Scaffold CLI Engine for Go Backend Engineers",
		Long: `go-aether is a lightning-fast dev-time CLI scaffolding engine built to standardize
Hexagonal / Ports & Adapters architectures, enforce zero runtime overhead, and eliminate legacy boilerplate.`,
		Version:       fmt.Sprintf("%s\ncommit: %s\nbuilt at: %s", version, commit, date),
		SilenceUsage:  true,
		SilenceErrors: false,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil || cfg == nil {
				return nil
			}

			// Apply context memory to flags if they were not explicitly provided
			if !cmd.Flags().Changed("engine") && cfg.Preferences.Engine != "" && cmd.Flags().Lookup("engine") != nil {
				_ = cmd.Flags().Set("engine", cfg.Preferences.Engine)
			}
			if !cmd.Flags().Changed("orm") && cfg.Preferences.ORM != "" && cmd.Flags().Lookup("orm") != nil {
				_ = cmd.Flags().Set("orm", cfg.Preferences.ORM)
			}
			if !cmd.Flags().Changed("architecture") && cfg.Preferences.Architecture != "" && cmd.Flags().Lookup("architecture") != nil {
				_ = cmd.Flags().Set("architecture", cfg.Preferences.Architecture)
			}
			return nil
		},
	}

	rootCmd.SetVersionTemplate("go-aether version {{.Version}}\n")

	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "enable verbose structured logging output")
	rootCmd.PersistentFlags().BoolVar(&flags.DryRun, "dry-run", false, "simulate execution without mutating physical filesystem")

	// Attach subcommands
	rootCmd.AddCommand(newCmdAPIGRPC(svc, &flags))
	rootCmd.AddCommand(newCmdAPIGRPCGateway(svc, &flags))
	rootCmd.AddCommand(newCmdAPIGraphql(svc, &flags))
	rootCmd.AddCommand(newCmdAPIIdempotency(svc, &flags))
	rootCmd.AddCommand(newCmdAPIMiddleware(svc, &flags))
	rootCmd.AddCommand(newCmdAPIOpenapi(svc, &flags))
	rootCmd.AddCommand(newCmdAPITransport(svc, &flags))
	rootCmd.AddCommand(newCmdAPIValidator(svc, &flags))
	rootCmd.AddCommand(newCmdAdopt(svc, &flags))
	rootCmd.AddCommand(newCmdArchAggregate(svc, &flags))
	rootCmd.AddCommand(newCmdArchCommand(svc, &flags))
	rootCmd.AddCommand(newCmdArchDI(svc, &flags))
	rootCmd.AddCommand(newCmdArchDomain(svc, &flags))
	rootCmd.AddCommand(newCmdArchError(svc, &flags))
	rootCmd.AddCommand(newCmdArchEvent(svc, &flags))
	rootCmd.AddCommand(newCmdArchHandler(svc, &flags))
	rootCmd.AddCommand(newCmdArchMock(svc, &flags))
	rootCmd.AddCommand(newCmdArchModule(svc, &flags))
	rootCmd.AddCommand(newCmdArchPipeline(svc, &flags))
	rootCmd.AddCommand(newCmdArchPort(svc, &flags))
	rootCmd.AddCommand(newCmdArchQuery(svc, &flags))
	rootCmd.AddCommand(newCmdArchRepository(svc, &flags))
	rootCmd.AddCommand(newCmdArchService(svc, &flags))
	rootCmd.AddCommand(newCmdArchSpecification(svc, &flags))
	rootCmd.AddCommand(newCmdArchValueobject(svc, &flags))
	rootCmd.AddCommand(newCmdAsyncCron(svc, &flags))
	rootCmd.AddCommand(newCmdAsyncEventing(svc, &flags))
	rootCmd.AddCommand(newCmdAsyncOutbox(svc, &flags))
	rootCmd.AddCommand(newCmdAsyncSaga(svc, &flags))
	rootCmd.AddCommand(newCmdAsyncWebhook(svc, &flags))
	rootCmd.AddCommand(newCmdAsyncWorker(svc, &flags))
	rootCmd.AddCommand(newCmdCacheBloom(svc, &flags))
	rootCmd.AddCommand(newCmdCacheDedup(svc, &flags))
	rootCmd.AddCommand(newCmdCacheMultilevel(svc, &flags))
	rootCmd.AddCommand(newCmdCacheRedis(svc, &flags))
	rootCmd.AddCommand(newCmdCloudS3(svc, &flags))
	rootCmd.AddCommand(newCmdCloudStorage(svc, &flags))
	rootCmd.AddCommand(newCmdDBMigration(svc, &flags))
	rootCmd.AddCommand(newCmdDBPaginator(svc, &flags))
	rootCmd.AddCommand(newCmdDBReadreplica(svc, &flags))
	rootCmd.AddCommand(newCmdDBSQLC(svc, &flags))
	rootCmd.AddCommand(newCmdDBSeeder(svc, &flags))
	rootCmd.AddCommand(newCmdDBUOW(svc, &flags))
	rootCmd.AddCommand(newCmdDoctor(svc, &flags))
	rootCmd.AddCommand(newCmdFintechDecimal(svc, &flags))
	rootCmd.AddCommand(newCmdFintechLedger(svc, &flags))
	rootCmd.AddCommand(newCmdFintechPricing(svc, &flags))
	rootCmd.AddCommand(newCmdFintechReconcile(svc, &flags))
	rootCmd.AddCommand(newCmdInfraCICD(svc, &flags))
	rootCmd.AddCommand(newCmdInfraConfig(svc, &flags))
	rootCmd.AddCommand(newCmdInfraDeploy(svc, &flags))
	rootCmd.AddCommand(newCmdInfraDrain(svc, &flags))
	rootCmd.AddCommand(newCmdInfraFeatureflags(svc, &flags))
	rootCmd.AddCommand(newCmdInfraHealthcheck(svc, &flags))
	rootCmd.AddCommand(newCmdInfraLint(svc, &flags))
	rootCmd.AddCommand(newCmdInfraProfiling(svc, &flags))
	rootCmd.AddCommand(newCmdInfraSearch(svc, &flags))
	rootCmd.AddCommand(newCmdInit(svc, &flags))
	rootCmd.AddCommand(newCmdLs(svc, &flags))
	rootCmd.AddCommand(newCmdNotifMail(svc, &flags))
	rootCmd.AddCommand(newCmdNotifPush(svc, &flags))
	rootCmd.AddCommand(newCmdNotifSms(svc, &flags))
	rootCmd.AddCommand(newCmdO11yLogger(svc, &flags))
	rootCmd.AddCommand(newCmdO11yMetrics(svc, &flags))
	rootCmd.AddCommand(newCmdO11yTracing(svc, &flags))
	rootCmd.AddCommand(newCmdPlatformAI(svc, &flags))
	rootCmd.AddCommand(newCmdPlatformCQRS(svc, &flags))
	rootCmd.AddCommand(newCmdPlatformDiscovery(svc, &flags))
	rootCmd.AddCommand(newCmdPlatformLock(svc, &flags))
	rootCmd.AddCommand(newCmdPlatformMultitenancy(svc, &flags))
	rootCmd.AddCommand(newCmdPlatformResilience(svc, &flags))
	rootCmd.AddCommand(newCmdPlatformTenant(svc, &flags))
	rootCmd.AddCommand(newCmdRealtimeMQTT(svc, &flags))
	rootCmd.AddCommand(newCmdRealtimeSSE(svc, &flags))
	rootCmd.AddCommand(newCmdRealtimeWS(svc, &flags))
	rootCmd.AddCommand(newCmdRealtimeWebRTC(svc, &flags))
	rootCmd.AddCommand(newCmdSecurityArgon2(svc, &flags))
	rootCmd.AddCommand(newCmdSecurityAuditlog(svc, &flags))
	rootCmd.AddCommand(newCmdSecurityAuth(svc, &flags))
	rootCmd.AddCommand(newCmdSecurityAuthz(svc, &flags))
	rootCmd.AddCommand(newCmdSecurityCrypto(svc, &flags))
	rootCmd.AddCommand(newCmdSecurityOauth2(svc, &flags))
	rootCmd.AddCommand(newCmdSecuritySecrets(svc, &flags))
	rootCmd.AddCommand(newCmdTestBenchmark(svc, &flags))
	rootCmd.AddCommand(newCmdTestChaos(svc, &flags))
	rootCmd.AddCommand(newCmdTestContainer(svc, &flags))
	rootCmd.AddCommand(newCmdTestFuzz(svc, &flags))
	rootCmd.AddCommand(newCmdTestIntegration(svc, &flags))
	rootCmd.AddCommand(newCmdTestMutation(svc, &flags))
	rootCmd.AddCommand(newCmdTestStress(svc, &flags))
	return rootCmd
}
