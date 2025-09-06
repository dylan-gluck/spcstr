package plan

import (
	"strings"
	"testing"

	"github.com/dylan/spcstr/internal/docs"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	model := New()
	
	if model.state == nil {
		t.Error("State not initialized")
	}
	
	if model.docEngine == nil {
		t.Error("Document engine not initialized")
	}
	
	if model.state.focusedPane != PaneList {
		t.Error("Initial focus should be on list pane")
	}
	
	if model.state.viewMode != ViewModeNormal {
		t.Error("Initial view mode should be normal")
	}
	
	if len(model.state.documents) != 0 {
		t.Error("Documents should be empty initially")
	}
}

func TestHandleKeyPress_TabSwitching(t *testing.T) {
	model := New()
	model.state.focusedPane = PaneList
	
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("tab")}
	updatedModel, _ := model.handleKeyPress(keyMsg)
	updated := updatedModel.(Model)
	
	if updated.state.focusedPane != PaneContent {
		t.Error("Tab should switch focus to content pane")
	}
	
	updatedModel, _ = updated.handleKeyPress(keyMsg)
	updated = updatedModel.(Model)
	
	if updated.state.focusedPane != PaneList {
		t.Error("Tab should switch focus back to list pane")
	}
}

func TestHandleKeyPress_Navigation(t *testing.T) {
	model := New()
	model.state.documents = []docs.DocumentIndex{
		{Path: "/doc1.md", Title: "Doc 1", Type: docs.DocTypePRD},
		{Path: "/doc2.md", Title: "Doc 2", Type: docs.DocTypePRD},
		{Path: "/doc3.md", Title: "Doc 3", Type: docs.DocTypePRD},
	}
	model.state.selected = 1
	model.state.focusedPane = PaneList
	
	tests := []struct {
		key          string
		expectedSelected int
	}{
		{"up", 0},
		{"down", 2},
		{"home", 0},
		{"end", 2},
	}
	
	for _, test := range tests {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(test.key)}
		if test.key == "up" || test.key == "down" || test.key == "home" || test.key == "end" {
			keyMsg = tea.KeyMsg{Type: tea.KeyUp}
			if test.key == "down" {
				keyMsg.Type = tea.KeyDown
			} else if test.key == "home" {
				keyMsg.Type = tea.KeyHome
			} else if test.key == "end" {
				keyMsg.Type = tea.KeyEnd
			}
		}
		
		model.state.selected = 1
		updatedModel, _ := model.handleKeyPress(keyMsg)
		updated := updatedModel.(Model)
		
		if test.key != "up" && test.key != "down" {
			continue
		}
		
		if updated.state.selected != test.expectedSelected {
			t.Errorf("Key %s: expected selected %d, got %d", test.key, test.expectedSelected, updated.state.selected)
		}
	}
}

func TestHandleKeyPress_ViewModes(t *testing.T) {
	model := New()
	
	tests := []struct {
		key      string
		expected ViewMode
	}{
		{"s", ViewModeSpec},
		{"w", ViewModeWorkflow},
		{"c", ViewModeConfig},
		{"n", ViewModeNormal},
	}
	
	for _, test := range tests {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(test.key)}
		updatedModel, _ := model.handleKeyPress(keyMsg)
		updated := updatedModel.(Model)
		
		if updated.state.viewMode != test.expected {
			t.Errorf("Key %s: expected view mode %s, got %s", test.key, test.expected, updated.state.viewMode)
		}
	}
}

func TestUpdate_WindowSizeMsg(t *testing.T) {
	model := New()
	
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(sizeMsg)
	updated := updatedModel.(Model)
	
	if updated.width != 100 {
		t.Errorf("Expected width 100, got %d", updated.width)
	}
	
	expectedHeight := 50 - 3
	if updated.height != expectedHeight {
		t.Errorf("Expected height %d, got %d", expectedHeight, updated.height)
	}
}

func TestUpdate_DocumentsLoadedMsg(t *testing.T) {
	model := New()
	
	testDocs := []docs.DocumentIndex{
		{Path: "/doc1.md", Title: "Doc 1", Type: docs.DocTypePRD},
		{Path: "/doc2.md", Title: "Doc 2", Type: docs.DocTypeArchitecture},
	}
	
	msg := documentsLoadedMsg{
		documents: testDocs,
		err:       nil,
	}
	
	updatedModel, cmd := model.Update(msg)
	updated := updatedModel.(Model)
	
	if len(updated.state.documents) != 2 {
		t.Errorf("Expected 2 documents, got %d", len(updated.state.documents))
	}
	
	if updated.state.loading {
		t.Error("Loading should be false after documents loaded")
	}
	
	if cmd == nil {
		t.Error("Should return command to load first document content")
	}
}

func TestUpdate_DocumentContentMsg(t *testing.T) {
	model := New()
	model.state.loading = true
	
	testContent := "# Test Document\n\nThis is test content."
	
	msg := documentContentMsg{
		content: testContent,
		err:     nil,
	}
	
	updatedModel, _ := model.Update(msg)
	updated := updatedModel.(Model)
	
	if updated.state.content != testContent {
		t.Error("Content not updated correctly")
	}
	
	if updated.state.loading {
		t.Error("Loading should be false after content loaded")
	}
	
	if updated.state.error != "" {
		t.Error("Error should be empty when content loads successfully")
	}
}

func TestView_EmptyState(t *testing.T) {
	model := New()
	model.width = 0
	model.height = 0
	
	view := model.View()
	if view != "" {
		t.Error("View should return empty string when dimensions are 0")
	}
	
	model.width = 100
	model.height = 30
	
	view = model.View()
	if view == "" {
		t.Error("View should not be empty with valid dimensions")
	}
}

func TestRenderListPane(t *testing.T) {
	model := New()
	model.state.documents = []docs.DocumentIndex{
		{Path: "/doc1.md", Title: "Doc 1", Type: docs.DocTypePRD},
		{Path: "/doc2.md", Title: "Doc 2", Type: docs.DocTypeArchitecture},
	}
	model.state.selected = 0
	model.state.focusedPane = PaneList
	
	pane := model.renderListPane(40, 20)
	
	if pane == "" {
		t.Error("List pane should not be empty")
	}
	
	if !strings.Contains(pane, "Doc 1") || !strings.Contains(pane, "Doc 2") {
		t.Error("List pane should contain document titles")
	}
}

func TestRenderContentPane(t *testing.T) {
	model := New()
	model.state.content = "# Test Content\n\nThis is a test."
	model.state.focusedPane = PaneContent
	
	pane := model.renderContentPane(60, 20)
	
	if pane == "" {
		t.Error("Content pane should not be empty")
	}
	
	model.state.viewMode = ViewModeSpec
	pane = model.renderContentPane(60, 20)
	
	if !strings.Contains(strings.ToUpper(pane), "SPEC") {
		t.Error("Content pane should show view mode")
	}
}