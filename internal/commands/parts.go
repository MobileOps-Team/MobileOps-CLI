package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var partsCmd = &cobra.Command{
	Use:   "parts",
	Short: "Manage vessel parts",
}

var partsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List parts",
	RunE:  runPartsList,
}

var partsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single part by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runPartsGet,
}

var partsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new part",
	RunE:  runPartsCreate,
}

var partsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a part",
	Args:  cobra.ExactArgs(1),
	RunE:  runPartsUpdate,
}

var partFlags = map[string]string{
	"name":                    "name",
	"vessel-id":               "vessel_id",
	"component-id":            "component_id",
	"model-id":                "model_id",
	"model-number":            "model_number",
	"code-id":                 "code_id",
	"ultra-field-template-id": "ultra_field_template_id",
	"hours":                   "hours",
	"specifications":          "specifications",
}

var partBoolFlags = map[string]string{
	"has-hours": "has_hours",
	"critical":  "critical",
}

func addPartWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Part name")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("component-id", "", "Component ID")
	cmd.Flags().String("model-id", "", "Model ID")
	cmd.Flags().String("model-number", "", "Model number")
	cmd.Flags().String("code-id", "", "Code ID")
	cmd.Flags().String("ultra-field-template-id", "", "Ultra field template ID")
	cmd.Flags().String("hours", "", "Hours")
	cmd.Flags().String("specifications", "", "Free-text specifications")
	cmd.Flags().Bool("has-hours", false, "Has hours")
	cmd.Flags().Bool("critical", false, "Critical")
}

func init() {
	partsListCmd.Flags().Int("page", 1, "Page number")
	partsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	partsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	partsListCmd.Flags().String("component-id", "", "Filter by component ID")

	addPartWriteFlags(partsCreateCmd)
	addPartWriteFlags(partsUpdateCmd)

	partsCmd.AddCommand(partsListCmd)
	partsCmd.AddCommand(partsGetCmd)
	partsCmd.AddCommand(partsCreateCmd)
	partsCmd.AddCommand(partsUpdateCmd)
}

func runPartsList(cmd *cobra.Command, args []string) error {
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

	response, err := c.Get("parts", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops parts get %s", 3)
	env := envelope.WrapCollection(response, "parts", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runPartsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("parts/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if componentID, ok := response["component_id"].(string); ok && componentID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops components get %s", componentID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops parts list")

	env := envelope.WrapRecord(response, "Part", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runPartsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, partFlags)
	bodyFromBoolFlags(cmd, body, partBoolFlags)

	response, err := c.Post("parts", wrapBody("part", body))
	if err != nil {
		return handleClientError(err)
	}

	partID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Part created", []string{
		fmt.Sprintf("mobileops parts get %s", partID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runPartsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, partFlags)
	bodyFromBoolFlags(cmd, body, partBoolFlags)

	response, err := c.Put(fmt.Sprintf("parts/%s", id), wrapBody("part", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Part updated", []string{
		fmt.Sprintf("mobileops parts get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
