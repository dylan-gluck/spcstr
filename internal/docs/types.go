package docs

import "time"

type DocumentType string

const (
	DocTypePRD          DocumentType = "prd"
	DocTypeArchitecture DocumentType = "architecture"
	DocTypeEpic         DocumentType = "epic"
	DocTypeStory        DocumentType = "story"
	DocTypeUnknown      DocumentType = "unknown"
)

type DocumentIndex struct {
	Path       string       `json:"path"`
	Title      string       `json:"title"`
	Type       DocumentType `json:"type"`
	ModifiedAt time.Time    `json:"modified_at"`
}

type Engine struct {
	scanner  *Scanner
	indexer  *Indexer
	Renderer *Renderer
}

func NewEngine(rootPath string) *Engine {
	return &Engine{
		scanner:  NewScanner(rootPath),
		indexer:  NewIndexer(),
		Renderer: NewRenderer(),
	}
}

func (e *Engine) ScanAndIndex() ([]DocumentIndex, error) {
	files, err := e.scanner.ScanForMarkdownFiles()
	if err != nil {
		return nil, err
	}
	
	return e.indexer.IndexDocuments(files, e.scanner)
}

func (e *Engine) RenderDocument(path string) (string, error) {
	return e.Renderer.RenderMarkdown(path)
}