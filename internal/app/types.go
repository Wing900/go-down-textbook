package app

import (
	"github.com/chenyb-go/go-down-textbook/internal/download"
	"github.com/chenyb-go/go-down-textbook/internal/models"
)

type Service struct {
	OutputDir string
}

type HomeData struct {
	OutputDir    string
	LoggedIn     bool
	HistoryCount int
}

type CatalogData struct {
	Token   string
	Entries []models.CatalogEntry
	Grades  []string
}

type SelectableBook struct {
	models.BookItem
	Downloaded bool
}

type DownloadUpdateType int

const (
	UpdateDownload DownloadUpdateType = iota
	UpdateBookmarkSuccess
	UpdateBookmarkError
	UpdateFinished
)

type DownloadUpdate struct {
	Type      DownloadUpdateType
	Event     download.DownloadEvent
	BookTitle string
	Error     error
}

func NewService(outputDir string) *Service {
	return &Service{OutputDir: outputDir}
}
