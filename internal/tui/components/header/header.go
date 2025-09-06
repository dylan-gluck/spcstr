package header

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width         int
	height        int
	currentView   string
	sessionStatus string
	styles        Styles
}

type Styles struct {
	Header lipgloss.Style
	Title  lipgloss.Style
	Status lipgloss.Style
}

func New() Model {
	return Model{
		height:        1,
		currentView:   "Plan",
		sessionStatus: "inactive",
		styles:        defaultStyles(),
	}
}

func defaultStyles() Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("248")),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = 1
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	left := m.styles.Title.Render(fmt.Sprintf("spcstr | %s View", m.currentView))

	statusIcon := "○"
	if m.sessionStatus == "active" {
		statusIcon = "●"
	}
	right := m.styles.Status.Render(fmt.Sprintf("%s Session: %s", statusIcon, m.sessionStatus))

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)

	// Calculate available width accounting for padding
	availableWidth := m.width - 2 // Header has padding of 1 on each side
	contentWidth := leftWidth + rightWidth
	paddingWidth := availableWidth - contentWidth
	if paddingWidth < 0 {
		paddingWidth = 0
	}

	content := left + strings.Repeat(" ", paddingWidth) + right

	return m.styles.Header.MaxWidth(m.width).Render(content)
}

func (m *Model) SetView(view string) {
	m.currentView = view
}

func (m *Model) SetSessionStatus(status string) {
	m.sessionStatus = status
}

func (m Model) Width() int {
	return m.width
}

func (m Model) Height() int {
	return m.height
}
