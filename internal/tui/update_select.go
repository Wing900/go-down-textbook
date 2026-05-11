package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *model) updateSelectKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.page = pageHome
		m.refreshHome()
		return m, nil
	case "tab":
		m.focus = (m.focus + 1) % 3
		return m, nil
	case "shift+tab":
		m.focus = (m.focus + 2) % 3
		return m, nil
	case " ":
		if m.focus == focusBooks {
			m.toggleCurrentBook()
		}
		return m, nil
	case "a":
		if m.focus == focusBooks {
			m.toggleAllVisibleBooks()
		}
		return m, nil
	case "d":
		return m, m.startDownloadFromSelection()
	}
	return m.routeSelectMovement(msg)
}

func (m *model) startDownloadFromSelection() tea.Cmd {
	if len(m.selected) == 0 || m.catalog == nil {
		m.selectHint = "还没选书，先用 空格 勾选几本吧。"
		return nil
	}

	books := selectedBooks(m.selected)
	m.startDownload(books)
	m.downloadCh = m.service.StartDownload(m.downloadCtx, m.catalog.Token, books)
	return tea.Batch(m.spinner.Tick, waitDownloadCmd(m.downloadCh))
}

func (m *model) routeSelectMovement(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.focus {
	case focusGrades:
		return m.updateGradeList(msg)
	case focusSubjects:
		return m.updateSubjectList(msg)
	case focusBooks:
		m.moveBookCursor(msg.String())
	}
	return m, nil
}
