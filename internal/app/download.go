package app

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/chenyb-go/go-down-textbook/internal/api"
	"github.com/chenyb-go/go-down-textbook/internal/download"
	"github.com/chenyb-go/go-down-textbook/internal/models"
	"github.com/chenyb-go/go-down-textbook/internal/pdf"
)

func (s *Service) StartDownload(ctx context.Context, token string, books []models.BookItem) <-chan DownloadUpdate {
	client := api.NewClient(token)
	source := download.NewManager(client, 3).DownloadBooks(ctx, books, s.OutputDir)
	updates := make(chan DownloadUpdate, 32)

	go func() {
		defer close(updates)
		var bookmarkWG sync.WaitGroup

		for event := range source {
			updates <- DownloadUpdate{Type: UpdateDownload, Event: event}
			if event.Type != download.EventDone {
				continue
			}
			bookmarkWG.Add(1)
			go func(evt download.DownloadEvent) {
				defer bookmarkWG.Done()
				s.tryAddBookmarks(client, filepath.Join(s.OutputDir, evt.Filename), evt.BookID, evt.BookTitle, updates)
			}(event)
		}

		bookmarkWG.Wait()
		updates <- DownloadUpdate{Type: UpdateFinished}
	}()
	return updates
}

func (s *Service) tryAddBookmarks(client *api.Client, pdfPath, contentID, title string, updates chan<- DownloadUpdate) {
	if !pdf.HasPDFCPU() {
		return
	}

	bookmarks, err := api.FetchBookmarks(client, contentID)
	if err != nil || len(bookmarks) == 0 {
		return
	}
	if err := pdf.AddBookmarks(pdfPath, bookmarks); err != nil {
		updates <- DownloadUpdate{Type: UpdateBookmarkError, BookTitle: title, Error: err}
		return
	}
	updates <- DownloadUpdate{Type: UpdateBookmarkSuccess, BookTitle: title}
}
