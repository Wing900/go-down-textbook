package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
)

type SelectData struct {
	Width, Focus, ListHeight                                          int
	Grade, Subject, Guide, GradesView, SubjectsView, BooksView, Queue string
	SelectedCount                                                     int
}

func RenderSelect(s Styles, d SelectData) string {
	top := s.Panel.Width(common.MaxInt(42, d.Width-6)).Render(strings.Join([]string{
		s.Title.Render("选择教材"),
		fmt.Sprintf("年级: [ %s ]    学科: [ %s ]", common.FallbackText(d.Grade), common.FallbackText(d.Subject)),
		d.Guide,
	}, "\n"))
	columns := lipgloss.JoinHorizontal(lipgloss.Top,
		renderListPanel(s, "年级", d.GradesView, d.Focus == 0, d.Width/5, d.ListHeight),
		renderListPanel(s, "学科", d.SubjectsView, d.Focus == 1, d.Width/5, d.ListHeight),
		renderBooksPanel(s, d.BooksView, d.Focus == 2, d.Width, d.ListHeight),
	)
	queue := s.Panel.Width(common.MaxInt(42, d.Width-6)).Render(strings.Join([]string{
		s.Subtitle.Render(fmt.Sprintf("待下载队列 (%d 本)", d.SelectedCount)),
		d.Queue, "", s.Accent.Render(fmt.Sprintf("已选 %d 本，准备下载", d.SelectedCount)),
	}, "\n"))
	footer := s.Footer.Width(common.MaxInt(42, d.Width-6)).Render("Tab 切换   ↑↓ 移动   空格勾选   A 全选   D 开始下载   Esc 返回")
	return s.Doc.Render(strings.Join([]string{top, columns, queue, footer}, "\n"))
}

func renderListPanel(s Styles, title, content string, focused bool, width, height int) string {
	style := s.Panel.Width(common.MaxInt(18, width)).Height(height)
	if focused {
		style = style.BorderForeground(ColorInfo())
	}
	return style.Render(s.Subtitle.Render(title) + "\n" + content)
}

func renderBooksPanel(s Styles, body string, focused bool, width, height int) string {
	style := s.Panel.Width(common.MaxInt(28, width-48)).Height(height)
	if focused {
		style = style.BorderForeground(ColorInfo())
	}
	return style.Render(body)
}
