package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var nearMissesCmd = &cobra.Command{
	Use:   "near-misses",
	Short: "View Near Miss Hazard Observation reports",
	Long: "View Near Miss Hazard Observation reports (near misses, hazard observations,\n" +
		"stop-work events). Archived reports are included unless --archived is set.",
}

var nearMissesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List near misses, most recent observation date first",
	RunE:  runNearMissesList,
}

var nearMissesGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single near miss by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runNearMissesGet,
}

var nearMissListFlags = map[string]string{
	"vessel-id":     "vessel_id",
	"from":          "start_date",
	"to":            "end_date",
	"updated-since": "updated_since",
}

var nearMissListBoolFlags = map[string]string{
	"archived": "archived",
}

func init() {
	nearMissesListCmd.Flags().Int("page", 1, "Page number")
	nearMissesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	nearMissesListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	nearMissesListCmd.Flags().String("from", "", "Reports observed on or after this date (YYYY-MM-DD)")
	nearMissesListCmd.Flags().String("to", "", "Reports observed on or before this date (YYYY-MM-DD)")
	nearMissesListCmd.Flags().String("updated-since", "", "Only reports updated at or after this time (ISO 8601)")
	nearMissesListCmd.Flags().Bool("archived", false, "Only archived reports; --archived=false for only unarchived")

	nearMissesCmd.AddCommand(nearMissesListCmd)
	nearMissesCmd.AddCommand(nearMissesGetCmd)
}

func runNearMissesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters(cmd, params, nearMissListFlags)
	boolListFilters(cmd, params, nearMissListBoolFlags)

	response, err := c.Get("near-misses", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops near-misses get %s", 3)
	env := envelope.WrapCollection(response, "near-misses", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runNearMissesGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("near-misses/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := []string{}
	if vesselID, ok := response["vessel_id"].(float64); ok {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %.0f", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops near-misses list")

	env := envelope.WrapRecord(response, "NearMiss", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}
