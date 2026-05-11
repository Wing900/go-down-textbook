package state

import (
	"github.com/chenyb-go/go-down-textbook/internal/app"
	"github.com/chenyb-go/go-down-textbook/internal/models"
	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
)

func VisibleBookRows(height int) int {
	rows := common.MinInt(5, height-20)
	return common.MaxInt(4, rows)
}

func MoveBookCursor(key string, bookIndex, bookTop, bookCount, rows int) (int, int) {
	if bookCount == 0 {
		return 0, 0
	}
	switch key {
	case "up", "k":
		if bookIndex > 0 {
			bookIndex--
		}
	case "down", "j":
		if bookIndex < bookCount-1 {
			bookIndex++
		}
	case "pgdown":
		bookIndex = common.MinInt(bookCount-1, bookIndex+rows)
	case "pgup":
		bookIndex = common.MaxInt(0, bookIndex-rows)
	}
	return SyncBookWindow(bookIndex, bookTop, rows), bookIndex
}

func SyncBookWindow(bookIndex, bookTop, rows int) int {
	if bookIndex < bookTop {
		return bookIndex
	}
	if bookIndex >= bookTop+rows {
		return bookIndex - rows + 1
	}
	return bookTop
}

func ToggleCurrentBook(selected map[string]models.BookItem, books []app.SelectableBook, bookIndex int) {
	if len(books) == 0 || bookIndex >= len(books) {
		return
	}
	book := books[bookIndex]
	if _, ok := selected[book.ID]; ok {
		delete(selected, book.ID)
		return
	}
	selected[book.ID] = book.BookItem
}

func ToggleAllVisibleBooks(selected map[string]models.BookItem, books []app.SelectableBook) {
	if len(books) == 0 {
		return
	}
	allSelected := true
	for _, book := range books {
		if _, ok := selected[book.ID]; !ok {
			allSelected = false
			break
		}
	}
	for _, book := range books {
		if allSelected {
			delete(selected, book.ID)
		} else {
			selected[book.ID] = book.BookItem
		}
	}
}
