package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var vesselsCmd = &cobra.Command{
	Use:   "vessels",
	Short: "Manage vessels/assets",
}

var vesselsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all vessels/assets",
	RunE:  runVesselsList,
}

var vesselsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single vessel by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselsGet,
}

var vesselsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vessel",
	RunE:  runVesselsCreate,
}

var vesselsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a vessel",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselsUpdate,
}

var vesselsSpecsCmd = &cobra.Command{
	Use:   "specs VESSEL_ID",
	Short: "List specs for a vessel",
	Args:  cobra.ExactArgs(1),
	RunE:  runVesselsSpecs,
}

var vesselFlags = map[string]string{
	"name":          "name",
	"division-id":   "division_id",
	"status":        "status",
	"category":      "category",
	"vessel-type-id": "vessel_type_id",
	"customer-id":   "customer_id",
	"customer-name": "customer_name",
	"imo-number":    "imo_number",
	"uscg-number":   "uscg_number",
	"mmsi-number":   "mmsi_number",
	"call-sign":     "call_sign",
	"color":         "color",
	"length":        "length",
	"height":        "height",
	"width":         "width",
	"dimension-unit": "dimension_unit",
}

var vesselBoolFlags = map[string]string{
	"active":         "active",
	"voyage-enabled": "voyage_enabled",
	"external":       "external",
}

func addVesselWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Vessel name")
	cmd.Flags().String("division-id", "", "Division ID")
	cmd.Flags().String("status", "", "Status")
	cmd.Flags().String("category", "", "Category")
	cmd.Flags().String("vessel-type-id", "", "Vessel type ID")
	cmd.Flags().String("customer-id", "", "Customer ID")
	cmd.Flags().String("customer-name", "", "Customer name")
	cmd.Flags().String("imo-number", "", "IMO number")
	cmd.Flags().String("uscg-number", "", "USCG number")
	cmd.Flags().String("mmsi-number", "", "MMSI number")
	cmd.Flags().String("call-sign", "", "Call sign")
	cmd.Flags().String("color", "", "Color")
	cmd.Flags().String("length", "", "Length")
	cmd.Flags().String("height", "", "Height")
	cmd.Flags().String("width", "", "Width")
	cmd.Flags().String("dimension-unit", "", "Dimension unit")
	cmd.Flags().Bool("active", false, "Active")
	cmd.Flags().Bool("voyage-enabled", false, "Voyage enabled")
	cmd.Flags().Bool("external", false, "External vessel")
}

func init() {
	vesselsListCmd.Flags().Int("page", 1, "Page number")
	vesselsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	vesselsListCmd.Flags().String("division-id", "", "Filter by division ID")

	addVesselWriteFlags(vesselsCreateCmd)
	addVesselWriteFlags(vesselsUpdateCmd)

	vesselsCmd.AddCommand(vesselsListCmd)
	vesselsCmd.AddCommand(vesselsGetCmd)
	vesselsCmd.AddCommand(vesselsCreateCmd)
	vesselsCmd.AddCommand(vesselsUpdateCmd)
	vesselsCmd.AddCommand(vesselsSpecsCmd)
}

func runVesselsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("assets", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops vessels get %s", 3)
	env := envelope.WrapCollection(response, "vessels", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("assets/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	vesselID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Vessel", []string{
		fmt.Sprintf("mobileops vessels specs %s", vesselID),
		fmt.Sprintf("mobileops components list --vessel-id %s", vesselID),
		fmt.Sprintf("mobileops jobs list --vessel-id %s", vesselID),
		fmt.Sprintf("mobileops work-requests list --vessel-id %s", vesselID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vesselFlags)
	bodyFromBoolFlags(cmd, body, vesselBoolFlags)

	response, err := c.Post("assets", wrapBody("asset", body))
	if err != nil {
		return handleClientError(err)
	}

	vesselID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Vessel created", []string{
		fmt.Sprintf("mobileops vessels get %s", vesselID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, vesselFlags)
	bodyFromBoolFlags(cmd, body, vesselBoolFlags)

	response, err := c.Put(fmt.Sprintf("assets/%s", id), wrapBody("asset", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Vessel updated", []string{
		fmt.Sprintf("mobileops vessels get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runVesselsSpecs(cmd *cobra.Command, args []string) error {
	vesselID := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get("vessel-specs", map[string]string{"vessel_id": vesselID})
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "vessel specs", []string{
		fmt.Sprintf("mobileops vessels get %s", vesselID),
		fmt.Sprintf("mobileops components list --vessel-id %s", vesselID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
