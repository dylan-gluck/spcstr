package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.initialized {
		return "Loading..."
	}

	if m.width == 0 || m.height == 0 {
		return "Terminal size detection in progress..."
	}

	header := renderHeader(m)
	content := renderContent(m)
	footer := renderFooter(m)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

func renderContent(m Model) string {
	var content string

	switch m.currentView {
	case ViewPlan:
		content = m.planContent
	case ViewObserve:
		content = m.observeContent
	default:
		content = "Unknown view"
	}

	contentHeight := m.height - 4
	if contentHeight < 1 {
		contentHeight = 1
	}

	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(contentHeight).
		Padding(1, 2)

	return contentStyle.Render(content)
}
