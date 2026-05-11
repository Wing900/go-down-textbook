package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	u "github.com/chenyb-go/go-down-textbook/internal/tui/update"
)

func (m *model) handleHomeClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	index := u.HomeMenuIndex(msg.Y, len(m.menuItems), m.homeErr != "")
	if index < 0 {
		return m, nil
	}
	m.menuIndex = index
	return m.runHomeAction()
}
