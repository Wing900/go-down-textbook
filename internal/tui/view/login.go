package view

import (
	"strings"

	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
)

type LoginData struct {
	Status, Detail, Spinner string
	Width                   int
}

func RenderLogin(s Styles, d LoginData) string {
	body := []string{
		s.Subtitle.Render("登录"),
		"",
		"正在打开浏览器，请在页面中完成登录。",
		"登录后会自动返回，不要关闭本窗口。",
		"",
		d.Detail,
		"",
		d.Spinner + " " + d.Status,
		"",
		common.RenderBar(0.35, common.MaxInt(20, d.Width/2), s.Accent.Render),
		"",
		s.Muted.Render("如果长时间没有反应，可以按 Esc 返回后重新尝试。"),
	}
	panel := s.Panel.Width(common.MaxInt(42, d.Width-6)).Render(strings.Join(body, "\n"))
	footer := s.Footer.Width(common.MaxInt(42, d.Width-6)).Render("Esc 返回")
	return s.Doc.Render(panel + "\n" + footer)
}
