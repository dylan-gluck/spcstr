package docs

import (
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	tea "github.com/charmbracelet/bubbletea"
)

type FileWatcher struct {
	watcher  *fsnotify.Watcher
	rootPath string
}

type FileChangeMsg struct {
	Path      string
	Operation string
}

func NewFileWatcher(rootPath string) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	
	docsPath := filepath.Join(rootPath, "docs")
	err = watcher.Add(docsPath)
	if err != nil {
		watcher.Close()
		return nil, err
	}
	
	err = addSubdirectories(watcher, docsPath)
	if err != nil {
		watcher.Close()
		return nil, err
	}
	
	return &FileWatcher{
		watcher:  watcher,
		rootPath: rootPath,
	}, nil
}

func addSubdirectories(watcher *fsnotify.Watcher, root string) error {
	subdirs := []string{"prd", "architecture", "epics", "stories"}
	for _, subdir := range subdirs {
		path := filepath.Join(root, subdir)
		watcher.Add(path)
	}
	return nil
}

func (fw *FileWatcher) Watch() tea.Cmd {
	return func() tea.Msg {
		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					return nil
				}
				
				if !strings.HasSuffix(strings.ToLower(event.Name), ".md") {
					continue
				}
				
				operation := "modified"
				if event.Op&fsnotify.Create == fsnotify.Create {
					operation = "created"
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					operation = "removed"
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					operation = "modified"
				}
				
				return FileChangeMsg{
					Path:      event.Name,
					Operation: operation,
				}
				
			case _, ok := <-fw.watcher.Errors:
				if !ok {
					return nil
				}
			}
		}
	}
}

func (fw *FileWatcher) Close() error {
	if fw.watcher != nil {
		return fw.watcher.Close()
	}
	return nil
}