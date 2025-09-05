# Data Models

## Session Model
**Purpose:** Represents a complete Claude Code session from start to finish

**Key Attributes:**
- SessionID: string - Unique identifier for the session
- Status: enum(active, completed, error) - Current session state
- StartTime: timestamp - When session began
- EndTime: timestamp - When session completed (nullable)
- ProjectPath: string - Absolute path to project directory
- Metadata: map[string]string - Additional session context

**Relationships:**
- Has many Agents (1-to-many)
- Has many Tasks (1-to-many)
- Has many FileOperations (1-to-many)
- Has many ToolExecutions (1-to-many)
- Has many Errors (1-to-many)

## Agent Model
**Purpose:** Represents an AI agent participating in the session

**Key Attributes:**
- AgentID: string - Unique identifier
- Name: string - Agent display name
- Type: string - Agent type (e.g., "dev", "architect")
- Status: enum(idle, active, error) - Current agent state
- LastActivity: timestamp - Last action timestamp
- SessionID: string - Parent session reference

**Relationships:**
- Belongs to Session (many-to-1)
- Has many Tasks (1-to-many)

## Task Model
**Purpose:** Tracks individual tasks within a session

**Key Attributes:**
- TaskID: string - Unique identifier
- Content: string - Task description
- Status: enum(pending, in_progress, completed) - Task state
- CreatedAt: timestamp - Task creation time
- CompletedAt: timestamp - Task completion time (nullable)
- AgentID: string - Assigned agent reference

**Relationships:**
- Belongs to Session (many-to-1)
- Belongs to Agent (many-to-1)

## FileOperation Model
**Purpose:** Records file system operations performed during session

**Key Attributes:**
- OperationID: string - Unique identifier
- FilePath: string - Absolute file path
- Operation: enum(created, edited, read, deleted) - Operation type
- Timestamp: timestamp - When operation occurred
- SessionID: string - Parent session reference
- LineCount: int - Number of lines affected (optional)

**Relationships:**
- Belongs to Session (many-to-1)

## ToolExecution Model
**Purpose:** Tracks tool/command executions

**Key Attributes:**
- ExecutionID: string - Unique identifier
- ToolName: string - Name of tool executed
- Command: string - Full command or parameters
- Timestamp: timestamp - Execution time
- Duration: int - Execution duration in milliseconds
- ExitCode: int - Command exit code (nullable)
- SessionID: string - Parent session reference

**Relationships:**
- Belongs to Session (many-to-1)

## PlanDocument Model
**Purpose:** Represents indexed planning documents (PRD, Architecture, etc.)

**Key Attributes:**
- DocumentID: string - Unique identifier
- FilePath: string - Absolute path to document
- Type: enum(prd, architecture, epic, story, config) - Document classification
- Title: string - Extracted document title
- LastModified: timestamp - File modification time
- Content: string - Cached content for search
- Metadata: map[string]string - Extracted frontmatter or metadata

**Relationships:**
- Independent entity (used by Plan View)

## Configuration Model
**Purpose:** Stores application and project configuration

**Key Attributes:**
- ConfigID: string - Unique identifier
- Scope: enum(global, project) - Configuration scope
- HooksPath: string - Path to hook scripts
- SessionsPath: string - Path to session storage
- DocsPath: []string - Paths to scan for documents
- UIPreferences: map[string]interface{} - UI settings
- Version: string - Config schema version

**Relationships:**
- Independent entity (singleton per scope)
