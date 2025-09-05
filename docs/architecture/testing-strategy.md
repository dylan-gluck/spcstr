# Testing Strategy

## Testing Pyramid
```
      E2E Tests
     /          \
   Integration Tests
  /              \
 Unit Tests    Hook Tests
```

## Test Organization

### Unit Tests Structure
```
internal/
├── state/
│   ├── manager_test.go       # State CRUD operations
│   ├── atomic_test.go        # Atomic write operations
│   └── watcher_test.go       # File watching logic
├── hooks/
│   ├── handlers/
│   │   ├── session_start_test.go
│   │   ├── pre_tool_use_test.go
│   │   └── post_tool_use_test.go
│   └── registry_test.go      # Hook registration
└── tui/
    ├── views/
    │   ├── plan_test.go      # Plan view logic
    │   └── observe_test.go   # Observe view logic
    └── components/
        └── dashboard_test.go  # Dashboard rendering
```

### Integration Tests Structure
```
tests/
├── integration/
│   ├── hook_integration_test.go    # End-to-end hook execution
│   ├── state_integration_test.go   # State management workflows
│   └── tui_integration_test.go     # TUI navigation flows
└── testdata/
    ├── sample_sessions/           # Test session data
    ├── sample_docs/              # Test documents
    └── expected_outputs/         # Expected test results
```

### Manual Testing Structure
```
tests/
└── manual/
    ├── init_test_steps.md        # spcstr init testing
    ├── hook_test_steps.md        # Hook execution testing
    ├── tui_test_steps.md         # TUI navigation testing
    └── performance_test_steps.md # Performance validation
```

## Test Examples

### State Management Unit Test
```go
func TestAtomicWrite(t *testing.T) {
    tests := []struct {
        name     string
        data     interface{}
        wantErr  bool
    }{
        {
            name: "valid session state",
            data: &SessionState{
                SessionID:     "test_session",
                CreatedAt:     time.Now(),
                SessionActive: true,
            },
            wantErr: false,
        },
        {
            name:    "nil data",
            data:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tempDir := t.TempDir()
            writer := NewAtomicWriter(tempDir)
            filePath := filepath.Join(tempDir, "test.json")
            
            err := writer.WriteJSON(filePath, tt.data)
            if (err != nil) != tt.wantErr {
                t.Errorf("WriteJSON() error = %v, wantErr %v", err, tt.wantErr)
            }
            
            if !tt.wantErr {
                // Verify file was created and contains expected data
                if _, err := os.Stat(filePath); os.IsNotExist(err) {
                    t.Error("Expected file was not created")
                }
            }
        })
    }
}
```

### Hook Handler Integration Test
```go
func TestSessionStartHook(t *testing.T) {
    tempDir := t.TempDir()
    os.Chdir(tempDir)
    
    // Create .spcstr directory
    os.MkdirAll(".spcstr/sessions", 0755)
    
    handler := &SessionStartHandler{
        stateManager: state.NewManager(),
    }
    
    input := `{"session_id": "test_session_123", "source": "startup"}`
    
    err := handler.Execute([]byte(input))
    if err != nil {
        t.Fatalf("Hook execution failed: %v", err)
    }
    
    // Verify state file was created
    statePath := ".spcstr/sessions/test_session_123/state.json"
    if _, err := os.Stat(statePath); os.IsNotExist(err) {
        t.Error("State file was not created")
    }
    
    // Verify state content
    var state SessionState
    data, err := os.ReadFile(statePath)
    if err != nil {
        t.Fatalf("Cannot read state file: %v", err)
    }
    
    if err := json.Unmarshal(data, &state); err != nil {
        t.Fatalf("Cannot parse state JSON: %v", err)
    }
    
    if state.SessionID != "test_session_123" {
        t.Errorf("Expected session ID 'test_session_123', got '%s'", state.SessionID)
    }
    
    if !state.SessionActive {
        t.Error("Expected session to be active")
    }
}
```

### TUI Component Test
```go
func TestObserveDashboard(t *testing.T) {
    // Create test session data
    sessionState := &SessionState{
        SessionID:     "test_session",
        SessionActive: true,
        Agents:        []string{"research-agent"},
        ToolsUsed:     map[string]int{"Read": 3, "Write": 1},
        Files: FileOperations{
            New:    []string{"/test/new.go"},
            Edited: []string{"/test/edited.md"},
            Read:   []string{"/test/readme.md"},
        },
    }
    
    dashboard := NewDashboard()
    dashboard.SetSession(sessionState)
    
    // Test rendering
    content := dashboard.View()
    
    // Verify key information is displayed
    if !strings.Contains(content, "test_session") {
        t.Error("Session ID not displayed")
    }
    
    if !strings.Contains(content, "research-agent") {
        t.Error("Active agent not displayed")
    }
    
    if !strings.Contains(content, "Read: 3") {
        t.Error("Tool usage not displayed")
    }
}
```
