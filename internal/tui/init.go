package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
	"github.com/chenyb-go/go-down-textbook/internal/tui/view"
)

func newSpinner() spinner.Model {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(view.ColorAccent())
	return spin
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) initLists() {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	m.grades = list.New(nil, delegate, 20, 10)
	m.subjects = list.New(nil, delegate, 20, 10)
	setupList(&m.grades, "年级")
	setupList(&m.subjects, "学科")
}

func setupList(target *list.Model, title string) {
	target.Title = title
	target.SetFilteringEnabled(false)
	target.SetShowStatusBar(false)
	target.SetShowHelp(false)
	target.SetShowPagination(false)
}

func (m *model) resize() {
	if m.width == 0 || m.height == 0 {
		return
	}
	listHeight := common.MaxInt(8, m.height-18)
	m.grades.SetSize(common.MaxInt(18, m.width/5), listHeight)
	m.subjects.SetSize(common.MaxInt(16, m.width/5), listHeight)
}
