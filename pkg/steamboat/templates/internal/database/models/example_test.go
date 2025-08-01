package models

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	
	schema := `
	CREATE TABLE examples (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`
	
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create test schema: %v", err)
	}
	
	return db
}

func TestExampleQueries_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	q := NewExampleQueries(db)
	ctx := context.Background()
	
	example := &Example{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err := q.Create(ctx, example)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}
	
	if example.ID == 0 {
		t.Error("Expected ID to be set after create")
	}
}

func TestExampleQueries_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	q := NewExampleQueries(db)
	ctx := context.Background()
	
	now := time.Now()
	example := &Example{
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	err := q.Create(ctx, example)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	retrieved, err := q.GetByID(ctx, example.ID)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}
	
	if retrieved.ID != example.ID {
		t.Errorf("Expected ID %d, got %d", example.ID, retrieved.ID)
	}
}

func TestExampleQueries_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	q := NewExampleQueries(db)
	ctx := context.Background()
	
	now := time.Now()
	example := &Example{
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	err := q.Create(ctx, example)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	newTime := time.Now().Add(time.Hour)
	example.UpdatedAt = newTime
	
	err = q.Update(ctx, example)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}
	
	retrieved, err := q.GetByID(ctx, example.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	
	if retrieved.UpdatedAt.Unix() != newTime.Unix() {
		t.Errorf("UpdatedAt not updated correctly")
	}
}

func TestExampleQueries_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	q := NewExampleQueries(db)
	ctx := context.Background()
	
	example := &Example{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err := q.Create(ctx, example)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	err = q.Delete(ctx, example.ID)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
	
	_, err = q.GetByID(ctx, example.ID)
	if err == nil {
		t.Error("Expected error when getting deleted record")
	}
}