package docs

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Indexer struct {
	cache map[string]DocumentIndex
}

func NewIndexer() *Indexer {
	return &Indexer{
		cache: make(map[string]DocumentIndex),
	}
}

func (i *Indexer) IndexDocuments(paths []string, scanner *Scanner) ([]DocumentIndex, error) {
	var documents []DocumentIndex
	
	for _, path := range paths {
		doc, err := i.indexDocument(path, scanner)
		if err != nil {
			continue
		}
		documents = append(documents, doc)
		i.cache[path] = doc
	}
	
	documents = i.groupAndSortDocuments(documents)
	
	return documents, nil
}

func (i *Indexer) indexDocument(path string, scanner *Scanner) (DocumentIndex, error) {
	info, err := os.Stat(path)
	if err != nil {
		return DocumentIndex{}, err
	}
	
	title := i.extractTitle(path)
	docType := scanner.GetDocumentType(path)
	
	return DocumentIndex{
		Path:       path,
		Title:      title,
		Type:       docType,
		ModifiedAt: info.ModTime(),
	}, nil
}

func (i *Indexer) extractTitle(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return filepath.Base(path)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	
	return filepath.Base(path)
}

func (i *Indexer) groupAndSortDocuments(documents []DocumentIndex) []DocumentIndex {
	typeOrder := map[DocumentType]int{
		DocTypePRD:          1,
		DocTypeArchitecture: 2,
		DocTypeEpic:         3,
		DocTypeStory:        4,
		DocTypeUnknown:      5,
	}
	
	sort.Slice(documents, func(a, b int) bool {
		if documents[a].Type != documents[b].Type {
			return typeOrder[documents[a].Type] < typeOrder[documents[b].Type]
		}
		return documents[a].Title < documents[b].Title
	})
	
	return documents
}

func (i *Indexer) GetCachedIndex(path string) (DocumentIndex, bool) {
	doc, exists := i.cache[path]
	return doc, exists
}

func (i *Indexer) ClearCache() {
	i.cache = make(map[string]DocumentIndex)
}