package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TemplateData contains variables for template processing
type TemplateData struct {
	ProjectName string
}

// ProcessTemplate processes a template file and replaces placeholders
func ProcessTemplate(content string, data TemplateData) string {
	content = strings.ReplaceAll(content, "<<!.ProjectName!>>", data.ProjectName)
	return content
}

// CopyTemplateDir copies a directory from templates and processes all files
func CopyTemplateDir(templateDir, targetDir string, data TemplateData) error {
	return filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from template dir
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Skip the template dir itself
		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			// Create directory with proper permissions
			return os.MkdirAll(targetPath, 0755)
		}

		// Handle files
		return copyTemplateFile(path, targetPath, data)
	})
}

func copyTemplateFile(srcPath, dstPath string, data TemplateData) error {
	// Read source file
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", srcPath, err)
	}

	// Process template
	processed := ProcessTemplate(string(content), data)

	// Handle .tmpl files - remove .tmpl extension
	if strings.HasSuffix(dstPath, ".tmpl") {
		dstPath = strings.TrimSuffix(dstPath, ".tmpl")
	}

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write processed content
	if err := os.WriteFile(dstPath, []byte(processed), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", dstPath, err)
	}

	return nil
}