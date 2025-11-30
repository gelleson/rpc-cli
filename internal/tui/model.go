package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"jsonrpc/internal/executor"
	"jsonrpc/internal/parser"
	"jsonrpc/pkg/types"
)

// View represents different UI states
type View int

const (
	ViewFileSelect View = iota
	ViewList
	ViewDetail
	ViewResults
	ViewHelp
)

// ExecutionHistory stores past execution results
type ExecutionHistory struct {
	Timestamp time.Time
	Results   []*types.ExecutionResult
	Selected  []string
}

// Model holds the state of our TUI application
type Model struct {
	// Data
	hclFile  *types.HCLFile
	requests []*types.Request
	selected map[int]struct{}
	results  []*types.ExecutionResult
	filename string
	history  []ExecutionHistory

	// File selection
	hclFiles      []string
	fileCursor    int
	fileSearching bool

	// UI State
	currentView View
	cursor      int
	width       int
	height      int
	error       error
	loading     bool
	ready       bool

	// Search/Filter
	searchMode   bool
	searchInput  textinput.Model
	filteredReqs []*types.Request
	filterMap    map[int]int

	// Components
	viewport viewport.Model
	spinner  spinner.Model
	help     help.Model
	keys     keyMap

	// Execution
	executor  *executor.Executor
	overrides *types.CLIOverrides

	// Styles
	styles *Styles
}

// keyMap defines keybindings
type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Select      key.Binding
	Enter       key.Binding
	Back        key.Binding
	Quit        key.Binding
	Run         key.Binding
	SelectAll   key.Binding
	DeselectAll key.Binding
	Search      key.Binding
	Help        key.Binding
	ClearSearch key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Enter},
		{k.Run, k.SelectAll, k.DeselectAll},
		{k.Search, k.Back, k.Help, k.Quit},
	}
}

func defaultKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter", "l", "right"),
			key.WithHelp("enter/l", "view details"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "h", "left"),
			key.WithHelp("esc/h", "go back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Run: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "run selected"),
		),
		SelectAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "select all"),
		),
		DeselectAll: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("A", "deselect all"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		ClearSearch: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear search"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

// NewModel creates a new TUI model
func NewModel() *Model {
	ti := textinput.New()
	ti.Placeholder = "Search requests..."
	ti.CharLimit = 50

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = DefaultStyles().AccentStyle

	return &Model{
		selected:    make(map[int]struct{}),
		filterMap:   make(map[int]int),
		currentView: ViewList,
		styles:      DefaultStyles(),
		executor:    executor.New(),
		overrides:   types.NewCLIOverrides(),
		searchInput: ti,
		spinner:     sp,
		help:        help.New(),
		keys:        defaultKeyMap(),
		history:     make([]ExecutionHistory, 0),
	}
}

// NewModelWithFile creates a new TUI model that loads a file on initialization
func NewModelWithFile(filename string) *Model {
	m := NewModel()
	m.filename = filename
	m.currentView = ViewList
	return m
}

// NewModelWithFileSelect creates a new TUI model that shows file selection first
func NewModelWithFileSelect(files []string) *Model {
	m := NewModel()
	m.hclFiles = files
	m.currentView = ViewFileSelect
	m.fileSearching = true
	return m
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick}
	if m.filename != "" && m.currentView != ViewFileSelect {
		cmds = append(cmds, m.LoadFile(m.filename))
	}
	return tea.Batch(cmds...)
}

// LoadFile returns a command to load a file into the model
func (m *Model) LoadFile(filename string) tea.Cmd {
	return func() tea.Msg {
		p := parser.New()
		hclFile, err := p.ParseFile(filename)
		if err != nil {
			return loadFileMsg{err: err}
		}
		return loadFileMsg{hclFile: hclFile}
	}
}

// Update handles incoming messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searchMode {
			return m.handleSearchInput(msg)
		}
		return m.handleKeyMsg(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-8)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 8
		}
		m.updateViewportContent()
		return m, nil

	case loadFileMsg:
		if msg.err != nil {
			m.error = msg.err
		} else {
			m.hclFile = msg.hclFile
			m.requests = msg.hclFile.Requests
			m.filteredReqs = msg.hclFile.Requests
			m.selected = make(map[int]struct{})
			m.cursor = 0
			m.currentView = ViewList
			m.buildFilterMap()
			m.updateViewportContent()
		}
		return m, nil

	case executionCompleteMsg:
		m.loading = false
		m.results = msg.results
		if msg.err != nil {
			m.error = msg.err
		} else {
			selectedNames := make([]string, 0, len(m.selected))
			for idx := range m.selected {
				if idx < len(m.requests) {
					selectedNames = append(selectedNames, m.requests[idx].Name)
				}
			}
			m.history = append(m.history, ExecutionHistory{
				Timestamp: time.Now(),
				Results:   msg.results,
				Selected:  selectedNames,
			})
		}
		m.currentView = ViewResults
		m.updateViewportContent()
		return m, nil

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.currentView != ViewHelp && !m.searchMode {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m *Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	switch m.currentView {
	case ViewFileSelect:
		return m.renderFileSelectView()
	case ViewList:
		return m.renderListView()
	case ViewDetail:
		return m.renderDetailView()
	case ViewResults:
		return m.renderResultsView()
	case ViewHelp:
		return m.renderHelpView()
	default:
		return "Unknown view"
	}
}

func (m *Model) updateViewportContent() {
	if !m.ready {
		return
	}

	var content string
	switch m.currentView {
	case ViewFileSelect:
		content = m.buildFileSelectContent()
	case ViewList:
		content = m.buildListContent()
	case ViewDetail:
		content = m.buildDetailContent()
	case ViewResults:
		content = m.buildResultsContent()
	}

	m.viewport.SetContent(content)
}

// handleKeyMsg processes keyboard input
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Quit) {
		return m, tea.Quit
	}

	if m.currentView != ViewFileSelect && key.Matches(msg, m.keys.Help) {
		if m.currentView == ViewHelp {
			m.currentView = ViewList
		} else {
			m.currentView = ViewHelp
		}
		return m, nil
	}

	switch m.currentView {
	case ViewFileSelect:
		return m.handleFileSelectKeys(msg)
	case ViewList:
		return m.handleListKeys(msg)
	case ViewDetail:
		return m.handleDetailKeys(msg)
	case ViewResults:
		return m.handleResultsKeys(msg)
	case ViewHelp:
		return m.handleHelpKeys(msg)
	}
	return m, nil
}

// handleFileSelectKeys handles keys in file selection view
func (m *Model) handleFileSelectKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Up) {
		if m.fileCursor > 0 {
			m.fileCursor--
		}
		m.updateViewportContent()
	}

	if key.Matches(msg, m.keys.Down) {
		if m.fileCursor < len(m.hclFiles)-1 {
			m.fileCursor++
		}
		m.updateViewportContent()
	}

	if key.Matches(msg, m.keys.Enter) {
		if len(m.hclFiles) > 0 && m.fileCursor < len(m.hclFiles) {
			m.filename = m.hclFiles[m.fileCursor]
			m.fileSearching = false
			return m, m.LoadFile(m.filename)
		}
	}

	return m, nil
}

// handleListKeys handles keys in list view
func (m *Model) handleListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Search):
		m.searchMode = true
		m.searchInput.Focus()
		return m, textinput.Blink
	case key.Matches(msg, m.keys.Up):
		m.moveCursorUp()
	case key.Matches(msg, m.keys.Down):
		m.moveCursorDown()
	case key.Matches(msg, m.keys.Select):
		m.toggleCurrentSelection()
	case key.Matches(msg, m.keys.Enter):
		m.viewCurrentDetails()
	case key.Matches(msg, m.keys.Run):
		if len(m.selected) > 0 {
			return m, m.executeSelected()
		}
	case key.Matches(msg, m.keys.SelectAll):
		m.selectAll()
	case key.Matches(msg, m.keys.DeselectAll):
		m.deselectAll()
	}
	return m, nil
}

func (m *Model) moveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
	m.updateViewportContent()
}

func (m *Model) moveCursorDown() {
	if m.cursor < len(m.filteredReqs)-1 {
		m.cursor++
	}
	m.updateViewportContent()
}

func (m *Model) toggleCurrentSelection() {
	actualIdx := m.getActualIndex(m.cursor)
	if actualIdx >= 0 {
		m.toggleSelection(actualIdx)
		m.updateViewportContent()
	}
}

func (m *Model) viewCurrentDetails() {
	if len(m.filteredReqs) > 0 && m.cursor < len(m.filteredReqs) {
		m.currentView = ViewDetail
		m.updateViewportContent()
	}
}

func (m *Model) selectAll() {
	for i := range m.requests {
		m.selected[i] = struct{}{}
	}
	m.updateViewportContent()
}

func (m *Model) deselectAll() {
	m.selected = make(map[int]struct{})
	m.updateViewportContent()
}

// handleDetailKeys handles keys in detail view
func (m *Model) handleDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Back) {
		m.currentView = ViewList
		m.updateViewportContent()
		return m, nil
	}

	if key.Matches(msg, m.keys.Select) {
		actualIdx := m.getActualIndex(m.cursor)
		if actualIdx >= 0 {
			m.toggleSelection(actualIdx)
			m.updateViewportContent()
		}
	}

	if key.Matches(msg, m.keys.Run) {
		if len(m.selected) > 0 {
			return m, m.executeSelected()
		}
	}

	return m, nil
}

// handleResultsKeys handles keys in results view
func (m *Model) handleResultsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Back) {
		m.currentView = ViewList
		m.results = nil
		m.error = nil
		m.updateViewportContent()
		return m, nil
	}

	if key.Matches(msg, m.keys.Run) {
		if len(m.selected) > 0 {
			return m, m.executeSelected()
		}
	}

	return m, nil
}

// handleHelpKeys handles keys in help view
func (m *Model) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Back) || key.Matches(msg, m.keys.Help) {
		m.currentView = ViewList
		return m, nil
	}
	return m, nil
}

// handleSearchInput handles input when in search mode
func (m *Model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEsc:
		m.searchMode = false
		m.searchInput.Blur()
		m.searchInput.SetValue("")
		m.filteredReqs = m.requests
		m.cursor = 0
		m.buildFilterMap()
		m.updateViewportContent()
		return m, nil

	case tea.KeyEnter:
		m.searchMode = false
		m.searchInput.Blur()
		return m, nil
	}

	m.searchInput, cmd = m.searchInput.Update(msg)
	m.filterRequests()
	m.cursor = 0
	m.buildFilterMap()
	m.updateViewportContent()

	return m, cmd
}

// toggleSelection toggles the selection state of a request
func (m *Model) toggleSelection(index int) {
	if _, ok := m.selected[index]; ok {
		delete(m.selected, index)
	} else {
		m.selected[index] = struct{}{}
	}
}

// executeSelected executes all selected requests
func (m *Model) executeSelected() tea.Cmd {
	if len(m.selected) == 0 {
		return nil
	}

	m.loading = true
	m.error = nil

	var selectedReqs []*types.Request
	for idx := range m.selected {
		if idx < len(m.requests) {
			selectedReqs = append(selectedReqs, m.requests[idx])
		}
	}

	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			results, err := m.executor.ExecuteAll(m.hclFile, selectedReqs, m.overrides)
			return executionCompleteMsg{results: results, err: err}
		},
	)
}

// filterRequests filters requests based on search input
func (m *Model) filterRequests() {
	query := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))
	if query == "" {
		m.filteredReqs = m.requests
		return
	}

	m.filteredReqs = make([]*types.Request, 0)
	for _, req := range m.requests {
		if strings.Contains(strings.ToLower(req.Name), query) ||
			strings.Contains(strings.ToLower(req.Method), query) {
			m.filteredReqs = append(m.filteredReqs, req)
		}
	}
}

// buildFilterMap builds mapping from filtered index to actual index
func (m *Model) buildFilterMap() {
	m.filterMap = make(map[int]int)
	for filteredIdx, req := range m.filteredReqs {
		for actualIdx, originalReq := range m.requests {
			if req == originalReq {
				m.filterMap[filteredIdx] = actualIdx
				break
			}
		}
	}
}

// getActualIndex returns the actual request index from filtered index
func (m *Model) getActualIndex(filteredIdx int) int {
	if actualIdx, ok := m.filterMap[filteredIdx]; ok {
		return actualIdx
	}
	return -1
}

// Message types
type loadFileMsg struct {
	hclFile *types.HCLFile
	err     error
}

type executionCompleteMsg struct {
	results []*types.ExecutionResult
	err     error
}
