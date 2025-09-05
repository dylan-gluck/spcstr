# Frontend Architecture

## Component Architecture

### Component Organization
```
internal/tui/
├── app/                        # Main application controller
│   └── app.go                 # Bubbletea app initialization
├── components/                 # Reusable UI components
│   ├── header/                # Header bar component
│   ├── footer/                # Footer/status bar component
│   ├── list/                  # Generic list component
│   └── dashboard/             # Session dashboard component
├── views/                     # Main view implementations
│   ├── plan/                  # Plan view (document browser)
│   └── observe/               # Observe view (session monitor)
├── styles/                    # Lipgloss styling definitions
│   └── theme.go              # Color scheme and layout styles
└── messages/                  # Custom Bubbletea messages
    └── events.go             # File change and update events
```

### Component Template
```go
// Standard Bubbletea component pattern
type Component struct {
    width  int
    height int
    styles lipgloss.Style
    // Component-specific fields
}

func NewComponent() Component {
    return Component{
        styles: styles.DefaultComponent(),
    }
}

func (c Component) Init() tea.Cmd {
    return nil
}

func (c Component) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        c.width = msg.Width
        c.height = msg.Height
    case tea.KeyMsg:
        return c.handleKeys(msg)
    }
    return c, nil
}

func (c Component) View() string {
    return c.styles.Render("Component content")
}
```

## State Management Architecture

### State Structure
```go
// TUI application state (separate from session state)
type AppState struct {
    currentView    ViewType      // plan, observe
    planState     *PlanState    // Plan view state
    observeState  *ObserveState // Observe view state
    windowSize    tea.WindowSizeMsg
    initialized   bool
}

type PlanState struct {
    documents     []DocumentIndex
    selected      int
    content       string
    focusedPane   PaneType // list, content
}

type ObserveState struct {
    sessions      []SessionState
    selected      int
    dashboard     DashboardData
    lastUpdate    time.Time
}
```

### State Management Patterns
- **Centralized App State:** Single AppState struct manages all TUI state
- **View-Specific State:** Each view maintains its own subset of state
- **Immutable Updates:** State changes create new state objects
- **Event-Driven Updates:** File watcher events trigger state refreshes
- **Local State Only:** No persistence of TUI state between runs

## Routing Architecture

### Route Organization
```
TUI Navigation Routes (key bindings):
├── Global Keys
│   ├── 'p' → Plan View
│   ├── 'o' → Observe View  
│   └── 'q' → Quit Application
├── Plan View Keys
│   ├── 'tab' → Switch Pane Focus
│   ├── '↑/↓' → Navigate Document List
│   ├── 'enter' → Select Document
│   └── 's/w/c' → Switch Modes (minimal for MVP)
└── Observe View Keys
    ├── '↑/↓' → Navigate Session List
    ├── 'enter' → Select Session
    └── 'r' → Manual Refresh
```

### Navigation Pattern
```go
// Key handler routing pattern
func (a App) handleGlobalKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "p":
        a.currentView = PlanView
        return a, nil
    case "o":
        a.currentView = ObserveView
        return a, tea.Cmd(loadSessions)
    case "q":
        return a, tea.Quit
    }
    
    // Route to view-specific handlers
    switch a.currentView {
    case PlanView:
        return a.planView.Update(msg)
    case ObserveView:
        return a.observeView.Update(msg)
    }
    
    return a, nil
}
```

## TUI Services Layer

### File Watching Service
```go
// File watcher integration for real-time updates
type FileWatcherService struct {
    watcher   *fsnotify.Watcher
    eventChan chan FileChangeEvent
}

func (f *FileWatcherService) WatchStateFiles() tea.Cmd {
    return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
        select {
        case event := <-f.eventChan:
            return FileChangeMsg{Event: event}
        default:
            return nil
        }
    })
}

type FileChangeMsg struct {
    Event FileChangeEvent
}
```

### State Service Integration
```go
// Service layer for state management
type StateService struct {
    stateManager *state.Manager
}

func (s *StateService) LoadSessionList() tea.Cmd {
    return func() tea.Msg {
        sessions, err := s.stateManager.ListSessions()
        if err != nil {
            return ErrorMsg{Err: err}
        }
        return SessionListMsg{Sessions: sessions}
    }
}

func (s *StateService) LoadSessionDetails(id string) tea.Cmd {
    return func() tea.Msg {
        session, err := s.stateManager.LoadState(id)
        if err != nil {
            return ErrorMsg{Err: err}
        }
        return SessionDetailMsg{Session: session}
    }
}
```
