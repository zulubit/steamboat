package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	_ "github.com/joho/godotenv/autoload"
)

var (
	port string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Steamboat server",
	Long:  `Start the Steamboat API server with the configured settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Override port if provided
		if port != "" {
			os.Setenv("PORT", port)
		}

		// Run the API server using exec
		apiCmd := exec.Command("go", "run", "cmd/api/main.go")
		apiCmd.Stdout = os.Stdout
		apiCmd.Stderr = os.Stderr
		apiCmd.Stdin = os.Stdin
		
		log.Println("Starting Steamboat server...")
		if err := apiCmd.Run(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	
	// Add flags
	serveCmd.Flags().StringVarP(&port, "port", "p", "", "Port to run the server on (overrides PORT env var)")
}