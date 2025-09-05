# Database Schema

Since Spec⭐️ uses **file-based JSON persistence** rather than a traditional database, here are the JSON schemas for our data structures:

## Session State Schema (`.spcstr/sessions/{session_id}/state.json`)
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["session_id", "created_at", "updated_at", "source", "project_path", "status"],
  "properties": {
    "session_id": {
      "type": "string",
      "pattern": "^sess_[a-zA-Z0-9-]+$",
      "description": "Unique session identifier with sess_ prefix"
    },
    "created_at": {
      "type": "string",
      "format": "date-time",
      "description": "ISO 8601 timestamp of session creation"
    },
    "updated_at": {
      "type": "string",
      "format": "date-time",
      "description": "ISO 8601 timestamp of last update"
    },
    "source": {
      "type": "string",
      "enum": ["startup", "resume", "clear"],
      "description": "How the session was initiated"
    },
    "project_path": {
      "type": "string",
      "description": "Absolute path to project directory"
    },
    "status": {
      "type": "string",
      "enum": ["active", "completed", "error"],
      "description": "Current session state"
    },
    "agents": {
      "type": "array",
      "items": {"type": "string"},
      "description": "Currently active agent names"
    },
    "agents_history": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["name", "started_at"],
        "properties": {
          "name": {
            "type": "string",
            "description": "Agent name/identifier"
          },
          "started_at": {
            "type": "string",
            "format": "date-time"
          },
          "ended_at": {
            "type": "string",
            "format": "date-time"
          }
        }
      },
      "description": "Complete history of all agents that have run"
    },
    "files": {
      "type": "object",
      "required": ["new", "edited", "read"],
      "properties": {
        "new": {
          "type": "array",
          "items": {"type": "string"},
          "description": "Absolute paths of created files"
        },
        "edited": {
          "type": "array",
          "items": {"type": "string"},
          "description": "Absolute paths of modified files"
        },
        "read": {
          "type": "array",
          "items": {"type": "string"},
          "description": "Absolute paths of read files"
        }
      },
      "description": "Categorized file operations"
    },
    "tools_used": {
      "type": "object",
      "additionalProperties": {
        "type": "integer",
        "minimum": 0
      },
      "description": "Map of tool names to usage counts"
    },
    "errors": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["timestamp", "hook", "message"],
        "properties": {
          "timestamp": {
            "type": "string",
            "format": "date-time"
          },
          "hook": {
            "type": "string",
            "description": "Hook name where error occurred"
          },
          "message": {
            "type": "string",
            "description": "Error description"
          }
        }
      },
      "description": "Structured error tracking"
    }
  }
}
```

## Configuration Schema (`.spcstr/config.json`)
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["version", "scope"],
  "properties": {
    "version": {
      "type": "string",
      "pattern": "^\\d+\\.\\d+\\.\\d+$"
    },
    "scope": {
      "type": "string",
      "enum": ["global", "project"]
    },
    "paths": {
      "type": "object",
      "properties": {
        "hooks": {"type": "string"},
        "sessions": {"type": "string"},
        "docs": {
          "type": "array",
          "items": {"type": "string"}
        }
      }
    },
    "ui": {
      "type": "object",
      "properties": {
        "theme": {
          "type": "string",
          "enum": ["default", "dark", "light", "high-contrast"]
        },
        "keyBindings": {
          "type": "object",
          "additionalProperties": {"type": "string"}
        },
        "refreshRate": {
          "type": "integer",
          "minimum": 100,
          "maximum": 5000
        }
      }
    },
    "session": {
      "type": "object",
      "properties": {
        "retentionDays": {
          "type": "integer",
          "minimum": 1
        },
        "autoArchive": {"type": "boolean"},
        "maxActiveSession": {"type": "integer"}
      }
    }
  }
}
```

## Document Index Cache Schema (`.spcstr/cache/document-index.json`)
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["version", "lastUpdated", "documents"],
  "properties": {
    "version": {"type": "string"},
    "lastUpdated": {
      "type": "string",
      "format": "date-time"
    },
    "documents": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["documentID", "filePath", "type", "title", "lastModified"],
        "properties": {
          "documentID": {"type": "string"},
          "filePath": {"type": "string"},
          "type": {
            "type": "string",
            "enum": ["prd", "architecture", "epic", "story", "config", "other"]
          },
          "title": {"type": "string"},
          "lastModified": {
            "type": "string",
            "format": "date-time"
          },
          "content": {"type": "string"},
          "metadata": {
            "type": "object",
            "additionalProperties": true
          },
          "searchTerms": {
            "type": "array",
            "items": {"type": "string"}
          }
        }
      }
    }
  }
}
```

## File System Layout
```
.spcstr/
├── config.json                 # Project configuration
├── sessions/                   # Session data directory
│   └── {session_id}/           # Per-session directory
│       ├── state.json          # Primary session state
│       ├── messages.json       # Message history (optional)
│       └── .lock              # Lock file for atomic operations
├── cache/                      # Application cache
│   └── document-index.json
└── logs/                       # Debug logs
    ├── debug.log
    └── hook-errors.log         # Hook-specific errors
```

## Index Strategies

Since we're using file-based storage, we implement indexing through:

1. **Session Index**: Directory-based organization with one directory per session
2. **Document Index**: Cached JSON file with pre-computed search terms
3. **Session State Files**: Atomic state.json per session for consistency
4. **Lock Files**: .lock files for concurrent access control

## Performance Considerations

- **Write Performance**: Atomic writes using temp file + rename (POSIX atomic)
- **Read Performance**: Direct JSON parsing with <10ms target
- **Concurrency**: File-based locking with exponential backoff (max 100ms wait)
- **Caching**: In-memory state cache with fsnotify invalidation
- **Hook Performance**: <10ms execution requirement, <2ms overhead target
