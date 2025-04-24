package model

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type header struct {
	homeDir    string
	workingDir string
	width      uint
}

var status = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingLeft(3).
	PaddingRight(3)

var statusExtra = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#403769"))

func (h header) View() string {
	if h.width == 0 {
		return ""
	}

	dir := strings.TrimPrefix(h.workingDir, h.homeDir)
	header := status.Render(dir)
	extender := strings.Repeat(" ", int(h.width-1))

	return header + statusExtra.Render(extender)
}
