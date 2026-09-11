package commands

import (
	"fmt"
	"os"

	"github.com/MobileOps-Team/mobileops-cli/internal/agent"
	"github.com/MobileOps-Team/mobileops-cli/internal/updater"
	"github.com/spf13/cobra"
)

var (
	jsonOutput     bool
	agentMode      bool
	envFlag        string
	versionCheckCh <-chan string
)

// rootCmd is the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "mobileops",
	Short: "CLI for the MobileOps API",
	Long:  "CLI for managing vessels, jobs, crew, and operations via the MobileOps API",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "update" && os.Getenv("MOBILEOPS_SKIP_VERSION_CHECK") == "" {
			versionCheckCh = updater.CheckVersionBackground(cliVersion)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if versionCheckCh == nil {
			return
		}
		latest := <-versionCheckCh
		if latest == "" {
			return
		}

		if updater.AutoUpdateEnabled() && updater.CanSelfUpdate() {
			if err := updater.AutoUpdate(latest); err == nil {
				fmt.Fprintf(os.Stderr, "\n\u2713 mobileops updated to v%s (takes effect on the next command)\n", latest)
				updateSkillQuiet()
				return
			}
		}

		fmt.Fprintf(os.Stderr, "\n\u26a0 Update available: v%s \u2192 v%s\n  Run: mobileops update\n", cliVersion, latest)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&agentMode, "agent", false, "Machine-readable help (use with --help)")
	rootCmd.PersistentFlags().StringVar(&envFlag, "env", "", "Environment: production, gamma, development, test")

	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if agentMode {
			agent.Generate(cmd.Root(), apiCoverage)
			return
		}
		defaultHelp(cmd, args)
	})

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(vesselsCmd)
	rootCmd.AddCommand(jobsCmd)
	rootCmd.AddCommand(crewCmd)
	rootCmd.AddCommand(componentsCmd)
	rootCmd.AddCommand(partsCmd)
	rootCmd.AddCommand(workRequestsCmd)
	rootCmd.AddCommand(deficienciesCmd)
	rootCmd.AddCommand(nonconformitiesCmd)
	rootCmd.AddCommand(observationsCmd)
	rootCmd.AddCommand(maintenanceReportsCmd)
	rootCmd.AddCommand(vesselDocumentsCmd)
	rootCmd.AddCommand(personnelDocumentsCmd)
	rootCmd.AddCommand(vesselSpecsCmd)
	rootCmd.AddCommand(vesselSpecTemplatesCmd)
	rootCmd.AddCommand(customersCmd)
	rootCmd.AddCommand(suppliersCmd)
	rootCmd.AddCommand(makesCmd)
	rootCmd.AddCommand(modelsCmd)
	rootCmd.AddCommand(divisionsCmd)
	rootCmd.AddCommand(employeePositionsCmd)
	rootCmd.AddCommand(vesselTypesCmd)
	rootCmd.AddCommand(workTypesCmd)
	rootCmd.AddCommand(locationsCmd)
	rootCmd.AddCommand(invoiceStatementsCmd)
	rootCmd.AddCommand(formInstancesCmd)
	rootCmd.AddCommand(formsCmd)
	rootCmd.AddCommand(routineTemplatesCmd)
	rootCmd.AddCommand(routineCalculationsCmd)
	rootCmd.AddCommand(positionReportsCmd)
	rootCmd.AddCommand(componentLogsCmd)
	rootCmd.AddCommand(partRequestsCmd)
	rootCmd.AddCommand(termsCmd)
	rootCmd.AddCommand(purchaseOrdersCmd)
	rootCmd.AddCommand(wheelhouseLogsCmd)
	rootCmd.AddCommand(cargoTypesCmd)
	rootCmd.AddCommand(auditsCmd)
	rootCmd.AddCommand(eventsCmd)
	rootCmd.AddCommand(workRestsCmd)
	rootCmd.AddCommand(codesCmd)
	rootCmd.AddCommand(fuelLevelReadingsCmd)
	rootCmd.AddCommand(bunkerPartitionsCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(treeCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// paginationParams extracts page and limit from command flags, clamping to valid ranges.
func paginationParams(cmd *cobra.Command) map[string]string {
	page, _ := cmd.Flags().GetInt("page")
	limit, _ := cmd.Flags().GetInt("limit")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	return map[string]string{
		"page":  intToStr(page),
		"limit": intToStr(limit),
	}
}

func intToStr(i int) string {
	return fmt.Sprintf("%d", i)
}
