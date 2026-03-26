package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var suppliersCmd = &cobra.Command{
	Use:   "suppliers",
	Short: "Manage suppliers",
}

var suppliersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List suppliers",
	RunE:  runSuppliersList,
}

var suppliersGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single supplier by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runSuppliersGet,
}

var suppliersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new supplier",
	RunE:  runSuppliersCreate,
}

var suppliersUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a supplier",
	Args:  cobra.ExactArgs(1),
	RunE:  runSuppliersUpdate,
}

var supplierFlags = map[string]string{
	"name":           "name",
	"notes":          "notes",
	"address":        "address",
	"city":           "city",
	"postal-code":    "postal_code",
	"country":        "country",
	"state-province": "state_province",
	"email-address":  "email_address",
	"phone":          "phone",
}

func addSupplierWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Supplier name")
	cmd.Flags().String("notes", "", "Notes")
	cmd.Flags().String("address", "", "Address")
	cmd.Flags().String("city", "", "City")
	cmd.Flags().String("postal-code", "", "Postal code")
	cmd.Flags().String("country", "", "Country")
	cmd.Flags().String("state-province", "", "State/Province")
	cmd.Flags().String("email-address", "", "Email address")
	cmd.Flags().String("phone", "", "Phone")
}

func init() {
	suppliersListCmd.Flags().Int("page", 1, "Page number")
	suppliersListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addSupplierWriteFlags(suppliersCreateCmd)
	addSupplierWriteFlags(suppliersUpdateCmd)

	suppliersCmd.AddCommand(suppliersListCmd)
	suppliersCmd.AddCommand(suppliersGetCmd)
	suppliersCmd.AddCommand(suppliersCreateCmd)
	suppliersCmd.AddCommand(suppliersUpdateCmd)
}

func runSuppliersList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("suppliers", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops suppliers get %s", 3)
	env := envelope.WrapCollection(response, "suppliers", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runSuppliersGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("suppliers/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Supplier", []string{
		fmt.Sprintf("mobileops terms list --supplier-id %s", id),
		"mobileops suppliers list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runSuppliersCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, supplierFlags)

	response, err := c.Post("suppliers", wrapBody("supplier", body))
	if err != nil {
		return handleClientError(err)
	}

	supplierID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Supplier created", []string{
		fmt.Sprintf("mobileops suppliers get %s", supplierID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runSuppliersUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, supplierFlags)

	response, err := c.Put(fmt.Sprintf("suppliers/%s", id), wrapBody("supplier", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Supplier updated", []string{
		fmt.Sprintf("mobileops suppliers get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
