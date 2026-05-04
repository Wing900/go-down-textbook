package download

// EventType 下载事件类型
type EventType int

const (
	EventStart    EventType = iota // 开始下载
	EventProgress                  // 下载进度更新
	EventDone                      // 下载完成
	EventError                     // 下载失败
)

// DownloadEvent 下载事件
type DownloadEvent struct {
	Type       EventType
	BookTitle  string
	BookID     string
	Filename   string
	BytesRead  int64
	TotalBytes int64
	Error      error
}

// ProgressInfo 下载进度信息
type ProgressInfo struct {
	BookTitle  string
	BytesRead  int64
	TotalBytes int64
	Percent    float64
	Speed      float64 // bytes per second
	ETA        int     // seconds remaining
}
