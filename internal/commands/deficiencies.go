package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var deficienciesCmd = &cobra.Command{
	Use:   "deficiencies",
	Short: "Manage deficiencies",
}

var deficienciesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List deficiencies",
	RunE:  runDeficienciesList,
}

var deficienciesGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single deficiency by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeficienciesGet,
}

var deficienciesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new deficiency",
	RunE:  runDeficienciesCreate,
}

var deficienciesUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a deficiency",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeficienciesUpdate,
}

var deficiencyFlags = map[string]string{
	"user-id":               "user_id",
	"date":                  "date",
	"resolution-date":       "resolution_date",
	"vessel-id":             "vessel_id",
	"component-id":          "component_id",
	"part-id":               "part_id",
	"description":           "description",
	"status":                "status",
	"priority":              "priority",
	"notification-group-id": "notification_group_id",
	"root-cause":            "root_cause",
	"corrective-actions":    "corrective_actions",
	"location":              "location",
	"manager-notes":         "manager_notes",
	"code-id":               "code_id",
	"reference-number":      "reference_number",
}

var deficiencyBoolFlags = map[string]string{
	"self-resolve":     "self_resolve",
	"modify-resolvers": "modify_resolvers",
}

var deficiencyArrayFlags = map[string]string{
	"tags":      "tags",
	"resolvers": "resolvers",
}

func addDeficiencyWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("user-id", "", "User ID")
	cmd.Flags().String("date", "", "Date")
	cmd.Flags().String("resolution-date", "", "Resolution date")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("component-id", "", "Component ID")
	cmd.Flags().String("part-id", "", "Part ID")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("status", "", "Status")
	cmd.Flags().String("priority", "", "Priority")
	cmd.Flags().String("notification-group-id", "", "Notification group ID")
	cmd.Flags().String("root-cause", "", "Root cause")
	cmd.Flags().String("corrective-actions", "", "Corrective actions")
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
	deficienciesListCmd.Flags().Int("page", 1, "Page number")
	deficienciesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	deficienciesListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	deficienciesListCmd.Flags().String("component-id", "", "Filter by component ID")
	deficienciesListCmd.Flags().String("part-id", "", "Filter by part ID")

	addDeficiencyWriteFlags(deficienciesCreateCmd)
	addDeficiencyWriteFlags(deficienciesUpdateCmd)

	deficienciesCmd.AddCommand(deficienciesListCmd)
	deficienciesCmd.AddCommand(deficienciesGetCmd)
	deficienciesCmd.AddCommand(deficienciesCreateCmd)
	deficienciesCmd.AddCommand(deficienciesUpdateCmd)
}

func runDeficienciesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("vessel-id"); v != "" {
		params["vessel_id"] = v
	}
	if v, _ := cmd.Flags().GetString("component-id"); v != "" {
		params["component_id"] = v
	}
	if v, _ := cmd.Flags().GetString("part-id"); v != "" {
		params["part_id"] = v
	}

	response, err := c.Get("deficiencies", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops deficiencies get %s", 3)
	env := envelope.WrapCollection(response, "deficiencies", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runDeficienciesGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("deficiencies/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if vesselID, ok := response["vessel_id"].(string); ok && vesselID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %s", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops deficiencies list")

	env := envelope.WrapRecord(response, "Deficiency", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runDeficienciesCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, deficiencyFlags)
	bodyFromBoolFlags(cmd, body, deficiencyBoolFlags)
	bodyFromArrayFlags(cmd, body, deficiencyArrayFlags)

	response, err := c.Post("deficiencies", wrapBody("deficiency", body))
	if err != nil {
		return handleClientError(err)
	}

	defID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Deficiency created", []string{
		fmt.Sprintf("mobileops deficiencies get %s", defID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runDeficienciesUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, deficiencyFlags)
	bodyFromBoolFlags(cmd, body, deficiencyBoolFlags)
	bodyFromArrayFlags(cmd, body, deficiencyArrayFlags)

	response, err := c.Put(fmt.Sprintf("deficiencies/%s", id), wrapBody("deficiency", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Deficiency updated", []string{
		fmt.Sprintf("mobileops deficiencies get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
