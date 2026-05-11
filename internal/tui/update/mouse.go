package update

import "github.com/chenyb-go/go-down-textbook/internal/tui/common"

const (
	AreaNone = iota
	AreaGrades
	AreaSubjects
	AreaBooks
	AreaFooter
)

type FooterAction string

const (
	FooterNone     FooterAction = ""
	FooterBack     FooterAction = "back"
	FooterStart    FooterAction = "start"
	FooterOpen     FooterAction = "open"
	FooterRetry    FooterAction = "retry"
	FooterContinue FooterAction = "continue"
)

func HomeMenuIndex(y, itemCount int, hasErr bool) int {
	startY := 9
	if hasErr {
		startY++
	}
	index := y - startY
	if index < 0 || index >= itemCount {
		return -1
	}
	return index
}

func SelectArea(x, y, width, listHeight int) int {
	gradesW := common.MaxInt(18, width/5)
	subjectsW := common.MaxInt(16, width/5)
	startY := 7
	if y < startY || y >= startY+listHeight {
		if y >= startY+listHeight+4 {
			return AreaFooter
		}
		return AreaNone
	}
	if x < 2+gradesW {
		return AreaGrades
	}
	if x < 2+gradesW+subjectsW {
		return AreaSubjects
	}
	return AreaBooks
}

func BookIndexByClick(x, y, width, visibleRows, bookTop, bookCount int) int {
	bookX := 2 + common.MaxInt(18, width/5) + common.MaxInt(16, width/5)
	bookY := 8
	if x < bookX || y < bookY || bookCount == 0 {
		return -1
	}
	row := y - bookY - 1
	if row < 0 || row >= visibleRows {
		return -1
	}
	index := bookTop + row
	if index < 0 || index >= bookCount {
		return -1
	}
	return index
}

func SelectFooterAction(x int) FooterAction {
	switch {
	case x >= 34 && x <= 46:
		return FooterStart
	case x >= 49 && x <= 56:
		return FooterBack
	default:
		return FooterNone
	}
}

func DownloadFooterAction(x int, done bool, hasFailed bool) FooterAction {
	switch {
	case x >= 2 && x <= 12:
		return FooterOpen
	case done && x >= 15 && x <= 29:
		return FooterContinue
	case hasFailed && x >= 15 && x <= 26:
		return FooterRetry
	default:
		return FooterNone
	}
}
