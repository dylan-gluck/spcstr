# Backend Architecture

## Service Architecture

Since spcstr is a CLI/TUI application, the "backend" consists of internal Go packages that handle state management, hook execution, and file operations.

### Service Organization
```
internal/
├── hooks/                      # Hook command implementations
│   ├── handlers/              # Individual hook handlers
│   ├── registry.go           # Hook command registration
│   └── executor.go           # Hook execution coordinator
├── state/                     # State management service
│   ├── manager.go            # State CRUD operations
│   ├── atomic.go             # Atomic file operations
│   └── watcher.go            # File change monitoring
├── docs/                      # Document management service
│   ├── scanner.go            # Document discovery
│   ├── indexer.go            # Document indexing
│   └── renderer.go           # Markdown rendering
└── config/                    # Configuration management
    ├── settings.go           # Application settings
    └── init.go               # Project initialization
```

### Hook Handler Architecture
```go
// Hook handler interface
type HookHandler interface {
    Name() string
    Execute(input []byte) error
}

// Hook registry for command routing
type HookRegistry struct {
    handlers map[string]HookHandler
}

func (r *HookRegistry) Register(handler HookHandler) {
    r.handlers[handler.Name()] = handler
}

func (r *HookRegistry) Execute(name string, input []byte) error {
    handler, exists := r.handlers[name]
    if !exists {
        return fmt.Errorf("unknown hook: %s", name)
    }
    return handler.Execute(input)
}

// Example hook handler implementation
type SessionStartHandler struct {
    stateManager *state.Manager
}

func (h *SessionStartHandler) Name() string {
    return "session_start"
}

func (h *SessionStartHandler) Execute(input []byte) error {
    var params SessionStartParams
    if err := json.Unmarshal(input, &params); err != nil {
        return err
    }
    
    return h.stateManager.InitializeState(params.SessionID)
}
```

## State Management Architecture

### Atomic Write Implementation
```go
// Atomic file operations for state safety
type AtomicWriter struct {
    basePath string
}

func (a *AtomicWriter) WriteJSON(filename string, data interface{}) error {
    // Marshal data to JSON
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }
    
    // Write to temporary file
    tempPath := filename + ".tmp"
    if err := os.WriteFile(tempPath, jsonData, 0644); err != nil {
        return err
    }
    
    // Atomic rename operation
    return os.Rename(tempPath, filename)
}

// State manager with atomic operations
type StateManager struct {
    writer *AtomicWriter
    logger *log.Logger
}

func (s *StateManager) UpdateState(sessionID string, update StateUpdate) error {
    // Load current state
    state, err := s.LoadState(sessionID)
    if err != nil {
        return err
    }
    
    // Apply update
    if err := update.Apply(state); err != nil {
        return err
    }
    
    // Update timestamp
    state.UpdatedAt = time.Now()
    
    // Atomic write
    statePath := filepath.Join(".spcstr", "sessions", sessionID, "state.json")
    return s.writer.WriteJSON(statePath, state)
}
```

### File Watching Integration
```go
// File system watcher for real-time updates
type StateWatcher struct {
    fsWatcher *fsnotify.Watcher
    eventChan chan StateChangeEvent
}

func (w *StateWatcher) WatchSessionDirectory() error {
    return w.fsWatcher.Add(filepath.Join(".spcstr", "sessions"))
}

func (w *StateWatcher) processEvents() {
    for {
        select {
        case event, ok := <-w.fsWatcher.Events:
            if !ok {
                return
            }
            if strings.HasSuffix(event.Name, "state.json") && event.Op&fsnotify.Write == fsnotify.Write {
                w.eventChan <- StateChangeEvent{
                    SessionID: extractSessionID(event.Name),
                    Path:      event.Name,
                }
            }
        }
    }
}
```

## Document Service Architecture

### Document Discovery and Indexing
```go
// Document scanner for Plan view
type DocumentScanner struct {
    basePaths []string
}

func (d *DocumentScanner) ScanDocuments() ([]DocumentIndex, error) {
    var documents []DocumentIndex
    
    for _, basePath := range d.basePaths {
        err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            
            if strings.HasSuffix(path, ".md") {
                doc, err := d.indexDocument(path, info)
                if err != nil {
                    return err
                }
                documents = append(documents, doc)
            }
            
            return nil
        })
        
        if err != nil && !os.IsNotExist(err) {
            return nil, err
        }
    }
    
    return documents, nil
}

func (d *DocumentScanner) indexDocument(path string, info os.FileInfo) (DocumentIndex, error) {
    // Extract title from markdown content
    content, err := os.ReadFile(path)
    if err != nil {
        return DocumentIndex{}, err
    }
    
    title := extractTitle(string(content))
    docType := classifyDocument(path)
    
    return DocumentIndex{
        Path:       path,
        Title:      title,
        Type:       docType,
        ModifiedAt: info.ModTime(),
    }, nil
}
```
