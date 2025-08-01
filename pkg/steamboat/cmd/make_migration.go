package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/zulubit/steamboat/pkg/steamboat/generator"
)

var makeMigrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Generate a new migration",
	Long:  `Generate a new migration file with up and down functions.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		
		log.Printf("Creating migration: %s", migrationName)
		
		filename, err := generator.GenerateMigration(migrationName)
		if err != nil {
			log.Fatalf("Failed to generate migration: %v", err)
		}
		
		fmt.Printf("✓ Migration created successfully\n")
		fmt.Printf("  Created: %s\n", filename)
	},
}

func init() {
	makeCmd.AddCommand(makeMigrationCmd)
}