package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Bold(true)

	statusActiveStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("86")).
				Bold(true)

	statusInactiveStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("241"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)
)

func renderHeader(m Model) string {
	var viewName string
	switch m.currentView {
	case ViewPlan:
		viewName = "PLAN"
	case ViewObserve:
		viewName = "OBSERVE"
	default:
		viewName = "UNKNOWN"
	}

	left := fmt.Sprintf(" %s ", viewName)

	var right string
	if m.sessionActive {
		right = " ● Active "
		right = statusActiveStyle.Render(right)
	} else {
		right = " ○ Inactive "
		right = statusInactiveStyle.Render(right)
	}

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(stripAnsi(right))
	gap := m.width - leftWidth - rightWidth

	if gap < 0 {
		gap = 0
	}

	middle := headerStyle.Render(strings.Repeat(" ", gap))
	left = headerStyle.Render(left)

	return left + middle + right
}

func renderFooter(m Model) string {
	var keys []string

	switch m.currentView {
	case ViewPlan:
		keys = []string{
			keyStyle.Render("o") + " observe",
			keyStyle.Render("q") + " quit",
		}
	case ViewObserve:
		keys = []string{
			keyStyle.Render("p") + " plan",
			keyStyle.Render("q") + " quit",
		}
	default:
		keys = []string{
			keyStyle.Render("q") + " quit",
		}
	}

	footer := " " + strings.Join(keys, " • ")

	return footerStyle.
		Width(m.width).
		Render(footer)
}

func stripAnsi(str string) string {
	var result strings.Builder
	var inEscape bool

	for _, r := range str {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}

	return result.String()
}
