package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var incidentsCmd = &cobra.Command{
	Use:   "incidents",
	Short: "View incident reports",
	Long: "View incident reports with their injury, lost-time and notification flags,\n" +
		"witness statements and comments. Archived reports are included unless\n" +
		"--archived is set.",
}

var incidentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List incidents, most recent incident date first",
	RunE:  runIncidentsList,
}

var incidentsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single incident by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runIncidentsGet,
}

var incidentListFlags = map[string]string{
	"vessel-id":     "vessel_id",
	"from":          "start_date",
	"to":            "end_date",
	"updated-since": "updated_since",
}

var incidentListBoolFlags = map[string]string{
	"archived": "archived",
}

func init() {
	incidentsListCmd.Flags().Int("page", 1, "Page number")
	incidentsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	incidentsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	incidentsListCmd.Flags().String("from", "", "Incidents that occurred on or after this date (YYYY-MM-DD)")
	incidentsListCmd.Flags().String("to", "", "Incidents that occurred on or before this date (YYYY-MM-DD)")
	incidentsListCmd.Flags().String("updated-since", "", "Only incidents updated at or after this time (ISO 8601)")
	incidentsListCmd.Flags().Bool("archived", false, "Only archived incidents; --archived=false for only unarchived")

	incidentsCmd.AddCommand(incidentsListCmd)
	incidentsCmd.AddCommand(incidentsGetCmd)
}

func runIncidentsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters(cmd, params, incidentListFlags)
	boolListFilters(cmd, params, incidentListBoolFlags)

	response, err := c.Get("incidents", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops incidents get %s", 3)
	env := envelope.WrapCollection(response, "incidents", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runIncidentsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("incidents/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := []string{}
	if vesselID, ok := response["vessel_id"].(float64); ok {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %.0f", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops incidents list")

	env := envelope.WrapRecord(response, "Incident", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}
