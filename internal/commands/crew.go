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

var crewFindByEmployeeNumberCmd = &cobra.Command{
	Use:   "find-by-employee-number NUMBER",
	Short: "Find a crew member by employee number",
	Args:  cobra.ExactArgs(1),
	RunE:  runCrewFindByEmployeeNumber,
}

func init() {
	crewListCmd.Flags().Int("page", 1, "Page number")
	crewListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	crewCmd.AddCommand(crewListCmd)
	crewCmd.AddCommand(crewGetCmd)
	crewCmd.AddCommand(crewFindByEmployeeNumberCmd)
}

func runCrewList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
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
