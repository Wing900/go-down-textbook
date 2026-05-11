package view

import (
	"fmt"
	"strings"

	"github.com/chenyb-go/go-down-textbook/internal/app"
	"github.com/chenyb-go/go-down-textbook/internal/models"
	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
)

func RenderBooks(s Styles, books []app.SelectableBook, selected map[string]models.BookItem, bookIndex, bookTop, rows, width int) string {
	lines := []string{s.Subtitle.Render("教材")}
	if len(books) == 0 {
		lines = append(lines, "", s.Muted.Render("当前条件下没有可显示的教材。"))
		return strings.Join(lines, "\n")
	}
	end := common.MinInt(len(books), bookTop+rows)
	for i := bookTop; i < end; i++ {
		lines = append(lines, renderBookLine(s, books[i], selected, i == bookIndex, width))
	}
	lines = append(lines, "", s.Muted.Render(fmt.Sprintf("显示 %d-%d / %d", bookTop+1, end, len(books))))
	return strings.Join(lines, "\n")
}

func renderBookLine(s Styles, book app.SelectableBook, selected map[string]models.BookItem, active bool, width int) string {
	check := "[ ]"
	if _, ok := selected[book.ID]; ok {
		check = s.Success.Render("[x]")
	}
	status := ""
	if book.Downloaded {
		status = "  " + s.Muted.Render("已下载")
	}
	line := fmt.Sprintf("%s %s%s", check, common.TruncateText(book.Title, width-12), status)
	if active {
		return s.SelectedLine.Width(width - 2).Render(line)
	}
	return line
}
