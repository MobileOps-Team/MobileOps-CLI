package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var customersCmd = &cobra.Command{
	Use:   "customers",
	Short: "Manage customers",
}

var customersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List customers",
	RunE:  runCustomersList,
}

var customersGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single customer by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runCustomersGet,
}

var customersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new customer",
	RunE:  runCustomersCreate,
}

var customersUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a customer",
	Args:  cobra.ExactArgs(1),
	RunE:  runCustomersUpdate,
}

var customerFlags = map[string]string{
	"name":           "name",
	"address":        "address",
	"city":           "city",
	"postal-code":    "postal_code",
	"state-province": "state_province",
	"country":        "country",
	"email-address":  "email_address",
	"phone":          "phone",
	"hex-color":      "hex_color",
}

var customerBoolFlags = map[string]string{
	"third-party": "third_party",
}

var customerArrayFlags = map[string]string{
	"add-vessels": "add_vessels",
}

func addCustomerWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Customer name")
	cmd.Flags().String("address", "", "Address")
	cmd.Flags().String("city", "", "City")
	cmd.Flags().String("postal-code", "", "Postal code")
	cmd.Flags().String("state-province", "", "State/Province")
	cmd.Flags().String("country", "", "Country")
	cmd.Flags().String("email-address", "", "Email address")
	cmd.Flags().String("phone", "", "Phone")
	cmd.Flags().String("hex-color", "", "Hex color")
	cmd.Flags().Bool("third-party", false, "Third party")
	cmd.Flags().String("add-vessels", "", "Vessel IDs to add (comma-separated)")
}

func init() {
	customersListCmd.Flags().Int("page", 1, "Page number")
	customersListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addCustomerWriteFlags(customersCreateCmd)
	addCustomerWriteFlags(customersUpdateCmd)

	customersCmd.AddCommand(customersListCmd)
	customersCmd.AddCommand(customersGetCmd)
	customersCmd.AddCommand(customersCreateCmd)
	customersCmd.AddCommand(customersUpdateCmd)
}

func runCustomersList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("customers", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops customers get %s", 3)
	env := envelope.WrapCollection(response, "customers", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runCustomersGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("customers/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Customer", []string{
		"mobileops customers list",
		"mobileops jobs list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCustomersCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, customerFlags)
	bodyFromBoolFlags(cmd, body, customerBoolFlags)
	bodyFromArrayFlags(cmd, body, customerArrayFlags)

	response, err := c.Post("customers", wrapBody("customer", body))
	if err != nil {
		return handleClientError(err)
	}

	customerID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Customer created", []string{
		fmt.Sprintf("mobileops customers get %s", customerID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runCustomersUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, customerFlags)
	bodyFromBoolFlags(cmd, body, customerBoolFlags)
	bodyFromArrayFlags(cmd, body, customerArrayFlags)

	response, err := c.Put(fmt.Sprintf("customers/%s", id), wrapBody("customer", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Customer updated", []string{
		fmt.Sprintf("mobileops customers get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
