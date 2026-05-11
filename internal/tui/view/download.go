package view

import (
	"strconv"
	"strings"

	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
)

type DownloadData struct {
	Width, DoneCount, TotalTarget int
	CurrentTitle, Speed, Spinner  string
	CurrentPercent                float64
	Downloading, Done             bool
	Completed, Failed, Logs       []string
}

func RenderDownload(s Styles, d DownloadData) string {
	left := s.Panel.Width(common.MaxInt(42, d.Width-6)).Render(strings.Join(summary(s, d), "\n"))
	right := s.Panel.Width(common.MaxInt(42, d.Width-6)).Render(strings.Join(sidebar(s, d), "\n"))
	footer := s.Footer.Width(common.MaxInt(42, d.Width-6)).Render(footerText(d))
	return s.Doc.Render(left + "\n" + right + "\n" + footer)
}

func summary(s Styles, d DownloadData) []string {
	total := 0.0
	if d.TotalTarget > 0 {
		total = float64(d.DoneCount) / float64(d.TotalTarget)
	}
	lines := []string{
		s.Title.Render(downloadTitle(d.Done)),
		"总进度: " + strconv.Itoa(d.DoneCount) + " / " + strconv.Itoa(d.TotalTarget),
		common.RenderBar(total, common.MaxInt(20, d.Width/2), s.Success.Render),
		"", d.Spinner + " " + d.CurrentTitle, common.FormatPercent(d.CurrentPercent),
		common.RenderBar(d.CurrentPercent, common.MaxInt(20, d.Width/2), s.Accent.Render),
	}
	if d.Downloading {
		lines = append(lines, "", s.Muted.Render("下载中，请不要关闭窗口。"))
	}
	if d.Speed != "" {
		lines = append(lines, s.Muted.Render("传输状态: "+d.Speed))
	}
	return lines
}

func sidebar(s Styles, d DownloadData) []string {
	return []string{
		s.Subtitle.Render("已完成"), common.RenderList(d.Completed, "  - 无"), "",
		s.Subtitle.Render("失败"), common.RenderList(d.Failed, "  - 无"), "",
		s.Subtitle.Render("日志"), common.RenderList(d.Logs, "  - 暂无"),
	}
}
