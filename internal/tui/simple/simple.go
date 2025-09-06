package simple

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int
	view   string
}

func New() Model {
	return Model{
		view: "plan",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "p":
			m.view = "plan"
			return m, nil
		case "o":
			m.view = "observe"
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "\n  Initializing...\n\n"
	}

	header := fmt.Sprintf("═══ spcstr | %s view ═══", m.view)
	footer := "  [p] Plan  [o] Observe  [q] Quit"

	var content string
	switch m.view {
	case "plan":
		content = "\n  📄 Plan View\n\n  Document browser will be here\n"
	case "observe":
		content = "\n  👁  Observe View\n\n  Session monitoring will be here\n"
	}

	return fmt.Sprintf("\n%s\n%s\n\n%s\n", header, content, footer)
}
