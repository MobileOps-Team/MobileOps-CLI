package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var vesselDocumentsCmd = &cobra.Command{
	Use:   "vessel-documents",
	Short: "Manage vessel documents",
}

var vesselDocumentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List vessel documents",
	RunE:  runVesselDocumentsList,
}

var vesselDocumentsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single vessel document by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselDocumentsGet,
}

var vesselDocumentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vessel document",
	RunE:  runVesselDocumentsCreate,
}

var vesselDocumentsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a vessel document",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselDocumentsUpdate,
}

var vesselDocumentsDeleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a vessel document",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselDocumentsDelete,
}

var vesselDocumentsAttachmentsCmd = &cobra.Command{
	Use:   "attachments ID",
	Short: "Get attachments for a vessel document",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselDocumentsAttachments,
}

var vesselDocumentFlags = map[string]string{
	"issue-date":            "issue_date",
	"expire-date":           "expire_date",
	"endorsement-date":      "endorsement_date",
	"name":                  "name",
	"notification-group-id": "notification_group_id",
	"reminder-interval":     "reminder_interval",
	"user-id":               "user_id",
	"vessel-id":             "vessel_id",
	"notes":                 "notes",
	"template-id":           "template_id",
}

var vesselDocumentBoolFlags = map[string]string{
	"expire-notification": "expire_notification",
	"captain-edit":        "captain_edit",
	"confidential":        "confidential",
}

var vesselDocumentArrayFlags = map[string]string{
	"tags": "tags",
}

func addVesselDocumentWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("issue-date", "", "Issue date")
	cmd.Flags().String("expire-date", "", "Expire date")
	cmd.Flags().String("endorsement-date", "", "Endorsement date")
	cmd.Flags().String("name", "", "Document name")
	cmd.Flags().String("notification-group-id", "", "Notification group ID")
	cmd.Flags().String("reminder-interval", "", "Reminder interval")
	cmd.Flags().String("user-id", "", "User ID")
	cmd.Flags().String("vessel-id", "", "Vessel ID")
	cmd.Flags().String("notes", "", "Notes")
	cmd.Flags().String("template-id", "", "Template ID")
	cmd.Flags().Bool("expire-notification", false, "Expire notification")
	cmd.Flags().Bool("captain-edit", false, "Captain edit")
	cmd.Flags().Bool("confidential", false, "Confidential")
	cmd.Flags().String("tags", "", "Tags (comma-separated)")
}

func init() {
	vesselDocumentsListCmd.Flags().Int("page", 1, "Page number")
	vesselDocumentsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	vesselDocumentsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")

	addVesselDocumentWriteFlags(vesselDocumentsCreateCmd)
	addVesselDocumentWriteFlags(vesselDocumentsUpdateCmd)

	vesselDocumentsCmd.AddCommand(vesselDocumentsListCmd)
	vesselDocumentsCmd.AddCommand(vesselDocumentsGetCmd)
	vesselDocumentsCmd.AddCommand(vesselDocumentsCreateCmd)
	vesselDocumentsCmd.AddCommand(vesselDocumentsUpdateCmd)
	vesselDocumentsCmd.AddCommand(vesselDocumentsDeleteCmd)
	vesselDocumentsCmd.AddCommand(vesselDocumentsAttachmentsCmd)
}

func runVesselDocumentsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("vessel-id"); v != "" {
		params["vessel_id"] = v
	}

	response, err := c.Get("vessel-documents", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops vessel-documents get %s", 3)
	env := envelope.WrapCollection(response, "vessel documents", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselDocumentsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("vessel-documents/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vessel document", []string{
		fmt.Sprintf("mobileops vessel-documents attachments %s", id),
		"mobileops vessel-documents list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselDocumentsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vesselDocumentFlags)
	bodyFromBoolFlags(cmd, body, vesselDocumentBoolFlags)
	bodyFromArrayFlags(cmd, body, vesselDocumentArrayFlags)

	response, err := c.Post("vessel-documents", wrapBody("vessel_document", body))
	if err != nil {
		return handleClientError(err)
	}

	docID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Vessel document created", []string{
		fmt.Sprintf("mobileops vessel-documents get %s", docID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselDocumentsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vesselDocumentFlags)
	bodyFromBoolFlags(cmd, body, vesselDocumentBoolFlags)
	bodyFromArrayFlags(cmd, body, vesselDocumentArrayFlags)

	response, err := c.Put(fmt.Sprintf("vessel-documents/%s", id), wrapBody("vessel_document", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vessel document updated", []string{
		fmt.Sprintf("mobileops vessel-documents get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselDocumentsDelete(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Delete(fmt.Sprintf("vessel-documents/%s", id))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vessel document deleted", []string{
		"mobileops vessel-documents list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselDocumentsAttachments(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("vessel-documents/%s/attachments", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "attachments", []string{
		fmt.Sprintf("mobileops vessel-documents get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
