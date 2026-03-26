package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var personnelDocumentsCmd = &cobra.Command{
	Use:   "personnel-documents",
	Short: "Manage personnel documents",
}

var personnelDocumentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List personnel documents",
	RunE:  runPersonnelDocumentsList,
}

var personnelDocumentsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single personnel document by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runPersonnelDocumentsGet,
}

var personnelDocumentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new personnel document",
	RunE:  runPersonnelDocumentsCreate,
}

var personnelDocumentsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a personnel document",
	Args:  cobra.ExactArgs(1),
	RunE:  runPersonnelDocumentsUpdate,
}

var personnelDocumentsDeleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a personnel document",
	Args:  cobra.ExactArgs(1),
	RunE:  runPersonnelDocumentsDelete,
}

var personnelDocumentsAttachmentsCmd = &cobra.Command{
	Use:   "attachments ID",
	Short: "Get attachments for a personnel document",
	Args:  cobra.ExactArgs(1),
	RunE:  runPersonnelDocumentsAttachments,
}

var personnelDocumentFlags = map[string]string{
	"user-id":     "user_id",
	"template-id": "template_id",
	"issue-date":  "issue_date",
	"expire-date": "expire_date",
	"notes":       "notes",
}

var personnelDocumentBoolFlags = map[string]string{
	"expire-notification": "expire_notification",
}

func addPersonnelDocumentWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("user-id", "", "User ID")
	cmd.Flags().String("template-id", "", "Template ID")
	cmd.Flags().String("issue-date", "", "Issue date")
	cmd.Flags().String("expire-date", "", "Expire date")
	cmd.Flags().String("notes", "", "Notes")
	cmd.Flags().Bool("expire-notification", false, "Expire notification")
}

func init() {
	personnelDocumentsListCmd.Flags().Int("page", 1, "Page number")
	personnelDocumentsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	personnelDocumentsListCmd.Flags().String("user-id", "", "Filter by user ID")
	personnelDocumentsListCmd.Flags().String("employee-number", "", "Filter by employee number")
	personnelDocumentsListCmd.Flags().String("template-id", "", "Filter by template ID")
	personnelDocumentsListCmd.Flags().Bool("include-ultra-fields", false, "Include custom ultra fields")

	addPersonnelDocumentWriteFlags(personnelDocumentsCreateCmd)
	addPersonnelDocumentWriteFlags(personnelDocumentsUpdateCmd)

	personnelDocumentsCmd.AddCommand(personnelDocumentsListCmd)
	personnelDocumentsCmd.AddCommand(personnelDocumentsGetCmd)
	personnelDocumentsCmd.AddCommand(personnelDocumentsCreateCmd)
	personnelDocumentsCmd.AddCommand(personnelDocumentsUpdateCmd)
	personnelDocumentsCmd.AddCommand(personnelDocumentsDeleteCmd)
	personnelDocumentsCmd.AddCommand(personnelDocumentsAttachmentsCmd)
}

func runPersonnelDocumentsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("user-id"); v != "" {
		params["user_id"] = v
	}
	if v, _ := cmd.Flags().GetString("employee-number"); v != "" {
		params["employee_number"] = v
	}
	if v, _ := cmd.Flags().GetString("template-id"); v != "" {
		params["template_id"] = v
	}
	if v, _ := cmd.Flags().GetBool("include-ultra-fields"); v {
		params["include_ultra_fields"] = "true"
	}

	response, err := c.Get("personnel-documents", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops personnel-documents get %s", 3)
	env := envelope.WrapCollection(response, "personnel documents", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runPersonnelDocumentsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("personnel-documents/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Personnel document", []string{
		fmt.Sprintf("mobileops personnel-documents attachments %s", id),
		"mobileops personnel-documents list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runPersonnelDocumentsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, personnelDocumentFlags)
	bodyFromBoolFlags(cmd, body, personnelDocumentBoolFlags)

	response, err := c.Post("personnel-documents", wrapBody("personnel_document", body))
	if err != nil {
		return handleClientError(err)
	}

	docID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Personnel document created", []string{
		fmt.Sprintf("mobileops personnel-documents get %s", docID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runPersonnelDocumentsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, personnelDocumentFlags)
	bodyFromBoolFlags(cmd, body, personnelDocumentBoolFlags)

	response, err := c.Put(fmt.Sprintf("personnel-documents/%s", id), wrapBody("personnel_document", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Personnel document updated", []string{
		fmt.Sprintf("mobileops personnel-documents get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runPersonnelDocumentsDelete(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Delete(fmt.Sprintf("personnel-documents/%s", id))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Personnel document deleted", []string{
		"mobileops personnel-documents list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runPersonnelDocumentsAttachments(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("personnel-documents/%s/attachments", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "attachments", []string{
		fmt.Sprintf("mobileops personnel-documents get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
