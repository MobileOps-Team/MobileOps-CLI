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
	Short: "List cargo types (read-only)",
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

func init() {
	cargoTypesListCmd.Flags().Int("page", 1, "Page number")
	cargoTypesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	cargoTypesCmd.AddCommand(cargoTypesListCmd)
	cargoTypesCmd.AddCommand(cargoTypesGetCmd)
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
