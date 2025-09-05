# Database Schema

*No traditional database - All data persisted as JSON files*

## File System Schema

```
.spcstr/
├── sessions/                    # Session state directory
│   └── {session-id}/
│       └── state.json          # SessionState structure
└── logs/                       # Hook execution logs
    ├── session_start.json      # Array of session start events
    ├── user_prompt_submit.json # Array of prompt events
    ├── pre_tool_use.json       # Array of tool invocation events
    ├── post_tool_use.json      # Array of tool completion events
    ├── notification.json       # Array of notification events
    ├── pre_compact.json        # Array of compaction events
    ├── session_end.json        # Array of session end events
    ├── stop.json              # Array of stop events
    └── subagent_stop.json     # Array of subagent stop events
```

## JSON Schema Examples

**state.json structure:**
```json
{
  "session_id": "claude_session_20250905_143022",
  "created_at": "2025-09-05T14:30:22Z",
  "updated_at": "2025-09-05T14:35:15Z",
  "session_active": true,
  "agents": ["research-agent"],
  "agents_history": [
    {
      "name": "research-agent",
      "started_at": "2025-09-05T14:32:10Z"
    }
  ],
  "files": {
    "new": ["/project/src/main.go"],
    "edited": ["/project/README.md"],
    "read": ["/project/go.mod"]
  },
  "tools_used": {
    "Read": 3,
    "Write": 1,
    "Task": 1
  },
  "errors": [],
  "prompts": [
    {
      "timestamp": "2025-09-05T14:30:25Z",
      "prompt": "Create a Go CLI application"
    }
  ],
  "notifications": []
}
```

**Log file structure (append-only arrays):**
```json
[
  {
    "timestamp": "2025-09-05T14:30:22Z",
    "session_id": "claude_session_20250905_143022",
    "source": "startup"
  },
  {
    "timestamp": "2025-09-05T15:20:15Z",
    "session_id": "claude_session_20250905_152010",
    "source": "resume"
  }
]
```
