package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var vesselSpecTemplatesCmd = &cobra.Command{
	Use:   "vessel-spec-templates",
	Short: "List vessel spec templates",
}

var vesselSpecTemplatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List vessel spec templates",
	RunE:  runVesselSpecTemplatesList,
}

var vesselSpecTemplatesGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single vessel spec template by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselSpecTemplatesGet,
}

func init() {
	vesselSpecTemplatesListCmd.Flags().Int("page", 1, "Page number")
	vesselSpecTemplatesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	vesselSpecTemplatesCmd.AddCommand(vesselSpecTemplatesListCmd)
	vesselSpecTemplatesCmd.AddCommand(vesselSpecTemplatesGetCmd)
}

func runVesselSpecTemplatesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("vessel-spec-templates", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops vessel-spec-templates get %s", 3)
	env := envelope.WrapCollection(response, "vessel spec templates", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselSpecTemplatesGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("vessel-spec-templates/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vessel spec template", []string{
		"mobileops vessel-spec-templates list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
