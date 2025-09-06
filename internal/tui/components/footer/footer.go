package footer

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Keybind struct {
	Key         string
	Description string
	Global      bool
}

type Model struct {
	width    int
	height   int
	keybinds []Keybind
	styles   Styles
}

type Styles struct {
	Footer  lipgloss.Style
	Key     lipgloss.Style
	Desc    lipgloss.Style
	Divider lipgloss.Style
}

func New() Model {
	return Model{
		height:   1,
		keybinds: defaultKeybinds(),
		styles:   defaultStyles(),
	}
}

func defaultStyles() Styles {
	return Styles{
		Footer: lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("248")).
			Padding(0, 1),
		Key: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true),
		Desc: lipgloss.NewStyle().
			Foreground(lipgloss.Color("248")),
		Divider: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),
	}
}

func defaultKeybinds() []Keybind {
	return []Keybind{
		{Key: "p", Description: "Plan", Global: true},
		{Key: "o", Description: "Observe", Global: true},
		{Key: "q", Description: "Quit", Global: true},
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

	var keybindStrs []string
	for _, kb := range m.keybinds {
		key := m.styles.Key.Render("[" + kb.Key + "]")
		desc := m.styles.Desc.Render(kb.Description)
		keybindStrs = append(keybindStrs, key+" "+desc)
	}

	content := strings.Join(keybindStrs, m.styles.Divider.Render("  "))

	return m.styles.Footer.MaxWidth(m.width).Render(content)
}

func (m *Model) SetKeybinds(keybinds []Keybind) {
	m.keybinds = keybinds
}

func (m *Model) UpdateForView(viewName string) {
	baseKeybinds := []Keybind{
		{Key: "p", Description: "Plan", Global: true},
		{Key: "o", Description: "Observe", Global: true},
		{Key: "q", Description: "Quit", Global: true},
	}

	switch viewName {
	case "Plan", "plan":
		m.keybinds = append(baseKeybinds,
			Keybind{Key: "↑/↓", Description: "Navigate", Global: false},
			Keybind{Key: "Enter", Description: "Open", Global: false},
		)
	case "Observe", "observe":
		m.keybinds = append(baseKeybinds,
			Keybind{Key: "↑/↓", Description: "Navigate", Global: false},
			Keybind{Key: "r", Description: "Refresh", Global: false},
		)
	default:
		m.keybinds = baseKeybinds
	}
}

func (m Model) Width() int {
	return m.width
}

func (m Model) Height() int {
	return m.height
}
