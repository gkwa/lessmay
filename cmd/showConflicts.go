package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gkwa/lessmay/core"
)

var defaultObsidianPath string

var showConflictsCmd = &cobra.Command{
	Use:     "show-conflicts [directories...]",
	Short:   "Resolve sync conflicts in Obsidian vault",
	Long:    `This command finds and displays differences between sync conflict files and their original versions in an Obsidian vault.`,
	Aliases: []string{"show-conflicts"},
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.Info("Running showConflicts command")

		if err := core.ShowConflicts(logger, args, defaultObsidianPath, skipPaths); err != nil {
			logger.Error(err, "Failed to resolve sync conflicts")
			cmd.PrintErrln("Error:", err)
			cmd.PrintErrln("Run with --verbose for more details.")
		}
	},
}

func init() {
	rootCmd.AddCommand(showConflictsCmd)

	defaultObsidianPath = core.GetDefaultObsidianPath()

	showConflictsCmd.Flags().
		StringVarP(&defaultObsidianPath, "default-path", "d", defaultObsidianPath, "Default Obsidian vault path")
}
