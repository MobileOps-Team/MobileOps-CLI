package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var nonconformitiesCmd = &cobra.Command{
	Use:   "nonconformities",
	Short: "Manage nonconformities",
}

var nonconformitiesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List nonconformities",
	RunE:  runNonconformitiesList,
}

var nonconformitiesGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single nonconformity by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runNonconformitiesGet,
}

var nonconformitiesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new nonconformity",
	RunE:  runNonconformitiesCreate,
}

var nonconformitiesUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a nonconformity",
	Args:  cobra.ExactArgs(1),
	RunE:  runNonconformitiesUpdate,
}

var nonconformityFlags = map[string]string{
	"user-id":               "user_id",
	"date":                  "date",
	"vessel-id":             "vessel_id",
	"plan-of-action":        "plan_of_action",
	"description":           "description",
	"status":                "status",
	"priority":              "priority",
	"notification-group-id": "notification_group_id",
	"root-cause":            "root_cause",
	"corrective-actions":    "corrective_actions",
	"location":              "location",
	"manager-notes":         "manager_notes",
	"resolution-date":       "resolution_date",
}

var nonconformityBoolFlags = map[string]string{
	"major":            "major",
	"external":         "external",
	"shoreside":        "shoreside",
	"self-resolve":     "self_resolve",
	"modify-resolvers": "modify_resolvers",
}

var nonconformityArrayFlags = map[string]string{
	"tags":      "tags",
	"resolvers": "resolvers",
}

func addNonconformityWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("user-id", "", "User ID")
	cmd.Flags().String("date", "", "Date")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("plan-of-action", "", "Plan of action")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("status", "", "Status")
	cmd.Flags().String("priority", "", "Priority")
	cmd.Flags().String("notification-group-id", "", "Notification group ID")
	cmd.Flags().String("root-cause", "", "Root cause")
	cmd.Flags().String("corrective-actions", "", "Corrective actions")
	cmd.Flags().String("location", "", "Location")
	cmd.Flags().String("manager-notes", "", "Manager notes")
	cmd.Flags().String("resolution-date", "", "Resolution date")
	cmd.Flags().Bool("major", false, "Major nonconformity")
	cmd.Flags().Bool("external", false, "External")
	cmd.Flags().Bool("shoreside", false, "Shoreside")
	cmd.Flags().Bool("self-resolve", false, "Self resolve")
	cmd.Flags().Bool("modify-resolvers", false, "Modify resolvers")
	cmd.Flags().String("tags", "", "Tags (comma-separated)")
	cmd.Flags().String("resolvers", "", "Resolver IDs (comma-separated)")
}

func init() {
	nonconformitiesListCmd.Flags().Int("page", 1, "Page number")
	nonconformitiesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	nonconformitiesListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")

	addNonconformityWriteFlags(nonconformitiesCreateCmd)
	addNonconformityWriteFlags(nonconformitiesUpdateCmd)

	nonconformitiesCmd.AddCommand(nonconformitiesListCmd)
	nonconformitiesCmd.AddCommand(nonconformitiesGetCmd)
	nonconformitiesCmd.AddCommand(nonconformitiesCreateCmd)
	nonconformitiesCmd.AddCommand(nonconformitiesUpdateCmd)
}

func runNonconformitiesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("vessel-id"); v != "" {
		params["vessel_id"] = v
	}

	response, err := c.Get("nonconformities", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops nonconformities get %s", 3)
	env := envelope.WrapCollection(response, "nonconformities", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runNonconformitiesGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("nonconformities/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if vesselID, ok := response["vessel_id"].(string); ok && vesselID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %s", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops nonconformities list")

	env := envelope.WrapRecord(response, "Nonconformity", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runNonconformitiesCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, nonconformityFlags)
	bodyFromBoolFlags(cmd, body, nonconformityBoolFlags)
	bodyFromArrayFlags(cmd, body, nonconformityArrayFlags)

	response, err := c.Post("nonconformities", wrapBody("nonconformity", body))
	if err != nil {
		return handleClientError(err)
	}

	ncID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Nonconformity created", []string{
		fmt.Sprintf("mobileops nonconformities get %s", ncID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runNonconformitiesUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, nonconformityFlags)
	bodyFromBoolFlags(cmd, body, nonconformityBoolFlags)
	bodyFromArrayFlags(cmd, body, nonconformityArrayFlags)

	response, err := c.Put(fmt.Sprintf("nonconformities/%s", id), wrapBody("nonconformity", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Nonconformity updated", []string{
		fmt.Sprintf("mobileops nonconformities get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
