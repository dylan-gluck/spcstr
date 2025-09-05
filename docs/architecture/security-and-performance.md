# Security and Performance

## Security Requirements

**Local Application Security:**
- Input Validation: All hook JSON input validated against schemas
- File Path Sanitization: Prevent directory traversal in state file operations  
- Safe Parsing: JSON unmarshaling with size limits and timeout constraints
- Permission Management: Minimal file system permissions (read/write only to `.spcstr/`)

**Data Privacy:**
- No Network Calls: Complete offline operation preserves user privacy
- Local Storage Only: All data remains in user's filesystem
- No Telemetry: Zero data collection or external reporting

**Process Security:**
- Hook Execution: Isolated hook processes with timeout constraints
- Resource Limits: Memory and CPU limits for hook operations
- Error Handling: Safe error messages that don't leak sensitive information

## Performance Optimization

**TUI Performance:**
- Render Optimization: Lazy rendering for large document lists and session data
- Memory Management: Efficient state management with garbage collection awareness
- Update Throttling: File watcher events throttled to prevent UI thrashing

**Hook Performance:**
- Fast Execution: Hook operations complete within 100ms to avoid blocking Claude Code
- Atomic Operations: State updates use filesystem-level atomic operations
- Efficient JSON: Minimal JSON parsing and serialization overhead

**File System Performance:**
- Efficient Indexing: Document scanning optimized with file modification time caching
- Minimal I/O: State loading/saving batched to reduce filesystem operations
- Watch Filtering: File watcher events filtered to relevant changes only
