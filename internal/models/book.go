package models

import "fmt"

// BookItem 表示一本可下载的教材
type BookItem struct {
	ID       string // content_id
	Title    string // 教材名称
	Grade    string // 年级
	Subject  string // 学科
	Filename string // 下载后的文件名
	Size     int64  // 文件大小（字节）
}

// BookOption 用于菜单展示的选择项
type BookOption struct {
	Index int
	Book  BookItem
	// Selected 标记是否被选中（多选用）
	Selected bool
}

// Label 返回菜单展示文本
func (b BookOption) Label() string {
	mark := "  "
	if b.Selected {
		mark = "✓ "
	}
	sizeMB := float64(b.Book.Size) / 1024 / 1024
	if sizeMB > 0 {
		return mark + b.Book.Title + "  (" + formatSize(b.Book.Size) + ")"
	}
	return mark + b.Book.Title
}

// formatSize 格式化文件大小
func formatSize(bytes int64) string {
	if bytes < 1024 {
		return "1KB"
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1fMB", float64(bytes)/1024/1024)
}
