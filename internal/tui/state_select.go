package tui

import (
	"github.com/chenyb-go/go-down-textbook/internal/models"
	"github.com/chenyb-go/go-down-textbook/internal/tui/state"
)

func (m *model) setGrades(grades []string) {
	m.grades.SetItems(state.ToListItems(grades))
	if len(grades) > 0 {
		m.grades.Select(0)
	}
}

func (m *model) setSubjects(subjects []string) {
	m.subjects.SetItems(state.ToListItems(subjects))
	if len(subjects) > 0 {
		m.subjects.Select(0)
	}
}

func selectedBooks(selected map[string]models.BookItem) []models.BookItem {
	return state.SelectedBooks(selected)
}

func joinSelectedBooks(selected map[string]models.BookItem, width int) string {
	return state.JoinSelectedBooks(selected, width)
}
