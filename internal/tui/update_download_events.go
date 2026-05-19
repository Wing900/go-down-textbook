package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chenyb-go/go-down-textbook/internal/app"
	"github.com/chenyb-go/go-down-textbook/internal/download"
	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
)

func (m *model) handleDownloadUpdate(update app.DownloadUpdate) (tea.Model, tea.Cmd) {
	switch update.Type {
	case app.UpdateDownload:
		m.applyDownloadEvent(update.Event)
	case app.UpdateBookmarkSuccess:
		m.pushLog("已添加目录书签: " + update.BookTitle)
	case app.UpdateBookmarkError:
		m.pushLog("添加目录书签失败: " + update.BookTitle + " - " + update.Error.Error())
	case app.UpdateFinished:
		m.downloading = false
		m.downloadDone = true
		m.currentPercent = 1
		if update.Error != nil {
			m.currentTitle = "下载流程提前结束"
			m.pushLog("下载流程失败: " + update.Error.Error())
		} else {
			m.currentTitle = "本次下载已完成"
		}
		m.refreshHome()
		return m, nil
	}
	return m, tea.Batch(m.spinner.Tick, waitDownloadCmd(m.downloadCh))
}

func (m *model) applyDownloadEvent(event download.DownloadEvent) {
	switch event.Type {
	case download.EventStart:
		m.currentTitle = event.BookTitle
		m.currentPercent = 0
		m.pushLog("开始下载: " + event.BookTitle)
	case download.EventProgress:
		m.currentTitle = event.BookTitle
		if event.TotalBytes > 0 {
			m.currentPercent = float64(event.BytesRead) / float64(event.TotalBytes)
		}
		m.lastSpeed = common.FormatBytes(event.BytesRead) + " / " + common.FormatBytes(common.Max64(event.TotalBytes, 1))
	case download.EventDone:
		m.doneCount++
		m.currentPercent = 1
		m.completed = append(m.completed, event.BookTitle)
		m.pushLog("下载完成: " + event.BookTitle)
	case download.EventError:
		m.doneCount++
		m.failed = append(m.failed, event.BookTitle)
		m.pushLog("下载失败: " + event.BookTitle + " - " + event.Error.Error())
	}
}
