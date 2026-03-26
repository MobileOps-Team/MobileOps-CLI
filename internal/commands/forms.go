package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var formsCmd = &cobra.Command{
	Use:   "forms",
	Short: "List form templates",
}

var formsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List form templates",
	RunE:  runFormsList,
}

func init() {
	formsListCmd.Flags().Int("page", 1, "Page number")
	formsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	formsListCmd.Flags().String("name", "", "Filter by form name")

	formsCmd.AddCommand(formsListCmd)
}

func runFormsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("name"); v != "" {
		params["form_name"] = v
	}

	response, err := c.Get("forms", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "forms", []string{
		"mobileops form-instances list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
