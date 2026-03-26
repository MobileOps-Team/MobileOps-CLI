package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var cargoTypesCmd = &cobra.Command{
	Use:   "cargo-types",
	Short: "Manage cargo types",
}

var cargoTypesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cargo types",
	RunE:  runCargoTypesList,
}

var cargoTypesGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single cargo type by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runCargoTypesGet,
}

var cargoTypesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cargo type",
	RunE:  runCargoTypesCreate,
}

var cargoTypesUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a cargo type",
	Args:  cobra.ExactArgs(1),
	RunE:  runCargoTypesUpdate,
}

var cargoTypeFlags = map[string]string{
	"name":             "name",
	"reference-number": "reference_number",
	"color":            "color",
}

var cargoTypeArrayFlags = map[string]string{
	"measurement-ids": "measurement_ids",
}

func addCargoTypeWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Cargo type name")
	cmd.Flags().String("reference-number", "", "Reference number")
	cmd.Flags().String("color", "", "Color")
	cmd.Flags().String("measurement-ids", "", "Measurement IDs (comma-separated)")
}

func init() {
	cargoTypesListCmd.Flags().Int("page", 1, "Page number")
	cargoTypesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addCargoTypeWriteFlags(cargoTypesCreateCmd)
	addCargoTypeWriteFlags(cargoTypesUpdateCmd)

	cargoTypesCmd.AddCommand(cargoTypesListCmd)
	cargoTypesCmd.AddCommand(cargoTypesGetCmd)
	cargoTypesCmd.AddCommand(cargoTypesCreateCmd)
	cargoTypesCmd.AddCommand(cargoTypesUpdateCmd)
}

func runCargoTypesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("cargo-types", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops cargo-types get %s", 3)
	env := envelope.WrapCollection(response, "cargo types", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runCargoTypesGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("cargo-types/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Cargo type", []string{
		"mobileops cargo-types list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCargoTypesCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, cargoTypeFlags)
	bodyFromArrayFlags(cmd, body, cargoTypeArrayFlags)

	response, err := c.Post("cargo-types", wrapBody("cargo_type", body))
	if err != nil {
		return handleClientError(err)
	}

	ctID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Cargo type created", []string{
		fmt.Sprintf("mobileops cargo-types get %s", ctID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCargoTypesUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, cargoTypeFlags)
	bodyFromArrayFlags(cmd, body, cargoTypeArrayFlags)

	response, err := c.Put(fmt.Sprintf("cargo-types/%s", id), wrapBody("cargo_type", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Cargo type updated", []string{
		fmt.Sprintf("mobileops cargo-types get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
