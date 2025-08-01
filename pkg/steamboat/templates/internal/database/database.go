package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"

	"<<!.ProjectName!>>/internal/utils"
)

// Service represents a service that interacts with a database.
type Service interface {
	Close() error
	// STEAMBOAT:QUERIES_START - Auto-generated query methods
	// STEAMBOAT:QUERIES_END
}

type service struct {
	db *sqlx.DB
	// STEAMBOAT:FIELDS_START - Auto-generated query fields
	// STEAMBOAT:FIELDS_END
}

var (
	dburl      = os.Getenv("DB_URL")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sqlx.Open("sqlite3", dburl)
	if err != nil {
		if utils.Logger != nil {
			utils.Logger.Error("Failed to open database", "error", err)
		}
		panic(err)
	}

	dbInstance = &service{
		db: db,
		// STEAMBOAT:INIT_START - Auto-generated query initialization
		// STEAMBOAT:INIT_END
	}
	return dbInstance
}

// STEAMBOAT:GETTERS_START - Auto-generated getter methods
// STEAMBOAT:GETTERS_END

func (s *service) Close() error {
	if utils.Logger != nil {
		utils.Logger.Info("Disconnected from database", "url", dburl)
	}
	return s.db.Close()
}

