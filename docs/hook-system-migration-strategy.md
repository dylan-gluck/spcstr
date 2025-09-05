# Hook System Migration Strategy

## Executive Summary

This document outlines the migration strategy from the current shell-based hook implementation to an integrated Go-based hook system within the main spcstr executable. The migration addresses critical issues with data model consistency, state management robustness, and performance requirements.

## Current State Analysis

### Issues with Current Implementation
1. **Fragile JSON Parsing**: Shell scripts use `grep`/`cut` for JSON parsing, breaking on nested objects or escaped quotes
2. **Race Conditions**: Direct file writes without locking cause state corruption during concurrent operations
3. **Incorrect Data Model**: Flat arrays instead of categorized structures, missing tool usage counts
4. **Performance Overhead**: 8-12ms shell spawning overhead vs <10ms requirement
5. **No Safety Checks**: Missing validation for dangerous operations

### Current vs Target State Comparison

| Component | Current (Shell) | Target (Go) |
|-----------|----------------|-------------|
| Execution Model | Separate shell scripts | Integrated `spcstr hook` command |
| JSON Handling | grep/cut parsing | Native Go marshaling |
| State Updates | Direct file writes | Atomic operations with temp files |
| Performance | 8-12ms overhead | <2ms overhead |
| Type Safety | None | Compile-time validation |
| Error Handling | Silent failures | Structured error management |

## Target Architecture

### Command Structure
```bash
# Claude settings.json configuration
spcstr hook <hook-name> [--cwd=$CLAUDE_PROJECT_DIR]

# Valid hook names (CLI arguments):
- pre-tool-use
- post-tool-use
- user-prompt-submit
- notification
- stop
- subagent-stop
- pre-compact
- session-start
- session-end
```

### Package Organization
```
spcstr/
├── cmd/
│   └── hook.go                 # Hook subcommand entry
├── pkg/
│   └── hooks/
│       ├── handler.go          # Main dispatcher
│       ├── state.go            # Session state management
│       ├── types.go            # Data structures
│       ├── persistence.go      # File I/O operations
│       └── handlers/
│           ├── pre_tool_use.go
│           ├── post_tool_use.go
│           └── ...
```

## Data Model Specification

### Session State Structure
```go
type SessionState struct {
    // Core identification
    SessionID      string                 `json:"session_id"`
    CreatedAt      time.Time             `json:"created_at"`
    UpdatedAt      time.Time             `json:"updated_at"`

    // Session metadata
    Source         string                `json:"source"`         // startup|resume|clear
    ProjectPath    string                `json:"project_path"`
    Status         string                `json:"status"`         // active|completed|error

    // Agent tracking
    Agents         []string              `json:"agents"`         // Currently active agents
    AgentsHistory  []AgentHistoryEntry   `json:"agents_history"` // All agents that have run

    // File operations (categorized)
    Files          FileOperations        `json:"files"`

    // Tool usage tracking
    ToolsUsed      map[string]int        `json:"tools_used"`     // Tool name -> count

    // Error tracking
    Errors         []ErrorEntry          `json:"errors"`

    // Internal state
    Modified       bool                  `json:"-"`              // Dirty flag (not persisted)
}

type FileOperations struct {
    New    []string `json:"new"`    // Created files (absolute paths)
    Edited []string `json:"edited"` // Modified files (absolute paths)
    Read   []string `json:"read"`   // Read files (absolute paths)
}

type AgentHistoryEntry struct {
    Name      string    `json:"name"`
    StartedAt time.Time `json:"started_at"`
    EndedAt   time.Time `json:"ended_at,omitempty"`
}

type ErrorEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Hook      string    `json:"hook"`
    Message   string    `json:"message"`
}
```

## State Persistence Strategy

### File Structure
```
.spcstr/
└── sessions/
    └── {session_id}/
        ├── state.json       # Primary session state
        ├── messages.json    # Message history (optional)
        └── .lock           # Lock file for atomic operations
```

### Atomic Write Pattern
```go
func (s *SessionState) SaveAtomic() error {
    // 1. Marshal state to JSON
    data, err := json.MarshalIndent(s, "", "  ")
    if err != nil {
        return err
    }

    // 2. Write to temporary file
    tempFile := filepath.Join(s.sessionDir(), "state.tmp")
    if err := os.WriteFile(tempFile, data, 0644); err != nil {
        return err
    }

    // 3. Atomic rename (POSIX atomic operation)
    stateFile := filepath.Join(s.sessionDir(), "state.json")
    return os.Rename(tempFile, stateFile)
}
```

### Concurrency Management
- Use file-based locking with `.lock` file
- Implement exponential backoff for lock acquisition
- Maximum 100ms wait time before proceeding (hooks must not block)

## Migration Steps

### Phase 1: Foundation (Week 1)
1. **Implement Go hook package**
   - Create `pkg/hooks` package structure
   - Implement `SessionState` type with methods
   - Add atomic file operations
   - Create hook handler dispatcher

2. **Add `spcstr hook` command**
   - Integrate with Cobra command structure
   - Add JSON stdin/stdout handling
   - Implement exit code logic (0 normal, 2 blocking)

3. **Create test harness**
   - Unit tests for state management
   - Integration tests for atomic operations
   - Benchmark tests for performance validation

### Phase 2: Hook Implementation (Week 2)
1. **Implement critical hooks first**
   - `session-start`: Initialize state structure
   - `pre-tool-use`: Add safety checks, track usage
   - `post-tool-use`: Track file operations
   - `session-end`: Finalize and archive state

2. **Implement remaining hooks**
   - `user-prompt-submit`: Prompt filtering/enhancement
   - `notification`: Log notifications
   - `stop`: Track completion
   - `subagent-stop`: Update agent history
   - `pre-compact`: Prepare for compaction

### Phase 3: Migration & Testing (Week 3)
1. **Direct cutover**
   ```bash
   #!/bin/bash
   # Remove old shell hooks
   rm -rf .spcstr/hooks

   # Update Claude settings.json
   spcstr configure-hooks

   # Clean up old session files
   rm -rf .spcstr/sessions/*.json
   ```

2. **Full deployment**
   - Deploy all hooks simultaneously
   - Fresh session state initialization
   - Remove all shell script dependencies

3. **Performance validation**
   - Measure hook execution time
   - Validate <10ms requirement
   - Check memory usage
   - Monitor for file descriptor leaks

### Phase 4: Documentation & Cleanup (Week 4)
1. **Update documentation**
   - Update architecture docs to reflect Go hooks
   - Revise PRD with correct data model
   - Create troubleshooting guide

2. **Remove legacy code**
   - Delete all shell hook scripts
   - Remove shell-specific utilities
   - Clean up old configuration files

## Acceptance Criteria

### Functional Requirements
- [ ] All 9 Claude hooks implemented in Go
- [ ] Session state persists across restarts
- [ ] File operations correctly categorized (new/edited/read)
- [ ] Tool usage counts tracked accurately
- [ ] Agent history maintained with timestamps
- [ ] Dangerous operations blocked (rm -rf, .env access)

### Performance Requirements
- [ ] Hook execution time <10ms (measured)
- [ ] Memory usage <10MB per session
- [ ] No file descriptor leaks
- [ ] Atomic writes prevent corruption
- [ ] Concurrent hook execution handled safely

### Integration Requirements
- [ ] TUI receives real-time updates via file watching
- [ ] Session files remain human-readable JSON
- [ ] Clean error messages for debugging
- [ ] Silent failures don't disrupt Claude

### Testing Requirements
- [ ] Unit test coverage >80%
- [ ] Integration tests for all hooks
- [ ] Stress test with 100+ concurrent operations
- [ ] Performance benchmarks documented

## Risk Mitigation

### Monitoring & Alerts
- Log hook execution times
- Track error rates
- Monitor session file sizes
- Alert on performance degradation

### Known Limitations
- Initial migration requires Claude restart
- Some edge cases in concurrent updates

## Success Metrics

### Immediate (Week 1)
- Hook execution time reduced by 75%
- Zero session corruption incidents
- All hooks responding within 10ms

### Short-term (Month 1)
- 100% of sessions tracked accurately
- File operations categorization 100% accurate
- No reported data loss incidents

### Long-term (Quarter)
- Maintenance time reduced by 50%
- New hook additions take <1 hour
- Zero critical bugs in production

## Appendix: Configuration Examples

### Claude settings.json
```json
{
  "hooks": {
    "PreToolUse": "spcstr hook pre-tool-use --cwd=$CLAUDE_PROJECT_DIR",
    "PostToolUse": "spcstr hook post-tool-use --cwd=$CLAUDE_PROJECT_DIR",
    "UserPromptSubmit": "spcstr hook user-prompt-submit --cwd=$CLAUDE_PROJECT_DIR",
    "Notification": "spcstr hook notification --cwd=$CLAUDE_PROJECT_DIR",
    "Stop": "spcstr hook stop --cwd=$CLAUDE_PROJECT_DIR",
    "SubagentStop": "spcstr hook subagent-stop --cwd=$CLAUDE_PROJECT_DIR",
    "PreCompact": "spcstr hook pre-compact --cwd=$CLAUDE_PROJECT_DIR",
    "SessionStart": "spcstr hook session-start --cwd=$CLAUDE_PROJECT_DIR",
    "SessionEnd": "spcstr hook session-end --cwd=$CLAUDE_PROJECT_DIR"
  }
}
```

### Example Session State
```json
{
  "session_id": "sess_abc123",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:45:00Z",
  "source": "startup",
  "project_path": "/Users/dylan/Workspace/projects/spcstr",
  "status": "active",
  "agents": ["meta-agent"],
  "agents_history": [
    {
      "name": "meta-agent",
      "started_at": "2025-01-15T10:30:00Z"
    }
  ],
  "files": {
    "new": ["/Users/dylan/Workspace/projects/spcstr/pkg/hooks/handler.go"],
    "edited": ["/Users/dylan/Workspace/projects/spcstr/cmd/hook.go"],
    "read": ["/Users/dylan/Workspace/projects/spcstr/docs/prd.md"]
  },
  "tools_used": {
    "Read": 5,
    "Write": 1,
    "Edit": 2,
    "Bash": 3
  },
  "errors": []
}
```

## Timeline

- **Week 1**: Foundation implementation
- **Week 2**: Hook handlers complete
- **Week 3**: Migration and testing
- **Week 4**: Documentation and cleanup
- **Week 5**: Production monitoring
- **Week 6**: Performance optimization

## Approval

- [ ] Engineering Lead
- [ ] Product Owner
- [ ] QA Lead
- [ ] DevOps Lead
