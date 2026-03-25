package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Manage jobs",
}

var jobsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List jobs",
	RunE:  runJobsList,
}

var jobsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single job by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runJobsGet,
}

func init() {
	jobsListCmd.Flags().Int("page", 1, "Page number")
	jobsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	jobsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	jobsListCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	jobsListCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
	jobsListCmd.Flags().Bool("active-only", false, "Only show active jobs")

	jobsCmd.AddCommand(jobsListCmd)
	jobsCmd.AddCommand(jobsGetCmd)
}

func runJobsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)

	if vesselID, _ := cmd.Flags().GetString("vessel-id"); vesselID != "" {
		params["vessel_id"] = vesselID
	}
	if from, _ := cmd.Flags().GetString("from"); from != "" {
		params["date_beginning"] = from
	}
	if to, _ := cmd.Flags().GetString("to"); to != "" {
		params["date_ending"] = to
	}
	if activeOnly, _ := cmd.Flags().GetBool("active-only"); activeOnly {
		params["scope_by_active_event_list_status"] = "true"
	}

	response, err := c.Get("jobs", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops jobs get %s", 3)
	env := envelope.WrapCollection(response, "jobs", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runJobsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("jobs/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Job", []string{
		"mobileops jobs list",
		"mobileops vessels list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
