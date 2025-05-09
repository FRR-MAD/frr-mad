package rib

const inputHeight = 3

var currentSubTabLocal = -1

func (m *Model) RibView(currentSubTab int) string {
	currentSubTabLocal = currentSubTab
	return m.View()
}

func (m *Model) View() string {
	if currentSubTabLocal == 0 {
		return m.renderFibTab()
	} else if currentSubTabLocal == 1 {
		return m.renderOSPFRoutesTab()
	} else if currentSubTabLocal == 2 {
		return m.renderConnectedRoutesTab()
	}
	return m.renderFibTab()
}

func (m *Model) renderFibTab() string {
	return "FIB"
}

func (m *Model) renderOSPFRoutesTab() string {
	return "OSPF learned routes"
}

func (m *Model) renderConnectedRoutesTab() string {

	return "Directly Connected Networks"
}
