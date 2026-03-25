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

func init() {
	componentsListCmd.Flags().Int("page", 1, "Page number")
	componentsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	componentsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")

	componentsCmd.AddCommand(componentsListCmd)
	componentsCmd.AddCommand(componentsGetCmd)
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
