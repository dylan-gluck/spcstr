package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ViewType int

const (
	ViewPlan ViewType = iota
	ViewObserve
)

type Model struct {
	currentView    ViewType
	width          int
	height         int
	initialized    bool
	projectPath    string
	planContent    string
	observeContent string
	sessionActive  bool
	lastSwitch     time.Time
}

func NewModel() Model {
	return Model{
		currentView:    ViewPlan,
		initialized:    false,
		sessionActive:  false,
		planContent:    "Plan View\n\nDocuments:\n• PRD.md\n• Architecture.md\n• Epic-1.md\n\nPress 'o' for Observe view",
		observeContent: "Observe View\n\nSession Monitor\n• No active session\n\nPress 'p' for Plan view",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
