package plan

import (
	"fmt"
	"os"
	"strings"

	"github.com/dylan/spcstr/internal/docs"
	"github.com/dylan/spcstr/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PaneType string

const (
	PaneList    PaneType = "list"
	PaneContent PaneType = "content"
)

type ViewMode string

const (
	ViewModeNormal   ViewMode = "normal"
	ViewModeSpec     ViewMode = "spec"
	ViewModeWorkflow ViewMode = "workflow"
	ViewModeConfig   ViewMode = "config"
)

type Model struct {
	width       int
	height      int
	baseStyles  styles.BaseStyles
	paneStyles  PaneStyles
	state       *PlanState
	docEngine   *docs.Engine
	fileWatcher *docs.FileWatcher
	initialized bool
}

type PlanState struct {
	documents     []docs.DocumentIndex
	selected      int
	content       string
	focusedPane   PaneType
	viewMode      ViewMode
	scrollOffset  int
	contentScroll int
	loading       bool
	error         string
}

type PaneStyles struct {
	ListPane       lipgloss.Style
	ContentPane    lipgloss.Style
	ListItem       lipgloss.Style
	SelectedItem   lipgloss.Style
	FocusedBorder  lipgloss.Style
	UnfocusedBorder lipgloss.Style
	CategoryHeader lipgloss.Style
	Loading        lipgloss.Style
}

type documentsLoadedMsg struct {
	documents []docs.DocumentIndex
	err       error
}

type documentContentMsg struct {
	content string
	err     error
}

func New() Model {
	cwd, _ := os.Getwd()
	baseStyles := styles.GetDefaultStyles()
	theme := styles.DefaultTheme
	
	fileWatcher, _ := docs.NewFileWatcher(cwd)
	
	return Model{
		baseStyles:  baseStyles,
		paneStyles:  createPaneStyles(theme),
		docEngine:   docs.NewEngine(cwd),
		fileWatcher: fileWatcher,
		state: &PlanState{
			documents:   []docs.DocumentIndex{},
			selected:    0,
			focusedPane: PaneList,
			viewMode:    ViewModeNormal,
		},
	}
}

func createPaneStyles(theme styles.Theme) PaneStyles {
	focusedBorder := lipgloss.RoundedBorder()
	unfocusedBorder := lipgloss.NormalBorder()
	
	return PaneStyles{
		ListPane: lipgloss.NewStyle().
			BorderStyle(unfocusedBorder).
			BorderForeground(theme.BorderMuted).
			Padding(1, 1),
		ContentPane: lipgloss.NewStyle().
			BorderStyle(unfocusedBorder).
			BorderForeground(theme.BorderMuted).
			Padding(1, 1),
		ListItem: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(theme.Text),
		SelectedItem: lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(theme.Secondary).
			Bold(true),
		FocusedBorder: lipgloss.NewStyle().
			BorderStyle(focusedBorder).
			BorderForeground(theme.Primary),
		UnfocusedBorder: lipgloss.NewStyle().
			BorderStyle(unfocusedBorder).
			BorderForeground(theme.BorderMuted),
		CategoryHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			MarginTop(1),
		Loading: lipgloss.NewStyle().
			Foreground(theme.TextMuted).
			Italic(true),
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.loadDocuments}
	
	if m.fileWatcher != nil {
		cmds = append(cmds, m.fileWatcher.Watch())
	}
	
	return tea.Batch(cmds...)
}

func (m Model) loadDocuments() tea.Msg {
	documents, err := m.docEngine.ScanAndIndex()
	return documentsLoadedMsg{
		documents: documents,
		err:       err,
	}
}

func (m Model) loadDocumentContent(path string) tea.Cmd {
	return func() tea.Msg {
		content, err := m.docEngine.RenderDocument(path)
		return documentContentMsg{
			content: content,
			err:     err,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 3  // Account for header/footer
		
		// Update renderer width (simpler calculation)
		if m.docEngine != nil && m.width > 40 {
			// Just use 2/3 of width for content, no complex math
			contentWidth := (m.width * 2) / 3
			m.docEngine.Renderer.SetWidth(contentWidth)
		}
		
	case documentsLoadedMsg:
		if msg.err != nil {
			m.state.error = fmt.Sprintf("Failed to load documents: %v", msg.err)
		} else {
			m.state.documents = msg.documents
			if len(msg.documents) > 0 && !m.initialized {
				m.initialized = true
				return m, m.loadDocumentContent(msg.documents[0].Path)
			}
		}
		m.state.loading = false
		
	case documentContentMsg:
		if msg.err != nil {
			m.state.error = fmt.Sprintf("Failed to load document: %v", msg.err)
		} else {
			m.state.content = msg.content
			m.state.error = ""
		}
		m.state.loading = false
		
	case docs.FileChangeMsg:
		cmds := []tea.Cmd{}
		
		if msg.Operation == "created" || msg.Operation == "removed" {
			cmds = append(cmds, m.loadDocuments)
		} else if msg.Operation == "modified" {
			if len(m.state.documents) > m.state.selected {
				currentDoc := m.state.documents[m.state.selected]
				if currentDoc.Path == msg.Path {
					cmds = append(cmds, m.loadDocumentContent(msg.Path))
				}
			}
		}
		
		if m.fileWatcher != nil {
			cmds = append(cmds, m.fileWatcher.Watch())
		}
		
		return m, tea.Batch(cmds...)
		
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		if m.state.focusedPane == PaneList {
			m.state.focusedPane = PaneContent
		} else {
			m.state.focusedPane = PaneList
		}
		
	case "up", "k":
		if m.state.focusedPane == PaneList && m.state.selected > 0 {
			m.state.selected--
			if len(m.state.documents) > m.state.selected {
				m.state.loading = true
				return m, m.loadDocumentContent(m.state.documents[m.state.selected].Path)
			}
		} else if m.state.focusedPane == PaneContent && m.state.contentScroll > 0 {
			m.state.contentScroll--
		}
		
	case "down", "j":
		if m.state.focusedPane == PaneList && m.state.selected < len(m.state.documents)-1 {
			m.state.selected++
			if len(m.state.documents) > m.state.selected {
				m.state.loading = true
				return m, m.loadDocumentContent(m.state.documents[m.state.selected].Path)
			}
		} else if m.state.focusedPane == PaneContent {
			m.state.contentScroll++
		}
		
	case "enter":
		if m.state.focusedPane == PaneList && len(m.state.documents) > m.state.selected {
			m.state.loading = true
			return m, m.loadDocumentContent(m.state.documents[m.state.selected].Path)
		}
		
	case "home", "g":
		if m.state.focusedPane == PaneList && len(m.state.documents) > 0 {
			m.state.selected = 0
			m.state.loading = true
			return m, m.loadDocumentContent(m.state.documents[0].Path)
		} else if m.state.focusedPane == PaneContent {
			m.state.contentScroll = 0
		}
		
	case "end", "G":
		if m.state.focusedPane == PaneList && len(m.state.documents) > 0 {
			m.state.selected = len(m.state.documents) - 1
			m.state.loading = true
			return m, m.loadDocumentContent(m.state.documents[m.state.selected].Path)
		}
		
	case "s":
		m.state.viewMode = ViewModeSpec
		
	case "w":
		m.state.viewMode = ViewModeWorkflow
		
	case "c":
		m.state.viewMode = ViewModeConfig
		
	case "n":
		m.state.viewMode = ViewModeNormal
	}
	
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	
	// Simple division: left pane gets 1/3 of width
	// Account for borders (2) and padding (2) per pane = 4 total per pane
	leftPaneWidth := (m.width / 3) - 4
	rightPaneWidth := m.width - (m.width / 3) - 4
	
	leftPane := m.renderListPane(leftPaneWidth, m.height-2)
	rightPane := m.renderContentPane(rightPaneWidth, m.height-2)
	
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

func (m Model) renderListPane(width, height int) string {
	var listItems []string
	
	if m.state.loading && len(m.state.documents) == 0 {
		listItems = append(listItems, m.paneStyles.Loading.Render("Loading documents..."))
	} else if len(m.state.documents) == 0 {
		listItems = append(listItems, m.baseStyles.TextMuted.Render("No documents found"))
	} else {
		currentType := docs.DocumentType("")
		for i, doc := range m.state.documents {
			if doc.Type != currentType {
				currentType = doc.Type
				header := fmt.Sprintf("── %s ──", strings.ToUpper(string(doc.Type)))
				listItems = append(listItems, m.paneStyles.CategoryHeader.Render(header))
			}
			
			prefix := "  "
			item := doc.Title
			if i == m.state.selected {
				prefix = "▸ "
				item = m.paneStyles.SelectedItem.Render(item)
			} else {
				item = m.paneStyles.ListItem.Render(item)
			}
			listItems = append(listItems, prefix+item)
		}
	}
	
	content := strings.Join(listItems, "\n")
	
	paneStyle := m.paneStyles.ListPane
	if m.state.focusedPane == PaneList {
		paneStyle = paneStyle.BorderStyle(m.paneStyles.FocusedBorder.GetBorderStyle()).
			BorderForeground(m.paneStyles.FocusedBorder.GetBorderTopForeground())
	}
	
	// KEEP MaxWidth/MaxHeight - they prevent overflow!
	return paneStyle.
		Width(width).
		Height(height).
		MaxWidth(width).    // Keep this!
		MaxHeight(height).  // Keep this!
		Render(content)
}

func (m Model) renderContentPane(width, height int) string {
	var content string
	
	if m.state.loading && m.state.content != "" {
		content = m.state.content
	} else if m.state.loading {
		content = m.paneStyles.Loading.Render("Loading document...")
	} else if m.state.error != "" {
		content = m.baseStyles.Error.Render(m.state.error)
	} else if m.state.content == "" {
		content = m.baseStyles.TextMuted.Render("Select a document to view")
	} else {
		content = m.state.content
	}
	
	if m.state.viewMode != ViewModeNormal {
		modeHeader := fmt.Sprintf("[ %s MODE ]", strings.ToUpper(string(m.state.viewMode)))
		content = m.baseStyles.Title.Render(modeHeader) + "\n\n" + content
	}
	
	lines := strings.Split(content, "\n")
	visibleHeight := height - 2
	if m.state.contentScroll > len(lines)-visibleHeight {
		m.state.contentScroll = len(lines) - visibleHeight
	}
	if m.state.contentScroll < 0 {
		m.state.contentScroll = 0
	}
	
	if m.state.contentScroll < len(lines) {
		endLine := m.state.contentScroll + visibleHeight
		if endLine > len(lines) {
			endLine = len(lines)
		}
		visibleLines := lines[m.state.contentScroll:endLine]
		content = strings.Join(visibleLines, "\n")
	}
	
	paneStyle := m.paneStyles.ContentPane
	if m.state.focusedPane == PaneContent {
		paneStyle = paneStyle.BorderStyle(m.paneStyles.FocusedBorder.GetBorderStyle()).
			BorderForeground(m.paneStyles.FocusedBorder.GetBorderTopForeground())
	}
	
	// KEEP MaxWidth/MaxHeight - they prevent overflow!
	return paneStyle.
		Width(width).
		Height(height).
		MaxWidth(width).    // Keep this!
		MaxHeight(height).  // Keep this!
		Render(content)
}