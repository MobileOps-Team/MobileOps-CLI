package commands

import (
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var locationsCmd = &cobra.Command{
	Use:   "locations",
	Short: "Manage locations",
}

var locationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List locations",
	RunE:  runLocationsList,
}

var locationsGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a single location by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runLocationsGet,
}

var locationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new location",
	RunE:  runLocationsCreate,
}

var locationsUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a location",
	Args:  cobra.ExactArgs(1),
	RunE:  runLocationsUpdate,
}

var locationFlags = map[string]string{
	"name":      "name",
	"type":      "type",
	"latitude":  "latitude",
	"longitude": "longitude",
}

func addLocationWriteFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Location name")
	cmd.Flags().String("type", "", "Location type")
	cmd.Flags().String("latitude", "", "Latitude")
	cmd.Flags().String("longitude", "", "Longitude")
}

func init() {
	locationsListCmd.Flags().Int("page", 1, "Page number")
	locationsListCmd.Flags().Int("limit", 10, "Items per page (max 100)")

	addLocationWriteFlags(locationsCreateCmd)
	addLocationWriteFlags(locationsUpdateCmd)

	locationsCmd.AddCommand(locationsListCmd)
	locationsCmd.AddCommand(locationsGetCmd)
	locationsCmd.AddCommand(locationsCreateCmd)
	locationsCmd.AddCommand(locationsUpdateCmd)
}

func runLocationsList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	response, err := c.Get("locations", params)
	if err != nil {
		return handleClientError(err)
	}

	breadcrumbs := collectionBreadcrumbs(response, "mobileops locations get %s", 3)
	env := envelope.WrapCollection(response, "locations", breadcrumbs)
	formatter.Output(env, jsonOutput)
	return nil
}

func runLocationsGet(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	response, err := c.Get(fmt.Sprintf("locations/%s", id), nil)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Location", []string{
		"mobileops locations list",
		"mobileops jobs list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runLocationsCreate(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, locationFlags)

	response, err := c.Post("locations", wrapBody("location", body))
	if err != nil {
		return handleClientError(err)
	}

	locationID := envelope.ExtractID(response)
	env := envelope.WrapRecord(response, "Location created", []string{
		fmt.Sprintf("mobileops locations get %s", locationID),
	})
	formatter.Output(env, jsonOutput)
	return nil
}

func runLocationsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	body := bodyFromFlags(cmd, locationFlags)

	response, err := c.Put(fmt.Sprintf("locations/%s", id), wrapBody("location", body))
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapRecord(response, "Location updated", []string{
		fmt.Sprintf("mobileops locations get %s", id),
	})
	formatter.Output(env, jsonOutput)
	return nil
}
