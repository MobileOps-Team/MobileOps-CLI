package commands

import (
	"fmt"
	"strings"

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

var jobsVesselAvailabilityCmd = &cobra.Command{
	Use:   "vessel-availability",
	Short: "List jobs already scheduled for a vessel in a date window (empty = vessel is free)",
	RunE:  runJobsVesselAvailability,
}

var jobsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Generate job report rows through the reporting engine (not Job objects)",
	RunE:  runJobsExport,
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

var jobJSONFlags = map[string]string{
	"customer-vessels": "customer_vessels",
	"booking-requests": "booking_requests",
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
	cmd.Flags().String("customer-vessels", "", "Customer-owned vessels as JSON, e.g. '[{\"customer_id\":\"...\",\"vessel_key\":\"...\"}]' (keys from the customer's vessels map)")
	cmd.Flags().String("booking-requests", "", "OpenTug booking requests as JSON, e.g. '[{\"id\":\"...\"}]'")
}

func init() {
	jobsListCmd.Flags().Int("page", 1, "Page number")
	jobsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	jobsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	jobsListCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	jobsListCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
	jobsListCmd.Flags().Bool("active-only", false, "Only show active jobs")
	jobsListCmd.Flags().String("status", "", "Event list statuses (comma-separated): NOT_READY, READY, LAUNCHED, IN_PROGRESS, COMPLETE, JOB_CANCELED, JOB_RESCHEDULED")

	addJobWriteFlags(jobsCreateCmd)
	addJobWriteFlags(jobsUpdateCmd)

	jobsVesselAvailabilityCmd.Flags().String("vessel-id", "", "Vessel ID (required)")
	jobsVesselAvailabilityCmd.Flags().String("from", "", "Window start (ISO 8601, required)")
	jobsVesselAvailabilityCmd.Flags().String("to", "", "Window end (ISO 8601, required)")
	jobsVesselAvailabilityCmd.MarkFlagRequired("vessel-id")
	jobsVesselAvailabilityCmd.MarkFlagRequired("from")
	jobsVesselAvailabilityCmd.MarkFlagRequired("to")

	jobsExportCmd.Flags().String("from", "", "Report start date (YYYY-MM-DD, required)")
	jobsExportCmd.Flags().String("to", "", "Report end date (YYYY-MM-DD, required)")
	jobsExportCmd.Flags().String("dispatch-segment-template-id", "", "Restrict to jobs using this dispatch segment template")
	jobsExportCmd.Flags().Bool("simple", false, "Simple dispatch data report (default is the detailed report)")
	jobsExportCmd.Flags().Bool("actuals", false, "Planned vs actual job segment report")
	jobsExportCmd.MarkFlagRequired("from")
	jobsExportCmd.MarkFlagRequired("to")

	jobsCmd.AddCommand(jobsListCmd)
	jobsCmd.AddCommand(jobsGetCmd)
	jobsCmd.AddCommand(jobsCreateCmd)
	jobsCmd.AddCommand(jobsUpdateCmd)
	jobsCmd.AddCommand(jobsVesselAvailabilityCmd)
	jobsCmd.AddCommand(jobsExportCmd)
}

func runJobsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if vesselID, _ := cmd.Flags().GetString("vessel-id"); vesselID != "" {
		params["job_asset_id"] = vesselID
	}
	if status, _ := cmd.Flags().GetString("status"); status != "" {
		params["event_list_status"] = strings.ToUpper(strings.ReplaceAll(status, " ", ""))
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
	if err := bodyFromJSONFlags(cmd, body, jobJSONFlags); err != nil {
		return err
	}

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
	if err := bodyFromJSONFlags(cmd, body, jobJSONFlags); err != nil {
		return err
	}

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

func runJobsVesselAvailability(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	vesselID, _ := cmd.Flags().GetString("vessel-id")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	response, err := c.Post("jobs/vessel-availability", map[string]interface{}{
		"vessel_id":      vesselID,
		"date_beginning": from,
		"date_ending":    to,
	})
	if err != nil {
		return handleClientError(err)
	}

	jobs, _ := response["jobs"].([]interface{})
	summary := fmt.Sprintf("%d overlapping jobs for vessel %s", len(jobs), vesselID)
	if len(jobs) == 0 {
		summary = fmt.Sprintf("Vessel %s is free between %s and %s", vesselID, from, to)
	}
	env := envelope.Wrap(jobs, summary, []string{
		fmt.Sprintf("mobileops jobs list --vessel-id %s", vesselID),
		fmt.Sprintf("mobileops vessels get %s", vesselID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runJobsExport(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	body := map[string]interface{}{
		"date_beginning": from,
		"date_ending":    to,
	}
	if v, _ := cmd.Flags().GetString("dispatch-segment-template-id"); v != "" {
		body["dispatch_segment_template_id"] = v
	}
	bodyFromBoolFlags(cmd, body, map[string]string{"simple": "simple", "actuals": "actuals"})

	response, err := c.Post("jobs/job-export", body)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "job report rows", []string{
		"mobileops jobs list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
