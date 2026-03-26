package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var vesselTypesCmd = &cobra.Command{
	Use:   "vessel-types",
	Short: "Manage vessel types",
}

var vesselTypesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List vessel types",
	RunE:  runVesselTypesList,
}

var vesselTypesGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single vessel type by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselTypesGet,
}

var vesselTypesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vessel type",
	RunE:  runVesselTypesCreate,
}

var vesselTypesUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a vessel type",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselTypesUpdate,
}

var vesselTypeFlags = map[string]string{
	"name": "name",
}

func addVesselTypeWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Vessel type name")
}

func init() {
	vesselTypesListCmd.Flags().Int("page", 1, "Page number")
	vesselTypesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addVesselTypeWriteFlags(vesselTypesCreateCmd)
	addVesselTypeWriteFlags(vesselTypesUpdateCmd)

	vesselTypesCmd.AddCommand(vesselTypesListCmd)
	vesselTypesCmd.AddCommand(vesselTypesGetCmd)
	vesselTypesCmd.AddCommand(vesselTypesCreateCmd)
	vesselTypesCmd.AddCommand(vesselTypesUpdateCmd)
}

func runVesselTypesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("vessel-types", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops vessel-types get %s", 3)
	env := envelope.WrapCollection(response, "vessel types", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselTypesGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("vessel-types/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vessel type", []string{
		"mobileops vessel-types list",
		"mobileops vessels list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselTypesCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vesselTypeFlags)

	response, err := c.Post("vessel-types", wrapBody("vessel_type", body))
	if err != nil {
		return handleClientError(err)
	}

	vtID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Vessel type created", []string{
		fmt.Sprintf("mobileops vessel-types get %s", vtID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselTypesUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vesselTypeFlags)

	response, err := c.Put(fmt.Sprintf("vessel-types/%s", id), wrapBody("vessel_type", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vessel type updated", []string{
		fmt.Sprintf("mobileops vessel-types get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
