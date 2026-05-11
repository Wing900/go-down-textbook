package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chenyb-go/go-down-textbook/internal/app"
	"github.com/chenyb-go/go-down-textbook/internal/models"
	"github.com/chenyb-go/go-down-textbook/internal/tui/view"
)

const version = "v3.0.2"

type page int
type focusArea int
type menuAction int

const (
	pageHome page = iota
	pageLogin
	pageSelect
	pageDownload
)

const (
	focusGrades focusArea = iota
	focusSubjects
	focusBooks
)

const (
	actionStart menuAction = iota
	actionOpenDir
	actionQuit
)

type menuItem struct {
	title  string
	action menuAction
}

type catalogLoadedMsg struct {
	requestID int
	data      *app.CatalogData
	err       error
}

type openDirMsg struct{ err error }
type downloadMsg struct{ update app.DownloadUpdate }

type model struct {
	service *app.Service
	styles  view.Styles
	width   int
	height  int
	page    page
	focus   focusArea

	homeData   app.HomeData
	homeErr    string
	menuIndex  int
	menuItems  []menuItem
	statusLine string

	spinner     spinner.Model
	loginStatus string
	loginDetail string
	catalog     *app.CatalogData
	catalogErr  string
	loadRequest int

	grades, subjects list.Model
	grade, subject   string
	books            []app.SelectableBook
	bookIndex        int
	bookTop          int
	selected         map[string]models.BookItem
	selectHint       string

	downloadCtx    context.Context
	cancelDownload context.CancelFunc
	downloadCh     <-chan app.DownloadUpdate
	downloading    bool
	downloadDone   bool
	totalTarget    int
	doneCount      int
	currentTitle   string
	currentPercent float64
	lastSpeed      string
	completed      []string
	failed         []string
	logs           []string
}

func NewModel(service *app.Service) tea.Model {
	m := &model{
		service: service,
		page:    pageHome,
		spinner: newSpinner(),
		menuItems: []menuItem{
			{title: "开始选书", action: actionStart},
			{title: "打开下载目录", action: actionOpenDir},
			{title: "退出", action: actionQuit},
		},
		selected: make(map[string]models.BookItem),
	}
	m.styles = view.NewStyles()
	m.initLists()
	m.refreshHome()
	return m
}
