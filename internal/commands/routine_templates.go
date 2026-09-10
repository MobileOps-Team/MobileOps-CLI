package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var routineTemplatesCmd = &cobra.Command{
	Use:   "routine-templates",
	Short: "List routine templates",
}

var routineTemplatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List routine templates",
	RunE:  runRoutineTemplatesList,
}

func init() {
	routineTemplatesListCmd.Flags().Int("page", 1, "Page number")
	routineTemplatesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	routineTemplatesListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	routineTemplatesListCmd.Flags().String("component-id", "", "Filter by component ID")
	routineTemplatesListCmd.Flags().String("part-id", "", "Filter by part ID")
	routineTemplatesListCmd.Flags().String("category", "", "Filter by category")
	routineTemplatesListCmd.Flags().String("type", "", "Filter by type")
	routineTemplatesListCmd.Flags().String("frequency-type", "", "Filter by frequency type")
	routineTemplatesListCmd.Flags().String("id", "", "Return a single template by ID (other filters ignored)")
	routineTemplatesListCmd.Flags().String("master", "", "Filter by master flag (true/false)")
	routineTemplatesListCmd.Flags().String("master-id", "", "Templates inheriting from this master template")
	routineTemplatesListCmd.Flags().String("universal", "", "Filter by universal flag (true/false)")

	routineTemplatesCmd.AddCommand(routineTemplatesListCmd)
}

func runRoutineTemplatesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters(cmd, params, map[string]string{
		"id":        "id",
		"master":    "master",
		"master-id": "master_id",
		"universal": "universal",
	})
	if v, _ := cmd.Flags().GetString("vessel-id"); v != "" {
		params["vessel_id"] = v
	}
	if v, _ := cmd.Flags().GetString("component-id"); v != "" {
		params["component_id"] = v
	}
	if v, _ := cmd.Flags().GetString("part-id"); v != "" {
		params["part_id"] = v
	}
	if v, _ := cmd.Flags().GetString("category"); v != "" {
		params["category"] = v
	}
	if v, _ := cmd.Flags().GetString("type"); v != "" {
		params["type"] = v
	}
	if v, _ := cmd.Flags().GetString("frequency-type"); v != "" {
		params["frequency_type"] = v
	}

	response, err := c.Get("routine-templates", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "routine templates", []string{
		"mobileops routine-calculations list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
