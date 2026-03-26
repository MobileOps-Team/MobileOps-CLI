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

var workRequestsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new work request",
	RunE:  runWorkRequestsCreate,
}

var workRequestsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a work request",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkRequestsUpdate,
}

var workRequestFlags = map[string]string{
	"user-id":                "user_id",
	"date":                   "date",
	"resolution-date":        "resolution_date",
	"vessel-id":              "vessel_id",
	"component-id":           "component_id",
	"part-id":                "part_id",
	"description":            "description",
	"status":                 "status",
	"plan-of-action":         "plan_of_action",
	"priority":               "priority",
	"notification-group-id":  "notification_group_id",
	"location":               "location",
	"manager-notes":          "manager_notes",
	"code-id":                "code_id",
	"reference-number":       "reference_number",
}

var workRequestBoolFlags = map[string]string{
	"self-resolve":     "self_resolve",
	"modify-resolvers": "modify_resolvers",
}

var workRequestArrayFlags = map[string]string{
	"tags":      "tags",
	"resolvers": "resolvers",
}

func addWorkRequestWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("user-id", "", "User ID")
	cmd.Flags().String("date", "", "Date")
	cmd.Flags().String("resolution-date", "", "Resolution date")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("component-id", "", "Component ID")
	cmd.Flags().String("part-id", "", "Part ID")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("status", "", "Status")
	cmd.Flags().String("plan-of-action", "", "Plan of action")
	cmd.Flags().String("priority", "", "Priority")
	cmd.Flags().String("notification-group-id", "", "Notification group ID")
	cmd.Flags().String("location", "", "Location")
	cmd.Flags().String("manager-notes", "", "Manager notes")
	cmd.Flags().String("code-id", "", "Code ID")
	cmd.Flags().String("reference-number", "", "Reference number")
	cmd.Flags().Bool("self-resolve", false, "Self resolve")
	cmd.Flags().Bool("modify-resolvers", false, "Modify resolvers")
	cmd.Flags().String("tags", "", "Tags (comma-separated)")
	cmd.Flags().String("resolvers", "", "Resolver IDs (comma-separated)")
}

func init() {
	workRequestsListCmd.Flags().Int("page", 1, "Page number")
	workRequestsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	workRequestsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")

	addWorkRequestWriteFlags(workRequestsCreateCmd)
	addWorkRequestWriteFlags(workRequestsUpdateCmd)

	workRequestsCmd.AddCommand(workRequestsListCmd)
	workRequestsCmd.AddCommand(workRequestsGetCmd)
	workRequestsCmd.AddCommand(workRequestsCreateCmd)
	workRequestsCmd.AddCommand(workRequestsUpdateCmd)
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

func runWorkRequestsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, workRequestFlags)
	bodyFromBoolFlags(cmd, body, workRequestBoolFlags)
	bodyFromArrayFlags(cmd, body, workRequestArrayFlags)

	response, err := c.Post("work-requests", wrapBody("work_request", body))
	if err != nil {
		return handleClientError(err)
	}

	wrID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Work request created", []string{
		fmt.Sprintf("mobileops work-requests get %s", wrID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runWorkRequestsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, workRequestFlags)
	bodyFromBoolFlags(cmd, body, workRequestBoolFlags)
	bodyFromArrayFlags(cmd, body, workRequestArrayFlags)

	response, err := c.Put(fmt.Sprintf("work-requests/%s", id), wrapBody("work_request", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Work request updated", []string{
		fmt.Sprintf("mobileops work-requests get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
