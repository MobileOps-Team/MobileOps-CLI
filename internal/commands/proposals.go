package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var proposalsCmd = &cobra.Command{
	Use:   "proposals",
	Short: "View sales proposals",
	Long: "View sales proposals with their deal fields and lifecycle rollup.\n\n" +
		"Outcome: won = Approved; lost = Rejected or Expired; every other status is open.\n" +
		"Quoted values and lifecycle amounts are integer US cents.",
}

var proposalsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List proposals, most recently updated first",
	RunE:  runProposalsList,
}

var proposalsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single proposal by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runProposalsGet,
}

var proposalListFlags = map[string]string{
	"status":        "status",
	"outcome":       "outcome",
	"customer-id":   "customer_id",
	"updated-since": "updated_since",
	"from":          "start_date",
	"to":            "end_date",
}

func init() {
	proposalsListCmd.Flags().Int("page", 1, "Page number")
	proposalsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	proposalsListCmd.Flags().String("status", "", "Filter by exact status (e.g. \"Pending Approval\")")
	proposalsListCmd.Flags().String("outcome", "", "Filter by outcome: open, won or lost")
	proposalsListCmd.Flags().String("customer-id", "", "Filter by customer ID")
	proposalsListCmd.Flags().String("updated-since", "", "Only proposals updated at or after this time (ISO 8601)")
	proposalsListCmd.Flags().String("from", "", "Service window ends on or after this date (YYYY-MM-DD)")
	proposalsListCmd.Flags().String("to", "", "Service window starts on or before this date (YYYY-MM-DD)")
	proposalsListCmd.Flags().Bool("include-archived", false, "Include archived proposals")

	proposalsCmd.AddCommand(proposalsListCmd)
	proposalsCmd.AddCommand(proposalsGetCmd)
}

func runProposalsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters(cmd, params, proposalListFlags)
	if v, _ := cmd.Flags().GetBool("include-archived"); v {
		params["include_archived"] = "true"
	}

	response, err := c.Get("proposals", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops proposals get %s", 3)
	env := envelope.WrapCollection(response, "proposals", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runProposalsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("proposals/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := []string{}
	if customerID, ok := response["customer_id"].(string); ok && customerID != "" {
		breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops customers get %s", customerID))
	}
	if jobIDs, ok := response["job_ids"].([]interface{}); ok {
		for i, jobID := range jobIDs {
			if i == 3 {
				break
			}
			breadcrumbs = append(breadcrumbs, fmt.Sprintf("mobileops jobs get %v", jobID))
		}
	}
	breadcrumbs = append(breadcrumbs, "mobileops proposals list")

	env := envelope.WrapRecord(response, "Proposal", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}
