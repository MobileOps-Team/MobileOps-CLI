package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var shiftsCmd = &cobra.Command{
	Use:   "shifts",
	Short: "Manage crew shifts (crew schedule)",
}

var shiftsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List shifts (newest start first)",
	RunE:  runShiftsList,
}

var shiftsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single shift by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runShiftsGet,
}

var shiftsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Schedule a crew member for a shift (--notify emails and texts them)",
	RunE:  runShiftsCreate,
}

var shiftsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a shift (re-notifies the crew member when the shift has notify set)",
	Args:  cobra.ExactArgs(1),
	RunE:  runShiftsUpdate,
}

var shiftsDeleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a shift without notifying the crew member",
	Args:  cobra.ExactArgs(1),
	RunE:  runShiftsDelete,
}

var shiftFlags = map[string]string{
	"user-id":              "user_id",
	"vessel-id":            "vessel_id",
	"from":                 "date_beginning",
	"to":                   "date_ending",
	"employee-position-id": "employee_position_id",
	"employee-rate-id":     "employee_rate_id",
	"shift-type-id":        "shift_type_id",
	"notes":                "notes",
}

var shiftBoolFlags = map[string]string{
	"notify": "notify",
	"cancel": "cancel",
}

var shiftArrayFlags = map[string]string{
	"wh-extra-rate-ids": "wh_extra_emp_rate_ids",
}

func addShiftWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("user-id", "", "Crew member (user) ID")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("from", "", "Shift start (ISO 8601)")
	cmd.Flags().String("to", "", "Shift end (ISO 8601)")
	cmd.Flags().String("employee-position-id", "", "Employee position ID")
	cmd.Flags().String("employee-rate-id", "", "Employee rate ID")
	cmd.Flags().String("shift-type-id", "", "Shift type ID")
	cmd.Flags().String("wh-extra-rate-ids", "", "Wheelhouse extra employee rate IDs (comma-separated)")
	cmd.Flags().String("notes", "", "Notes shown to the crew member")
	cmd.Flags().Bool("notify", false, "Email and text the crew member about the shift and later changes")
	cmd.Flags().Bool("cancel", false, "Mark the shift cancelled (--cancel=false restores it)")
}

func init() {
	shiftsListCmd.Flags().Int("page", 1, "Page number")
	shiftsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	shiftsListCmd.Flags().String("from", "", "Only shifts ending on or after this time (YYYY-MM-DD or ISO 8601)")
	shiftsListCmd.Flags().String("to", "", "Only shifts starting on or before this time (YYYY-MM-DD or ISO 8601)")
	shiftsListCmd.Flags().String("user-id", "", "Filter by crew member (user) ID")
	shiftsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	shiftsListCmd.Flags().String("employee-position-id", "", "Filter by employee position ID")

	addShiftWriteFlags(shiftsCreateCmd)
	addShiftWriteFlags(shiftsUpdateCmd)
	shiftsCreateCmd.MarkFlagRequired("user-id")
	shiftsCreateCmd.MarkFlagRequired("from")
	shiftsCreateCmd.MarkFlagRequired("to")

	shiftsCmd.AddCommand(shiftsListCmd)
	shiftsCmd.AddCommand(shiftsGetCmd)
	shiftsCmd.AddCommand(shiftsCreateCmd)
	shiftsCmd.AddCommand(shiftsUpdateCmd)
	shiftsCmd.AddCommand(shiftsDeleteCmd)
}

func runShiftsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters(cmd, params, map[string]string{
		"from":                 "start_date",
		"to":                   "end_date",
		"user-id":              "user_id",
		"vessel-id":            "vessel_id",
		"employee-position-id": "employee_position_id",
	})

	response, err := c.Get("shifts", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops shifts get %s", 3)
	env := envelope.WrapCollection(response, "shifts", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runShiftsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("shifts/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := []string{"mobileops shifts list"}
	if userID, ok := response["user_id"].(string); ok && userID != "" {
		breadcrumbs = append([]string{fmt.Sprintf("mobileops crew get %s", userID)}, breadcrumbs...)
	}

	env := envelope.WrapRecord(response, "Shift", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func shiftBody(cmd *cobra.Command) map[string]interface{} {
	body := bodyFromFlags(cmd, shiftFlags)
	bodyFromBoolFlags(cmd, body, shiftBoolFlags)
	bodyFromArrayFlags(cmd, body, shiftArrayFlags)
	return body
}

func runShiftsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Post("shifts", wrapBody("shift", shiftBody(cmd)))
	if err != nil {
		return handleClientError(err)
	}

	shiftID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Shift created", []string{
		fmt.Sprintf("mobileops shifts get %s", shiftID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runShiftsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Put(fmt.Sprintf("shifts/%s", id), wrapBody("shift", shiftBody(cmd)))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Shift updated", []string{
		fmt.Sprintf("mobileops shifts get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runShiftsDelete(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Delete(fmt.Sprintf("shifts/%s", id))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Shift deleted", []string{
		"mobileops shifts list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
