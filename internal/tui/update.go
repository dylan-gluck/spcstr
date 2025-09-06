package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.initialized = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			if m.currentView != ViewPlan {
				m.currentView = ViewPlan
				m.lastSwitch = time.Now()
			}
			return m, nil

		case "o":
			if m.currentView != ViewObserve {
				m.currentView = ViewObserve
				m.lastSwitch = time.Now()
			}
			return m, nil

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}
