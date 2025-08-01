package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// CreateProject creates a new Steamboat project from templates
func CreateProject(projectName string, targetDir string, force bool) error {
	// Check if target directory exists
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		if !force {
			return fmt.Errorf("directory '%s' already exists. Use --force to overwrite", targetDir)
		}
		// Remove existing directory if force is true
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Get template directory (relative to this file)
	templateDir, err := getTemplateDir()
	if err != nil {
		return fmt.Errorf("failed to find templates: %w", err)
	}

	// Prepare template data
	data := TemplateData{
		ProjectName: projectName,
	}

	// Copy and process templates
	if err := CopyTemplateDir(templateDir, targetDir, data); err != nil {
		return fmt.Errorf("failed to copy templates: %w", err)
	}

	// Initialize go module
	if err := initGoModule(targetDir, projectName); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	return nil
}

func getTemplateDir() (string, error) {
	// Get the directory where this file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// Navigate to templates directory: ../templates from current file
	templateDir := filepath.Join(filepath.Dir(filename), "..", "templates")
	
	// Convert to absolute path
	absPath, err := filepath.Abs(templateDir)
	if err != nil {
		return "", err
	}

	// Check if templates directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("templates directory not found at %s", absPath)
	}

	return absPath, nil
}

func initGoModule(projectDir, projectName string) error {
	// For now, we skip go mod tidy since the CLI package may not be published yet
	// Users should run 'go mod tidy' manually after creating the project
	return nil
}