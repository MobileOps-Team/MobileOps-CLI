package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var crewCmd = &cobra.Command{
	Use:   "crew",
	Short: "Manage crew members",
}

var crewListCmd = &cobra.Command{
	Use:   "list",
	Short: "List crew members",
	RunE:  runCrewList,
}

var crewGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single crew member by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runCrewGet,
}

var crewCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new crew member",
	RunE:  runCrewCreate,
}

var crewUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a crew member",
	Args:  cobra.ExactArgs(1),
	RunE:  runCrewUpdate,
}

var crewDeleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a crew member",
	Args:  cobra.ExactArgs(1),
	RunE:  runCrewDelete,
}

var crewFindByEmployeeNumberCmd = &cobra.Command{
	Use:   "find-by-employee-number NUMBER",
	Short: "Find a crew member by employee number",
	Args:  cobra.ExactArgs(1),
	RunE:  runCrewFindByEmployeeNumber,
}

var crewFlags = map[string]string{
	"first-name":      "first_name",
	"last-name":       "last_name",
	"email":           "email",
	"password":        "password",
	"phone":           "phone",
	"time-zone":       "time_zone",
	"address":         "address",
	"birthday":        "birthday",
	"passport-number": "passport_number",
	"mmc-number":      "mmc_number",
	"employee-number": "employee_number",
}

var crewBoolFlags = map[string]string{
	"receive-notifications": "receive_notifications",
	"login-disabled":        "login_disabled",
	"archived":              "archived",
}

var crewArrayFlags = map[string]string{
	"division-ids":       "division_ids",
	"employee-positions": "employee_positions",
	"primary-assets":     "primary_assets",
}

func addCrewWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("first-name", "", "First name")
	cmd.Flags().String("last-name", "", "Last name")
	cmd.Flags().String("email", "", "Email")
	cmd.Flags().String("password", "", "Password")
	cmd.Flags().String("phone", "", "Phone")
	cmd.Flags().String("time-zone", "", "Time zone")
	cmd.Flags().String("address", "", "Address")
	cmd.Flags().String("birthday", "", "Birthday")
	cmd.Flags().String("passport-number", "", "Passport number")
	cmd.Flags().String("mmc-number", "", "MMC number")
	cmd.Flags().String("employee-number", "", "Employee number")
	cmd.Flags().Bool("receive-notifications", false, "Receive notifications")
	cmd.Flags().Bool("login-disabled", false, "Disable login")
	cmd.Flags().Bool("archived", false, "Archived")
	cmd.Flags().String("division-ids", "", "Division IDs (comma-separated)")
	cmd.Flags().String("employee-positions", "", "Employee positions (comma-separated)")
	cmd.Flags().String("primary-assets", "", "Primary asset IDs (comma-separated)")
}

func init() {
	crewListCmd.Flags().Int("page", 1, "Page number")
	crewListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	crewListCmd.Flags().Bool("include-ultra-fields", false, "Include custom ultra fields")

	addCrewWriteFlags(crewCreateCmd)
	addCrewWriteFlags(crewUpdateCmd)

	crewCmd.AddCommand(crewListCmd)
	crewCmd.AddCommand(crewGetCmd)
	crewCmd.AddCommand(crewCreateCmd)
	crewCmd.AddCommand(crewUpdateCmd)
	crewCmd.AddCommand(crewDeleteCmd)
	crewCmd.AddCommand(crewFindByEmployeeNumberCmd)
}

func runCrewList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetBool("include-ultra-fields"); v {
		params["include_ultra_fields"] = "true"
	}

	response, err := c.Get("users", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops crew get %s", 3)
	env := envelope.WrapCollection(response, "crew members", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runCrewGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("users/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Crew member", []string{
		"mobileops crew list",
		"mobileops jobs list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCrewCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, crewFlags)
	bodyFromBoolFlags(cmd, body, crewBoolFlags)
	bodyFromArrayFlags(cmd, body, crewArrayFlags)

	response, err := c.Post("users", wrapBody("user", body))
	if err != nil {
		return handleClientError(err)
	}

	userID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Crew member created", []string{
		fmt.Sprintf("mobileops crew get %s", userID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCrewUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, crewFlags)
	bodyFromBoolFlags(cmd, body, crewBoolFlags)
	bodyFromArrayFlags(cmd, body, crewArrayFlags)

	response, err := c.Put(fmt.Sprintf("users/%s", id), wrapBody("user", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Crew member updated", []string{
		fmt.Sprintf("mobileops crew get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCrewDelete(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Delete(fmt.Sprintf("users/%s", id))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Crew member deleted", []string{
		"mobileops crew list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCrewFindByEmployeeNumber(cmd *cobra.Command, args []string) error {
	number := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get("users/by_employee_number", map[string]string{"employee_number": number})
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Crew member", []string{
		"mobileops crew list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
