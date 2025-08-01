package handlers

import (
	"<<!.ProjectName!>>/internal/database"
)

type Handlers struct {
	db database.Service
}

func New(db database.Service) *Handlers {
	return &Handlers{
		db: db,
	}
}