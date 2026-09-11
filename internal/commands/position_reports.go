package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var positionReportsCmd = &cobra.Command{
	Use:   "position-reports",
	Short: "Manage vessel position reports",
}

var positionReportsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List position reports",
	RunE:  runPositionReportsList,
}

var positionReportsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new position report",
	RunE:  runPositionReportsCreate,
}

var positionReportFlags = map[string]string{
	"asset-id":  "asset_id",
	"latitude":  "latitude",
	"longitude": "longitude",
	"heading":   "heading",
}

func addPositionReportWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("asset-id", "", "Asset/Vessel ID (required)")
	cmd.Flags().String("latitude", "", "Latitude, -90 to 90 (required)")
	cmd.Flags().String("longitude", "", "Longitude, -180 to 180 (required)")
	cmd.Flags().String("heading", "", "Heading")
	cmd.MarkFlagRequired("asset-id")
	cmd.MarkFlagRequired("latitude")
	cmd.MarkFlagRequired("longitude")
}

func init() {
	positionReportsListCmd.Flags().Int("page", 1, "Page number")
	positionReportsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	positionReportsListCmd.Flags().String("vessel-id", "", "Filter by vessel ID")
	positionReportsListCmd.Flags().Bool("include-ultra-fields", false, "Include custom ultra fields")

	addPositionReportWriteFlags(positionReportsCreateCmd)

	positionReportsCmd.AddCommand(positionReportsListCmd)
	positionReportsCmd.AddCommand(positionReportsCreateCmd)
}

func runPositionReportsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("vessel-id"); v != "" {
		params["vessel_id"] = v
	}
	if v, _ := cmd.Flags().GetBool("include-ultra-fields"); v {
		params["include_ultra_fields"] = "true"
	}

	response, err := c.Get("position-reports", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "position reports", []string{
		"mobileops vessels list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runPositionReportsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	// Fields are sent at the top level, not inside a wrapper.
	body := bodyFromFlags(cmd, positionReportFlags)

	response, err := c.Post("position-reports", body)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Position report created", []string{
		"mobileops position-reports list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
