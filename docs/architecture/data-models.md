# Data Models

## SessionState

**Purpose:** Core session tracking data structure shared between hooks and TUI

**Key Attributes:**
- session_id: string - Unique identifier for Claude Code session
- created_at: time.Time - Session initialization timestamp
- updated_at: time.Time - Last modification timestamp
- session_active: bool - Current session status flag
- agents: []string - Currently executing agents
- agents_history: []AgentExecution - Complete agent execution log

### TypeScript Interface
```typescript
interface SessionState {
  session_id: string;
  created_at: string; // ISO8601
  updated_at: string; // ISO8601
  session_active: boolean;
  agents: string[];
  agents_history: AgentExecution[];
  files: FileOperations;
  tools_used: Record<string, number>;
  errors: ErrorEntry[];
  prompts: PromptEntry[];
  notifications: NotificationEntry[];
}

interface AgentExecution {
  name: string;
  started_at: string; // ISO8601
  completed_at?: string; // ISO8601
}

interface FileOperations {
  new: string[];
  edited: string[];
  read: string[];
}
```

### Relationships
- One-to-many with LogEntry (session has multiple log events)
- Aggregates multiple FileOperations for comprehensive tracking
- Contains multiple AgentExecution records for multi-agent sessions

## DocumentIndex

**Purpose:** Document discovery and navigation for Plan view

**Key Attributes:**
- path: string - Absolute file path to markdown document
- title: string - Extracted document title
- type: DocumentType - Category (PRD, Architecture, Epic, Story)
- modified_at: time.Time - File modification timestamp

### TypeScript Interface
```typescript
interface DocumentIndex {
  path: string;
  title: string;
  type: 'prd' | 'architecture' | 'epic' | 'story';
  modified_at: string; // ISO8601
}
```

### Relationships
- Grouped by DocumentType for hierarchical navigation
- Links to actual markdown files in docs/ directory

## HookEvent

**Purpose:** Individual hook execution tracking for comprehensive logging

**Key Attributes:**
- timestamp: time.Time - Event occurrence time
- session_id: string - Associated session identifier
- hook_name: string - Which hook generated the event
- input_data: interface{} - Raw hook input parameters
- success: bool - Execution success status

### TypeScript Interface
```typescript
interface HookEvent {
  timestamp: string; // ISO8601
  session_id: string;
  hook_name: string;
  input_data: Record<string, any>;
  success: boolean;
}
```

### Relationships
- Many-to-one with SessionState (multiple events per session)
- Aggregated for activity timeline in Observe view
