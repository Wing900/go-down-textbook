package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	u "github.com/chenyb-go/go-down-textbook/internal/tui/update"
)

func (m *model) updateMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch m.page {
	case pageHome:
		return m.handleHomeMouse(msg)
	case pageSelect:
		return m.handleSelectMouse(msg)
	case pageDownload:
		return m.handleDownloadMouse(msg)
	case pageLogin:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft && m.mouseInFooter(msg.Y) {
			m.loadRequest++
			m.page = pageHome
			m.refreshHome()
		}
	}
	return m, nil
}

func (m *model) mouseInFooter(y int) bool {
	return y >= m.height-3
}

func (m *model) handleHomeMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if u.IsLeftClick(msg.String(), mouseButtonName(msg.Button)) {
		return m.handleHomeClick(msg)
	}
	return m, nil
}

func (m *model) handleDownloadMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if !u.IsLeftClick(msg.String(), mouseButtonName(msg.Button)) || !m.mouseInFooter(msg.Y) {
		return m, nil
	}
	switch u.DownloadFooterAction(msg.X, m.downloadDone, len(m.failed) > 0) {
	case u.FooterOpen:
		return m, openDirCmd(m.service)
	case u.FooterContinue:
		if !m.downloading {
			m.page = pageSelect
			m.downloadDone = false
			m.statusLine = "可以继续选书下载。"
			m.refreshHome()
		}
	case u.FooterRetry:
		if !m.downloading && len(m.failed) > 0 {
			m.selectHint = "失败项已经保留在日志里，请回选书页重新勾选后下载。"
			m.page = pageSelect
		}
	}
	return m, nil
}

func mouseButtonName(button tea.MouseButton) string {
	switch button {
	case tea.MouseButtonLeft:
		return "left"
	default:
		return ""
	}
}
