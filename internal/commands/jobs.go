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

var jobsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new job",
	RunE:  runJobsCreate,
}

var jobsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a job",
	Args:  cobra.ExactArgs(1),
	RunE:  runJobsUpdate,
}

var jobFlags = map[string]string{
	"from":        "date_beginning",
	"to":          "date_ending",
	"notes":       "notes",
	"user-id":     "user_id",
	"division-id": "division_id",
	"po-number":   "po_number",
	"ref-number":  "ref_number",
}

var jobBoolFlags = map[string]string{
	"tbd":    "tbd",
	"cancel": "cancel",
}

var jobArrayFlags = map[string]string{
	"vessels":    "vessels",
	"customers":  "customers",
	"work-types": "work_types",
	"locations":  "locations",
}

func addJobWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
	cmd.Flags().String("notes", "", "Notes")
	cmd.Flags().String("user-id", "", "User ID")
	cmd.Flags().String("division-id", "", "Division ID")
	cmd.Flags().String("po-number", "", "PO number")
	cmd.Flags().String("ref-number", "", "Reference number")
	cmd.Flags().Bool("tbd", false, "TBD")
	cmd.Flags().Bool("cancel", false, "Cancel job")
	cmd.Flags().String("vessels", "", "Vessel IDs (comma-separated)")
	cmd.Flags().String("customers", "", "Customer IDs (comma-separated)")
	cmd.Flags().String("work-types", "", "Work type IDs (comma-separated)")
	cmd.Flags().String("locations", "", "Location IDs (comma-separated)")
}

func init() {
	jobsListCmd.Flags().Int("page", 1, "Page number")
	jobsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	jobsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	jobsListCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	jobsListCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
	jobsListCmd.Flags().Bool("active-only", false, "Only show active jobs")

	addJobWriteFlags(jobsCreateCmd)
	addJobWriteFlags(jobsUpdateCmd)

	jobsCmd.AddCommand(jobsListCmd)
	jobsCmd.AddCommand(jobsGetCmd)
	jobsCmd.AddCommand(jobsCreateCmd)
	jobsCmd.AddCommand(jobsUpdateCmd)
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

func runJobsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, jobFlags)
	bodyFromBoolFlags(cmd, body, jobBoolFlags)
	bodyFromArrayFlags(cmd, body, jobArrayFlags)

	response, err := c.Post("jobs", wrapBody("job", body))
	if err != nil {
		return handleClientError(err)
	}

	jobID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Job created", []string{
		fmt.Sprintf("mobileops jobs get %s", jobID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runJobsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, jobFlags)
	bodyFromBoolFlags(cmd, body, jobBoolFlags)
	bodyFromArrayFlags(cmd, body, jobArrayFlags)

	response, err := c.Put(fmt.Sprintf("jobs/%s", id), wrapBody("job", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Job updated", []string{
		fmt.Sprintf("mobileops jobs get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
