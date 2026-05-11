package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chenyb-go/go-down-textbook/internal/app"
)

func loadCatalogCmd(service *app.Service, requestID int) tea.Cmd {
	return func() tea.Msg {
		data, err := service.EnsureCatalog()
		return catalogLoadedMsg{requestID: requestID, data: data, err: err}
	}
}

func openDirCmd(service *app.Service) tea.Cmd {
	return func() tea.Msg {
		return openDirMsg{err: service.OpenDir()}
	}
}

func waitDownloadCmd(ch <-chan app.DownloadUpdate) tea.Cmd {
	if ch == nil {
		return nil
	}
	return func() tea.Msg {
		update, ok := <-ch
		if !ok {
			return nil
		}
		return downloadMsg{update: update}
	}
}
