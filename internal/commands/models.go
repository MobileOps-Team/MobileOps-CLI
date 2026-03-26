package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Manage models",
}

var modelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List models",
	RunE:  runModelsList,
}

var modelsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single model by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runModelsGet,
}

var modelsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new model",
	RunE:  runModelsCreate,
}

var modelsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a model",
	Args:  cobra.ExactArgs(1),
	RunE:  runModelsUpdate,
}

var modelFlags = map[string]string{
	"make-id":                  "make_id",
	"serial-number":            "serial_number",
	"name":                     "name",
	"part-number":              "part_number",
	"notes":                    "notes",
	"vessel-id":                "vessel_id",
	"unit-cost":                "unit_cost",
	"expected-lifetime-hours":  "expected_lifetime_hours",
}

var modelBoolFlags = map[string]string{
	"critical": "critical",
}

func addModelWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("make-id", "", "Make ID")
	cmd.Flags().String("serial-number", "", "Serial number")
	cmd.Flags().String("name", "", "Model name")
	cmd.Flags().String("part-number", "", "Part number")
	cmd.Flags().String("notes", "", "Notes")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("unit-cost", "", "Unit cost")
	cmd.Flags().String("expected-lifetime-hours", "", "Expected lifetime hours")
	cmd.Flags().Bool("critical", false, "Critical")
}

func init() {
	modelsListCmd.Flags().Int("page", 1, "Page number")
	modelsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	modelsListCmd.Flags().String("make-id", "", "Filter by make ID")

	modelsGetCmd.Flags().String("make-id", "", "Parent make ID (required)")

	addModelWriteFlags(modelsCreateCmd)
	addModelWriteFlags(modelsUpdateCmd)

	modelsCmd.AddCommand(modelsListCmd)
	modelsCmd.AddCommand(modelsGetCmd)
	modelsCmd.AddCommand(modelsCreateCmd)
	modelsCmd.AddCommand(modelsUpdateCmd)
}

func runModelsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("make-id"); v != "" {
		params["make_id"] = v
	}

	response, err := c.Get("models", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops models get %s", 3)
	env := envelope.WrapCollection(response, "models", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runModelsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := map[string]string{}
	if v, _ := cmd.Flags().GetString("make-id"); v != "" {
		params["make_id"] = v
	}

	response, err := c.Get(fmt.Sprintf("models/%s", id), params)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if makeID, ok := response["make_id"].(string); ok && makeID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops makes get %s", makeID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops models list")

	env := envelope.WrapRecord(response, "Model", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runModelsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, modelFlags)
	bodyFromBoolFlags(cmd, body, modelBoolFlags)

	response, err := c.Post("models", wrapBody("model", body))
	if err != nil {
		return handleClientError(err)
	}

	modelID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Model created", []string{
		fmt.Sprintf("mobileops models get %s", modelID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runModelsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, modelFlags)
	bodyFromBoolFlags(cmd, body, modelBoolFlags)

	response, err := c.Put(fmt.Sprintf("models/%s", id), wrapBody("model", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Model updated", []string{
		fmt.Sprintf("mobileops models get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
