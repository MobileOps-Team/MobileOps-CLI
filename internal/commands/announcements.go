package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var announcementsCmd = &cobra.Command{
	Use:   "announcements",
	Short: "View announcements and their acknowledgements",
	Long: "View announcements with the employees who acknowledged each one.\n" +
		"Drafts and archived announcements are included unless --published or\n" +
		"--archived is set.",
}

var announcementsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List announcements, newest first",
	RunE:  runAnnouncementsList,
}

var announcementsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single announcement by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runAnnouncementsGet,
}

var announcementListFlags = map[string]string{
	"from":          "start_date",
	"to":            "end_date",
	"updated-since": "updated_since",
}

var announcementListBoolFlags = map[string]string{
	"published": "published",
	"safety":    "safety",
	"archived":  "archived",
}

func init() {
	announcementsListCmd.Flags().Int("page", 1, "Page number")
	announcementsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	announcementsListCmd.Flags().String("from", "", "Announcements created on or after this date (YYYY-MM-DD, UTC)")
	announcementsListCmd.Flags().String("to", "", "Announcements created on or before this date (YYYY-MM-DD, UTC)")
	announcementsListCmd.Flags().String("updated-since", "", "Only announcements updated at or after this time (ISO 8601)")
	announcementsListCmd.Flags().Bool("published", false, "Only published announcements; --published=false for only drafts")
	announcementsListCmd.Flags().Bool("safety", false, "Only safety announcements; --safety=false for only the others")
	announcementsListCmd.Flags().Bool("archived", false, "Only archived announcements; --archived=false for only unarchived")

	announcementsCmd.AddCommand(announcementsListCmd)
	announcementsCmd.AddCommand(announcementsGetCmd)
}

func runAnnouncementsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters(cmd, params, announcementListFlags)
	boolListFilters(cmd, params, announcementListBoolFlags)

	response, err := c.Get("announcements", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops announcements get %s", 3)
	env := envelope.WrapCollection(response, "announcements", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runAnnouncementsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("announcements/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := []string{"mobileops announcements list"}

	env := envelope.WrapRecord(response, "Announcement", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}
