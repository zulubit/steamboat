package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// These variables are set during build time
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information for Steamboat CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Steamboat CLI\n")
		fmt.Printf("  Version:    %s\n", Version)
		fmt.Printf("  Git Commit: %s\n", GitCommit)
		fmt.Printf("  Build Time: %s\n", BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}