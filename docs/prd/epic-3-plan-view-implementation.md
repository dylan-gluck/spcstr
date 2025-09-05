# Epic 3: Plan View Implementation

Build the Plan View that provides access to planning documents with multiple modes (Spec, Workflow, Config). Implement document indexing with search capabilities and markdown preview with syntax highlighting.

## Story 3.1: Document Indexing and Discovery

As a user,
I want spcstr to automatically find all planning documents,
so that I can quickly navigate to any document.

### Acceptance Criteria
1: Recursive search finds all markdown files in configured paths
2: Documents are organized by type (PRD, Architecture, Epic, Story)
3: Search supports fuzzy matching on file names and content
4: Index updates automatically when files change
5: Performance remains fast with hundreds of documents

## Story 3.2: Markdown Preview with Syntax Highlighting

As a user,
I want to preview markdown documents with proper formatting,
so that I can read planning documents without leaving the TUI.

### Acceptance Criteria
1: Markdown renders with headers, lists, and emphasis
2: Code blocks have syntax highlighting
3: Links are visually distinguished
4: Preview updates in real-time during edits
5: Large documents scroll smoothly

## Story 3.3: Plan View Mode Navigation

As a user,
I want to switch between Spec, Workflow, and Config modes,
so that I can access different types of planning information.

### Acceptance Criteria
1: Mode switching via keyboard shortcuts (s, w, c)
2: Each mode shows relevant document categories
3: Mode state persists across view switches
4: Visual indicator shows current mode
5: Mode-specific keybindings are available

## Story 3.4: Document Navigation and Selection

As a user,
I want to navigate and select documents efficiently,
so that I can quickly access the information I need.

### Acceptance Criteria
1: Arrow keys navigate document list
2: Tab switches focus between list and preview
3: Enter opens document for detailed view
4: Search narrows list in real-time
5: Recently accessed documents appear first
