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
	Short: "List suppliers (read-only)",
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

func init() {
	suppliersListCmd.Flags().Int("page", 1, "Page number")
	suppliersListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	suppliersCmd.AddCommand(suppliersListCmd)
	suppliersCmd.AddCommand(suppliersGetCmd)
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
