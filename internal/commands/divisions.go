package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var divisionsCmd = &cobra.Command{
	Use:   "divisions",
	Short: "Manage divisions",
}

var divisionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List divisions",
	RunE:  runDivisionsList,
}

var divisionsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single division by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runDivisionsGet,
}

var divisionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new division",
	RunE:  runDivisionsCreate,
}

var divisionsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a division",
	Args:  cobra.ExactArgs(1),
	RunE:  runDivisionsUpdate,
}

var divisionFlags = map[string]string{
	"name":             "name",
	"description":      "description",
	"reference-number": "reference_number",
	"hex-color":        "hex_color",
}

var divisionBoolFlags = map[string]string{
	"active": "active",
}

func addDivisionWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Division name")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("reference-number", "", "Reference number")
	cmd.Flags().String("hex-color", "", "Hex color")
	cmd.Flags().Bool("active", false, "Active")
}

func init() {
	divisionsListCmd.Flags().Int("page", 1, "Page number")
	divisionsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addDivisionWriteFlags(divisionsCreateCmd)
	addDivisionWriteFlags(divisionsUpdateCmd)

	divisionsCmd.AddCommand(divisionsListCmd)
	divisionsCmd.AddCommand(divisionsGetCmd)
	divisionsCmd.AddCommand(divisionsCreateCmd)
	divisionsCmd.AddCommand(divisionsUpdateCmd)
}

func runDivisionsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("divisions", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops divisions get %s", 3)
	env := envelope.WrapCollection(response, "divisions", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runDivisionsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("divisions/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Division", []string{
		fmt.Sprintf("mobileops vessels list --division-id %s", id),
		"mobileops divisions list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runDivisionsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, divisionFlags)
	bodyFromBoolFlags(cmd, body, divisionBoolFlags)

	response, err := c.Post("divisions", wrapBody("division", body))
	if err != nil {
		return handleClientError(err)
	}

	divisionID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Division created", []string{
		fmt.Sprintf("mobileops divisions get %s", divisionID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runDivisionsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, divisionFlags)
	bodyFromBoolFlags(cmd, body, divisionBoolFlags)

	response, err := c.Put(fmt.Sprintf("divisions/%s", id), wrapBody("division", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Division updated", []string{
		fmt.Sprintf("mobileops divisions get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
