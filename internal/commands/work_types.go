package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var workTypesCmd = &cobra.Command{
	Use:   "work-types",
	Short: "Manage work types",
}

var workTypesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List work types",
	RunE:  runWorkTypesList,
}

var workTypesGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single work type by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkTypesGet,
}

var workTypesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new work type",
	RunE:  runWorkTypesCreate,
}

var workTypesUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a work type",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkTypesUpdate,
}

var workTypeFlags = map[string]string{
	"name": "name",
}

func addWorkTypeWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Work type name")
}

func init() {
	workTypesListCmd.Flags().Int("page", 1, "Page number")
	workTypesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addWorkTypeWriteFlags(workTypesCreateCmd)
	addWorkTypeWriteFlags(workTypesUpdateCmd)

	workTypesCmd.AddCommand(workTypesListCmd)
	workTypesCmd.AddCommand(workTypesGetCmd)
	workTypesCmd.AddCommand(workTypesCreateCmd)
	workTypesCmd.AddCommand(workTypesUpdateCmd)
}

func runWorkTypesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("work-types", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops work-types get %s", 3)
	env := envelope.WrapCollection(response, "work types", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runWorkTypesGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("work-types/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Work type", []string{
		"mobileops work-types list",
		"mobileops jobs list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runWorkTypesCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, workTypeFlags)

	response, err := c.Post("work-types", wrapBody("work_type", body))
	if err != nil {
		return handleClientError(err)
	}

	wtID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Work type created", []string{
		fmt.Sprintf("mobileops work-types get %s", wtID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runWorkTypesUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, workTypeFlags)

	response, err := c.Put(fmt.Sprintf("work-types/%s", id), wrapBody("work_type", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Work type updated", []string{
		fmt.Sprintf("mobileops work-types get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
