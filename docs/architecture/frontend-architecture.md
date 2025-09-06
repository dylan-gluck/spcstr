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
// Standard Bubbletea component pattern with lessons learned
type Component struct {
    width  int
    height int
    styles ComponentStyles  // Use structured styles, not single style
    // Component-specific fields
}

type ComponentStyles struct {
    Container      lipgloss.Style
    FocusedBorder  lipgloss.Style
    UnfocusedBorder lipgloss.Style
}

func NewComponent() Component {
    return Component{
        styles: createComponentStyles(theme),
    }
}

func (c Component) Init() tea.Cmd {
    return nil
}

func (c Component) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        c.width = msg.Width
        c.height = msg.Height - 3  // Account for header/footer
    case tea.KeyMsg:
        return c.handleKeys(msg)
    }
    return c, nil
}

func (c Component) View() string {
    if c.width == 0 || c.height == 0 {
        return ""  // Don't render without dimensions
    }
    
    // LESSON LEARNED: Use MaxWidth/MaxHeight to prevent overflow
    return c.styles.Container.
        Width(c.width).
        Height(c.height).
        MaxWidth(c.width).    // Critical for preventing overflow
        MaxHeight(c.height).  // Critical for preventing overflow
        Render("Component content")
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

## Layout Best Practices (Lessons Learned)

### Critical Layout Patterns

#### 1. Simple Division Over Complex Calculations
```go
// ❌ BAD: Complex percentage calculations with conditionals
leftPaneWidth := (m.width * 30) / 100
if leftPaneWidth < 20 {
    leftPaneWidth = 20
}
rightPaneWidth := m.width - leftPaneWidth - 2
if rightPaneWidth < 40 {
    rightPaneWidth = 40
}

// ✅ GOOD: Simple division
leftPaneWidth := (m.width / 3) - 4  // Account for borders/padding
rightPaneWidth := m.width - (m.width / 3) - 4
```

#### 2. Always Use MaxWidth/MaxHeight
```go
// ❌ BAD: Only setting Width/Height
return style.
    Width(width).
    Height(height).
    Render(content)

// ✅ GOOD: MaxWidth/MaxHeight prevent overflow
return style.
    Width(width).
    Height(height).
    MaxWidth(width).    // CRITICAL: Prevents horizontal overflow
    MaxHeight(height).  // CRITICAL: Prevents vertical overflow
    Render(content)
```

#### 3. Account for Chrome in Calculations
```go
// ❌ BAD: Using raw terminal dimensions
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

// ✅ GOOD: Account for UI chrome
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height - 3  // Account for header/footer
```

#### 4. Use lipgloss.Place for Complex Layouts
```go
// For complex content positioning, use lipgloss.Place
content := lipgloss.Place(
    width, 
    height,
    lipgloss.Left,    // horizontal alignment
    lipgloss.Top,     // vertical alignment  
    renderedContent,
)
```

### Common Pitfalls to Avoid

1. **Don't Fight lipgloss:** Work with its layout system, not against it
2. **Avoid Manual Pixel Counting:** Let lipgloss handle border/padding calculations
3. **No Complex Minimum Checks:** Simple divisions work better than conditional minimums
4. **Test Multiple Terminal Sizes:** Always test with 80x24, 120x40, and edge cases

### Recommended Pane Layout Pattern
```go
func (m Model) View() string {
    if m.width == 0 || m.height == 0 {
        return ""
    }
    
    // Simple division for multi-pane layouts
    leftPaneWidth := (m.width / 3) - 4  // 1/3 for list
    rightPaneWidth := m.width - (m.width / 3) - 4  // 2/3 for content
    
    leftPane := m.renderPane(leftPaneWidth, m.height-2, leftContent)
    rightPane := m.renderPane(rightPaneWidth, m.height-2, rightContent)
    
    return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

func (m Model) renderPane(width, height int, content string) string {
    return m.paneStyle.
        Width(width).
        Height(height).
        MaxWidth(width).    // Always include
        MaxHeight(height).  // Always include
        Render(content)
}
```
