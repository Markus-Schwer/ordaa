package tui

import "github.com/charmbracelet/lipgloss"

func WithBorderAndCorner(base lipgloss.Style, cornerChar string, isFocued bool) lipgloss.Style {
	tabBorder := lipgloss.NormalBorder()
	tabBorder.TopLeft = cornerChar
	borderStyle := base.Border(tabBorder)
	if isFocued {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("10"))
	} else {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("8"))
	}
	return borderStyle
}
