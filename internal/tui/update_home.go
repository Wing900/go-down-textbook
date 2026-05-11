package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *model) updateHomeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case "down", "j":
		if m.menuIndex < len(m.menuItems)-1 {
			m.menuIndex++
		}
	case "enter":
		return m.runHomeAction()
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m *model) runHomeAction() (tea.Model, tea.Cmd) {
	switch m.menuItems[m.menuIndex].action {
	case actionStart:
		m.page = pageLogin
		m.loadRequest++
		m.loginStatus = "正在准备选书流程..."
		m.loginDetail = "这一步会检查登录状态，并动态获取目录。"
		m.catalogErr = ""
		return m, tea.Batch(m.spinner.Tick, loadCatalogCmd(m.service, m.loadRequest))
	case actionOpenDir:
		return m, openDirCmd(m.service)
	case actionQuit:
		return m, tea.Quit
	default:
		return m, nil
	}
}
