package hooks

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Generator handles hook script generation
type Generator struct {
	hooksPath   string
	projectPath string
}

// NewGenerator creates a new hook generator
func NewGenerator(hooksPath, projectPath string) *Generator {
	return &Generator{
		hooksPath:   hooksPath,
		projectPath: projectPath,
	}
}

// GenerateHooks generates all hook scripts
func (g *Generator) GenerateHooks() error {
	hooks := []struct {
		name     string
		template string
		desc     string
	}{
		{
			name:     "pre-command.sh",
			template: preCommandTemplate,
			desc:     "pre-command hook",
		},
		{
			name:     "post-command.sh",
			template: postCommandTemplate,
			desc:     "post-command hook",
		},
		{
			name:     "file-modified.sh",
			template: fileModifiedTemplate,
			desc:     "file-modified hook",
		},
		{
			name:     "session-end.sh",
			template: sessionEndTemplate,
			desc:     "session-end hook",
		},
	}

	// Ensure hooks directory exists
	if err := os.MkdirAll(g.hooksPath, 0755); err != nil {
		return fmt.Errorf("creating hooks directory: %w", err)
	}

	// Create common functions file first
	commonPath := filepath.Join(g.hooksPath, "common.sh")
	if err := g.writeFile(commonPath, commonFunctions, nil); err != nil {
		return fmt.Errorf("creating common.sh: %w", err)
	}

	// Generate each hook script
	for _, hook := range hooks {
		hookPath := filepath.Join(g.hooksPath, hook.name)
		if err := g.generateHook(hookPath, hook.template, hook.desc); err != nil {
			return fmt.Errorf("generating %s: %w", hook.desc, err)
		}
	}

	return nil
}

// generateHook generates a single hook script
func (g *Generator) generateHook(path, tmplContent, desc string) error {
	data := struct {
		ProjectPath string
		HooksPath   string
		LogFile     string
	}{
		ProjectPath: g.projectPath,
		HooksPath:   g.hooksPath,
		LogFile:     filepath.Join(g.projectPath, "logs", "hook-errors.log"),
	}

	return g.writeFile(path, tmplContent, data)
}

// writeFile writes content to a file with proper permissions
func (g *Generator) writeFile(path, tmplContent string, data interface{}) error {
	// Parse template
	tmpl, err := template.New("hook").Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	// Create temporary file
	tempPath := path + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	// Ensure temp file is removed on error
	defer func() {
		if err != nil {
			os.Remove(tempPath)
		}
	}()

	// Execute template
	if err = tmpl.Execute(file, data); err != nil {
		file.Close()
		return fmt.Errorf("executing template: %w", err)
	}

	// Close file
	if err = file.Close(); err != nil {
		return fmt.Errorf("closing file: %w", err)
	}

	// Make script executable
	if err = os.Chmod(tempPath, 0755); err != nil {
		return fmt.Errorf("setting permissions: %w", err)
	}

	// Atomic rename
	if err = os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("renaming file: %w", err)
	}

	return nil
}