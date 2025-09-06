package styles

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Background  lipgloss.Color
	Surface     lipgloss.Color
	Error       lipgloss.Color
	Warning     lipgloss.Color
	Success     lipgloss.Color
	Text        lipgloss.Color
	TextMuted   lipgloss.Color
	Border      lipgloss.Color
	BorderMuted lipgloss.Color
}

var DefaultTheme = Theme{
	Primary:     lipgloss.Color("62"),
	Secondary:   lipgloss.Color("205"),
	Background:  lipgloss.Color("236"),
	Surface:     lipgloss.Color("238"),
	Error:       lipgloss.Color("196"),
	Warning:     lipgloss.Color("214"),
	Success:     lipgloss.Color("82"),
	Text:        lipgloss.Color("230"),
	TextMuted:   lipgloss.Color("248"),
	Border:      lipgloss.Color("62"),
	BorderMuted: lipgloss.Color("241"),
}

type BaseStyles struct {
	App       lipgloss.Style
	Header    lipgloss.Style
	Footer    lipgloss.Style
	Container lipgloss.Style
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	Text      lipgloss.Style
	TextMuted lipgloss.Style
	Error     lipgloss.Style
	Warning   lipgloss.Style
	Success   lipgloss.Style
	Border    lipgloss.Style
	Key       lipgloss.Style
}

func NewBaseStyles(theme Theme) BaseStyles {
	return BaseStyles{
		App: lipgloss.NewStyle().
			Background(lipgloss.Color("0")),

		Header: lipgloss.NewStyle().
			Background(theme.Primary).
			Foreground(theme.Text).
			Padding(0, 1),

		Footer: lipgloss.NewStyle().
			Background(theme.Background).
			Foreground(theme.TextMuted).
			Padding(0, 1),

		Container: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(1, 2),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Secondary),

		Subtitle: lipgloss.NewStyle().
			Foreground(theme.Text),

		Text: lipgloss.NewStyle().
			Foreground(theme.Text),

		TextMuted: lipgloss.NewStyle().
			Foreground(theme.TextMuted),

		Error: lipgloss.NewStyle().
			Foreground(theme.Error).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(theme.Warning),

		Success: lipgloss.NewStyle().
			Foreground(theme.Success),

		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border),

		Key: lipgloss.NewStyle().
			Foreground(theme.Secondary).
			Bold(true),
	}
}

func GetDefaultStyles() BaseStyles {
	return NewBaseStyles(DefaultTheme)
}
