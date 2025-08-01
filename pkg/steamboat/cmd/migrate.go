package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zulubit/steamboat/pkg/steamboat/migrate"
)

var (
	rollback   bool
	showStatus bool
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run all pending database migrations, rollback the last migration, or show status.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showStatus {
			version, dirty, err := migrate.Status()
			if err != nil {
				log.Fatalf("Failed to get migration status: %v", err)
			}
			
			if version == 0 {
				log.Println("No migrations have been applied")
			} else {
				status := "clean"
				if dirty {
					status = "dirty"
				}
				log.Printf("Current migration version: %d (%s)", version, status)
			}
			return
		}
		
		if rollback {
			log.Println("Rolling back last migration...")
			if err := migrate.Rollback(); err != nil {
				log.Fatalf("Rollback failed: %v", err)
			}
			log.Println("Rollback completed successfully")
		} else {
			log.Println("Running migrations...")
			if err := migrate.Run(); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			log.Println("All migrations completed successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	
	// Add flags
	migrateCmd.Flags().BoolVarP(&rollback, "rollback", "r", false, "Rollback the last migration")
	migrateCmd.Flags().BoolVarP(&showStatus, "status", "s", false, "Show current migration status")
}