package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var auditsCmd = &cobra.Command{
	Use:   "audits",
	Short: "List audits (SIRE inspections, internal/external audits and surveys)",
}

var auditsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List audits, newest first; each embeds its observations, deficiencies and nonconformities",
	RunE:  runAuditsList,
}

func init() {
	auditsListCmd.Flags().Int("page", 1, "Page number")
	auditsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	auditsCmd.AddCommand(auditsListCmd)
}

func runAuditsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get("audits", paginationParams(cmd))
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops observations list --audit-id %s", 3)
	env := envelope.WrapCollection(response, "audits", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}
