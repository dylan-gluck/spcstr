package observe

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dylan/spcstr/internal/state"
	"github.com/dylan/spcstr/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
)

type PaneType string

const (
	PaneSessionList PaneType = "sessions"
	PaneDashboard   PaneType = "dashboard"
)

type Model struct {
	width        int
	height       int
	baseStyles   styles.BaseStyles
	paneStyles   PaneStyles
	state        *ObserveState
	stateManager *state.StateManager
	fileWatcher  *fsnotify.Watcher
	initialized  bool
	basePath     string
}

type ObserveState struct {
	sessions       []SessionInfo
	selected       int
	dashboard      *DashboardData
	focusedPane    PaneType
	scrollOffset   int
	dashboardScroll int
	loading        bool
	error          string
	lastUpdate     time.Time
}

type SessionInfo struct {
	ID            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Active        bool
	AgentCount    int
	FileCount     int
	ToolCount     int
	TodoSummary   string
}

type DashboardData struct {
	Session       *state.SessionState
	FormattedData map[string]interface{}
}

type PaneStyles struct {
	ListPane        lipgloss.Style
	DashboardPane   lipgloss.Style
	ListItem        lipgloss.Style
	SelectedItem    lipgloss.Style
	ActiveIndicator lipgloss.Style
	InactiveIndicator lipgloss.Style
	FocusedBorder   lipgloss.Style
	UnfocusedBorder lipgloss.Style
	SectionHeader   lipgloss.Style
	Loading         lipgloss.Style
	StatLabel       lipgloss.Style
	StatValue       lipgloss.Style
}

type sessionsLoadedMsg struct {
	sessions []SessionInfo
	err      error
}

type sessionDataMsg struct {
	dashboard *DashboardData
	err       error
}

type fileChangedMsg struct {
	path string
}

func New() Model {
	// Use .spcstr relative to where spcstr is run (project root)
	basePath := ".spcstr"
	baseStyles := styles.GetDefaultStyles()
	theme := styles.DefaultTheme
	
	return Model{
		baseStyles:   baseStyles,
		paneStyles:   createPaneStyles(theme),
		stateManager: state.NewStateManager(basePath),
		basePath:     basePath,
		state: &ObserveState{
			sessions:    []SessionInfo{},
			selected:    0,
			focusedPane: PaneSessionList,
		},
	}
}

func createPaneStyles(theme styles.Theme) PaneStyles {
	focusedBorder := lipgloss.RoundedBorder()
	unfocusedBorder := lipgloss.NormalBorder()
	
	return PaneStyles{
		ListPane: lipgloss.NewStyle().
			BorderStyle(unfocusedBorder).
			BorderForeground(theme.BorderMuted).
			Padding(1, 1),
		DashboardPane: lipgloss.NewStyle().
			BorderStyle(unfocusedBorder).
			BorderForeground(theme.BorderMuted).
			Padding(1, 1),
		ListItem: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(theme.Text),
		SelectedItem: lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(theme.Secondary).
			Bold(true),
		ActiveIndicator: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true),
		InactiveIndicator: lipgloss.NewStyle().
			Foreground(theme.TextMuted),
		FocusedBorder: lipgloss.NewStyle().
			BorderStyle(focusedBorder).
			BorderForeground(theme.Primary),
		UnfocusedBorder: lipgloss.NewStyle().
			BorderStyle(unfocusedBorder).
			BorderForeground(theme.BorderMuted),
		SectionHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			MarginTop(1).
			MarginBottom(1),
		Loading: lipgloss.NewStyle().
			Foreground(theme.TextMuted).
			Italic(true),
		StatLabel: lipgloss.NewStyle().
			Foreground(theme.TextMuted),
		StatValue: lipgloss.NewStyle().
			Foreground(theme.Text).
			Bold(true),
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadSessions
}

func (m Model) IsInitialized() bool {
	return m.initialized
}

func (m Model) loadSessions() tea.Msg {
	ctx := context.Background()
	sessionIDs, err := m.stateManager.ListSessions(ctx)
	if err != nil {
		return sessionsLoadedMsg{nil, fmt.Errorf("error listing sessions: %w", err)}
	}
	
	var sessions []SessionInfo
	for _, id := range sessionIDs {
		sessionState, err := m.stateManager.LoadState(ctx, id)
		if err != nil {
			continue
		}
		
		info := SessionInfo{
			ID:         id,
			CreatedAt:  sessionState.CreatedAt,
			UpdatedAt:  sessionState.UpdatedAt,
			Active:     sessionState.SessionActive,
			AgentCount: len(sessionState.Agents),
			FileCount:  len(sessionState.Files.New) + len(sessionState.Files.Edited) + len(sessionState.Files.Read),
			ToolCount:  sumToolUsage(sessionState.ToolsUsed),
		}
		
		if sessionState.Todos.Total > 0 {
			info.TodoSummary = fmt.Sprintf("%d/%d", sessionState.Todos.Completed, sessionState.Todos.Total)
		}
		
		sessions = append(sessions, info)
	}
	
	// Sort by created date, newest first
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})
	
	return sessionsLoadedMsg{sessions, nil}
}

func (m Model) loadSessionData(sessionID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		sessionState, err := m.stateManager.LoadState(ctx, sessionID)
		if err != nil {
			return sessionDataMsg{nil, err}
		}
		
		dashboard := &DashboardData{
			Session:       sessionState,
			FormattedData: formatSessionData(sessionState),
		}
		
		return sessionDataMsg{dashboard, nil}
	}
}

func (m Model) watchSessionFile(sessionID string) tea.Cmd {
	return func() tea.Msg {
		// Clean up existing watcher
		if m.fileWatcher != nil {
			m.fileWatcher.Close()
		}
		
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return nil
		}
		
		sessionPath := filepath.Join(m.basePath, "sessions", sessionID, "state.json")
		err = watcher.Add(sessionPath)
		if err != nil {
			watcher.Close()
			return nil
		}
		
		m.fileWatcher = watcher
		
		// Start watching in background
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						// File changed, but we'll handle this in Update
					}
				case _, ok := <-watcher.Errors:
					if !ok {
						return
					}
				}
			}
		}()
		
		return nil
	}
}

func (m Model) checkForFileChanges() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		if m.fileWatcher == nil {
			return nil
		}
		
		select {
		case event := <-m.fileWatcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				return fileChangedMsg{path: event.Name}
			}
		default:
			return nil
		}
		return nil
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 3 // Account for header/footer
		
	case sessionsLoadedMsg:
		if msg.err != nil {
			m.state.error = fmt.Sprintf("Failed to load sessions: %v", msg.err)
		} else {
			m.state.sessions = msg.sessions
			if len(msg.sessions) > 0 && !m.initialized {
				m.initialized = true
				return m, tea.Batch(
					m.loadSessionData(msg.sessions[0].ID),
					m.watchSessionFile(msg.sessions[0].ID),
					m.checkForFileChanges(),
				)
			}
		}
		m.state.loading = false
		
	case sessionDataMsg:
		if msg.err != nil {
			m.state.error = fmt.Sprintf("Failed to load session data: %v", msg.err)
		} else {
			m.state.dashboard = msg.dashboard
			m.state.error = ""
			m.state.lastUpdate = time.Now()
		}
		m.state.loading = false
		
	case fileChangedMsg:
		// Reload current session data when file changes
		if m.state.dashboard != nil && m.state.dashboard.Session != nil {
			cmds := []tea.Cmd{
				m.loadSessionData(m.state.dashboard.Session.SessionID),
				m.checkForFileChanges(),
			}
			return m, tea.Batch(cmds...)
		}
		return m, m.checkForFileChanges()
		
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	
	// Keep checking for file changes
	if m.fileWatcher != nil {
		return m, m.checkForFileChanges()
	}
	
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		if m.state.focusedPane == PaneSessionList {
			m.state.focusedPane = PaneDashboard
		} else {
			m.state.focusedPane = PaneSessionList
		}
		
	case "up", "k":
		if m.state.focusedPane == PaneSessionList && m.state.selected > 0 {
			m.state.selected--
			if len(m.state.sessions) > m.state.selected {
				m.state.loading = true
				return m, tea.Batch(
					m.loadSessionData(m.state.sessions[m.state.selected].ID),
					m.watchSessionFile(m.state.sessions[m.state.selected].ID),
				)
			}
		} else if m.state.focusedPane == PaneDashboard && m.state.dashboardScroll > 0 {
			m.state.dashboardScroll--
		}
		
	case "down", "j":
		if m.state.focusedPane == PaneSessionList && m.state.selected < len(m.state.sessions)-1 {
			m.state.selected++
			if len(m.state.sessions) > m.state.selected {
				m.state.loading = true
				return m, tea.Batch(
					m.loadSessionData(m.state.sessions[m.state.selected].ID),
					m.watchSessionFile(m.state.sessions[m.state.selected].ID),
				)
			}
		} else if m.state.focusedPane == PaneDashboard {
			m.state.dashboardScroll++
		}
		
	case "enter":
		if m.state.focusedPane == PaneSessionList && len(m.state.sessions) > m.state.selected {
			m.state.loading = true
			return m, tea.Batch(
				m.loadSessionData(m.state.sessions[m.state.selected].ID),
				m.watchSessionFile(m.state.sessions[m.state.selected].ID),
			)
		}
		
	case "r":
		// Manual refresh
		if m.state.dashboard != nil && m.state.dashboard.Session != nil {
			m.state.loading = true
			return m, m.loadSessionData(m.state.dashboard.Session.SessionID)
		}
		
	case "R":
		// Refresh session list
		m.state.loading = true
		return m, m.loadSessions
		
	case "home", "g":
		if m.state.focusedPane == PaneSessionList && len(m.state.sessions) > 0 {
			m.state.selected = 0
			m.state.loading = true
			return m, tea.Batch(
				m.loadSessionData(m.state.sessions[0].ID),
				m.watchSessionFile(m.state.sessions[0].ID),
			)
		} else if m.state.focusedPane == PaneDashboard {
			m.state.dashboardScroll = 0
		}
		
	case "end", "G":
		if m.state.focusedPane == PaneSessionList && len(m.state.sessions) > 0 {
			m.state.selected = len(m.state.sessions) - 1
			m.state.loading = true
			return m, tea.Batch(
				m.loadSessionData(m.state.sessions[m.state.selected].ID),
				m.watchSessionFile(m.state.sessions[m.state.selected].ID),
			)
		}
	}
	
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Waiting for window size..."
	}
	
	// Simple division: left pane gets 1/3 of width
	leftPaneWidth := (m.width / 3) - 4
	rightPaneWidth := m.width - (m.width / 3) - 4
	
	leftPane := m.renderSessionList(leftPaneWidth, m.height-2)
	rightPane := m.renderDashboard(rightPaneWidth, m.height-2)
	
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

func (m Model) renderSessionList(width, height int) string {
	var listItems []string
	
	header := m.paneStyles.SectionHeader.Render("── SESSIONS ──")
	listItems = append(listItems, header)
	
	if m.state.loading && len(m.state.sessions) == 0 {
		listItems = append(listItems, m.paneStyles.Loading.Render("Loading sessions..."))
	} else if len(m.state.sessions) == 0 {
		listItems = append(listItems, m.baseStyles.TextMuted.Render("No sessions found"))
		listItems = append(listItems, m.baseStyles.TextMuted.Render(fmt.Sprintf("Path: %s", m.basePath)))
		if m.state.error != "" {
			listItems = append(listItems, m.baseStyles.Error.Render(fmt.Sprintf("Error: %s", m.state.error)))
		}
	} else {
		for i, session := range m.state.sessions {
			statusIcon := "○"
			statusStyle := m.paneStyles.InactiveIndicator
			if session.Active {
				statusIcon = "●"
				statusStyle = m.paneStyles.ActiveIndicator
			}
			
			// Format session info
			sessionID := session.ID
			if len(sessionID) > 20 {
				sessionID = sessionID[:8] + "..." + sessionID[len(sessionID)-8:]
			}
			
			timeAgo := formatTimeAgo(session.UpdatedAt)
			
			item := fmt.Sprintf("%s %s", statusStyle.Render(statusIcon), sessionID)
			subInfo := fmt.Sprintf("   %s", timeAgo)
			
			if session.TodoSummary != "" {
				subInfo += fmt.Sprintf(" | Tasks: %s", session.TodoSummary)
			}
			
			if i == m.state.selected {
				item = m.paneStyles.SelectedItem.Render("▸ " + item)
				subInfo = m.paneStyles.SelectedItem.Render(subInfo)
			} else {
				item = m.paneStyles.ListItem.Render(item)
				subInfo = m.baseStyles.TextMuted.Render(subInfo)
			}
			
			listItems = append(listItems, item)
			listItems = append(listItems, subInfo)
		}
	}
	
	content := strings.Join(listItems, "\n")
	
	paneStyle := m.paneStyles.ListPane
	if m.state.focusedPane == PaneSessionList {
		paneStyle = paneStyle.BorderStyle(m.paneStyles.FocusedBorder.GetBorderStyle()).
			BorderForeground(m.paneStyles.FocusedBorder.GetBorderTopForeground())
	}
	
	return paneStyle.
		Width(width).
		Height(height).
		MaxWidth(width).
		MaxHeight(height).
		Render(content)
}

func (m Model) renderDashboard(width, height int) string {
	var content string
	
	if m.state.loading && m.state.dashboard != nil {
		// Show existing content while loading
		content = m.formatDashboardContent()
	} else if m.state.loading {
		content = m.paneStyles.Loading.Render("Loading session data...")
	} else if m.state.error != "" {
		content = m.baseStyles.Error.Render(m.state.error)
	} else if m.state.dashboard == nil {
		content = m.baseStyles.TextMuted.Render("Select a session to view details")
	} else {
		content = m.formatDashboardContent()
	}
	
	// Handle scrolling
	lines := strings.Split(content, "\n")
	visibleHeight := height - 2
	if m.state.dashboardScroll > len(lines)-visibleHeight {
		m.state.dashboardScroll = len(lines) - visibleHeight
	}
	if m.state.dashboardScroll < 0 {
		m.state.dashboardScroll = 0
	}
	
	if m.state.dashboardScroll < len(lines) {
		endLine := m.state.dashboardScroll + visibleHeight
		if endLine > len(lines) {
			endLine = len(lines)
		}
		visibleLines := lines[m.state.dashboardScroll:endLine]
		content = strings.Join(visibleLines, "\n")
	}
	
	paneStyle := m.paneStyles.DashboardPane
	if m.state.focusedPane == PaneDashboard {
		paneStyle = paneStyle.BorderStyle(m.paneStyles.FocusedBorder.GetBorderStyle()).
			BorderForeground(m.paneStyles.FocusedBorder.GetBorderTopForeground())
	}
	
	return paneStyle.
		Width(width).
		Height(height).
		MaxWidth(width).
		MaxHeight(height).
		Render(content)
}

func (m Model) formatDashboardContent() string {
	if m.state.dashboard == nil || m.state.dashboard.Session == nil {
		return m.baseStyles.TextMuted.Render("No session data available")
	}
	
	session := m.state.dashboard.Session
	var sections []string
	
	// Session Header
	headerText := fmt.Sprintf("SESSION: %s", session.SessionID[:16])
	if session.SessionActive {
		headerText += " " + m.paneStyles.ActiveIndicator.Render("[ACTIVE]")
	} else {
		headerText += " " + m.paneStyles.InactiveIndicator.Render("[INACTIVE]")
	}
	sections = append(sections, m.paneStyles.SectionHeader.Render(headerText))
	
	// Timing info
	duration := session.UpdatedAt.Sub(session.CreatedAt)
	lastUpdate := formatTimeAgo(session.UpdatedAt)
	sections = append(sections, fmt.Sprintf("%s %s | %s %s | %s %s",
		m.paneStyles.StatLabel.Render("Started:"),
		m.paneStyles.StatValue.Render(session.CreatedAt.Format("15:04:05")),
		m.paneStyles.StatLabel.Render("Duration:"),
		m.paneStyles.StatValue.Render(formatDuration(duration)),
		m.paneStyles.StatLabel.Render("Updated:"),
		m.paneStyles.StatValue.Render(lastUpdate),
	))
	
	// Agents Section
	if len(session.Agents) > 0 || len(session.AgentsHistory) > 0 {
		sections = append(sections, "")
		sections = append(sections, m.paneStyles.SectionHeader.Render("── AGENTS ──"))
		
		if len(session.Agents) > 0 {
			sections = append(sections, fmt.Sprintf("%s %s",
				m.paneStyles.StatLabel.Render("Active:"),
				m.paneStyles.StatValue.Render(strings.Join(session.Agents, ", ")),
			))
		}
		
		if len(session.AgentsHistory) > 0 {
			sections = append(sections, fmt.Sprintf("%s %d agents executed",
				m.paneStyles.StatLabel.Render("History:"),
				len(session.AgentsHistory),
			))
			// Show last 3 agents
			start := len(session.AgentsHistory) - 3
			if start < 0 {
				start = 0
			}
			for _, agent := range session.AgentsHistory[start:] {
				dur := "running"
				if agent.CompletedAt != nil {
					dur = formatDuration(agent.CompletedAt.Sub(agent.StartedAt))
				}
				sections = append(sections, fmt.Sprintf("  • %s (%s)",
					agent.Name, dur,
				))
			}
		}
	}
	
	// Files Section
	totalFiles := len(session.Files.New) + len(session.Files.Edited) + len(session.Files.Read)
	if totalFiles > 0 {
		sections = append(sections, "")
		sections = append(sections, m.paneStyles.SectionHeader.Render("── FILES ──"))
		
		fileStats := fmt.Sprintf("%s %s | %s %s | %s %s",
			m.paneStyles.StatLabel.Render("New:"),
			m.paneStyles.StatValue.Render(fmt.Sprintf("%d", len(session.Files.New))),
			m.paneStyles.StatLabel.Render("Edited:"),
			m.paneStyles.StatValue.Render(fmt.Sprintf("%d", len(session.Files.Edited))),
			m.paneStyles.StatLabel.Render("Read:"),
			m.paneStyles.StatValue.Render(fmt.Sprintf("%d", len(session.Files.Read))),
		)
		sections = append(sections, fileStats)
		
		// Show recent files
		recentFiles := []string{}
		for i := len(session.Files.New) - 1; i >= 0 && len(recentFiles) < 2; i-- {
			recentFiles = append(recentFiles, fmt.Sprintf("  + %s", filepath.Base(session.Files.New[i])))
		}
		for i := len(session.Files.Edited) - 1; i >= 0 && len(recentFiles) < 4; i-- {
			recentFiles = append(recentFiles, fmt.Sprintf("  ~ %s", filepath.Base(session.Files.Edited[i])))
		}
		sections = append(sections, recentFiles...)
	}
	
	// Tools Section
	if len(session.ToolsUsed) > 0 {
		sections = append(sections, "")
		sections = append(sections, m.paneStyles.SectionHeader.Render("── TOOLS USED ──"))
		
		// Sort tools by usage
		type toolUsage struct {
			name  string
			count int
		}
		var tools []toolUsage
		for name, count := range session.ToolsUsed {
			tools = append(tools, toolUsage{name, count})
		}
		sort.Slice(tools, func(i, j int) bool {
			return tools[i].count > tools[j].count
		})
		
		for i, tool := range tools {
			if i >= 5 {
				break // Show top 5
			}
			sections = append(sections, fmt.Sprintf("  %s: %s",
				m.paneStyles.StatLabel.Render(tool.name),
				m.paneStyles.StatValue.Render(fmt.Sprintf("%d", tool.count)),
			))
		}
	}
	
	// Tasks Section
	if session.Todos.Total > 0 {
		sections = append(sections, "")
		sections = append(sections, m.paneStyles.SectionHeader.Render("── TASKS ──"))
		
		taskStats := fmt.Sprintf("%s %s | %s %s | %s %s | %s %s",
			m.paneStyles.StatLabel.Render("Total:"),
			m.paneStyles.StatValue.Render(fmt.Sprintf("%d", session.Todos.Total)),
			m.paneStyles.StatLabel.Render("Done:"),
			m.paneStyles.StatValue.Render(fmt.Sprintf("%d", session.Todos.Completed)),
			m.paneStyles.StatLabel.Render("Active:"),
			m.paneStyles.StatValue.Render(fmt.Sprintf("%d", session.Todos.InProgress)),
			m.paneStyles.StatLabel.Render("Pending:"),
			m.paneStyles.StatValue.Render(fmt.Sprintf("%d", session.Todos.Pending)),
		)
		sections = append(sections, taskStats)
		
		if len(session.Todos.Recent) > 0 {
			sections = append(sections, m.baseStyles.TextMuted.Render("Recent:"))
			for i, todo := range session.Todos.Recent {
				if i >= 3 {
					break
				}
				statusIcon := "○"
				if todo.Status == "completed" {
					statusIcon = "✓"
				} else if todo.Status == "in_progress" {
					statusIcon = "◐"
				}
				
				content := todo.Content
				if len(content) > 50 {
					content = content[:47] + "..."
				}
				sections = append(sections, fmt.Sprintf("  %s %s", statusIcon, content))
			}
		}
	}
	
	// Activity Section
	if len(session.Prompts) > 0 || len(session.Notifications) > 0 {
		sections = append(sections, "")
		sections = append(sections, m.paneStyles.SectionHeader.Render("── RECENT ACTIVITY ──"))
		
		// Combine and sort by time
		type activity struct {
			time    time.Time
			content string
			typ     string
		}
		var activities []activity
		
		for _, prompt := range session.Prompts {
			content := prompt.Prompt
			if len(content) > 60 {
				content = content[:57] + "..."
			}
			activities = append(activities, activity{
				time:    prompt.Timestamp,
				content: content,
				typ:     "prompt",
			})
		}
		
		for _, notif := range session.Notifications {
			activities = append(activities, activity{
				time:    notif.Timestamp,
				content: notif.Message,
				typ:     notif.Type,
			})
		}
		
		sort.Slice(activities, func(i, j int) bool {
			return activities[i].time.After(activities[j].time)
		})
		
		// Show last 5 activities
		for i, act := range activities {
			if i >= 5 {
				break
			}
			icon := "•"
			if act.typ == "prompt" {
				icon = ">"
			} else if act.typ == "hook" {
				icon = "⚡"
			}
			sections = append(sections, fmt.Sprintf("  %s %s %s",
				icon,
				m.baseStyles.TextMuted.Render(act.time.Format("15:04")),
				act.content,
			))
		}
	}
	
	// Last refresh indicator
	if m.state.lastUpdate.After(time.Time{}) {
		sections = append(sections, "")
		sections = append(sections, m.baseStyles.TextMuted.Render(
			fmt.Sprintf("Last refresh: %s", formatTimeAgo(m.state.lastUpdate)),
		))
	}
	
	return strings.Join(sections, "\n")
}

// Helper functions

func sumToolUsage(tools map[string]int) int {
	total := 0
	for _, count := range tools {
		total += count
	}
	return total
}

func formatSessionData(session *state.SessionState) map[string]interface{} {
	data := make(map[string]interface{})
	
	data["id"] = session.SessionID
	data["active"] = session.SessionActive
	data["created"] = session.CreatedAt
	data["updated"] = session.UpdatedAt
	data["agents"] = session.Agents
	data["files"] = session.Files
	data["tools"] = session.ToolsUsed
	data["todos"] = session.Todos
	
	return data
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	
	if duration < time.Minute {
		return fmt.Sprintf("%ds ago", int(duration.Seconds()))
	} else if duration < time.Hour {
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	} else {
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	} else {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
}