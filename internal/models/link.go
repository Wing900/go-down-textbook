package models

// LinkData 表示一个教材的下载链接信息
type LinkData struct {
	ContentID string // 资源 ID
	Title     string // 教材名称
	URL       string // 主下载链接
	BackupURL string // 备用下载链接
	Size      int64  // 文件大小（字节）
}
