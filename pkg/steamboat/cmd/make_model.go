package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/zulubit/steamboat/pkg/steamboat/generator"
)

var makeModelCmd = &cobra.Command{
	Use:   "model [name]",
	Short: "Generate a new model",
	Long:  `Generate a new model with database struct and query methods.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modelName := args[0]
		
		log.Printf("Creating model: %s", modelName)
		
		if err := generator.GenerateModel(modelName); err != nil {
			log.Fatalf("Failed to generate model: %v", err)
		}
		
		fmt.Printf("✓ Model '%s' created successfully\n", modelName)
		fmt.Printf("  Created: internal/database/models/%s.go\n", modelName)
		fmt.Printf("  Created: internal/database/models/%s_test.go\n", modelName)
	},
}

func init() {
	makeCmd.AddCommand(makeModelCmd)
}