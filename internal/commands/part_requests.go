package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var partRequestsCmd = &cobra.Command{
	Use:   "part-requests",
	Short: "Manage part requests",
}

var partRequestsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List part requests",
	RunE:  runPartRequestsList,
}

var partRequestsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single part request by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runPartRequestsGet,
}

var partRequestsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new part request",
	RunE:  runPartRequestsCreate,
}

var partRequestsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a part request",
	Args:  cobra.ExactArgs(1),
	RunE:  runPartRequestsUpdate,
}

var partRequestFlags = map[string]string{
	"vessel-id":    "vessel_id",
	"component-id": "component_id",
	"part-id":      "part_id",
	"model-id":     "model_id",
	"quantity":     "quantity",
	"status":       "status",
	"notes":        "notes",
	"user-id":      "user_id",
}

func addPartRequestWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("component-id", "", "Component ID")
	cmd.Flags().String("part-id", "", "Part ID")
	cmd.Flags().String("model-id", "", "Model ID")
	cmd.Flags().String("quantity", "", "Quantity")
	cmd.Flags().String("status", "", "Status")
	cmd.Flags().String("notes", "", "Notes")
	cmd.Flags().String("user-id", "", "User ID")
}

func init() {
	partRequestsListCmd.Flags().Int("page", 1, "Page number")
	partRequestsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	partRequestsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	partRequestsListCmd.Flags().String("component-id", "", "Filter by component ID")
	partRequestsListCmd.Flags().String("part-id", "", "Filter by part ID")

	addPartRequestWriteFlags(partRequestsCreateCmd)
	addPartRequestWriteFlags(partRequestsUpdateCmd)

	partRequestsCmd.AddCommand(partRequestsListCmd)
	partRequestsCmd.AddCommand(partRequestsGetCmd)
	partRequestsCmd.AddCommand(partRequestsCreateCmd)
	partRequestsCmd.AddCommand(partRequestsUpdateCmd)
}

func runPartRequestsList(cmd *cobra.Command, args []string) error {
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
	if v, _ := cmd.Flags().GetString("part-id"); v != "" {
		params["part_id"] = v
	}

	response, err := c.Get("part-requests", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops part-requests get %s", 3)
	env := envelope.WrapCollection(response, "part requests", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runPartRequestsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("part-requests/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if vesselID, ok := response["vessel_id"].(string); ok && vesselID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %s", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops part-requests list")

	env := envelope.WrapRecord(response, "Part request", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runPartRequestsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, partRequestFlags)

	response, err := c.Post("part-requests", wrapBody("part_request", body))
	if err != nil {
		return handleClientError(err)
	}

	prID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Part request created", []string{
		fmt.Sprintf("mobileops part-requests get %s", prID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runPartRequestsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, partRequestFlags)

	response, err := c.Put(fmt.Sprintf("part-requests/%s", id), wrapBody("part_request", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Part request updated", []string{
		fmt.Sprintf("mobileops part-requests get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
