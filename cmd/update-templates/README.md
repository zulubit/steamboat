# Template Update Workflow

This tool helps maintain the Steamboat templates by allowing you to work on a real project and sync changes back.

## Workflow

1. **Generate a working copy:**
   ```bash
   go run ./cmd/steamboat/main.go create workingcopy
   ```

2. **Make your changes:**
   Work on the `workingcopy` project as a normal Go project. Test your changes, add features, etc.

3. **Update templates:**
   ```bash
   go run cmd/update-templates/main.go
   ```

   This will:
   - Remove the current templates in `pkg/steamboat/templates`
   - Copy all files from `workingcopy` to the templates directory
   - Replace all occurrences of "workingcopy" with `<<!.ProjectName!>>`

4. **Clean up:**
   ```bash
   rm -rf workingcopy
   ```

## Notes

- All "workingcopy" strings in file contents are replaced with the template placeholder
- `go.sum` files are skipped (will be generated fresh for each project)
- `go.mod` files are renamed to `go.mod.tpl` (the .tpl extension is removed by the generator)
- Make sure to test the template generation after updating
