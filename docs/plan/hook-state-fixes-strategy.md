# Hook State Management Fixes Strategy

## Current Issues Identified

### 1. Missing Data Extraction from Claude Events
**Problem**: Hook handlers are only extracting basic fields (session_id, tool_name, agent_name) but ignoring the rich data available in the actual events.

**Evidence**: 
- TodoWrite events contain full `tool_input.todos` array with task details
- File operations have `tool_input.file_path` but handlers expect different field names
- Task tool events contain agent details in `tool_input` not being captured

### 2. Incorrect Field Mapping
**Problem**: Handlers expect simplified field names that don't match actual Claude event structure.

**Current Handler Expectations**:
```go
type PostToolUseParams struct {
    SessionID    string   `json:"session_id"`
    ToolName     string   `json:"tool_name"`
    AgentName    string   `json:"agent_name,omitempty"`
    FilesCreated []string `json:"files_created,omitempty"`
    FilesEdited  []string `json:"files_edited,omitempty"`
    FilesRead    []string `json:"files_read,omitempty"`
}
```

**Actual Event Structure**:
```json
{
  "session_id": "...",
  "tool_name": "Write",
  "tool_input": {
    "file_path": "/path/to/file",
    "content": "..."
  },
  "tool_response": {
    "filePath": "/path/to/file",
    "type": "create|edit"
  }
}
```

### 3. No TodoWrite Support
**Problem**: TodoWrite tool events are not being processed to track todo items.

**Required Data**:
- Total todos count
- Pending count (status: "pending")
- In-progress count (status: "in_progress")
- Completed count (status: "completed")
- Current/recent todo items for display

### 4. Missing Error Tracking
**Problem**: Failed pre_tool_use events (blocked operations) are not being recorded as errors in state.

### 5. No Subagent/Task Tracking
**Problem**: Task tool invocations contain agent information in `tool_input` that's not being extracted.

## Implementation Strategy

### Phase 1: Fix Event Structure Parsing

#### 1.1 Create Comprehensive Event Types
```go
// internal/hooks/events/types.go
type ClaudeEvent struct {
    SessionID       string          `json:"session_id"`
    HookEventName   string          `json:"hook_event_name"`
    ToolName        string          `json:"tool_name"`
    ToolInput       json.RawMessage `json:"tool_input"`
    ToolResponse    json.RawMessage `json:"tool_response"`
    PermissionMode  string          `json:"permission_mode"`
    TranscriptPath  string          `json:"transcript_path"`
    CWD            string          `json:"cwd"`
}

// Tool-specific input structures
type TodoWriteInput struct {
    Todos []TodoItem `json:"todos"`
}

type TodoItem struct {
    Content    string `json:"content"`
    Status     string `json:"status"`
    ActiveForm string `json:"activeForm"`
}

type FileOperationInput struct {
    FilePath string `json:"file_path"`
    Content  string `json:"content,omitempty"`
}

type TaskInput struct {
    Description   string `json:"description"`
    Prompt       string `json:"prompt"`
    SubagentType string `json:"subagent_type"`
}
```

#### 1.2 Update Handler Processing Logic
Each handler should:
1. Parse the full ClaudeEvent structure
2. Extract tool-specific data based on tool_name
3. Process according to tool type

### Phase 2: Enhance State Management

#### 2.1 Add Todo Tracking to State
```go
type SessionState struct {
    // ... existing fields ...
    
    Todos TodoState `json:"todos"`
}

type TodoState struct {
    Total       int        `json:"total"`
    Pending     int        `json:"pending"`
    InProgress  int        `json:"in_progress"`
    Completed   int        `json:"completed"`
    Recent      []TodoItem `json:"recent"` // Last 5 items
    LastUpdated string     `json:"last_updated"`
}
```

#### 2.2 File Operation Detection
Parse tool_response to determine operation type:
- Write tool with `type: "create"` → new file
- Write/Edit tools with `type: "edit"` → edited file
- Read tool → read file

### Phase 3: Implement Specific Handlers

#### 3.1 PostToolUse Handler Updates
```go
func (h *PostToolUseHandler) Execute(input []byte) error {
    var event ClaudeEvent
    json.Unmarshal(input, &event)
    
    switch event.ToolName {
    case "TodoWrite":
        return h.processTodoWrite(event)
    case "Write", "Edit", "MultiEdit":
        return h.processFileWrite(event)
    case "Read":
        return h.processFileRead(event)
    case "Task":
        return h.processTask(event)
    default:
        // Just track tool usage
    }
}
```

#### 3.2 PreToolUse Handler Updates
- Track tool usage counts
- For Task tool: Extract agent details from tool_input
- Handle blocking/errors appropriately

### Phase 4: Error Handling

#### 4.1 Track Hook Failures
- When pre_tool_use returns error (blocks operation)
- Record in state.errors array with timestamp and details
- Ensure error logging doesn't block operation

#### 4.2 Fix Logging Issues
Current issue: "Warning: Failed to log hook event: failed to rename temp log"
- Investigate atomic write operations in logger
- Ensure proper file permissions
- Add retry logic for transient failures

### Phase 5: Remove Incorrect Logic

#### 5.1 Remove agent_name from subagent_stop
- This field doesn't exist in the event
- Agent tracking should only happen via Task tool events

#### 5.2 Remove Blocking Logic
- No security filtering in hooks
- Let Claude handle permission decisions
- Hooks should only track/log, never block

## Implementation Order

1. **Update event structures** (events/types.go)
2. **Fix post_tool_use handler** to parse actual event structure
3. **Add TodoWrite support** with todo state tracking
4. **Fix file operation detection** using tool_response
5. **Add error tracking** for failed operations
6. **Update state structure** with new fields
7. **Fix logging issues** preventing proper event recording
8. **Test with real Claude sessions** to verify

## Testing Strategy

1. Create test fixtures from actual log data
2. Unit test each handler with real event structures
3. Integration test with mock Claude events
4. Manual testing with live Claude sessions
5. Verify state files contain expected data

## Success Criteria

- [ ] Todo counts appear in session state
- [ ] File operations correctly categorized (new/edited/read)
- [ ] Errors tracked in state
- [ ] No blocking of legitimate operations
- [ ] Agent/subagent tracking works for Task tool
- [ ] Logging works without errors
- [ ] TUI can display all tracked data