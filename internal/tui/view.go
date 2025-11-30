package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"jsonrpc/pkg/types"
)

func (m *Model) renderFileSelectView() string {
	var b strings.Builder

	header := m.styles.HeaderStyle.Render(" Select HCL File ")
	b.WriteString(header)
	b.WriteString("\n")

	// Info line outside viewport
	cwd, _ := os.Getwd()
	intro := fmt.Sprintf("Found %d HCL file(s) in %s", len(m.hclFiles), filepath.Base(cwd))
	b.WriteString(m.styles.InstructionsStyle.Render(intro))
	b.WriteString("\n")

	content := strings.TrimSpace(m.viewport.View())
	if content != "" {
		b.WriteString(content)
		b.WriteString("\n")
	}

	footer := m.styles.FooterStyle.Render("↑/k: up | ↓/j: down | enter: select | q: quit")
	b.WriteString(footer)

	return b.String()
}

func (m *Model) buildFileSelectContent() string {
	if len(m.hclFiles) == 0 {
		return m.styles.InstructionsStyle.Render("No HCL files found in current directory.")
	}

	var b strings.Builder

	for i, file := range m.hclFiles {
		var line strings.Builder

		cursor := " "
		fileStyle := m.styles.ItemNameStyle
		if i == m.fileCursor {
			cursor = m.styles.AccentStyle.Render("→")
			fileStyle = fileStyle.Background(m.styles.CursorColor).Foreground(lipgloss.Color("0"))
		}

		line.WriteString(cursor)
		line.WriteString(" ")
		line.WriteString(fileStyle.Render(file))

		fileInfo, err := os.Stat(file)
		if err == nil {
			size := fileInfo.Size()
			sizeStr := formatFileSize(size)
			line.WriteString("  ")
			line.WriteString(m.styles.InstructionsStyle.Render(fmt.Sprintf("(%s)", sizeStr)))
		}

		b.WriteString(line.String())
		b.WriteString("\n")
	}

	return b.String()
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func (m *Model) renderListView() string {
	var b strings.Builder

	header := m.styles.HeaderStyle.Render(" JSON-RPC Interactive TUI ")
	b.WriteString(header)
	b.WriteString("\n")

	if m.searchMode {
		b.WriteString(m.styles.SearchStyle.Render("🔍 Search: " + m.searchInput.View()))
		b.WriteString("\n")
	}

	// Status line outside viewport
	selectedCount := len(m.selected)
	totalCount := len(m.requests)
	filteredCount := len(m.filteredReqs)
	statusLine := fmt.Sprintf("Total: %d | Filtered: %d | Selected: %d", totalCount, filteredCount, selectedCount)
	b.WriteString(m.styles.InstructionsStyle.Render(statusLine))
	b.WriteString("\n")

	content := strings.TrimSpace(m.viewport.View())
	if content != "" {
		b.WriteString(content)
		b.WriteString("\n")
	}

	footer := m.buildFooter()
	b.WriteString(footer)

	if m.error != nil {
		b.WriteString("\n")
		b.WriteString(m.styles.ErrorStyle.Render("⚠ Error: " + m.error.Error()))
	}

	return b.String()
}

func (m *Model) buildListContent() string {
	if len(m.filteredReqs) == 0 {
		return m.renderEmptyState()
	}

	var b strings.Builder

	for i, req := range m.filteredReqs {
		actualIdx := m.getActualIndex(i)
		line := m.renderRequestLine(i, actualIdx, req)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

func (m *Model) renderRequestLine(filteredIdx, actualIdx int, req *types.Request) string {
	var parts []string

	if _, selected := m.selected[actualIdx]; selected {
		parts = append(parts, m.styles.SelectedStyle.Render("◉"))
	} else {
		parts = append(parts, "○")
	}

	cursor := " "
	nameStyle := m.styles.ItemNameStyle
	if filteredIdx == m.cursor {
		cursor = m.styles.AccentStyle.Render("→")
		nameStyle = nameStyle.Background(m.styles.CursorColor).Foreground(lipgloss.Color("0"))
	}

	parts = append(parts, cursor)
	parts = append(parts, nameStyle.Render(req.Name))
	parts = append(parts, m.styles.MethodStyle.Render("("+req.Method+")"))

	if req.Config != "" {
		parts = append(parts, m.styles.ConfigStyle.Render("["+req.Config+"]"))
	} else if req.URL != "" {
		parts = append(parts, m.styles.URLStyle.Render("[custom-url]"))
	}

	return strings.Join(parts, " ")
}

func (m *Model) renderDetailView() string {
	var b strings.Builder

	if len(m.filteredReqs) == 0 || m.cursor >= len(m.filteredReqs) {
		return "No request selected"
	}

	req := m.filteredReqs[m.cursor]
	actualIdx := m.getActualIndex(m.cursor)

	header := m.styles.HeaderStyle.Render(fmt.Sprintf(" Request Details: %s ", req.Name))
	b.WriteString(header)
	b.WriteString("\n")

	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	footer := m.buildFooter()
	b.WriteString(footer)

	if _, selected := m.selected[actualIdx]; selected {
		b.WriteString(" ")
		b.WriteString(m.styles.SelectedStyle.Render("[SELECTED]"))
	}

	return b.String()
}

func (m *Model) buildDetailContent() string {
	if len(m.filteredReqs) == 0 || m.cursor >= len(m.filteredReqs) {
		return "No request selected"
	}

	req := m.filteredReqs[m.cursor]
	var b strings.Builder

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.BorderColor).
		Padding(1, 2).
		Width(m.width - 4)

	var details []string

	details = append(details, m.renderDetailField("Name", req.Name, m.styles.ItemNameStyle))
	details = append(details, m.renderDetailField("Method", req.Method, m.styles.MethodStyle))

	if req.Config != "" {
		if config, exists := m.hclFile.Configs[req.Config]; exists {
			details = append(details, m.renderDetailField("Config", req.Config, m.styles.ConfigStyle))
			details = append(details, m.renderDetailField("URL", config.URL, m.styles.ValueStyle))
			details = append(details, m.renderDetailField("Timeout", fmt.Sprintf("%ds", config.Timeout), m.styles.ValueStyle))

			if len(config.Headers) > 0 {
				details = append(details, "")
				details = append(details, m.styles.SectionHeader.Render("Headers:"))
				for k, v := range config.Headers {
					details = append(details, fmt.Sprintf("  %s: %s",
						m.styles.KeyStyle.Render(k),
						m.styles.ValueStyle.Render(v)))
				}
			}
		}
	} else if req.URL != "" {
		details = append(details, m.renderDetailField("URL", req.URL, m.styles.URLStyle))
	}

	if len(req.Headers) > 0 {
		details = append(details, "")
		details = append(details, m.styles.SectionHeader.Render("Request Headers:"))
		for k, v := range req.Headers {
			details = append(details, fmt.Sprintf("  %s: %s",
				m.styles.KeyStyle.Render(k),
				m.styles.ValueStyle.Render(v)))
		}
	}

	if req.Timeout > 0 {
		details = append(details, m.renderDetailField("Timeout", fmt.Sprintf("%ds", req.Timeout), m.styles.ValueStyle))
	}

	if req.ProcessedParams != nil {
		details = append(details, "")
		details = append(details, m.styles.SectionHeader.Render("Parameters:"))
		paramsJSON, err := json.MarshalIndent(req.ProcessedParams, "  ", "  ")
		if err == nil {
			details = append(details, "  "+m.highlightJSON(string(paramsJSON)))
		} else {
			details = append(details, fmt.Sprintf("  %v", req.ProcessedParams))
		}
	}

	b.WriteString(boxStyle.Render(strings.Join(details, "\n")))

	return b.String()
}

func (m *Model) renderDetailField(key, value string, valueStyle lipgloss.Style) string {
	return fmt.Sprintf("%s: %s",
		m.styles.KeyStyle.Render(key),
		valueStyle.Render(value))
}

func (m *Model) renderResultsView() string {
	var b strings.Builder

	header := m.styles.HeaderStyle.Render(" Execution Results ")
	b.WriteString(header)
	b.WriteString("\n")

	if m.loading {
		b.WriteString("\n")
		b.WriteString(m.spinner.View() + " Executing requests...\n")
		return b.String()
	}

	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	footer := m.buildFooter()
	b.WriteString(footer)

	return b.String()
}

func (m *Model) buildResultsContent() string {
	var b strings.Builder

	if m.error != nil {
		errorBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.styles.ErrorColor).
			Foreground(m.styles.ErrorColor).
			Padding(1, 2).
			Width(m.width - 4).
			Render("⚠ Error: " + m.error.Error())

		b.WriteString(errorBox)
		return b.String()
	}

	if len(m.results) == 0 {
		return "No results to display"
	}

	successCount := 0
	for _, result := range m.results {
		if result.IsSuccess() {
			successCount++
		}
	}

	summary := fmt.Sprintf("Executed %d requests: %d succeeded, %d failed\n\n",
		len(m.results), successCount, len(m.results)-successCount)
	b.WriteString(m.styles.SectionHeader.Render(summary))

	for i, result := range m.results {
		b.WriteString(m.renderResult(i, result))
		b.WriteString("\n\n")
	}

	if len(m.history) > 1 {
		b.WriteString("\n")
		b.WriteString(m.styles.InstructionsStyle.Render(
			fmt.Sprintf("History: %d executions recorded", len(m.history))))
	}

	return b.String()
}

func (m *Model) renderResult(index int, result *types.ExecutionResult) string {
	var b strings.Builder

	var icon string
	var statusStyle lipgloss.Style
	if result.IsSuccess() {
		icon = m.styles.SelectedStyle.Render("✓")
		statusStyle = m.styles.SelectedStyle
	} else {
		icon = m.styles.ErrorStyle.Render("✗")
		statusStyle = m.styles.ErrorStyle
	}

	title := fmt.Sprintf("%s %s (%s)",
		icon,
		m.styles.ItemNameStyle.Render(result.Request.Name),
		m.styles.MethodStyle.Render(result.Request.Method))

	b.WriteString(title)
	b.WriteString("\n")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(statusStyle.GetForeground()).
		Padding(1, 2).
		Width(m.width - 4)

	var details []string

	if result.Error != nil {
		details = append(details, m.styles.ErrorStyle.Render("Error: "+result.Error.Error()))
	} else {
		details = append(details, m.styles.KeyStyle.Render("Status: ")+statusStyle.Render("Success"))

		if result.Duration > 0 {
			details = append(details, m.renderDetailField("Response Time",
				fmt.Sprintf("%dms", result.Duration.Milliseconds()),
				m.styles.ValueStyle))
		}

		if result.Response != nil {
			details = append(details, "")
			details = append(details, m.styles.SectionHeader.Render("Response:"))

			responseJSON, err := json.MarshalIndent(result.Response, "", "  ")
			if err == nil {
				details = append(details, m.highlightJSON(string(responseJSON)))
			} else {
				details = append(details, fmt.Sprintf("%v", result.Response))
			}
		}
	}

	b.WriteString(boxStyle.Render(strings.Join(details, "\n")))

	return b.String()
}

func (m *Model) renderHelpView() string {
	var b strings.Builder

	header := m.styles.HeaderStyle.Render(" Help - Keyboard Shortcuts ")
	b.WriteString(header)
	b.WriteString("\n\n")

	helpContent := m.help.View(m.keys)
	b.WriteString(helpContent)

	b.WriteString("\n\n")
	b.WriteString(m.styles.InstructionsStyle.Render("Press ? or ESC to close help"))

	return b.String()
}

func (m *Model) buildFooter() string {
	var parts []string

	if m.searchMode {
		parts = append(parts, "ESC: exit search")
	} else {
		switch m.currentView {
		case ViewList:
			parts = append(parts, "?: help")
			parts = append(parts, "/: search")
			parts = append(parts, "space: select")
			parts = append(parts, "r: run")
		case ViewDetail:
			parts = append(parts, "ESC: back")
			parts = append(parts, "space: select")
			parts = append(parts, "r: run")
		case ViewResults:
			parts = append(parts, "ESC: back")
			parts = append(parts, "r: rerun")
		}
		parts = append(parts, "q: quit")
	}

	return m.styles.FooterStyle.Render(strings.Join(parts, " | "))
}

func (m *Model) renderEmptyState() string {
	emptyMsg := `
  No requests found

  Make sure your HCL file contains request blocks.
`
	return m.styles.InstructionsStyle.Render(emptyMsg)
}

func (m *Model) highlightJSON(jsonStr string) string {
	var b strings.Builder
	indent := 0

	for i := 0; i < len(jsonStr); i++ {
		ch := jsonStr[i]

		switch ch {
		case '{', '[':
			indent = m.handleOpenBracket(&b, jsonStr, i, indent, ch)
		case '}', ']':
			indent = m.handleCloseBracket(&b, jsonStr, i, indent, ch)
		case ',':
			m.handleComma(&b, indent)
		case '"':
			i = m.handleString(&b, jsonStr, i)
		case ':':
			b.WriteString(": ")
		case ' ', '\t', '\n', '\r':
			continue
		default:
			i = m.handleValue(&b, jsonStr, i)
		}
	}

	return b.String()
}

func (m *Model) handleOpenBracket(b *strings.Builder, jsonStr string, i, indent int, ch byte) int {
	b.WriteRune(rune(ch))
	if i+1 < len(jsonStr) && (jsonStr[i+1] == '}' || jsonStr[i+1] == ']') {
		return indent
	}
	indent++
	b.WriteRune('\n')
	b.WriteString(strings.Repeat("  ", indent))
	return indent
}

func (m *Model) handleCloseBracket(b *strings.Builder, jsonStr string, i, indent int, ch byte) int {
	if jsonStr[i-1] != '{' && jsonStr[i-1] != '[' {
		b.WriteRune('\n')
		indent--
		b.WriteString(strings.Repeat("  ", indent))
	}
	b.WriteRune(rune(ch))
	return indent
}

func (m *Model) handleComma(b *strings.Builder, indent int) {
	b.WriteRune(',')
	b.WriteRune('\n')
	b.WriteString(strings.Repeat("  ", indent))
}

func (m *Model) handleString(b *strings.Builder, jsonStr string, i int) int {
	start := i
	i++
	for i < len(jsonStr) && jsonStr[i] != '"' {
		if jsonStr[i] == '\\' {
			i++
		}
		i++
	}
	if i < len(jsonStr) {
		i++
	}

	str := jsonStr[start:i]
	if i < len(jsonStr) && jsonStr[i] == ':' {
		b.WriteString(m.styles.JSONKeyStyle.Render(str))
	} else {
		b.WriteString(m.styles.JSONStringStyle.Render(str))
	}
	return i - 1
}

func (m *Model) handleValue(b *strings.Builder, jsonStr string, i int) int {
	start := i
	for i < len(jsonStr) && !strings.ContainsRune("{[]},:\" \t\n\r", rune(jsonStr[i])) {
		i++
	}
	value := jsonStr[start:i]

	switch value {
	case "true", "false":
		b.WriteString(m.styles.JSONBoolStyle.Render(value))
	case "null":
		b.WriteString(m.styles.JSONNullStyle.Render(value))
	default:
		b.WriteString(m.styles.JSONNumberStyle.Render(value))
	}
	return i - 1
}
