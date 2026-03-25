package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var workRequestsCmd = &cobra.Command{
	Use:   "work-requests",
	Short: "Manage work requests",
}

var workRequestsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List work requests",
	RunE:  runWorkRequestsList,
}

var workRequestsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single work request by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkRequestsGet,
}

func init() {
	workRequestsListCmd.Flags().Int("page", 1, "Page number")
	workRequestsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	workRequestsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")

	workRequestsCmd.AddCommand(workRequestsListCmd)
	workRequestsCmd.AddCommand(workRequestsGetCmd)
}

func runWorkRequestsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if vesselID, _ := cmd.Flags().GetString("vessel-id"); vesselID != "" {
		params["vessel_id"] = vesselID
	}

	response, err := c.Get("work-requests", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops work-requests get %s", 3)
	env := envelope.WrapCollection(response, "work requests", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runWorkRequestsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("work-requests/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if vesselID, ok := response["vessel_id"].(string); ok && vesselID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %s", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops work-requests list")

	env := envelope.WrapRecord(response, "Work request", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}
