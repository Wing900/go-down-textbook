package util

import (
	"fmt"
	"os"
	"strings"
)

// supportsHyperlink 检测终端是否支持 OSC 8 超链接
func supportsHyperlink() bool {
	term := os.Getenv("TERM_PROGRAM")
	// Windows Terminal, VS Code, WezTerm 等都支持
	if term == "vscode" || term == "WezTerm" || term == "Hyper" {
		return true
	}
	// Windows Terminal 通过 WT_SESSION 检测
	if os.Getenv("WT_SESSION") != "" {
		return true
	}
	// ConEmu / Cmder
	if os.Getenv("ConEmuPID") != "" || os.Getenv("CMDER_ROOT") != "" {
		return true
	}
	return false
}

// FileLink 生成可点击的文件路径链接
func FileLink(path, text string) string {
	if !supportsHyperlink() {
		return text
	}
	// 将反斜杠替换为正斜杠，以便在 OSC 8 中使用
	uri := "file:///" + strings.ReplaceAll(path, "\\", "/")
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", uri, text)
}

// HTTPLink 生成可点击的 URL 链接
func HTTPLink(url, text string) string {
	if !supportsHyperlink() {
		return text
	}
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, text)
}

// Header 生成带样式的标题文本
func Header(text string) string {
	return "\033[1;36m" + text + "\033[0m"
}

// Success 生成绿色成功文本
func Success(text string) string {
	return "\033[1;32m" + text + "\033[0m"
}

// Error 生成红色错误文本
func Error(text string) string {
	return "\033[1;31m" + text + "\033[0m"
}

// Warn 生成黄色警告文本
func Warn(text string) string {
	return "\033[1;33m" + text + "\033[0m"
}

// Dim 生成暗淡文本
func Dim(text string) string {
	return "\033[2m" + text + "\033[0m"
}

// Bold 生成粗体文本
func Bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

// Progress 生成进度条文本
func Progress(current, total int, width int) string {
	if total == 0 {
		return ""
	}
	filled := current * width / total
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s] %d/%d", bar, current, total)
}
