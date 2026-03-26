package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var maintenanceReportsCmd = &cobra.Command{
	Use:   "maintenance-reports",
	Short: "Manage maintenance reports",
}

var maintenanceReportsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List maintenance reports",
	RunE:  runMaintenanceReportsList,
}

var maintenanceReportsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single maintenance report by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runMaintenanceReportsGet,
}

var maintenanceReportsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new maintenance report",
	RunE:  runMaintenanceReportsCreate,
}

var maintenanceReportsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a maintenance report",
	Args:  cobra.ExactArgs(1),
	RunE:  runMaintenanceReportsUpdate,
}

var maintenanceReportFlags = map[string]string{
	"description":     "description",
	"date":            "date",
	"vessel-id":       "vessel_id",
	"component-id":    "component_id",
	"component-hours": "component_hours",
	"part-id":         "part_id",
	"user-id":         "user_id",
}

var maintenanceReportArrayFlags = map[string]string{
	"users": "users",
}

func addMaintenanceReportWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("date", "", "Date")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("component-id", "", "Component ID")
	cmd.Flags().String("component-hours", "", "Component hours")
	cmd.Flags().String("part-id", "", "Part ID")
	cmd.Flags().String("user-id", "", "User ID")
	cmd.Flags().String("users", "", "User IDs (comma-separated)")
}

func init() {
	maintenanceReportsListCmd.Flags().Int("page", 1, "Page number")
	maintenanceReportsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	maintenanceReportsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	maintenanceReportsListCmd.Flags().String("component-id", "", "Filter by component ID")
	maintenanceReportsListCmd.Flags().String("part-id", "", "Filter by part ID")

	addMaintenanceReportWriteFlags(maintenanceReportsCreateCmd)
	addMaintenanceReportWriteFlags(maintenanceReportsUpdateCmd)

	maintenanceReportsCmd.AddCommand(maintenanceReportsListCmd)
	maintenanceReportsCmd.AddCommand(maintenanceReportsGetCmd)
	maintenanceReportsCmd.AddCommand(maintenanceReportsCreateCmd)
	maintenanceReportsCmd.AddCommand(maintenanceReportsUpdateCmd)
}

func runMaintenanceReportsList(cmd *cobra.Command, args []string) error {
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

	response, err := c.Get("maintenance-reports", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops maintenance-reports get %s", 3)
	env := envelope.WrapCollection(response, "maintenance reports", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runMaintenanceReportsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("maintenance-reports/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if vesselID, ok := response["vessel_id"].(string); ok && vesselID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops vessels get %s", vesselID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops maintenance-reports list")

	env := envelope.WrapRecord(response, "Maintenance report", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runMaintenanceReportsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, maintenanceReportFlags)
	bodyFromArrayFlags(cmd, body, maintenanceReportArrayFlags)

	response, err := c.Post("maintenance-reports", wrapBody("maintenance_report", body))
	if err != nil {
		return handleClientError(err)
	}

	mrID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Maintenance report created", []string{
		fmt.Sprintf("mobileops maintenance-reports get %s", mrID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runMaintenanceReportsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, maintenanceReportFlags)
	bodyFromArrayFlags(cmd, body, maintenanceReportArrayFlags)

	response, err := c.Put(fmt.Sprintf("maintenance-reports/%s", id), wrapBody("maintenance_report", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Maintenance report updated", []string{
		fmt.Sprintf("mobileops maintenance-reports get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
