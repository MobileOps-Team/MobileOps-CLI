package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var bunkerPartitionsCmd = &cobra.Command{
	Use:   "bunker-partitions",
	Short: "List bunker partitions (fuel lifts and the consumption drawn from them)",
}

var bunkerPartitionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List bunker partitions, oldest lift first",
	RunE:  runBunkerPartitionsList,
}

func init() {
	bunkerPartitionsListCmd.Flags().Int("page", 1, "Page number")
	bunkerPartitionsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	bunkerPartitionsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	bunkerPartitionsListCmd.Flags().String("job-id", "", "Only partitions with consumption attributed to this job")
	bunkerPartitionsListCmd.Flags().String("from", "", "Earliest lift timestamp (ISO 8601, company time zone)")
	bunkerPartitionsListCmd.Flags().String("to", "", "Latest lift timestamp (ISO 8601, company time zone)")

	bunkerPartitionsCmd.AddCommand(bunkerPartitionsListCmd)
}

func runBunkerPartitionsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters(cmd, params, map[string]string{
		"vessel-id": "asset_id",
		"job-id":    "job_id",
		"from":      "start_time",
		"to":        "end_time",
	})

	response, err := c.Get("bunker-partitions", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "bunker partitions", []string{
		"mobileops fuel-level-readings list",
		"mobileops vessels list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
