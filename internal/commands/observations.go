package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var observationsCmd = &cobra.Command{
	Use:   "observations",
	Short: "Manage observations (e.g. audit findings)",
}

var observationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List observations (filter by --audit-type to pull a SIRE/audit report)",
	RunE:  runObservationsList,
}

var observationsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single observation by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runObservationsGet,
}

var observationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new observation",
	RunE:  runObservationsCreate,
}

var observationsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update an observation",
	Args:  cobra.ExactArgs(1),
	RunE:  runObservationsUpdate,
}

var observationFlags = map[string]string{
	"user-id":               "user_id",
	"date":                  "date",
	"created-at":            "created_at",
	"vessel-id":             "vessel_id",
	"description":           "description",
	"corrective-actions":    "corrective_actions",
	"plan-of-action":        "plan_of_action",
	"status":                "status",
	"notification-group-id": "notification_group_id",
	"resolution-date":       "resolution_date",
	"manager-notes":         "manager_notes",
	"reference-number":      "reference_number",
	"assigned-to-id":        "assigned_to_id",
	"assigned-to":           "assigned_to",
	"audit-id":              "audit_id",
	"routine-id":            "routine_id",
}

var observationBoolFlags = map[string]string{
	"self-resolve":     "self_resolve",
	"modify-resolvers": "modify_resolvers",
}

var observationArrayFlags = map[string]string{
	"tags":      "tags",
	"resolvers": "resolvers",
}

func addObservationWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("user-id", "", "User ID")
	cmd.Flags().String("date", "", "Date (YYYY-MM-DD)")
	cmd.Flags().String("created-at", "", "Creation timestamp (ISO 8601); backdates the record")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("corrective-actions", "", "Corrective actions")
	cmd.Flags().String("plan-of-action", "", "Plan of action")
	cmd.Flags().String("status", "", "Status (New, In Progress, Deferred, Awaiting Approval, Resolved)")
	cmd.Flags().String("notification-group-id", "", "Notification group ID")
	cmd.Flags().String("resolution-date", "", "Resolution date (YYYY-MM-DD)")
	cmd.Flags().String("manager-notes", "", "Manager notes")
	cmd.Flags().String("reference-number", "", "Reference number")
	cmd.Flags().String("assigned-to-id", "", "Assigned-to user ID")
	cmd.Flags().String("assigned-to", "", "Assigned-to user name")
	cmd.Flags().String("audit-id", "", "Parent Audit ID (links observation to an audit)")
	cmd.Flags().String("routine-id", "", "Routine ID")
	cmd.Flags().Bool("self-resolve", false, "Self resolve")
	cmd.Flags().Bool("modify-resolvers", false, "Modify resolvers")
	cmd.Flags().String("tags", "", "Tags (comma-separated)")
	cmd.Flags().String("resolvers", "", "Resolver IDs (comma-separated)")
}

func init() {
	observationsListCmd.Flags().Int("page", 1, "Page number")
	observationsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	observationsListCmd.Flags().String("audit-type", "", "Filter by parent Audit type, e.g. \"SIRE Inspection\"")
	observationsListCmd.Flags().String("audit-id", "", "Filter by a specific Audit ID")
	observationsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	observationsListCmd.Flags().String("start-date", "", "Earliest observation date (YYYY-MM-DD)")
	observationsListCmd.Flags().String("end-date", "", "Latest observation date (YYYY-MM-DD)")
	observationsListCmd.Flags().String("status", "", "Filter by status (e.g. Resolved)")

	addObservationWriteFlags(observationsCreateCmd)
	addObservationWriteFlags(observationsUpdateCmd)

	observationsCmd.AddCommand(observationsListCmd)
	observationsCmd.AddCommand(observationsGetCmd)
	observationsCmd.AddCommand(observationsCreateCmd)
	observationsCmd.AddCommand(observationsUpdateCmd)
}

func runObservationsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters := map[string]string{
		"audit-type": "audit_type",
		"audit-id":   "audit_id",
		"vessel-id":  "vessel_id",
		"start-date": "start_date",
		"end-date":   "end_date",
		"status":     "status",
	}
	for flag, paramKey := range listFilters {
		if v, _ := cmd.Flags().GetString(flag); v != "" {
			params[paramKey] = v
		}
	}

	response, err := c.Get("observations", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops observations get %s", 3)
	env := envelope.WrapCollection(response, "observations", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runObservationsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("observations/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if auditID, ok := response["audit_id"].(string); ok && auditID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops observations list --audit-id %s", auditID))
	}
	if vesselID, ok := response["vessel_id"].(string); ok && vesselID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %s", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops observations list")

	env := envelope.WrapRecord(response, "Observation", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runObservationsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, observationFlags)
	bodyFromBoolFlags(cmd, body, observationBoolFlags)
	bodyFromArrayFlags(cmd, body, observationArrayFlags)

	response, err := c.Post("observations", wrapBody("observation", body))
	if err != nil {
		return handleClientError(err)
	}

	obsID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Observation created", []string{
		fmt.Sprintf("mobileops observations get %s", obsID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runObservationsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, observationFlags)
	bodyFromBoolFlags(cmd, body, observationBoolFlags)
	bodyFromArrayFlags(cmd, body, observationArrayFlags)

	response, err := c.Put(fmt.Sprintf("observations/%s", id), wrapBody("observation", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Observation updated", []string{
		fmt.Sprintf("mobileops observations get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
