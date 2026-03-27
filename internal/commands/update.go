package commands

import (
	"fmt"
	"os"

	"github.com/MobileOps-Team/mobileops-cli/internal/updater"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update mobileops-cli to the latest version",
	Long:  "Check GitHub Releases for the latest version and self-update the CLI binary.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := updater.SelfUpdate(cliVersion); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}
