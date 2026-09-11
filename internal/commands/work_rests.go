package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var workRestsCmd = &cobra.Command{
	Use:   "work-rests",
	Short: "List work/rest hour logs",
}

var workRestsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List work rest logs (one record per user per day, newest first)",
	RunE:  runWorkRestsList,
}

func init() {
	workRestsListCmd.Flags().Int("page", 1, "Page number")
	workRestsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	workRestsListCmd.Flags().String("from", "", "Earliest record date (YYYY-MM-DD)")
	workRestsListCmd.Flags().String("to", "", "Latest record date (YYYY-MM-DD)")
	workRestsListCmd.Flags().String("user-id", "", "Filter by user ID")
	workRestsListCmd.Flags().String("vessel-ids", "", "Vessel IDs (comma-separated); matches any period on one of them")
	workRestsListCmd.Flags().String("employee-position-ids", "", "Employee position IDs (comma-separated)")
	workRestsListCmd.Flags().String("employee-rate-ids", "", "Employee rate IDs (comma-separated)")

	workRestsCmd.AddCommand(workRestsListCmd)
}

func runWorkRestsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters(cmd, params, map[string]string{
		"from":                  "start_date",
		"to":                    "end_date",
		"user-id":               "user_id",
		"vessel-ids":            "asset_ids",
		"employee-position-ids": "employee_position_ids",
		"employee-rate-ids":     "employee_rate_ids",
	})

	response, err := c.Get("work-rests", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "work rest records", []string{
		"mobileops crew list",
		"mobileops employee-positions list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
