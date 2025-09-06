package plan

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width  int
	height int
	styles Styles
}

type Styles struct {
	Container lipgloss.Style
	Title     lipgloss.Style
	Content   lipgloss.Style
}

func New() Model {
	return Model{
		styles: defaultStyles(),
	}
}

func defaultStyles() Styles {
	return Styles{
		Container: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1),
		Content: lipgloss.NewStyle().
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
		m.height = msg.Height - 3 // Account for header and footer
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	title := m.styles.Title.Render("Plan View")
	content := m.styles.Content.Render("Document browser will be displayed here\n\nAvailable documents:\n• PRD\n• Architecture\n• Epics\n• Stories")

	innerContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
	)

	// Use MaxWidth for responsive layout
	// Account for border (2) and padding (4)
	maxWidth := m.width - 2
	if maxWidth < 20 {
		maxWidth = 20
	}

	maxHeight := m.height - 2
	if maxHeight < 5 {
		maxHeight = 5
	}

	return m.styles.Container.
		MaxWidth(maxWidth).
		MaxHeight(maxHeight).
		Render(innerContent)
}
