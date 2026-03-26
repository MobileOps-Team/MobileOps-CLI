package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var wheelhouseLogsCmd = &cobra.Command{
	Use:   "wheelhouse-logs",
	Short: "List wheelhouse logs",
}

var wheelhouseLogsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List wheelhouse logs",
	RunE:  runWheelhouseLogsList,
}

func init() {
	wheelhouseLogsListCmd.Flags().Int("page", 1, "Page number")
	wheelhouseLogsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	wheelhouseLogsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	wheelhouseLogsListCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	wheelhouseLogsListCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")

	wheelhouseLogsCmd.AddCommand(wheelhouseLogsListCmd)
}

func runWheelhouseLogsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("vessel-id"); v != "" {
		params["asset_id"] = v
	}
	if v, _ := cmd.Flags().GetString("from"); v != "" {
		params["date_gte"] = v
	}
	if v, _ := cmd.Flags().GetString("to"); v != "" {
		params["date_lte"] = v
	}

	response, err := c.Get("wheelhouse-logs", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "wheelhouse logs", []string{
		"mobileops vessels list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
