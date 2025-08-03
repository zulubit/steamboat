package updatetemplates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	workingCopyName = "workingcopy"
	templateMarker  = "<<!.ProjectName!>>"
)

// Run executes the template update process
func Run() error {
	// Verify workingcopy exists
	if _, err := os.Stat(workingCopyName); err != nil {
		return fmt.Errorf("workingcopy directory not found: %w", err)
	}

	// Remove existing templates
	templatesDir := "pkg/steamboat/templates"
	if err := os.RemoveAll(templatesDir); err != nil {
		return fmt.Errorf("failed to remove templates directory: %w", err)
	}

	// Create fresh templates directory
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	// Copy and process files from workingcopy
	return processWorkingCopy(workingCopyName, templatesDir)
}

func processWorkingCopy(srcDir, dstDir string) error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}


		// Calculate relative path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		// Skip go.sum files
		if d.Name() == "go.sum" {
			return nil
		}

		// Skip _templ.go files in internal/views directory
		if strings.Contains(relPath, "internal/views") && strings.HasSuffix(d.Name(), "_templ.go") {
			return nil
		}

		targetPath := filepath.Join(dstDir, relPath)

		// Rename go.mod to go.mod.tpl
		if d.Name() == "go.mod" && !d.IsDir() {
			targetPath = targetPath + ".tpl"
		}

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Process and copy file
		return processFile(path, targetPath)
	})
}

func processFile(srcPath, dstPath string) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", srcPath, err)
	}

	// Replace workingcopy references with template marker
	processed := strings.ReplaceAll(string(content), workingCopyName, templateMarker)

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