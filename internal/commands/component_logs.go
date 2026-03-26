package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var componentLogsCmd = &cobra.Command{
	Use:   "component-logs",
	Short: "List component hour logs",
}

var componentLogsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List component hour logs",
	RunE:  runComponentLogsList,
}

func init() {
	componentLogsListCmd.Flags().Int("page", 1, "Page number")
	componentLogsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	componentLogsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	componentLogsListCmd.Flags().String("created-after", "", "Created after date (YYYY-MM-DD)")
	componentLogsListCmd.Flags().String("created-before", "", "Created before date (YYYY-MM-DD)")
	componentLogsListCmd.Flags().String("updated-after", "", "Updated after date (YYYY-MM-DD)")
	componentLogsListCmd.Flags().String("updated-before", "", "Updated before date (YYYY-MM-DD)")
	componentLogsListCmd.Flags().Bool("sum", false, "Aggregate hours by component")

	componentLogsCmd.AddCommand(componentLogsListCmd)
}

func runComponentLogsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("vessel-id"); v != "" {
		params["vessel_id"] = v
	}
	if v, _ := cmd.Flags().GetString("created-after"); v != "" {
		params["created_at_gte"] = v
	}
	if v, _ := cmd.Flags().GetString("created-before"); v != "" {
		params["created_at_lte"] = v
	}
	if v, _ := cmd.Flags().GetString("updated-after"); v != "" {
		params["updated_at_gte"] = v
	}
	if v, _ := cmd.Flags().GetString("updated-before"); v != "" {
		params["updated_at_lte"] = v
	}
	if v, _ := cmd.Flags().GetBool("sum"); v {
		params["sum"] = "true"
	}

	response, err := c.Get("component-logs", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "component hour logs", []string{
		"mobileops components list",
		"mobileops vessels list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
