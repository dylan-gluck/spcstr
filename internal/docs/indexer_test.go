package docs

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIndexer_IndexDocuments(t *testing.T) {
	tempDir := t.TempDir()
	
	testFiles := []struct {
		path    string
		content string
		title   string
	}{
		{
			path:    filepath.Join(tempDir, "test1.md"),
			content: "# Test Document 1\n\nContent here",
			title:   "Test Document 1",
		},
		{
			path:    filepath.Join(tempDir, "test2.md"),
			content: "## Not a title\n# Actual Title\n\nMore content",
			title:   "Actual Title",
		},
		{
			path:    filepath.Join(tempDir, "test3.md"),
			content: "No title in this document",
			title:   "test3.md",
		},
	}
	
	for _, tf := range testFiles {
		os.WriteFile(tf.path, []byte(tf.content), 0644)
	}
	
	indexer := NewIndexer()
	scanner := NewScanner(tempDir)
	
	paths := []string{testFiles[0].path, testFiles[1].path, testFiles[2].path}
	documents, err := indexer.IndexDocuments(paths, scanner)
	
	if err != nil {
		t.Fatalf("IndexDocuments returned error: %v", err)
	}
	
	if len(documents) != 3 {
		t.Errorf("Expected 3 documents, got %d", len(documents))
	}
	
	titleMap := make(map[string]string)
	for _, tf := range testFiles {
		titleMap[tf.path] = tf.title
	}
	
	for i, doc := range documents {
		expectedTitle := titleMap[doc.Path]
		if doc.Title != expectedTitle {
			t.Errorf("Document %d: expected title '%s', got '%s'", i, expectedTitle, doc.Title)
		}
		
		if doc.ModifiedAt.IsZero() {
			t.Errorf("Document %d: ModifiedAt is zero", i)
		}
	}
}

func TestIndexer_ExtractTitle(t *testing.T) {
	tempDir := t.TempDir()
	indexer := NewIndexer()
	
	tests := []struct {
		content      string
		expectedTitle string
	}{
		{"# Main Title\n\nContent", "Main Title"},
		{"## Subtitle\n# Title\n\nContent", "Title"},
		{"No markdown header", "test.md"},
		{"#NoSpace", "test.md"},
		{"# Title with # symbols #", "Title with # symbols #"},
	}
	
	for i, test := range tests {
		filepath := filepath.Join(tempDir, "test.md")
		os.WriteFile(filepath, []byte(test.content), 0644)
		
		title := indexer.extractTitle(filepath)
		if title != test.expectedTitle {
			t.Errorf("Test %d: expected title '%s', got '%s'", i, test.expectedTitle, title)
		}
	}
}

func TestIndexer_GroupAndSortDocuments(t *testing.T) {
	indexer := NewIndexer()
	
	now := time.Now()
	documents := []DocumentIndex{
		{Path: "/story2.md", Title: "Story 2", Type: DocTypeStory, ModifiedAt: now},
		{Path: "/prd.md", Title: "PRD", Type: DocTypePRD, ModifiedAt: now},
		{Path: "/epic1.md", Title: "Epic 1", Type: DocTypeEpic, ModifiedAt: now},
		{Path: "/story1.md", Title: "Story 1", Type: DocTypeStory, ModifiedAt: now},
		{Path: "/arch.md", Title: "Architecture", Type: DocTypeArchitecture, ModifiedAt: now},
	}
	
	sorted := indexer.groupAndSortDocuments(documents)
	
	expectedOrder := []DocumentType{
		DocTypePRD,
		DocTypeArchitecture,
		DocTypeEpic,
		DocTypeStory,
		DocTypeStory,
	}
	
	for i, doc := range sorted {
		if doc.Type != expectedOrder[i] {
			t.Errorf("Document %d: expected type %s, got %s", i, expectedOrder[i], doc.Type)
		}
	}
	
	if sorted[3].Title > sorted[4].Title {
		t.Error("Stories not sorted alphabetically")
	}
}

func TestIndexer_Cache(t *testing.T) {
	indexer := NewIndexer()
	
	doc := DocumentIndex{
		Path:       "/test.md",
		Title:      "Test",
		Type:       DocTypePRD,
		ModifiedAt: time.Now(),
	}
	
	indexer.cache[doc.Path] = doc
	
	cached, exists := indexer.GetCachedIndex(doc.Path)
	if !exists {
		t.Error("Document not found in cache")
	}
	
	if cached.Title != doc.Title {
		t.Errorf("Cached title mismatch: expected '%s', got '%s'", doc.Title, cached.Title)
	}
	
	indexer.ClearCache()
	
	_, exists = indexer.GetCachedIndex(doc.Path)
	if exists {
		t.Error("Cache not cleared properly")
	}
}