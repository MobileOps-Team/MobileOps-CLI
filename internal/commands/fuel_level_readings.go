package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var fuelLevelReadingsCmd = &cobra.Command{
	Use:   "fuel-level-readings",
	Short: "List fuel tank level readings",
}

var fuelLevelReadingsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List fuel level readings, oldest first",
	RunE:  runFuelLevelReadingsList,
}

func init() {
	fuelLevelReadingsListCmd.Flags().Int("page", 1, "Page number")
	fuelLevelReadingsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	fuelLevelReadingsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	fuelLevelReadingsListCmd.Flags().String("job-id", "", "Filter by the job the reading was recorded on")
	fuelLevelReadingsListCmd.Flags().String("from", "", "Earliest reading time (ISO 8601)")
	fuelLevelReadingsListCmd.Flags().String("to", "", "Latest reading time (ISO 8601)")

	fuelLevelReadingsCmd.AddCommand(fuelLevelReadingsListCmd)
}

func runFuelLevelReadingsList(cmd *cobra.Command, args []string) error {
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

	response, err := c.Get("fuel-level-readings", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "fuel level readings", []string{
		"mobileops bunker-partitions list",
		"mobileops vessels list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
