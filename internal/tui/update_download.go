package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *model) updateDownloadKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		if m.cancelDownload != nil && m.downloading {
			m.cancelDownload()
		}
		return m, tea.Quit
	case "o":
		return m, openDirCmd(m.service)
	case "enter":
		if !m.downloading {
			m.page = pageSelect
			m.downloadDone = false
			m.statusLine = "可以继续选书下载。"
			m.refreshHome()
		}
	case "r":
		if !m.downloading && len(m.failed) > 0 {
			m.selectHint = "失败项已经保留在日志里，请回选书页重新勾选后下载。"
			m.page = pageSelect
		}
	}
	return m, nil
}
