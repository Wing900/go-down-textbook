package view

import "github.com/charmbracelet/lipgloss"

var (
	colorText       = lipgloss.Color("#DCDCDC")
	colorTitle      = lipgloss.Color("#10B981")
	colorSelectedBg = lipgloss.Color("#1E3A5F")
	colorSuccess    = lipgloss.Color("#34D399")
	colorAccent     = lipgloss.Color("#FBBF24")
	colorInfo       = lipgloss.Color("#2DD4BF")
	colorDanger     = lipgloss.Color("#F87171")
	colorBorder     = lipgloss.Color("#4B5563")
	colorMuted      = lipgloss.Color("#9CA3AF")
)

type Styles struct {
	Doc, Panel, Title, Subtitle, Muted lipgloss.Style
	Success, Accent, Danger, Info      lipgloss.Style
	SelectedLine, Footer               lipgloss.Style
	Label, Value                       lipgloss.Style
}

func NewStyles() Styles {
	border := lipgloss.NormalBorder()
	return Styles{
		Doc:          lipgloss.NewStyle().Foreground(colorText).Padding(1, 2),
		Panel:        lipgloss.NewStyle().Border(border).BorderForeground(colorBorder).Padding(0, 1),
		Title:        lipgloss.NewStyle().Foreground(colorTitle).Bold(true),
		Subtitle:     lipgloss.NewStyle().Foreground(colorAccent).Bold(true),
		Muted:        lipgloss.NewStyle().Foreground(colorMuted),
		Success:      lipgloss.NewStyle().Foreground(colorSuccess).Bold(true),
		Accent:       lipgloss.NewStyle().Foreground(colorAccent).Bold(true),
		Danger:       lipgloss.NewStyle().Foreground(colorDanger).Bold(true),
		Info:         lipgloss.NewStyle().Foreground(colorInfo),
		SelectedLine: lipgloss.NewStyle().Background(colorSelectedBg).Foreground(lipgloss.Color("#FFFFFF")),
		Footer:       lipgloss.NewStyle().Foreground(colorMuted).BorderTop(true).BorderForeground(colorBorder).PaddingTop(1),
		Label:        lipgloss.NewStyle().Foreground(colorMuted),
		Value:        lipgloss.NewStyle().Foreground(colorText).Bold(true),
	}
}

func ColorAccent() lipgloss.Color { return colorAccent }
func ColorInfo() lipgloss.Color   { return colorInfo }
