package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var employeePositionsCmd = &cobra.Command{
	Use:   "employee-positions",
	Short: "Manage employee positions",
}

var employeePositionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List employee positions",
	RunE:  runEmployeePositionsList,
}

var employeePositionsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single employee position by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runEmployeePositionsGet,
}

var employeePositionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new employee position",
	RunE:  runEmployeePositionsCreate,
}

var employeePositionsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update an employee position",
	Args:  cobra.ExactArgs(1),
	RunE:  runEmployeePositionsUpdate,
}

var employeePositionFlags = map[string]string{
	"name":             "name",
	"description":      "description",
	"reference-number": "reference_number",
	"hex-color":        "hex_color",
}

var employeePositionBoolFlags = map[string]string{
	"active": "active",
}

func addEmployeePositionWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Position name")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("reference-number", "", "Reference number")
	cmd.Flags().String("hex-color", "", "Hex color")
	cmd.Flags().Bool("active", false, "Active")
}

func init() {
	employeePositionsListCmd.Flags().Int("page", 1, "Page number")
	employeePositionsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addEmployeePositionWriteFlags(employeePositionsCreateCmd)
	addEmployeePositionWriteFlags(employeePositionsUpdateCmd)

	employeePositionsCmd.AddCommand(employeePositionsListCmd)
	employeePositionsCmd.AddCommand(employeePositionsGetCmd)
	employeePositionsCmd.AddCommand(employeePositionsCreateCmd)
	employeePositionsCmd.AddCommand(employeePositionsUpdateCmd)
}

func runEmployeePositionsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("employee-positions", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops employee-positions get %s", 3)
	env := envelope.WrapCollection(response, "employee positions", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runEmployeePositionsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("employee-positions/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Employee position", []string{
		"mobileops employee-positions list",
		"mobileops crew list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runEmployeePositionsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, employeePositionFlags)
	bodyFromBoolFlags(cmd, body, employeePositionBoolFlags)

	response, err := c.Post("employee-positions", wrapBody("employee_position", body))
	if err != nil {
		return handleClientError(err)
	}

	epID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Employee position created", []string{
		fmt.Sprintf("mobileops employee-positions get %s", epID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runEmployeePositionsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, employeePositionFlags)
	bodyFromBoolFlags(cmd, body, employeePositionBoolFlags)

	response, err := c.Put(fmt.Sprintf("employee-positions/%s", id), wrapBody("employee_position", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Employee position updated", []string{
		fmt.Sprintf("mobileops employee-positions get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
