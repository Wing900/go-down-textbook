package state

import (
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/chenyb-go/go-down-textbook/internal/models"
	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
)

func ToListItems(values []string) []list.Item {
	items := make([]list.Item, 0, len(values))
	for _, value := range values {
		items = append(items, listItem(value))
	}
	return items
}

type listItem string

func (i listItem) FilterValue() string { return string(i) }
func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "" }

func SelectedBooks(selected map[string]models.BookItem) []models.BookItem {
	books := make([]models.BookItem, 0, len(selected))
	for _, book := range selected {
		books = append(books, book)
	}
	return books
}

func JoinSelectedBooks(selected map[string]models.BookItem, width int) string {
	if len(selected) == 0 {
		return "还没选书，试着用方向键浏览一下吧。"
	}
	names := make([]string, 0, len(selected))
	for _, book := range selected {
		names = append(names, book.Title)
	}
	sort.Strings(names)
	return common.TruncateText(strings.Join(names, "  |  "), width)
}

func PushLog(logs []string, line string, maxLen int) []string {
	timestamp := time.Now().Format("15:04:05")
	logs = append([]string{timestamp + "  " + line}, logs...)
	if len(logs) > maxLen {
		logs = logs[:maxLen]
	}
	return logs
}
