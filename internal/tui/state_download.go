package tui

import (
	"context"
	"fmt"

	"github.com/chenyb-go/go-down-textbook/internal/models"
	"github.com/chenyb-go/go-down-textbook/internal/tui/state"
)

func (m *model) startDownload(books []models.BookItem) {
	m.page = pageDownload
	m.downloading = true
	m.downloadDone = false
	m.doneCount = 0
	m.totalTarget = len(books)
	m.currentTitle = "准备开始下载"
	m.currentPercent = 0
	m.lastSpeed = ""
	m.completed = nil
	m.failed = nil
	m.logs = nil
	m.downloadCtx, m.cancelDownload = context.WithCancel(context.Background())
	m.pushLog(fmt.Sprintf("准备下载 %d 本教材", len(books)))
}

func (m *model) pushLog(line string) {
	m.logs = state.PushLog(m.logs, line, 8)
}
