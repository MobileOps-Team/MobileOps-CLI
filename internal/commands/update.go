package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/MobileOps-Team/mobileops-cli/internal/updater"
	"github.com/spf13/cobra"
)

const skillUpdateTimeout = 90 * time.Second

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
	runSkillUpdate(false)
}

func updateSkillQuiet() {
	runSkillUpdate(true)
}

func runSkillUpdate(quiet bool) {
	npx, err := exec.LookPath("npx")
	if err != nil {
		if !quiet {
			fmt.Fprintln(os.Stderr, "\nSkill update skipped: npx not found. Update manually with:")
			fmt.Fprintln(os.Stderr, "  npx skills update MobileOps-Team/mobileops-cli")
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), skillUpdateTimeout)
	defer cancel()

	update := exec.CommandContext(ctx, npx, "-y", "skills", "update", "MobileOps-Team/mobileops-cli")
	update.Stdin = nil
	if quiet {
		update.Stdout = io.Discard
		update.Stderr = io.Discard
	} else {
		fmt.Println("\nUpdating skill...")
		update.Stdout = os.Stdout
		update.Stderr = os.Stderr
	}

	if err := update.Run(); err != nil {
		if !quiet {
			fmt.Fprintln(os.Stderr, "Skill update failed:", err)
			fmt.Fprintln(os.Stderr, "  Update manually: npx skills update MobileOps-Team/mobileops-cli")
		}
		return
	}
	if !quiet {
		fmt.Println("\u2713 Skill updated")
	}
}
