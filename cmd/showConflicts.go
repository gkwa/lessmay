package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
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

		var paths []string
		if len(args) == 0 {
			paths = append(paths, defaultObsidianPath)
		} else {
			for _, arg := range args {
				expandedPath, err := homedir.Expand(arg)
				if err != nil {
					logger.Error(err, "Failed to expand path", "path", arg)
					os.Exit(1)
				}
				paths = append(paths, expandedPath)
			}
		}

		resolver := core.NewSyncConflictResolver(logger)
		if err := resolver.ResolveSyncConflicts(paths, skipPaths); err != nil {
			logger.Error(err, "Failed to resolve sync conflicts")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(showConflictsCmd)

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	defaultObsidianPath = filepath.Join(home, "Documents", "Obsidian Vault")

	showConflictsCmd.Flags().StringVarP(&defaultObsidianPath, "default-path", "d", defaultObsidianPath, "Default Obsidian vault path")
}
