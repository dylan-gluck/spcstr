package docs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dylan/spcstr/internal/config"
)

type Scanner struct {
	rootPath string
	config   *config.CoreConfig
}

func NewScanner(rootPath string) *Scanner {
	cfg, _ := config.LoadCoreConfig(rootPath)
	return &Scanner{
		rootPath: rootPath,
		config:   cfg,
	}
}

func (s *Scanner) ScanForMarkdownFiles() ([]string, error) {
	var markdownFiles []string
	
	docsPath := filepath.Join(s.rootPath, "docs")
	
	if _, err := os.Stat(docsPath); os.IsNotExist(err) {
		return markdownFiles, nil
	}
	
	if s.config != nil {
		if !s.config.PRD.PRDSharded {
			prdPath := filepath.Join(s.rootPath, s.config.PRD.PRDFile)
			if _, err := os.Stat(prdPath); err == nil {
				markdownFiles = append(markdownFiles, prdPath)
			}
		} else {
			s.scanDirectory(filepath.Join(s.rootPath, s.config.PRD.PRDShardedLocation), &markdownFiles)
		}
		
		if !s.config.Architecture.ArchitectureSharded {
			archPath := filepath.Join(s.rootPath, s.config.Architecture.ArchitectureFile)
			if _, err := os.Stat(archPath); err == nil {
				markdownFiles = append(markdownFiles, archPath)
			}
		} else {
			s.scanDirectory(filepath.Join(s.rootPath, s.config.Architecture.ArchitectureShardedLocation), &markdownFiles)
		}
		
		s.scanDirectory(filepath.Join(s.rootPath, s.config.DevStoryLocation), &markdownFiles)
		s.scanDirectory(filepath.Join(s.rootPath, "docs/epics"), &markdownFiles)
	}
	
	err := filepath.Walk(docsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if info.IsDir() {
			return nil
		}
		
		if strings.HasSuffix(strings.ToLower(path), ".md") {
			for _, existing := range markdownFiles {
				if existing == path {
					return nil
				}
			}
			markdownFiles = append(markdownFiles, path)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return markdownFiles, nil
}

func (s *Scanner) scanDirectory(dir string, files *[]string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return
	}
	
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		
		if strings.HasSuffix(strings.ToLower(path), ".md") {
			*files = append(*files, path)
		}
		return nil
	})
}

func (s *Scanner) GetDocumentType(path string) DocumentType {
	normalizedPath := strings.ToLower(path)
	
	if strings.Contains(normalizedPath, "/prd") || strings.Contains(normalizedPath, "prd.md") {
		return DocTypePRD
	}
	if strings.Contains(normalizedPath, "/architecture") || strings.Contains(normalizedPath, "architecture.md") {
		return DocTypeArchitecture
	}
	if strings.Contains(normalizedPath, "/epics") || strings.Contains(normalizedPath, "epic") {
		return DocTypeEpic
	}
	if strings.Contains(normalizedPath, "/stories") || strings.Contains(normalizedPath, "story") {
		return DocTypeStory
	}
	
	return DocTypeUnknown
}