package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.page {
	case pageHome:
		return m.updateHomeKeys(msg)
	case pageLogin:
		return m.updateLoginKeys(msg)
	case pageSelect:
		return m.updateSelectKeys(msg)
	case pageDownload:
		return m.updateDownloadKeys(msg)
	default:
		return m, nil
	}
}

func (m *model) updateSpinner(msg spinner.TickMsg) (tea.Model, tea.Cmd) {
	if m.page != pageLogin && !(m.page == pageDownload && m.downloading) {
		return m, nil
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m *model) updateCatalogLoaded(msg catalogLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.requestID != m.loadRequest {
		return m, nil
	}
	if msg.err != nil {
		m.catalogErr = msg.err.Error()
		m.loginStatus = "加载失败"
		m.loginDetail = "可以按 Esc 返回首页后重试。"
		return m, nil
	}
	m.catalog = msg.data
	m.catalogErr = ""
	m.page = pageSelect
	m.focus = focusGrades
	m.setGrades(msg.data.Grades)
	if len(msg.data.Grades) > 0 {
		m.grade = msg.data.Grades[0]
		m.updateSubjectsForGrade()
	}
	m.selectHint = "选择你想下载的教材，按 空格 勾选，按 D 开始下载。"
	return m, nil
}

func (m *model) updateOpenDir(msg openDirMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.statusLine = "打开目录失败: " + msg.err.Error()
	} else {
		m.statusLine = "已尝试打开下载目录。"
	}
	return m, nil
}
