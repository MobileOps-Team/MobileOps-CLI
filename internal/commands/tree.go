package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Print tree of all available commands",
	Run: func(cmd *cobra.Command, args []string) {
		printCommandTree(rootCmd, "")
	},
}

func printCommandTree(cmd *cobra.Command, indent string) {
	if cmd.Name() == "help" || cmd.Name() == "completion" {
		return
	}

	line := cmd.Name()
	if cmd.Short != "" {
		line += " - " + cmd.Short
	}
	fmt.Println(indent + line)

	children := cmd.Commands()
	for i, child := range children {
		if child.Name() == "help" || child.Name() == "completion" {
			continue
		}
		prefix := "├── "
		childIndent := "│   "
		if isLastVisible(children, i) {
			prefix = "└── "
			childIndent = "    "
		}
		printSubTree(child, indent+prefix, indent+childIndent)
	}
}

func printSubTree(cmd *cobra.Command, prefix string, indent string) {
	line := cmd.Name()
	if cmd.Short != "" {
		line += " - " + cmd.Short
	}
	// Replace indent portion with prefix for first line
	fmt.Println(strings.TrimRight(prefix, " ") + " " + line)

	children := cmd.Commands()
	for i, child := range children {
		if child.Name() == "help" || child.Name() == "completion" {
			continue
		}
		childPrefix := "├── "
		childIndent := "│   "
		if isLastVisible(children, i) {
			childPrefix = "└── "
			childIndent = "    "
		}
		printSubTree(child, indent+childPrefix, indent+childIndent)
	}
}

func isLastVisible(cmds []*cobra.Command, idx int) bool {
	for i := idx + 1; i < len(cmds); i++ {
		if cmds[i].Name() != "help" && cmds[i].Name() != "completion" {
			return false
		}
	}
	return true
}
