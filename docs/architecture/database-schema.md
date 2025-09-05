# Database Schema

Since Spec⭐️ uses **file-based JSON persistence** rather than a traditional database, here are the JSON schemas for our data structures:

## Session Schema (`.spcstr/sessions/{session-id}.json`)
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["sessionID", "status", "startTime", "projectPath"],
  "properties": {
    "sessionID": {
      "type": "string",
      "pattern": "^[a-zA-Z0-9-]+$"
    },
    "status": {
      "type": "string",
      "enum": ["active", "completed", "error"]
    },
    "startTime": {
      "type": "string",
      "format": "date-time"
    },
    "endTime": {
      "type": ["string", "null"],
      "format": "date-time"
    },
    "projectPath": {
      "type": "string"
    },
    "metadata": {
      "type": "object",
      "additionalProperties": {"type": "string"}
    },
    "agents": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["agentID", "name", "type", "status"],
        "properties": {
          "agentID": {"type": "string"},
          "name": {"type": "string"},
          "type": {"type": "string"},
          "status": {
            "type": "string",
            "enum": ["idle", "active", "error"]
          },
          "lastActivity": {
            "type": "string",
            "format": "date-time"
          }
        }
      }
    },
    "tasks": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["taskID", "content", "status", "createdAt"],
        "properties": {
          "taskID": {"type": "string"},
          "content": {"type": "string"},
          "status": {
            "type": "string",
            "enum": ["pending", "in_progress", "completed"]
          },
          "createdAt": {
            "type": "string",
            "format": "date-time"
          },
          "completedAt": {
            "type": ["string", "null"],
            "format": "date-time"
          },
          "agentID": {"type": "string"}
        }
      }
    },
    "fileOperations": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["operationID", "filePath", "operation", "timestamp"],
        "properties": {
          "operationID": {"type": "string"},
          "filePath": {"type": "string"},
          "operation": {
            "type": "string",
            "enum": ["created", "edited", "read", "deleted"]
          },
          "timestamp": {
            "type": "string",
            "format": "date-time"
          },
          "lineCount": {"type": "integer"}
        }
      }
    },
    "toolExecutions": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["executionID", "toolName", "timestamp"],
        "properties": {
          "executionID": {"type": "string"},
          "toolName": {"type": "string"},
          "command": {"type": "string"},
          "timestamp": {
            "type": "string",
            "format": "date-time"
          },
          "duration": {"type": "integer"},
          "exitCode": {"type": ["integer", "null"]}
        }
      }
    },
    "errors": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["errorID", "message", "timestamp"],
        "properties": {
          "errorID": {"type": "string"},
          "message": {"type": "string"},
          "timestamp": {
            "type": "string",
            "format": "date-time"
          },
          "severity": {
            "type": "string",
            "enum": ["warning", "error", "fatal"]
          },
          "context": {"type": "object"}
        }
      }
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
│   ├── active/                 # Active sessions
│   │   └── {session-id}.json
│   └── archive/                # Completed sessions
│       └── 2025-01-05/
│           └── {session-id}.json
├── hooks/                      # Generated hook scripts
│   ├── pre-command.sh
│   ├── post-command.sh
│   ├── file-modified.sh
│   └── session-end.sh
├── cache/                      # Application cache
│   └── document-index.json
└── logs/                       # Debug logs
    └── debug.log
```

## Index Strategies

Since we're using file-based storage, we implement indexing through:

1. **Session Index**: Directory listing with file naming convention for quick access
2. **Document Index**: Cached JSON file with pre-computed search terms
3. **Active Sessions**: Separate directory for O(1) active session queries
4. **Date-based Archive**: Directory structure for efficient historical queries

## Performance Considerations

- **Write Performance**: Atomic writes using temp file + rename
- **Read Performance**: Memory-mapped files for large sessions
- **Concurrency**: File locking for write operations
- **Caching**: In-memory cache with fsnotify invalidation
- **Compression**: Optional gzip for archived sessions
