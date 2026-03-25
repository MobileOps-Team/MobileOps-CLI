package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

const cliVersion = "0.1.0"

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
