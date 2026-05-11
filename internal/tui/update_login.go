package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *model) updateLoginKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.loadRequest++
		m.page = pageHome
		m.refreshHome()
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}
