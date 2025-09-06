package footer

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewFooter(t *testing.T) {
	f := New()

	if f.height != 1 {
		t.Errorf("Footer height should be 1, got %d", f.height)
	}

	if len(f.keybinds) != 3 {
		t.Errorf("Default keybinds should be 3, got %d", len(f.keybinds))
	}

	// Check default keybinds
	expectedKeys := []string{"p", "o", "q"}
	for i, kb := range f.keybinds {
		if kb.Key != expectedKeys[i] {
			t.Errorf("Keybind %d should be %s, got %s", i, expectedKeys[i], kb.Key)
		}
	}
}

func TestFooterUpdate(t *testing.T) {
	f := New()

	// Test window size update
	sizeMsg := tea.WindowSizeMsg{
		Width:  120,
		Height: 40,
	}

	model, _ := f.Update(sizeMsg)
	updated := model.(Model)

	if updated.width != 120 {
		t.Errorf("Footer width should be 120, got %d", updated.width)
	}

	if updated.height != 1 {
		t.Errorf("Footer height should remain 1, got %d", updated.height)
	}
}

func TestFooterView(t *testing.T) {
	f := New()
	f.width = 80

	view := f.View()

	if view == "" {
		t.Error("Footer view should not be empty")
	}

	// Check for default keybinds in view
	expectedStrings := []string{"[p]", "Plan", "[o]", "Observe", "[q]", "Quit"}
	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("Footer should contain '%s'", expected)
		}
	}
}

func TestUpdateForView(t *testing.T) {
	tests := []struct {
		name             string
		viewName         string
		expectedKeybinds int
		shouldContain    []string
	}{
		{
			name:             "Plan view",
			viewName:         "Plan",
			expectedKeybinds: 5, // 3 global + 2 view-specific
			shouldContain:    []string{"Navigate", "Open"},
		},
		{
			name:             "Observe view",
			viewName:         "observe",
			expectedKeybinds: 5, // 3 global + 2 view-specific
			shouldContain:    []string{"Navigate", "Refresh"},
		},
		{
			name:             "Unknown view",
			viewName:         "Unknown",
			expectedKeybinds: 3, // Only global keybinds
			shouldContain:    []string{"Plan", "Observe", "Quit"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New()
			f.UpdateForView(tt.viewName)

			if len(f.keybinds) != tt.expectedKeybinds {
				t.Errorf("Expected %d keybinds, got %d", tt.expectedKeybinds, len(f.keybinds))
			}

			// Check keybind descriptions
			var descriptions []string
			for _, kb := range f.keybinds {
				descriptions = append(descriptions, kb.Description)
			}

			for _, expected := range tt.shouldContain {
				found := false
				for _, desc := range descriptions {
					if desc == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Keybinds should contain '%s'", expected)
				}
			}
		})
	}
}

func TestSetKeybinds(t *testing.T) {
	f := New()

	customKeybinds := []Keybind{
		{Key: "h", Description: "Help", Global: true},
		{Key: "s", Description: "Settings", Global: true},
	}

	f.SetKeybinds(customKeybinds)

	if len(f.keybinds) != 2 {
		t.Errorf("Keybinds should be 2, got %d", len(f.keybinds))
	}

	if f.keybinds[0].Key != "h" {
		t.Errorf("First keybind should be 'h', got %s", f.keybinds[0].Key)
	}

	if f.keybinds[1].Description != "Settings" {
		t.Errorf("Second keybind description should be 'Settings', got %s", f.keybinds[1].Description)
	}
}
