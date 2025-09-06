package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dylan/spcstr/internal/tui/components/footer"
	"github.com/dylan/spcstr/internal/tui/components/header"
	"github.com/dylan/spcstr/internal/tui/styles"
	"github.com/dylan/spcstr/internal/tui/views/observe"
	"github.com/dylan/spcstr/internal/tui/views/plan"
)

type ViewType string

const (
	ViewPlan    ViewType = "plan"
	ViewObserve ViewType = "observe"
)

type AppState struct {
	currentView    ViewType
	planView       tea.Model
	observeView    tea.Model
	header         tea.Model
	footer         tea.Model
	windowWidth    int
	windowHeight   int
	initialized    bool
	projectPath    string
	lastSwitchTime time.Time
}

type App struct {
	state *AppState
}

func New() *App {
	return &App{
		state: &AppState{
			currentView: ViewPlan,
			initialized: false,
		},
	}
}

func (a *App) Init() tea.Cmd {
	a.checkInitialization()

	if a.state.initialized {
		a.initializeViews()
	}

	return tea.EnterAltScreen
}

func (a *App) checkInitialization() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	spcstrDir := filepath.Join(cwd, ".spcstr")
	if info, err := os.Stat(spcstrDir); err == nil && info.IsDir() {
		a.state.initialized = true
		a.state.projectPath = cwd
	}
}

func (a *App) initializeViews() {
	// Initialize header with size
	headerModel := header.New()
	headerModel.SetSessionStatus("active")
	if a.state.windowWidth > 0 {
		headerModel.Update(tea.WindowSizeMsg{Width: a.state.windowWidth, Height: a.state.windowHeight})
	}
	a.state.header = headerModel

	// Initialize footer with size
	footerModel := footer.New()
	footerModel.UpdateForView(string(a.state.currentView))
	if a.state.windowWidth > 0 {
		footerModel.Update(tea.WindowSizeMsg{Width: a.state.windowWidth, Height: a.state.windowHeight})
	}
	a.state.footer = footerModel

	// Initialize views with size
	planModel := plan.New()
	if a.state.windowWidth > 0 {
		planModel.Update(tea.WindowSizeMsg{Width: a.state.windowWidth, Height: a.state.windowHeight})
	}
	a.state.planView = planModel

	observeModel := observe.New()
	if a.state.windowWidth > 0 {
		observeModel.Update(tea.WindowSizeMsg{Width: a.state.windowWidth, Height: a.state.windowHeight})
	}
	a.state.observeView = observeModel
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.state.windowWidth = msg.Width
		a.state.windowHeight = msg.Height

		// Initialize views on first size message if not already done
		if a.state.initialized && a.state.header == nil {
			a.initializeViews()
		}

		return a, a.propagateSizeUpdate(msg)

	case tea.KeyMsg:
		return a.handleGlobalKeys(msg)
	}

	if !a.state.initialized {
		return a, nil
	}

	return a.updateCurrentView(msg)
}

func (a *App) handleGlobalKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return a, tea.Quit

	case "p":
		if a.state.currentView != ViewPlan {
			a.switchView(ViewPlan)
		}
		return a, nil

	case "o":
		if a.state.currentView != ViewObserve {
			a.switchView(ViewObserve)
		}
		return a, nil
	}

	return a.updateCurrentView(msg)
}

func (a *App) switchView(view ViewType) {
	start := time.Now()
	a.state.currentView = view

	if a.state.footer != nil {
		footerModel := a.state.footer.(footer.Model)
		footerModel.UpdateForView(string(view))
		a.state.footer = footerModel
	}

	if a.state.header != nil {
		headerModel := a.state.header.(header.Model)
		headerModel.SetView(string(view))
		a.state.header = headerModel
	}

	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		log.Printf("WARNING: View switch took %v (exceeded 100ms requirement)", elapsed)
	}
	a.state.lastSwitchTime = start
}

func (a *App) propagateSizeUpdate(msg tea.WindowSizeMsg) tea.Cmd {
	var cmds []tea.Cmd

	if a.state.header != nil {
		_, cmd := a.state.header.Update(msg)
		cmds = append(cmds, cmd)
	}

	if a.state.footer != nil {
		_, cmd := a.state.footer.Update(msg)
		cmds = append(cmds, cmd)
	}

	if a.state.planView != nil && a.state.currentView == ViewPlan {
		_, cmd := a.state.planView.Update(msg)
		cmds = append(cmds, cmd)
	}

	if a.state.observeView != nil && a.state.currentView == ViewObserve {
		_, cmd := a.state.observeView.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (a *App) updateCurrentView(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch a.state.currentView {
	case ViewPlan:
		if a.state.planView != nil {
			a.state.planView, cmd = a.state.planView.Update(msg)
		}
	case ViewObserve:
		if a.state.observeView != nil {
			a.state.observeView, cmd = a.state.observeView.Update(msg)
		}
	}

	return a, cmd
}

func (a *App) View() string {
	// Wait for initial window size before rendering
	if a.state.windowWidth == 0 || a.state.windowHeight == 0 {
		return "\n  Initializing TUI... (Press 'q' to quit)\n"
	}

	if !a.state.initialized {
		return a.renderInitPrompt()
	}

	// Always use default renderers if components aren't ready
	header := a.renderDefaultHeader()
	footer := a.renderDefaultFooter()

	// Try to use component views if available
	if a.state.header != nil {
		if v := a.state.header.View(); v != "" {
			header = v
		}
	}

	if a.state.footer != nil {
		if v := a.state.footer.View(); v != "" {
			footer = v
		}
	}

	// Main content
	mainHeight := a.state.windowHeight - 3
	if mainHeight < 1 {
		mainHeight = 1
	}

	var mainContent string
	switch a.state.currentView {
	case ViewPlan:
		if a.state.planView != nil {
			if v := a.state.planView.View(); v != "" {
				mainContent = v
			} else {
				mainContent = a.renderPlaceholderView("Plan View", mainHeight)
			}
		} else {
			mainContent = a.renderPlaceholderView("Plan View", mainHeight)
		}
	case ViewObserve:
		if a.state.observeView != nil {
			if v := a.state.observeView.View(); v != "" {
				mainContent = v
			} else {
				mainContent = a.renderPlaceholderView("Observe View", mainHeight)
			}
		} else {
			mainContent = a.renderPlaceholderView("Observe View", mainHeight)
		}
	default:
		mainContent = a.renderPlaceholderView("Unknown View", mainHeight)
	}

	return header + "\n" + mainContent + "\n" + footer
}

func (a *App) renderInitPrompt() string {
	baseStyles := styles.GetDefaultStyles()

	style := lipgloss.NewStyle().
		Width(a.state.windowWidth).
		Height(a.state.windowHeight).
		Align(lipgloss.Center, lipgloss.Center)

	message := baseStyles.Error.Render("Project not initialized") + "\n\n" +
		baseStyles.Text.Render("Run 'spcstr init' to initialize the project") + "\n\n" +
		baseStyles.TextMuted.Render("Press 'q' to quit")

	return style.Render(message)
}

func (a *App) renderDefaultHeader() string {
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		MaxWidth(a.state.windowWidth).
		Padding(0, 1)

	viewName := string(a.state.currentView)
	if viewName == "plan" {
		viewName = "Plan"
	} else if viewName == "observe" {
		viewName = "Observe"
	}

	status := "inactive"
	if a.state.initialized {
		status = "active"
	}

	left := fmt.Sprintf("spcstr | %s View", viewName)
	right := fmt.Sprintf("Session: %s", status)

	// Calculate available space for padding
	availableWidth := a.state.windowWidth - 2 // Account for padding
	contentWidth := len(left) + len(right)
	padding := availableWidth - contentWidth
	if padding < 0 {
		padding = 0
	}

	spacer := strings.Repeat(" ", padding)
	return headerStyle.Render(left + spacer + right)
}

func (a *App) renderDefaultFooter() string {
	footerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("248")).
		MaxWidth(a.state.windowWidth).
		Padding(0, 1)

	keybinds := "[p] Plan  [o] Observe  [q] Quit"

	return footerStyle.Render(keybinds)
}

func (a *App) renderPlaceholderView(name string, height int) string {
	// Use MaxWidth and account for border
	maxWidth := a.state.windowWidth - 2
	if maxWidth < 20 {
		maxWidth = 20
	}

	style := lipgloss.NewStyle().
		MaxWidth(maxWidth).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("241"))

	content := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf("%s\n(placeholder)", name))

	return style.Render(content)
}

func (a *App) Run(ctx context.Context) error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}
