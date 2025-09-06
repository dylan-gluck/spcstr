package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dylan/spcstr/internal/docs"
	"github.com/dylan/spcstr/internal/tui/views/plan"
	tea "github.com/charmbracelet/bubbletea"
)

func TestPlanView_CompleteFlow(t *testing.T) {
	tempDir := t.TempDir()
	
	docsDir := filepath.Join(tempDir, "docs")
	os.MkdirAll(docsDir, 0755)
	os.MkdirAll(filepath.Join(docsDir, "prd"), 0755)
	os.MkdirAll(filepath.Join(docsDir, "architecture"), 0755)
	os.MkdirAll(filepath.Join(docsDir, "stories"), 0755)
	
	testFiles := []struct {
		path    string
		content string
	}{
		{
			path: filepath.Join(docsDir, "prd.md"),
			content: `# Product Requirements Document

## Overview
This is the main PRD for the project.

## Features
- Feature 1
- Feature 2
- Feature 3`,
		},
		{
			path: filepath.Join(docsDir, "architecture.md"),
			content: `# Architecture Document

## System Design
The system uses a modular architecture.

## Components
- Component A
- Component B`,
		},
		{
			path: filepath.Join(docsDir, "stories", "story-1.md"),
			content: `# Story 1: User Authentication

## Description
Implement user authentication system.`,
		},
	}
	
	for _, tf := range testFiles {
		err := os.WriteFile(tf.path, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.path, err)
		}
	}
	
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)
	
	model := plan.New()
	
	// First initialize the model
	initCmd := model.Init()
	if initCmd == nil {
		t.Error("Init should return a command to load documents")
	}
	
	// Then update with window size
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(plan.Model)
	
	time.Sleep(100 * time.Millisecond)
	
	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}
}

func TestPlanView_DocumentScanning(t *testing.T) {
	tempDir := t.TempDir()
	
	docsDir := filepath.Join(tempDir, "docs")
	os.MkdirAll(docsDir, 0755)
	os.MkdirAll(filepath.Join(docsDir, "prd"), 0755)
	
	testFiles := []string{
		filepath.Join(docsDir, "prd.md"),
		filepath.Join(docsDir, "architecture.md"),
		filepath.Join(docsDir, "prd", "features.md"),
		filepath.Join(docsDir, "prd", "requirements.md"),
	}
	
	for _, path := range testFiles {
		content := "# " + filepath.Base(path)
		os.WriteFile(path, []byte(content), 0644)
	}
	
	engine := docs.NewEngine(tempDir)
	documents, err := engine.ScanAndIndex()
	
	if err != nil {
		t.Fatalf("ScanAndIndex failed: %v", err)
	}
	
	if len(documents) != 4 {
		t.Errorf("Expected 4 documents, found %d", len(documents))
	}
	
	prdCount := 0
	archCount := 0
	for _, doc := range documents {
		switch doc.Type {
		case docs.DocTypePRD:
			prdCount++
		case docs.DocTypeArchitecture:
			archCount++
		}
	}
	
	if prdCount != 3 {
		t.Errorf("Expected 3 PRD documents, found %d", prdCount)
	}
	
	if archCount != 1 {
		t.Errorf("Expected 1 Architecture document, found %d", archCount)
	}
}

func TestPlanView_MarkdownRendering(t *testing.T) {
	tempDir := t.TempDir()
	
	testFile := filepath.Join(tempDir, "test.md")
	markdown := `# Test Document

This is a test with **bold** and *italic* text.

## Code Example

` + "```go" + `
func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

### Lists

1. First item
2. Second item
3. Third item

- Bullet 1
- Bullet 2
`
	
	os.WriteFile(testFile, []byte(markdown), 0644)
	
	engine := docs.NewEngine(tempDir)
	rendered, err := engine.RenderDocument(testFile)
	
	if err != nil {
		t.Fatalf("RenderDocument failed: %v", err)
	}
	
	if rendered == "" {
		t.Error("Rendered content should not be empty")
	}
	
	if len(rendered) <= len(markdown) {
		t.Error("Rendered content should include styling/formatting")
	}
}

func TestPlanView_FileWatching(t *testing.T) {
	tempDir := t.TempDir()
	
	docsDir := filepath.Join(tempDir, "docs")
	os.MkdirAll(docsDir, 0755)
	
	testFile := filepath.Join(docsDir, "test.md")
	os.WriteFile(testFile, []byte("# Initial Content"), 0644)
	
	watcher, err := docs.NewFileWatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create file watcher: %v", err)
	}
	defer watcher.Close()
	
	msgChan := make(chan tea.Msg, 1)
	go func() {
		cmd := watcher.Watch()
		msg := cmd()
		if msg != nil {
			msgChan <- msg
		}
	}()
	
	time.Sleep(100 * time.Millisecond)
	
	os.WriteFile(testFile, []byte("# Updated Content"), 0644)
	
	select {
	case msg := <-msgChan:
		if fileMsg, ok := msg.(docs.FileChangeMsg); ok {
			if fileMsg.Operation != "modified" {
				t.Errorf("Expected 'modified' operation, got '%s'", fileMsg.Operation)
			}
			if fileMsg.Path != testFile {
				t.Errorf("Expected path '%s', got '%s'", testFile, fileMsg.Path)
			}
		} else {
			t.Error("Expected FileChangeMsg type")
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for file change event")
	}
}

func TestPlanView_Performance(t *testing.T) {
	tempDir := t.TempDir()
	
	docsDir := filepath.Join(tempDir, "docs")
	os.MkdirAll(docsDir, 0755)
	storiesDir := filepath.Join(docsDir, "stories")
	os.MkdirAll(storiesDir, 0755)
	
	for i := 0; i < 50; i++ {
		filename := filepath.Join(storiesDir, fmt.Sprintf("story-%d.md", i))
		content := fmt.Sprintf("# Story %d\n\nContent for story %d", i, i)
		os.WriteFile(filename, []byte(content), 0644)
	}
	
	start := time.Now()
	engine := docs.NewEngine(tempDir)
	documents, err := engine.ScanAndIndex()
	elapsed := time.Since(start)
	
	if err != nil {
		t.Fatalf("ScanAndIndex failed: %v", err)
	}
	
	if len(documents) != 50 {
		t.Errorf("Expected 50 documents, found %d", len(documents))
	}
	
	if elapsed > 500*time.Millisecond {
		t.Errorf("Document indexing took %v, exceeds 500ms requirement", elapsed)
	}
	
	if len(documents) > 0 {
		start = time.Now()
		_, err = engine.RenderDocument(documents[0].Path)
		elapsed = time.Since(start)
		
		if err != nil {
			t.Fatalf("RenderDocument failed: %v", err)
		}
		
		if elapsed > 200*time.Millisecond {
			t.Errorf("Markdown rendering took %v, exceeds 200ms requirement", elapsed)
		}
	}
}