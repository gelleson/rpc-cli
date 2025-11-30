package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles holds the visual styling configuration for the TUI
type Styles struct {
	// Colors
	PrimaryColor   lipgloss.Color
	SecondaryColor lipgloss.Color
	AccentColor    lipgloss.Color
	SuccessColor   lipgloss.Color
	ErrorColor     lipgloss.Color
	WarningColor   lipgloss.Color
	TextColor      lipgloss.Color
	BorderColor    lipgloss.Color
	HeaderColor    lipgloss.Color
	CursorColor    lipgloss.Color

	// Text styles
	InstructionsStyle lipgloss.Style
	SectionHeader     lipgloss.Style
	ItemNameStyle     lipgloss.Style
	MethodStyle       lipgloss.Style
	ConfigStyle       lipgloss.Style
	URLStyle          lipgloss.Style
	SelectedStyle     lipgloss.Style
	ErrorStyle        lipgloss.Style
	AccentStyle       lipgloss.Style
	HeaderStyle       lipgloss.Style
	FooterStyle       lipgloss.Style
	SearchStyle       lipgloss.Style
	HelpStyle         lipgloss.Style
	KeyStyle          lipgloss.Style
	ValueStyle        lipgloss.Style
	JSONKeyStyle      lipgloss.Style
	JSONStringStyle   lipgloss.Style
	JSONNumberStyle   lipgloss.Style
	JSONBoolStyle     lipgloss.Style
	JSONNullStyle     lipgloss.Style
}

// DefaultStyles returns a new Styles instance with sensible defaults
func DefaultStyles() *Styles {
	return &Styles{
		PrimaryColor:   lipgloss.Color("69"),
		SecondaryColor: lipgloss.Color("240"),
		AccentColor:    lipgloss.Color("86"),
		SuccessColor:   lipgloss.Color("42"),
		ErrorColor:     lipgloss.Color("196"),
		WarningColor:   lipgloss.Color("208"),
		TextColor:      lipgloss.Color("15"),
		BorderColor:    lipgloss.Color("59"),
		HeaderColor:    lipgloss.Color("21"),
		CursorColor:    lipgloss.Color("226"),

		InstructionsStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("242")),

		SectionHeader: lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")).
			Bold(true).
			Underline(true),

		ItemNameStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true),

		MethodStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Italic(true),

		ConfigStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true),

		URLStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("147")).
			Italic(true),

		SelectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true),

		ErrorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true),

		AccentStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true),

		HeaderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("21")).
			Bold(true).
			Padding(0, 1),

		FooterStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("242")),

		SearchStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("237")).
			Padding(0, 1),

		HelpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("242")),

		KeyStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true),

		ValueStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("147")),

		JSONKeyStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")),

		JSONStringStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")),

		JSONNumberStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")),

		JSONBoolStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")),

		JSONNullStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
	}
}

// DarkStyles returns styles optimized for dark terminal backgrounds
func DarkStyles() *Styles {
	styles := DefaultStyles()
	styles.TextColor = lipgloss.Color("251")
	styles.BorderColor = lipgloss.Color("59")
	styles.InstructionsStyle = styles.InstructionsStyle.
		Foreground(lipgloss.Color("245"))
	return styles
}

// LightStyles returns styles optimized for light terminal backgrounds
func LightStyles() *Styles {
	styles := DefaultStyles()
	styles.TextColor = lipgloss.Color("16")
	styles.BorderColor = lipgloss.Color("249")
	styles.InstructionsStyle = styles.InstructionsStyle.
		Foreground(lipgloss.Color("240"))
	styles.PrimaryColor = lipgloss.Color("26")
	styles.SecondaryColor = lipgloss.Color("243")
	styles.CursorColor = lipgloss.Color("226")
	return styles
}

// HighContrastStyles returns styles optimized for accessibility
func HighContrastStyles() *Styles {
	styles := DefaultStyles()
	styles.PrimaryColor = lipgloss.Color("15")
	styles.SecondaryColor = lipgloss.Color("7")
	styles.TextColor = lipgloss.Color("15")
	styles.CursorColor = lipgloss.Color("11")
	styles.SuccessColor = lipgloss.Color("10")
	styles.ErrorColor = lipgloss.Color("9")
	styles.WarningColor = lipgloss.Color("14")
	styles.BorderColor = lipgloss.Color("15")
	styles.HeaderColor = lipgloss.Color("15")
	return styles
}
