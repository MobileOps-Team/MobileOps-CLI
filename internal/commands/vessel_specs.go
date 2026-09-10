package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var vesselSpecsCmd = &cobra.Command{
	Use:   "vessel-specs",
	Short: "Manage vessel specs",
}

var vesselSpecsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List vessel specs",
	RunE:  runVesselSpecsList,
}

var vesselSpecsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single vessel spec by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselSpecsGet,
}

var vesselSpecsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vessel spec",
	RunE:  runVesselSpecsCreate,
}

var vesselSpecsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a vessel spec",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselSpecsUpdate,
}

var vesselSpecFlags = map[string]string{
	"vessel-id":   "vessel_id",
	"template-id": "template_id",
}

var vesselSpecJSONFlags = map[string]string{
	"spec-fields": "spec_fields",
}

func addVesselSpecWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("vessel-id", "", "Vessel ID (one spec per vessel)")
	cmd.Flags().String("template-id", "", "Vessel spec template ID (required; filters spec-fields)")
	cmd.Flags().String("spec-fields", "", "Spec values as a JSON object keyed by template field, e.g. '{\"length\":\"120 ft\"}'")
}

func init() {
	vesselSpecsListCmd.Flags().Int("page", 1, "Page number")
	vesselSpecsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	vesselSpecsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")

	addVesselSpecWriteFlags(vesselSpecsCreateCmd)
	addVesselSpecWriteFlags(vesselSpecsUpdateCmd)

	vesselSpecsCmd.AddCommand(vesselSpecsListCmd)
	vesselSpecsCmd.AddCommand(vesselSpecsGetCmd)
	vesselSpecsCmd.AddCommand(vesselSpecsCreateCmd)
	vesselSpecsCmd.AddCommand(vesselSpecsUpdateCmd)
}

func runVesselSpecsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("vessel-id"); v != "" {
		params["vessel_id"] = v
	}

	response, err := c.Get("vessel-specs", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops vessel-specs get %s", 3)
	env := envelope.WrapCollection(response, "vessel specs", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselSpecsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("vessel-specs/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if vesselID, ok := response["vessel_id"].(string); ok && vesselID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %s", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops vessel-specs list")

	env := envelope.WrapRecord(response, "Vessel spec", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselSpecsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vesselSpecFlags)
	if err := bodyFromJSONFlags(cmd, body, vesselSpecJSONFlags); err != nil {
		return err
	}

	response, err := c.Post("vessel-specs", wrapBody("vessel_spec", body))
	if err != nil {
		return handleClientError(err)
	}

	specID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Vessel spec created", []string{
		fmt.Sprintf("mobileops vessel-specs get %s", specID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselSpecsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vesselSpecFlags)
	if err := bodyFromJSONFlags(cmd, body, vesselSpecJSONFlags); err != nil {
		return err
	}

	response, err := c.Put(fmt.Sprintf("vessel-specs/%s", id), wrapBody("vessel_spec", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vessel spec updated", []string{
		fmt.Sprintf("mobileops vessel-specs get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
