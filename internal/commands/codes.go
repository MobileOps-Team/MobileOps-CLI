package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var codesCmd = &cobra.Command{
	Use:   "codes",
	Short: "Manage codes (accounting/billing/job codes, company-wide or per vessel)",
}

var codesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List codes (archived codes are not returned)",
	RunE:  runCodesList,
}

var codesGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single code by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runCodesGet,
}

var codesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a code",
	RunE:  runCodesCreate,
}

var codesUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a code",
	Args:  cobra.ExactArgs(1),
	RunE:  runCodesUpdate,
}

var codeFlags = map[string]string{
	"name":        "name",
	"description": "description",
	"type":        "type",
	"vessel-id":   "vessel_id",
}

func addCodeWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Code name")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("type", "", "Category (the app offers Accounting, Billing, Job, Other)")
	cmd.Flags().String("vessel-id", "", "Scope the code to a vessel; omit for a company-wide code")
}

func init() {
	codesListCmd.Flags().Int("page", 1, "Page number")
	codesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addCodeWriteFlags(codesCreateCmd)
	addCodeWriteFlags(codesUpdateCmd)
	codesUpdateCmd.Flags().Bool("company-wide", false, "Clear the vessel scope (sends vessel_id: null)")

	codesCmd.AddCommand(codesListCmd)
	codesCmd.AddCommand(codesGetCmd)
	codesCmd.AddCommand(codesCreateCmd)
	codesCmd.AddCommand(codesUpdateCmd)
}

func runCodesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get("codes", paginationParams(cmd))
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops codes get %s", 3)
	env := envelope.WrapCollection(response, "codes", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runCodesGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("codes/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Code", []string{
		fmt.Sprintf("mobileops codes update %s --name \"...\"", id),
		"mobileops codes list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCodesCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, codeFlags)

	response, err := c.Post("codes", wrapBody("code", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Code created", []string{
		fmt.Sprintf("mobileops codes get %s", envelope.ExtractID(response)),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCodesUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, codeFlags)
	if v, _ := cmd.Flags().GetBool("company-wide"); v {
		body["vessel_id"] = nil
	}

	response, err := c.Put(fmt.Sprintf("codes/%s", id), wrapBody("code", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Code updated", []string{
		fmt.Sprintf("mobileops codes get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
