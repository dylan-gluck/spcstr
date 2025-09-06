package docs

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

type Renderer struct {
	glamourRenderer *glamour.TermRenderer
}

func NewRenderer() *Renderer {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		renderer, _ = glamour.NewTermRenderer(
			glamour.WithStandardStyle("notty"),
		)
	}
	
	return &Renderer{
		glamourRenderer: renderer,
	}
}

func (r *Renderer) RenderMarkdown(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	rendered, err := r.glamourRenderer.Render(string(content))
	if err != nil {
		return r.fallbackRender(string(content)), nil
	}
	
	return rendered, nil
}

func (r *Renderer) RenderMarkdownContent(content string) (string, error) {
	rendered, err := r.glamourRenderer.Render(content)
	if err != nil {
		return r.fallbackRender(content), nil
	}
	
	return rendered, nil
}

func (r *Renderer) fallbackRender(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	
	for _, line := range lines {
		result = append(result, line)
	}
	
	return strings.Join(result, "\n")
}

func (r *Renderer) SetWidth(width int) error {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return err
	}
	
	r.glamourRenderer = renderer
	return nil
}