package database


import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"

	"<<!.ProjectName!>>/internal/database/models"
	"<<!.ProjectName!>>/internal/utils"
)

// Service represents a service that interacts with a database.
type Service interface {
	Close() error
	Example() *models.ExampleQueries
}

type service struct {
	db      *sqlx.DB
	example *models.ExampleQueries
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
		utils.Logger.Error("Failed to open database", "error", err)
		panic(err)
	}

	dbInstance = &service{
		db:      db,
		example: models.NewExampleQueries(db),
	}
	return dbInstance
}


func (s *service) Example() *models.ExampleQueries {
	return s.example
}

func (s *service) Close() error {
	utils.Logger.Info("Disconnected from database", "url", dburl)
	return s.db.Close()
}