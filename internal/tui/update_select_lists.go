package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	u "github.com/chenyb-go/go-down-textbook/internal/tui/update"
)

func (m *model) updateGradeList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	prev := m.grade
	var cmd tea.Cmd
	m.grades, cmd = m.grades.Update(msg)
	if title := u.SelectedItemTitle(m.grades.SelectedItem()); title != "" {
		m.grade = title
		if m.grade != prev {
			m.updateSubjectsForGrade()
		}
	}
	return m, cmd
}

func (m *model) updateSubjectList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	prev := m.subject
	var cmd tea.Cmd
	m.subjects, cmd = m.subjects.Update(msg)
	if title := u.SelectedItemTitle(m.subjects.SelectedItem()); title != "" {
		m.subject = title
		if m.subject != prev {
			m.updateBooksForSubject()
		}
	}
	return m, cmd
}

func (m *model) updateSelectLists(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.grades, cmd = m.grades.Update(msg)
	cmds = append(cmds, cmd)
	m.subjects, cmd = m.subjects.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
