package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var partRequestsCmd = &cobra.Command{
	Use:   "part-requests",
	Short: "List part requests (read-only)",
}

var partRequestsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List part requests",
	RunE:  runPartRequestsList,
}

var partRequestsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single part request by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runPartRequestsGet,
}

func init() {
	partRequestsListCmd.Flags().Int("page", 1, "Page number")
	partRequestsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	partRequestsCmd.AddCommand(partRequestsListCmd)
	partRequestsCmd.AddCommand(partRequestsGetCmd)
}

func runPartRequestsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get("part-requests", paginationParams(cmd))
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops part-requests get %s", 3)
	env := envelope.WrapCollection(response, "part requests", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runPartRequestsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("part-requests/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if vesselID, ok := response["vessel_id"].(string); ok && vesselID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %s", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops part-requests list")

	env := envelope.WrapRecord(response, "Part request", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}
