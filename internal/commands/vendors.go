package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var vendorsCmd = &cobra.Command{
	Use:   "vendors",
	Short: "Manage vendors",
}

var vendorsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List vendors",
	RunE:  runVendorsList,
}

var vendorsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single vendor by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runVendorsGet,
}

var vendorsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vendor",
	RunE:  runVendorsCreate,
}

var vendorsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a vendor",
	Args:  cobra.ExactArgs(1),
	RunE:  runVendorsUpdate,
}

var vendorFlags = map[string]string{
	"name":           "name",
	"short-name":     "short_name",
	"address":        "address",
	"city":           "city",
	"postal-code":    "postal_code",
	"state-province": "state_province",
	"country":        "country",
	"email-address":  "email_address",
	"phone":          "phone",
	"rating":         "rating",
	"rating-notes":   "rating_notes",
}

func addVendorWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Vendor name")
	cmd.Flags().String("short-name", "", "Short name")
	cmd.Flags().String("address", "", "Address")
	cmd.Flags().String("city", "", "City")
	cmd.Flags().String("postal-code", "", "Postal code")
	cmd.Flags().String("state-province", "", "State/Province")
	cmd.Flags().String("country", "", "Country")
	cmd.Flags().String("email-address", "", "Email address")
	cmd.Flags().String("phone", "", "Phone")
	cmd.Flags().String("rating", "", "Rating")
	cmd.Flags().String("rating-notes", "", "Rating notes")
}

func init() {
	vendorsListCmd.Flags().Int("page", 1, "Page number")
	vendorsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addVendorWriteFlags(vendorsCreateCmd)
	addVendorWriteFlags(vendorsUpdateCmd)

	vendorsCmd.AddCommand(vendorsListCmd)
	vendorsCmd.AddCommand(vendorsGetCmd)
	vendorsCmd.AddCommand(vendorsCreateCmd)
	vendorsCmd.AddCommand(vendorsUpdateCmd)
}

func runVendorsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("vendors", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops vendors get %s", 3)
	env := envelope.WrapCollection(response, "vendors", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runVendorsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("vendors/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vendor", []string{
		"mobileops vendors list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVendorsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vendorFlags)

	response, err := c.Post("vendors", wrapBody("vendor", body))
	if err != nil {
		return handleClientError(err)
	}

	vendorID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Vendor created", []string{
		fmt.Sprintf("mobileops vendors get %s", vendorID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVendorsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vendorFlags)

	response, err := c.Put(fmt.Sprintf("vendors/%s", id), wrapBody("vendor", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vendor updated", []string{
		fmt.Sprintf("mobileops vendors get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
