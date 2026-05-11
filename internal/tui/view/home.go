package view

import (
	"fmt"
	"strings"

	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
)

type HomeData struct {
	Version, OutputDir, StatusLine string
	LoggedIn                       bool
	HistoryCount, MenuIndex, Width int
	HomeErr                        string
	MenuItems                      []string
}

func RenderHome(s Styles, d HomeData) string {
	header := s.Title.Render("BoooookDown") + "  " + s.Muted.Render(d.Version)
	info := []string{
		fmt.Sprintf("%s  %s", s.Label.Render("保存到"), s.Value.Render(d.OutputDir)),
		fmt.Sprintf("%s  %s", s.Label.Render("登录状态"), common.RenderBool(s.Success, s.Accent, d.LoggedIn)),
		fmt.Sprintf("%s  %s", s.Label.Render("下载记录"), s.Value.Render(fmt.Sprintf("已经下载了 %d 本教材", d.HistoryCount))),
	}
	if d.HomeErr != "" {
		info = append(info, s.Danger.Render("读取状态失败: "+d.HomeErr))
	}
	body := s.Panel.Width(common.MaxInt(42, d.Width-6)).Render(strings.Join([]string{
		header, "BoooookDown 教材下载工具", "", strings.Join(info, "\n"), "",
		renderMenu(s, d.MenuItems, d.MenuIndex), "", s.Info.Render("进入后会动态获取目录，并按年级、学科和教材滚动选择。"),
	}, "\n"))
	return s.Doc.Render(body + "\n" + renderFooter(s, d.StatusLine, d.Width))
}

func renderMenu(s Styles, items []string, selected int) string {
	lines := make([]string, 0, len(items))
	for i, item := range items {
		line := "  " + item
		if i == selected {
			line = s.SelectedLine.Render("> " + item)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func renderFooter(s Styles, status string, width int) string {
	text := "↑↓ 移动   Enter 确认   q 退出"
	if status != "" {
		text = status + "\n" + text
	}
	return s.Footer.Width(common.MaxInt(42, width-6)).Render(text)
}
