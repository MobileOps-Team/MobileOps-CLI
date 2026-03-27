package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/MobileOps-Team/mobileops-cli/internal/updater"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update mobileops-cli and skills to the latest version",
	Long:  "Check GitHub Releases for the latest version, self-update the CLI binary, and update the installed skill.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := updater.SelfUpdate(cliVersion); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		updateSkill()
	},
}

func updateSkill() {
	npx, err := exec.LookPath("npx")
	if err != nil {
		fmt.Fprintln(os.Stderr, "\nSkill update skipped: npx not found. Update manually with:")
		fmt.Fprintln(os.Stderr, "  npx skills update MobileOps-Team/mobileops-cli")
		return
	}

	fmt.Println("\nUpdating skill...")
	update := exec.Command(npx, "skills", "update", "MobileOps-Team/mobileops-cli")
	update.Stdout = os.Stdout
	update.Stderr = os.Stderr
	if err := update.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Skill update failed:", err)
		fmt.Fprintln(os.Stderr, "  Update manually: npx skills update MobileOps-Team/mobileops-cli")
		return
	}
	fmt.Println("✓ Skill updated")
}
