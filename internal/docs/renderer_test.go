package docs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderer_RenderMarkdown(t *testing.T) {
	tempDir := t.TempDir()
	renderer := NewRenderer()
	
	markdown := `# Test Document

This is a **bold** text and this is *italic*.

## Code Block

` + "```go" + `
func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

- Item 1
- Item 2
- Item 3
`
	
	filepath := filepath.Join(tempDir, "test.md")
	os.WriteFile(filepath, []byte(markdown), 0644)
	
	rendered, err := renderer.RenderMarkdown(filepath)
	if err != nil {
		t.Fatalf("RenderMarkdown returned error: %v", err)
	}
	
	if rendered == "" {
		t.Error("Rendered content is empty")
	}
	
	if !strings.Contains(rendered, "Test Document") {
		t.Error("Rendered content missing title")
	}
}

func TestRenderer_RenderMarkdownContent(t *testing.T) {
	renderer := NewRenderer()
	
	markdown := "# Title\n\nThis is a test."
	
	rendered, err := renderer.RenderMarkdownContent(markdown)
	if err != nil {
		t.Fatalf("RenderMarkdownContent returned error: %v", err)
	}
	
	if rendered == "" {
		t.Error("Rendered content is empty")
	}
	
	if !strings.Contains(rendered, "Title") {
		t.Error("Rendered content missing title")
	}
}

func TestRenderer_FallbackRender(t *testing.T) {
	renderer := NewRenderer()
	
	content := "Line 1\nLine 2\nLine 3"
	
	fallback := renderer.fallbackRender(content)
	
	if fallback != content {
		t.Errorf("Fallback render altered content: expected '%s', got '%s'", content, fallback)
	}
}

func TestRenderer_SetWidth(t *testing.T) {
	renderer := NewRenderer()
	
	err := renderer.SetWidth(80)
	if err != nil {
		t.Errorf("SetWidth returned error: %v", err)
	}
	
	err = renderer.SetWidth(120)
	if err != nil {
		t.Errorf("SetWidth returned error: %v", err)
	}
}

func TestRenderer_FileNotFound(t *testing.T) {
	renderer := NewRenderer()
	
	_, err := renderer.RenderMarkdown("/nonexistent/file.md")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}