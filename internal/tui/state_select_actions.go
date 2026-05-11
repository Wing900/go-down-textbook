package tui

import "github.com/chenyb-go/go-down-textbook/internal/tui/state"

func (m *model) syncBookWindow() {
	m.bookTop = state.SyncBookWindow(m.bookIndex, m.bookTop, m.visibleBookRows())
}

func (m *model) toggleCurrentBook() {
	state.ToggleCurrentBook(m.selected, m.books, m.bookIndex)
}

func (m *model) toggleAllVisibleBooks() {
	state.ToggleAllVisibleBooks(m.selected, m.books)
}
