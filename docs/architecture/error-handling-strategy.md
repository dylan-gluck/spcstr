# Error Handling Strategy

## General Approach
- **Error Model:** Go's explicit error handling with wrapped errors for context
- **Exception Hierarchy:** Structured error types with error codes for classification
- **Error Propagation:** Errors bubble up with context using fmt.Errorf with %w verb

## Logging Standards
- **Library:** Standard library log/slog
- **Format:** JSON for production, text for development
- **Levels:** DEBUG, INFO, WARN, ERROR, FATAL
- **Required Context:**
  - Correlation ID: UUID per session (format: sess_{uuid})
  - Service Context: Component name and operation
  - User Context: Never log sensitive data, only session ID

## Error Handling Patterns

### External API Errors
- **Retry Policy:** N/A - No external APIs in this application
- **Circuit Breaker:** N/A - No external services
- **Timeout Configuration:** File operations timeout after 5 seconds
- **Error Translation:** File system errors wrapped with operation context

### Business Logic Errors
- **Custom Exceptions:** 
  ```go
  type SessionError (InvalidState, NotFound, Corrupted)
  type ConfigError (Invalid, Missing, VersionMismatch)
  type HookError (ScriptFailed, PermissionDenied)
  ```
- **User-Facing Errors:** Human-readable messages with suggested actions
- **Error Codes:** SPCSTR-XXXX format (e.g., SPCSTR-1001: Session not found)

### Data Consistency
- **Transaction Strategy:** Atomic file writes via temp file + rename
- **Compensation Logic:** Automatic rollback on partial writes
- **Idempotency:** All hook operations are idempotent by design

## Error Categories and Handling

**Fatal Errors (Application exits):**
- Unable to access .spcstr directory
- Configuration corruption
- Terminal initialization failure
- Action: Log error, display message, exit with code 1

**Recoverable Errors (Operation retried):**
- File lock contention
- Temporary file system issues
- Action: Exponential backoff retry (3 attempts max)

**User Errors (User notification):**
- Invalid command syntax
- Missing required files
- Action: Display helpful error with usage instructions

**Warning Level (Logged but continues):**
- Outdated session files
- Missing optional configuration
- Action: Log warning, use defaults, continue operation

## Hook Script Error Handling

**Critical Requirement:** Hook scripts must NEVER block Claude Code

```bash
// Error handling in Go hooks
func (h *HookHandler) Execute() error {
    // Operations that may fail are logged but don't block
    if err := h.updateState(); err != nil {
        h.logError(err) // Log but continue
    }
    return nil // Always succeed to not block Claude
}
```

## TUI Error Display

**Error Notification Levels:**
1. **Status Bar Warning:** Yellow text for non-critical issues
2. **Modal Dialog:** Red border for errors requiring user action
3. **Error Section:** Persistent error log in dashboard

**Error Recovery Options:**
- Retry operation
- Skip and continue
- View detailed log
- Export error report

## Debug Logging

**Debug Log Location:** `.spcstr/logs/debug.log`

**Debug Information Captured:**
- Full stack traces for errors
- File operation details
- Event bus message flow
- Performance metrics
- Session state transitions

**Log Rotation:**
- Maximum file size: 10MB
- Keep last 3 files
- Automatic compression of old logs
