package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var fuelRemovalsCmd = &cobra.Command{
	Use:   "fuel-removals",
	Short: "List fuel taken off a vessel without being burned",
}

var fuelRemovalsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List fuel removals, oldest first",
	RunE:  runFuelRemovalsList,
}

func init() {
	fuelRemovalsListCmd.Flags().Int("page", 1, "Page number")
	fuelRemovalsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	fuelRemovalsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	fuelRemovalsListCmd.Flags().String("job-id", "", "Filter by the job the removal was recorded on")
	fuelRemovalsListCmd.Flags().String("from", "", "Earliest removal time (ISO 8601)")
	fuelRemovalsListCmd.Flags().String("to", "", "Latest removal time (ISO 8601)")

	fuelRemovalsCmd.AddCommand(fuelRemovalsListCmd)
}

func runFuelRemovalsList(cmd *cobra.Command, args []string) error {
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

	response, err := c.Get("fuel-removals", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "fuel removals", []string{
		"mobileops bunker-partitions list",
		"mobileops fuel-level-readings list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
