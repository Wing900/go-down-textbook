package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chenyb-go/go-down-textbook/internal/tui/common"
	"github.com/chenyb-go/go-down-textbook/internal/tui/view"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.resize()
		return m, nil
	case tea.KeyMsg:
		return m.updateKey(msg)
	case tea.MouseMsg:
		return m.updateMouse(msg)
	case spinner.TickMsg:
		return m.updateSpinner(msg)
	case catalogLoadedMsg:
		return m.updateCatalogLoaded(msg)
	case openDirMsg:
		return m.updateOpenDir(msg)
	case downloadMsg:
		return m.handleDownloadUpdate(msg.update)
	}
	if m.page == pageSelect {
		return m.updateSelectLists(msg)
	}
	return m, nil
}

func (m *model) View() string {
	switch m.page {
	case pageHome:
		return view.RenderHome(m.styles, view.HomeData{
			Version: version, OutputDir: m.homeData.OutputDir, LoggedIn: m.homeData.LoggedIn,
			HistoryCount: m.homeData.HistoryCount, MenuIndex: m.menuIndex, Width: m.width,
			HomeErr: m.homeErr, StatusLine: m.statusLine, MenuItems: m.menuTitles(),
		})
	case pageLogin:
		return view.RenderLogin(m.styles, view.LoginData{
			Status: m.loginStatusOrDefault(), Detail: m.loginDetailOrDefault(),
			Spinner: m.spinner.View(), Width: m.width,
		})
	case pageSelect:
		return view.RenderSelect(m.styles, view.SelectData{
			Width: m.width, Focus: int(m.focus), ListHeight: m.visibleBookRows() + 3, Grade: m.grade, Subject: m.subject,
			Guide: selectGuide(m), GradesView: m.grades.View(), SubjectsView: m.subjects.View(),
			BooksView: view.RenderBooks(m.styles, m.books, m.selected, m.bookIndex, m.bookTop, m.visibleBookRows(), common.MaxInt(28, m.width-48)),
			Queue:     joinSelectedBooks(m.selected, common.MaxInt(20, m.width-12)), SelectedCount: len(m.selected),
		})
	case pageDownload:
		return view.RenderDownload(m.styles, view.DownloadData{
			Width: m.width, DoneCount: m.doneCount, TotalTarget: m.totalTarget, Spinner: m.spinner.View(),
			CurrentTitle: m.currentTitle, CurrentPercent: m.currentPercent, Speed: m.lastSpeed,
			Downloading: m.downloading, Done: m.downloadDone, Completed: m.completed, Failed: m.failed, Logs: m.logs,
		})
	default:
		return ""
	}
}
