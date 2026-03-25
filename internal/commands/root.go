package commands

import (
	"fmt"
	"os"

	"github.com/MobileOps-Team/mobileops-cli/internal/agent"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	agentMode  bool
	envFlag    string
)

// rootCmd is the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "mobileops",
	Short: "CLI for the MobileOps API",
	Long:  "CLI for managing vessels, jobs, crew, and operations via the MobileOps API",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// If --agent and --help are both set, print agent manifest and exit
		if agentMode {
			agent.Generate()
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&agentMode, "agent", false, "Machine-readable help (use with --help)")
	rootCmd.PersistentFlags().StringVar(&envFlag, "env", "", "Environment: production, gamma, development, test")

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(vesselsCmd)
	rootCmd.AddCommand(jobsCmd)
	rootCmd.AddCommand(crewCmd)
	rootCmd.AddCommand(componentsCmd)
	rootCmd.AddCommand(workRequestsCmd)
	rootCmd.AddCommand(versionCmd)
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
