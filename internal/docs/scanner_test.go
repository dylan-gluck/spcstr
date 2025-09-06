package docs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanner_ScanForMarkdownFiles(t *testing.T) {
	tempDir := t.TempDir()
	
	docsDir := filepath.Join(tempDir, "docs")
	os.MkdirAll(docsDir, 0755)
	os.MkdirAll(filepath.Join(docsDir, "prd"), 0755)
	os.MkdirAll(filepath.Join(docsDir, "stories"), 0755)
	
	testFiles := []string{
		filepath.Join(docsDir, "prd.md"),
		filepath.Join(docsDir, "architecture.md"),
		filepath.Join(docsDir, "prd", "overview.md"),
		filepath.Join(docsDir, "stories", "story1.md"),
		filepath.Join(docsDir, "readme.txt"),
	}
	
	for _, file := range testFiles[:4] {
		os.WriteFile(file, []byte("# Test Document"), 0644)
	}
	os.WriteFile(testFiles[4], []byte("Not markdown"), 0644)
	
	scanner := NewScanner(tempDir)
	files, err := scanner.ScanForMarkdownFiles()
	
	if err != nil {
		t.Fatalf("ScanForMarkdownFiles returned error: %v", err)
	}
	
	if len(files) != 4 {
		t.Errorf("Expected 4 markdown files, got %d", len(files))
	}
	
	for _, file := range files {
		if filepath.Ext(file) != ".md" {
			t.Errorf("Non-markdown file included: %s", file)
		}
	}
}

func TestScanner_GetDocumentType(t *testing.T) {
	scanner := NewScanner("/test")
	
	tests := []struct {
		path     string
		expected DocumentType
	}{
		{"/docs/prd.md", DocTypePRD},
		{"/docs/prd/overview.md", DocTypePRD},
		{"/docs/architecture.md", DocTypeArchitecture},
		{"/docs/architecture/components.md", DocTypeArchitecture},
		{"/docs/epics/epic-1.md", DocTypeEpic},
		{"/docs/stories/story-1.md", DocTypeStory},
		{"/docs/random.md", DocTypeUnknown},
	}
	
	for _, test := range tests {
		result := scanner.GetDocumentType(test.path)
		if result != test.expected {
			t.Errorf("GetDocumentType(%s) = %s, expected %s", test.path, result, test.expected)
		}
	}
}