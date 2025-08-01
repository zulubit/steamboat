package cmd

import (
	"github.com/spf13/cobra"
)

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate code for models and migrations",
	Long:  `Generate boilerplate code for models, migrations, and other components.`,
}

func init() {
	rootCmd.AddCommand(makeCmd)
}