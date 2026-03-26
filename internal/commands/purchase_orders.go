package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var purchaseOrdersCmd = &cobra.Command{
	Use:   "purchase-orders",
	Short: "List purchase orders",
}

var purchaseOrdersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List purchase orders",
	RunE:  runPurchaseOrdersList,
}

func init() {
	purchaseOrdersListCmd.Flags().Int("page", 1, "Page number")
	purchaseOrdersListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	purchaseOrdersListCmd.Flags().String("status", "", "Filter by status")
	purchaseOrdersListCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	purchaseOrdersListCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")

	purchaseOrdersCmd.AddCommand(purchaseOrdersListCmd)
}

func runPurchaseOrdersList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("status"); v != "" {
		params["status"] = v
	}
	if v, _ := cmd.Flags().GetString("from"); v != "" {
		params["start"] = v
	}
	if v, _ := cmd.Flags().GetString("to"); v != "" {
		params["end"] = v
	}

	response, err := c.Get("purchase-orders", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "purchase orders", []string{
		"mobileops suppliers list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
