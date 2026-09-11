package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Manage vessel components",
}

var componentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List components",
	RunE:  runComponentsList,
}

var componentsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single component by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runComponentsGet,
}

var componentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new component",
	RunE:  runComponentsCreate,
}

var componentsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a component",
	Args:  cobra.ExactArgs(1),
	RunE:  runComponentsUpdate,
}

var componentFlags = map[string]string{
	"name":                    "name",
	"vessel-id":               "vessel_id",
	"model-id":                "model_id",
	"model-number":            "model_number",
	"ultra-field-template-id": "ultra_field_template_id",
	"hours":                   "hours",
	"lifetime-hours":          "lifetime_hours",
	"specifications":          "specifications",
}

var componentBoolFlags = map[string]string{
	"has-hours":   "has_hours",
	"critical":    "critical",
	"reset-hours": "reset_hours",
}

func addComponentWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Component name")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("model-id", "", "Model ID")
	cmd.Flags().String("model-number", "", "Model number")
	cmd.Flags().String("ultra-field-template-id", "", "Ultra field template ID")
	cmd.Flags().String("hours", "", "Hours")
	cmd.Flags().String("lifetime-hours", "", "Lifetime hours")
	cmd.Flags().String("specifications", "", "Free-text specifications")
	cmd.Flags().Bool("has-hours", false, "Has hours")
	cmd.Flags().Bool("critical", false, "Critical")
	cmd.Flags().Bool("reset-hours", false, "Reset hours")
}

func init() {
	componentsListCmd.Flags().Int("page", 1, "Page number")
	componentsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	componentsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")

	addComponentWriteFlags(componentsCreateCmd)
	addComponentWriteFlags(componentsUpdateCmd)

	componentsCmd.AddCommand(componentsListCmd)
	componentsCmd.AddCommand(componentsGetCmd)
	componentsCmd.AddCommand(componentsCreateCmd)
	componentsCmd.AddCommand(componentsUpdateCmd)
}

func runComponentsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if vesselID, _ := cmd.Flags().GetString("vessel-id"); vesselID != "" {
		params["vessel_id"] = vesselID
	}

	response, err := c.Get("components", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops components get %s", 3)
	env := envelope.WrapCollection(response, "components", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runComponentsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("components/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if vesselID, ok := response["vessel_id"].(string); ok && vesselID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %s", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops components list")

	env := envelope.WrapRecord(response, "Component", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runComponentsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, componentFlags)
	bodyFromBoolFlags(cmd, body, componentBoolFlags)

	response, err := c.Post("components", wrapBody("component", body))
	if err != nil {
		return handleClientError(err)
	}

	componentID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Component created", []string{
		fmt.Sprintf("mobileops components get %s", componentID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runComponentsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, componentFlags)
	bodyFromBoolFlags(cmd, body, componentBoolFlags)

	response, err := c.Put(fmt.Sprintf("components/%s", id), wrapBody("component", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Component updated", []string{
		fmt.Sprintf("mobileops components get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
