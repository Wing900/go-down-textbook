package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderBool(success, accent lipgloss.Style, ok bool) string {
	if ok {
		return success.Render("已登录，可以下载")
	}
	return accent.Render("还未登录")
}

func RenderBar(percent float64, width int, colorize func(...string) string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}
	filled := int(percent * float64(width))
	if filled > width {
		filled = width
	}
	return "[" + colorize(strings.Repeat("=", filled)) + strings.Repeat("-", width-filled) + "]"
}

func RenderList(items []string, empty string) string {
	if len(items) == 0 {
		return empty
	}
	return "  - " + strings.Join(items, "\n  - ")
}

func FallbackText(v string) string {
	if v == "" {
		return "未选择"
	}
	return v
}

func FormatPercent(v float64) string {
	return fmt.Sprintf("当前进度: %.1f%%", v*100)
}
