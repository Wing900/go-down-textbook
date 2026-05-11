package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chenyb-go/go-down-textbook/internal/tui/state"
)

func (m *model) updateSubjectsForGrade() {
	if m.catalog == nil {
		return
	}
	subjects := m.service.Subjects(m.catalog.Entries, m.grade)
	m.setSubjects(subjects)
	if len(subjects) > 0 {
		m.subject = subjects[0]
	}
	m.updateBooksForSubject()
}

func (m *model) updateBooksForSubject() {
	if m.catalog == nil || m.grade == "" || m.subject == "" {
		m.books = nil
		return
	}
	m.books = m.service.Books(m.catalog.Entries, m.grade, m.subject)
	m.bookIndex = 0
	m.bookTop = 0
}

func (m *model) visibleBookRows() int {
	return state.VisibleBookRows(m.height)
}

func (m *model) moveBookCursor(key string) {
	if len(m.books) == 0 {
		return
	}
	m.bookTop, m.bookIndex = state.MoveBookCursor(key, m.bookIndex, m.bookTop, len(m.books), m.visibleBookRows())
}

func (m *model) scrollBooks(msg tea.MouseMsg) {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		m.bookTop, m.bookIndex = state.MoveBookCursor("up", m.bookIndex, m.bookTop, len(m.books), m.visibleBookRows())
	case tea.MouseButtonWheelDown:
		m.bookTop, m.bookIndex = state.MoveBookCursor("down", m.bookIndex, m.bookTop, len(m.books), m.visibleBookRows())
	}
}
