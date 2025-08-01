package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zulubit/steamboat/pkg/steamboat/generator"
)

var (
	force bool
)

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new Steamboat project",
	Long:  `Bootstrap a new Steamboat project with all necessary files and structure.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		
		// Use project name as directory name by default
		targetDir := projectName
		
		// Convert to absolute path
		absPath, err := filepath.Abs(targetDir)
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}
		
		log.Printf("Creating Steamboat project: %s", projectName)
		log.Printf("Target directory: %s", absPath)
		
		if err := generator.CreateProject(projectName, absPath, force); err != nil {
			log.Fatalf("Failed to create project: %v", err)
		}
		
		fmt.Printf("✅ Project '%s' created successfully!\n\n", projectName)
		fmt.Printf("Next steps:\n")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Printf("  go run cmd/cli/main.go migrate\n")
		fmt.Printf("  go run cmd/cli/main.go serve\n\n")
		fmt.Printf("Available commands:\n")
		fmt.Printf("  go run cmd/cli/main.go make model [name]     # Generate a model\n")
		fmt.Printf("  go run cmd/cli/main.go make migration [name] # Generate a migration\n")
		fmt.Printf("  go run cmd/cli/main.go migrate              # Run migrations\n")
		fmt.Printf("  go run cmd/cli/main.go serve                # Start the server\n")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	
	// Add flags
	createCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing directory")
}