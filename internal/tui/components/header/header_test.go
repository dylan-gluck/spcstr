package header

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewHeader(t *testing.T) {
	h := New()

	if h.height != 1 {
		t.Errorf("Header height should be 1, got %d", h.height)
	}

	if h.currentView != "Plan" {
		t.Errorf("Default view should be Plan, got %s", h.currentView)
	}

	if h.sessionStatus != "inactive" {
		t.Errorf("Default session status should be inactive, got %s", h.sessionStatus)
	}
}

func TestHeaderUpdate(t *testing.T) {
	h := New()

	// Test window size update
	sizeMsg := tea.WindowSizeMsg{
		Width:  100,
		Height: 30,
	}

	model, _ := h.Update(sizeMsg)
	updated := model.(Model)

	if updated.width != 100 {
		t.Errorf("Header width should be 100, got %d", updated.width)
	}

	if updated.height != 1 {
		t.Errorf("Header height should remain 1, got %d", updated.height)
	}
}

func TestHeaderView(t *testing.T) {
	h := New()
	h.width = 80

	// Test with default values
	view := h.View()

	if view == "" {
		t.Error("Header view should not be empty")
	}

	if !strings.Contains(view, "spcstr") {
		t.Error("Header should contain 'spcstr'")
	}

	if !strings.Contains(view, "Plan View") {
		t.Error("Header should contain 'Plan View'")
	}

	if !strings.Contains(view, "Session: inactive") {
		t.Error("Header should contain 'Session: inactive'")
	}
}

func TestHeaderSetters(t *testing.T) {
	h := New()

	// Test SetView
	h.SetView("Observe")
	if h.currentView != "Observe" {
		t.Errorf("View should be Observe, got %s", h.currentView)
	}

	// Test SetSessionStatus
	h.SetSessionStatus("active")
	if h.sessionStatus != "active" {
		t.Errorf("Session status should be active, got %s", h.sessionStatus)
	}

	// Verify changes appear in view
	h.width = 80
	view := h.View()

	if !strings.Contains(view, "Observe View") {
		t.Error("Header should contain 'Observe View' after SetView")
	}

	if !strings.Contains(view, "Session: active") {
		t.Error("Header should contain 'Session: active' after SetSessionStatus")
	}
}
