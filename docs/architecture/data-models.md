# Data Models

## SessionState Model
**Purpose:** Represents the complete state of a Claude Code session with integrated tracking

**Key Attributes:**
- SessionID: string - Unique identifier for the session (format: sess_{uuid})
- CreatedAt: timestamp - Session creation time
- UpdatedAt: timestamp - Last state update time
- Source: enum(startup, resume, clear) - Session initiation source
- ProjectPath: string - Absolute path to project directory
- Status: enum(active, completed, error) - Current session state
- Agents: []string - Currently active agent names
- AgentsHistory: []AgentHistoryEntry - Complete agent execution history
- Files: FileOperations - Categorized file operations (new/edited/read)
- ToolsUsed: map[string]int - Tool name to usage count mapping
- Errors: []ErrorEntry - Structured error tracking
- Modified: bool - Internal dirty flag (not persisted)

**Relationships:**
- Contains AgentHistoryEntry records (embedded)
- Contains FileOperations structure (embedded)
- Contains ErrorEntry records (embedded)

## AgentHistoryEntry Model
**Purpose:** Tracks complete history of agent executions within a session

**Key Attributes:**
- Name: string - Agent name/identifier
- StartedAt: timestamp - When agent execution began
- EndedAt: timestamp - When agent execution completed (nullable for active agents)

**Relationships:**
- Embedded within SessionState (many-to-1)

## FileOperations Model
**Purpose:** Categorizes all file operations performed during session

**Key Attributes:**
- New: []string - Absolute paths of created files
- Edited: []string - Absolute paths of modified files  
- Read: []string - Absolute paths of files read

**Relationships:**
- Embedded within SessionState (1-to-1)

## ErrorEntry Model
**Purpose:** Structured error tracking with context

**Key Attributes:**
- Timestamp: timestamp - When error occurred
- Hook: string - Hook name where error originated
- Message: string - Error description

**Relationships:**
- Embedded within SessionState (many-to-1)

## Tool Usage Tracking
**Purpose:** Maintains count of tool invocations during session

**Structure:** map[string]int - Tool name mapped to usage count

**Common Tools Tracked:**
- Read, Write, Edit, MultiEdit - File operations
- Bash, BashOutput, KillBash - Command execution
- Task - Agent invocations  
- Grep, Glob - Search operations
- WebSearch, WebFetch - Web operations

**Relationships:**
- Embedded within SessionState as ToolsUsed field



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
