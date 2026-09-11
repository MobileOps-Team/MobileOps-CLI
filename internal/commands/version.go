package commands

import (
	"encoding/json"
	"fmt"

	"github.com/MobileOps-Team/mobileops-cli/internal/version"
	"github.com/spf13/cobra"
)

const cliVersion = version.Version

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		if jsonOutput {
			data, _ := json.MarshalIndent(map[string]string{"version": cliVersion}, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Printf("mobileops-cli %s\n", cliVersion)
		}
	},
}
