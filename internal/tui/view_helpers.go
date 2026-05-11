package tui

func selectGuide(m *model) string {
	if m.selectHint != "" {
		return m.selectHint
	}
	return "选择你想下载的教材，按 空格 勾选，按 D 开始下载。"
}

func (m *model) menuTitles() []string {
	items := make([]string, 0, len(m.menuItems))
	for _, item := range m.menuItems {
		items = append(items, item.title)
	}
	return items
}

func (m *model) loginStatusOrDefault() string {
	if m.loginStatus == "" {
		return "正在处理..."
	}
	return m.loginStatus
}

func (m *model) loginDetailOrDefault() string {
	if m.catalogErr != "" {
		return m.catalogErr
	}
	if m.loginDetail == "" {
		return "这一步只需要做一次，之后就可以直接下载。"
	}
	return m.loginDetail
}
