package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	u "github.com/chenyb-go/go-down-textbook/internal/tui/update"
)

func (m *model) handleSelectMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown {
		return m.handleSelectWheel(msg)
	}
	if u.IsLeftClick(msg.String(), mouseButtonName(msg.Button)) {
		return m.handleSelectClick(msg)
	}
	return m, nil
}

func (m *model) handleSelectClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch u.SelectArea(msg.X, msg.Y, m.width, m.visibleBookRows()+3) {
	case u.AreaGrades:
		m.focus = focusGrades
		m.grades, _ = m.grades.Update(msg)
		m.syncGradeSelection()
	case u.AreaSubjects:
		m.focus = focusSubjects
		m.subjects, _ = m.subjects.Update(msg)
		m.syncSubjectSelection()
	case u.AreaBooks:
		if m.clickBook(msg) {
			return m, nil
		}
	case u.AreaFooter:
		switch u.SelectFooterAction(msg.X) {
		case u.FooterStart:
			return m, m.startDownloadFromSelection()
		case u.FooterBack:
			m.page = pageHome
			m.refreshHome()
		}
	}
	return m, nil
}

func (m *model) handleSelectWheel(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch u.SelectArea(msg.X, msg.Y, m.width, m.visibleBookRows()+3) {
	case u.AreaGrades:
		m.focus = focusGrades
		m.grades, _ = m.grades.Update(msg)
		m.syncGradeSelection()
	case u.AreaSubjects:
		m.focus = focusSubjects
		m.subjects, _ = m.subjects.Update(msg)
		m.syncSubjectSelection()
	case u.AreaBooks:
		m.focus = focusBooks
		m.scrollBooks(msg)
	}
	return m, nil
}

func (m *model) clickBook(msg tea.MouseMsg) bool {
	index := u.BookIndexByClick(msg.X, msg.Y, m.width, m.visibleBookRows(), m.bookTop, len(m.books))
	if index < 0 {
		return false
	}
	m.focus = focusBooks
	m.bookIndex = index
	m.toggleCurrentBook()
	return true
}

func (m *model) syncGradeSelection() {
	if title := u.SelectedItemTitle(m.grades.SelectedItem()); title != "" && title != m.grade {
		m.grade = title
		m.updateSubjectsForGrade()
	}
}

func (m *model) syncSubjectSelection() {
	if title := u.SelectedItemTitle(m.subjects.SelectedItem()); title != "" && title != m.subject {
		m.subject = title
		m.updateBooksForSubject()
	}
}
