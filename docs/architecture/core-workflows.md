# Core Workflows

## Session Creation and Tracking Workflow
```mermaid
sequenceDiagram
    participant CC as Claude Code
    participant HS as Hook Script
    participant FS as File System
    participant FW as File Watcher
    participant EB as Event Bus
    participant SM as Session Manager
    participant TUI as TUI Display
    
    CC->>HS: Trigger session start hook
    HS->>HS: Generate session ID
    HS->>HS: Validate paths
    HS->>FS: Create session.json
    FS-->>FW: File created event
    FW->>EB: Publish NewSession event
    EB->>SM: Handle NewSession
    SM->>SM: Load session data
    SM->>EB: Publish SessionUpdated
    EB->>TUI: Update session list
    TUI->>TUI: Render new session
    
    loop During Session
        CC->>HS: Trigger tool/file hooks
        HS->>FS: Update session.json
        FS-->>FW: File modified event
        FW->>EB: Publish SessionModified
        EB->>SM: Update in-memory state
        SM->>EB: Publish updates
        EB->>TUI: Real-time updates
    end
    
    CC->>HS: Trigger session end hook
    HS->>FS: Finalize session.json
    FW->>EB: Publish SessionCompleted
    EB->>SM: Mark session complete
    EB->>TUI: Update status
```

## Plan View Document Discovery Workflow
```mermaid
sequenceDiagram
    participant User as User
    participant TUI as TUI Engine
    participant PV as Plan View
    participant DI as Document Indexer
    participant FS as File System
    participant MR as Markdown Renderer
    
    User->>TUI: Press 'p' (Plan View)
    TUI->>PV: SwitchToView()
    PV->>DI: RequestIndex()
    
    alt Index not cached
        DI->>FS: ScanDirectory(docs/)
        FS-->>DI: Return file list
        DI->>DI: ClassifyDocuments()
        DI->>DI: ExtractMetadata()
        DI->>PV: Return index
    else Index cached
        DI->>PV: Return cached index
    end
    
    PV->>PV: RenderDocumentList()
    PV->>TUI: Display view
    
    User->>PV: Navigate with arrows
    PV->>PV: UpdateSelection()
    PV->>FS: ReadDocument()
    FS-->>PV: Document content
    PV->>MR: RenderMarkdown()
    MR-->>PV: Formatted output
    PV->>TUI: Update preview pane
    
    User->>PV: Press '/' (search)
    PV->>PV: ShowSearchInput()
    User->>PV: Type search query
    PV->>DI: FuzzySearch(query)
    DI-->>PV: Filtered results
    PV->>TUI: Update list
```

## Real-time Dashboard Update Workflow
```mermaid
sequenceDiagram
    participant HS as Hook Script
    participant FS as File System
    participant FW as File Watcher
    participant EB as Event Bus
    participant SM as Session Manager
    participant OV as Observe View
    participant DR as Dashboard Renderer
    
    HS->>FS: Update task status
    FW->>EB: FileModified event
    EB->>SM: UpdateSession()
    SM->>SM: Parse changes
    SM->>EB: TaskUpdated event
    
    EB->>OV: HandleTaskUpdate()
    OV->>DR: UpdateTaskSection()
    DR->>DR: Calculate progress
    DR->>DR: Format display
    DR->>OV: Rendered section
    OV->>OV: Merge updates
    OV->>OV: TUI refresh
    
    Note over OV: Sub-100ms update cycle
    
    HS->>FS: Add error entry
    FW->>EB: FileModified event
    EB->>SM: UpdateSession()
    SM->>EB: ErrorAdded event
    EB->>OV: HandleError()
    OV->>DR: UpdateErrorLog()
    DR->>OV: Rendered errors
    OV->>OV: Flash error indicator
```

## Session Recovery on Restart Workflow
```mermaid
sequenceDiagram
    participant User as User
    participant CLI as CLI
    participant CM as Config Manager
    participant SM as Session Manager
    participant PL as Persistence Layer
    participant EB as Event Bus
    participant TUI as TUI Engine
    
    User->>CLI: spcstr run
    CLI->>CM: LoadConfig()
    CM->>CM: Merge global/project
    CM-->>CLI: Configuration
    
    CLI->>SM: Initialize()
    SM->>PL: ListSessions()
    PL->>PL: Read .spcstr/sessions/
    PL-->>SM: Session files
    
    loop For each session
        SM->>PL: LoadSession()
        PL->>PL: Parse JSON
        PL-->>SM: Session data
        SM->>SM: Rebuild state
        SM->>SM: Check if active
    end
    
    SM->>EB: PublishAllSessions()
    CLI->>TUI: Start()
    TUI->>EB: Subscribe()
    EB-->>TUI: Session events
    TUI->>TUI: Render UI
    TUI-->>User: Display recovered state
```

## Error Handling Flow
```mermaid
sequenceDiagram
    participant Component as Any Component
    participant Err as Error
    participant Log as Logger
    participant EB as Event Bus
    participant TUI as TUI
    participant FS as File System
    
    Component->>Err: Operation fails
    Err->>Log: Log error details
    Log->>FS: Write to debug log
    
    alt Recoverable Error
        Err->>Component: Return error
        Component->>Component: Retry/fallback
        Component->>EB: Publish warning
        EB->>TUI: Display warning
        TUI->>TUI: Show in status bar
    else Fatal Error
        Err->>EB: Publish fatal error
        EB->>TUI: Display error modal
        TUI->>TUI: Offer recovery options
        TUI-->>Component: User choice
    end
    
    Note over FS: Errors never block hooks
```
