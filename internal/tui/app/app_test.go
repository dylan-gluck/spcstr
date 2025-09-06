package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAppInitialization(t *testing.T) {
	app := New()

	if app.state == nil {
		t.Fatal("App state should be initialized")
	}

	if app.state.currentView != ViewPlan {
		t.Errorf("Initial view should be Plan, got %s", app.state.currentView)
	}

	if app.state.initialized {
		t.Error("App should not be initialized before Init()")
	}
}

func TestCheckInitialization(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	app := New()

	// Test without .spcstr directory
	app.checkInitialization()
	if app.state.initialized {
		t.Error("Should not be initialized without .spcstr directory")
	}

	// Create .spcstr directory
	os.Mkdir(filepath.Join(tmpDir, ".spcstr"), 0755)

	// Test with .spcstr directory
	app.checkInitialization()
	if !app.state.initialized {
		t.Error("Should be initialized with .spcstr directory")
	}

	// Resolve symlinks for comparison
	expectedPath, _ := filepath.EvalSymlinks(tmpDir)
	actualPath, _ := filepath.EvalSymlinks(app.state.projectPath)
	if actualPath != expectedPath {
		t.Errorf("Project path should be %s, got %s", expectedPath, actualPath)
	}
}

func TestHandleGlobalKeys(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		initialView ViewType
		wantView    ViewType
		wantQuit    bool
	}{
		{
			name:        "Switch to Plan view",
			key:         "p",
			initialView: ViewObserve,
			wantView:    ViewPlan,
			wantQuit:    false,
		},
		{
			name:        "Switch to Observe view",
			key:         "o",
			initialView: ViewPlan,
			wantView:    ViewObserve,
			wantQuit:    false,
		},
		{
			name:        "Stay in Plan when already in Plan",
			key:         "p",
			initialView: ViewPlan,
			wantView:    ViewPlan,
			wantQuit:    false,
		},
		{
			name:        "Quit with q",
			key:         "q",
			initialView: ViewPlan,
			wantView:    ViewPlan,
			wantQuit:    true,
		},
		{
			name:        "Quit with ctrl+c",
			key:         "ctrl+c",
			initialView: ViewPlan,
			wantView:    ViewPlan,
			wantQuit:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New()
			app.state.currentView = tt.initialView
			app.state.initialized = true
			app.initializeViews()

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "ctrl+c" {
				keyMsg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}

			model, cmd := app.handleGlobalKeys(keyMsg)

			if tt.wantQuit {
				if cmd == nil {
					t.Error("Expected quit command")
				}
			}

			if updatedApp, ok := model.(*App); ok {
				if updatedApp.state.currentView != tt.wantView {
					t.Errorf("View should be %s, got %s", tt.wantView, updatedApp.state.currentView)
				}
			}
		})
	}
}

func TestViewSwitchingPerformance(t *testing.T) {
	app := New()
	app.state.initialized = true
	app.initializeViews()

	// Test view switch performance
	start := time.Now()
	app.switchView(ViewObserve)
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("View switch took %v, exceeded 100ms requirement", elapsed)
	}

	if app.state.currentView != ViewObserve {
		t.Error("View should have switched to Observe")
	}
}

func TestWindowResize(t *testing.T) {
	app := New()
	app.state.initialized = true
	app.initializeViews()

	// Send window size message
	sizeMsg := tea.WindowSizeMsg{
		Width:  80,
		Height: 24,
	}

	model, _ := app.Update(sizeMsg)
	updatedApp := model.(*App)

	if updatedApp.state.windowWidth != 80 {
		t.Errorf("Window width should be 80, got %d", updatedApp.state.windowWidth)
	}

	if updatedApp.state.windowHeight != 24 {
		t.Errorf("Window height should be 24, got %d", updatedApp.state.windowHeight)
	}
}

func TestRenderInitPrompt(t *testing.T) {
	app := New()
	app.state.windowWidth = 80
	app.state.windowHeight = 24

	output := app.renderInitPrompt()

	if output == "" {
		t.Error("Init prompt should not be empty")
	}

	// Check for expected content
	expectedStrings := []string{
		"Project not initialized",
		"spcstr init",
		"quit",
	}

	for _, expected := range expectedStrings {
		if !contains(output, expected) {
			t.Errorf("Init prompt should contain '%s'", expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) &&
		(s == substr || s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
