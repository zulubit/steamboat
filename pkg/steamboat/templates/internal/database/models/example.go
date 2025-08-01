package models

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Example struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ExampleQueries struct {
	db *sqlx.DB
}

func NewExampleQueries(db *sqlx.DB) *ExampleQueries {
	return &ExampleQueries{db: db}
}

func (q *ExampleQueries) GetByID(ctx context.Context, id int) (*Example, error) {
	var example Example
	query := `SELECT id, created_at, updated_at FROM examples WHERE id = ?`
	err := q.db.GetContext(ctx, &example, query, id)
	if err != nil {
		return nil, err
	}
	return &example, nil
}

func (q *ExampleQueries) GetAll(ctx context.Context) ([]Example, error) {
	var examples []Example
	query := `SELECT id, created_at, updated_at FROM examples ORDER BY created_at DESC`
	err := q.db.SelectContext(ctx, &examples, query)
	if err != nil {
		return nil, err
	}
	return examples, nil
}

func (q *ExampleQueries) Create(ctx context.Context, example *Example) error {
	query := `
		INSERT INTO examples (created_at, updated_at) 
		VALUES (:created_at, :updated_at)
	`
	result, err := q.db.NamedExecContext(ctx, query, example)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	example.ID = int(id)
	return nil
}

func (q *ExampleQueries) Update(ctx context.Context, example *Example) error {
	query := `
		UPDATE examples 
		SET updated_at = :updated_at
		WHERE id = :id
	`
	_, err := q.db.NamedExecContext(ctx, query, example)
	return err
}

func (q *ExampleQueries) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM examples WHERE id = ?`
	_, err := q.db.ExecContext(ctx, query, id)
	return err
}