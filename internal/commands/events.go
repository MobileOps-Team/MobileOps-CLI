package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Search wheelhouse events (job/order timeline entries)",
}

var eventsSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search events, newest first; every filter is optional and comma-separated",
	Long: `Search your company's events. The response also carries
wheelhouse_template_event_fields, the field definitions needed to read
each event's values. Events cannot be created or edited through the API.`,
	RunE: runEventsSearch,
}

var eventSearchArrayFlags = map[string]string{
	"job-ref-numbers":   "job_ref_numbers",
	"job-ids":           "job_ids",
	"order-ref-numbers": "order_ref_numbers",
	"order-ids":         "order_ids",
	"types":             "types",
	"work-type-ids":     "work_type_ids",
	"event-type-ids":    "event_type_ids",
	"vessel-ids":        "acting_asset_ids",
}

func init() {
	eventsSearchCmd.Flags().Int("page", 1, "Page number")
	eventsSearchCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	eventsSearchCmd.Flags().String("job-ref-numbers", "", "Job reference numbers (comma-separated); takes precedence over --job-ids")
	eventsSearchCmd.Flags().String("job-ids", "", "Job IDs (comma-separated)")
	eventsSearchCmd.Flags().String("order-ref-numbers", "", "Order reference numbers (comma-separated); takes precedence over --order-ids")
	eventsSearchCmd.Flags().String("order-ids", "", "Order IDs (comma-separated)")
	eventsSearchCmd.Flags().String("types", "", "Event types (comma-separated)")
	eventsSearchCmd.Flags().String("work-type-ids", "", "Work type IDs (comma-separated)")
	eventsSearchCmd.Flags().String("event-type-ids", "", "Event type IDs (comma-separated); takes precedence over --work-type-ids")
	eventsSearchCmd.Flags().String("vessel-ids", "", "Acting vessel IDs (comma-separated)")

	eventsCmd.AddCommand(eventsSearchCmd)
}

func runEventsSearch(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := map[string]interface{}{}
	bodyFromArrayFlags(cmd, body, eventSearchArrayFlags)

	response, err := c.PostWithParams("events/search", paginationParams(cmd), body)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "events", []string{
		"mobileops jobs list",
		"mobileops work-types list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
