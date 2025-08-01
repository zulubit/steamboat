package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const modelTemplate = `package models

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type {{.StructName}} struct {
	ID        int       ` + "`" + `db:"id" json:"id"` + "`" + `
	CreatedAt time.Time ` + "`" + `db:"created_at" json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `db:"updated_at" json:"updated_at"` + "`" + `
}

type {{.StructName}}Queries struct {
	db *sqlx.DB
}

func New{{.StructName}}Queries(db *sqlx.DB) *{{.StructName}}Queries {
	return &{{.StructName}}Queries{db: db}
}

func (q *{{.StructName}}Queries) GetByID(ctx context.Context, id int) (*{{.StructName}}, error) {
	var {{.VarName}} {{.StructName}}
	query := ` + "`" + `SELECT id, created_at, updated_at FROM {{.TableName}} WHERE id = ?` + "`" + `
	err := q.db.GetContext(ctx, &{{.VarName}}, query, id)
	if err != nil {
		return nil, err
	}
	return &{{.VarName}}, nil
}

func (q *{{.StructName}}Queries) GetAll(ctx context.Context) ([]{{.StructName}}, error) {
	var {{.PluralVarName}} []{{.StructName}}
	query := ` + "`" + `SELECT id, created_at, updated_at FROM {{.TableName}} ORDER BY created_at DESC` + "`" + `
	err := q.db.SelectContext(ctx, &{{.PluralVarName}}, query)
	if err != nil {
		return nil, err
	}
	return {{.PluralVarName}}, nil
}

func (q *{{.StructName}}Queries) Create(ctx context.Context, {{.VarName}} *{{.StructName}}) error {
	query := ` + "`" + `
		INSERT INTO {{.TableName}} (created_at, updated_at) 
		VALUES (:created_at, :updated_at)
	` + "`" + `
	result, err := q.db.NamedExecContext(ctx, query, {{.VarName}})
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	{{.VarName}}.ID = int(id)
	return nil
}

func (q *{{.StructName}}Queries) Update(ctx context.Context, {{.VarName}} *{{.StructName}}) error {
	query := ` + "`" + `
		UPDATE {{.TableName}} 
		SET updated_at = :updated_at
		WHERE id = :id
	` + "`" + `
	_, err := q.db.NamedExecContext(ctx, query, {{.VarName}})
	return err
}

func (q *{{.StructName}}Queries) Delete(ctx context.Context, id int) error {
	query := ` + "`" + `DELETE FROM {{.TableName}} WHERE id = ?` + "`" + `
	_, err := q.db.ExecContext(ctx, query, id)
	return err
}
`

const modelTestTemplate = `package models

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setup{{.StructName}}Test(t *testing.T) (*{{.StructName}}Queries, func()) {
	// Create in-memory SQLite database
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create table
	_, err = db.Exec(` + "`" + `CREATE TABLE {{.TableName}} (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)` + "`" + `)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	queries := New{{.StructName}}Queries(db)
	
	cleanup := func() {
		db.Close()
	}

	return queries, cleanup
}

func Test{{.StructName}}Queries_Create(t *testing.T) {
	queries, cleanup := setup{{.StructName}}Test(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()
	
	{{.VarName}} := &{{.StructName}}{
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := queries.Create(ctx, {{.VarName}})
	if err != nil {
		t.Fatalf("Failed to create {{.VarName}}: %v", err)
	}

	if {{.VarName}}.ID == 0 {
		t.Error("Expected ID to be set after creation")
	}
}

func Test{{.StructName}}Queries_GetByID(t *testing.T) {
	queries, cleanup := setup{{.StructName}}Test(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()
	
	// Create a {{.VarName}} first
	{{.VarName}} := &{{.StructName}}{
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := queries.Create(ctx, {{.VarName}})
	if err != nil {
		t.Fatalf("Failed to create {{.VarName}}: %v", err)
	}

	// Get the {{.VarName}} by ID
	retrieved, err := queries.GetByID(ctx, {{.VarName}}.ID)
	if err != nil {
		t.Fatalf("Failed to get {{.VarName}} by ID: %v", err)
	}

	if retrieved.ID != {{.VarName}}.ID {
		t.Errorf("Expected ID %d, got %d", {{.VarName}}.ID, retrieved.ID)
	}
}

func Test{{.StructName}}Queries_GetAll(t *testing.T) {
	queries, cleanup := setup{{.StructName}}Test(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	// Create multiple {{.PluralVarName}}
	for i := 0; i < 3; i++ {
		{{.VarName}} := &{{.StructName}}{
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := queries.Create(ctx, {{.VarName}})
		if err != nil {
			t.Fatalf("Failed to create {{.VarName}} %d: %v", i, err)
		}
	}

	// Get all {{.PluralVarName}}
	{{.PluralVarName}}, err := queries.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get all {{.PluralVarName}}: %v", err)
	}

	if len({{.PluralVarName}}) != 3 {
		t.Errorf("Expected 3 {{.PluralVarName}}, got %d", len({{.PluralVarName}}))
	}
}

func Test{{.StructName}}Queries_Update(t *testing.T) {
	queries, cleanup := setup{{.StructName}}Test(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()
	
	// Create a {{.VarName}} first
	{{.VarName}} := &{{.StructName}}{
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := queries.Create(ctx, {{.VarName}})
	if err != nil {
		t.Fatalf("Failed to create {{.VarName}}: %v", err)
	}

	// Update the {{.VarName}}
	{{.VarName}}.UpdatedAt = time.Now().Add(time.Hour)
	err = queries.Update(ctx, {{.VarName}})
	if err != nil {
		t.Fatalf("Failed to update {{.VarName}}: %v", err)
	}

	// Verify the update
	retrieved, err := queries.GetByID(ctx, {{.VarName}}.ID)
	if err != nil {
		t.Fatalf("Failed to get updated {{.VarName}}: %v", err)
	}

	if retrieved.UpdatedAt.Equal(now) {
		t.Error("Expected UpdatedAt to be changed after update")
	}
}

func Test{{.StructName}}Queries_Delete(t *testing.T) {
	queries, cleanup := setup{{.StructName}}Test(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()
	
	// Create a {{.VarName}} first
	{{.VarName}} := &{{.StructName}}{
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := queries.Create(ctx, {{.VarName}})
	if err != nil {
		t.Fatalf("Failed to create {{.VarName}}: %v", err)
	}

	// Delete the {{.VarName}}
	err = queries.Delete(ctx, {{.VarName}}.ID)
	if err != nil {
		t.Fatalf("Failed to delete {{.VarName}}: %v", err)
	}

	// Verify deletion
	_, err = queries.GetByID(ctx, {{.VarName}}.ID)
	if err == nil {
		t.Error("Expected error when getting deleted {{.VarName}}")
	}
}
`

type ModelData struct {
	StructName    string
	VarName       string
	PluralVarName string
	TableName     string
}

func GenerateModel(name string) error {
	// Convert name to proper case formats
	structName := toPascalCase(name)
	varName := toCamelCase(name)
	pluralVarName := pluralize(varName)
	tableName := toSnakeCase(pluralize(name))

	data := ModelData{
		StructName:    structName,
		VarName:       varName,
		PluralVarName: pluralVarName,
		TableName:     tableName,
	}

	// Create the model file
	modelPath := filepath.Join("internal", "database", "models", fmt.Sprintf("%s.go", name))
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(modelPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Execute template
	tmpl, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create the test file
	testPath := filepath.Join("internal", "database", "models", fmt.Sprintf("%s_test.go", name))
	testFile, err := os.Create(testPath)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	defer testFile.Close()

	// Execute test template
	testTmpl, err := template.New("modelTest").Parse(modelTestTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse test template: %w", err)
	}

	if err := testTmpl.Execute(testFile, data); err != nil {
		return fmt.Errorf("failed to execute test template: %w", err)
	}

	// Update database.go to include the new model
	if err := updateDatabaseFile(structName); err != nil {
		return fmt.Errorf("failed to update database.go: %w", err)
	}

	return nil
}

func updateDatabaseFile(structName string) error {
	dbPath := filepath.Join("internal", "database", "database.go")
	
	// Read the existing file
	content, err := os.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("failed to read database.go: %w", err)
	}
	
	fileContent := string(content)
	varName := toCamelCase(structName)
	
	// Check if already exists
	if strings.Contains(fileContent, fmt.Sprintf("%s() *models.%sQueries", structName, structName)) {
		fmt.Printf("✓ Model %s already exists in database.go\n", structName)
		return nil
	}
	
	// Check if models import exists, add it if it doesn't
	if !strings.Contains(fileContent, "/internal/database/models") {
		// Find the utils import line and add models import before it
		// Extract the project name from the existing utils import
		re := regexp.MustCompile(`"([^"]+)/internal/utils"`)
		matches := re.FindStringSubmatch(fileContent)
		if len(matches) > 1 {
			projectName := matches[1]
			utilsImport := fmt.Sprintf(`"%s/internal/utils"`, projectName)
			modelsImport := fmt.Sprintf(`"%s/internal/database/models"
	%s`, projectName, utilsImport)
			fileContent = strings.Replace(fileContent, utilsImport, modelsImport, 1)
		}
	}
	
	// Add to Service interface using markers
	queriesStart := "// STEAMBOAT:QUERIES_START - Auto-generated query methods"
	queriesEnd := "// STEAMBOAT:QUERIES_END"
	interfaceAddition := fmt.Sprintf("\t%s() *models.%sQueries", structName, structName)
	fileContent = addBetweenMarkers(fileContent, queriesStart, queriesEnd, interfaceAddition)
	
	// Add to service struct using markers
	fieldsStart := "// STEAMBOAT:FIELDS_START - Auto-generated query fields"
	fieldsEnd := "// STEAMBOAT:FIELDS_END"
	structAddition := fmt.Sprintf("\t%s *models.%sQueries", varName, structName)
	fileContent = addBetweenMarkers(fileContent, fieldsStart, fieldsEnd, structAddition)
	
	// Add to New() function initialization using markers
	initStart := "// STEAMBOAT:INIT_START - Auto-generated query initialization"
	initEnd := "// STEAMBOAT:INIT_END"
	initAddition := fmt.Sprintf("\t\t%s: models.New%sQueries(db),", varName, structName)
	fileContent = addBetweenMarkers(fileContent, initStart, initEnd, initAddition)
	
	// Add getter method using markers
	gettersStart := "// STEAMBOAT:GETTERS_START - Auto-generated getter methods"
	gettersEnd := "// STEAMBOAT:GETTERS_END"
	getterAddition := fmt.Sprintf(`// %s returns the %sQueries instance
func (s *service) %s() *models.%sQueries {
	return s.%s
}`, structName, structName, structName, structName, varName)
	fileContent = addBetweenMarkers(fileContent, gettersStart, gettersEnd, getterAddition)
	
	// Write the updated content back
	if err := os.WriteFile(dbPath, []byte(fileContent), 0644); err != nil {
		return fmt.Errorf("failed to write database.go: %w", err)
	}
	
	fmt.Printf("✓ Updated internal/database/database.go\n")
	return nil
}

func addBetweenMarkers(content, startMarker, endMarker, addition string) string {
	startIndex := strings.Index(content, startMarker)
	endIndex := strings.Index(content, endMarker)
	
	if startIndex == -1 || endIndex == -1 {
		// Fallback to old method if markers not found
		return content
	}
	
	// Find the end of the start marker line
	lineEnd := strings.Index(content[startIndex:], "\n")
	if lineEnd == -1 {
		return content
	}
	
	insertPos := startIndex + lineEnd
	return content[:insertPos] + "\n" + addition + content[insertPos:]
}


// Helper functions for string conversions
func toPascalCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) == 0 {
		return ""
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func pluralize(s string) string {
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "ch") {
		return s + "es"
	}
	return s + "s"
}