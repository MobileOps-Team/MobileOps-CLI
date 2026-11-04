package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var makesCmd = &cobra.Command{
	Use:   "makes",
	Short: "Manage makes (manufacturers)",
}

var makesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List makes",
	RunE:  runMakesList,
}

var makesGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single make by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runMakesGet,
}

var makesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new make",
	RunE:  runMakesCreate,
}

var makesUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a make",
	Args:  cobra.ExactArgs(1),
	RunE:  runMakesUpdate,
}

var makeFlags = map[string]string{
	"name":           "name",
	"address":        "address",
	"city":           "city",
	"postal-code":    "postal_code",
	"country":        "country",
	"state-province": "state_province",
	"email-address":  "email_address",
	"phone":          "phone",
}

func addMakeWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Make name")
	cmd.Flags().String("address", "", "Address")
	cmd.Flags().String("city", "", "City")
	cmd.Flags().String("postal-code", "", "Postal code")
	cmd.Flags().String("country", "", "Country")
	cmd.Flags().String("state-province", "", "State/Province")
	cmd.Flags().String("email-address", "", "Email address")
	cmd.Flags().String("phone", "", "Phone")
}

func init() {
	makesListCmd.Flags().Int("page", 1, "Page number")
	makesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addMakeWriteFlags(makesCreateCmd)
	addMakeWriteFlags(makesUpdateCmd)

	makesCmd.AddCommand(makesListCmd)
	makesCmd.AddCommand(makesGetCmd)
	makesCmd.AddCommand(makesCreateCmd)
	makesCmd.AddCommand(makesUpdateCmd)
}

func runMakesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("makes", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops makes get %s", 3)
	env := envelope.WrapCollection(response, "makes", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runMakesGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("makes/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Make", []string{
		fmt.Sprintf("mobileops models list --make-id %s", id),
		"mobileops makes list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runMakesCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, makeFlags)

	response, err := c.Post("makes", wrapBody("make", body))
	if err != nil {
		return handleClientError(err)
	}

	makeID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Make created", []string{
		fmt.Sprintf("mobileops makes get %s", makeID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runMakesUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, makeFlags)

	response, err := c.Put(fmt.Sprintf("makes/%s", id), wrapBody("make", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Make updated", []string{
		fmt.Sprintf("mobileops makes get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
