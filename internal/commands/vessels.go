package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var vesselsCmd = &cobra.Command{
	Use:   "vessels",
	Short: "Manage vessels/assets",
}

var vesselsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all vessels/assets",
	RunE:  runVesselsList,
}

var vesselsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single vessel by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselsGet,
}

var vesselsSpecsCmd = &cobra.Command{
	Use:   "specs VESSEL_ID",
	Short: "List specs for a vessel",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselsSpecs,
}

func init() {
	vesselsListCmd.Flags().Int("page", 1, "Page number")
	vesselsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	vesselsCmd.AddCommand(vesselsListCmd)
	vesselsCmd.AddCommand(vesselsGetCmd)
	vesselsCmd.AddCommand(vesselsSpecsCmd)
}

func runVesselsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("assets", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops vessels get %s", 3)
	env := envelope.WrapCollection(response, "vessels", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("assets/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	vesselID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Vessel", []string{
		fmt.Sprintf("mobileops vessels specs %s", vesselID),
		fmt.Sprintf("mobileops components list --vessel-id %s", vesselID),
		fmt.Sprintf("mobileops jobs list --vessel-id %s", vesselID),
		fmt.Sprintf("mobileops work-requests list --vessel-id %s", vesselID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselsSpecs(cmd *cobra.Command, args []string) error {
	vesselID := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get("vessel-specs", map[string]string{"vessel_id": vesselID})
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "vessel specs", []string{
		fmt.Sprintf("mobileops vessels get %s", vesselID),
		fmt.Sprintf("mobileops components list --vessel-id %s", vesselID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
