package tui

func (m *model) refreshHome() {
	data, err := m.service.HomeData()
	if err != nil {
		m.homeErr = err.Error()
		return
	}
	m.homeErr = ""
	m.homeData = data
}
