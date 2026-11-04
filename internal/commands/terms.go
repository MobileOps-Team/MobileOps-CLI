package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var termsCmd = &cobra.Command{
	Use:   "terms",
	Short: "Manage supplier terms",
}

var termsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List terms",
	RunE:  runTermsList,
}

var termsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single term by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runTermsGet,
}

var termsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new term",
	RunE:  runTermsCreate,
}

var termsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a term",
	Args:  cobra.ExactArgs(1),
	RunE:  runTermsUpdate,
}

var termFlags = map[string]string{
	"make-id":        "make_id",
	"model-id":       "model_id",
	"model-title":    "model_title",
	"supplier-id":    "supplier_id",
	"price":          "price",
	"notes":          "notes",
	"date-beginning": "date_beginning",
	"date-ending":    "date_ending",
}

func addTermWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("make-id", "", "Make ID")
	cmd.Flags().String("model-id", "", "Model ID")
	cmd.Flags().String("model-title", "", "Model title")
	cmd.Flags().String("supplier-id", "", "Supplier ID")
	cmd.Flags().String("price", "", "Price")
	cmd.Flags().String("notes", "", "Notes")
	cmd.Flags().String("date-beginning", "", "Date beginning")
	cmd.Flags().String("date-ending", "", "Date ending")
}

func init() {
	termsListCmd.Flags().Int("page", 1, "Page number")
	termsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	termsListCmd.Flags().String("supplier-id", "", "Filter by supplier ID")

	termsGetCmd.Flags().String("supplier-id", "", "Parent supplier ID (required)")

	addTermWriteFlags(termsCreateCmd)
	addTermWriteFlags(termsUpdateCmd)

	termsCmd.AddCommand(termsListCmd)
	termsCmd.AddCommand(termsGetCmd)
	termsCmd.AddCommand(termsCreateCmd)
	termsCmd.AddCommand(termsUpdateCmd)
}

func runTermsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("supplier-id"); v != "" {
		params["supplier_id"] = v
	}

	response, err := c.Get("terms", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops terms get %s", 3)
	env := envelope.WrapCollection(response, "terms", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runTermsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := map[string]string{}
	if v, _ := cmd.Flags().GetString("supplier-id"); v != "" {
		params["supplier_id"] = v
	}

	response, err := c.Get(fmt.Sprintf("terms/%s", id), params)
	if err != nil {
		return handleClientError(err)
	}

	var breadcrumbs []string
	if supplierID, ok := response["supplier_id"].(string); ok && supplierID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops suppliers get %s", supplierID))
	}
	breadcrumbs = append(breadcrumbs, "mobileops terms list")

	env := envelope.WrapRecord(response, "Term", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runTermsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, termFlags)

	response, err := c.Post("terms", wrapBody("term", body))
	if err != nil {
		return handleClientError(err)
	}

	termID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Term created", []string{
		fmt.Sprintf("mobileops terms get %s", termID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runTermsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, termFlags)

	response, err := c.Put(fmt.Sprintf("terms/%s", id), wrapBody("term", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Term updated", []string{
		fmt.Sprintf("mobileops terms get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
