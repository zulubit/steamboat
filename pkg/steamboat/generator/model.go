package generator

import (
	"fmt"
	"os"
	"path/filepath"
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
	
	// Add getter method before the closing of the file
	getterMethod := fmt.Sprintf(`
// %s returns the %sQueries instance
func (s *service) %s() *models.%sQueries {
	return s.%s
}`, structName, structName, structName, structName, varName)
	
	// Find a good place to insert the getter - after the last getter method or before Close()
	closeIndex := strings.Index(fileContent, "// Close closes the database connection.")
	if closeIndex > 0 {
		fileContent = fileContent[:closeIndex] + getterMethod + "\n\n" + fileContent[closeIndex:]
	} else {
		// Fallback: add at the end
		fileContent = strings.TrimRight(fileContent, "\n") + "\n" + getterMethod + "\n"
	}
	
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