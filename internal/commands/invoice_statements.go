package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var invoiceStatementsCmd = &cobra.Command{
	Use:   "invoice-statements",
	Short: "Manage invoice statements",
}

var invoiceStatementsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List invoice statements",
	RunE:  runInvoiceStatementsList,
}

var invoiceStatementsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update an invoice statement",
	Args:  cobra.ExactArgs(1),
	RunE:  runInvoiceStatementsUpdate,
}

var invoiceStatementFlags = map[string]string{
	"integration-status": "integration_status",
}

func init() {
	invoiceStatementsListCmd.Flags().Int("page", 1, "Page number")
	invoiceStatementsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	invoiceStatementsListCmd.Flags().String("status", "", "Filter by status")
	invoiceStatementsListCmd.Flags().String("integration-status", "", "Filter by integration status")
	invoiceStatementsListCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	invoiceStatementsListCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
	invoiceStatementsListCmd.Flags().Bool("include-jobs", false, "Include job data")

	invoiceStatementsUpdateCmd.Flags().String("integration-status", "", "Integration status")

	invoiceStatementsCmd.AddCommand(invoiceStatementsListCmd)
	invoiceStatementsCmd.AddCommand(invoiceStatementsUpdateCmd)
}

func runInvoiceStatementsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("status"); v != "" {
		params["status"] = v
	}
	if v, _ := cmd.Flags().GetString("integration-status"); v != "" {
		params["integration_status"] = v
	}
	if v, _ := cmd.Flags().GetString("from"); v != "" {
		params["start"] = v
	}
	if v, _ := cmd.Flags().GetString("to"); v != "" {
		params["end"] = v
	}
	if v, _ := cmd.Flags().GetBool("include-jobs"); v {
		params["include_jobs"] = "true"
	}

	response, err := c.Get("invoice-statements", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "invoice statements", []string{
		"mobileops jobs list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runInvoiceStatementsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, invoiceStatementFlags)

	response, err := c.Put(fmt.Sprintf("invoice-statements/%s", id), wrapBody("invoice_statement", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Invoice statement updated", []string{
		"mobileops invoice-statements list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
