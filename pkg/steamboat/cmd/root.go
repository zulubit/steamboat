package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	_ "github.com/joho/godotenv/autoload"
)

var rootCmd = &cobra.Command{
	Use:   "steamboat",
	Short: "Steamboat CLI - Manage your Steamboat application",
	Long: `Steamboat CLI is a command line tool for managing your Steamboat application.
	
It provides commands for database migrations, server management, and more.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
}