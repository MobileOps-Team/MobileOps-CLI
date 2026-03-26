package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var routineCalculationsCmd = &cobra.Command{
	Use:   "routine-calculations",
	Short: "List routine calculations",
}

var routineCalculationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List routine calculations",
	RunE:  runRoutineCalculationsList,
}

func init() {
	routineCalculationsListCmd.Flags().Int("page", 1, "Page number")
	routineCalculationsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	routineCalculationsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	routineCalculationsListCmd.Flags().String("component-id", "", "Filter by component ID")
	routineCalculationsListCmd.Flags().String("part-id", "", "Filter by part ID")
	routineCalculationsListCmd.Flags().String("division-id", "", "Filter by division ID")
	routineCalculationsListCmd.Flags().String("frequency-type", "", "Filter by frequency type")
	routineCalculationsListCmd.Flags().String("value-lte", "", "Value less than or equal")
	routineCalculationsListCmd.Flags().String("value-gte", "", "Value greater than or equal")
	routineCalculationsListCmd.Flags().String("due-date-lte", "", "Due date before (YYYY-MM-DD)")
	routineCalculationsListCmd.Flags().String("due-date-gte", "", "Due date after (YYYY-MM-DD)")
	routineCalculationsListCmd.Flags().String("routine-template-id", "", "Filter by routine template ID")
	routineCalculationsListCmd.Flags().String("master-template-id", "", "Filter by master template ID")

	routineCalculationsCmd.AddCommand(routineCalculationsListCmd)
}

func runRoutineCalculationsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	flagMap := map[string]string{
		"vessel-id":           "vessel_id",
		"component-id":        "component_id",
		"part-id":             "part_id",
		"division-id":         "division_id",
		"frequency-type":      "frequency_type",
		"value-lte":           "value_lte",
		"value-gte":           "value_gte",
		"due-date-lte":        "due_date_lte",
		"due-date-gte":        "due_date_gte",
		"routine-template-id": "routine_template_id",
		"master-template-id":  "master_template_id",
	}
	for flag, param := range flagMap {
		if v, _ := cmd.Flags().GetString(flag); v != "" {
			params[param] = v
		}
	}

	response, err := c.Get("routine-calculations", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "routine calculations", []string{
		"mobileops routine-templates list",
		"mobileops vessels list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
