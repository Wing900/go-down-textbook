package common

import "fmt"

func FormatBytes(v int64) string {
	switch {
	case v >= 1024*1024:
		return fmt.Sprintf("%.1fMB", float64(v)/1024/1024)
	case v >= 1024:
		return fmt.Sprintf("%.1fKB", float64(v)/1024)
	default:
		return fmt.Sprintf("%dB", v)
	}
}

func Max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TruncateText(s string, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width <= 1 {
		return string(r[:width])
	}
	return string(r[:width-1]) + "…"
}
